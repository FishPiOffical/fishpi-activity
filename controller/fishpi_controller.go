package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"bless-activity/model"

	types2 "github.com/FishPiOffical/golang-sdk/types"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	ctxFishpiLoginUser = "login_user"
	ctxFishpiNext      = "next"
	ctxFishpiOpenId    = "openid"
	ctxFishpiUserInfo  = "fishpi_user_info"
	ctxFishpiClientId  = "client_id"
)

type FishPiController struct {
	*BaseController

	group *router.RouterGroup[*core.RequestEvent]

	logger *slog.Logger
}

func NewFishPiController(base *BaseController, group *router.RouterGroup[*core.RequestEvent]) *FishPiController {
	logger := base.app.Logger().WithGroup("controller.fishpi")

	controller := &FishPiController{
		BaseController: base,

		group:  group,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *FishPiController) registerRoutes() {

	fishpiGroup := controller.group.Group("/fishpi")

	fishpiGroup.GET("/login", controller.Login)
	fishpiGroup.GET("/callback", controller.Callback).BindFunc(
		controller.CallbackVerify,
	)
	fishpiGroup.GET("/verify", controller.Verify)
	fishpiGroup.GET("/redirect", controller.Redirect)

}

func (controller *FishPiController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

func (controller *FishPiController) Login(event *core.RequestEvent) error {

	appUrl := event.App.Settings().Meta.AppURL

	// 获取原始页面URL（redirect 或 next 参数）
	redirectUrl := event.Request.URL.Query().Get("redirect")
	if redirectUrl == "" {
		redirectUrl = event.Request.URL.Query().Get("next")
	}
	if redirectUrl == "" {
		redirectUrl = event.Request.Header.Get("Referer")
	}

	// 将redirect参数添加到callback URL中
	callbackUrl := fmt.Sprintf("%s/backend/fishpi/callback", appUrl)
	callbackParams := url.Values{}
	if redirectUrl != "" && redirectUrl != "/" && redirectUrl != appUrl && redirectUrl != appUrl+"/" {
		callbackParams.Set("redirect", redirectUrl)
	}
	if clientId := event.Request.URL.Query().Get("client_id"); clientId != "" {
		callbackParams.Set("client_id", clientId)
	}
	if callbackParams.Encode() != "" {
		callbackUrl = fmt.Sprintf("%s?%s", callbackUrl, callbackParams.Encode())
	}

	link := controller.fishPiSdk.GetOpenIdUrl(appUrl, callbackUrl)

	return event.Redirect(http.StatusFound, link)
}

func (controller *FishPiController) CallbackVerify(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("callback verify").With(
		slog.String("path", event.Request.URL.String()),
	)

	info, err := event.RequestInfo()
	if err != nil {
		logger.Error("获取请求信息失败", slog.Any("err", err))
		return err
	}

	query := info.Query

	var openIdPtr *string
	if openIdPtr, err = controller.fishPiSdk.PostOpenIdVerify(query); err != nil {
		logger.Error("发起验证请求失败", slog.Any("err", err))
		return err
	}

	openid := convertor.ToString(openIdPtr)

	resp := new(types2.ApiResponse[*types2.GetUserInfoByIdData])
	if resp, err = controller.fishPiSdk.GetUserInfoById(openid); err != nil {
		logger.Error("发起获取用户信息请求失败", slog.Any("err", err))
		return err
	}
	if resp.Code != 0 {
		logger.Error("获取用户信息失败", slog.String("resp", fmt.Sprintf("%+v", resp)))
		return errors.New(resp.Msg)
	}

	// 保存redirect参数
	redirectUrl := event.Request.URL.Query().Get("redirect")
	if redirectUrl != "" {
		event.Set("redirect_url", redirectUrl)
	}
	if clientId := event.Request.URL.Query().Get("client_id"); clientId != "" {
		event.Set(ctxFishpiClientId, clientId)
	}

	user := new(model.User)
	if err = event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{model.UsersFieldOId: openid}).One(user); err == nil {
		event.Set(ctxFishpiLoginUser, user)
		event.Set(ctxFishpiUserInfo, resp.Data)
		event.Set(ctxFishpiNext, "login")

		return event.Next()
	} else if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("查询用户信息失败", slog.Any("err", err))
		return err
	}

	event.Set(ctxFishpiOpenId, openid)
	event.Set(ctxFishpiUserInfo, resp.Data)
	event.Set(ctxFishpiNext, "register")

	return event.Next()
}

