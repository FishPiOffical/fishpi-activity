//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameUserTokens       = "user_tokens" // 用户令牌表
	UserTokensFieldUserId  = "user_id"     // 用户ID
	UserTokensFieldToken   = "token"       // 令牌
	UserTokensFieldState   = "state"       // 令牌状态
	UserTokenFieldExpired  = "expired"     // 过期时间
	UserTokensFieldCreated = "created"     // 创建时间
	UserTokensFieldUpdated = "updated"     // 更新时间
)

type UserToken struct {
	core.BaseRecordProxy
}

func NewUserToken(record *core.Record) *UserToken {
	userToken := new(UserToken)
	userToken.SetProxyRecord(record)
	return userToken
}

func NewUserTokenFromCollection(collection *core.Collection) *UserToken {
	record := core.NewRecord(collection)
	return NewUserToken(record)
}

func (userToken *UserToken) GetUserId() string {
	return userToken.GetString(UserTokensFieldUserId)
}

func (userToken *UserToken) SetUserId(value string) {
	userToken.Set(UserTokensFieldUserId, value)
}

func (userToken *UserToken) GetToken() string {
	return userToken.GetString(UserTokensFieldToken)
}

func (userToken *UserToken) SetToken(value string) {
	userToken.Set(UserTokensFieldToken, value)
}

// UserTokenState 用户令牌状态
/*
ENUM(
unverified // 未验证
verified // 已验证
revoked // 已撤销
)
*/
type UserTokenState string

func (userToken *UserToken) GetState() UserTokenState {
	return UserTokenState(userToken.GetString(UserTokensFieldState))
}

func (userToken *UserToken) SetState(value UserTokenState) {
	userToken.Set(UserTokensFieldState, string(value))
}

func (userToken *UserToken) GetExpired() types.DateTime {
	return userToken.GetDateTime(UserTokenFieldExpired)
}

func (userToken *UserToken) SetExpired(value types.DateTime) {
	userToken.Set(UserTokenFieldExpired, value)
}

func (userToken *UserToken) GetCreated() types.DateTime {
	return userToken.GetDateTime(UserTokensFieldCreated)
}

func (userToken *UserToken) GetUpdated() types.DateTime {
	return userToken.GetDateTime(UserTokensFieldUpdated)
}
