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
	DbNameActivities           = "activities"
	ActivitiesFieldName        = "name"
	ActivitiesFieldSlug        = "slug"
	ActivitiesFieldArticleUrl  = "articleUrl"
	ActivitiesFieldExternalUrl = "externalUrl"
	ActivitiesFieldDesc        = "desc"
	ActivitiesFieldStart       = "start"
	ActivitiesFieldEnd         = "end"
	ActivitiesFieldVoteId      = "voteId"
	ActivitiesFieldRewardGroup = "rewardGroup"
	ActivitiesFieldHideInList  = "hideInList"
	ActivitiesFieldCreated     = "created"
	ActivitiesFieldUpdated     = "updated"
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

func (activity *Activity) RewardGroup() string {
	return activity.GetString(ActivitiesFieldRewardGroup)
}

func (activity *Activity) SetRewardGroup(value string) {
	activity.Set(ActivitiesFieldRewardGroup, value)
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
	DbNameRewards     = "rewards"
	RewardsFieldGroup = "group"
	RewardsFieldMin   = "min"
	RewardsFieldMax   = "max"
	RewardsFieldPoint = "point"
	RewardsFieldMore  = "more"
)
