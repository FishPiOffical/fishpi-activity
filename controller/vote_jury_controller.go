package controller

import (
	"bless-activity/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

type VoteJuryController struct {
	*BaseController

	group  *router.RouterGroup[*core.RequestEvent]
	logger *slog.Logger
}

func NewVoteJuryController(base *BaseController, group *router.RouterGroup[*core.RequestEvent]) *VoteJuryController {
	logger := base.app.Logger().WithGroup("controller.vote_jury")

	controller := &VoteJuryController{
		BaseController: base,

		group:  group,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *VoteJuryController) registerRoutes() {
	juryGroup := controller.group.Group("/vote/jury")

	// 管理员接口
	juryGroup.GET("/info/{voteId}", controller.GetJuryInfo)
	juryGroup.POST("/user/create", controller.CreateUser).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/member/add", controller.AddMember).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/member/remove", controller.RemoveMember).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/apply/audit", controller.AuditApply).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/status/switch", controller.SwitchStatus).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/calculate", controller.Calculate).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.GET("/vote-details/{voteId}", controller.GetVoteDetails).BindFunc(controller.RequireAuth, controller.RequireAdminByPath)

	// 用户接口
	juryGroup.POST("/apply", controller.Apply).BindFunc(controller.RequireAuth)
	juryGroup.POST("/vote", controller.Vote).BindFunc(controller.RequireAuth)
	juryGroup.POST("/vote/cancel", controller.CancelVote).BindFunc(controller.RequireAuth)
	juryGroup.GET("/result/{voteId}", controller.GetResult)
	juryGroup.GET("/my-apply/{voteId}", controller.GetMyApply).BindFunc(controller.RequireAuth)
	juryGroup.GET("/candidates/{voteId}", controller.GetCandidates).BindFunc(controller.RequireAuth)
}

// RequireAuth 要求登录
func (controller *VoteJuryController) RequireAuth(event *core.RequestEvent) error {
	if event.Auth == nil {
		return event.UnauthorizedError("未登录", nil)
	}
	return event.Next()
}

