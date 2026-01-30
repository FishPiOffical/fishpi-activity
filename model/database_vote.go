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
