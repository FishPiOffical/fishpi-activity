//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	DbNameVoteJuryApplyLogs      = "voteJuryApplyLogs" // 评审团申请日志表
	VoteJuryApplyLogFieldVoteId  = "voteId"            // 关联投票ID
	VoteJuryApplyLogFieldUserId  = "userId"            // 申请用户ID
	VoteJuryApplyLogFieldReason  = "reason"            // 申请理由
	VoteJuryApplyLogFieldStatus  = "status"            // 申请状态 待审核、已通过、已拒绝
	VoteJuryApplyLogFieldAdminId = "adminId"           // 审核管理员用户ID
	VoteJuryApplyLogFieldCreated = "created"           // 创建时间
	VoteJuryApplyLogFieldUpdated = "updated"           // 更新时间
)

// VoteJuryApplyLogStatus 评审团申请日志状态
/*
ENUM(
pending  // 待审核
approved // 已通过
rejected // 已拒绝
)
*/
type VoteJuryApplyLogStatus string

// VoteJuryApplyLog wrapper type
type VoteJuryApplyLog struct {
	core.BaseRecordProxy
}

func NewVoteJuryApplyLog(record *core.Record) *VoteJuryApplyLog {
	log := new(VoteJuryApplyLog)
	log.SetProxyRecord(record)
	return log
}

func NewVoteJuryApplyLogFromCollection(collection *core.Collection) *VoteJuryApplyLog {
	record := core.NewRecord(collection)
	return NewVoteJuryApplyLog(record)
}

func (log *VoteJuryApplyLog) VoteId() string {
	return log.GetString(VoteJuryApplyLogFieldVoteId)
}

func (log *VoteJuryApplyLog) SetVoteId(value string) {
	log.Set(VoteJuryApplyLogFieldVoteId, value)
}

func (log *VoteJuryApplyLog) UserId() string {
	return log.GetString(VoteJuryApplyLogFieldUserId)
}

func (log *VoteJuryApplyLog) SetUserId(value string) {
	log.Set(VoteJuryApplyLogFieldUserId, value)
}

func (log *VoteJuryApplyLog) Reason() string {
	return log.GetString(VoteJuryApplyLogFieldReason)
}

func (log *VoteJuryApplyLog) SetReason(value string) {
	log.Set(VoteJuryApplyLogFieldReason, value)
}

func (log *VoteJuryApplyLog) Status() VoteJuryApplyLogStatus {
	return MustParseVoteJuryApplyLogStatus(log.GetString(VoteJuryApplyLogFieldStatus))
}

func (log *VoteJuryApplyLog) SetStatus(value VoteJuryApplyLogStatus) {
	log.Set(VoteJuryApplyLogFieldStatus, value)
}

func (log *VoteJuryApplyLog) AdminId() string {
	return log.GetString(VoteJuryApplyLogFieldAdminId)
}

func (log *VoteJuryApplyLog) SetAdminId(value string) {
	log.Set(VoteJuryApplyLogFieldAdminId, value)
}