// RequireAdmin 要求管理员权限
func (controller *VoteJuryController) RequireAdmin(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")
	if voteId == "" {
		// 从请求体中获取
		data := struct {
			VoteId string `json:"voteId"`
		}{}
		if err := event.BindBody(&data); err == nil {
			voteId = data.VoteId
		}
	}

	if voteId == "" {
		return event.BadRequestError("投票ID不能为空", nil)
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: voteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查是否是管理员
	admins := rule.Admins()
	userId := event.Auth.Id
	if !slices.Contains(admins, userId) {
		return event.ForbiddenError("无管理员权限", nil)
	}

	// 将规则存入上下文
	event.Set("jury_rule", rule)

	return event.Next()
}

// RequireAdminByPath 通过路径参数校验管理员权限
func (controller *VoteJuryController) RequireAdminByPath(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")
	if voteId == "" {
		return event.BadRequestError("投票ID不能为空", nil)
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: voteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查是否是管理员
	admins := rule.Admins()
	userId := event.Auth.Id
	if !slices.Contains(admins, userId) {
		return event.ForbiddenError("无管理员权限", nil)
	}

	// 将规则存入上下文
	event.Set("jury_rule", rule)

	return event.Next()
}

// GetJuryInfo 获取评审团详情
func (controller *VoteJuryController) GetJuryInfo(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")

	// 获取投票信息
	vote := new(model.Vote)
	if err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{model.CommonFieldId: voteId}).
		One(vote); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("投票不存在", nil)
		}
		return event.InternalServerError("获取投票失败", err)
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: voteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查当前用户是否是管理员
	isAdmin := false
	if event.Auth != nil {
		admins := rule.Admins()
		isAdmin = slices.Contains(admins, event.Auth.Id)
	}

	// 获取已通过的评审团成员
	var juryUsers []*model.VoteJuryUser
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		All(&juryUsers); err != nil {
		return event.InternalServerError("获取评审团成员失败", err)
	}

	// 扩展用户信息
	members := make([]map[string]any, 0, len(juryUsers))
	for _, juryUser := range juryUsers {
		user := new(model.User)
		if err := controller.app.RecordQuery(model.DbNameUsers).
			Where(dbx.HashExp{model.CommonFieldId: juryUser.UserId()}).
			One(user); err == nil {
			members = append(members, map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			})
		}
	}

	// 获取申请列表（仅管理员可见）
	var applyLogs []map[string]any
	if isAdmin {
		var logs []*model.VoteJuryApplyLog
		if err := controller.app.RecordQuery(model.DbNameVoteJuryApplyLogs).
			Where(dbx.HashExp{model.VoteJuryApplyLogFieldVoteId: voteId}).
			OrderBy(fmt.Sprintf("%s DESC", model.VoteJuryApplyLogFieldCreated)).
			All(&logs); err == nil {
			for _, log := range logs {
				user := new(model.User)
				var userData map[string]any
				if err := controller.app.RecordQuery(model.DbNameUsers).
					Where(dbx.HashExp{model.CommonFieldId: log.UserId()}).
					One(user); err == nil {
					userData = map[string]any{
						"id":       user.Id,
						"name":     user.Name(),
						"nickname": user.Nickname(),
						"avatar":   user.Avatar(),
					}
				}

				applyLogs = append(applyLogs, map[string]any{
					"id":      log.Id,
					"userId":  log.UserId(),
					"reason":  log.Reason(),
					"status":  log.Status(),
					"adminId": log.AdminId(),
					"user":    userData,
					"created": log.GetDateTime("created").String(),
				})
			}
		}
	}

	// 获取所有轮次的投票结果
	var results []*model.VoteJuryResult
	if err := controller.app.RecordQuery(model.DbNameVoteJuryResults).
		Where(dbx.HashExp{model.VoteJuryResultFieldVoteId: voteId}).
		OrderBy("round").
		All(&results); err != nil {
		controller.logger.Error("获取投票结果失败", slog.Any("err", err))
	}

	// 格式化结果
	roundResults := make([]map[string]any, 0, len(results))
	for _, result := range results {
		// 解析results JSON
		var voteResults map[string]int
		if err := json.Unmarshal([]byte(result.Results()), &voteResults); err != nil {
			controller.logger.Error("解析投票结果失败", slog.Any("err", err))
			continue
		}

		// 统计该轮投票的人数
		var roundVoteLogs []*model.VoteJuryLog
		controller.app.RecordQuery(model.DbNameVoteJuryLogs).
			Where(dbx.HashExp{
				model.VoteJuryLogFieldVoteId: voteId,
				model.VoteJuryLogFieldRound:  result.Round(),
			}).
			All(&roundVoteLogs)

		votedUserIds := make(map[string]bool)
		for _, log := range roundVoteLogs {
			votedUserIds[log.FromUserId()] = true
		}
		votedCount := len(votedUserIds)
		abstainCount := len(juryUsers) - votedCount // 弃票人数

		// 构建弃票用户列表
		abstainUsers := make([]map[string]any, 0)
		for _, juryUser := range juryUsers {
			if !votedUserIds[juryUser.UserId()] {
				user := new(model.User)
				if err := controller.app.RecordQuery(model.DbNameUsers).
					Where(dbx.HashExp{model.CommonFieldId: juryUser.UserId()}).
					One(user); err == nil {
					abstainUsers = append(abstainUsers, map[string]any{
						"id":       user.Id,
						"name":     user.Name(),
						"nickname": user.Nickname(),
						"avatar":   user.Avatar(),
					})
				}
			}
		}

		// 构建投票详情：谁投了谁，投了几票
		// 结构: toUserId -> [{fromUserId, fromUser, times}]
		voteDetailsMap := make(map[string][]map[string]any)
		for _, log := range roundVoteLogs {
			fromUser := new(model.User)
			var fromUserData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: log.FromUserId()}).
				One(fromUser); err == nil {
				fromUserData = map[string]any{
					"id":       fromUser.Id,
					"name":     fromUser.Name(),
					"nickname": fromUser.Nickname(),
					"avatar":   fromUser.Avatar(),
				}
			}
			voteDetailsMap[log.ToUserId()] = append(voteDetailsMap[log.ToUserId()], map[string]any{
				"fromUserId": log.FromUserId(),
				"fromUser":   fromUserData,
				"times":      log.Times(),
			})
		}

		// 扩展用户信息
		resultWithUsers := make([]map[string]any, 0)
		for userId, count := range voteResults {
			user := new(model.User)
			var userData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: userId}).
				One(user); err == nil {
				userData = map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
				}
			}
			resultWithUsers = append(resultWithUsers, map[string]any{
				"userId":      userId,
				"count":       count,
				"user":        userData,
				"voteDetails": voteDetailsMap[userId], // 添加投票详情
			})
		}

		roundResults = append(roundResults, map[string]any{
			"round":        result.Round(),
			"results":      resultWithUsers,
			"continue":     result.Continue(),
			"userIds":      result.UserIds(),
			"votedCount":   votedCount,
			"abstainCount": abstainCount,
			"abstainUsers": abstainUsers,
			"totalMembers": len(juryUsers),
		})
	}

	// 获取投票进度（仅评审中状态时）
	var votingProgress map[string]any
	if rule.Status() == model.VoteJuryRuleStatusVoting {
		currentRound := rule.CurrentRound()
		if currentRound == 0 {
			currentRound = 1
		}

		// 统计已投票人数
		var votedLogs []*model.VoteJuryLog
		if err := controller.app.RecordQuery(model.DbNameVoteJuryLogs).
			Where(dbx.HashExp{
				model.VoteJuryLogFieldVoteId: voteId,
				model.VoteJuryLogFieldRound:  currentRound,
			}).
			All(&votedLogs); err == nil {
			// 统计唯一投票人
			votedUsers := make(map[string]bool)
			for _, log := range votedLogs {
				votedUsers[log.FromUserId()] = true
			}
			votingProgress = map[string]any{
				"voted":   len(votedUsers),
				"total":   len(juryUsers),
				"unvoted": len(juryUsers) - len(votedUsers),
			}
		}
	}

	// 检查是否已有最终获胜者（最后一轮结果的continue为false）
	var finalWinner map[string]any
	isVoteCompleted := false
	if len(results) > 0 {
		lastResult := results[len(results)-1]
		if !lastResult.Continue() && len(lastResult.UserIds()) == 1 {
			isVoteCompleted = true
			winnerId := lastResult.UserIds()[0]
			user := new(model.User)
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: winnerId}).
				One(user); err == nil {
				// 从结果中获取票数
				var voteResults map[string]int
				json.Unmarshal([]byte(lastResult.Results()), &voteResults)

				// 获取获胜者的文章
				var winnerArticles []map[string]any
				activity := new(model.Activity)
				if err := controller.app.RecordQuery(model.DbNameActivities).
					Where(dbx.HashExp{model.ActivitiesFieldVoteId: voteId}).
					One(activity); err == nil {
					var articles []*model.RelArticle
					if err := controller.app.RecordQuery(model.DbNameRelArticles).
						Where(dbx.HashExp{
							model.RelArticlesFieldActivityId: activity.Id,
							model.RelArticlesFieldUserId:     winnerId,
						}).
						All(&articles); err == nil {
						for _, art := range articles {
							winnerArticles = append(winnerArticles, map[string]any{
								"id":           art.Id,
								"oId":          art.OId(),
								"title":        art.Title(),
								"viewCount":    art.ViewCount(),
								"goodCnt":      art.GoodCnt(),
								"commentCount": art.CommentCount(),
								"collectCnt":   art.CollectCnt(),
								"thankCnt":     art.ThankCnt(),
							})
						}
					}
				}

				finalWinner = map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
					"votes":    voteResults[winnerId],
					"articles": winnerArticles,
				}
			}
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"vote": map[string]any{
			"id":    vote.Id,
			"name":  vote.Name(),
			"desc":  vote.Desc(),
			"type":  vote.Type(),
			"times": vote.Times(),
			"start": vote.Start().String(),
			"end":   vote.End().String(),
		},
		"rule": map[string]any{
			"id":            rule.Id,
			"count":         rule.Count(),
			"admins":        rule.Admins(),
			"decisions":     rule.Decisions(),
			"status":        rule.Status(),
			"applyTime":     rule.ApplyTime().String(),
			"publicityTime": rule.PublicityTime().String(),
			"currentRound":  rule.CurrentRound(),
		},
		"members":         members,
		"applyLogs":       applyLogs,
		"results":         roundResults,
		"votingProgress":  votingProgress,
		"isAdmin":         isAdmin,
		"isVoteCompleted": isVoteCompleted,
		"finalWinner":     finalWinner,
	})
}

