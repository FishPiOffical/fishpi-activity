package controller

import (
	"bless-activity/model"
	"bless-activity/service/events"
	"time"

	"github.com/FishPiOffical/golang-sdk/sdk"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
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

// RequireAdminRole 验证用户是否拥有管理员角色
// 此中间件会先验证用户是否已登录，然后检查用户的 role 字段是否为 admin
func RequireAdminRole() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: "require_admin_role",
		Func: func(event *core.RequestEvent) error {
			// 获取当前用户记录
			authRecord := event.Auth
			if authRecord == nil {
				return event.UnauthorizedError("未登录", nil)
			}

			// 检查用户角色
			role := authRecord.GetString(model.UsersFieldRole)
			if role != string(model.UserRoleAdmin) {
				return event.ForbiddenError("需要管理员权限", nil)
			}

			return event.Next()
		},
	}
}
