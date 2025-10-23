package application

import (
	"bless-activity/controller"
	"bless-activity/service/fishpi"
	"log/slog"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

/*
	活动：年终征文徽章配色、周岁徽章关键词、周岁徽章设计（除名称）
	页面：nav: 活动名称 各活动跳转 登录/退出; 介绍：活动信息介绍; 每个活动内容
	后端：


    /login 登陆
    /callback 回调


*/

type Application struct {
	app *pocketbase.PocketBase

	fishPiService *fishpi.Service

	baseController   *controller.BaseController
	fishPiController *controller.FishPiController
	userController   *controller.UserController
}

func NewApp() *Application {
	app := pocketbase.New()

	application := &Application{
		app: app,
	}

	return application
}

func (application *Application) Start() error {

	// 初始化
	application.app.OnBootstrap().BindFunc(func(event *core.BootstrapEvent) error {

		if err := event.Next(); err != nil {
			return err
		}

		return application.init(event)
	})

	return application.app.Start()
}

func (application *Application) init(event *core.BootstrapEvent) error {
	event.App.Logger().Debug("初始化程序")

	var err error
	if application.fishPiService, err = fishpi.NewService(event.App); err != nil {
		event.App.Logger().Error("创建fishPi Service失败", slog.Any("err", err))
		return err
	}

	// 问题修复
	if err = application.fixBug(event); err != nil {
		return err
	}

	// 注册路由
	application.app.OnServe().BindFunc(application.registerRoutes)

	event.App.Logger().Debug("初始化完成")
	return nil
}

func (application *Application) registerRoutes(event *core.ServeEvent) error {

	event.Router.Bind(
		cookieToHeader,
		queryToHeader,
	)

	application.baseController = controller.NewBaseController(event)
	application.fishPiController = controller.NewFishPiController(event)
	application.userController = controller.NewUserController(event)

	event.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

	event.Router.GET("/status", func(e *core.RequestEvent) error {
		return e.String(http.StatusOK, "ok.")
	})

	return event.Next()
}