// CreateUser 通过用户名创建用户
func (controller *VoteJuryController) CreateUser(event *core.RequestEvent) error {
	data := struct {
		VoteId   string `json:"voteId"`
		Username string `json:"username"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.Username == "" {
		return event.BadRequestError("用户名不能为空", nil)
	}

	// 检查用户是否已存在
	existingUser := new(model.User)
	if err := controller.app.RecordQuery(model.DbNameUsers).
		Where(dbx.HashExp{model.UsersFieldName: data.Username}).
		One(existingUser); err == nil {
		return event.JSON(http.StatusOK, map[string]any{
			"message": "用户已存在",
			"user": map[string]any{
				"id":       existingUser.Id,
				"name":     existingUser.Name(),
				"nickname": existingUser.Nickname(),
				"avatar":   existingUser.Avatar(),
			},
		})
	}

	// 从 FishPi 获取用户信息
	fishpiUser, err := controller.fishPiSdk.GetUserInfoByUsername(data.Username)
	if err != nil {
		controller.logger.Error("获取鱼派用户信息失败", slog.String("username", data.Username), slog.Any("err", err))
		return event.BadRequestError("鱼派用户不存在或网络请求失败", err)
	}

	if fishpiUser == nil {
		return event.BadRequestError("鱼派用户不存在", nil)
	}

	// 创建本地用户
	userCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameUsers)
	if err != nil {
		return event.InternalServerError("获取用户集合失败", err)
	}

	user := model.NewUserFromCollection(userCollection)
	user.SetEmail(fmt.Sprintf("%s@fishpi.cn", fishpiUser.OId))
	user.SetEmailVisibility(true)
	user.SetVerified(true)
	user.SetOId(fishpiUser.OId)
	user.SetName(fishpiUser.UserName)
	user.SetNickname(fishpiUser.UserNickname)
	user.SetAvatar(fishpiUser.UserAvatarURL)
	user.SetRandomPassword()

	if err := controller.app.Save(user); err != nil {
		return event.InternalServerError("创建用户失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message": "用户创建成功",
		"user": map[string]any{
			"id":       user.Id,
			"name":     user.Name(),
			"nickname": user.Nickname(),
			"avatar":   user.Avatar(),
		},
	})
}

// AddMember 管理员添加评审团成员
func (controller *VoteJuryController) AddMember(event *core.RequestEvent) error {
	data := struct {
		VoteId   string `json:"voteId"`
		Username string `json:"username"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.Username == "" {
		return event.BadRequestError("用户名不能为空", nil)
	}

	// 查找用户
	user := new(model.User)
	if err := controller.app.RecordQuery(model.DbNameUsers).
		Where(dbx.HashExp{model.UsersFieldName: data.Username}).
		One(user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.BadRequestError("用户不存在，请先创建用户", nil)
		}
		return event.InternalServerError("查询用户失败", err)
	}

	// 检查是否已经是评审团成员
	existingMember := new(model.VoteJuryUser)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: data.VoteId,
			model.VoteJuryUserFieldUserId: user.Id,
		}).
		One(existingMember); err == nil {
		return event.BadRequestError("该用户已经是评审团成员", nil)
	}

	// 获取评审团规则检查席位
	rule := event.Get("jury_rule").(*model.VoteJuryRule)

	// 统计已通过的成员数量
	var approvedMembers []*model.VoteJuryUser
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: data.VoteId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		All(&approvedMembers); err != nil {
		return event.InternalServerError("获取评审团成员失败", err)
	}

	if len(approvedMembers) >= rule.Count() {
		return event.BadRequestError("评审团席位已满", nil)
	}

	// 创建评审团成员记录
	juryUserCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteJuryUsers)
	if err != nil {
		return event.InternalServerError("获取评审团成员集合失败", err)
	}

	juryUser := model.NewVoteJuryUserFromCollection(juryUserCollection)
	juryUser.SetVoteId(data.VoteId)
	juryUser.SetUserId(user.Id)
	juryUser.SetStatus(model.VoteJuryUserStatusApproved)

	if err := controller.app.Save(juryUser); err != nil {
		return event.InternalServerError("添加评审团成员失败", err)
	}

	// 创建申请日志记录
	applyLogCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteJuryApplyLogs)
	if err != nil {
		return event.InternalServerError("获取申请日志集合失败", err)
	}

	applyLog := model.NewVoteJuryApplyLogFromCollection(applyLogCollection)
	applyLog.SetVoteId(data.VoteId)
	applyLog.SetUserId(user.Id)
	applyLog.SetReason("管理员指定")
	applyLog.SetStatus(model.VoteJuryApplyLogStatusApproved)
	applyLog.SetAdminId(event.Auth.Id)

	if err := controller.app.Save(applyLog); err != nil {
		controller.logger.Error("创建申请日志失败", slog.Any("err", err))
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message": "添加评审团成员成功",
		"member": map[string]any{
			"id":       user.Id,
			"name":     user.Name(),
			"nickname": user.Nickname(),
			"avatar":   user.Avatar(),
		},
	})
}

// AuditApply 审核申请
func (controller *VoteJuryController) AuditApply(event *core.RequestEvent) error {
	data := struct {
		VoteId  string `json:"voteId"`
		ApplyId string `json:"applyId"`
		Status  string `json:"status"` // approved / rejected
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.ApplyId == "" {
		return event.BadRequestError("申请ID不能为空", nil)
	}

	if data.Status != "approved" && data.Status != "rejected" {
		return event.BadRequestError("状态值无效", nil)
	}

	// 获取申请记录
	applyLog := new(model.VoteJuryApplyLog)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryApplyLogs).
		Where(dbx.HashExp{model.CommonFieldId: data.ApplyId}).
		One(applyLog); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("申请记录不存在", nil)
		}
		return event.InternalServerError("获取申请记录失败", err)
	}

	// 检查申请状态
	if applyLog.Status() != model.VoteJuryApplyLogStatusPending {
		return event.BadRequestError("该申请已处理", nil)
	}

	// 更新申请状态
	newStatus := model.VoteJuryApplyLogStatusApproved
	if data.Status == "rejected" {
		newStatus = model.VoteJuryApplyLogStatusRejected
	}
	applyLog.SetStatus(newStatus)
	applyLog.SetAdminId(event.Auth.Id)

	if err := controller.app.Save(applyLog); err != nil {
		return event.InternalServerError("更新申请状态失败", err)
	}

	// 如果通过，创建评审团成员记录
	if newStatus == model.VoteJuryApplyLogStatusApproved {
		// 获取评审团规则检查席位
		rule := event.Get("jury_rule").(*model.VoteJuryRule)

		// 统计已通过的成员数量
		var approvedMembers []*model.VoteJuryUser
		if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
			Where(dbx.HashExp{
				model.VoteJuryUserFieldVoteId: data.VoteId,
				model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
			}).
			All(&approvedMembers); err != nil {
			return event.InternalServerError("获取评审团成员失败", err)
		}

		if len(approvedMembers) >= rule.Count() {
			// 席位已满，将申请状态改回待审核
			applyLog.SetStatus(model.VoteJuryApplyLogStatusPending)
			applyLog.SetAdminId("")
			_ = controller.app.Save(applyLog)
			return event.BadRequestError("评审团席位已满", nil)
		}

		// 检查是否已经是成员
		existingMember := new(model.VoteJuryUser)
		if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
			Where(dbx.HashExp{
				model.VoteJuryUserFieldVoteId: data.VoteId,
				model.VoteJuryUserFieldUserId: applyLog.UserId(),
			}).
			One(existingMember); err == nil {
			// 更新状态为通过
			existingMember.SetStatus(model.VoteJuryUserStatusApproved)
			if err := controller.app.Save(existingMember); err != nil {
				return event.InternalServerError("更新评审团成员状态失败", err)
			}
		} else {
			// 创建新成员记录
			juryUserCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteJuryUsers)
			if err != nil {
				return event.InternalServerError("获取评审团成员集合失败", err)
			}

			juryUser := model.NewVoteJuryUserFromCollection(juryUserCollection)
			juryUser.SetVoteId(data.VoteId)
			juryUser.SetUserId(applyLog.UserId())
			juryUser.SetStatus(model.VoteJuryUserStatusApproved)

			if err := controller.app.Save(juryUser); err != nil {
				return event.InternalServerError("创建评审团成员失败", err)
			}
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message": "审核成功",
	})
}

