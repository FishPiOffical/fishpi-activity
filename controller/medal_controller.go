package controller

import (
	"bless-activity/model"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	types2 "github.com/FishPiOffical/golang-sdk/types"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/types"
)

type MedalController struct {
	*BaseController

	event *core.ServeEvent
	group *router.RouterGroup[*core.RequestEvent]
	app   core.App

	logger *slog.Logger
}

func NewMedalController(event *core.ServeEvent, group *router.RouterGroup[*core.RequestEvent], base *BaseController) *MedalController {
	logger := event.App.Logger().With(
		slog.String("controller", "medal"),
	)

	controller := &MedalController{
		BaseController: base,
		event:          event,
		group:          group,
		app:            event.App,
		logger:         logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *MedalController) registerRoutes() {
	group := controller.group.Group("/admin/medal").Bind(
		RequireAdminRole(),
	)

	// 勋章列表
	group.GET("/list", controller.List)
	// 勋章详情
	group.GET("/detail/{medalId}", controller.Detail)
	// 创建勋章
	group.POST("/create", controller.Create)
	// 编辑勋章
	group.PUT("/edit/{medalId}", controller.Edit)
	// 删除勋章
	group.DELETE("/delete/{medalId}", controller.Delete)

	// 同步相关
	group.POST("/sync/all", controller.SyncAllMedals)
	group.POST("/sync/{medalId}", controller.SyncSingleMedal)
	group.POST("/sync/owners/all", controller.SyncAllMedalOwners)
	group.POST("/sync/owners/{medalId}", controller.SyncSingleMedalOwners)
	group.POST("/sync/user/{userId}", controller.SyncUserMedals)

	// 勋章拥有者列表
	group.GET("/owners/{medalId}", controller.GetMedalOwners)

	// 给用户授予/撤销勋章
	group.POST("/grant", controller.GrantMedal)
	group.POST("/grant/batch", controller.GrantMedalBatch)
	group.POST("/revoke", controller.RevokeMedal)

	// 搜索勋章
	group.GET("/search", controller.Search)

	// 用户选择相关
	group.GET("/users/search", controller.SearchUsers)
	group.GET("/activities", controller.GetActivities)
	group.GET("/activity/{activityId}/participants", controller.GetActivityParticipants)
	group.GET("/vote/{voteId}/jury", controller.GetVoteJuryMembers)
}

func (controller *MedalController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

// List 获取勋章列表
func (controller *MedalController) List(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("list")

	page, _ := strconv.Atoi(event.Request.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(event.Request.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var medals []*model.Medal
	var total int

	// 先查询总数
	if err := event.App.RecordQuery(model.DbNameMedals).Select("count(*)").Row(&total); err != nil {
		logger.Error("查询勋章总数失败", slog.Any("err", err))
		return event.InternalServerError("查询勋章总数失败", err)
	}

	// 再查询列表
	if err := event.App.RecordQuery(model.DbNameMedals).OrderBy("rowid DESC").Limit(int64(pageSize)).Offset(int64((page - 1) * pageSize)).All(&medals); err != nil {
		logger.Error("查询勋章列表失败", slog.Any("err", err))
		return event.InternalServerError("查询勋章列表失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items":      medals,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + pageSize - 1) / pageSize,
	})
}

// Detail 获取勋章详情
func (controller *MedalController) Detail(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("detail")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

	medal := new(model.Medal)
	if err := event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: medalId,
	}).One(medal); err != nil {
		logger.Error("查询勋章详情失败", slog.Any("err", err), slog.String("medal_id", medalId))
		return event.NotFoundError("勋章不存在", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"medal": medal,
	})
}

// Create 创建勋章（先在鱼排创建，然后同步到本地）
func (controller *MedalController) Create(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("create")

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Attr        string `json:"attr"`
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if req.Name == "" {
		return event.BadRequestError("勋章名称不能为空", nil)
	}

	// 先在鱼排创建
	resp, err := controller.fishPiSdk.PostMedalAdminCreate(req.Name, types2.MedalType(req.Type), req.Description, req.Attr)
	if err != nil {
		logger.Error("在鱼排创建勋章失败", slog.Any("err", err))
		return event.InternalServerError("在鱼排创建勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("在鱼排创建勋章失败", slog.Any("resp", resp))
		return event.InternalServerError("在鱼排创建勋章失败: "+resp.Msg, nil)
	}

	// 只返回了resp.Data.OId，没有什么用，获取不到新创建的勋章详情，所以只能结束。

	//medalData := resp.Data
	//
	//// 保存到本地数据库
	//medalCollection, err := event.App.FindCollectionByNameOrId(model.DbNameMedals)
	//if err != nil {
	//	logger.Error("获取勋章集合失败", slog.Any("err", err))
	//	return event.InternalServerError("获取勋章集合失败", err)
	//}
	//
	//medal := model.NewMedalFromCollection(medalCollection)
	//medal.SetOId(medalData.OId)
	//medal.SetMedalId(medalData.MedalId)
	//medal.SetType(medalData.MedalType)
	//medal.SetName(medalData.MedalName)
	//medal.SetDescription(medalData.MedalDescription)
	//medal.SetAttr(medalData.MedalAttr)
	//
	//if err = event.App.Save(medal); err != nil {
	//	logger.Error("保存勋章失败", slog.Any("err", err))
	//	return event.InternalServerError("保存勋章失败", err)
	//}

	logger.Info("创建勋章成功", slog.Any("medal", resp.Data))

	return event.JSON(http.StatusOK, map[string]any{
		"medal": resp.Data,
	})
}

// Edit 编辑勋章
func (controller *MedalController) Edit(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("edit")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Attr        string `json:"attr"`
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	// 先在鱼排编辑
	resp, err := controller.fishPiSdk.PostMedalAdminEdit(medalId, req.Name, types2.MedalType(req.Type), req.Description, req.Attr)
	if err != nil {
		logger.Error("在鱼排编辑勋章失败", slog.Any("err", err))
		return event.InternalServerError("在鱼排编辑勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("在鱼排编辑勋章失败", slog.Any("resp", resp))
		return event.InternalServerError("在鱼排编辑勋章失败: "+resp.Msg, nil)
	}

	// 更新本地数据库
	medal := new(model.Medal)
	if err = event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: medalId,
	}).One(medal); err != nil {
		logger.Error("查询本地勋章失败", slog.Any("err", err))
		return event.NotFoundError("本地勋章不存在", err)
	}

	medal.SetType(types2.MedalType(req.Type))
	medal.SetName(req.Name)
	medal.SetDescription(req.Description)
	medal.SetAttr(req.Attr)

	if err = event.App.Save(medal); err != nil {
		logger.Error("更新本地勋章失败", slog.Any("err", err))
		return event.InternalServerError("更新本地勋章失败", err)
	}

	logger.Info("编辑勋章成功", slog.Any("medal", medal))

	return event.JSON(http.StatusOK, map[string]any{
		"medal": medal,
	})
}

// Delete 删除勋章
func (controller *MedalController) Delete(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("delete")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

	// 先在鱼排删除
	resp, err := controller.fishPiSdk.PostMedalAdminDelete(medalId)
	if err != nil {
		logger.Error("在鱼排删除勋章失败", slog.Any("err", err))
		return event.InternalServerError("在鱼排删除勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("在鱼排删除勋章失败", slog.Any("resp", resp))
		return event.InternalServerError("在鱼排删除勋章失败: "+resp.Msg, nil)
	}

	// 删除本地勋章
	medal := new(model.Medal)
	if err = event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: medalId,
	}).One(medal); err != nil {
		logger.Warn("本地勋章不存在", slog.Any("err", err))
		return event.JSON(http.StatusOK, map[string]any{
			"deleted": true,
			"message": "鱼排勋章已删除，本地不存在",
		})
	}

	if err = event.App.Delete(medal); err != nil {
		logger.Error("删除本地勋章失败", slog.Any("err", err))
		return event.InternalServerError("删除本地勋章失败", err)
	}

	// 删除关联的勋章拥有者记录
	if _, err = event.App.DB().Delete(model.DbNameMedalOwners, dbx.HashExp{
		model.MedalOwnersFieldMedalId: medalId,
	}).Execute(); err != nil {
		logger.Warn("删除勋章拥有者记录失败", slog.Any("err", err))
	}

	logger.Info("删除勋章成功", slog.String("medal_id", medalId))

	return event.JSON(http.StatusOK, map[string]any{
		"deleted": true,
	})
}

