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
	juryGroup.POST("/apply/audit", controller.AuditApply).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/status/switch", controller.SwitchStatus).BindFunc(controller.RequireAuth, controller.RequireAdmin)
	juryGroup.POST("/calculate", controller.Calculate).BindFunc(controller.RequireAuth, controller.RequireAdmin)

	// 用户接口
	juryGroup.POST("/apply", controller.Apply).BindFunc(controller.RequireAuth)
	juryGroup.POST("/vote", controller.Vote).BindFunc(controller.RequireAuth)
	juryGroup.GET("/result/{voteId}", controller.GetResult)
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
			model.VoteJuryUserFieldStatus: model.JuryUserStatusApproved,
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
			OrderBy("-created").
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
				"userId": userId,
				"count":  count,
				"user":   userData,
			})
		}

		roundResults = append(roundResults, map[string]any{
			"round":    result.Round(),
			"results":  resultWithUsers,
			"continue": result.Continue(),
			"userIds":  result.UserIds(),
		})
	}

	// 获取投票进度（仅评审中状态时）
	var votingProgress map[string]any
	if rule.Status() == model.JuryStatusVoting {
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
		"members":        members,
		"applyLogs":      applyLogs,
		"results":        roundResults,
		"votingProgress": votingProgress,
		"isAdmin":        isAdmin,
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
			model.VoteJuryUserFieldStatus: model.JuryUserStatusApproved,
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
	juryUser.SetStatus(model.JuryUserStatusApproved)

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
	applyLog.SetStatus(model.JuryApplyStatusApproved)
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
	if applyLog.Status() != model.JuryApplyStatusPending {
		return event.BadRequestError("该申请已处理", nil)
	}

	// 更新申请状态
	newStatus := model.JuryApplyStatusApproved
	if data.Status == "rejected" {
		newStatus = model.JuryApplyStatusRejected
	}
	applyLog.SetStatus(newStatus)
	applyLog.SetAdminId(event.Auth.Id)

	if err := controller.app.Save(applyLog); err != nil {
		return event.InternalServerError("更新申请状态失败", err)
	}

	// 如果通过，创建评审团成员记录
	if newStatus == model.JuryApplyStatusApproved {
		// 获取评审团规则检查席位
		rule := event.Get("jury_rule").(*model.VoteJuryRule)

		// 统计已通过的成员数量
		var approvedMembers []*model.VoteJuryUser
		if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
			Where(dbx.HashExp{
				model.VoteJuryUserFieldVoteId: data.VoteId,
				model.VoteJuryUserFieldStatus: model.JuryUserStatusApproved,
			}).
			All(&approvedMembers); err != nil {
			return event.InternalServerError("获取评审团成员失败", err)
		}

		if len(approvedMembers) >= rule.Count() {
			// 席位已满，将申请状态改回待审核
			applyLog.SetStatus(model.JuryApplyStatusPending)
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
			existingMember.SetStatus(model.JuryUserStatusApproved)
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
			juryUser.SetStatus(model.JuryUserStatusApproved)

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
	newStatus, err := model.ParseJuryStatus(data.NewStatus)
	if err != nil {
		return event.BadRequestError("状态值无效", nil)
	}

	rule := event.Get("jury_rule").(*model.VoteJuryRule)
	currentStatus := rule.Status()

	// 状态流转校验
	validTransitions := map[model.JuryStatus][]model.JuryStatus{
		model.JuryStatusPending:   {model.JuryStatusApplying},
		model.JuryStatusApplying:  {model.JuryStatusPublicity},
		model.JuryStatusPublicity: {model.JuryStatusVoting},
		model.JuryStatusVoting:    {model.JuryStatusCompleted},
		model.JuryStatusCompleted: {}, // 终态
	}

	validNextStatuses := validTransitions[currentStatus]
	if !slices.Contains(validNextStatuses, newStatus) {
		return event.BadRequestError(fmt.Sprintf("不能从 %s 切换到 %s", currentStatus, newStatus), nil)
	}

	// 进入评审状态时，初始化轮次
	if newStatus == model.JuryStatusVoting && rule.CurrentRound() == 0 {
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
	if rule.Status() != model.JuryStatusVoting {
		return event.BadRequestError("当前状态不是评审中", nil)
	}

	currentRound := rule.CurrentRound()
	if currentRound == 0 {
		currentRound = 1
	}

	// 获取投票配置
	vote := new(model.Vote)
	if err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{model.CommonFieldId: data.VoteId}).
		One(vote); err != nil {
		return event.InternalServerError("获取投票配置失败", err)
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

	// 判断是否需要进入下一轮
	needNextRound := false
	var nextRoundUsers []string

	if len(topUsers) > 1 {
		// 有平票，检查决策票
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

		if len(decisionTopUsers) > 1 || maxDecisionVotes == 0 {
			// 决策票也无法决定，进入下一轮
			needNextRound = true
			nextRoundUsers = topUsers
		}
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
		result.SetUserIds(nextRoundUsers)
	}

	if err := controller.app.Save(result); err != nil {
		return event.InternalServerError("保存投票结果失败", err)
	}

	// 如果需要下一轮
	if needNextRound {
		rule.SetCurrentRound(currentRound + 1)
		if err := controller.app.Save(rule); err != nil {
			return event.InternalServerError("更新轮次失败", err)
		}

		return event.JSON(http.StatusOK, map[string]any{
			"message":        "平票，进入下一轮投票",
			"needNextRound":  true,
			"nextRound":      currentRound + 1,
			"nextRoundUsers": nextRoundUsers,
			"results":        voteCount,
		})
	}

	return event.JSON(http.StatusOK, map[string]any{
		"message":       "投票结果已统计",
		"needNextRound": false,
		"results":       voteCount,
		"winner":        topUsers[0],
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
	if rule.Status() != model.JuryStatusApplying {
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
	applyLog.SetStatus(model.JuryApplyStatusPending)

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
	if rule.Status() != model.JuryStatusVoting {
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
			model.VoteJuryUserFieldStatus: model.JuryUserStatusApproved,
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

// GetResult 获取投票结果
func (controller *VoteJuryController) GetResult(event *core.RequestEvent) error {
	voteId := event.Request.PathValue("voteId")

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
				"userId": userId,
				"count":  count,
				"user":   userData,
			})
		}

		roundResults = append(roundResults, map[string]any{
			"round":    result.Round(),
			"results":  resultWithUsers,
			"continue": result.Continue(),
			"userIds":  result.UserIds(),
		})
	}

	// 获取评审团规则
	rule := new(model.VoteJuryRule)
	if err := controller.app.RecordQuery(model.DbNameVoteJuryRules).
		Where(dbx.HashExp{model.VoteJuryRuleFieldVoteId: voteId}).
		One(rule); err != nil {
		return event.InternalServerError("获取评审团规则失败", err)
	}

	// 获取已通过的评审团成员
	var juryUsers []*model.VoteJuryUser
	if err := controller.app.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldStatus: model.JuryUserStatusApproved,
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

	return event.JSON(http.StatusOK, map[string]any{
		"status":       rule.Status(),
		"currentRound": rule.CurrentRound(),
		"results":      roundResults,
		"members":      members,
	})
}