// SwitchStatus 切换状态
func (controller *VoteJuryController) SwitchStatus(event *core.RequestEvent) error {
	data := struct {
		VoteId    string `json:"voteId"`
		NewStatus string `json:"newStatus"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	// 解析新状态
	newStatus, err := model.ParseVoteJuryRuleStatus(data.NewStatus)
	if err != nil {
		return event.BadRequestError("状态值无效", nil)
	}

	rule := event.Get("jury_rule").(*model.VoteJuryRule)
	currentStatus := rule.Status()

	// 状态流转校验（支持前进和回退）
	validTransitions := map[model.VoteJuryRuleStatus][]model.VoteJuryRuleStatus{
		model.VoteJuryRuleStatusPending:   {model.VoteJuryRuleStatusApplying},
		model.VoteJuryRuleStatusApplying:  {model.VoteJuryRuleStatusPublicity, model.VoteJuryRuleStatusPending},   // 可回退到未开启
		model.VoteJuryRuleStatusPublicity: {model.VoteJuryRuleStatusVoting, model.VoteJuryRuleStatusApplying},     // 可回退到申请中
		model.VoteJuryRuleStatusVoting:    {model.VoteJuryRuleStatusCompleted, model.VoteJuryRuleStatusPublicity}, // 可回退到公示中
		model.VoteJuryRuleStatusCompleted: {model.VoteJuryRuleStatusVoting},                               // 可回退到评审中
	}

	validNextStatuses := validTransitions[currentStatus]
	if !slices.Contains(validNextStatuses, newStatus) {
		return event.BadRequestError(fmt.Sprintf("不能从 %s 切换到 %s", currentStatus, newStatus), nil)
	}

	// 进入评审状态时，初始化轮次
	if newStatus == model.VoteJuryRuleStatusVoting && rule.CurrentRound() == 0 {
		rule.SetCurrentRound(1)
	}

	rule.SetStatus(newStatus)

	if err := controller.app.Save(rule); err != nil {
		return event.InternalServerError("更新状态失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message":   "状态切换成功",
		"newStatus": newStatus,
	})
}

// Calculate 手动触发算票
func (controller *VoteJuryController) Calculate(event *core.RequestEvent) error {
	data := struct {
		VoteId string `json:"voteId"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	rule := event.Get("jury_rule").(*model.VoteJuryRule)

	// 检查状态
	if rule.Status() != model.VoteJuryRuleStatusVoting {
		return event.BadRequestError("当前状态不是评审中", nil)
	}

	currentRound := rule.CurrentRound()
	if currentRound == 0 {
		currentRound = 1
	}

	// 检查当前轮次是否已经算过票
	existingResult := new(model.VoteJuryResult)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryResults).
		Where(dbx.HashExp{
			model.VoteJuryResultFieldVoteId: data.VoteId,
			model.VoteJuryResultFieldRound:  currentRound,
		}).
		One(existingResult); err == nil {
		// 已经算过票了
		if !existingResult.Continue() {
			return event.BadRequestError("投票已结束，无需再次算票", nil)
		}
		return event.BadRequestError("当前轮次已算票，请等待下一轮投票完成后再算票", nil)
	}

	// 统计当前轮次的投票
	var voteLogs []*model.VoteJuryLog
	if err := controller.app.RecordQuery(model.DbNameVoteJuryLogs).
		Where(dbx.HashExp{
			model.VoteJuryLogFieldVoteId: data.VoteId,
			model.VoteJuryLogFieldRound:  currentRound,
		}).
		All(&voteLogs); err != nil {
		return event.InternalServerError("获取投票记录失败", err)
	}

	if len(voteLogs) == 0 {
		return event.BadRequestError("当前轮次没有投票记录", nil)
	}

	// 统计每个候选人的得票数
	voteCount := make(map[string]int)
	for _, log := range voteLogs {
		voteCount[log.ToUserId()] += log.Times()
	}

	// 找出最高票数
	maxVotes := 0
	for _, count := range voteCount {
		if count > maxVotes {
			maxVotes = count
		}
	}

	// 找出所有得到最高票的用户
	topUsers := make([]string, 0)
	for userId, count := range voteCount {
		if count == maxVotes {
			topUsers = append(topUsers, userId)
		}
	}

	// 判断是否需要进入下一轮（有平票）
	needNextRound := len(topUsers) > 1

	// 如果有平票，检查决策票能否决定
	var winner string
	if needNextRound {
		decisions := rule.Decisions()
		decisionVotes := make(map[string]int)

		for _, log := range voteLogs {
			if slices.Contains(decisions, log.FromUserId()) && slices.Contains(topUsers, log.ToUserId()) {
				decisionVotes[log.ToUserId()] += log.Times()
			}
		}

		// 找出决策票最高的
		maxDecisionVotes := 0
		for _, count := range decisionVotes {
			if count > maxDecisionVotes {
				maxDecisionVotes = count
			}
		}

		decisionTopUsers := make([]string, 0)
		for userId, count := range decisionVotes {
			if count == maxDecisionVotes {
				decisionTopUsers = append(decisionTopUsers, userId)
			}
		}

		// 如果决策票能决定，则不需要下一轮
		if len(decisionTopUsers) == 1 && maxDecisionVotes > 0 {
			needNextRound = false
			winner = decisionTopUsers[0]
		}
	} else {
		// 没有平票，直接确定获胜者
		winner = topUsers[0]
	}

	// 保存本轮结果
	resultsJson, _ := json.Marshal(voteCount)

	resultCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteJuryResults)
	if err != nil {
		return event.InternalServerError("获取结果集合失败", err)
	}

	result := model.NewVoteJuryResultFromCollection(resultCollection)
	result.SetVoteId(data.VoteId)
	result.SetRound(currentRound)
	result.SetResults(string(resultsJson))
	result.SetContinue(needNextRound)
	if needNextRound {
		result.SetUserIds(topUsers) // 平票的用户进入下一轮
	} else if winner != "" {
		result.SetUserIds([]string{winner}) // 最终获胜者
	}

	if err := controller.app.Save(result); err != nil {
		return event.InternalServerError("保存投票结果失败", err)
	}

	// 如果需要下一轮投票
	if needNextRound {
		rule.SetCurrentRound(currentRound + 1)
		if err := controller.app.Save(rule); err != nil {
			return event.InternalServerError("更新轮次失败", err)
		}

		// 扩展平票用户信息
		tieUsers := make([]map[string]any, 0, len(topUsers))
		for _, userId := range topUsers {
			user := new(model.User)
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: userId}).
				One(user); err == nil {
				tieUsers = append(tieUsers, map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
					"votes":    voteCount[userId],
				})
			}
		}

		return event.JSON(http.StatusOK, map[string]any{
			"message":       fmt.Sprintf("第 %d 轮投票结束，有 %d 人平票（%d票），进入第 %d 轮投票", currentRound, len(topUsers), maxVotes, currentRound+1),
			"needNextRound": true,
			"currentRound":  currentRound,
			"nextRound":     currentRound + 1,
			"tieUsers":      tieUsers,
			"results":       voteCount,
		})
	}

	// 投票结束，更新状态为计票完成
	rule.SetStatus(model.VoteJuryRuleStatusCompleted)
	if err := controller.app.Save(rule); err != nil {
		return event.InternalServerError("更新状态失败", err)
	}

	// 获取获胜者信息
	winnerUser := new(model.User)
	var winnerInfo map[string]any
	winnerNickname := "未知用户"
	winnerVotes := voteCount[winner]

	if err := controller.app.RecordQuery(model.DbNameUsers).
		Where(dbx.HashExp{model.CommonFieldId: winner}).
		One(winnerUser); err == nil {
		winnerNickname = winnerUser.Nickname()
		winnerInfo = map[string]any{
			"id":       winnerUser.Id,
			"name":     winnerUser.Name(),
			"nickname": winnerUser.Nickname(),
			"avatar":   winnerUser.Avatar(),
			"votes":    winnerVotes,
		}
	} else {
		controller.logger.Error("获取获胜者信息失败", slog.String("winnerId", winner), slog.Any("err", err))
		winnerInfo = map[string]any{
			"id":    winner,
			"votes": winnerVotes,
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message":       fmt.Sprintf("投票结束！获胜者: %s (%d票)", winnerNickname, winnerVotes),
		"needNextRound": false,
		"currentRound":  currentRound,
		"winner":        winnerInfo,
		"results":       voteCount,
	})
}

