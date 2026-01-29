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
	group.POST("/revoke", controller.RevokeMedal)

	// 搜索勋章
	group.GET("/search", controller.Search)
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

	var result syncOwnersResult

	if err = app.RunInTransaction(func(txApp core.App) error {
		// 标记现有记录为待删除
		_, txErr := txApp.DB().Update(model.DbNameMedalOwners, dbx.Params{
			model.MedalOwnersFieldUserId: dbx.NewExp(fmt.Sprintf("'pending_delete_' || `%s`", model.MedalOwnersFieldUserId)),
		}, dbx.HashExp{
			model.MedalOwnersFieldMedalId: medalId,
		}).Execute()
		if txErr != nil {
			return txErr
		}

		for _, ownerData := range allOwners {
			ownerRecord := new(model.MedalOwner)
			if txErr = txApp.RecordQuery(model.DbNameMedalOwners).Where(dbx.And(
				dbx.HashExp{model.MedalOwnersFieldMedalId: medalId},
				dbx.Or(
					dbx.HashExp{model.MedalOwnersFieldUserId: ownerData.UserId},
					dbx.HashExp{model.MedalOwnersFieldUserId: "pending_delete_" + ownerData.UserId},
				),
			)).One(ownerRecord); txErr != nil {
				// 创建新记录
				ownerRecord = model.NewMedalOwnerFromCollection(ownerCollection)
				ownerRecord.SetMedalId(medalId)
				ownerRecord.SetUserId(ownerData.UserId)
				ownerRecord.SetDisplay(ownerData.Display)
				ownerRecord.SetDisplayOrder(ownerData.DisplayOrder)
				ownerRecord.SetData(ownerData.Data)
				if ownerData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(ownerData.ExpireTime)); parseErr == nil {
						ownerRecord.SetExpired(expired)
					}
				}
				if txErr = txApp.Save(ownerRecord); txErr != nil {
					return txErr
				}
				result.Created++
			} else {
				// 更新记录
				ownerRecord.SetUserId(ownerData.UserId)
				ownerRecord.SetDisplay(ownerData.Display)
				ownerRecord.SetDisplayOrder(ownerData.DisplayOrder)
				ownerRecord.SetData(ownerData.Data)
				if ownerData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(ownerData.ExpireTime)); parseErr == nil {
						ownerRecord.SetExpired(expired)
					}
				}
				if txErr = txApp.Save(ownerRecord); txErr != nil {
					return txErr
				}
				result.Updated++
			}
		}

		// 删除待删除的记录
		var res sql.Result
		if res, txErr = txApp.DB().Delete(model.DbNameMedalOwners, dbx.And(
			dbx.HashExp{model.MedalOwnersFieldMedalId: medalId},
			dbx.Like(model.MedalOwnersFieldUserId, "pending_delete_"),
		)).Execute(); txErr != nil {
			return txErr
		}
		if result.Deleted, txErr = res.RowsAffected(); txErr != nil {
			return txErr
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
	)

	if err = event.App.RunInTransaction(func(txApp core.App) error {
		// 标记该用户的所有勋章为待删除
		_, txErr := txApp.DB().Update(model.DbNameMedalOwners, dbx.Params{
			model.MedalOwnersFieldMedalId: dbx.NewExp(fmt.Sprintf("'pending_delete_' || `%s`", model.MedalOwnersFieldMedalId)),
		}, dbx.HashExp{
			model.MedalOwnersFieldUserId: userId,
		}).Execute()
		if txErr != nil {
			return txErr
		}

		for _, medalData := range medals {
			ownerRecord := new(model.MedalOwner)
			if txErr = txApp.RecordQuery(model.DbNameMedalOwners).Where(dbx.And(
				dbx.HashExp{model.MedalOwnersFieldUserId: userId},
				dbx.Or(
					dbx.HashExp{model.MedalOwnersFieldMedalId: medalData.MedalId},
					dbx.HashExp{model.MedalOwnersFieldMedalId: "pending_delete_" + medalData.MedalId},
				),
			)).One(ownerRecord); txErr != nil {
				// 创建新记录
				ownerRecord = model.NewMedalOwnerFromCollection(ownerCollection)
				ownerRecord.SetMedalId(medalData.MedalId)
				ownerRecord.SetUserId(userId)
				ownerRecord.SetDisplay(medalData.Display)
				ownerRecord.SetDisplayOrder(medalData.DisplayOrder)
				ownerRecord.SetData(medalData.Data)
				if medalData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(medalData.ExpireTime)); parseErr == nil {
						ownerRecord.SetExpired(expired)
					}
				}
				if txErr = txApp.Save(ownerRecord); txErr != nil {
					return txErr
				}
				createdCount++
			} else {
				// 更新记录
				ownerRecord.SetMedalId(medalData.MedalId)
				ownerRecord.SetDisplay(medalData.Display)
				ownerRecord.SetDisplayOrder(medalData.DisplayOrder)
				ownerRecord.SetData(medalData.Data)
				if medalData.ExpireTime > 0 {
					if expired, parseErr := types.ParseDateTime(time.UnixMilli(medalData.ExpireTime)); parseErr == nil {
						ownerRecord.SetExpired(expired)
					}
				}
				if txErr = txApp.Save(ownerRecord); txErr != nil {
					return txErr
				}
				updatedCount++
			}
		}

		// 删除待删除的记录
		var res sql.Result
		if res, txErr = txApp.DB().Delete(model.DbNameMedalOwners, dbx.And(
			dbx.HashExp{model.MedalOwnersFieldUserId: userId},
			dbx.Like(model.MedalOwnersFieldMedalId, "pending_delete_"),
		)).Execute(); txErr != nil {
			return txErr
		}
		if deletedCount, txErr = res.RowsAffected(); txErr != nil {
			return txErr
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
	)

	return event.JSON(http.StatusOK, map[string]any{
		"total_medals": len(medals),
		"created":      createdCount,
		"updated":      updatedCount,
		"deleted":      deletedCount,
	})
}

