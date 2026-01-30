//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	DbNameVoteJuryUsers      = "voteJuryUsers" // 评审团成员表
	VoteJuryUserFieldVoteId  = "voteId"        // 关联投票ID
	VoteJuryUserFieldUserId  = "userId"        // 评审团成员用户ID
	VoteJuryUserFieldStatus  = "status"        // 评审团成员状态 待审核、已通过、已拒绝
	VoteJuryUserFieldCreated = "created"       // 创建时间
	VoteJuryUserFieldUpdated = "updated"       // 更新时间
)

// VoteJuryUserStatus 评审团成员状态
/*
ENUM(
pending  // 待审核
approved // 已通过
rejected // 已拒绝
)
*/
type VoteJuryUserStatus string

// VoteJuryUser wrapper type
type VoteJuryUser struct {
	core.BaseRecordProxy
}

func NewVoteJuryUser(record *core.Record) *VoteJuryUser {
	user := new(VoteJuryUser)
	user.SetProxyRecord(record)
	return user
}

func NewVoteJuryUserFromCollection(collection *core.Collection) *VoteJuryUser {
	record := core.NewRecord(collection)
	return NewVoteJuryUser(record)
}

func (user *VoteJuryUser) VoteId() string {
	return user.GetString(VoteJuryUserFieldVoteId)
}

func (user *VoteJuryUser) SetVoteId(value string) {
	user.Set(VoteJuryUserFieldVoteId, value)
}

func (user *VoteJuryUser) UserId() string {
	return user.GetString(VoteJuryUserFieldUserId)
}

func (user *VoteJuryUser) SetUserId(value string) {
	user.Set(VoteJuryUserFieldUserId, value)
}

func (user *VoteJuryUser) Status() VoteJuryUserStatus {
	return MustParseVoteJuryUserStatus(user.GetString(VoteJuryUserFieldStatus))
}

func (user *VoteJuryUser) SetStatus(value VoteJuryUserStatus) {
	user.Set(VoteJuryUserFieldStatus, value)
}