// SyncAllMedals 同步所有勋章
func (controller *MedalController) SyncAllMedals(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("sync_all_medals")

	const pageSize = 100
	var page = 1

	var medals []*types2.Medal
	for {
		resp, err := controller.fishPiSdk.PostMedalAdminList(page, pageSize)
		if err != nil {
			logger.Error("查询勋章列表失败", slog.Any("err", err), slog.Int("page", page), slog.Int("page_size", pageSize))
			return event.InternalServerError("查询勋章列表失败", err)
		}
		if resp.Code != 0 {
			logger.Error("查询勋章列表失败", slog.Any("resp", resp), slog.Int("page", page), slog.Int("page_size", pageSize))
			return event.InternalServerError("查询勋章列表失败: "+resp.Msg, nil)
		}
		if len(resp.Data) == 0 {
			break
		}
		medals = append(medals, resp.Data...)

		if len(resp.Data) < pageSize {
			break
		}

		page++
		time.Sleep(time.Second)
	}

	medalCollection, err := event.App.FindCollectionByNameOrId(model.DbNameMedals)
	if err != nil {
		logger.Error("获取勋章集合失败", slog.Any("err", err))
		return event.InternalServerError("获取勋章集合失败", err)
	}

	var (
		deletedCount int64
		createdCount int64
		updatedCount int64
	)

	if err = event.App.RunInTransaction(func(txApp core.App) error {
		// 标记所有勋章为待删除
		_, txErr := txApp.DB().Update(model.DbNameMedals, dbx.Params{
			model.MedalsFieldMedalId: dbx.NewExp(fmt.Sprintf("'pending_delete_' || `%s`", model.MedalsFieldMedalId)),
			model.MedalsFieldName:    dbx.NewExp(fmt.Sprintf("'pending_delete_' || `%s`", model.MedalsFieldName)),
		}, dbx.Not(dbx.HashExp{
			model.MedalsFieldOId: "",
		})).Execute()
		if txErr != nil {
			logger.Error("标记待删除勋章失败", slog.Any("err", txErr))
			return txErr
		}

		medalRecord := new(model.Medal)
		for _, medalData := range medals {
			if txErr = txApp.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{model.MedalsFieldOId: medalData.OId}).One(medalRecord); txErr != nil {
				logger.Debug("查询勋章失败，创建新勋章", slog.String("medal_oId", medalData.OId))
				medalRecord = model.NewMedalFromCollection(medalCollection)
				medalRecord.SetOId(medalData.OId)
				medalRecord.SetMedalId(medalData.MedalId)
				medalRecord.SetType(medalData.MedalType)
				medalRecord.SetName(medalData.MedalName)
				medalRecord.SetDescription(medalData.MedalDescription)
				medalRecord.SetAttr(medalData.MedalAttr)
				if txErr = txApp.Save(medalRecord); txErr != nil {
					logger.Error("保存勋章失败", slog.Any("medal", medalRecord), slog.Any("err", txErr))
					return txErr
				}
				createdCount++
			} else {
				medalRecord.SetMedalId(medalData.MedalId)
				medalRecord.SetType(medalData.MedalType)
				medalRecord.SetName(medalData.MedalName)
				medalRecord.SetDescription(medalData.MedalDescription)
				medalRecord.SetAttr(medalData.MedalAttr)
				if txErr = txApp.Save(medalRecord); txErr != nil {
					logger.Error("更新勋章失败", slog.Any("medal", medalRecord), slog.Any("err", txErr))
					return txErr
				}
				updatedCount++
			}
		}

		// 删除待删除的勋章
		var res sql.Result
		if res, txErr = txApp.DB().Delete(model.DbNameMedals, dbx.Like(model.MedalsFieldMedalId, "pending_delete_")).Execute(); txErr != nil {
			logger.Error("删除勋章失败", slog.Any("err", txErr))
			return txErr
		}
		if deletedCount, txErr = res.RowsAffected(); txErr != nil {
			logger.Error("获取删除勋章数量失败", slog.Any("err", txErr))
			return txErr
		}

		return nil
	}); err != nil {
		return event.InternalServerError("同步勋章失败", err)
	}

	logger.Info("同步所有勋章完成",
		slog.Int("synced_count", len(medals)),
		slog.Int64("created_count", createdCount),
		slog.Int64("updated_count", updatedCount),
		slog.Int64("deleted_count", deletedCount),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"synced_count":  len(medals),
		"created_count": createdCount,
		"updated_count": updatedCount,
		"deleted_count": deletedCount,
	})
}

