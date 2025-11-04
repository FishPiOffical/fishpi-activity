package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type ActivityController struct {
	event *core.ServeEvent
	app   core.App
}

func NewActivityController(event *core.ServeEvent) *ActivityController {
	controller := &ActivityController{
		event: event,
		app:   event.App,
	}

	controller.registerRoutes()

	return controller
}

func (controller *ActivityController) registerRoutes() {
	controller.event.Router.GET("/activity-api/activities", controller.GetActivities)
	controller.event.Router.GET("/activity-api/activities/{id}", controller.GetActivityRewards)
	controller.event.Router.GET("/activity-api/yearly-histories", controller.GetYearlyHistories)
	controller.event.Router.GET("/activity-api/recent", controller.GetActivityList)
}

func (controller *ActivityController) GetActivities(e *core.RequestEvent) error {
	// 查询所有未隐藏的活动
	activities, err := controller.app.FindRecordsByFilter(
		model.DbNameActivities,
		model.ActivitiesFieldHideInList+" = false",
		"-"+model.ActivitiesFieldStart, // 按开始时间倒序排列
		0,
		0,
	)

	if err != nil {
		return e.InternalServerError("Failed to load activities", err)
	}

	// 准备返回数据
	type RewardItem struct {
		Name  string `json:"name"`
		Min   int    `json:"min"`
		Max   int    `json:"max"`
		Point int    `json:"point"`
		More  string `json:"more"`
	}

	type ActivityResponse struct {
		ID          string       `json:"id"`
		Name        string       `json:"name"`
		Slug        string       `json:"slug"`
		ArticleUrl  string       `json:"articleUrl"`
		ExternalUrl string       `json:"externalUrl"`
		Desc        string       `json:"desc"`
		Start       string       `json:"start"`
		End         string       `json:"end"`
		Rewards     []RewardItem `json:"rewards,omitempty"`
	}

	activityList := make([]ActivityResponse, 0, len(activities))

	for _, record := range activities {
		activity := model.NewActivity(record)

		activityResp := ActivityResponse{
			ID:          activity.ProxyRecord().Id,
			Name:        activity.Name(),
			Slug:        activity.Slug(),
			ArticleUrl:  activity.ArticleUrl(),
			ExternalUrl: activity.ExternalUrl(),
			Desc:        activity.Desc(),
			Start:       activity.Start().String(),
			End:         activity.End().String(),
		}

		// 查询活动关联的奖励信息
		rewardGroupId := activity.RewardGroupId()
		slog.Info("活动信息", slog.String("activity", activity.Name()), slog.String("rewardGroupId", rewardGroupId))
		if rewardGroupId != "" {
			// 根据 rewardGroupId 查询 rewards 表
			rewards, err := controller.app.FindRecordsByFilter(
				model.DbNameRewards,
				model.RewardsFieldRewardGroupId+" = {:rewardGroupId}",
				model.RewardsFieldMin, // 按最小名次排序
				0,
				0,
				map[string]any{"rewardGroupId": rewardGroupId},
			)

			slog.Info("奖励内容 ", slog.String("activity", activity.Name()), slog.Any("err", err), slog.Int("count", len(rewards)))
			if err == nil && len(rewards) > 0 {
				rewardItems := make([]RewardItem, 0, len(rewards))
				for _, rewardRecord := range rewards {
					reward := model.NewReward(rewardRecord)
					rewardItems = append(rewardItems, RewardItem{
						Name:  reward.Name(),
						Min:   reward.Min(),
						Max:   reward.Max(),
						Point: reward.Point(),
						More:  reward.More(),
					})
				}
				activityResp.Rewards = rewardItems
			}
		}

		activityList = append(activityList, activityResp)
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"items": activityList,
	})
}