type FishpiUserInfoResult struct {
	Msg  string          `json:"msg"`
	Code int             `json:"code"`
	Data *FishpiUserInfo `json:"data"`
}

type FishpiUserInfo struct {
	UserAvatarURL string `json:"userAvatarURL"`
	UserNickname  string `json:"userNickname"`
	UserName      string `json:"userName"`
}

func (fishpiUserInfo *FishpiUserInfo) Name() string {
	return fmt.Sprintf("(%s)%s", fishpiUserInfo.UserName, fishpiUserInfo.UserNickname)
}

func (controller *FishPiController) Callback(event *core.RequestEvent) error {
	if event.Get(ctxFishpiNext) == "login" {
		return controller.login(event)
	}
	return controller.register(event)
}

func (controller *FishPiController) login(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("callback login").With(
		slog.String("path", event.Request.URL.String()),
	)
	user := event.Get(ctxFishpiLoginUser).(*model.User)
	fishpiUserInfo := event.Get(ctxFishpiUserInfo).(*types2.GetUserInfoByIdData)

	logger = logger.With(slog.String("id", user.Id), slog.String("name", user.GetString("name")))

	// 更新用户资料
	if fishpiUserInfo.UserName != user.Name() || fishpiUserInfo.UserNickname != user.Nickname() || fishpiUserInfo.UserAvatarURL != user.Avatar() {
		if fishpiUserInfo.UserName != user.Name() {
			user.SetName(fishpiUserInfo.UserName)
		}
		if fishpiUserInfo.UserNickname != user.Nickname() {
			user.SetNickname(fishpiUserInfo.UserNickname)
		}
		if fishpiUserInfo.UserAvatarURL != user.Avatar() {
			user.SetAvatar(fishpiUserInfo.UserAvatarURL)
		}

		if err := event.App.Save(user); err != nil {
			logger.Error("更新用户资料失败", slog.Any("user", user), slog.Any("fishpi_user_info", fishpiUserInfo), slog.Any("err", err))
			return err
		}
	}

	//token, err := user.NewAuthToken()
	//if err != nil {
	//	logger.Error("生成token失败", slog.Any("err", err))
	//	return err
	//}

	// 创建userToken记录
	userTokenCollection, err := event.App.FindCollectionByNameOrId(model.DbNameUserTokens)
	if err != nil {
		logger.Error("查找user_token集合失败", slog.Any("user", user), slog.Any("err", err))
		return err
	}
	userToken := model.NewUserTokenFromCollection(userTokenCollection)
	userToken.SetUserId(user.Id)
	//userToken.SetToken("") // 请求接口时再生成token
	userToken.SetState(model.UserTokenStateUnverified)
	userToken.SetExpired(types.NowDateTime().Add(time.Minute))
	if err = event.App.Save(userToken); err != nil {
		logger.Error("创建user_token记录失败", slog.Any("user", user), slog.Any("err", err))
		return err
	}

	// 获取重定向URL
	var values url.Values
	values.Add("id", userToken.Id)
	if redirect := event.Get("redirect_url"); redirect != nil {
		if redirectStr, ok := redirect.(string); ok && redirectStr != "" {
			values.Add("redirect", redirectStr)
		}
	}
	redirectUrl := "/redirect?" + values.Encode()

	// 发送客户端登录通知

	return event.Redirect(http.StatusFound, redirectUrl)
}

