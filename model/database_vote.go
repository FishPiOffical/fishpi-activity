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
	VoteJuryRuleCreatedAt          = "createdAt"     // 创建时间
	VoteJuryRuleUpdatedAt          = "updatedAt"     // 更新时间

	DbNameVoteJuryUsers     = "voteJuryUsers" // 评审团成员表
	VoteJuryUserFieldVoteId = "voteId"        // 关联投票ID
	VoteJuryUserFieldUserId = "userId"        // 评审团成员用户ID
	VoteJuryUserFieldStatus = "status"        // 评审团成员状态 待审核、已通过、已拒绝
	VoteJuryUserCreatedAt   = "createdAt"     // 创建时间
	VoteJuryUserUpdatedAt   = "updatedAt"     // 更新时间

	DbNameVoteJuryApplyLogs      = "voteJuryApplyLogs" // 评审团申请日志表
	VoteJuryApplyLogFieldVoteId  = "voteId"            // 关联投票ID
	VoteJuryApplyLogFieldUserId  = "userId"            // 申请用户ID
	VoteJuryApplyLogFieldReason  = "reason"            // 申请理由
	VoteJuryApplyLogFieldStatus  = "status"            // 申请状态 待审核、已通过、已拒绝
	VoteJuryApplyLogFieldAdminId = "adminId"           // 审核管理员用户ID
	VoteJuryApplyLogCreatedAt    = "createdAt"         // 创建时间
	VoteJuryApplyLogUpdatedAt    = "updatedAt"         // 更新时间

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

// JuryStatus 评审团状态
/*
ENUM(
pending    // 未开启
applying   // 申请中
publicity  // 公示中
voting     // 评审中
completed  // 计票完成
)
*/
type JuryStatus string

// JuryUserStatus 评审团成员状态
/*
ENUM(
pending  // 待审核
approved // 已通过
rejected // 已拒绝
)
*/
type JuryUserStatus string

// JuryApplyStatus 评审团申请状态
/*
ENUM(
pending  // 待审核
approved // 已通过
rejected // 已拒绝
)
*/
type JuryApplyStatus string

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

func (rule *VoteJuryRule) Status() JuryStatus {
	return MustParseJuryStatus(rule.GetString(VoteJuryRuleFieldStatus))
}

func (rule *VoteJuryRule) SetStatus(value JuryStatus) {
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

func (user *VoteJuryUser) Status() JuryUserStatus {
	return MustParseJuryUserStatus(user.GetString(VoteJuryUserFieldStatus))
}

func (user *VoteJuryUser) SetStatus(value JuryUserStatus) {
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

func (log *VoteJuryApplyLog) Status() JuryApplyStatus {
	return MustParseJuryApplyStatus(log.GetString(VoteJuryApplyLogFieldStatus))
}

func (log *VoteJuryApplyLog) SetStatus(value JuryApplyStatus) {
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

// VoteValid
/*
ENUM(
valid   // 有效
invalid // 无效
)
*/
type VoteValid string

func (voteLog *VoteLog) Valid() VoteValid {
	return MustParseVoteValid(voteLog.GetString(VoteLogsFieldValid))
}

func (voteLog *VoteLog) SetValid(value VoteValid) {
	voteLog.Set(VoteLogsFieldValid, value)
}

func (voteLog *VoteLog) Created() types.DateTime {
	return voteLog.GetDateTime(VoteLogsFieldCreated)
}

func (voteLog *VoteLog) Updated() types.DateTime {
	return voteLog.GetDateTime(VoteLogsFieldUpdated)
}