// GetMedalOwners 获取勋章拥有者列表
func (controller *MedalController) GetMedalOwners(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_medal_owners")

	medalId := event.Request.PathValue("medalId")
	if medalId == "" {
		return event.BadRequestError("缺少勋章ID", nil)
	}

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
		model.MedalOwnersFieldMedalId: medalId,
	}).Select("count(*)").Row(&total); err != nil {
		logger.Error("查询勋章拥有者总数失败", slog.Any("err", err))
		return event.InternalServerError("查询勋章拥有者总数失败", err)
	}

	// 再查询列表
	if err := event.App.RecordQuery(model.DbNameMedalOwners).Where(dbx.HashExp{
		model.MedalOwnersFieldMedalId: medalId,
	}).OrderBy(fmt.Sprintf("%s DESC", model.MedalOwnersFieldCreated)).Limit(int64(pageSize)).Offset(int64((page - 1) * pageSize)).All(&owners); err != nil {
		logger.Error("查询勋章拥有者列表失败", slog.Any("err", err))
		return event.InternalServerError("查询勋章拥有者列表失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items":      owners,
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
		dbx.HashExp{model.MedalOwnersFieldMedalId: req.MedalId},
		dbx.HashExp{model.MedalOwnersFieldUserId: req.UserId},
	)).One(ownerRecord); err != nil {
		ownerRecord = model.NewMedalOwnerFromCollection(ownerCollection)
		ownerRecord.SetMedalId(req.MedalId)
		ownerRecord.SetUserId(req.UserId)
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
	if _, err = event.App.DB().Delete(model.DbNameMedalOwners, dbx.And(
		dbx.HashExp{model.MedalOwnersFieldMedalId: req.MedalId},
		dbx.HashExp{model.MedalOwnersFieldUserId: req.UserId},
	)).Execute(); err != nil {
		logger.Warn("删除本地勋章拥有者记录失败", slog.Any("err", err))
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