// Apply 用户申请加入评审团
func (controller *VoteJuryController) Apply(event *core.RequestEvent) error {
	data := struct {
		VoteId string `json:"voteId"`
		Reason string `json:"reason"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.VoteId == "" {
		return event.BadRequestError("投票ID不能为空", nil)
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: data.VoteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查状态是否允许申请
	if rule.Status() != model.VoteJuryRuleStatusApplying {
		return event.BadRequestError("当前不在申请阶段", nil)
	}

	userId := event.Auth.Id

	// 检查是否已申请
	existingLog := new(model.VoteJuryApplyLog)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryApplyLogs).
		Where(dbx.HashExp{
			model.VoteJuryApplyLogFieldVoteId: data.VoteId,
			model.VoteJuryApplyLogFieldUserId: userId,
		}).
		AndWhere(dbx.NotIn(model.VoteJuryApplyLogFieldStatus, model.VoteJuryApplyLogStatusRejected)).
		One(existingLog); err == nil {
		return event.BadRequestError("您已经申请过了", nil)
	}

	// 创建申请记录
	applyLogCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteJuryApplyLogs)
	if err != nil {
		return event.InternalServerError("获取申请日志集合失败", err)
	}

	applyLog := model.NewVoteJuryApplyLogFromCollection(applyLogCollection)
	applyLog.SetVoteId(data.VoteId)
	applyLog.SetUserId(userId)
	applyLog.SetReason(data.Reason)
	applyLog.SetStatus(model.VoteJuryApplyLogStatusPending)

	if err := controller.app.Save(applyLog); err != nil {
		return event.InternalServerError("提交申请失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message": "申请提交成功",
	})
}

// Vote 评审团成员投票
func (controller *VoteJuryController) Vote(event *core.RequestEvent) error {
	data := struct {
		VoteId   string `json:"voteId"`
		ToUserId string `json:"toUserId"`
		Comment  string `json:"comment"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.VoteId == "" || data.ToUserId == "" {
		return event.BadRequestError("参数不完整", nil)
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: data.VoteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查状态
	if rule.Status() != model.VoteJuryRuleStatusVoting {
		return event.BadRequestError("当前不在投票阶段", nil)
	}

	userId := event.Auth.Id
	currentRound := rule.CurrentRound()
	if currentRound == 0 {
		currentRound = 1
	}

	// 检查是否是评审团成员
	juryUser := new(model.VoteJuryUser)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: data.VoteId,
			model.VoteJuryUserFieldUserId: userId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		One(juryUser); err != nil {
		return event.ForbiddenError("您不是评审团成员", nil)
	}

	// 获取投票配置
	vote := new(model.Vote)
	if err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{model.CommonFieldId: data.VoteId}).
		One(vote); err != nil {
		return event.InternalServerError("获取投票配置失败", err)
	}

	// 检查本轮是否已投票（根据投票次数配置）
	var existingVotes []*model.VoteJuryLog
	if err := controller.app.RecordQuery(model.DbNameVoteJuryLogs).
		Where(dbx.HashExp{
			model.VoteJuryLogFieldVoteId:     data.VoteId,
			model.VoteJuryLogFieldFromUserId: userId,
			model.VoteJuryLogFieldRound:      currentRound,
		}).
		All(&existingVotes); err != nil {
		return event.InternalServerError("检查投票记录失败", err)
	}

	// 统计已使用的票数
	usedVotes := 0
	for _, v := range existingVotes {
		usedVotes += v.Times()
	}

	if usedVotes >= vote.Times() {
		return event.BadRequestError("您的投票次数已用完", nil)
	}

	// 如果不允许重复投票，检查是否已给该用户投过票
	if !vote.Repeat() {
		for _, v := range existingVotes {
			if v.ToUserId() == data.ToUserId {
				return event.BadRequestError("您已经给该用户投过票了", nil)
			}
		}
	}

	// 如果是第2轮及以后，检查被投用户是否在候选名单中
	if currentRound > 1 {
		// 获取上一轮结果
		lastResult := new(model.VoteJuryResult)
		if err := controller.app.RecordQuery(model.DbNameVoteJuryResults).
			Where(dbx.HashExp{
				model.VoteJuryResultFieldVoteId: data.VoteId,
				model.VoteJuryResultFieldRound:  currentRound - 1,
			}).
			One(lastResult); err != nil {
			return event.InternalServerError("获取上一轮结果失败", err)
		}

		if !slices.Contains(lastResult.UserIds(), data.ToUserId) {
			return event.BadRequestError("该用户不在本轮候选名单中", nil)
		}
	}

	// 创建投票记录
	voteLogCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteJuryLogs)
	if err != nil {
		return event.InternalServerError("获取投票日志集合失败", err)
	}

	voteLog := model.NewVoteJuryLogFromCollection(voteLogCollection)
	voteLog.SetVoteId(data.VoteId)
	voteLog.SetFromUserId(userId)
	voteLog.SetToUserId(data.ToUserId)
	voteLog.SetTimes(1)
	voteLog.SetRound(currentRound)
	voteLog.SetComment(data.Comment)

	if err := controller.app.Save(voteLog); err != nil {
		return event.InternalServerError("保存投票记录失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message":   "投票成功",
		"remaining": vote.Times() - usedVotes - 1,
	})
}

