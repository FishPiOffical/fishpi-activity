package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

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
