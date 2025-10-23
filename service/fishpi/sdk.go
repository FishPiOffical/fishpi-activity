package fishpi

import (
	"fmt"

	"github.com/duke-git/lancet/v2/xerror"
)

// GetChatroomNodeGet 获取节点列表
func (service *Service) GetChatroomNodeGet() (*GetChatroomNodeGetResponse, error) {
	res, err := service.client.NewRequest().
		SetQueryParam("apiKey", service.config.ApiKey).
		Get("/chat-room/node/get")
	if err != nil {
		return nil, err
	}
	response := new(GetChatroomNodeGetResponse)
	if err = res.Unmarshal(&response); err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, xerror.New("%s", response.Msg)
	}
	for _, node := range response.Avaliable {
		node.Node += fmt.Sprintf("?apiKey=%s", service.config.ApiKey)
	}
	return response, nil
}

// PostChatroomSend 发送聊天室消息
func (service *Service) PostChatroomSend(req *PostChatroomSendRequest) (*PostChatroomSendResponse, error) {
	req.ApiKey = service.config.ApiKey
	req.Client = "Golang/v0.0.3"

	res, err := service.client.NewRequest().
		SetBodyJsonMarshal(req).
		Post("/chat-room/send")
	if err != nil {
		return nil, err
	}
	response := new(PostChatroomSendResponse)
	if err = res.Unmarshal(&response); err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, xerror.New("%s", response.Msg)
	}
	return response, nil
}

// GetApiArticlesTag 获取帖子列表根据标签
func (service *Service) GetApiArticlesTag(tagName string, page int, size int) (*GetApiArticlesTagResponse, error) {
	page = max(page, 1)
	size = max(size, 1)
	res, err := service.client.NewRequest().
		SetQueryParam("apiKey", service.config.ApiKey).
		SetQueryParamsAnyType(map[string]any{
			"p":    page,
			"size": size,
		}).
		SetPathParam("tag", tagName).
		Get("/api/articles/tag/{tag}")
	if err != nil {
		return nil, err
	}

	response := new(GetApiArticlesTagResponse)
	if err = res.Unmarshal(&response); err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, xerror.New("%s", response.Msg)
	}

	return response, nil
}
