package controller

import (
	"bless-activity/model"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/FishPiOffical/golang-sdk/sdk"
	"github.com/FishPiOffical/golang-sdk/types"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type RewardDistributionController struct {
	*BaseController
	fishpiSdk *sdk.FishPiSDK
	event     *core.ServeEvent
}

func NewRewardDistributionController(event *core.ServeEvent, base *BaseController) *RewardDistributionController {
	controller := &RewardDistributionController{
		BaseController: base,
		event:          event,
	}

	controller.registerRoutes()
	return controller
}

func (c *RewardDistributionController) registerRoutes() {
	// todo 发放完成后没有更改活动表中奖励发放状态 需要检查
	rewardGroup := c.event.Router.Group("/activity-api/reward").Bind(
		apis.RequireSuperuserAuth(),
	)
	rewardGroup.POST("/distribute", c.DistributeRewards)
	rewardGroup.POST("/retry", c.RetryFailedDistributions)
}

// DistributeRequest 发放请求参数
type DistributeRequest struct {
	ActivityId string `json:"activityId"`
}

// UserRewardDistribution 用户奖励发放信息（内部使用）
type UserRewardDistribution struct {
	UserId string
	Rank   int
	Point  int
}

// DistributeRewards 发放奖励接口
func (c *RewardDistributionController) DistributeRewards(event *core.RequestEvent) error {
	req := new(DistributeRequest)
	if err := event.BindBody(req); err != nil {
		return event.BadRequestError("Invalid request body", err)
	}

	// 验证活动是否存在
	activity := model.NewActivity(nil)
	if err := c.app.RecordQuery(model.DbNameActivities).
		AndWhere(dbx.HashExp{model.CommonFieldId: req.ActivityId}).
		One(activity); err != nil {
		return event.NotFoundError("Activity not found", err)
	}

	logger := c.app.Logger().With(
		slog.String("controller", "RewardDistribution"),
		slog.String("action", "DistributeRewards"),
		slog.String("activityId", req.ActivityId),
	)

	// 获取活动关联的投票
	voteId := activity.GetVoteId()
	if voteId == "" {
		return event.BadRequestError("Activity has no vote associated", nil)
	}

	vote := model.NewVote(nil)
	if err := c.app.RecordQuery(model.DbNameVotes).
		AndWhere(dbx.HashExp{model.CommonFieldId: voteId}).
		One(vote); err != nil {
		return event.NotFoundError("Vote not found", err)
	}

	// 获取活动关联的奖励组
	rewardGroupId := activity.GetRewardGroupId()
	if rewardGroupId == "" {
		return event.BadRequestError("Activity has no reward group associated", nil)
	}

	logger = logger.With(slog.String("voteId", voteId), slog.String("rewardGroupId", rewardGroupId))

	// 从VoteLog中统计得票数 (只统计有效票)
	records, err := c.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId} && valid = {:valid}",
		"",
		0,
		0,
		map[string]any{
			"voteId": voteId,
			"valid":  model.VoteLogValidValid,
		},
	)
	if err != nil {
		logger.Error("Failed to fetch vote logs", slog.Any("error", err))
		return event.InternalServerError("Failed to fetch vote logs", err)
	}

	// 统计每个用户获得的有效票数和最后一张票的时间
	type voteInfo struct {
		count        int
		lastVoteTime time.Time
	}
	voteStats := make(map[string]*voteInfo)
	for _, record := range records {
		voteLog := model.NewVoteLog(record)
		toUserId := voteLog.ToUserId()
		created := voteLog.Created().Time()

		if info, exists := voteStats[toUserId]; exists {
			info.count++
			// 更新最后一张票的时间（取最晚的时间）
			if created.After(info.lastVoteTime) {
				info.lastVoteTime = created
			}
		} else {
			voteStats[toUserId] = &voteInfo{
				count:        1,
				lastVoteTime: created,
			}
		}
	}

	if len(voteStats) == 0 {
		return event.BadRequestError("No votes found for this vote", nil)
	}

	// 从数据库获取奖励配置
	var rewardRecords []*core.Record
	err = c.app.RecordQuery(model.DbNameRewards).
		AndWhere(dbx.HashExp{model.RewardsFieldRewardGroupId: rewardGroupId}).
		OrderBy("min ASC").
		All(&rewardRecords)

	if err != nil {
		logger.Error("Failed to fetch reward config", slog.Any("error", err))
		return event.InternalServerError("Failed to fetch reward config", err)
	}

	if len(rewardRecords) == 0 {
		return event.BadRequestError("No reward configuration found for this vote", nil)
	}

	// 按得票数排序，票数相同时按最后一张票的时间排序
	type userVote struct {
		userId       string
		votes        int
		lastVoteTime time.Time
	}
	var userVotes []userVote
	for userId, info := range voteStats {
		userVotes = append(userVotes, userVote{
			userId:       userId,
			votes:        info.count,
			lastVoteTime: info.lastVoteTime,
		})
	}

	// 排序：得票数从高到低，票数相同时按最后一张票的时间从早到晚
	for i := 0; i < len(userVotes); i++ {
		for j := i + 1; j < len(userVotes); j++ {
			// 如果票数不同，按票数降序
			if userVotes[j].votes > userVotes[i].votes {
				userVotes[i], userVotes[j] = userVotes[j], userVotes[i]
			} else if userVotes[j].votes == userVotes[i].votes {
				// 票数相同时，按最后一张票的时间升序（越早越靠前）
				if userVotes[j].lastVoteTime.Before(userVotes[i].lastVoteTime) {
					userVotes[i], userVotes[j] = userVotes[j], userVotes[i]
				}
			}
		}
	}

	// 查找参与奖配置（min > 0 且 max = 0）
	var participationReward *model.Reward
	for _, rec := range rewardRecords {
		reward := model.NewReward(rec)
		logger.Debug("Checking reward config",
			slog.Int("min", reward.Min()),
			slog.Int("max", reward.Max()),
			slog.Int("point", reward.Point()))

		if reward.Min() > 0 && reward.Max() == 0 && reward.Point() > 0 {
			participationReward = reward
			logger.Info("Found participation reward config",
				slog.Int("min", reward.Min()),
				slog.Int("max", reward.Max()),
				slog.Int("point", reward.Point()))
			break
		}
	}

	// 构建用户奖励列表，根据排名范围匹配奖励
	var usersToReward []UserRewardDistribution
	rankedUserIds := make(map[string]bool) // 记录已获得名次奖励的用户

	for rank, uv := range userVotes {
		rankNum := rank + 1 // 排名从1开始
		matched := false

		// 查找匹配的名次奖励配置（min <= rankNum <= max，且 max > 0）
		for _, rec := range rewardRecords {
			reward := model.NewReward(rec)
			if reward.Point() == 0 || reward.Max() == 0 {
				continue // 跳过无奖励配置和参与奖配置
			}
			if rankNum >= reward.Min() && rankNum <= reward.Max() {
				usersToReward = append(usersToReward, UserRewardDistribution{
					UserId: uv.userId,
					Rank:   rankNum,
					Point:  reward.Point(),
				})
				rankedUserIds[uv.userId] = true
				matched = true
				break
			}
		}

		// 如果没有匹配到名次奖励，且存在参与奖配置，则标记为参与奖候选
		if !matched && participationReward != nil {
			// 暂不添加，稍后统一处理参与奖
		}
	}

	// 发放参与奖：从articles表获取所有参与活动的用户
	if participationReward != nil {
		// 获取所有提交文章的用户（参与活动的用户）
		articleRecords, err := c.app.FindRecordsByFilter(
			model.DbNameArticles,
			"activityId = {:activityId}",
			"",
			0,
			0,
			map[string]any{
				"activityId": req.ActivityId,
			},
		)
		if err != nil {
			logger.Error("Failed to fetch articles for participation reward", slog.Any("error", err))
		} else {
			logger.Info("Found articles for participation reward",
				slog.Int("articleCount", len(articleRecords)))

			// 遍历所有提交文章的用户
			for _, articleRec := range articleRecords {
				article := model.NewArticle(articleRec)
				userId := article.UserId()

				// 只给未获得名次奖励的用户发放参与奖
				if !rankedUserIds[userId] {
					usersToReward = append(usersToReward, UserRewardDistribution{
						UserId: userId,
						Rank:   0, // 参与奖名次为0
						Point:  participationReward.Point(),
					})
					logger.Info("Added user for participation reward",
						slog.String("userId", userId),
						slog.Int("point", participationReward.Point()))
				} else {
					logger.Debug("User already has ranked reward, skipping participation reward",
						slog.String("userId", userId))
				}
			}
		}
	}

	if len(usersToReward) == 0 {
		return event.BadRequestError("No users eligible for rewards based on config", nil)
	}

	logger.Info("Starting reward distribution",
		slog.Int("totalVoters", len(voteStats)),
		slog.Int("rewardRecipients", len(usersToReward)))

	// 更新活动状态为发放中
	activity.SetRewardDistributionStatus(model.DistributionStatusDistributing)
	if err = c.app.Save(activity); err != nil {
		logger.Error("Failed to update activity status", slog.Any("error", err))
		return event.InternalServerError("Failed to update activity status", err)
	}

	// 开始发放流程
	successCount := 0
	failedCount := 0

	for _, user := range usersToReward {
		if err = c.distributeToUser(voteId, user, logger); err != nil {
			logger.Error("Failed to distribute to user",
				slog.String("userId", user.UserId),
				slog.Any("error", err))
			failedCount++
		} else {
			successCount++
		}
		time.Sleep(time.Millisecond * 200) // 避免请求过快被摸鱼派拒绝
	}

	// 更新活动的最终状态
	if failedCount == 0 {
		activity.SetRewardDistributionStatus(model.DistributionStatusSuccess)
	} else if successCount == 0 {
		activity.SetRewardDistributionStatus(model.DistributionStatusFailed)
	} else {
		// 部分成功部分失败,保持发放中状态
		activity.SetRewardDistributionStatus(model.DistributionStatusDistributing)
	}

	if err = c.app.Save(activity); err != nil {
		logger.Error("Failed to update final activity status", slog.Any("error", err))
	}

	logger.Info("Distribution completed",
		slog.Int("success", successCount),
		slog.Int("failed", failedCount))

	return event.JSON(http.StatusOK, map[string]any{
		"success":        successCount,
		"failed":         failedCount,
		"totalUsers":     len(usersToReward),
		"activityStatus": activity.GetRewardDistributionStatus(),
	})
}

