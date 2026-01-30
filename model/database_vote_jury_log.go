package model

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	DbNameVoteJuryLogs         = "voteJuryLogs" // 评审团投票日志表
	VoteJuryLogFieldVoteId     = "voteId"       // 关联投票ID
	VoteJuryLogFieldFromUserId = "fromUserId"   // 评审团成员用户ID
	VoteJuryLogFieldToUserId   = "toUserId"     // 被投票用户ID
	VoteJuryLogFieldTimes      = "times"        // 投票次数
	VoteJuryLogFieldRound      = "round"        // 评审轮次
	VoteJuryLogFieldComment    = "comment"      // 投票备注
)

// VoteJuryLog wrapper type
type VoteJuryLog struct {
	core.BaseRecordProxy
}

func NewVoteJuryLog(record *core.Record) *VoteJuryLog {
	log := new(VoteJuryLog)
	log.SetProxyRecord(record)
	return log
}

func NewVoteJuryLogFromCollection(collection *core.Collection) *VoteJuryLog {
	record := core.NewRecord(collection)
	return NewVoteJuryLog(record)
}

func (log *VoteJuryLog) VoteId() string {
	return log.GetString(VoteJuryLogFieldVoteId)
}

func (log *VoteJuryLog) SetVoteId(value string) {
	log.Set(VoteJuryLogFieldVoteId, value)
}

func (log *VoteJuryLog) FromUserId() string {
	return log.GetString(VoteJuryLogFieldFromUserId)
}

func (log *VoteJuryLog) SetFromUserId(value string) {
	log.Set(VoteJuryLogFieldFromUserId, value)
}

func (log *VoteJuryLog) ToUserId() string {
	return log.GetString(VoteJuryLogFieldToUserId)
}

func (log *VoteJuryLog) SetToUserId(value string) {
	log.Set(VoteJuryLogFieldToUserId, value)
}

func (log *VoteJuryLog) Times() int {
	return log.GetInt(VoteJuryLogFieldTimes)
}

func (log *VoteJuryLog) SetTimes(value int) {
	log.Set(VoteJuryLogFieldTimes, value)
}

func (log *VoteJuryLog) Round() int {
	return log.GetInt(VoteJuryLogFieldRound)
}

func (log *VoteJuryLog) SetRound(value int) {
	log.Set(VoteJuryLogFieldRound, value)
}

func (log *VoteJuryLog) Comment() string {
	return log.GetString(VoteJuryLogFieldComment)
}

func (log *VoteJuryLog) SetComment(value string) {
	log.Set(VoteJuryLogFieldComment, value)
}