func (controller *FishPiController) register(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("callback register").With(
		slog.String("path", event.Request.URL.String()),
	)

	openid := event.Get(ctxFishpiOpenId).(string)
	fishpiUserInfo := event.Get(ctxFishpiUserInfo).(*types2.GetUserInfoByIdData)

	// 创建用户
	userCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameUsers)
	if err != nil {
		return err
	}
	user := model.NewUserFromCollection(userCollection)
	user.SetEmail(fmt.Sprintf("%s@fishpi.cn", openid))
	user.SetEmailVisibility(true)
	user.SetVerified(true)
	user.SetOId(openid)
	user.SetName(fishpiUserInfo.UserName)
	user.SetNickname(fishpiUserInfo.UserNickname)
	user.SetAvatar(fishpiUserInfo.UserAvatarURL)
	user.SetRandomPassword()
	if err = controller.app.Save(user); err != nil {
		return err
	}

	// 创建userToken记录
	userTokenCollection, err := event.App.FindCollectionByNameOrId(model.DbNameUserTokens)
	if err != nil {
		logger.Error("查找user_token集合失败", slog.Any("user", user), slog.Any("err", err))
		return err
	}
	userToken := model.NewUserTokenFromCollection(userTokenCollection)
	userToken.SetUserId(user.Id)
	//userToken.SetToken("") // 请求接口时再生成token
	userToken.SetState(model.UserTokenStateUnverified)
	userToken.SetExpired(types.NowDateTime().Add(time.Minute))
	if err = event.App.Save(userToken); err != nil {
		logger.Error("创建user_token记录失败", slog.Any("user", user), slog.Any("err", err))
		return err
	}

	// 获取重定向URL
	var values url.Values
	values.Add("id", userToken.Id)
	if redirect := event.Get("redirect_url"); redirect != nil {
		if redirectStr, ok := redirect.(string); ok && redirectStr != "" {
			values.Add("redirect", redirectStr)
		}
	}
	redirectUrl := "/redirect?" + values.Encode()

	// 发送客户端登录通知

	return event.Redirect(http.StatusFound, redirectUrl)
}

func (controller *FishPiController) Verify(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("verify").With(
		slog.String("path", event.Request.URL.String()),
	)

	id := event.Request.URL.Query().Get("id")
	if id == "" {
		logger.Error("缺少参数id")
		return event.BadRequestError("请求异常", nil)
	}

	userToken := new(model.UserToken)
	if err := event.App.RecordQuery(model.DbNameUserTokens).Where(dbx.HashExp{model.CommonFieldId: id}).One(userToken); err != nil {
		logger.Error("查询user_token记录失败", slog.String("id", id), slog.Any("err", err))
		return event.BadRequestError("参数错误", nil)
	}
	if userToken.GetExpired().After(types.NowDateTime()) {
		logger.Error("user_token已过期", slog.String("id", id), slog.Any("expired", userToken.GetExpired()))
		return event.BadRequestError("链接已过期，请重新登录", nil)
	}
	if userToken.GetState() != model.UserTokenStateUnverified {
		logger.Error("user_token状态异常", slog.String("id", id), slog.String("state", string(userToken.GetState())))
		return event.BadRequestError("请勿重复操作", nil)
	}

	user := model.NewUser(event.Auth)
	token, err := user.NewAuthToken()
	if err != nil {
		logger.Error("生成token失败", slog.Any("err", err))
		return event.InternalServerError("服务器错误", nil)
	}
	userToken.SetToken(token)
	userToken.SetState(model.UsersFieldVerified)

	if err = event.App.Save(userToken); err != nil {
		logger.Error("更新user_token记录失败", slog.String("id", id), slog.Any("err", err))
		return event.InternalServerError("服务器错误", nil)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (controller *FishPiController) Redirect(event *core.RequestEvent) error {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>登录成功</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .container {
            text-align: center;
            background: white;
            padding: 60px 40px;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            max-width: 500px;
            width: 90%;
        }
        .success-icon {
            font-size: 80px;
            color: #52c41a;
            margin-bottom: 20px;
        }
        h1 {
            font-size: 28px;
            color: #333;
            margin-bottom: 15px;
        }
        p {
            font-size: 16px;
            color: #666;
            margin-bottom: 35px;
            line-height: 1.6;
        }
        .close-btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 15px 40px;
            font-size: 16px;
            border-radius: 50px;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }
        .close-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
        }
        .close-btn:active {
            transform: translateY(0);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">✓</div>
        <h1>您已成功登录</h1>
        <p>点击下方按钮可关闭此界面</p>
        <button class="close-btn" onclick="closeWindow()">关闭窗口</button>
    </div>
    <script>
        function closeWindow() {
            // 尝试关闭窗口
            window.close();
            
            // 如果无法关闭，提示用户
            setTimeout(function() {
                if (!window.closed) {
                    alert('请手动关闭此标签页');
                }
            }, 100);
        }
    </script>
</body>
</html>`
	return event.HTML(http.StatusOK, html)
}
