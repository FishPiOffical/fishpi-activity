package controller

import (
	"bless-activity/model"
	"net/http"

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
	controller.event.Router.GET("/api/activities", controller.GetActivities)
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
	type ActivityResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		ArticleUrl  string `json:"articleUrl"`
		ExternalUrl string `json:"externalUrl"`
		Desc        string `json:"desc"`
		Start       string `json:"start"`
		End         string `json:"end"`
	}

	activityList := make([]ActivityResponse, 0, len(activities))

	for _, record := range activities {
		activity := model.NewActivity(record)

		activityList = append(activityList, ActivityResponse{
			ID:          activity.ProxyRecord().Id,
			Name:        activity.Name(),
			Slug:        activity.Slug(),
			ArticleUrl:  activity.ArticleUrl(),
			ExternalUrl: activity.ExternalUrl(),
			Desc:        activity.Desc(),
			Start:       activity.Start().String(),
			End:         activity.End().String(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"items": activityList,
	})
}