// SyncSingleMedal 同步单个勋章
func (controller *MedalController) SyncSingleMedal(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("sync_single_medal")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

	resp, err := controller.fishPiSdk.PostMedalAdminDetail(medalId)
	if err != nil {
		logger.Error("查询勋章详情失败", slog.Any("err", err), slog.String("medal_id", medalId))
		return event.InternalServerError("查询勋章详情失败", err)
	}
	if resp.Code != 0 {
		logger.Error("查询勋章详情失败", slog.Any("resp", resp), slog.String("medal_id", medalId))
		return event.InternalServerError("查询勋章详情失败: "+resp.Msg, nil)
	}

	medalData := resp.Data

	medal := new(model.Medal)
	created := false
	if err = event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: medalId,
	}).One(medal); err != nil {
		logger.Debug("查询本地勋章失败，创建新勋章", slog.String("medal_id", medalId))

		var medalCollection *core.Collection
		if medalCollection, err = event.App.FindCollectionByNameOrId(model.DbNameMedals); err != nil {
			logger.Error("获取勋章集合失败", slog.Any("err", err))
			return event.InternalServerError("获取勋章集合失败", err)
		}

		medal = model.NewMedalFromCollection(medalCollection)
		medal.SetOId(medalData.OId)
		created = true
	}

	medal.SetMedalId(medalData.MedalId)
	medal.SetType(medalData.MedalType)
	medal.SetName(medalData.MedalName)
	medal.SetDescription(medalData.MedalDescription)
	medal.SetAttr(medalData.MedalAttr)

	if err = event.App.Save(medal); err != nil {
		logger.Error("保存勋章失败", slog.Any("medal", medal), slog.Any("err", err))
		return event.InternalServerError("保存勋章失败", err)
	}

	logger.Info("同步单个勋章完成", slog.Any("medal", medal), slog.Bool("created", created))

	return event.JSON(http.StatusOK, map[string]any{
		"created": created,
		"medal":   medal,
	})
}

// SyncAllMedalOwners 同步所有勋章的拥有者
func (controller *MedalController) SyncAllMedalOwners(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("sync_all_medal_owners")

	// 获取所有本地勋章
	var medals []*model.Medal
	if err := event.App.RecordQuery(model.DbNameMedals).All(&medals); err != nil {
		logger.Error("查询本地勋章列表失败", slog.Any("err", err))
		return event.InternalServerError("查询本地勋章列表失败", err)
	}

	var (
		totalCreated int64
		totalUpdated int64
		totalDeleted int64
		syncedCount  int
		skippedCount int
	)

	for _, medal := range medals {
		medalId := medal.MedalId()
		if medalId == "" {
			skippedCount++
			continue
		}

		result, err := controller.syncSingleMedalOwners(event.App, medalId)
		if err != nil {
			logger.Warn("同步勋章拥有者失败", slog.Any("err", err), slog.String("medal_id", medalId))
			continue
		}

		totalCreated += result.Created
		totalUpdated += result.Updated
		totalDeleted += result.Deleted
		syncedCount++

		// 避免请求过于频繁
		time.Sleep(500 * time.Millisecond)
	}

	logger.Info("同步所有勋章拥有者完成",
		slog.Int("synced_count", syncedCount),
		slog.Int("skipped_count", skippedCount),
		slog.Int64("total_created", totalCreated),
		slog.Int64("total_updated", totalUpdated),
		slog.Int64("total_deleted", totalDeleted),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"synced_count":  syncedCount,
		"skipped_count": skippedCount,
		"total_created": totalCreated,
		"total_updated": totalUpdated,
		"total_deleted": totalDeleted,
	})
}

type syncOwnersResult struct {
	Created int64
	Updated int64
	Deleted int64
}

