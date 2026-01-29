//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameUsers               = "users"
	UsersFieldEmail           = "email"
	UsersFieldEmailVisibility = "emailVisibility"
	UsersFieldVerified        = "verified"
	UsersFieldName            = "name"
	UsersFieldNickname        = "nickname"
	UsersFieldAvatar          = "avatar"
	UsersFieldOId             = "oId"
	UsersFieldRole            = "role"
	UsersFieldCreated         = "created"
	UsersFieldUpdated         = "updated"
)

type User struct {
	core.BaseRecordProxy
}

func NewUser(record *core.Record) *User {
	user := new(User)
	user.SetProxyRecord(record)
	return user
}

func NewUserFromCollection(collection *core.Collection) *User {
	record := core.NewRecord(collection)
	return NewUser(record)
}

func (user *User) Name() string {
	return user.GetString(UsersFieldName)
}

func (user *User) SetName(value string) {
	user.Set(UsersFieldName, value)
}

func (user *User) Nickname() string {
	return user.GetString(UsersFieldNickname)
}

func (user *User) SetNickname(value string) {
	user.Set(UsersFieldNickname, value)
}

func (user *User) Avatar() string {
	return user.GetString(UsersFieldAvatar)
}

func (user *User) SetAvatar(value string) {
	user.Set(UsersFieldAvatar, value)
}

func (user *User) OId() string {
	return user.GetString(UsersFieldOId)
}

func (user *User) SetOId(value string) {
	user.Set(UsersFieldOId, value)
}

// UserRole 用户角色
/*
ENUM(
admin // 管理员
user  // 普通用户
)
*/
type UserRole string

func (user *User) Role() UserRole {
	return UserRole(user.GetString(UsersFieldRole))
}

func (user *User) SetRole(value UserRole) {
	user.Set(UsersFieldRole, string(value))
}

func (user *User) Created() types.DateTime {
	return user.GetDateTime(UsersFieldCreated)
}

func (user *User) Updated() types.DateTime {
	return user.GetDateTime(UsersFieldUpdated)
}

func (user *User) RegisteredAt() types.DateTime {
	oId := user.OId()
	if oId == "" {
		return types.DateTime{}
	}

	// Parse oId as milliseconds timestamp
	timestamp, err := strconv.ParseInt(oId, 10, 64)
	if err != nil {
		return types.DateTime{}
	}

	// Convert milliseconds to time.Time
	t := time.UnixMilli(timestamp)

	// Convert to types.DateTime

	dt := types.DateTime{}
	_ = dt.Scan(t)
	return dt
}
