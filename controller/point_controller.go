package controller

import (
	"bless-activity/model"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

type PointController struct {
	*BaseController

	event  *core.ServeEvent
	group  *router.RouterGroup[*core.RequestEvent]
	app    core.App
	logger *slog.Logger
}

func NewPointController(event *core.ServeEvent, group *router.RouterGroup[*core.RequestEvent], base *BaseController) *PointController {
	logger := event.App.Logger().With(
		slog.String("controller", "point"),
	)

	controller := &PointController{
		BaseController: base,
		event:          event,
		group:          group,
		app:            event.App,
		logger:         logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *PointController) registerRoutes() {
	group := controller.group.Group("/admin/point").Bind(
		RequireAdminRole(),
	)

	// 积分记录列表
	group.GET("/list", controller.List)
	// 批量创建积分记录
	group.POST("/batch/create", controller.BatchCreate)
	// 批量发放积分
	group.POST("/batch/distribute", controller.BatchDistribute)
	// 批量重试发放
	group.POST("/batch/retry", controller.BatchRetry)
	// 删除积分记录
	group.DELETE("/delete/{id}", controller.Delete)
	// 批量删除积分记录
	group.POST("/batch/delete", controller.BatchDelete)
}

func (controller *PointController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

// List 获取积分记录列表
func (controller *PointController) List(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("list")

	page, _ := strconv.Atoi(event.Request.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(event.Request.URL.Query().Get("pageSize"))
	status := event.Request.URL.Query().Get("status")
	group := event.Request.URL.Query().Get("group")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var points []*model.Point
	var total int

	// 构建查询条件
	query := event.App.RecordQuery(model.DbNamePoints)
	countQuery := event.App.RecordQuery(model.DbNamePoints)

	if status != "" {
		query = query.AndWhere(dbx.HashExp{model.PointsFieldStatus: status})
		countQuery = countQuery.AndWhere(dbx.HashExp{model.PointsFieldStatus: status})
	}
	if group != "" {
		query = query.AndWhere(dbx.HashExp{model.PointsFieldGroup: group})
		countQuery = countQuery.AndWhere(dbx.HashExp{model.PointsFieldGroup: group})
	}

	// 查询总数
	if err := countQuery.Select("count(*)").Row(&total); err != nil {
		logger.Error("查询积分记录总数失败", slog.Any("err", err))
		return event.InternalServerError("查询积分记录总数失败", err)
	}

	// 查询列表
	if err := query.OrderBy(fmt.Sprintf("%s DESC", model.PointsFieldCreated)).
		Limit(int64(pageSize)).
		Offset(int64((page - 1) * pageSize)).
		All(&points); err != nil {
		logger.Error("查询积分记录列表失败", slog.Any("err", err))
		return event.InternalServerError("查询积分记录列表失败", err)
	}

	// 批量查询关联用户信息
	userIds := make([]any, 0, len(points))
	for _, p := range points {
		if userId := p.UserId(); userId != "" {
			userIds = append(userIds, userId)
		}
	}

	usersMap := make(map[string]*model.User)
	if len(userIds) > 0 {
		var users []*model.User
		if err := event.App.RecordQuery(model.DbNameUsers).
			Where(dbx.In(model.CommonFieldId, userIds...)).
			All(&users); err != nil {
			logger.Warn("批量查询用户失败", slog.Any("err", err))
		} else {
			for _, user := range users {
				usersMap[user.Id] = user
			}
		}
	}

	// 构建响应
	type pointItem struct {
		Id      string `json:"id"`
		Group   string `json:"group"`
		UserId  string `json:"userId"`
		Point   int    `json:"point"`
		Status  string `json:"status"`
		Memo    string `json:"memo"`
		Created string `json:"created"`
		Updated string `json:"updated"`
		Expand  *struct {
			UserId *struct {
				Id       string `json:"id"`
				OId      string `json:"oId"`
				Name     string `json:"name"`
				Nickname string `json:"nickname"`
				Avatar   string `json:"avatar"`
			} `json:"userId"`
		} `json:"expand,omitempty"`
	}

	items := make([]*pointItem, 0, len(points))
	for _, p := range points {
		item := &pointItem{
			Id:      p.Id,
			Group:   p.Group(),
			UserId:  p.UserId(),
			Point:   p.Point(),
			Status:  string(p.Status()),
			Memo:    p.Memo(),
			Created: p.Created().String(),
			Updated: p.Updated().String(),
		}
		if user, exists := usersMap[p.UserId()]; exists {
			item.Expand = &struct {
				UserId *struct {
					Id       string `json:"id"`
					OId      string `json:"oId"`
					Name     string `json:"name"`
					Nickname string `json:"nickname"`
					Avatar   string `json:"avatar"`
				} `json:"userId"`
			}{
				UserId: &struct {
					Id       string `json:"id"`
					OId      string `json:"oId"`
					Name     string `json:"name"`
					Nickname string `json:"nickname"`
					Avatar   string `json:"avatar"`
				}{
					Id:       user.Id,
					OId:      user.OId(),
					Name:     user.Name(),
					Nickname: user.Nickname(),
					Avatar:   user.Avatar(),
				},
			}
		}
		items = append(items, item)
	}

	// 获取所有分组
	var groups []string
	rows, err := event.App.DB().Select("DISTINCT " + model.PointsFieldGroup).From(model.DbNamePoints).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var g string
			if err := rows.Scan(&g); err == nil && g != "" {
				groups = append(groups, g)
			}
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items":      items,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + pageSize - 1) / pageSize,
		"groups":     groups,
	})
}

// BatchCreate 批量创建积分记录
func (controller *PointController) BatchCreate(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("batch_create")

	var req struct {
		UserIds []string `json:"userIds"` // 用户oId列表
		Point   int      `json:"point"`   // 积分数量
		Memo    string   `json:"memo"`    // 备注
		Group   string   `json:"group"`   // 分组标识
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if len(req.UserIds) == 0 {
		return event.BadRequestError("用户列表不能为空", nil)
	}
	if req.Point == 0 {
		return event.BadRequestError("积分数量不能为0", nil)
	}

	pointCollection, err := event.App.FindCollectionByNameOrId(model.DbNamePoints)
	if err != nil {
		logger.Error("获取积分集合失败", slog.Any("err", err))
		return event.InternalServerError("获取积分集合失败", err)
	}

	var (
		createdCount int
		skippedCount int
		failedCount  int
	)

	for _, userOId := range req.UserIds {
		// 查找本地用户记录
		localUser := new(model.User)
		if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
			model.UsersFieldOId: userOId,
		}).One(localUser); err != nil {
			logger.Warn("本地用户不存在", slog.String("user_id", userOId))
			failedCount++
			continue
		}

		// 检查是否已存在相同分组和用户的待发放记录
		if req.Group != "" {
			existingPoint := new(model.Point)
			if err := event.App.RecordQuery(model.DbNamePoints).Where(dbx.And(
				dbx.HashExp{model.PointsFieldGroup: req.Group},
				dbx.HashExp{model.PointsFieldUserId: localUser.Id},
				dbx.HashExp{model.PointsFieldStatus: string(model.PointStatusPending)},
			)).One(existingPoint); err == nil {
				logger.Info("用户已存在相同分组的待发放记录，跳过", slog.String("user_id", userOId), slog.String("group", req.Group))
				skippedCount++
				continue
			}
		}

		// 创建积分记录
		pointRecord := model.NewPointFromCollection(pointCollection)
		pointRecord.SetGroup(req.Group)
		pointRecord.SetUserId(localUser.Id)
		pointRecord.SetPoint(req.Point)
		pointRecord.SetStatus(model.PointStatusPending)
		pointRecord.SetMemo(req.Memo)

		if err := event.App.Save(pointRecord); err != nil {
			logger.Error("创建积分记录失败", slog.Any("err", err), slog.String("user_id", userOId))
			failedCount++
			continue
		}
		createdCount++
	}

	logger.Info("批量创建积分记录完成",
		slog.Int("total", len(req.UserIds)),
		slog.Int("created", createdCount),
		slog.Int("skipped", skippedCount),
		slog.Int("failed", failedCount),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"total":   len(req.UserIds),
		"created": createdCount,
		"skipped": skippedCount,
		"failed":  failedCount,
	})
}