// CancelVote 撤销投票
func (controller *VoteJuryController) CancelVote(event *core.RequestEvent) error {
	data := struct {
		VoteId   string `json:"voteId"`
		ToUserId string `json:"toUserId"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.VoteId == "" {
		return event.BadRequestError("投票ID不能为空", nil)
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: data.VoteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查状态
	if rule.Status() != model.VoteJuryRuleStatusVoting {
		return event.BadRequestError("当前不在投票阶段", nil)
	}

	currentRound := rule.CurrentRound()
	if currentRound == 0 {
		currentRound = 1
	}

	// 检查当前轮次是否已经算过票
	existingResult := new(model.VoteJuryResult)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryResults).
		Where(dbx.HashExp{
			model.VoteJuryResultFieldVoteId: data.VoteId,
			model.VoteJuryResultFieldRound:  currentRound,
		}).
		One(existingResult); err == nil {
		return event.BadRequestError("当前轮次已算票，无法撤销投票", nil)
	}

	userId := event.Auth.Id

	// 检查是否是评审团成员
	juryUser := new(model.VoteJuryUser)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: data.VoteId,
			model.VoteJuryUserFieldUserId: userId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		One(juryUser); err != nil {
		return event.ForbiddenError("您不是评审团成员", nil)
	}

	// 查找要撤销的投票记录
	conditions := dbx.HashExp{
		model.VoteJuryLogFieldVoteId:     data.VoteId,
		model.VoteJuryLogFieldFromUserId: userId,
		model.VoteJuryLogFieldRound:      currentRound,
	}
	if data.ToUserId != "" {
		conditions[model.VoteJuryLogFieldToUserId] = data.ToUserId
	}

	var voteLogs []*model.VoteJuryLog
	if err := controller.app.RecordQuery(model.DbNameVoteJuryLogs).
		Where(conditions).
		All(&voteLogs); err != nil {
		return event.InternalServerError("查询投票记录失败", err)
	}

	if len(voteLogs) == 0 {
		return event.BadRequestError("没有找到可撤销的投票记录", nil)
	}

	// 删除投票记录
	cancelledCount := 0
	for _, log := range voteLogs {
		if err := controller.app.Delete(log); err != nil {
			controller.logger.Error("删除投票记录失败", slog.Any("err", err))
		} else {
			cancelledCount++
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message":        fmt.Sprintf("已撤销 %d 条投票记录", cancelledCount),
		"cancelledCount": cancelledCount,
	})
}

// GetResult 获取投票结果
func (controller *VoteJuryController) GetResult(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")

	// 先获取已通过的评审团成员（用于计算弃票）
	var juryUsers []*model.VoteJuryUser
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		All(&juryUsers); err != nil {
		return event.InternalServerError("获取评审团成员失败", err)
	}

	// 获取所有轮次的投票结果
	var results []*model.VoteJuryResult
	if err := controller.app.RecordQuery(model.DbNameVoteJuryResults).
		Where(dbx.HashExp{model.VoteJuryResultFieldVoteId: voteId}).
		OrderBy("round").
		All(&results); err != nil {
		return event.InternalServerError("获取投票结果失败", err)
	}

	// 格式化结果
	roundResults := make([]map[string]any, 0, len(results))
	for _, result := range results {
		// 解析results JSON
		var voteResults map[string]int
		if err := json.Unmarshal([]byte(result.Results()), &voteResults); err != nil {
			controller.logger.Error("解析投票结果失败", slog.Any("err", err))
			continue
		}

		// 统计该轮投票的人数
		var roundVoteLogs []*model.VoteJuryLog
		controller.app.RecordQuery(model.DbNameVoteJuryLogs).
			Where(dbx.HashExp{
				model.VoteJuryLogFieldVoteId: voteId,
				model.VoteJuryLogFieldRound:  result.Round(),
			}).
			All(&roundVoteLogs)

		votedUserIds := make(map[string]bool)
		for _, log := range roundVoteLogs {
			votedUserIds[log.FromUserId()] = true
		}
		votedCount := len(votedUserIds)
		abstainCount := len(juryUsers) - votedCount

		// 构建弃票用户列表
		abstainUsers := make([]map[string]any, 0)
		for _, juryUser := range juryUsers {
			if !votedUserIds[juryUser.UserId()] {
				user := new(model.User)
				if err := controller.app.RecordQuery(model.DbNameUsers).
					Where(dbx.HashExp{model.CommonFieldId: juryUser.UserId()}).
					One(user); err == nil {
					abstainUsers = append(abstainUsers, map[string]any{
						"id":       user.Id,
						"name":     user.Name(),
						"nickname": user.Nickname(),
						"avatar":   user.Avatar(),
					})
				}
			}
		}

		// 构建投票详情：谁投了谁，投了几票
		// 结构: toUserId -> [{fromUserId, fromUser, times}]
		voteDetailsMap := make(map[string][]map[string]any)
		for _, log := range roundVoteLogs {
			fromUser := new(model.User)
			var fromUserData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: log.FromUserId()}).
				One(fromUser); err == nil {
				fromUserData = map[string]any{
					"id":       fromUser.Id,
					"name":     fromUser.Name(),
					"nickname": fromUser.Nickname(),
					"avatar":   fromUser.Avatar(),
				}
			}
			voteDetailsMap[log.ToUserId()] = append(voteDetailsMap[log.ToUserId()], map[string]any{
				"fromUserId": log.FromUserId(),
				"fromUser":   fromUserData,
				"times":      log.Times(),
			})
		}

		// 扩展用户信息
		resultWithUsers := make([]map[string]any, 0)
		for userId, count := range voteResults {
			user := new(model.User)
			var userData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: userId}).
				One(user); err == nil {
				userData = map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
				}
			}
			resultWithUsers = append(resultWithUsers, map[string]any{
				"userId":      userId,
				"count":       count,
				"user":        userData,
				"voteDetails": voteDetailsMap[userId], // 添加投票详情
			})
		}

		roundResults = append(roundResults, map[string]any{
			"round":        result.Round(),
			"results":      resultWithUsers,
			"continue":     result.Continue(),
			"userIds":      result.UserIds(),
			"votedCount":   votedCount,
			"abstainCount": abstainCount,
			"abstainUsers": abstainUsers,
			"totalMembers": len(juryUsers),
		})
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: voteId}).
		One(rule); err != nil {
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 扩展用户信息
	members := make([]map[string]any, 0, len(juryUsers))
	for _, juryUser := range juryUsers {
		user := new(model.User)
		if err := controller.app.RecordQuery(model.DbNameUsers).
			Where(dbx.HashExp{model.CommonFieldId: juryUser.UserId()}).
			One(user); err == nil {
			members = append(members, map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			})
		}
	}

	// 检查是否已有最终获胜者
	var finalWinner map[string]any
	isVoteCompleted := false
	if len(results) > 0 {
		lastResult := results[len(results)-1]
		if !lastResult.Continue() && len(lastResult.UserIds()) == 1 {
			isVoteCompleted = true
			winnerId := lastResult.UserIds()[0]
			user := new(model.User)
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: winnerId}).
				One(user); err == nil {
				var voteResults map[string]int
				json.Unmarshal([]byte(lastResult.Results()), &voteResults)

				// 获取获胜者的文章
				var winnerArticles []map[string]any
				activity := new(model.Activity)
				if err := controller.app.RecordQuery(model.DbNameActivities).
					Where(dbx.HashExp{model.ActivitiesFieldVoteId: voteId}).
					One(activity); err == nil {
					var articles []*model.RelArticle
					if err := controller.app.RecordQuery(model.DbNameRelArticles).
						Where(dbx.HashExp{
							model.RelArticlesFieldActivityId: activity.Id,
							model.RelArticlesFieldUserId:     winnerId,
						}).
						All(&articles); err == nil {
						for _, art := range articles {
							winnerArticles = append(winnerArticles, map[string]any{
								"id":           art.Id,
								"oId":          art.OId(),
								"title":        art.Title(),
								"viewCount":    art.ViewCount(),
								"goodCnt":      art.GoodCnt(),
								"commentCount": art.CommentCount(),
								"collectCnt":   art.CollectCnt(),
								"thankCnt":     art.ThankCnt(),
							})
						}
					}
				}

				finalWinner = map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
					"votes":    voteResults[winnerId],
					"articles": winnerArticles,
				}
			}
		}
	}

	// 检查当前用户是否是管理员
	isAdmin := false
	if event.Auth != nil {
		admins := rule.Admins()
		isAdmin = slices.Contains(admins, event.Auth.Id)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"status":          rule.Status(),
		"currentRound":    rule.CurrentRound(),
		"results":         roundResults,
		"members":         members,
		"totalMembers":    len(juryUsers),
		"isVoteCompleted": isVoteCompleted,
		"finalWinner":     finalWinner,
		"isAdmin":         isAdmin,
	})
}

// RemoveMember 删除评审团成员
func (controller *VoteJuryController) RemoveMember(event *core.RequestEvent) error {
	data := struct {
		VoteId string `json:"voteId"`
		UserId string `json:"userId"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("参数错误", err)
	}

	if data.UserId == "" {
		return event.BadRequestError("用户ID不能为空", nil)
	}

	// 查找评审团成员记录
	juryUser := new(model.VoteJuryUser)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: data.VoteId,
			model.VoteJuryUserFieldUserId: data.UserId,
		}).
		One(juryUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("该用户不是评审团成员", nil)
		}
		return event.InternalServerError("查询评审团成员失败", err)
	}

	// 更新成员状态为已拒绝
	juryUser.SetStatus(model.VoteJuryUserStatusRejected)
	if err := controller.app.Save(juryUser); err != nil {
		return event.InternalServerError("更新评审团成员状态失败", err)
	}

	// 查找该用户最新的申请记录并更新
	var applyLogs []*model.VoteJuryApplyLog
	if err := controller.app.RecordQuery(model.DbNameVoteJuryApplyLogs).
		Where(dbx.HashExp{
			model.VoteJuryApplyLogFieldVoteId: data.VoteId,
			model.VoteJuryApplyLogFieldUserId: data.UserId,
		}).
		OrderBy(fmt.Sprintf("%s DESC", model.VoteJuryApplyLogFieldCreated)).
		Limit(1).
		All(&applyLogs); err == nil && len(applyLogs) > 0 {
		applyLog := applyLogs[0]
		applyLog.SetStatus(model.VoteJuryApplyLogStatusRejected)
		applyLog.SetReason("管理员删除")
		applyLog.SetAdminId(event.Auth.Id)
		if err := controller.app.Save(applyLog); err != nil {
			controller.logger.Error("更新申请日志失败", slog.Any("err", err))
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message": "删除评审团成员成功",
	})
}

