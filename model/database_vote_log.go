//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameVoteLogs          = "voteLogs"   // 投票日志表
	VoteLogsFieldVoteId     = "voteId"     // 关联投票ID
	VoteLogsFieldFromUserId = "fromUserId" // 投票用户ID
	VoteLogsFieldToUserId   = "toUserId"   // 被投票用户ID
	VoteLogsFieldComment    = "comment"    // 投票备注
	VoteLogsFieldValid      = "valid"      // 投票有效性
	VoteLogsFieldCreated    = "created"    // 创建时间
	VoteLogsFieldUpdated    = "updated"    // 更新时间
)

// VoteLogValid 投票日志有效性
/*
ENUM(
valid   // 有效
invalid // 无效
)
*/
type VoteLogValid string

// VoteLog wrapper type
type VoteLog struct {
	core.BaseRecordProxy
}

func NewVoteLog(record *core.Record) *VoteLog {
	voteLog := new(VoteLog)
	voteLog.SetProxyRecord(record)
	return voteLog
}

func NewVoteLogFromCollection(collection *core.Collection) *VoteLog {
	record := core.NewRecord(collection)
	return NewVoteLog(record)
}

func (voteLog *VoteLog) VoteId() string {
	return voteLog.GetString(VoteLogsFieldVoteId)
}

func (voteLog *VoteLog) SetVoteId(value string) {
	voteLog.Set(VoteLogsFieldVoteId, value)
}

func (voteLog *VoteLog) FromUserId() string {
	return voteLog.GetString(VoteLogsFieldFromUserId)
}

func (voteLog *VoteLog) SetFromUserId(value string) {
	voteLog.Set(VoteLogsFieldFromUserId, value)
}

func (voteLog *VoteLog) ToUserId() string {
	return voteLog.GetString(VoteLogsFieldToUserId)
}

func (voteLog *VoteLog) SetToUserId(value string) {
	voteLog.Set(VoteLogsFieldToUserId, value)
}

func (voteLog *VoteLog) Comment() string {
	return voteLog.GetString(VoteLogsFieldComment)
}

func (voteLog *VoteLog) SetComment(value string) {
	voteLog.Set(VoteLogsFieldComment, value)
}

func (voteLog *VoteLog) Valid() VoteLogValid {
	return MustParseVoteLogValid(voteLog.GetString(VoteLogsFieldValid))
}

func (voteLog *VoteLog) SetValid(value VoteLogValid) {
	voteLog.Set(VoteLogsFieldValid, value)
}

func (voteLog *VoteLog) Created() types.DateTime {
	return voteLog.GetDateTime(VoteLogsFieldCreated)
}

func (voteLog *VoteLog) Updated() types.DateTime {
	return voteLog.GetDateTime(VoteLogsFieldUpdated)
}
