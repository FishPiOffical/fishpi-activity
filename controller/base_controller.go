package controller

import (
	"bless-activity/service/events"
	"time"

	"github.com/FishPiOffical/golang-sdk/sdk"
	"github.com/pocketbase/pocketbase/core"
)

type BaseController struct {
	event *core.ServeEvent
	app   core.App

	fishPiSdk *sdk.FishPiSDK
	eventbus  *events.Service
}

func NewBaseController(event *core.ServeEvent, eventbus *events.Service, fishPiSdk *sdk.FishPiSDK) *BaseController {
	controller := &BaseController{
		event: event,
		app:   event.App,

		fishPiSdk: fishPiSdk,
		eventbus:  eventbus,
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
