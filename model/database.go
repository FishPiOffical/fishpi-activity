package model

import (
	"strconv"
	"time"

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
	_ core.RecordProxy = (*Article)(nil)
	_ core.RecordProxy = (*Shield)(nil)
	_ core.RecordProxy = (*Vote)(nil)
	_ core.RecordProxy = (*VoteLog)(nil)
	_ core.RecordProxy = (*YearlyHistory)(nil)
	_ core.RecordProxy = (*RewardGroup)(nil)
	_ core.RecordProxy = (*Reward)(nil)
	_ core.RecordProxy = (*RewardDistribution)(nil)
	_ core.RecordProxy = (*RelArticle)(nil)
	_ core.RecordProxy = (*UserToken)(nil)
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

func (user *User) RegisteredAt() types.DateTime {
	oId := user.OId()
	if oId == "" {
		return types.DateTime{}
	}

	// Parse oId as milliseconds timestamp
	timestamp, err := strconv.ParseInt(oId, 10, 64)
	if err != nil {
		return types.DateTime{}
	}

	// Convert milliseconds to time.Time
	t := time.UnixMilli(timestamp)

	// Convert to types.DateTime

	dt := types.DateTime{}
	_ = dt.Scan(t)
	return dt
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
	DbNameArticles          = "Articles"
	ArticlesFieldActivityId = "activityId"
	ArticlesFieldUserId     = "userId"
	ArticlesFieldTitle      = "title"
	ArticlesFieldContent    = "content"
	ArticlesFieldShieldId   = "shieldId"
	ArticlesFieldCreated    = "created"
	ArticlesFieldUpdated    = "updated"
)

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

func (article *Article) Created() types.DateTime {
	return article.GetDateTime(ArticlesFieldCreated)
}

func (article *Article) Updated() types.DateTime {
	return article.GetDateTime(ArticlesFieldUpdated)
}

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

func (shield *Shield) Img() string {
	return shield.GetString(ShieldsFieldImg)
}

func (shield *Shield) SetImg(value any) {
	shield.Set(ShieldsFieldImg, value)
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

func (shield *Shield) Ver() string {
	return shield.GetString(ShieldsFieldVer)
}

func (shield *Shield) SetVer(value string) {
	shield.Set(ShieldsFieldVer, value)
}

func (shield *Shield) Scale() string {
	return shield.GetString(ShieldsFieldScale)
}

func (shield *Shield) SetScale(value string) {
	shield.Set(ShieldsFieldScale, value)
}

func (shield *Shield) Size() string {
	return shield.GetString(ShieldsFieldSize)
}

func (shield *Shield) SetSize(value string) {
	shield.Set(ShieldsFieldSize, value)
}

func (shield *Shield) Border() string {
	return shield.GetString(ShieldsFieldBorder)
}

func (shield *Shield) SetBorder(value string) {
	shield.Set(ShieldsFieldBorder, value)
}

func (shield *Shield) BarLen() string {
	return shield.GetString(ShieldsFieldBarLen)
}

func (shield *Shield) SetBarLen(value string) {
	shield.Set(ShieldsFieldBarLen, value)
}

func (shield *Shield) Fontsize() string {
	return shield.GetString(ShieldsFieldFontsize)
}

func (shield *Shield) SetFontsize(value string) {
	shield.Set(ShieldsFieldFontsize, value)
}

func (shield *Shield) BarRadius() string {
	return shield.GetString(ShieldsFieldBarRadius)
}

func (shield *Shield) SetBarRadius(value string) {
	shield.Set(ShieldsFieldBarRadius, value)
}

func (shield *Shield) Shadow() string {
	return shield.GetString(ShieldsFieldShadow)
}

func (shield *Shield) SetShadow(value string) {
	shield.Set(ShieldsFieldShadow, value)
}

func (shield *Shield) Anime() string {
	return shield.GetString(ShieldsFieldAnime)
}

func (shield *Shield) SetAnime(value string) {
	shield.Set(ShieldsFieldAnime, value)
}

func (shield *Shield) Created() types.DateTime {
	return shield.GetDateTime(ShieldsFieldCreated)
}

func (shield *Shield) Updated() types.DateTime {
	return shield.GetDateTime(ShieldsFieldUpdated)
}

const (
	DbNameYearlyHistories                 = "yearlyHistories"
	YearlyHistoriesFieldYear              = "year"              // 年份
	YearlyHistoriesFieldKeyword           = "keyword"           // 关键词
	YearlyHistoriesFieldArticleShieldId   = "articleShieldId"   // 年终征文徽章
	YearlyHistoriesFieldAgeShieldId       = "ageShieldId"       // 年份徽章
	YearlyHistoriesFieldArticleUrl        = "articleUrl"        // 推文链接
	YearlyHistoriesFieldPostArticleUrl    = "postArticleUrl"    // 投稿文章汇总链接
	YearlyHistoriesFieldCollectArticleUrl = "collectArticleUrl" // 征文汇总链接
	YearlyHistoriesFieldActivityId        = "activityId"        // 关联活动ID
	YearlyHistoriesFieldStart             = "start"             // 活动开始时间
	YearlyHistoriesFieldEnd               = "end"               // 活动结束时间
)

type YearlyHistory struct {
	core.BaseRecordProxy
}

func NewYearlyHistory(record *core.Record) *YearlyHistory {
	yearlyHistory := new(YearlyHistory)
	yearlyHistory.SetProxyRecord(record)
	return yearlyHistory
}

func NewYearlyHistoryFromCollection(collection *core.Collection) *YearlyHistory {
	record := core.NewRecord(collection)
	return NewYearlyHistory(record)
}

func (yearlyHistory *YearlyHistory) Year() int {
	return yearlyHistory.GetInt(YearlyHistoriesFieldYear)
}

func (yearlyHistory *YearlyHistory) Keyword() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldKeyword)
}