// distributeToUser 为单个用户发放奖励(幂等性处理)
func (c *RewardDistributionController) distributeToUser(voteId string, userReward UserRewardDistribution, logger *slog.Logger) error {
	// 检查是否已经成功发放过
	existingRecord := model.NewRewardDistribution(nil)
	err := c.app.RecordQuery(model.DbNameRewardDistributions).
		AndWhere(dbx.HashExp{
			model.RewardDistributionsFieldVoteId: voteId,
			model.RewardDistributionsFieldUserId: userReward.UserId,
		}).
		One(existingRecord)

	if err == nil {
		// 记录已存在
		if existingRecord.Status() == model.DistributionStatusSuccess {
			logger.Info("Reward already distributed successfully",
				slog.String("userId", userReward.UserId),
				slog.String("recordId", existingRecord.Id))
			return nil // 已经成功发放,跳过
		}
		// 如果是失败或待发放状态,继续尝试发放
	}

	// 创建或更新发放记录
	var record *model.RewardDistribution
	if err != nil {
		// 记录不存在,创建新记录
		collection, err := c.app.FindCollectionByNameOrId(model.DbNameRewardDistributions)
		if err != nil {
			return fmt.Errorf("collection not found: %w", err)
		}
		record = model.NewRewardDistributionFromCollection(collection)
		record.SetVoteId(voteId)
		record.SetUserId(userReward.UserId)
		record.SetRank(userReward.Rank)
		record.SetPoint(userReward.Point)
	} else {
		// 记录存在,使用已有记录
		record = existingRecord
	}

	// 设置为发放中状态
	record.SetStatus(model.DistributionStatusDistributing)
	if err = c.app.Save(record); err != nil {
		return fmt.Errorf("failed to save distributing status: %w", err)
	}

	// 获取用户信息
	user := model.NewUser(nil)
	if err = c.app.RecordQuery(model.DbNameUsers).
		AndWhere(dbx.HashExp{model.CommonFieldId: userReward.UserId}).
		One(user); err != nil {
		record.SetStatus(model.DistributionStatusFailed)
		record.SetMemo(fmt.Sprintf("User not found: %v", err))
		if err1 := c.app.Save(record); err1 != nil {
			logger.Error("Failed to save failed status", slog.Any("error", err1))
		}
		return fmt.Errorf("user not found: %w", err)
	}

	// 获取投票信息用于memo
	vote := model.NewVote(nil)
	if err = c.app.RecordQuery(model.DbNameVotes).
		AndWhere(dbx.HashExp{model.CommonFieldId: voteId}).
		One(vote); err != nil {
		record.SetStatus(model.DistributionStatusFailed)
		record.SetMemo(fmt.Sprintf("Vote not found: %v", err))
		if err1 := c.app.Save(record); err1 != nil {
			logger.Error("Failed to save failed status", slog.Any("error", err1))
		}
		return fmt.Errorf("vote not found: %w", err)
	}

	// 构建memo：您在活动《{votes.name}》中取得第x名 交易单号：{RewardDistributions.id}
	// 参与奖（rank=0）显示"感谢参与"
	var memo string
	if userReward.Rank == 0 {
		memo = fmt.Sprintf("感谢参与活动《%s》 交易单号：%s", vote.Name(), record.Id)
	} else {
		memo = fmt.Sprintf("您在活动《%s》中取得第%d名 交易单号：%s", vote.Name(), userReward.Rank, record.Id)
	}

	// 调用摸鱼派接口发放积分
	if !c.app.IsDev() {
		var resp *types.SimpleResponse
		resp, err = c.fishpiSdk.PostUserEditPoints(user.Name(), userReward.Point, memo)
		if err == nil && resp.Code != 0 {
			err = errors.New(resp.Msg)
		}
	}

	if err != nil {
		// 发放失败
		record.SetStatus(model.DistributionStatusFailed)
		record.SetMemo(fmt.Sprintf("Distribution failed: %v", err))
		if err1 := c.app.Save(record); err1 != nil {
			logger.Error("Failed to save failed status", slog.Any("error", err1))
		}
		return fmt.Errorf("fishpi distribute failed: %w", err)
	}

	// 发放成功
	record.SetStatus(model.DistributionStatusSuccess)
	record.SetMemo(memo)
	if err = c.app.Save(record); err != nil {
		logger.Error("Failed to save success status", slog.Any("error", err))
		return fmt.Errorf("failed to save success status: %w", err)
	}

	logger.Info("Successfully distributed reward",
		slog.String("userId", userReward.UserId),
		slog.String("username", user.Name()),
		slog.Int("point", userReward.Point),
		slog.Int("rank", userReward.Rank))

	return nil
}

