//go:generate go-enum --names --values --ptr --mustparse
package model

// CronKey 定时任务key
/*
ENUM(
fetch_article // 爬取文章
)
*/
type CronKey string

func NewCronKeyFetchArticle(activityId string) string {
	return CronKeyFetchArticle.String() + "_" + activityId
}