// BatchDistribute 批量发放积分
func (controller *PointController) BatchDistribute(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("batch_distribute")

	var req struct {
		Ids []string `json:"ids"` // 积分记录ID列表
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if len(req.Ids) == 0 {
		return event.BadRequestError("记录ID列表不能为空", nil)
	}

	isDev := controller.app.IsDev()

	var (
		successCount int
		failedCount  int
		skippedCount int
		results      []map[string]any
	)

	for _, id := range req.Ids {
		result := map[string]any{
			"id":      id,
			"success": false,
		}

		// 查询积分记录
		pointRecord := new(model.Point)
		if err := event.App.RecordQuery(model.DbNamePoints).Where(dbx.HashExp{
			model.CommonFieldId: id,
		}).One(pointRecord); err != nil {
			logger.Warn("积分记录不存在", slog.String("id", id))
			result["error"] = "记录不存在"
			failedCount++
			results = append(results, result)
			continue
		}

		// 检查状态
		if pointRecord.Status() == model.PointStatusSuccess {
			result["error"] = "已发放成功，无需重复发放"
			skippedCount++
			results = append(results, result)
			continue
		}

		if pointRecord.Status() == model.PointStatusDistributing {
			result["error"] = "正在发放中"
			skippedCount++
			results = append(results, result)
			continue
		}

		// 查询用户
		user := new(model.User)
		if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
			model.CommonFieldId: pointRecord.UserId(),
		}).One(user); err != nil {
			logger.Error("用户不存在", slog.String("user_id", pointRecord.UserId()))
			pointRecord.SetStatus(model.PointStatusFailed)
			pointRecord.SetMemo(fmt.Sprintf("用户不存在: %v", err))
			_ = event.App.Save(pointRecord)
			result["error"] = "用户不存在"
			failedCount++
			results = append(results, result)
			continue
		}

		// 更新状态为发放中
		pointRecord.SetStatus(model.PointStatusDistributing)
		if err := event.App.Save(pointRecord); err != nil {
			logger.Error("更新状态失败", slog.Any("err", err))
			result["error"] = "更新状态失败"
			failedCount++
			results = append(results, result)
			continue
		}

		// 构建发送给用户的memo（包含交易单号）
		userMemo := pointRecord.Memo()
		if userMemo != "" {
			userMemo = fmt.Sprintf("%s 交易单号：%s", userMemo, pointRecord.Id)
		} else {
			userMemo = fmt.Sprintf("积分发放 交易单号：%s", pointRecord.Id)
		}

		// 发放积分
		if isDev {
			logger.Warn("[DEV] 开发模式，跳过实际积分发放",
				slog.String("user_name", user.Name()),
				slog.Int("point", pointRecord.Point()),
				slog.String("memo", userMemo),
			)
		} else {
			resp, err := controller.fishPiSdk.PostUserEditPoints(user.Name(), pointRecord.Point(), userMemo)
			if err != nil {
				logger.Error("发放积分失败", slog.Any("err", err))
				pointRecord.SetStatus(model.PointStatusFailed)
				pointRecord.SetMemo(fmt.Sprintf("%s | 发放失败: %v", pointRecord.Memo(), err))
				_ = event.App.Save(pointRecord)
				result["error"] = "调用接口失败"
				failedCount++
				results = append(results, result)
				continue
			}
			if resp.Code != 0 {
				logger.Error("发放积分失败", slog.Any("resp", resp))
				pointRecord.SetStatus(model.PointStatusFailed)
				pointRecord.SetMemo(fmt.Sprintf("%s | 发放失败: %s", pointRecord.Memo(), resp.Msg))
				_ = event.App.Save(pointRecord)
				result["error"] = resp.Msg
				failedCount++
				results = append(results, result)
				continue
			}

			// 增加时间间隔，防止请求过于频繁
			time.Sleep(500 * time.Millisecond)
		}

		// 发放成功
		pointRecord.SetStatus(model.PointStatusSuccess)
		if err := event.App.Save(pointRecord); err != nil {
			logger.Error("更新成功状态失败", slog.Any("err", err))
		}

		result["success"] = true
		successCount++
		results = append(results, result)
	}

	logger.Info("批量发放积分完成",
		slog.Int("total", len(req.Ids)),
		slog.Int("success", successCount),
		slog.Int("failed", failedCount),
		slog.Int("skipped", skippedCount),
		slog.Bool("dev_mode", isDev),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"total":    len(req.Ids),
		"success":  successCount,
		"failed":   failedCount,
		"skipped":  skippedCount,
		"dev_mode": isDev,
		"results":  results,
	})
}

