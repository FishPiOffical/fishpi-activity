package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
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