// GetActivityRewards 获取指定活动的奖励信息
func (controller *ActivityController) GetActivityRewards(e *core.RequestEvent) error {
	activityId := e.Request.PathValue("id")

	if activityId == "" {
		return e.BadRequestError("Activity ID is required", nil)
	}

	// 查询活动信息
	activity, err := controller.app.FindRecordById(model.DbNameActivities, activityId)
	if err != nil {
		return e.NotFoundError("Activity not found", err)
	}

	activityModel := model.NewActivity(activity)
	rewardGroupId := activityModel.RewardGroupId()

	type RewardItem struct {
		Name  string `json:"name"`
		Min   int    `json:"min"`
		Max   int    `json:"max"`
		Point int    `json:"point"`
		More  string `json:"more"`
	}

	var rewardItems []RewardItem

	if rewardGroupId != "" {
		// 根据 rewardGroupId 查询 rewards 表
		rewards, err := controller.app.FindRecordsByFilter(
			model.DbNameRewards,
			model.RewardsFieldRewardGroupId+" = {:rewardGroupId}",
			model.RewardsFieldMin, // 按最小名次排序
			0,
			0,
			map[string]any{"rewardGroupId": rewardGroupId},
		)

		if err == nil && len(rewards) > 0 {
			rewardItems = make([]RewardItem, 0, len(rewards))
			for _, rewardRecord := range rewards {
				reward := model.NewReward(rewardRecord)
				rewardItems = append(rewardItems, RewardItem{
					Name:  reward.Name(),
					Min:   reward.Min(),
					Max:   reward.Max(),
					Point: reward.Point(),
					More:  reward.More(),
				})
			}
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"activityId": activityId,
		"name":       activityModel.Name(),
		"desc":       activityModel.Desc(),
		"rewards":    rewardItems,
		"start":      activityModel.Start(),
		"end":        activityModel.End(),
	})
}

// GetActivityList 获取活动列表（正在进行的所有活动 + 即将开始的最多5个活动）
func (controller *ActivityController) GetActivityList(e *core.RequestEvent) error {
	// 获取expand参数
	expandParam := e.Request.URL.Query().Get("expand")
	expandMap := make(map[string]bool)
	if expandParam != "" {
		for _, item := range parseExpandParam(expandParam) {
			expandMap[item] = true
		}
	}

	// 查询所有未隐藏的活动
	activities, err := controller.app.FindRecordsByFilter(
		model.DbNameActivities,
		model.ActivitiesFieldHideInList+" = false",
		model.ActivitiesFieldStart, // 按开始时间正序排列
		0,
		0,
	)

	if err != nil {
		return e.InternalServerError("Failed to load activities", err)
	}

	type RewardItem struct {
		Name  string `json:"name"`
		Min   int    `json:"min"`
		Max   int    `json:"max"`
		Point int    `json:"point"`
		More  string `json:"more"`
	}

	type ActivityItem struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Slug        string      `json:"slug,omitempty"`
		SlugUrl     string      `json:"slugUrl,omitempty"`
		ArticleUrl  string      `json:"articleUrl,omitempty"`
		ExternalUrl string      `json:"externalUrl,omitempty"`
		Desc        string      `json:"desc,omitempty"`
		Start       string      `json:"start"`
		End         string      `json:"end"`
		FirstReward *RewardItem `json:"firstReward,omitempty"`
	}

	activeActivities := make([]ActivityItem, 0)
	upcomingActivities := make([]ActivityItem, 0)

	// 获取应用URL
	appURL := controller.app.Settings().Meta.AppURL

	for _, record := range activities {
		activity := model.NewActivity(record)

		startTime := activity.Start().Time()
		endTime := activity.End().Time()

		activityItem := ActivityItem{
			ID:    activity.ProxyRecord().Id,
			Name:  activity.Name(),
			Start: activity.Start().String(),
			End:   activity.End().String(),
		}

		// links - 包含slug, articleUrl, externalUrl
		if expandMap["links"] {
			slug := activity.Slug()
			if slug != "" {
				activityItem.Slug = slug
				activityItem.SlugUrl = appURL + "/" + slug + ".html"
			}
			activityItem.ArticleUrl = activity.ArticleUrl()
			activityItem.ExternalUrl = activity.ExternalUrl()
		}

		// details - 包含desc
		if expandMap["details"] {
			activityItem.Desc = activity.Desc()
		}

		// rewards - 包含奖励信息
		if expandMap["rewards"] {
			rewardGroupId := activity.RewardGroupId()
			if rewardGroupId != "" {
				rewards, err := controller.app.FindRecordsByFilter(
					model.DbNameRewards,
					model.RewardsFieldRewardGroupId+" = {:rewardGroupId}",
					model.RewardsFieldMin, // 按最小名次排序
					1,                     // 只取第一个
					0,
					map[string]any{"rewardGroupId": rewardGroupId},
				)

				if err == nil && len(rewards) > 0 {
					reward := model.NewReward(rewards[0])
					activityItem.FirstReward = &RewardItem{
						Name:  reward.Name(),
						Min:   reward.Min(),
						Max:   reward.Max(),
						Point: reward.Point(),
						More:  reward.More(),
					}
				}
			}
		}

		// 判断活动状态
		nowTime := time.Now()

		if (startTime.Before(nowTime) || startTime.Equal(nowTime)) && endTime.After(nowTime) {
			// 正在进行
			activeActivities = append(activeActivities, activityItem)
		} else if startTime.After(nowTime) {
			// 即将开始
			if len(upcomingActivities) < 5 {
				upcomingActivities = append(upcomingActivities, activityItem)
			}
		}
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"active":   activeActivities,
		"upcoming": upcomingActivities,
	})
}

