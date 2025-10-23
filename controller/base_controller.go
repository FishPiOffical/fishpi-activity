package controller

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type BaseController struct {
	event *core.ServeEvent
	app   core.App
}

func NewBaseController(event *core.ServeEvent) *BaseController {
	controller := &BaseController{
		event: event,
		app:   event.App,
	}
	return controller
}

func (controller *BaseController) CheckActivity(event *core.RequestEvent) error {

	endTime := time.Date(2025, 10, 20, 0, 0, 0, 0, time.Local)
	if time.Now().After(endTime) {
		return event.ForbiddenError("活动已结束", nil)
	}

	return event.Next()
}
