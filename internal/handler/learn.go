package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/metrics"
	"github.com/chenzanhong/formallanglab-master/internal/service/learn_s"
	"github.com/chenzanhong/zlog"
	"github.com/gin-gonic/gin"
)

// LearnHandler 学习资源处理器
type LearnHandler struct {
	learnService learn_s.LearnService
}

// NewLearnHandler 创建学习资源处理器
func NewLearnHandler(learnService learn_s.LearnService) *LearnHandler {
	return &LearnHandler{
		learnService: learnService,
	}
}

// LearnList 获取学习资料列表
// @Summary 获取学习资料列表
// @Tags learn
// @Param category query string false "分类参数（暂时不使用，返回所有文件）"
// @Success 200 {array} dto.LearnMaterialResponse
// @Router /learn [get]
func (h *LearnHandler) LearnList(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "list", time.Since(start).Seconds())
	}()

	// 获取查询参数
	category := c.DefaultQuery("category", "")

	// 调用服务
	materials, err := h.learnService.ListMaterials(c.Request.Context(), category)
	if err != nil {
		metrics.IncOperation("learn", "list", "error")
		zlog.Errorw("获取学习资源列表失败", "error", err, "category", category)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学习资源列表失败"})

		return
	}

	metrics.IncOperation("learn", "list", "success")
	zlog.Infow("获取学习资源列表成功", "count", len(materials), "category", category)
	c.JSON(http.StatusOK, materials)
}

// LearnGetByID 获取单个学习资源详情
// @Summary 获取单个学习资源详情
// @Tags learn
// @Param id path int true "资源ID"
// @Success 200 {object} dto.LearnMaterialDetailResponse
// @Router /learn/{id} [get]
func (h *LearnHandler) LearnGetByID(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "get_by_id", time.Since(start).Seconds())
	}()

	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		metrics.IncOperation("learn", "get_by_id", "error")
		zlog.Errorw("无效的资源ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资源ID"})

		return
	}

	// 调用服务
	material, err := h.learnService.GetMaterialByID(c.Request.Context(), id)
	if err != nil {
		metrics.IncOperation("learn", "get_by_id", "error")
		zlog.Errorw("获取学习资源详情失败", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学习资源详情失败"})

		return
	}

	if material == nil {
		metrics.IncOperation("learn", "get_by_id", "error")
		zlog.Warnw("学习资源不存在", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "学习资源不存在"})

		return
	}

	metrics.IncOperation("learn", "get_by_id", "success")
	zlog.Infow("获取学习资源详情成功", "id", id, "title", material.Title)
	c.JSON(http.StatusOK, material)
}

// LearnAddMaterial 添加学习资源（管理员直接在 OSS 上传后调用）
// @Summary 添加学习资源（管理员直接在 OSS 上传后调用）
// @Tags learn
// @Accept json
// @Produce json
// @Param request body dto.AddMaterialRequest true "添加学习资源请求"
// @Success 200 {object} map[string]string
// @Router /learn [post]
func (h *LearnHandler) LearnAddMaterial(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "add", time.Since(start).Seconds())
	}()

	// TODO: 验证用户是否为管理员
	// 暂时注释，开发阶段可以跳过权限检查
	// if !isAdmin(c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "无权限执行此操作"})
	// 	return
	// }

	// 绑定请求参数
	var req dto.AddMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("learn", "add", "error")
		zlog.Errorw("无效的请求参数", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// 调用服务
	if err := h.learnService.AddMaterial(c.Request.Context(), &req); err != nil {
		metrics.IncOperation("learn", "add", "error")
		zlog.Errorw("添加学习资源失败", "error", err, "title", req.Title)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加学习资源失败"})

		return
	}

	metrics.IncOperation("learn", "add", "success")
	zlog.Infow("添加学习资源成功", "title", req.Title, "file_key", req.FileKey)
	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// LearnDeleteMaterial 删除学习资源
// @Summary 删除学习资源
// @Tags learn
// @Param id path int true "资源ID"
// @Success 200 {object} map[string]string
// @Router /learn/{id} [delete]
func (h *LearnHandler) LearnDeleteMaterial(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "delete", time.Since(start).Seconds())
	}()

	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		metrics.IncOperation("learn", "delete", "error")
		zlog.Errorw("无效的资源ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资源ID"})

		return
	}

	// 调用服务
	if err := h.learnService.DeleteMaterial(c.Request.Context(), id); err != nil {
		metrics.IncOperation("learn", "delete", "error")
		zlog.Errorw("删除学习资源失败", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除学习资源失败"})

		return
	}

	metrics.IncOperation("learn", "delete", "success")
	zlog.Infow("删除学习资源成功", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// LearnSyncOSSFiles 同步OSS文件到数据库（自动识别新增文件）
// @Summary 同步OSS文件到数据库（管理员在OSS上传后调用）
// @Tags learn
// @Success 200 {object} dto.SyncOSSFilesResponse
// @Router /learn/sync [post]
func (h *LearnHandler) LearnSyncOSSFiles(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "sync", time.Since(start).Seconds())
	}()

	// 调用服务
	response, err := h.learnService.SyncOSSFiles(c.Request.Context())
	if err != nil {
		metrics.IncOperation("learn", "sync", "error")
		zlog.Errorw("同步OSS文件失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步OSS文件失败"})

		return
	}

	metrics.IncOperation("learn", "sync", "success")
	zlog.Infow("同步OSS文件成功", "total_files", response.TotalFiles, "new_files", response.NewFiles, "existing_files", response.ExistingFiles)
	c.JSON(http.StatusOK, response)
}