func (yearlyHistory *YearlyHistory) ArticleShieldId() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldArticleShieldId)
}

func (yearlyHistory *YearlyHistory) AgeShieldId() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldAgeShieldId)
}

func (yearlyHistory *YearlyHistory) ArticleUrl() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldArticleUrl)
}

func (yearlyHistory *YearlyHistory) PostArticleUrl() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldPostArticleUrl)
}

func (yearlyHistory *YearlyHistory) CollectArticleUrl() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldCollectArticleUrl)
}

func (yearlyHistory *YearlyHistory) ActivityId() string {
	return yearlyHistory.GetString(YearlyHistoriesFieldActivityId)
}

func (yearlyHistory *YearlyHistory) Start() types.DateTime {
	return yearlyHistory.GetDateTime(YearlyHistoriesFieldStart)
}

func (yearlyHistory *YearlyHistory) End() types.DateTime {
	return yearlyHistory.GetDateTime(YearlyHistoriesFieldEnd)
}

const (
	DbNameRewardGroups       = "rewardGroups"
	RewardGroupsFieldName    = "name"
	RewardGroupsFieldCreated = "created"
	RewardGroupsFieldUpdated = "updated"
)

type RewardGroup struct {
	core.BaseRecordProxy
}

func NewRewardGroup(record *core.Record) *RewardGroup {
	rewardGroup := new(RewardGroup)
	rewardGroup.SetProxyRecord(record)
	return rewardGroup
}

func NewRewardGroupFromCollection(collection *core.Collection) *RewardGroup {
	record := core.NewRecord(collection)
	return NewRewardGroup(record)
}

func (rewardGroup *RewardGroup) Name() string {
	return rewardGroup.GetString(RewardGroupsFieldName)
}

func (rewardGroup *RewardGroup) SetName(value string) {
	rewardGroup.Set(RewardGroupsFieldName, value)
}

func (rewardGroup *RewardGroup) Created() types.DateTime {
	return rewardGroup.GetDateTime(RewardGroupsFieldCreated)
}

func (rewardGroup *RewardGroup) Updated() types.DateTime {
	return rewardGroup.GetDateTime(RewardGroupsFieldUpdated)
}

const (
	DbNameRewards             = "rewards"
	RewardsFieldRewardGroupId = "rewardGroupId"
	RewardsFieldName          = "name"
	RewardsFieldMin           = "min"
	RewardsFieldMax           = "max"
	RewardsFieldPoint         = "point"
	RewardsFieldShieldIds     = "shieldIds"
	RewardsFieldMore          = "more"
)

type Reward struct {
	core.BaseRecordProxy
}

func NewReward(record *core.Record) *Reward {
	reward := new(Reward)
	reward.SetProxyRecord(record)
	return reward
}

func NewRewardFromCollection(collection *core.Collection) *Reward {
	record := core.NewRecord(collection)
	return NewReward(record)
}

func (reward *Reward) RewardGroupId() string {
	return reward.GetString(RewardsFieldRewardGroupId)
}

