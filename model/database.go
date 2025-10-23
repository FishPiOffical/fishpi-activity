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
	ActivitiesFieldDesc        = "desc"
	ActivitiesFieldStart       = "start"
	ActivitiesFieldEnd         = "end"
	ActivitiesFieldVoteId      = "voteId"
	ActivitiesFieldRewardGroup = "rewardGroup"
	ActivitiesFieldHideInList  = "hideInList"
)

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
