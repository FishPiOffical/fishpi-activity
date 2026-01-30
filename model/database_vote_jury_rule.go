//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameVoteJuryRules            = "voteJuryRules" // 评审团投票规则表
	VoteJuryRuleFieldVoteId        = "voteId"        // 关联投票ID
	VoteJuryRuleFieldCount         = "count"         // 评审团成员数量
	VoteJuryRuleFieldAdmins        = "admins"        // 评审团管理员用户ID列表
	VoteJuryRuleFieldDecisions     = "decisions"     // 评审团决策者用户ID列表
	VoteJuryRuleFieldStatus        = "status"        // 评审团状态 未开启、开放申请中、公示中、评审中、计票完成
	VoteJuryRuleFieldCurrentRound  = "currentRound"  // 当前轮次
	VoteJuryRuleFieldApplyTime     = "applyTime"     // 开放申请时间
	VoteJuryRuleFieldPublicityTime = "publicityTime" // 公示时间
	VoteJuryRuleFieldCreated       = "created"       // 创建时间
	VoteJuryRuleFieldUpdated       = "updated"       // 更新时间
)

// VoteJuryRuleStatus 评审团规则状态
/*
ENUM(
pending    // 未开启
applying   // 申请中
publicity  // 公示中
voting     // 评审中
completed  // 计票完成
)
*/
type VoteJuryRuleStatus string

// VoteJuryRule wrapper type
type VoteJuryRule struct {
	core.BaseRecordProxy
}

func NewVoteJuryRule(record *core.Record) *VoteJuryRule {
	rule := new(VoteJuryRule)
	rule.SetProxyRecord(record)
	return rule
}

func NewVoteJuryRuleFromCollection(collection *core.Collection) *VoteJuryRule {
	record := core.NewRecord(collection)
	return NewVoteJuryRule(record)
}

func (rule *VoteJuryRule) VoteId() string {
	return rule.GetString(VoteJuryRuleFieldVoteId)
}

func (rule *VoteJuryRule) SetVoteId(value string) {
	rule.Set(VoteJuryRuleFieldVoteId, value)
}

func (rule *VoteJuryRule) Count() int {
	return rule.GetInt(VoteJuryRuleFieldCount)
}

func (rule *VoteJuryRule) SetCount(value int) {
	rule.Set(VoteJuryRuleFieldCount, value)
}

func (rule *VoteJuryRule) Admins() []string {
	return rule.GetStringSlice(VoteJuryRuleFieldAdmins)
}

func (rule *VoteJuryRule) SetAdmins(value []string) {
	rule.Set(VoteJuryRuleFieldAdmins, value)
}

func (rule *VoteJuryRule) Decisions() []string {
	return rule.GetStringSlice(VoteJuryRuleFieldDecisions)
}

func (rule *VoteJuryRule) SetDecisions(value []string) {
	rule.Set(VoteJuryRuleFieldDecisions, value)
}

func (rule *VoteJuryRule) Status() VoteJuryRuleStatus {
	return MustParseVoteJuryRuleStatus(rule.GetString(VoteJuryRuleFieldStatus))
}

func (rule *VoteJuryRule) SetStatus(value VoteJuryRuleStatus) {
	rule.Set(VoteJuryRuleFieldStatus, value)
}

func (rule *VoteJuryRule) ApplyTime() types.DateTime {
	return rule.GetDateTime(VoteJuryRuleFieldApplyTime)
}

func (rule *VoteJuryRule) SetApplyTime(value types.DateTime) {
	rule.Set(VoteJuryRuleFieldApplyTime, value)
}

func (rule *VoteJuryRule) PublicityTime() types.DateTime {
	return rule.GetDateTime(VoteJuryRuleFieldPublicityTime)
}

func (rule *VoteJuryRule) SetPublicityTime(value types.DateTime) {
	rule.Set(VoteJuryRuleFieldPublicityTime, value)
}

func (rule *VoteJuryRule) CurrentRound() int {
	return rule.GetInt(VoteJuryRuleFieldCurrentRound)
}

func (rule *VoteJuryRule) SetCurrentRound(value int) {
	rule.Set(VoteJuryRuleFieldCurrentRound, value)
}
