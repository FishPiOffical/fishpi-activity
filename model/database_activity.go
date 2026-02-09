//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameActivities                        = "activities"               // 活动表
	ActivitiesFieldName                     = "name"                     // 活动名称
	ActivitiesFieldTemplate                 = "template"                 // 活动模板
	ActivitiesFieldSlug                     = "slug"                     // 活动标识
	ActivitiesFieldArticleUrl               = "articleUrl"               // 鱼排文章链接
	ActivitiesFieldExternalUrl              = "externalUrl"              // 活动外部链接
	ActivitiesFieldDesc                     = "desc"                     // 活动描述
	ActivitiesFieldTag                      = "tag"                      // 活动绑定鱼排标签
	ActivitiesFieldStart                    = "start"                    // 活动开始时间
	ActivitiesFieldEnd                      = "end"                      // 活动结束时间
	ActivitiesFieldVoteId                   = "voteId"                   // 投票ID
	ActivitiesFieldRewardGroupId            = "rewardGroupId"            // 奖励组ID
	ActivitiesFieldRewardDistributionStatus = "rewardDistributionStatus" // 奖励发放状态
	ActivitiesFieldHideInList               = "hideInList"               // 是否在列表隐藏
	ActivitiesFieldChildActivityIds         = "childActivityIds"         // 子活动ID列表
	ActivitiesFieldImage                    = "image"                    // 活动图片
	ActivitiesFieldImages                   = "images"                   // 活动图片(多张)
	ActivitiesFieldMetadata                 = "metadata"                 // 元数据(JSON)
	ActivitiesFieldCreated                  = "created"                  // 创建时间
	ActivitiesFieldUpdated                  = "updated"                  // 更新时间
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

func (activity *Activity) GetName() string {
	return activity.GetString(ActivitiesFieldName)
}

func (activity *Activity) SetName(value string) {
	activity.Set(ActivitiesFieldName, value)
}

func (activity *Activity) GetMetadata() any {
	return activity.Get(ActivitiesFieldMetadata)
}

func (activity *Activity) SetMetadata(value any) {
	activity.Set(ActivitiesFieldMetadata, value)
}

func (activity *Activity) GetImages() []string {
	return activity.GetStringSlice(ActivitiesFieldImages)
}

func (activity *Activity) SetImages(value []string) {
	activity.Set(ActivitiesFieldImages, value)
}

// ActivityTemplate 活动模版
/*
ENUM(
article // 征文活动
redirect // 外链活动
)
*/
type ActivityTemplate string

func (activity *Activity) GetTemplate() ActivityTemplate {
	return MustParseActivityTemplate(activity.GetString(ActivitiesFieldTemplate))
}

func (activity *Activity) SetTemplate(value ActivityTemplate) {
	activity.Set(ActivitiesFieldTemplate, value)
}

func (activity *Activity) GetSlug() string {
	return activity.GetString(ActivitiesFieldSlug)
}

func (activity *Activity) SetSlug(value string) {
	activity.Set(ActivitiesFieldSlug, value)
}

func (activity *Activity) GetArticleUrl() string {
	return activity.GetString(ActivitiesFieldArticleUrl)
}

func (activity *Activity) SetArticleUrl(value string) {
	activity.Set(ActivitiesFieldArticleUrl, value)
}

func (activity *Activity) GetExternalUrl() string {
	return activity.GetString(ActivitiesFieldExternalUrl)
}

func (activity *Activity) SetExternalUrl(value string) {
	activity.Set(ActivitiesFieldExternalUrl, value)
}

func (activity *Activity) GetDesc() string {
	return activity.GetString(ActivitiesFieldDesc)
}

func (activity *Activity) SetDesc(value string) {
	activity.Set(ActivitiesFieldDesc, value)
}

func (activity *Activity) GetTag() string {
	return activity.GetString(ActivitiesFieldTag)
}

func (activity *Activity) SetTag(value string) {
	activity.Set(ActivitiesFieldTag, value)
}

func (activity *Activity) GetStart() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldStart)
}

func (activity *Activity) SetStart(value types.DateTime) {
	activity.Set(ActivitiesFieldStart, value)
}

func (activity *Activity) GetEnd() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldEnd)
}

func (activity *Activity) SetEnd(value types.DateTime) {
	activity.Set(ActivitiesFieldEnd, value)
}

func (activity *Activity) GetVoteId() string {
	return activity.GetString(ActivitiesFieldVoteId)
}

func (activity *Activity) SetVoteId(value string) {
	activity.Set(ActivitiesFieldVoteId, value)
}

func (activity *Activity) GetRewardGroupId() string {
	return activity.GetString(ActivitiesFieldRewardGroupId)
}

func (activity *Activity) SetRewardGroupId(value string) {
	activity.Set(ActivitiesFieldRewardGroupId, value)
}

// DistributionStatus
/*
ENUM(
pending      // 待发放
distributing // 发放中
failed       // 发放失败
success      // 发放成功
)
*/
type DistributionStatus string

func (activity *Activity) GetRewardDistributionStatus() DistributionStatus {
	return MustParseDistributionStatus(activity.GetString(ActivitiesFieldRewardDistributionStatus))
}

func (activity *Activity) SetRewardDistributionStatus(value DistributionStatus) {
	activity.Set(ActivitiesFieldRewardDistributionStatus, value)
}

func (activity *Activity) GetHideInList() bool {
	return activity.GetBool(ActivitiesFieldHideInList)
}

func (activity *Activity) SetHideInList(value bool) {
	activity.Set(ActivitiesFieldHideInList, value)
}

func (activity *Activity) GetChildActivityIds() []string {
	return activity.GetStringSlice(ActivitiesFieldChildActivityIds)
}

func (activity *Activity) SetChildActivityIds(childActivityIds []string) {
	activity.Set(ActivitiesFieldChildActivityIds, childActivityIds)
}

func (activity *Activity) GetImage() string {
	return activity.GetString(ActivitiesFieldImage)
}

func (activity *Activity) SetImage(value string) {
	activity.Set(ActivitiesFieldImage, value)
}

func (activity *Activity) GetCreated() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldCreated)
}

func (activity *Activity) GetUpdated() types.DateTime {
	return activity.GetDateTime(ActivitiesFieldUpdated)
}