func (controller *MedalController) syncSingleMedalOwners(app core.App, medalId string) (*syncOwnersResult, error) {
	logger := controller.makeActionLogger("sync_medal_owners_internal").With(slog.String("medal_id", medalId))

	// 先查找本地勋章记录，获取其 PocketBase ID
	localMedal := new(model.Medal)
	if err := app.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: medalId,
	}).One(localMedal); err != nil {
		logger.Error("本地勋章不存在，请先同步勋章", slog.Any("err", err))
		return nil, fmt.Errorf("本地勋章不存在: %s", medalId)
	}
	localMedalRecordId := localMedal.Id

	const pageSize = 100
	var page = 1
	var allOwners []*types2.MedalOwner

	for {
		resp, err := controller.fishPiSdk.PostMedalAdminOwners(medalId, page, pageSize)
		if err != nil {
			logger.Error("查询勋章拥有者失败", slog.Any("err", err))
			return nil, err
		}
		if resp.Code != 0 {
			logger.Error("查询勋章拥有者失败", slog.Any("resp", resp))
			return nil, fmt.Errorf("查询勋章拥有者失败: %s", resp.Msg)
		}
		if resp.Data == nil || len(resp.Data.Items) == 0 {
			break
		}

		allOwners = append(allOwners, resp.Data.Items...)

		if len(resp.Data.Items) < pageSize {
			break
		}

		page++
		time.Sleep(300 * time.Millisecond)
	}

	ownerCollection, err := app.FindCollectionByNameOrId(model.DbNameMedalOwners)
	if err != nil {
		return nil, err
	}

	userCollection, err := app.FindCollectionByNameOrId(model.DbNameUsers)
	if err != nil {
		return nil, err
	}

	var result syncOwnersResult

	if err = app.RunInTransaction(func(txApp core.App) error {
		// 获取该勋章现有的所有拥有者记录ID
		var existingOwners []*model.MedalOwner
		if txErr := txApp.RecordQuery(model.DbNameMedalOwners).Where(dbx.HashExp{
			model.MedalOwnersFieldMedalId: localMedalRecordId,
		}).All(&existingOwners); txErr != nil {
			logger.Warn("查询现有拥有者失败", slog.Any("err", txErr))
		}

		// 构建现有拥有者映射 (userId -> record)
		existingMap := make(map[string]*model.MedalOwner)
		for _, owner := range existingOwners {
			existingMap[owner.UserId()] = owner
		}

		// 记录本次同步处理的用户ID
		processedUserIds := make(map[string]bool)

		for _, ownerData := range allOwners {
			// 根据鱼排用户ID查找本地用户记录
			localUser := new(model.User)
			if txErr := txApp.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
				model.UsersFieldOId: ownerData.UserId,
			}).One(localUser); txErr != nil {
				// 本地用户不存在，创建一个
				localUser = model.NewUserFromCollection(userCollection)
				localUser.SetOId(ownerData.UserId)
				localUser.SetName(ownerData.UserName)
				localUser.SetNickname(ownerData.UserName)
				localUser.Set(model.UsersFieldEmail, fmt.Sprintf("%s@fishpi.cn", ownerData.UserId))
				localUser.SetRaw("password", fmt.Sprintf("temp_%s_%d", ownerData.UserId, time.Now().UnixNano()))
				if txErr = txApp.Save(localUser); txErr != nil {
					logger.Warn("创建本地用户失败", slog.Any("err", txErr), slog.String("user_id", ownerData.UserId))
					continue
				}
			}
			localUserRecordId := localUser.Id
			processedUserIds[localUserRecordId] = true

			// 查找或创建勋章拥有者记录
			if existingOwner, exists := existingMap[localUserRecordId]; exists {
				// 更新记录
				existingOwner.SetDisplay(ownerData.Display)
				existingOwner.SetDisplayOrder(ownerData.DisplayOrder)
				existingOwner.SetData(ownerData.Data)
				if ownerData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(ownerData.ExpireTime)); parseErr == nil {
						existingOwner.SetExpired(expired)
					}
				}
				if txErr := txApp.Save(existingOwner); txErr != nil {
					logger.Warn("更新拥有者记录失败", slog.Any("err", txErr))
					continue
				}
				result.Updated++
			} else {
				// 创建新记录
				ownerRecord := model.NewMedalOwnerFromCollection(ownerCollection)
				ownerRecord.SetMedalId(localMedalRecordId)
				ownerRecord.SetUserId(localUserRecordId)
				ownerRecord.SetDisplay(ownerData.Display)
				ownerRecord.SetDisplayOrder(ownerData.DisplayOrder)
				ownerRecord.SetData(ownerData.Data)
				if ownerData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(ownerData.ExpireTime)); parseErr == nil {
						ownerRecord.SetExpired(expired)
					}
				}
				if txErr := txApp.Save(ownerRecord); txErr != nil {
					logger.Warn("创建拥有者记录失败", slog.Any("err", txErr))
					continue
				}
				result.Created++
			}
		}

		// 删除不在本次同步中的记录
		for userRecordId, owner := range existingMap {
			if !processedUserIds[userRecordId] {
				if txErr := txApp.Delete(owner); txErr != nil {
					logger.Warn("删除拥有者记录失败", slog.Any("err", txErr))
					continue
				}
				result.Deleted++
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

// SyncSingleMedalOwners 同步单个勋章的拥有者
func (controller *MedalController) SyncSingleMedalOwners(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("sync_single_medal_owners")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

	result, err := controller.syncSingleMedalOwners(event.App, medalId)
	if err != nil {
		logger.Error("同步勋章拥有者失败", slog.Any("err", err))
		return event.InternalServerError("同步勋章拥有者失败", err)
	}

	logger.Info("同步单个勋章拥有者完成",
		slog.Int64("created", result.Created),
		slog.Int64("updated", result.Updated),
		slog.Int64("deleted", result.Deleted),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"created": result.Created,
		"updated": result.Updated,
		"deleted": result.Deleted,
	})
}

// SyncUserMedals 同步某用户的所有勋章
func (controller *MedalController) SyncUserMedals(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("sync_user_medals")

	userId := event.Request.PathValue("userId")
	if userId == "" {
		return event.BadRequestError("缺少用户ID", nil)
	}

	// 先查找本地用户记录
	localUser := new(model.User)
	if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
		model.UsersFieldOId: userId,
	}).One(localUser); err != nil {
		logger.Error("本地用户不存在", slog.Any("err", err), slog.String("user_id", userId))
		return event.BadRequestError("本地用户不存在，请先同步用户", nil)
	}
	localUserRecordId := localUser.Id

	// 查询用户的所有勋章
	resp, err := controller.fishPiSdk.PostMedalAdminUserMedals(&types2.PostMedalAdminUserMedalsRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error("查询用户勋章失败", slog.Any("err", err), slog.String("user_id", userId))
		return event.InternalServerError("查询用户勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("查询用户勋章失败", slog.Any("resp", resp), slog.String("user_id", userId))
		return event.InternalServerError("查询用户勋章失败: "+resp.Msg, nil)
	}

	medals := resp.Data
	ownerCollection, err := event.App.FindCollectionByNameOrId(model.DbNameMedalOwners)
	if err != nil {
		logger.Error("获取勋章拥有者集合失败", slog.Any("err", err))
		return event.InternalServerError("获取勋章拥有者集合失败", err)
	}

	var (
		createdCount int64
		updatedCount int64
		deletedCount int64
		skippedCount int64
	)

	if err = event.App.RunInTransaction(func(txApp core.App) error {
		// 获取该用户现有的所有勋章记录
		var existingOwners []*model.MedalOwner
		if txErr := txApp.RecordQuery(model.DbNameMedalOwners).Where(dbx.HashExp{
			model.MedalOwnersFieldUserId: localUserRecordId,
		}).All(&existingOwners); txErr != nil {
			logger.Warn("查询现有勋章失败", slog.Any("err", txErr))
		}

		// 构建现有勋章映射 (medalRecordId -> record)
		existingMap := make(map[string]*model.MedalOwner)
		for _, owner := range existingOwners {
			existingMap[owner.MedalId()] = owner
		}

		// 记录本次同步处理的勋章ID
		processedMedalIds := make(map[string]bool)

		for _, medalData := range medals {
			// 根据鱼排勋章ID查找本地勋章记录
			localMedal := new(model.Medal)
			if txErr := txApp.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
				model.MedalsFieldMedalId: medalData.MedalId,
			}).One(localMedal); txErr != nil {
				logger.Warn("本地勋章不存在，跳过", slog.String("medal_id", medalData.MedalId))
				skippedCount++
				continue
			}
			localMedalRecordId := localMedal.Id
			processedMedalIds[localMedalRecordId] = true

			// 查找或创建勋章拥有者记录
			if existingOwner, exists := existingMap[localMedalRecordId]; exists {
				// 更新记录
				existingOwner.SetDisplay(medalData.Display)
				existingOwner.SetDisplayOrder(medalData.DisplayOrder)
				existingOwner.SetData(medalData.Data)
				if medalData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(medalData.ExpireTime)); parseErr == nil {
						existingOwner.SetExpired(expired)
					}
				}
				if txErr := txApp.Save(existingOwner); txErr != nil {
					logger.Warn("更新勋章拥有记录失败", slog.Any("err", txErr))
					continue
				}
				updatedCount++
			} else {
				// 创建新记录
				ownerRecord := model.NewMedalOwnerFromCollection(ownerCollection)
				ownerRecord.SetMedalId(localMedalRecordId)
				ownerRecord.SetUserId(localUserRecordId)
				ownerRecord.SetDisplay(medalData.Display)
				ownerRecord.SetDisplayOrder(medalData.DisplayOrder)
				ownerRecord.SetData(medalData.Data)
				if medalData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(medalData.ExpireTime)); parseErr == nil {
						ownerRecord.SetExpired(expired)
					}
				}
				if txErr := txApp.Save(ownerRecord); txErr != nil {
					logger.Warn("创建勋章拥有记录失败", slog.Any("err", txErr))
					continue
				}
				createdCount++
			}
		}

		// 删除不在本次同步中的记录
		for medalRecordId, owner := range existingMap {
			if !processedMedalIds[medalRecordId] {
				if txErr := txApp.Delete(owner); txErr != nil {
					logger.Warn("删除勋章拥有记录失败", slog.Any("err", txErr))
					continue
				}
				deletedCount++
			}
		}

		return nil
	}); err != nil {
		logger.Error("同步用户勋章失败", slog.Any("err", err))
		return event.InternalServerError("同步用户勋章失败", err)
	}

	logger.Info("同步用户勋章完成",
		slog.String("user_id", userId),
		slog.Int("total_medals", len(medals)),
		slog.Int64("created", createdCount),
		slog.Int64("updated", updatedCount),
		slog.Int64("deleted", deletedCount),
		slog.Int64("skipped", skippedCount),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"total_medals": len(medals),
		"created":      createdCount,
		"updated":      updatedCount,
		"deleted":      deletedCount,
		"skipped":      skippedCount,
	})
}

