package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type UserController struct {
	event *core.ServeEvent
	app   core.App

	logger *slog.Logger
}

func NewUserController(event *core.ServeEvent) *UserController {
	logger := event.App.Logger().With(
		slog.String("controller", "user"),
	)

	controller := &UserController{
		event:  event,
		app:    event.App,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *UserController) registerRoutes() {
	group := controller.event.Router.Group("/user")
	group.GET("/me", controller.GetMe).BindFunc(
		controller.CheckLogin,
	)
	// 后端登出，清除 token cookie 并重定向到首页
	group.GET("/logout", controller.Logout)
}

func (controller *UserController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

func (controller *UserController) CheckLogin(event *core.RequestEvent) error {
	if event.Auth == nil {
		return event.UnauthorizedError("未登录", nil)
	}
	if event.HasSuperuserAuth() {
		return event.ForbiddenError("请登录普通用户账号", nil)
	}
	return event.Next()
}

func (controller *UserController) GetMe(event *core.RequestEvent) error {
	//logger := controller.makeActionLogger("get_me")

	user := model.NewUser(event.Auth)

	return event.JSON(http.StatusOK, map[string]any{
		"id":       user.Id,
		"o_id":     user.OId(),
		"name":     user.Name(),
		"nickname": user.Nickname(),
		"avatar":   user.Avatar(),
	})
}

// Logout 清除 token cookie 并重定向到首页
func (controller *UserController) Logout(event *core.RequestEvent) error {
	// 设置过期 cookie 清除客户端 token（兼容HttpOnly和非HttpOnly情形）
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  (func() time.Time { t := time.Unix(0, 0); return t })(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	event.SetCookie(cookie)

	// 也尝试清除带点域名的 cookie（部分场景需要）
	// 注意: Go 的 http.SetCookie 无法直接带 domain 设置为 .example.com 进行删除
	// 如果需要，请在前端或部署环境中统一处理 domain

	return event.Redirect(http.StatusFound, "/")
}
