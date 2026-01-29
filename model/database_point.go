//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNamePoints       = "points"  // 积分操作表
	PointsFieldGroup   = "group"   // 分组标识
	PointsFieldUserId  = "userId"  // 用户ID (关联users表)
	PointsFieldPoint   = "point"   // 积分数量
	PointsFieldStatus  = "status"  // 状态
	PointsFieldMemo    = "memo"    // 备注
	PointsFieldCreated = "created" // 创建时间
	PointsFieldUpdated = "updated" // 更新时间
)

// PointStatus 积分发放状态
/*
ENUM(
pending      // 待发放
distributing // 发放中
success      // 发放成功
failed       // 发放失败
)
*/
type PointStatus string

type Point struct {
	core.BaseRecordProxy
}

func NewPoint(record *core.Record) *Point {
	point := new(Point)
	point.SetProxyRecord(record)
	return point
}

func NewPointFromCollection(collection *core.Collection) *Point {
	record := core.NewRecord(collection)
	return NewPoint(record)
}

func (p *Point) Group() string {
	return p.GetString(PointsFieldGroup)
}

func (p *Point) SetGroup(value string) {
	p.Set(PointsFieldGroup, value)
}

func (p *Point) UserId() string {
	return p.GetString(PointsFieldUserId)
}

func (p *Point) SetUserId(value string) {
	p.Set(PointsFieldUserId, value)
}

func (p *Point) Point() int {
	return p.GetInt(PointsFieldPoint)
}

func (p *Point) SetPoint(value int) {
	p.Set(PointsFieldPoint, value)
}

func (p *Point) Status() PointStatus {
	return PointStatus(p.GetString(PointsFieldStatus))
}

func (p *Point) SetStatus(value PointStatus) {
	p.Set(PointsFieldStatus, string(value))
}

func (p *Point) Memo() string {
	return p.GetString(PointsFieldMemo)
}

func (p *Point) SetMemo(value string) {
	p.Set(PointsFieldMemo, value)
}

func (p *Point) Created() types.DateTime {
	return p.GetDateTime(PointsFieldCreated)
}

func (p *Point) Updated() types.DateTime {
	return p.GetDateTime(PointsFieldUpdated)
}
