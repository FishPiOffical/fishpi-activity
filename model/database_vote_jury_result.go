package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameVoteJuryResults       = "voteJuryResults" // 评审团结果表
	VoteJuryResultFieldVoteId   = "voteId"          // 关联投票ID
	VoteJuryResultFieldRound    = "round"           // 评审轮次
	VoteJuryResultFieldResults  = "results"         // 评审结果 JSON 用户ID与得票数映射
	VoteJuryResultFieldContinue = "continue"        // 是否进入下一轮评审
	VoteJuryResultFieldUserIds  = "userIds"         // 进入下一轮评审的用户ID列表
	VoteJuryResultFieldCreated  = "created"         // 创建时间
)

// VoteJuryResult wrapper type
type VoteJuryResult struct {
	core.BaseRecordProxy
}

func NewVoteJuryResult(record *core.Record) *VoteJuryResult {
	result := new(VoteJuryResult)
	result.SetProxyRecord(record)
	return result
}

func NewVoteJuryResultFromCollection(collection *core.Collection) *VoteJuryResult {
	record := core.NewRecord(collection)
	return NewVoteJuryResult(record)
}

func (result *VoteJuryResult) VoteId() string {
	return result.GetString(VoteJuryResultFieldVoteId)
}

func (result *VoteJuryResult) SetVoteId(value string) {
	result.Set(VoteJuryResultFieldVoteId, value)
}

func (result *VoteJuryResult) Round() int {
	return result.GetInt(VoteJuryResultFieldRound)
}

func (result *VoteJuryResult) SetRound(value int) {
	result.Set(VoteJuryResultFieldRound, value)
}

func (result *VoteJuryResult) Results() string {
	return result.GetString(VoteJuryResultFieldResults)
}

func (result *VoteJuryResult) SetResults(value string) {
	result.Set(VoteJuryResultFieldResults, value)
}

func (result *VoteJuryResult) Continue() bool {
	return result.GetBool(VoteJuryResultFieldContinue)
}

func (result *VoteJuryResult) SetContinue(value bool) {
	result.Set(VoteJuryResultFieldContinue, value)
}

func (result *VoteJuryResult) UserIds() []string {
	return result.GetStringSlice(VoteJuryResultFieldUserIds)
}

func (result *VoteJuryResult) SetUserIds(value []string) {
	result.Set(VoteJuryResultFieldUserIds, value)
}

func (result *VoteJuryResult) Created() types.DateTime {
	return result.GetDateTime(VoteJuryResultFieldCreated)
}
