package application

import (
	"bless-activity/controller"
	"bless-activity/pkg/fishpi_sdk"
	"bless-activity/service/events"
	"bless-activity/service/fetch_article"
	"log/slog"
	"net/http"
	"os"

	"github.com/FishPiOffical/golang-sdk/sdk"
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

	fishPiSdk *sdk.FishPiSDK

	fetchArticleService *fetch_article.Service

	baseController               *controller.BaseController
	fishPiController             *controller.FishPiController
	userController               *controller.UserController
	activityController           *controller.ActivityController
	shieldFiveYearController     *controller.ShieldFiveYearController
	rewardDistributionController *controller.RewardDistributionController
	voteJuryController           *controller.VoteJuryController
	medalController              *controller.MedalController
	pointController              *controller.PointController

	eventbus *events.Service
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
	var provider *fishpi_sdk.Provider

	if provider, err = fishpi_sdk.NewProvider(event.App); err != nil {
		event.App.Logger().Error("创建fishPi SDK Provider失败", slog.Any("err", err))
		return err
	}

	application.fishPiSdk = sdk.NewSDK(
		provider,
		sdk.WithLogDir("_tmp/logs/"),
	)

	application.fetchArticleService = fetch_article.NewService(application.app, application.fishPiSdk)
	if !application.app.IsDev() {
		if err = application.fetchArticleService.Run(); err != nil {
			event.App.Logger().Error("启动文章爬取服务失败", slog.Any("err", err))
			return err
		}
	}

	// 问题修复
	if err = application.fixBug(event); err != nil {
		return err
	}

	// 注册钩子
	application.registerHooks()

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

	application.eventbus = events.NewService(event.App)

	// 调整
	application.baseController = controller.NewBaseController(event, application.eventbus, application.fishPiSdk)

	backendGroup := event.Router.Group("/backend")

	// fishpi 鱼排相关
	application.fishPiController = controller.NewFishPiController(application.baseController, backendGroup)

	// 评审团投票
	application.voteJuryController = controller.NewVoteJuryController(application.baseController, backendGroup)

	// 勋章管理
	application.medalController = controller.NewMedalController(event, backendGroup, application.baseController)

	// 积分管理
	application.pointController = controller.NewPointController(event, backendGroup, application.baseController)

	// 待定
	application.userController = controller.NewUserController(event)
	application.activityController = controller.NewActivityController(event)
	application.shieldFiveYearController = controller.NewShieldFiveYearController(event)
	application.rewardDistributionController = controller.NewRewardDistributionController(event, application.baseController)

	event.Router.GET("/status", func(e *core.RequestEvent) error {
		return e.String(http.StatusOK, "ok.")
	})

	event.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), true))

	return event.Next()
}