func (reward *Reward) SetRewardGroupId(value string) {
	reward.Set(RewardsFieldRewardGroupId, value)
}
func (reward *Reward) Name() string {
	return reward.GetString(RewardsFieldName)
}
func (reward *Reward) SetName(value string) {
	reward.Set(RewardsFieldName, value)
}
func (reward *Reward) Min() int {
	return reward.GetInt(RewardsFieldMin)
}
func (reward *Reward) SetMin(value int) {
	reward.Set(RewardsFieldMin, value)
}
func (reward *Reward) Max() int {
	return reward.GetInt(RewardsFieldMax)
}
func (reward *Reward) SetMax(value int) {
	reward.Set(RewardsFieldMax, value)
}
func (reward *Reward) Point() int {
	return reward.GetInt(RewardsFieldPoint)
}
func (reward *Reward) SetPoint(value int) {
	reward.Set(RewardsFieldPoint, value)
}

func (reward *Reward) ShieldIds() []string {
	return reward.GetStringSlice(RewardsFieldShieldIds)
}

func (reward *Reward) SetShieldIds(value []string) {
	reward.Set(RewardsFieldShieldIds, value)
}

func (reward *Reward) More() string {
	return reward.GetString(RewardsFieldMore)
}
func (reward *Reward) SetMore(value string) {
	reward.Set(RewardsFieldMore, value)
}

const (
	DbNameRewardDistributions       = "rewardDistributions"
	RewardDistributionsFieldVoteId  = "voteId"
	RewardDistributionsFieldUserId  = "userId"
	RewardDistributionsFieldRank    = "rank"
	RewardDistributionsFieldPoint   = "point"
	RewardDistributionsFieldStatus  = "status"
	RewardDistributionsFieldMemo    = "memo"
	RewardDistributionsFieldCreated = "created"
	RewardDistributionsFieldUpdated = "updated"
)

type RewardDistribution struct {
	core.BaseRecordProxy
}

func NewRewardDistribution(record *core.Record) *RewardDistribution {
	rewardDistribution := new(RewardDistribution)
	rewardDistribution.SetProxyRecord(record)
	return rewardDistribution
}

func NewRewardDistributionFromCollection(collection *core.Collection) *RewardDistribution {
	record := core.NewRecord(collection)
	return NewRewardDistribution(record)
}

func (rd *RewardDistribution) VoteId() string {
	return rd.GetString(RewardDistributionsFieldVoteId)
}

func (rd *RewardDistribution) SetVoteId(value string) {
	rd.Set(RewardDistributionsFieldVoteId, value)
}

func (rd *RewardDistribution) UserId() string {
	return rd.GetString(RewardDistributionsFieldUserId)
}

func (rd *RewardDistribution) SetUserId(value string) {
	rd.Set(RewardDistributionsFieldUserId, value)
}

func (rd *RewardDistribution) Rank() int {
	return rd.GetInt(RewardDistributionsFieldRank)
}

func (rd *RewardDistribution) SetRank(value int) {
	rd.Set(RewardDistributionsFieldRank, value)
}

func (rd *RewardDistribution) Point() int {
	return rd.GetInt(RewardDistributionsFieldPoint)
}

func (rd *RewardDistribution) SetPoint(value int) {
	rd.Set(RewardDistributionsFieldPoint, value)
}

func (rd *RewardDistribution) Status() DistributionStatus {
	return MustParseDistributionStatus(rd.GetString(RewardDistributionsFieldStatus))
}

func (rd *RewardDistribution) SetStatus(value DistributionStatus) {
	rd.Set(RewardDistributionsFieldStatus, value)
}

func (rd *RewardDistribution) Memo() string {
	return rd.GetString(RewardDistributionsFieldMemo)
}

func (rd *RewardDistribution) SetMemo(value string) {
	rd.Set(RewardDistributionsFieldMemo, value)
}

func (rd *RewardDistribution) Created() types.DateTime {
	return rd.GetDateTime(RewardDistributionsFieldCreated)
}

func (rd *RewardDistribution) Updated() types.DateTime {
	return rd.GetDateTime(RewardDistributionsFieldUpdated)
}