// GetMedalOwners 获取勋章拥有者列表
func (controller *MedalController) GetMedalOwners(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_medal_owners")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

	// 先查找本地勋章记录
	localMedal := new(model.Medal)
	if err := event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: medalId,
	}).One(localMedal); err != nil {
		logger.Error("本地勋章不存在", slog.Any("err", err), slog.String("medal_id", medalId))
		return event.NotFoundError("本地勋章不存在", err)
	}
	localMedalRecordId := localMedal.Id

	page, _ := strconv.Atoi(event.Request.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(event.Request.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var owners []*model.MedalOwner
	var total int

	// 先查询总数
	if err := event.App.RecordQuery(model.DbNameMedalOwners).Where(dbx.HashExp{
		model.MedalOwnersFieldMedalId: localMedalRecordId,
	}).Select("count(*)").Row(&total); err != nil {
		logger.Error("查询勋章拥有者总数失败", slog.Any("err", err))
		return event.InternalServerError("查询勋章拥有者总数失败", err)
	}

	// 再查询列表
	if err := event.App.RecordQuery(model.DbNameMedalOwners).Where(dbx.HashExp{
		model.MedalOwnersFieldMedalId: localMedalRecordId,
	}).OrderBy(fmt.Sprintf("%s DESC", model.MedalOwnersFieldCreated)).Limit(int64(pageSize)).Offset(int64((page - 1) * pageSize)).All(&owners); err != nil {
		logger.Error("查询勋章拥有者列表失败", slog.Any("err", err))
		return event.InternalServerError("查询勋章拥有者列表失败", err)
	}

	// 查询关联的用户信息
	userIds := make([]string, 0, len(owners))
	for _, owner := range owners {
		if userId := owner.UserId(); userId != "" {
			userIds = append(userIds, userId)
		}
	}

	// 批量查询用户
	usersMap := make(map[string]*model.User)
	if len(userIds) > 0 {
		var users []*model.User
		// 使用 IN 查询批量获取用户
		userIdsAny := make([]any, len(userIds))
		for i, id := range userIds {
			userIdsAny[i] = id
		}
		if err := event.App.RecordQuery(model.DbNameUsers).Where(
			dbx.In(model.CommonFieldId, userIdsAny...),
		).All(&users); err != nil {
			logger.Warn("批量查询用户失败", slog.Any("err", err))
		} else {
			for _, user := range users {
				usersMap[user.Id] = user
			}
		}
	}

	// 构建包含用户信息的响应
	type ownerWithUser struct {
		Id           string `json:"id"`
		MedalId      string `json:"medalId"`
		UserId       string `json:"userId"`
		Display      bool   `json:"display"`
		DisplayOrder int    `json:"displayOrder"`
		Data         string `json:"data"`
		Expired      string `json:"expired"`
		Created      string `json:"created"`
		Updated      string `json:"updated"`
		Expand       *struct {
			UserId *struct {
				Id       string `json:"id"`
				Name     string `json:"name"`
				Nickname string `json:"nickname"`
				Avatar   string `json:"avatar"`
			} `json:"userId"`
		} `json:"expand,omitempty"`
	}

	items := make([]*ownerWithUser, 0, len(owners))
	for _, owner := range owners {
		item := &ownerWithUser{
			Id:           owner.Id,
			MedalId:      owner.MedalId(),
			UserId:       owner.UserId(),
			Display:      owner.Display(),
			DisplayOrder: owner.DisplayOrder(),
			Data:         owner.Data(),
			Expired:      owner.Expired().String(),
			Created:      owner.Created().String(),
			Updated:      owner.Updated().String(),
		}
		if user, exists := usersMap[owner.UserId()]; exists {
			item.Expand = &struct {
				UserId *struct {
					Id       string `json:"id"`
					Name     string `json:"name"`
					Nickname string `json:"nickname"`
					Avatar   string `json:"avatar"`
				} `json:"userId"`
			}{
				UserId: &struct {
					Id       string `json:"id"`
					Name     string `json:"name"`
					Nickname string `json:"nickname"`
					Avatar   string `json:"avatar"`
				}{
					Id:       user.Id,
					Name:     user.Name(),
					Nickname: user.Nickname(),
					Avatar:   user.Avatar(),
				},
			}
		}
		items = append(items, item)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items":      items,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + pageSize - 1) / pageSize,
	})
}