// RetryFailedDistributions 重试失败的发放记录
func (c *RewardDistributionController) RetryFailedDistributions(event *core.RequestEvent) error {
	activityId := event.Request.URL.Query().Get("activityId")
	if activityId == "" {
		return event.BadRequestError("activityId is required", nil)
	}

	logger := c.app.Logger().With(
		slog.String("controller", "RewardDistribution"),
		slog.String("action", "RetryFailedDistributions"),
		slog.String("activityId", activityId),
	)

	// 获取活动关联的投票
	activity := model.NewActivity(nil)
	if err := c.app.RecordQuery(model.DbNameActivities).
		AndWhere(dbx.HashExp{model.CommonFieldId: activityId}).
		One(activity); err != nil {
		return event.NotFoundError("Activity not found", err)
	}

	voteId := activity.GetVoteId()
	if voteId == "" {
		return event.BadRequestError("Activity has no vote associated", nil)
	}

	logger = logger.With(slog.String("voteId", voteId))

	// 查找失败的发放记录
	var failedRecords []*core.Record
	err := c.app.RecordQuery(model.DbNameRewardDistributions).
		AndWhere(dbx.HashExp{
			model.RewardDistributionsFieldVoteId: voteId,
			model.RewardDistributionsFieldStatus: model.DistributionStatusFailed,
		}).
		All(&failedRecords)

	if err != nil {
		return event.InternalServerError("Failed to query failed records", err)
	}

	if len(failedRecords) == 0 {
		return event.JSON(http.StatusOK, map[string]any{
			"message": "No failed distributions to retry",
			"retried": 0,
		})
	}

	successCount := 0
	stillFailedCount := 0

	for _, rec := range failedRecords {
		record := model.NewRewardDistribution(rec)
		userReward := UserRewardDistribution{
			UserId: record.UserId(),
			Rank:   record.Rank(),
			Point:  record.Point(),
		}

		if err = c.distributeToUser(voteId, userReward, logger); err != nil {
			stillFailedCount++
		} else {
			successCount++
		}

		time.Sleep(time.Millisecond * 200) // 避免请求过快被摸鱼派拒绝
	}

	logger.Info("Retry completed",
		slog.Int("success", successCount),
		slog.Int("stillFailed", stillFailedCount))

	return event.JSON(http.StatusOK, map[string]any{
		"totalRetried": len(failedRecords),
		"success":      successCount,
		"stillFailed":  stillFailedCount,
	})
}