const (
	DbNameRelArticles              = "relArticles"
	RelArticlesFieldUserId         = "userId"         // 用户
	RelArticlesFieldActivityId     = "activityId"     // 活动ID
	RelArticlesFieldOId            = "oId"            // 文章ID
	RelArticlesFieldTitle          = "title"          // 标题
	RelArticlesFieldPreviewContent = "previewContent" // 预览内容
	RelArticlesFieldViewCount      = "viewCount"      // 浏览量
	RelArticlesFieldGoodCnt        = "goodCnt"        // 点赞量
	RelArticlesFieldCommentCount   = "commentCount"   // 评论数
	RelArticlesFieldCollectCnt     = "collectCnt"     // 收藏数
	RelArticlesFieldThankCnt       = "thankCnt"       // 感谢数
	RelArticlesFieldCreatedAt      = "createdAt"      // 发表时间
	RelArticlesFieldUpdatedAt      = "updatedAt"      // 更新时间
	RelArticlesFieldCreated        = "created"        // 爬取时间
	RelArticlesFieldUpdated        = "updated"        // 爬取更新时间
)

type RelArticle struct {
	core.BaseRecordProxy
}

func NewRelArticle(record *core.Record) *RelArticle {
	relArticle := new(RelArticle)
	relArticle.SetProxyRecord(record)
	return relArticle
}

func NewRelArticleFromCollection(collection *core.Collection) *RelArticle {
	record := core.NewRecord(collection)
	return NewRelArticle(record)
}

func (ra *RelArticle) UserId() string {
	return ra.GetString(RelArticlesFieldUserId)
}

func (ra *RelArticle) SetUserId(value string) {
	ra.Set(RelArticlesFieldUserId, value)
}

func (ra *RelArticle) ActivityId() string {
	return ra.GetString(RelArticlesFieldActivityId)
}

func (ra *RelArticle) SetActivityId(value string) {
	ra.Set(RelArticlesFieldActivityId, value)
}

func (ra *RelArticle) OId() string {
	return ra.GetString(RelArticlesFieldOId)
}

func (ra *RelArticle) SetOId(value string) {
	ra.Set(RelArticlesFieldOId, value)
}

func (ra *RelArticle) Title() string {
	return ra.GetString(RelArticlesFieldTitle)
}

func (ra *RelArticle) SetTitle(value string) {
	ra.Set(RelArticlesFieldTitle, value)
}

func (ra *RelArticle) PreviewContent() string {
	return ra.GetString(RelArticlesFieldPreviewContent)
}

func (ra *RelArticle) SetPreviewContent(value string) {
	ra.Set(RelArticlesFieldPreviewContent, value)
}

func (ra *RelArticle) ViewCount() int {
	return ra.GetInt(RelArticlesFieldViewCount)
}

func (ra *RelArticle) SetViewCount(value int) {
	ra.Set(RelArticlesFieldViewCount, value)
}

func (ra *RelArticle) GoodCnt() int {
	return ra.GetInt(RelArticlesFieldGoodCnt)
}

func (ra *RelArticle) SetGoodCnt(value int) {
	ra.Set(RelArticlesFieldGoodCnt, value)
}

func (ra *RelArticle) CommentCount() int {
	return ra.GetInt(RelArticlesFieldCommentCount)
}

func (ra *RelArticle) SetCommentCount(value int) {
	ra.Set(RelArticlesFieldCommentCount, value)
}

func (ra *RelArticle) CollectCnt() int {
	return ra.GetInt(RelArticlesFieldCollectCnt)
}

func (ra *RelArticle) SetCollectCnt(value int) {
	ra.Set(RelArticlesFieldCollectCnt, value)
}

func (ra *RelArticle) ThankCnt() int {
	return ra.GetInt(RelArticlesFieldThankCnt)
}

func (ra *RelArticle) SetThankCnt(value int) {
	ra.Set(RelArticlesFieldThankCnt, value)
}

func (ra *RelArticle) CreatedAt() types.DateTime {
	return ra.GetDateTime(RelArticlesFieldCreatedAt)
}

func (ra *RelArticle) SetCreatedAt(value types.DateTime) {
	ra.Set(RelArticlesFieldCreatedAt, value)
}

func (ra *RelArticle) UpdatedAt() types.DateTime {
	return ra.GetDateTime(RelArticlesFieldUpdatedAt)
}

func (ra *RelArticle) SetUpdatedAt(value types.DateTime) {
	ra.Set(RelArticlesFieldUpdatedAt, value)
}

func (ra *RelArticle) Created() types.DateTime {
	return ra.GetDateTime(RelArticlesFieldCreated)
}

func (ra *RelArticle) Updated() types.DateTime {
	return ra.GetDateTime(RelArticlesFieldUpdated)
}
