package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

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