// GrantMedal 给用户授予勋章
func (controller *MedalController) GrantMedal(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("grant_medal")

	var req struct {
		UserId     string `json:"userId"`
		MedalId    string `json:"medalId"`
		ExpireTime int64  `json:"expireTime"` // 毫秒时间戳，0表示永不过期
		Data       string `json:"data"`
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if req.UserId == "" || req.MedalId == "" {
		return event.BadRequestError("用户ID和勋章ID不能为空", nil)
	}

	// 查找本地勋章记录
	localMedal := new(model.Medal)
	if err := event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: req.MedalId,
	}).One(localMedal); err != nil {
		logger.Error("本地勋章不存在", slog.Any("err", err), slog.String("medal_id", req.MedalId))
		return event.BadRequestError("本地勋章不存在，请先同步勋章", nil)
	}
	localMedalRecordId := localMedal.Id

	// 查找本地用户记录
	localUser := new(model.User)
	if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
		model.UsersFieldOId: req.UserId,
	}).One(localUser); err != nil {
		logger.Error("本地用户不存在", slog.Any("err", err), slog.String("user_id", req.UserId))
		return event.BadRequestError("本地用户不存在，请先同步用户", nil)
	}
	localUserRecordId := localUser.Id

	// 调用鱼排接口授予勋章
	resp, err := controller.fishPiSdk.PostMedalAdminGrant(req.UserId, req.MedalId, req.ExpireTime, req.Data)
	if err != nil {
		logger.Error("授予勋章失败", slog.Any("err", err))
		return event.InternalServerError("授予勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("授予勋章失败", slog.Any("resp", resp))
		return event.InternalServerError("授予勋章失败: "+resp.Msg, nil)
	}

	// 保存到本地数据库
	ownerCollection, err := event.App.FindCollectionByNameOrId(model.DbNameMedalOwners)
	if err != nil {
		logger.Error("获取勋章拥有者集合失败", slog.Any("err", err))
		return event.InternalServerError("获取勋章拥有者集合失败", err)
	}

	ownerRecord := new(model.MedalOwner)
	created := false
	if err = event.App.RecordQuery(model.DbNameMedalOwners).Where(dbx.And(
		dbx.HashExp{model.MedalOwnersFieldMedalId: localMedalRecordId},
		dbx.HashExp{model.MedalOwnersFieldUserId: localUserRecordId},
	)).One(ownerRecord); err != nil {
		ownerRecord = model.NewMedalOwnerFromCollection(ownerCollection)
		ownerRecord.SetMedalId(localMedalRecordId)
		ownerRecord.SetUserId(localUserRecordId)
		created = true
	}

	ownerRecord.SetData(req.Data)
	ownerRecord.SetDisplay(true)
	if req.ExpireTime > 0 {
		if expired, parseErr := types.ParseDateTime(time.UnixMilli(req.ExpireTime)); parseErr == nil {
			ownerRecord.SetExpired(expired)
		}
	}

	if err = event.App.Save(ownerRecord); err != nil {
		logger.Error("保存勋章拥有者记录失败", slog.Any("err", err))
		return event.InternalServerError("保存勋章拥有者记录失败", err)
	}

	logger.Info("授予勋章成功",
		slog.String("user_id", req.UserId),
		slog.String("medal_id", req.MedalId),
		slog.Bool("created", created),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"success": true,
		"created": created,
	})
}

// RevokeMedal 撤销用户勋章
func (controller *MedalController) RevokeMedal(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("revoke_medal")

	var req struct {
		UserId  string `json:"userId"`
		MedalId string `json:"medalId"`
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if req.UserId == "" || req.MedalId == "" {
		return event.BadRequestError("用户ID和勋章ID不能为空", nil)
	}

	// 查找本地勋章记录
	localMedal := new(model.Medal)
	if err := event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: req.MedalId,
	}).One(localMedal); err != nil {
		logger.Warn("本地勋章不存在", slog.Any("err", err), slog.String("medal_id", req.MedalId))
	}

	// 查找本地用户记录
	localUser := new(model.User)
	if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
		model.UsersFieldOId: req.UserId,
	}).One(localUser); err != nil {
		logger.Warn("本地用户不存在", slog.Any("err", err), slog.String("user_id", req.UserId))
	}

	// 调用鱼排接口撤销勋章
	resp, err := controller.fishPiSdk.PostMedalAdminRevoke(req.UserId, req.MedalId)
	if err != nil {
		logger.Error("撤销勋章失败", slog.Any("err", err))
		return event.InternalServerError("撤销勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("撤销勋章失败", slog.Any("resp", resp))
		return event.InternalServerError("撤销勋章失败: "+resp.Msg, nil)
	}

	// 删除本地记录
	if localMedal.Id != "" && localUser.Id != "" {
		if _, err = event.App.DB().Delete(model.DbNameMedalOwners, dbx.And(
			dbx.HashExp{model.MedalOwnersFieldMedalId: localMedal.Id},
			dbx.HashExp{model.MedalOwnersFieldUserId: localUser.Id},
		)).Execute(); err != nil {
			logger.Warn("删除本地勋章拥有者记录失败", slog.Any("err", err))
		}
	}

	logger.Info("撤销勋章成功",
		slog.String("user_id", req.UserId),
		slog.String("medal_id", req.MedalId),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}

// Search 搜索勋章
func (controller *MedalController) Search(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("search")

	keyword := event.Request.URL.Query().Get("keyword")
	if keyword == "" {
		return event.BadRequestError("搜索关键词不能为空", nil)
	}

	// 从鱼排搜索
	resp, err := controller.fishPiSdk.PostMedalAdminSearch(keyword)
	if err != nil {
		logger.Error("搜索勋章失败", slog.Any("err", err))
		return event.InternalServerError("搜索勋章失败", err)
	}
	if resp.Code != 0 {
		logger.Error("搜索勋章失败", slog.Any("resp", resp))
		return event.InternalServerError("搜索勋章失败: "+resp.Msg, nil)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": resp.Data,
	})
}

