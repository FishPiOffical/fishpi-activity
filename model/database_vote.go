//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameVotes                = "votes"            // 投票表
	VotesFieldName             = "name"             // 投票名称
	VotesFieldDesc             = "desc"             // 投票描述
	VotesFieldType             = "type"             // 投票类型 普通投票、评审团投票
	VotesFieldTimes            = "times"            // 可投票次数
	VotesFieldRepeat           = "repeat"           // 是否允许重复投票
	VotesFieldUserRegisterDays = "userRegisterDays" // 用户注册天数限制
	VotesFieldStart            = "start"            // 开始时间
	VotesFieldEnd              = "end"              // 结束时间
)

// VoteType 投票类型
/*
ENUM(
normal // 普通投票
jury   // 评审团投票
)
*/
type VoteType string

// Vote wrapper type
type Vote struct {
	core.BaseRecordProxy
}

func NewVote(record *core.Record) *Vote {
	vote := new(Vote)
	vote.SetProxyRecord(record)
	return vote
}

func NewVoteFromCollection(collection *core.Collection) *Vote {
	record := core.NewRecord(collection)
	return NewVote(record)
}

func (vote *Vote) Name() string {
	return vote.GetString(VotesFieldName)
}

func (vote *Vote) SetName(value string) {
	vote.Set(VotesFieldName, value)
}

func (vote *Vote) Desc() string {
	return vote.GetString(VotesFieldDesc)
}

func (vote *Vote) SetDesc(value string) {
	vote.Set(VotesFieldDesc, value)
}

func (vote *Vote) Type() VoteType {
	typeStr := vote.GetString(VotesFieldType)
	if typeStr == "" {
		return VoteTypeNormal
	}
	return MustParseVoteType(typeStr)
}

func (vote *Vote) SetType(value VoteType) {
	vote.Set(VotesFieldType, value)
}

func (vote *Vote) Times() int {
	return vote.GetInt(VotesFieldTimes)
}

func (vote *Vote) SetTimes(value int) {
	vote.Set(VotesFieldTimes, value)
}

func (vote *Vote) Repeat() bool {
	return vote.GetBool(VotesFieldRepeat)
}

func (vote *Vote) SetRepeat(value bool) {
	vote.Set(VotesFieldRepeat, value)
}

func (vote *Vote) UserRegisterDays() int {
	return vote.GetInt(VotesFieldUserRegisterDays)
}

func (vote *Vote) SetUserRegisterDays(value int) {
	vote.Set(VotesFieldUserRegisterDays, value)
}

func (vote *Vote) Start() types.DateTime {
	return vote.GetDateTime(VotesFieldStart)
}

func (vote *Vote) SetStart(value types.DateTime) {
	vote.Set(VotesFieldStart, value)
}

func (vote *Vote) End() types.DateTime {
	return vote.GetDateTime(VotesFieldEnd)
}

func (vote *Vote) SetEnd(value types.DateTime) {
	vote.Set(VotesFieldEnd, value)
}

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