// GetMyApply 获取当前用户的申请记录
func (controller *VoteJuryController) GetMyApply(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")
	userId := event.Auth.Id

	// 查找用户的申请记录
	var applyLogs []*model.VoteJuryApplyLog
	if err := controller.app.RecordQuery(model.DbNameVoteJuryApplyLogs).
		Where(dbx.HashExp{
			model.VoteJuryApplyLogFieldVoteId: voteId,
			model.VoteJuryApplyLogFieldUserId: userId,
		}).
		OrderBy(fmt.Sprintf("%s DESC", model.VoteJuryApplyLogFieldCreated)).
		All(&applyLogs); err != nil {
		return event.InternalServerError("获取申请记录失败", err)
	}

	// 格式化结果
	result := make([]map[string]any, 0, len(applyLogs))
	for _, log := range applyLogs {
		result = append(result, map[string]any{
			"id":      log.Id,
			"reason":  log.Reason(),
			"status":  log.Status(),
			"adminId": log.AdminId(),
			"created": log.GetDateTime("created").String(),
		})
	}

	// 检查是否已是评审团成员
	isMember := false
	juryUser := new(model.VoteJuryUser)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldUserId: userId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		One(juryUser); err == nil {
		isMember = true
	}

	return event.JSON(http.StatusOK, map[string]any{
		"applies":  result,
		"isMember": isMember,
	})
}

