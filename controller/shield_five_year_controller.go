package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

type ShieldFiveYearController struct {
	event  *core.ServeEvent
	app    core.App
	logger *slog.Logger
}

func NewShieldFiveYearController(event *core.ServeEvent) *ShieldFiveYearController {
	logger := event.App.Logger().With(
		slog.String("controller", "shield_five_year"),
	)

	controller := &ShieldFiveYearController{
		event:  event,
		app:    event.App,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *ShieldFiveYearController) registerRoutes() {
	slog.Info("注册路由")
	group := controller.event.Router.Group("/activity-api/shield-five-year")

	// 徽章相关接口
	group.POST("/shields", controller.CreateShield).BindFunc(controller.CheckLogin)
	group.GET("/shields/{activityId}", controller.GetShieldsByActivity)
	group.PATCH("/shields/{id}", controller.UpdateShield).BindFunc(controller.CheckLogin)

	// 文章相关接口（关键词活动）
	group.POST("/articles", controller.CreateArticle).BindFunc(controller.CheckLogin)
	group.GET("/articles/{activityId}", controller.GetArticlesByActivity)
	group.PATCH("/articles/{id}", controller.UpdateArticle).BindFunc(controller.CheckLogin)
	group.GET("/my-articles", controller.GetMyArticles).BindFunc(controller.CheckLogin)

	// 投票相关接口
	group.POST("/vote", controller.Vote).BindFunc(controller.CheckLogin)
	group.DELETE("/vote/{id}", controller.DeleteVote).BindFunc(controller.CheckLogin)
	group.GET("/votes/{activityId}", controller.GetVotesByActivity)
	group.GET("/vote-stats/{activityId}", controller.GetVoteStats)
	group.GET("/my-votes", controller.GetMyVotes).BindFunc(controller.CheckLogin)
	group.GET("/vote-quota/{activityId}", controller.GetVoteQuota).BindFunc(controller.CheckLogin)
	group.GET("/vote-details/{activityId}/{userId}", controller.GetVoteDetails)
}

func (controller *ShieldFiveYearController) CheckLogin(event *core.RequestEvent) error {
	if event.Auth == nil {
		return event.UnauthorizedError("未登录", nil)
	}
	if event.HasSuperuserAuth() {
		return event.ForbiddenError("请登录普通用户账号", nil)
	}
	return event.Next()
}

// CreateShield 创建徽章
func (controller *ShieldFiveYearController) CreateShield(e *core.RequestEvent) error {
	data := struct {
		ActivityId string `json:"activityId"`
		Text       string `json:"text"`
		Url        string `json:"url"`
		Backcolor  string `json:"backcolor"`
		Fontcolor  string `json:"fontcolor"`
		Ver        string `json:"ver"`
		Scale      string `json:"scale"`
		Size       string `json:"size"`
		Border     string `json:"border"`
		BarLen     string `json:"barlen"`
		Fontsize   string `json:"fontsize"`
		BarRadius  string `json:"barradius"`
		Shadow     string `json:"shadow"`
		Anime      string `json:"anime"`
		Title      string `json:"title"`
		Note       string `json:"note"`
	}{}

	if err := e.BindBody(&data); err != nil {
		return e.BadRequestError("参数错误", err)
	}

	if data.ActivityId == "" || data.Text == "" {
		return e.BadRequestError("活动ID和文本不能为空", nil)
	}

	user := model.NewUser(e.Auth)

	// 检查用户是否已经为该活动创建过徽章
	existing, err := controller.app.FindFirstRecordByFilter(
		model.DbNameShields,
		"activityId = {:activityId} && userId = {:userId}",
		map[string]any{
			"activityId": data.ActivityId,
			"userId":     user.Id,
		},
	)

	var shield *model.Shield
	var collection *core.Collection

	collection, err = controller.app.FindCollectionByNameOrId(model.DbNameShields)
	if err != nil {
		return e.InternalServerError("获取集合失败", err)
	}

	if existing != nil {
		// 更新现有徽章
		shield = model.NewShield(existing)
	} else {
		// 创建新徽章
		shield = model.NewShieldFromCollection(collection)
		shield.Set("activityId", data.ActivityId)
		shield.Set("userId", user.Id)
	}

	shield.SetText(data.Text)
	shield.SetUrl(data.Url)
	shield.SetBackcolor(data.Backcolor)
	shield.SetFontcolor(data.Fontcolor)

	// 保存标题和设计思路
	if data.Title != "" {
		shield.Set("title", data.Title)
	}
	if data.Note != "" {
		shield.Set("note", data.Note)
	}

	if data.Ver != "" {
		shield.Set(model.ShieldsFieldVer, data.Ver)
	}
	if data.Scale != "" {
		shield.Set(model.ShieldsFieldScale, data.Scale)
	}
	if data.Size != "" {
		shield.Set(model.ShieldsFieldSize, data.Size)
	}
	if data.Border != "" {
		shield.Set(model.ShieldsFieldBorder, data.Border)
	}
	if data.BarLen != "" {
		shield.Set(model.ShieldsFieldBarLen, data.BarLen)
	}
	if data.Fontsize != "" {
		shield.Set(model.ShieldsFieldFontsize, data.Fontsize)
	}
	if data.BarRadius != "" {
		shield.Set(model.ShieldsFieldBarRadius, data.BarRadius)
	}
	if data.Shadow != "" {
		shield.Set(model.ShieldsFieldShadow, data.Shadow)
	}
	if data.Anime != "" {
		shield.Set(model.ShieldsFieldAnime, data.Anime)
	}

	if err := controller.app.Save(shield.ProxyRecord()); err != nil {
		return e.InternalServerError("保存徽章失败", err)
	}

	// 返回完整的徽章信息，包括预览URL
	return e.JSON(http.StatusOK, map[string]any{
		"id":      shield.Id,
		"message": "徽章保存成功",
		"shield": map[string]any{
			"id":        shield.Id,
			"text":      shield.Text(),
			"url":       shield.Url(),
			"backcolor": shield.Backcolor(),
			"fontcolor": shield.Fontcolor(),
			"title":     shield.GetString("title"),
			"note":      shield.GetString("note"),
		},
	})
}

// GetShieldsByActivity 获取活动的所有徽章
func (controller *ShieldFiveYearController) GetShieldsByActivity(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("activityId")

	records, err := controller.app.FindRecordsByFilter(
		model.DbNameShields,
		"activityId = {:activityId}",
		"-created",
		0,
		0,
		map[string]any{
			"activityId": activityId,
		},
	)

	if err != nil {
		return e.InternalServerError("获取徽章列表失败", err)
	}

	// 扩展用户信息
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		shield := model.NewShield(record)
		userId := shield.GetString("userId")

		userRecord, _ := controller.app.FindRecordById(model.DbNameUsers, userId)
		var userData map[string]any
		if userRecord != nil {
			user := model.NewUser(userRecord)
			userData = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}

		result = append(result, map[string]any{
			"id":        shield.Id,
			"text":      shield.Text(),
			"url":       shield.Url(),
			"backcolor": shield.Backcolor(),
			"fontcolor": shield.Fontcolor(),
			"title":     shield.GetString("title"),
			"note":      shield.GetString("note"),
			"created":   shield.GetDateTime(model.ShieldsFieldCreated).String(),
			"user":      userData,
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"items": result,
	})
}

// CreateArticle 创建文章（关键词）
func (controller *ShieldFiveYearController) CreateArticle(e *core.RequestEvent) error {
	data := struct {
		ActivityId string `json:"activityId"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		ShieldId   string `json:"shieldId"`
	}{}

	if err := e.BindBody(&data); err != nil {
		return e.BadRequestError("参数错误", err)
	}

	if data.ActivityId == "" || data.Title == "" || data.Content == "" {
		return e.BadRequestError("活动ID、标题和内容不能为空", nil)
	}

	user := model.NewUser(e.Auth)

	// 检查用户是否已经为该活动创建过文章
	existing, err := controller.app.FindFirstRecordByFilter(
		model.DbNameArticles,
		"activityId = {:activityId} && userId = {:userId}",
		map[string]any{
			"activityId": data.ActivityId,
			"userId":     user.Id,
		},
	)

	var article *model.Article
	var collection *core.Collection

	collection, err = controller.app.FindCollectionByNameOrId(model.DbNameArticles)
	if err != nil {
		return e.InternalServerError("获取集合失败", err)
	}

	if existing != nil {
		// 更新现有文章
		article = model.NewArticle(existing)
	} else {
		// 创建新文章
		article = model.NewArticleFromCollection(collection)
		article.SetActivityId(data.ActivityId)
		article.SetUserId(user.Id)
	}

	article.SetTitle(data.Title)
	article.SetContent(data.Content)
	if data.ShieldId != "" {
		article.SetShieldId(data.ShieldId)
	}

	if err := controller.app.Save(article.ProxyRecord()); err != nil {
		return e.InternalServerError("保存文章失败", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"id":      article.Id,
		"message": "文章保存成功",
	})
}

// GetArticlesByActivity 获取活动的所有文章
func (controller *ShieldFiveYearController) GetArticlesByActivity(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("activityId")

	records, err := controller.app.FindRecordsByFilter(
		model.DbNameArticles,
		"activityId = {:activityId}",
		"-created",
		0,
		0,
		map[string]any{
			"activityId": activityId,
		},
	)

	if err != nil {
		return e.InternalServerError("获取文章列表失败", err)
	}

	// 扩展用户信息和徽章信息
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		article := model.NewArticle(record)
		userId := article.UserId()

		userRecord, _ := controller.app.FindRecordById(model.DbNameUsers, userId)
		var userData map[string]any
		if userRecord != nil {
			user := model.NewUser(userRecord)
			userData = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}

		// 如果文章关联了徽章，获取徽章信息
		var shieldData map[string]any
		shieldId := article.ShieldId()
		if shieldId != "" {
			shieldRecord, err := controller.app.FindRecordById(model.DbNameShields, shieldId)
			if err == nil && shieldRecord != nil {
				shield := model.NewShield(shieldRecord)
				shieldData = map[string]any{
					"id":        shield.Id,
					"text":      shield.Text(),
					"url":       shield.Url(),
					"backcolor": shield.Backcolor(),
					"fontcolor": shield.Fontcolor(),
					"title":     shield.GetString("title"),
					"note":      shield.GetString("note"),
				}
			}
		}

		result = append(result, map[string]any{
			"id":       article.Id,
			"title":    article.Title(),
			"content":  article.Content(),
			"shieldId": shieldId,
			"shield":   shieldData,
			"created":  article.GetDateTime(model.ArticlesFieldCreated).String(),
			"user":     userData,
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"items": result,
	})
}

// Vote 投票
func (controller *ShieldFiveYearController) Vote(e *core.RequestEvent) error {
	data := struct {
		VoteId     string `json:"voteId"`
		ActivityId string `json:"activityId"`
		ToUserId   string `json:"toUserId"`
		TargetType string `json:"targetType"` // "article" or "shield"
		TargetId   string `json:"targetId"`
		Comment    string `json:"comment"`
	}{}

	if err := e.BindBody(&data); err != nil {
		return e.BadRequestError("参数错误", err)
	}

	if data.ToUserId == "" {
		return e.BadRequestError("目标用户ID不能为空", nil)
	}

	user := model.NewUser(e.Auth)

	// 如果没有提供voteId，通过activityId查找
	voteId := data.VoteId
	if voteId == "" && data.ActivityId != "" {
		activity, err := controller.app.FindRecordById(model.DbNameActivities, data.ActivityId)
		if err != nil {
			return e.BadRequestError("活动不存在", err)
		}
		activityModel := model.NewActivity(activity)
		voteId = activityModel.VoteId()

		if voteId == "" {
			return e.BadRequestError("该活动未配置投票", nil)
		}
	}

	if voteId == "" {
		return e.BadRequestError("投票ID不能为空", nil)
	}

	// 检查是否已经投过票
	existing, _ := controller.app.FindFirstRecordByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId} && fromUserId = {:fromUserId} && toUserId = {:toUserId}",
		map[string]any{
			"voteId":     voteId,
			"fromUserId": user.Id,
			"toUserId":   data.ToUserId,
		},
	)

	if existing != nil {
		return e.BadRequestError("您已经为该作品投过票了", nil)
	}

	// 获取投票配置，检查投票次数限制
	vote, err := controller.app.FindRecordById(model.DbNameVotes, voteId)
	if err != nil {
		return e.BadRequestError("投票活动不存在", err)
	}

	voteConfig := model.NewVote(vote)
	voteTimes := voteConfig.GetInt(model.VotesFieldTimes)

	// 检查用户已投票次数
	userVoteLogs, err := controller.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId} && fromUserId = {:fromUserId}",
		"",
		0,
		0,
		map[string]any{
			"voteId":     voteId,
			"fromUserId": user.Id,
		},
	)

	if err != nil {
		return e.InternalServerError("检查投票次数失败", err)
	}

	if len(userVoteLogs) >= voteTimes {
		return e.BadRequestError("您的投票次数已用完", nil)
	}

	// 创建投票记录
	collection, err := controller.app.FindCollectionByNameOrId(model.DbNameVoteLogs)
	if err != nil {
		return e.InternalServerError("获取集合失败", err)
	}

	voteLog := model.NewVoteLogFromCollection(collection)
	voteLog.SetVoteId(voteId)
	voteLog.SetFromUserId(user.Id)
	voteLog.SetToUserId(data.ToUserId)
	voteLog.SetComment(data.Comment)

	if err := controller.app.Save(voteLog.ProxyRecord()); err != nil {
		return e.InternalServerError("保存投票失败", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "投票成功",
	})
}

// GetVotesByActivity 获取活动的投票记录
func (controller *ShieldFiveYearController) GetVotesByActivity(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("activityId")

	// 先获取活动关联的投票ID
	activity, err := controller.app.FindRecordById(model.DbNameActivities, activityId)
	if err != nil {
		return e.BadRequestError("活动不存在", err)
	}

	activityModel := model.NewActivity(activity)
	voteId := activityModel.VoteId()

	if voteId == "" {
		return e.JSON(http.StatusOK, map[string]any{
			"items": []any{},
		})
	}

	records, err := controller.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId}",
		"-created",
		0,
		0,
		map[string]any{
			"voteId": voteId,
		},
	)

	if err != nil {
		return e.InternalServerError("获取投票记录失败", err)
	}

	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		voteLog := model.NewVoteLog(record)

		fromUser, _ := controller.app.FindRecordById(model.DbNameUsers, voteLog.FromUserId())
		toUser, _ := controller.app.FindRecordById(model.DbNameUsers, voteLog.ToUserId())

		var fromUserData, toUserData map[string]any
		if fromUser != nil {
			user := model.NewUser(fromUser)
			fromUserData = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}
		if toUser != nil {
			user := model.NewUser(toUser)
			toUserData = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}

		result = append(result, map[string]any{
			"id":       voteLog.Id,
			"fromUser": fromUserData,
			"toUser":   toUserData,
			"comment":  voteLog.Comment(),
			"created":  voteLog.GetDateTime(model.VoteLogsFieldCreated).String(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"items": result,
	})
}

// GetVoteStats 获取投票统计
func (controller *ShieldFiveYearController) GetVoteStats(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("activityId")

	// 先获取活动关联的投票ID
	activity, err := controller.app.FindRecordById(model.DbNameActivities, activityId)
	if err != nil {
		return e.BadRequestError("活动不存在", err)
	}

	activityModel := model.NewActivity(activity)
	voteId := activityModel.VoteId()

	if voteId == "" {
		return e.JSON(http.StatusOK, map[string]any{
			"stats": map[string]int{},
		})
	}

	// 获取所有投票记录
	records, err := controller.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId}",
		"",
		0,
		0,
		map[string]any{
			"voteId": voteId,
		},
	)

	if err != nil {
		return e.InternalServerError("获取投票统计失败", err)
	}

	// 统计每个用户获得的票数
	stats := make(map[string]int)
	for _, record := range records {
		voteLog := model.NewVoteLog(record)
		toUserId := voteLog.ToUserId()
		stats[toUserId]++
	}

	return e.JSON(http.StatusOK, map[string]any{
		"stats": stats,
	})
}

// UpdateShield 更新徽章
func (controller *ShieldFiveYearController) UpdateShield(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	authRecord := e.Auth

	// 查找徽章记录
	record, err := controller.app.FindRecordById(model.DbNameShields, id)
	if err != nil {
		return e.NotFoundError("徽章不存在", err)
	}

	shield := model.NewShield(record)

	// 检查权限：只能更新自己的徽章
	if shield.GetString("userId") != authRecord.Id {
		return e.ForbiddenError("无权限更新此徽章", nil)
	}

	// 解析请求体
	var data map[string]any
	if err := e.BindBody(&data); err != nil {
		return e.BadRequestError("参数错误", err)
	}

	// 更新字段
	for key, value := range data {
		record.Set(key, value)
	}

	if err := controller.app.Save(record); err != nil {
		return e.InternalServerError("更新徽章失败", err)
	}

	return e.JSON(http.StatusOK, record)
}

// UpdateArticle 更新文章
func (controller *ShieldFiveYearController) UpdateArticle(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	authRecord := e.Auth

	// 查找文章记录
	record, err := controller.app.FindRecordById(model.DbNameArticles, id)
	if err != nil {
		return e.NotFoundError("文章不存在", err)
	}

	article := model.NewArticle(record)

	// 检查权限：只能更新自己的文章
	if article.UserId() != authRecord.Id {
		return e.ForbiddenError("无权限更新此文章", nil)
	}

	// 解析请求体
	var data map[string]any
	if err := e.BindBody(&data); err != nil {
		return e.BadRequestError("参数错误", err)
	}

	// 更新字段
	for key, value := range data {
		record.Set(key, value)
	}

	if err := controller.app.Save(record); err != nil {
		return e.InternalServerError("更新文章失败", err)
	}

	return e.JSON(http.StatusOK, record)
}

// GetMyArticles 获取当前用户的投稿
func (controller *ShieldFiveYearController) GetMyArticles(e *core.RequestEvent) error {
	authRecord := e.Auth

	records, err := controller.app.FindRecordsByFilter(
		model.DbNameArticles,
		"userId = {:userId}",
		"",
		0,
		0,
		map[string]any{
			"userId": authRecord.Id,
		},
	)

	if err != nil {
		return e.InternalServerError("获取投稿失败", err)
	}

	// 扩展用户和徽章信息
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		article := model.NewArticle(record)
		data := record.PublicExport()

		// 获取用户信息
		if userRecord, err := controller.app.FindRecordById(model.DbNameUsers, article.UserId()); err == nil {
			user := model.NewUser(userRecord)
			data["user"] = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}

		// 获取徽章信息
		if article.ShieldId() != "" {
			if shieldRecord, err := controller.app.FindRecordById(model.DbNameShields, article.ShieldId()); err == nil {
				data["shield"] = shieldRecord.PublicExport()
			}
		}

		result = append(result, data)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"items": result,
	})
}

// GetMyVotes 获取当前用户的投票记录
func (controller *ShieldFiveYearController) GetMyVotes(e *core.RequestEvent) error {
	authRecord := e.Auth

	records, err := controller.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"fromUserId = {:userId}",
		"",
		0,
		0,
		map[string]any{
			"userId": authRecord.Id,
		},
	)

	if err != nil {
		return e.InternalServerError("获取投票记录失败", err)
	}

	// 扩展信息：获取每个投票对应的活动ID
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		voteLog := model.NewVoteLog(record)
		data := record.PublicExport()

		// 通过voteId查找对应的活动
		activities, err := controller.app.FindRecordsByFilter(
			model.DbNameActivities,
			"voteId = {:voteId}",
			"",
			1,
			0,
			map[string]any{
				"voteId": voteLog.VoteId(),
			},
		)

		if err == nil && len(activities) > 0 {
			data["activityId"] = activities[0].Id
		}

		result = append(result, data)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"items": result,
	})
}

// GetVoteQuota 获取投票配额
func (controller *ShieldFiveYearController) GetVoteQuota(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("activityId")
	authRecord := e.Auth

	// 获取活动信息
	activity, err := controller.app.FindRecordById(model.DbNameActivities, activityId)
	if err != nil {
		return e.BadRequestError("活动不存在", err)
	}

	activityModel := model.NewActivity(activity)
	voteId := activityModel.VoteId()

	if voteId == "" {
		return e.JSON(http.StatusOK, map[string]any{
			"total":     0,
			"used":      0,
			"remaining": 0,
		})
	}

	// 获取投票配置
	vote, err := controller.app.FindRecordById(model.DbNameVotes, voteId)
	if err != nil {
		return e.BadRequestError("投票不存在", err)
	}

	voteModel := model.NewVote(vote)
	maxVotes := voteModel.GetInt(model.VotesFieldTimes)

	// 统计已使用的票数
	usedVotes, err := controller.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId} && fromUserId = {:userId}",
		"",
		0,
		0,
		map[string]any{
			"voteId": voteId,
			"userId": authRecord.Id,
		},
	)

	if err != nil {
		return e.InternalServerError("获取已投票数失败", err)
	}

	remaining := maxVotes - len(usedVotes)
	if remaining < 0 {
		remaining = 0
	}

	return e.JSON(http.StatusOK, map[string]any{
		"total":     maxVotes,
		"used":      len(usedVotes),
		"remaining": remaining,
	})
}

// GetVoteDetails 获取投票详情
func (controller *ShieldFiveYearController) GetVoteDetails(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("activityId")
	userId := e.Request.PathValue("userId")

	// 获取活动信息
	activity, err := controller.app.FindRecordById(model.DbNameActivities, activityId)
	if err != nil {
		return e.BadRequestError("活动不存在", err)
	}

	activityModel := model.NewActivity(activity)
	voteId := activityModel.VoteId()

	if voteId == "" {
		return e.JSON(http.StatusOK, map[string]any{
			"voters": []any{},
		})
	}

	// 获取给该用户投票的所有记录
	records, err := controller.app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"voteId = {:voteId} && toUserId = {:toUserId}",
		"-created",
		0,
		0,
		map[string]any{
			"voteId":   voteId,
			"toUserId": userId,
		},
	)

	if err != nil {
		return e.InternalServerError("获取投票详情失败", err)
	}

	// 扩展投票人信息
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		voteLog := model.NewVoteLog(record)
		data := record.PublicExport()

		// 获取投票人信息
		if userRecord, err := controller.app.FindRecordById(model.DbNameUsers, voteLog.FromUserId()); err == nil {
			user := model.NewUser(userRecord)
			data["user"] = map[string]any{
				"id":       user.Id,
				"name":     user.Name(),
				"nickname": user.Nickname(),
				"avatar":   user.Avatar(),
			}
		}

		result = append(result, data)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"voters": result,
	})
}

// DeleteVote 删除投票（取消投票）
func (controller *ShieldFiveYearController) DeleteVote(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	authRecord := e.Auth

	// 查找投票记录
	record, err := controller.app.FindRecordById(model.DbNameVoteLogs, id)
	if err != nil {
		return e.NotFoundError("投票记录不存在", err)
	}

	voteLog := model.NewVoteLog(record)

	// 检查权限：只能删除自己的投票
	if voteLog.FromUserId() != authRecord.Id {
		return e.ForbiddenError("无权限删除此投票", nil)
	}

	// 删除投票记录
	if err := controller.app.Delete(record); err != nil {
		return e.InternalServerError("删除投票失败", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "投票已取消",
	})
}
