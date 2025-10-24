package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	CommonFieldId = "id"
)

var (
	_ core.RecordProxy = (*User)(nil)
	_ core.RecordProxy = (*Config)(nil)
	_ core.RecordProxy = (*Activity)(nil)
)

const (
	DbNameUsers               = "users"
	UsersFieldEmail           = "email"
	UsersFieldEmailVisibility = "emailVisibility"
	UsersFieldVerified        = "verified"
	UsersFieldName            = "name"
	UsersFieldNickname        = "nickname"
	UsersFieldAvatar          = "avatar"
	UsersFieldOId             = "oId"
	UsersFieldCreated         = "created"
	UsersFieldUpdated         = "updated"
)

type User struct {
	core.BaseRecordProxy
}

func NewUser(record *core.Record) *User {
	user := new(User)
	user.SetProxyRecord(record)
	return user
}

func NewUserFromCollection(collection *core.Collection) *User {
	record := core.NewRecord(collection)
	return NewUser(record)
}

func (user *User) Name() string {
	return user.GetString(UsersFieldName)
}

func (user *User) SetName(value string) {
	user.Set(UsersFieldName, value)
}

func (user *User) Nickname() string {
	return user.GetString(UsersFieldNickname)
}

func (user *User) SetNickname(value string) {
	user.Set(UsersFieldNickname, value)
}

func (user *User) Avatar() string {
	return user.GetString(UsersFieldAvatar)
}

func (user *User) SetAvatar(value string) {
	user.Set(UsersFieldAvatar, value)
}

func (user *User) OId() string {
	return user.GetString(UsersFieldOId)
}

func (user *User) SetOId(value string) {
	user.Set(UsersFieldOId, value)
}

func (user *User) Created() types.DateTime {
	return user.GetDateTime(UsersFieldCreated)
}

func (user *User) Updated() types.DateTime {
	return user.GetDateTime(UsersFieldUpdated)
}

const (
	DbNameConfigs     = "configs"
	ConfigsFieldKey   = "key"
	ConfigsFieldValue = "value"
)

type Config struct {
	core.BaseRecordProxy
}

func NewConfig(record *core.Record) *Config {
	config := new(Config)
	config.SetProxyRecord(record)
	return config
}

func NewConfigFromCollection(collection *core.Collection) *Config {
	record := core.NewRecord(collection)
	return NewConfig(record)
}

func (config *Config) Key() ConfigKey {
	return MustParseConfigKey(config.GetString(ConfigsFieldKey))
}

func (config *Config) SetKey(value ConfigKey) {
	config.Set(ConfigsFieldKey, value)
}

func (config *Config) Value() string {
	return config.GetString(ConfigsFieldValue)
}

func (config *Config) SetValue(value string) {
	config.Set(ConfigsFieldValue, value)
}

const (
	DbNameActivities             = "activities"
	ActivitiesFieldName          = "name"
	ActivitiesFieldSlug          = "slug"
	ActivitiesFieldArticleUrl    = "articleUrl"
	ActivitiesFieldExternalUrl   = "externalUrl"
	ActivitiesFieldDesc          = "desc"
	ActivitiesFieldStart         = "start"
	ActivitiesFieldEnd           = "end"
	ActivitiesFieldVoteId        = "voteId"
	ActivitiesFieldRewardGroupId = "rewardGroupId"
	ActivitiesFieldHideInList    = "hideInList"
	ActivitiesFieldCreated       = "created"
	ActivitiesFieldUpdated       = "updated"
)

type Activity struct {
	core.BaseRecordProxy
}

func NewActivity(record *core.Record) *Activity {
	activity := new(Activity)
	activity.SetProxyRecord(record)
	return activity
}

func NewActivityFromCollection(collection *core.Collection) *Activity {
	record := core.NewRecord(collection)
	return NewActivity(record)
}

func (activity *Activity) Name() string {
	return activity.GetString(ActivitiesFieldName)
}

func (activity *Activity) SetName(value string) {
	activity.Set(ActivitiesFieldName, value)
}

func (activity *Activity) Slug() string {
	return activity.GetString(ActivitiesFieldSlug)
}

func (activity *Activity) SetSlug(value string) {
	activity.Set(ActivitiesFieldSlug, value)
}

func (activity *Activity) ArticleUrl() string {
	return activity.GetString(ActivitiesFieldArticleUrl)
}

func (activity *Activity) SetArticleUrl(value string) {
	activity.Set(ActivitiesFieldArticleUrl, value)
}

func (activity *Activity) ExternalUrl() string {
	return activity.GetString(ActivitiesFieldExternalUrl)
}

func (activity *Activity) SetExternalUrl(value string) {
	activity.Set(ActivitiesFieldExternalUrl, value)
}

func (activity *Activity) Desc() string {
	return activity.GetString(ActivitiesFieldDesc)
}

func (activity *Activity) SetDesc(value string) {
	activity.Set(ActivitiesFieldDesc, value)
}

func (activity *Activity) Start() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldStart)
}

func (activity *Activity) SetStart(value types.DateTime) {
	activity.Set(ActivitiesFieldStart, value)
}

func (activity *Activity) End() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldEnd)
}

func (activity *Activity) SetEnd(value types.DateTime) {
	activity.Set(ActivitiesFieldEnd, value)
}

func (activity *Activity) VoteId() string {
	return activity.GetString(ActivitiesFieldVoteId)
}

func (activity *Activity) SetVoteId(value string) {
	activity.Set(ActivitiesFieldVoteId, value)
}

func (activity *Activity) RewardGroupId() string {
	return activity.GetString(ActivitiesFieldRewardGroupId)
}

func (activity *Activity) SetRewardGroupId(value string) {
	activity.Set(ActivitiesFieldRewardGroupId, value)
}

