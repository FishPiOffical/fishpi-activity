//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameMedalOwners            = "medalOwners"  // 勋章用户表
	MedalOwnersFieldMedalId      = "medalId"      // 勋章ID
	MedalOwnersFieldUserId       = "userId"       // 用户ID
	MedalOwnersFieldDisplay      = "display"      // 是否展示
	MedalOwnersFieldDisplayOrder = "displayOrder" // 展示顺序
	MedalOwnersFieldData         = "data"         // 勋章数据
	MedalOwnersFieldExpired      = "expired"      // 过期时间
	MedalOwnersFieldCreated      = "created"      // 创建时间
	MedalOwnersFieldUpdated      = "updated"      // 更新时间
)

type MedalOwner struct {
	core.BaseRecordProxy
}

func NewMedalOwner(record *core.Record) *MedalOwner {
	owner := new(MedalOwner)
	owner.SetProxyRecord(record)
	return owner
}

func NewMedalOwnerFromCollection(collection *core.Collection) *MedalOwner {
	record := core.NewRecord(collection)
	return NewMedalOwner(record)
}

func (medal *MedalOwner) MedalId() string {
	return medal.GetString(MedalOwnersFieldMedalId)
}

func (medal *MedalOwner) SetMedalId(value string) {
	medal.Set(MedalOwnersFieldMedalId, value)
}

func (medal *MedalOwner) UserId() string {
	return medal.GetString(MedalOwnersFieldUserId)
}

func (medal *MedalOwner) SetUserId(value string) {
	medal.Set(MedalOwnersFieldUserId, value)
}

func (medal *MedalOwner) Display() bool {
	return medal.GetBool(MedalOwnersFieldDisplay)
}

func (medal *MedalOwner) SetDisplay(value bool) {
	medal.Set(MedalOwnersFieldDisplay, value)
}

func (medal *MedalOwner) DisplayOrder() int {
	return medal.GetInt(MedalOwnersFieldDisplayOrder)
}

func (medal *MedalOwner) SetDisplayOrder(value int) {
	medal.Set(MedalOwnersFieldDisplayOrder, value)
}

func (medal *MedalOwner) Data() string {
	return medal.GetString(MedalOwnersFieldData)
}

func (medal *MedalOwner) SetData(value string) {
	medal.Set(MedalOwnersFieldData, value)
}

func (medal *MedalOwner) Expired() types.DateTime {
	return medal.GetDateTime(MedalOwnersFieldExpired)
}

func (medal *MedalOwner) SetExpired(value types.DateTime) {
	medal.Set(MedalOwnersFieldExpired, value)
}

func (medal *MedalOwner) Created() types.DateTime {
	return medal.GetDateTime(MedalOwnersFieldCreated)
}

func (medal *MedalOwner) Updated() types.DateTime {
	return medal.GetDateTime(MedalOwnersFieldUpdated)
}
