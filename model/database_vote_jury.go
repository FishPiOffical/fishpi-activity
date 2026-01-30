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

const (
	DbNameVoteJuryLogs         = "voteJuryLogs" // 评审团投票日志表
	VoteJuryLogFieldVoteId     = "voteId"       // 关联投票ID
	VoteJuryLogFieldFromUserId = "fromUserId"   // 评审团成员用户ID
	VoteJuryLogFieldToUserId   = "toUserId"     // 被投票用户ID
	VoteJuryLogFieldTimes      = "times"        // 投票次数
	VoteJuryLogFieldRound      = "round"        // 评审轮次
	VoteJuryLogFieldComment    = "comment"      // 投票备注

	DbNameVoteJuryResults       = "voteJuryResults" // 评审团结果表
	VoteJuryResultFieldVoteId   = "voteId"          // 关联投票ID
	VoteJuryResultFieldRound    = "round"           // 评审轮次
	VoteJuryResultFieldResults  = "results"         // 评审结果 JSON 用户ID与得票数映射
	VoteJuryResultFieldContinue = "continue"        // 是否进入下一轮评审
	VoteJuryResultFieldUserIds  = "userIds"         // 进入下一轮评审的用户ID列表
	VoteJuryResultFieldCreated  = "created"         // 创建时间
)

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