// parseExpandParam 解析expand参数
func parseExpandParam(expand string) []string {
	result := make([]string, 0)
	for _, item := range strings.Split(expand, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetYearlyHistories 获取历年数据
func (controller *ActivityController) GetYearlyHistories(e *core.RequestEvent) error {
	// 查询所有历年数据，按年份倒序排列
	histories, err := controller.app.FindRecordsByFilter(
		model.DbNameYearlyHistories,
		"",
		"-"+model.YearlyHistoriesFieldYear, // 按年份倒序
		0,
		0,
	)

	if err != nil {
		return e.InternalServerError("Failed to load yearly histories", err)
	}

	type ShieldInfo struct {
		Text      string `json:"text"`
		Url       string `json:"url,omitempty"`
		Backcolor string `json:"backcolor"`
		Fontcolor string `json:"fontcolor"`
		Ver       string `json:"ver,omitempty"`
		Scale     string `json:"scale,omitempty"`
		Size      string `json:"size,omitempty"`
		Border    string `json:"border,omitempty"`
		Barlen    string `json:"barlen,omitempty"`
		Fontsize  string `json:"fontsize,omitempty"`
		Barradius string `json:"barradius,omitempty"`
		Shadow    string `json:"shadow,omitempty"`
		Anime     string `json:"anime,omitempty"`
	}

	type YearlyHistoryResponse struct {
		Year          int         `json:"year"`
		Keyword       string      `json:"keyword"`
		ArticleUrl    string      `json:"articleUrl"`
		ArticleShield *ShieldInfo `json:"articleShield,omitempty"`
		AgeShield     *ShieldInfo `json:"ageShield,omitempty"`
	}

	historyList := make([]YearlyHistoryResponse, 0, len(histories))

	for _, record := range histories {
		history := model.NewYearlyHistory(record)

		response := YearlyHistoryResponse{
			Year:       history.Year(),
			Keyword:    history.Keyword(),
			ArticleUrl: history.ArticleUrl(),
		}

		// 查询流行色徽章信息
		if articleShieldId := history.ArticleShieldId(); articleShieldId != "" {
			if shieldRecord, err := controller.app.FindRecordById(model.DbNameShields, articleShieldId); err == nil {
				shield := model.NewShield(shieldRecord)
				response.ArticleShield = &ShieldInfo{
					Text:      shield.Text(),
					Url:       shield.Url(),
					Backcolor: shield.Backcolor(),
					Fontcolor: shield.Fontcolor(),
					Ver:       shield.Ver(),
					Scale:     shield.Scale(),
					Size:      shield.Size(),
					Border:    shield.Border(),
					Barlen:    shield.BarLen(),
					Fontsize:  shield.Fontsize(),
					Barradius: shield.BarRadius(),
					Shadow:    shield.Shadow(),
					Anime:     shield.Anime(),
				}
			}
		}

		// 查询周岁徽章信息
		if ageShieldId := history.AgeShieldId(); ageShieldId != "" {
			if shieldRecord, err := controller.app.FindRecordById(model.DbNameShields, ageShieldId); err == nil {
				shield := model.NewShield(shieldRecord)
				response.AgeShield = &ShieldInfo{
					Text:      shield.Text(),
					Url:       shield.Url(),
					Backcolor: shield.Backcolor(),
					Fontcolor: shield.Fontcolor(),
					Ver:       shield.Ver(),
					Scale:     shield.Scale(),
					Size:      shield.Size(),
					Border:    shield.Border(),
					Barlen:    shield.BarLen(),
					Fontsize:  shield.Fontsize(),
					Barradius: shield.BarRadius(),
					Shadow:    shield.Shadow(),
					Anime:     shield.Anime(),
				}
			}
		}

		historyList = append(historyList, response)
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"items": historyList,
	})
}