// SearchUsers 搜索用户
func (controller *MedalController) SearchUsers(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("search_users")

	keyword := event.Request.URL.Query().Get("keyword")
	if keyword == "" {
		return event.BadRequestError("搜索关键词不能为空", nil)
	}

	var users []*model.User
	// 按用户名或昵称模糊搜索，或者精确匹配oId
	if err := event.App.RecordQuery(model.DbNameUsers).
		Where(dbx.Or(
			dbx.NewExp(fmt.Sprintf("%s LIKE {:keyword}", model.UsersFieldName), dbx.Params{"keyword": "%" + keyword + "%"}),
			dbx.NewExp(fmt.Sprintf("%s LIKE {:keyword}", model.UsersFieldNickname), dbx.Params{"keyword": "%" + keyword + "%"}),
			dbx.NewExp(fmt.Sprintf("%s = {:oId}", model.UsersFieldOId), dbx.Params{"oId": keyword}),
		)).
		Limit(50).
		All(&users); err != nil {
		logger.Error("搜索用户失败", slog.Any("err", err))
		return event.InternalServerError("搜索用户失败", err)
	}

	type userItem struct {
		Id       string `json:"id"`
		OId      string `json:"oId"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	items := make([]*userItem, 0, len(users))
	for _, user := range users {
		items = append(items, &userItem{
			Id:       user.Id,
			OId:      user.OId(),
			Name:     user.Name(),
			Nickname: user.Nickname(),
			Avatar:   user.Avatar(),
		})
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": items,
	})
}

// GetActivities 获取活动列表
func (controller *MedalController) GetActivities(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_activities")

	var activities []*model.Activity
	if err := event.App.RecordQuery(model.DbNameActivities).
		OrderBy(fmt.Sprintf("%s DESC", model.ActivitiesFieldCreated)).
		Limit(100).
		All(&activities); err != nil {
		logger.Error("获取活动列表失败", slog.Any("err", err))
		return event.InternalServerError("获取活动列表失败", err)
	}

	type activityItem struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Slug   string `json:"slug"`
		VoteId string `json:"voteId"`
	}

	items := make([]*activityItem, 0, len(activities))
	for _, activity := range activities {
		items = append(items, &activityItem{
			Id:     activity.Id,
			Name:   activity.GetName(),
			Slug:   activity.GetString(model.ActivitiesFieldSlug),
			VoteId: activity.GetString(model.ActivitiesFieldVoteId),
		})
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": items,
	})
}

// GetActivityParticipants 获取活动参与者
func (controller *MedalController) GetActivityParticipants(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_activity_participants")

	activityId := event.Request.PathValue("activityId")
	if activityId == "" {
		return event.BadRequestError("缺少活动ID", nil)
	}

	// 查找活动
	activity := new(model.Activity)
	if err := event.App.RecordQuery(model.DbNameActivities).
		Where(dbx.HashExp{model.CommonFieldId: activityId}).
		One(activity); err != nil {
		logger.Error("活动不存在", slog.Any("err", err))
		return event.NotFoundError("活动不存在", err)
	}

	voteId := activity.GetString(model.ActivitiesFieldVoteId)
	if voteId == "" {
		return event.JSON(http.StatusOK, map[string]any{
			"items": []any{},
		})
	}

	// 查询投票日志中的用户
	var voteJuryLogs []*model.VoteJuryLog
	if err := event.App.RecordQuery(model.DbNameVoteJuryLogs).
		Where(dbx.HashExp{model.VoteJuryLogFieldVoteId: voteId}).
		All(&voteJuryLogs); err != nil {
		logger.Warn("查询投票日志失败", slog.Any("err", err))
	}

	// 收集所有参与的用户ID (toUserId 是被投票的参与者)
	userIdSet := make(map[string]bool)
	for _, log := range voteJuryLogs {
		if toUserId := log.ToUserId(); toUserId != "" {
			userIdSet[toUserId] = true
		}
	}

	// 批量查询用户
	userIds := make([]any, 0, len(userIdSet))
	for userId := range userIdSet {
		userIds = append(userIds, userId)
	}

	type userItem struct {
		Id       string `json:"id"`
		OId      string `json:"oId"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	items := make([]*userItem, 0)
	if len(userIds) > 0 {
		var users []*model.User
		if err := event.App.RecordQuery(model.DbNameUsers).
			Where(dbx.In(model.CommonFieldId, userIds...)).
			All(&users); err != nil {
			logger.Warn("批量查询用户失败", slog.Any("err", err))
		} else {
			for _, user := range users {
				items = append(items, &userItem{
					Id:       user.Id,
					OId:      user.OId(),
					Name:     user.Name(),
					Nickname: user.Nickname(),
					Avatar:   user.Avatar(),
				})
			}
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": items,
	})
}

// GetVoteJuryMembers 获取评审团成员
func (controller *MedalController) GetVoteJuryMembers(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_vote_jury_members")

	voteId := event.Request.PathValue("voteId")
	if voteId == "" {
		return event.BadRequestError("缺少投票ID", nil)
	}

	// 查询评审团成员
	var juryUsers []*model.VoteJuryUser
	if err := event.App.RecordQuery(model.DbNameVoteJuryUsers).
		Where(dbx.HashExp{
			model.VoteJuryUserFieldVoteId: voteId,
			model.VoteJuryUserFieldStatus: string(model.JuryUserStatusApproved),
		}).
		All(&juryUsers); err != nil {
		logger.Error("查询评审团成员失败", slog.Any("err", err))
		return event.InternalServerError("查询评审团成员失败", err)
	}

	// 收集用户ID
	userIds := make([]any, 0, len(juryUsers))
	for _, juryUser := range juryUsers {
		if userId := juryUser.UserId(); userId != "" {
			userIds = append(userIds, userId)
		}
	}

	type userItem struct {
		Id       string `json:"id"`
		OId      string `json:"oId"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	items := make([]*userItem, 0)
	if len(userIds) > 0 {
		var users []*model.User
		if err := event.App.RecordQuery(model.DbNameUsers).
			Where(dbx.In(model.CommonFieldId, userIds...)).
			All(&users); err != nil {
			logger.Warn("批量查询用户失败", slog.Any("err", err))
		} else {
			for _, user := range users {
				items = append(items, &userItem{
					Id:       user.Id,
					OId:      user.OId(),
					Name:     user.Name(),
					Nickname: user.Nickname(),
					Avatar:   user.Avatar(),
				})
			}
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": items,
	})
}

// GrantMedalBatch 批量授予勋章
func (controller *MedalController) GrantMedalBatch(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("grant_medal_batch")

	var req struct {
		UserIds    []string `json:"userIds"`    // 用户oId列表
		MedalId    string   `json:"medalId"`    // 勋章ID (鱼排的medalId)
		ExpireTime int64    `json:"expireTime"` // 毫秒时间戳，0表示永不过期
		Data       string   `json:"data"`
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if len(req.UserIds) == 0 || req.MedalId == "" {
		return event.BadRequestError("用户ID列表和勋章ID不能为空", nil)
	}

	// 查找本地勋章记录
	localMedal := new(model.Medal)
	if err := event.App.RecordQuery(model.DbNameMedals).Where(dbx.HashExp{
		model.MedalsFieldMedalId: req.MedalId,
	}).One(localMedal); err != nil {
		logger.Error("本地勋章不存在", slog.Any("err", err), slog.String("medal_id", req.MedalId))
		return event.BadRequestError("本地勋章不存在，请先同步勋章", nil)
	}
	localMedalRecordId := localMedal.Id

	ownerCollection, err := event.App.FindCollectionByNameOrId(model.DbNameMedalOwners)
	if err != nil {
		logger.Error("获取勋章拥有者集合失败", slog.Any("err", err))
		return event.InternalServerError("获取勋章拥有者集合失败", err)
	}

	// 检查是否是开发模式
	isDev := controller.app.IsDev()

	var (
		successCount int
		failedCount  int
		skippedCount int
		results      []map[string]any
	)

	for _, userOId := range req.UserIds {
		result := map[string]any{
			"userId":  userOId,
			"success": false,
		}

		// 查找本地用户记录
		localUser := new(model.User)
		if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
			model.UsersFieldOId: userOId,
		}).One(localUser); err != nil {
			logger.Warn("本地用户不存在", slog.String("user_id", userOId))
			result["error"] = "本地用户不存在"
			failedCount++
			results = append(results, result)
			continue
		}
		localUserRecordId := localUser.Id

		// 检查是否已拥有该勋章
		existingOwner := new(model.MedalOwner)
		if err := event.App.RecordQuery(model.DbNameMedalOwners).Where(dbx.And(
			dbx.HashExp{model.MedalOwnersFieldMedalId: localMedalRecordId},
			dbx.HashExp{model.MedalOwnersFieldUserId: localUserRecordId},
		)).One(existingOwner); err == nil {
			logger.Info("用户已拥有该勋章，跳过", slog.String("user_id", userOId))
			result["error"] = "用户已拥有该勋章"
			skippedCount++
			results = append(results, result)
			continue
		}

		// 开发模式下不实际发放勋章
		if isDev {
			logger.Warn("[DEV] 开发模式，跳过实际勋章发放",
				slog.String("user_id", userOId),
				slog.String("medal_id", req.MedalId),
				slog.String("medal_name", localMedal.Name()),
			)
		} else {
			// 调用鱼排接口授予勋章
			resp, err := controller.fishPiSdk.PostMedalAdminGrant(userOId, req.MedalId, req.ExpireTime, req.Data)
			if err != nil {
				logger.Error("授予勋章失败", slog.Any("err", err), slog.String("user_id", userOId))
				result["error"] = "调用鱼排接口失败"
				failedCount++
				results = append(results, result)
				continue
			}
			if resp.Code != 0 {
				logger.Error("授予勋章失败", slog.Any("resp", resp), slog.String("user_id", userOId))
				result["error"] = resp.Msg
				failedCount++
				results = append(results, result)
				continue
			}
		}

		// 保存到本地数据库
		ownerRecord := model.NewMedalOwnerFromCollection(ownerCollection)
		ownerRecord.SetMedalId(localMedalRecordId)
		ownerRecord.SetUserId(localUserRecordId)
		ownerRecord.SetData(req.Data)
		ownerRecord.SetDisplay(true)
		if req.ExpireTime > 0 {
			if expired, parseErr := types.ParseDateTime(time.UnixMilli(req.ExpireTime)); parseErr == nil {
				ownerRecord.SetExpired(expired)
			}
		}

		if err := event.App.Save(ownerRecord); err != nil {
			logger.Error("保存勋章拥有者记录失败", slog.Any("err", err), slog.String("user_id", userOId))
			result["error"] = "保存本地记录失败"
			failedCount++
			results = append(results, result)
			continue
		}

		result["success"] = true
		successCount++
		results = append(results, result)
	}

	logger.Info("批量授予勋章完成",
		slog.Int("total", len(req.UserIds)),
		slog.Int("success", successCount),
		slog.Int("failed", failedCount),
		slog.Int("skipped", skippedCount),
		slog.Bool("dev_mode", isDev),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"total":    len(req.UserIds),
		"success":  successCount,
		"failed":   failedCount,
		"skipped":  skippedCount,
		"dev_mode": isDev,
		"results":  results,
	})
}