// BatchRetry 批量重试发放失败的记录
func (controller *PointController) BatchRetry(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("batch_retry")

	var req struct {
		Ids []string `json:"ids"` // 积分记录ID列表
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if len(req.Ids) == 0 {
		return event.BadRequestError("记录ID列表不能为空", nil)
	}

	// 筛选出失败的记录
	var failedIds []string
	for _, id := range req.Ids {
		pointRecord := new(model.Point)
		if err := event.App.RecordQuery(model.DbNamePoints).Where(dbx.HashExp{
			model.CommonFieldId: id,
		}).One(pointRecord); err != nil {
			continue
		}
		if pointRecord.Status() == model.PointStatusFailed {
			failedIds = append(failedIds, id)
		}
	}

	if len(failedIds) == 0 {
		return event.JSON(http.StatusOK, map[string]any{
			"message": "没有需要重试的失败记录",
			"total":   0,
		})
	}

	logger.Info("开始重试发放", slog.Int("count", len(failedIds)))

	// 复用 BatchDistribute 逻辑
	return controller.batchDistributeInternal(event, failedIds)
}

func (controller *PointController) batchDistributeInternal(event *core.RequestEvent, ids []string) error {
	logger := controller.makeActionLogger("batch_distribute_internal")

	isDev := controller.app.IsDev()

	var (
		successCount int
		failedCount  int
		results      []map[string]any
	)

	for _, id := range ids {
		result := map[string]any{
			"id":      id,
			"success": false,
		}

		// 查询积分记录
		pointRecord := new(model.Point)
		if err := event.App.RecordQuery(model.DbNamePoints).Where(dbx.HashExp{
			model.CommonFieldId: id,
		}).One(pointRecord); err != nil {
			result["error"] = "记录不存在"
			failedCount++
			results = append(results, result)
			continue
		}

		// 查询用户
		user := new(model.User)
		if err := event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{
			model.CommonFieldId: pointRecord.UserId(),
		}).One(user); err != nil {
			pointRecord.SetStatus(model.PointStatusFailed)
			pointRecord.SetMemo(fmt.Sprintf("用户不存在: %v", err))
			_ = event.App.Save(pointRecord)
			result["error"] = "用户不存在"
			failedCount++
			results = append(results, result)
			continue
		}

		// 更新状态为发放中
		pointRecord.SetStatus(model.PointStatusDistributing)
		if err := event.App.Save(pointRecord); err != nil {
			result["error"] = "更新状态失败"
			failedCount++
			results = append(results, result)
			continue
		}

		// 构建发送给用户的memo
		userMemo := pointRecord.Memo()
		// 移除之前的错误信息
		if idx := len(userMemo); idx > 0 {
			if sepIdx := findLastIndex(userMemo, " | "); sepIdx > 0 {
				userMemo = userMemo[:sepIdx]
			}
		}
		if userMemo != "" {
			userMemo = fmt.Sprintf("%s 交易单号：%s", userMemo, pointRecord.Id)
		} else {
			userMemo = fmt.Sprintf("积分发放 交易单号：%s", pointRecord.Id)
		}

		// 发放积分
		if isDev {
			logger.Warn("[DEV] 开发模式，跳过实际积分发放",
				slog.String("user_name", user.Name()),
				slog.Int("point", pointRecord.Point()),
			)
		} else {
			resp, err := controller.fishPiSdk.PostUserEditPoints(user.Name(), pointRecord.Point(), userMemo)
			if err != nil || resp.Code != 0 {
				errMsg := ""
				if err != nil {
					errMsg = err.Error()
				} else {
					errMsg = resp.Msg
				}
				pointRecord.SetStatus(model.PointStatusFailed)
				pointRecord.SetMemo(fmt.Sprintf("%s | 发放失败: %s", pointRecord.Memo(), errMsg))
				_ = event.App.Save(pointRecord)
				result["error"] = errMsg
				failedCount++
				results = append(results, result)
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}

		// 发放成功
		pointRecord.SetStatus(model.PointStatusSuccess)
		// 清理memo中的错误信息
		originalMemo := pointRecord.Memo()
		if sepIdx := findLastIndex(originalMemo, " | "); sepIdx > 0 {
			pointRecord.SetMemo(originalMemo[:sepIdx])
		}
		_ = event.App.Save(pointRecord)

		result["success"] = true
		successCount++
		results = append(results, result)
	}

	logger.Info("批量发放完成",
		slog.Int("total", len(ids)),
		slog.Int("success", successCount),
		slog.Int("failed", failedCount),
		slog.Bool("dev_mode", isDev),
	)

	return event.JSON(http.StatusOK, map[string]any{
		"total":    len(ids),
		"success":  successCount,
		"failed":   failedCount,
		"dev_mode": isDev,
		"results":  results,
	})
}

// findLastIndex 查找最后一个子串的位置
func findLastIndex(s, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Delete 删除积分记录
func (controller *PointController) Delete(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("delete")

	id := event.Request.PathValue("id")
	if id == "" {
		return event.BadRequestError("缺少记录ID", nil)
	}

	pointRecord := new(model.Point)
	if err := event.App.RecordQuery(model.DbNamePoints).Where(dbx.HashExp{
		model.CommonFieldId: id,
	}).One(pointRecord); err != nil {
		return event.NotFoundError("记录不存在", err)
	}

	// 检查状态，已发放成功的不允许删除
	if pointRecord.Status() == model.PointStatusSuccess {
		return event.BadRequestError("已发放成功的记录不能删除", nil)
	}

	if err := event.App.Delete(pointRecord); err != nil {
		logger.Error("删除记录失败", slog.Any("err", err))
		return event.InternalServerError("删除记录失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"deleted": true,
	})
}

// BatchDelete 批量删除积分记录
func (controller *PointController) BatchDelete(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("batch_delete")

	var req struct {
		Ids []string `json:"ids"`
	}

	if err := event.BindBody(&req); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	if len(req.Ids) == 0 {
		return event.BadRequestError("记录ID列表不能为空", nil)
	}

	var (
		deletedCount int
		skippedCount int
	)

	for _, id := range req.Ids {
		pointRecord := new(model.Point)
		if err := event.App.RecordQuery(model.DbNamePoints).Where(dbx.HashExp{
			model.CommonFieldId: id,
		}).One(pointRecord); err != nil {
			continue
		}

		// 已发放成功的不允许删除
		if pointRecord.Status() == model.PointStatusSuccess {
			skippedCount++
			continue
		}

		if err := event.App.Delete(pointRecord); err != nil {
			logger.Warn("删除记录失败", slog.Any("err", err), slog.String("id", id))
			continue
		}
		deletedCount++
	}

	return event.JSON(http.StatusOK, map[string]any{
		"total":   len(req.Ids),
		"deleted": deletedCount,
		"skipped": skippedCount,
	})
}
