package fishpi

import (
	"fmt"
	"log/slog"
)

func (service *Service) GetInfo(openid string) (*UserInfo, error) {
	logger := service.logger.With(
		slog.String("service_action", "获取用户信息"),
		slog.String("openid", openid),
	)

	result := new(UserInfoResult)
	resp, err := service.client.NewRequest().
		SetSuccessResult(result).
		SetQueryParam("userId", openid).
		Get("/api/user/getInfoById")
	if err != nil {
		logger.Error("请求失败", slog.Any("err", err))
		return nil, err
	}
	if resp.IsErrorState() {
		logger.Error("请求状态码异常", slog.String("resp", resp.String()))
		return nil, fmt.Errorf("status:%d", resp.GetStatusCode())
	}
	if result.Code != 0 {
		logger.Error("状态码异常", slog.String("resp", resp.String()))
		return nil, fmt.Errorf("code:%d,message:%s", result.Code, result.Msg)
	}
	return result.Data, nil
}

func (service *Service) EditPoint(req *EditPointReq) (*EditPointReply, error) {
	logger := service.logger.With(
		slog.String("service_action", "编辑积分"),
		slog.Any("req", req),
	)

	result := new(EditPointReply)
	resp, err := service.client.NewRequest().
		SetBodyJsonMarshal(map[string]any{
			"goldFingerKey": service.config.GoldFingerKey,
			"userName":      req.UserName,
			"point":         req.Point,
			"memo":          req.Memo,
		}).
		SetSuccessResult(result).
		Post(`/user/edit/points`)
	if err != nil {
		logger.Error("请求失败", slog.Any("err", err))
		return nil, err
	}
	if resp.IsErrorState() {
		logger.Error("请求状态码异常", slog.String("resp", resp.String()))
		return nil, fmt.Errorf("status:%d", resp.GetStatusCode())
	}
	if result.Code != 0 {
		logger.Error("状态码异常", slog.String("resp", resp.String()))
		return nil, fmt.Errorf("code:%d,message:%s", result.Code, result.Msg)
	}
	return result, nil
}

func (service *Service) GetUser(username string) (*GetUserReply, error) {
	logger := service.logger.With(
		slog.String("service_action", "get_user"),
		slog.String("username", username),
	)
	result := new(GetUserReply)
	resp, err := service.client.NewRequest().
		SetPathParam("username", username).
		SetSuccessResult(result).
		Get("/user/{username}")
	if err != nil {
		logger.Error("请求失败", slog.Any("err", err))
		return nil, err
	}
	if resp.IsErrorState() {
		logger.Error("请求状态码异常", slog.String("resp", resp.String()))
		return nil, fmt.Errorf("status:%d", resp.GetStatusCode())
	}

	return result, nil
}

// Distribute 发放积分的便捷方法
func (service *Service) Distribute(username string, point int, memo string) error {
	req := &EditPointReq{
		UserName: username,
		Point:    point,
		Memo:     memo,
	}
	_, err := service.EditPoint(req)
	return err
}
