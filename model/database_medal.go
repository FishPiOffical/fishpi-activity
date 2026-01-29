//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	types2 "github.com/FishPiOffical/golang-sdk/types"
	"github.com/pocketbase/pocketbase/core"
)

const (
	DbNameMedals           = "medals"      // 勋章表
	MedalsFieldOId         = "oId"         // 唯一ID 创建时间戳文本
	MedalsFieldMedalId     = "medalId"     // 唯一ID 递增数字本文
	MedalsFieldType        = "type"        // 类型 精良
	MedalsFieldName        = "name"        // 名称
	MedalsFieldDescription = "description" // 描述
	MedalsFieldAttr        = "attr"        // 属性
)

type Medal struct {
	core.BaseRecordProxy
}

func NewMedal(record *core.Record) *Medal {
	permission := new(Medal)
	permission.SetProxyRecord(record)
	return permission
}

func NewMedalFromCollection(collection *core.Collection) *Medal {
	record := core.NewRecord(collection)
	return NewMedal(record)
}

func (medal *Medal) OId() string {
	return medal.GetString(MedalsFieldOId)
}

func (medal *Medal) SetOId(value string) {
	medal.Set(MedalsFieldOId, value)
}

func (medal *Medal) MedalId() string {
	return medal.GetString(MedalsFieldMedalId)
}

func (medal *Medal) SetMedalId(value string) {
	medal.Set(MedalsFieldMedalId, value)
}

func (medal *Medal) Type() types2.MedalType {
	return types2.MedalType(medal.GetString(MedalsFieldType))
}

func (medal *Medal) SetType(value types2.MedalType) {
	medal.Set(MedalsFieldType, string(value))
}

func (medal *Medal) Name() string {
	return medal.GetString(MedalsFieldName)
}

func (medal *Medal) SetName(value string) {
	medal.Set(MedalsFieldName, value)
}

func (medal *Medal) Description() string {
	return medal.GetString(MedalsFieldDescription)
}

func (medal *Medal) SetDescription(value string) {
	medal.Set(MedalsFieldDescription, value)
}

func (medal *Medal) Attr() string {
	return medal.GetString(MedalsFieldAttr)
}

func (medal *Medal) SetAttr(value string) {
	medal.Set(MedalsFieldAttr, value)
}