func (activity *Activity) HideInList() bool {
	return activity.GetBool(ActivitiesFieldHideInList)
}

func (activity *Activity) SetHideInList(value bool) {
	activity.Set(ActivitiesFieldHideInList, value)
}

func (activity *Activity) Created() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldCreated)
}

func (activity *Activity) Updated() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldUpdated)
}

const (
	DbNameArticles          = "Articles"
	ArticlesFieldActivityId = "activityId"
	ArticlesFieldUserId     = "userId"
	ArticlesFieldTitle      = "title"
	ArticlesFieldContent    = "content"
	ArticlesFieldShieldId   = "shieldId"
	ArticlesFieldCreated    = "created"
	ArticlesFieldUpdated    = "updated"
)

const (
	DbNameShields         = "shields"
	ShieldsFieldText      = "text"
	ShieldsFieldImg       = "img"
	ShieldsFieldUrl       = "url"
	ShieldsFieldBackcolor = "backcolor"
	ShieldsFieldFontcolor = "fontcolor"
	ShieldsFieldVer       = "ver"
	ShieldsFieldScale     = "scale"
	ShieldsFieldSize      = "size"
	ShieldsFieldBorder    = "border"
	ShieldsFieldBarLen    = "barlen"
	ShieldsFieldFontsize  = "fontsize"
	ShieldsFieldBarRadius = "barradius"
	ShieldsFieldShadow    = "shadow"
	ShieldsFieldAnime     = "anime"
	ShieldsFieldCreated   = "created"
	ShieldsFieldUpdated   = "updated"
)

const (
	DbNameVotes     = "votes"
	VotesFieldName  = "name"
	VotesFieldDesc  = "desc"
	VotesFieldTimes = "times"
	VotesFieldStart = "start"
	VotesFieldEnd   = "end"
)

const (
	DbNameVoteLogs          = "voteLogs"
	VoteLogsFieldVoteId     = "voteId"
	VoteLogsFieldFromUserId = "fromUserId"
	VoteLogsFieldToUserId   = "toUserId"
	VoteLogsFieldComment    = "comment"
	VoteLogsFieldCreated    = "created"
	VoteLogsFieldUpdated    = "updated"
)

const (
	DbNameYearlyHistories               = "yearlyHistories"
	YearlyHistoriesFieldYear            = "year"
	YearlyHistoriesFieldKeyword         = "keyword"
	YearlyHistoriesFieldArticleShieldId = "articleShieldId"
	YearlyHistoriesFieldAgeShieldId     = "ageShieldId"
	YearlyHistoriesFieldArticleUrl      = "articleUrl"
)

const (
	DbNameRewardGroups       = "rewardGroups"
	RewardGroupsFieldName    = "name"
	RewardGroupsFieldCreated = "created"
	RewardGroupsFieldUpdated = "updated"
)

const (
	DbNameRewards             = "rewards"
	RewardsFieldRewardGroupId = "rewardGroupId"
	RewardsFieldName          = "name"
	RewardsFieldMin           = "min"
	RewardsFieldMax           = "max"
	RewardsFieldPoint         = "point"
	RewardsFieldMore          = "more"
)

// Shield wrapper type
type Shield struct {
	core.BaseRecordProxy
}

func NewShield(record *core.Record) *Shield {
	shield := new(Shield)
	shield.SetProxyRecord(record)
	return shield
}

func NewShieldFromCollection(collection *core.Collection) *Shield {
	record := core.NewRecord(collection)
	return NewShield(record)
}

func (shield *Shield) Text() string {
	return shield.GetString(ShieldsFieldText)
}

func (shield *Shield) SetText(value string) {
	shield.Set(ShieldsFieldText, value)
}

func (shield *Shield) Url() string {
	return shield.GetString(ShieldsFieldUrl)
}

func (shield *Shield) SetUrl(value string) {
	shield.Set(ShieldsFieldUrl, value)
}

func (shield *Shield) Backcolor() string {
	return shield.GetString(ShieldsFieldBackcolor)
}

func (shield *Shield) SetBackcolor(value string) {
	shield.Set(ShieldsFieldBackcolor, value)
}

func (shield *Shield) Fontcolor() string {
	return shield.GetString(ShieldsFieldFontcolor)
}

func (shield *Shield) SetFontcolor(value string) {
	shield.Set(ShieldsFieldFontcolor, value)
}

// Article wrapper type
type Article struct {
	core.BaseRecordProxy
}

func NewArticle(record *core.Record) *Article {
	article := new(Article)
	article.SetProxyRecord(record)
	return article
}

func NewArticleFromCollection(collection *core.Collection) *Article {
	record := core.NewRecord(collection)
	return NewArticle(record)
}

func (article *Article) ActivityId() string {
	return article.GetString(ArticlesFieldActivityId)
}

func (article *Article) SetActivityId(value string) {
	article.Set(ArticlesFieldActivityId, value)
}

func (article *Article) UserId() string {
	return article.GetString(ArticlesFieldUserId)
}

func (article *Article) SetUserId(value string) {
	article.Set(ArticlesFieldUserId, value)
}

func (article *Article) Title() string {
	return article.GetString(ArticlesFieldTitle)
}

func (article *Article) SetTitle(value string) {
	article.Set(ArticlesFieldTitle, value)
}

func (article *Article) Content() string {
	return article.GetString(ArticlesFieldContent)
}

func (article *Article) SetContent(value string) {
	article.Set(ArticlesFieldContent, value)
}

func (article *Article) ShieldId() string {
	return article.GetString(ArticlesFieldShieldId)
}

func (article *Article) SetShieldId(value string) {
	article.Set(ArticlesFieldShieldId, value)
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
