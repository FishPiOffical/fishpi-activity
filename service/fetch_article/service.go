package fetch_article

import (
	"bless-activity/model"
	"fmt"
	"log/slog"
	"time"

	"github.com/FishPiOffical/golang-sdk/sdk"
	types2 "github.com/FishPiOffical/golang-sdk/types"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Service struct {
	app core.App
	sdk *sdk.FishPiSDK

	userMap    *maputil.ConcurrentMap[string, *model.User]
	articleMap *maputil.ConcurrentMap[string, *model.RelArticle]

	logger *slog.Logger
}

func NewService(app core.App, sdk *sdk.FishPiSDK) *Service {

	service := &Service{
		app: app,
		sdk: sdk,

		userMap:    maputil.NewConcurrentMap[string, *model.User](100),
		articleMap: maputil.NewConcurrentMap[string, *model.RelArticle](100),

		logger: app.Logger().WithGroup("service_fetch_article"),
	}

	return service
}

func (service *Service) Run() error {

	var activities []*model.Activity
	if err := service.app.RecordQuery(model.DbNameActivities).Where(dbx.HashExp{
		model.ActivitiesFieldHideInList: false,
	}).AndWhere(dbx.Not(dbx.HashExp{
		model.ActivitiesFieldTag: "",
	})).AndWhere(dbx.NewExp("{:now} >= start and {:now} <= end", dbx.Params{
		"now": types.NowDateTime(),
	})).OrderBy(fmt.Sprintf("%s desc", model.ActivitiesFieldEnd)).All(&activities); err != nil {
		return err
	}

	service.logger.Debug("未结束的活动列表", slog.Any("activities", activities))

	service.cacheAuthors()

	cron := service.app.Cron()
	for i, activity := range activities {
		cron.Remove(model.NewCronKeyFetchArticle(activity.Id))
		service.cacheArticles(activity)

		expr := fmt.Sprintf("%d-59/%d * * * *", i, max(5, len(activities)))
		if err := cron.Add(model.NewCronKeyFetchArticle(activity.Id), expr, service.FetchArticlesFunc(activity)); err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) FetchArticlesFunc(activity *model.Activity) func() {
	return func() {
		service.FetchArticles(activity)
	}
}

func (service *Service) FetchArticles(activity *model.Activity) {

	const size = 50
	var page = 1
	for {
		response, err := service.sdk.GetArticles(&types2.GetArticlesRequest{
			Type:    types2.GetArticleTypeTag,
			Keyword: activity.GetTag(),
			Page:    page,
			Size:    size,
		})
		if err != nil {
			service.logger.Error("爬取文章失败", slog.String("activity_id", activity.Id), slog.Any("err", err))
			return
		}
		if response.Code != 0 {
			service.logger.Error("爬取文章失败", slog.String("activity_id", activity.Id), slog.Int("code", response.Code), slog.String("msg", response.Msg))
			return
		}
		if response.Data == nil {
			service.logger.Error("爬取文章失败，返回数据为空", slog.String("activity_id", activity.Id))
			return
		}

		service.logger.Debug("爬取文章结果", slog.String("activity_id", activity.Id), slog.Int("length", len(response.Data.Articles)))

		if len(response.Data.Articles) == 0 {
			break
		}

		// 处理文章 和 作者信息
		for _, article := range response.Data.Articles {
			service.HandleArticle(activity, article)
		}

		if page >= response.Data.Pagination.PaginationPageCount {
			break
		}

		page++
	}
	service.logger.Debug("活动文章爬取完成", slog.String("activity_id", activity.Id))

}

func (service *Service) cacheAuthors() {
	var users []*model.User
	if err := service.app.RecordQuery(model.DbNameUsers).All(&users); err != nil {
		service.logger.Error("缓存作者失败", slog.Any("err", err))
		return
	}
	for _, user := range users {
		service.userMap.Set(user.OId(), user)
	}
}

func (service *Service) cacheArticles(activity *model.Activity) {
	var articles []*model.RelArticle
	if err := service.app.RecordQuery(model.DbNameRelArticles).Where(dbx.HashExp{
		model.RelArticlesFieldActivityId: activity.Id,
	}).All(&articles); err != nil {
		return
	}
	for _, article := range articles {
		service.articleMap.Set(article.OId(), article)
	}
}

func (service *Service) HandleArticle(activity *model.Activity, responseArticle *types2.ArticleInfo) {

	if err := service.HandleAuthor(responseArticle.ArticleAuthor); err != nil {
		service.logger.Error("处理作者失败", slog.String("author_oid", responseArticle.ArticleAuthor.OId), slog.Any("err", err))
		return
	}

	article, exist := service.articleMap.Get(responseArticle.OId)
	if exist {
		// 更新文章
		article.SetActivityId(activity.Id)
		article.SetTitle(responseArticle.ArticleTitle)
		article.SetPreviewContent(responseArticle.ArticlePreviewContent)
		article.SetViewCount(responseArticle.ArticleViewCount)
		article.SetGoodCnt(responseArticle.ArticleGoodCnt)
		article.SetCommentCount(responseArticle.ArticleCommentCount)
		article.SetCollectCnt(responseArticle.ArticleCollectCnt)
		article.SetThankCnt(responseArticle.ArticleThankCnt)
		updatedTime, _ := time.ParseInLocation(time.DateTime, responseArticle.ArticleUpdateTimeStr, time.Local)
		updated, _ := types.ParseDateTime(updatedTime)
		article.SetUpdatedAt(updated)
		if err := service.app.Save(article); err != nil {
			service.logger.Error("更新文章失败", slog.String("article_oid", responseArticle.OId), slog.Any("err", err))
			return
		}
		return
	}

	// 创建文章
	user, userExist := service.userMap.Get(responseArticle.ArticleAuthor.OId)
	if !userExist {
		return
	}

	articleCollection, err := service.app.FindCollectionByNameOrId(model.DbNameRelArticles)
	if err != nil {
		service.logger.Error("查找文章集合失败", slog.Any("err", err))
		return
	}
	article = model.NewRelArticleFromCollection(articleCollection)
	article.SetUserId(user.Id)
	article.SetActivityId(activity.Id)
	article.SetOId(responseArticle.OId)
	article.SetTitle(responseArticle.ArticleTitle)
	article.SetPreviewContent(responseArticle.ArticlePreviewContent)
	article.SetViewCount(responseArticle.ArticleViewCount)
	article.SetGoodCnt(responseArticle.ArticleGoodCnt)
	article.SetCommentCount(responseArticle.ArticleCommentCount)
	article.SetCollectCnt(responseArticle.ArticleCollectCnt)
	article.SetThankCnt(responseArticle.ArticleThankCnt)
	createdTime, _ := time.ParseInLocation(time.DateTime, responseArticle.ArticleCreateTimeStr, time.Local)
	created, _ := types.ParseDateTime(createdTime)
	article.SetCreatedAt(created)
	updatedTime, _ := time.ParseInLocation(time.DateTime, responseArticle.ArticleUpdateTimeStr, time.Local)
	updated, _ := types.ParseDateTime(updatedTime)
	article.SetUpdatedAt(updated)
	if err = service.app.Save(article); err != nil {
		service.logger.Error("创建文章失败", slog.String("article_oid", responseArticle.OId), slog.Any("err", err))
		return
	}
	service.articleMap.Set(article.OId(), article)
}

func (service *Service) HandleAuthor(author *types2.ArticleAuthor) error {
	user, exist := service.userMap.Get(author.OId)
	if exist {
		// 更新用户
		user.SetName(author.UserName)
		user.SetNickname(author.UserNickname)
		user.SetAvatar(author.UserAvatarURL)
		if err := service.app.Save(user); err != nil {
			return err
		}
		return nil
	}
	// 创建用户
	userCollection, err := service.app.FindCollectionByNameOrId(model.DbNameUsers)
	if err != nil {
		return err
	}
	user = model.NewUserFromCollection(userCollection)
	user.SetEmail(fmt.Sprintf("%s@fishpi.cn", author.OId))
	user.SetEmailVisibility(true)
	user.SetVerified(true)
	user.SetOId(author.OId)
	user.SetName(author.UserName)
	user.SetNickname(author.UserNickname)
	user.SetAvatar(author.UserAvatarURL)
	user.SetRandomPassword()
	if err = service.app.Save(user); err != nil {
		return err
	}
	service.userMap.Set(user.OId(), user)
	return nil
}