// GetCandidates 获取候选人列表（用于投票）
func (controller *VoteJuryController) GetCandidates(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")
	userId := event.Auth.Id

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: voteId}).
		One(rule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.NotFoundError("评审团规则不存在", nil)
		}
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 检查是否是评审团成员
	juryUser := new(model.VoteJuryUser)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldUserId: userId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		One(juryUser); err != nil {
		return event.ForbiddenError("您不是评审团成员", nil)
	}

	// 获取投票配置
	vote := new(model.Vote)
	if err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{model.CommonFieldId: voteId}).
		One(vote); err != nil {
		return event.InternalServerError("获取投票配置失败", err)
	}

	currentRound := rule.CurrentRound()
	if currentRound == 0 {
		currentRound = 1
	}

	// 获取关联的活动
	activity := new(model.Activity)
	if err := controller.app.RecordQuery(model.DbNameActivities).
		Where(dbx.HashExp{model.ActivitiesFieldVoteId: voteId}).
		One(activity); err != nil {
		return event.InternalServerError("获取关联活动失败", err)
	}

	var candidates []map[string]any

	if currentRound == 1 {
		// 第一轮：从活动关联的文章中获取候选人
		var articles []*model.RelArticle
		if err := controller.app.RecordQuery(model.DbNameRelArticles).
			Where(dbx.HashExp{model.RelArticlesFieldActivityId: activity.Id}).
			All(&articles); err != nil {
			return event.InternalServerError("获取文章列表失败", err)
		}

		// 按用户分组
		userArticles := make(map[string][]*model.RelArticle)
		for _, article := range articles {
			userArticles[article.UserId()] = append(userArticles[article.UserId()], article)
		}

		for candidateUserId, arts := range userArticles {
			user := new(model.User)
			var userData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: candidateUserId}).
				One(user); err == nil {
				userData = map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
				}
			}

			// 格式化文章信息
			articleList := make([]map[string]any, 0, len(arts))
			for _, art := range arts {
				articleList = append(articleList, map[string]any{
					"id":           art.Id,
					"oId":          art.OId(),
					"title":        art.Title(),
					"viewCount":    art.ViewCount(),
					"goodCnt":      art.GoodCnt(),
					"commentCount": art.CommentCount(),
					"collectCnt":   art.CollectCnt(),
					"thankCnt":     art.ThankCnt(),
				})
			}

			candidates = append(candidates, map[string]any{
				"userId":   candidateUserId,
				"user":     userData,
				"articles": articleList,
			})
		}
	} else {
		// 后续轮次：从上一轮结果中获取候选人
		lastResult := new(model.VoteJuryResult)
		if err := controller.app.RecordQuery(model.DbNameVoteJuryResults).
			Where(dbx.HashExp{
				model.VoteJuryResultFieldVoteId: voteId,
				model.VoteJuryResultFieldRound:  currentRound - 1,
			}).
			One(lastResult); err != nil {
			return event.InternalServerError("获取上一轮结果失败", err)
		}

		for _, candidateUserId := range lastResult.UserIds() {
			user := new(model.User)
			var userData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: candidateUserId}).
				One(user); err == nil {
				userData = map[string]any{
					"id":       user.Id,
					"name":     user.Name(),
					"nickname": user.Nickname(),
					"avatar":   user.Avatar(),
				}
			}

			// 获取该用户的文章
			var articles []*model.RelArticle
			controller.app.RecordQuery(model.DbNameRelArticles).
				Where(dbx.HashExp{
					model.RelArticlesFieldActivityId: activity.Id,
					model.RelArticlesFieldUserId:     candidateUserId,
				}).
				All(&articles)

			articleList := make([]map[string]any, 0, len(articles))
			for _, art := range articles {
				articleList = append(articleList, map[string]any{
					"id":           art.Id,
					"oId":          art.OId(),
					"title":        art.Title(),
					"viewCount":    art.ViewCount(),
					"goodCnt":      art.GoodCnt(),
					"commentCount": art.CommentCount(),
					"collectCnt":   art.CollectCnt(),
					"thankCnt":     art.ThankCnt(),
				})
			}

			candidates = append(candidates, map[string]any{
				"userId":   candidateUserId,
				"user":     userData,
				"articles": articleList,
			})
		}
	}

	// 获取当前用户本轮已投票记录
	var myVotes []*model.VoteJuryLog
	controller.app.RecordQuery(model.DbNameVoteJuryLogs).
		Where(dbx.HashExp{
			model.VoteJuryLogFieldVoteId:     voteId,
			model.VoteJuryLogFieldFromUserId: userId,
			model.VoteJuryLogFieldRound:      currentRound,
		}).
		All(&myVotes)

	votedUsers := make(map[string]int)
	for _, v := range myVotes {
		votedUsers[v.ToUserId()] += v.Times()
	}

	usedVotes := 0
	for _, times := range votedUsers {
		usedVotes += times
	}

	return event.JSON(http.StatusOK, map[string]any{
		"candidates":     candidates,
		"currentRound":   currentRound,
		"totalVotes":     vote.Times(),
		"usedVotes":      usedVotes,
		"remainingVotes": vote.Times() - usedVotes,
		"allowRepeat":    vote.Repeat(),
		"votedUsers":     votedUsers,
	})
}

// GetVoteDetails 获取投票详情（管理员）
func (controller *VoteJuryController) GetVoteDetails(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")

	rule := event.Get("jury_rule").(*model.VoteJuryRule)
	currentRound := rule.CurrentRound()
	if currentRound == 0 {
		currentRound = 1
	}

	// 获取所有评审团成员
	var juryUsers []*model.VoteJuryUser
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldStatus: model.VoteJuryUserStatusApproved,
		}).
		All(&juryUsers); err != nil {
		return event.InternalServerError("获取评审团成员失败", err)
	}

	// 获取当前轮次的所有投票记录
	var voteLogs []*model.VoteJuryLog
	if err := controller.app.RecordQuery(model.DbNameVoteJuryLogs).
		Where(dbx.HashExp{
			model.VoteJuryLogFieldVoteId: voteId,
			model.VoteJuryLogFieldRound:  currentRound,
		}).
		All(&voteLogs); err != nil {
		return event.InternalServerError("获取投票记录失败", err)
	}

	// 按投票人分组
	votesByUser := make(map[string][]*model.VoteJuryLog)
	for _, log := range voteLogs {
		votesByUser[log.FromUserId()] = append(votesByUser[log.FromUserId()], log)
	}

	// 构建详细的投票信息
	memberDetails := make([]map[string]any, 0, len(juryUsers))
	for _, juryUser := range juryUsers {
		user := new(model.User)
		var userData map[string]any
		if err := controller.app.RecordQuery(model.DbNameUsers).
			Where(dbx.HashExp{model.CommonFieldId: juryUser.UserId()}).
			One(user); err == nil {
			userData = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}

		// 该成员的投票记录
		userVotes := votesByUser[juryUser.UserId()]
		voteRecords := make([]map[string]any, 0, len(userVotes))
		for _, v := range userVotes {
			toUser := new(model.User)
			var toUserData map[string]any
			if err := controller.app.RecordQuery(model.DbNameUsers).
				Where(dbx.HashExp{model.CommonFieldId: v.ToUserId()}).
				One(toUser); err == nil {
				toUserData = map[string]any{
					"id":       toUser.Id,
					"name":     toUser.Name(),
					"nickname": toUser.Nickname(),
					"avatar":   toUser.Avatar(),
				}
			}

			voteRecords = append(voteRecords, map[string]any{
				"toUserId": v.ToUserId(),
				"toUser":   toUserData,
				"times":    v.Times(),
				"comment":  v.Comment(),
				"created":  v.GetDateTime("created").String(),
			})
		}

		memberDetails = append(memberDetails, map[string]any{
			"userId":    juryUser.UserId(),
			"user":      userData,
			"hasVoted":  len(userVotes) > 0,
			"voteCount": len(userVotes),
			"votes":     voteRecords,
		})
	}

	return event.JSON(http.StatusOK, map[string]any{
		"currentRound":  currentRound,
		"memberDetails": memberDetails,
	})
}
