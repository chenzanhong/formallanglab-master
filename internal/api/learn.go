package api

import (
	"backend/internal/domain/dto"
	"backend/internal/metrics"
	"backend/internal/service/learn_s"
	"net/http"
	"strconv"
	"time"

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
// @Param category query string false "分类: grammar/automaton/regex/general"
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

// LearnGeneratePresignedURL 生成上传预签名URL
// @Summary 生成上传预签名URL（仅管理员）
// @Tags learn
// @Accept json
// @Produce json
// @Param request body dto.PresignedURLRequest true "预签名URL请求"
// @Success 200 {object} dto.PresignedURLResponse
// @Router /learn/presigned-url [post]
func (h *LearnHandler) LearnGeneratePresignedURL(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "generate_presigned_url", time.Since(start).Seconds())
	}()

	// TODO: 验证用户是否为管理员
	// 暂时注释，开发阶段可以跳过权限检查
	// if !isAdmin(c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "无权限执行此操作"})
	// 	return
	// }

	// 绑定请求参数
	var req dto.PresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("learn", "generate_presigned_url", "error")
		zlog.Errorw("无效的请求参数", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	response, err := h.learnService.GeneratePresignedURL(c.Request.Context(), &req)
	if err != nil {
		metrics.IncOperation("learn", "generate_presigned_url", "error")
		zlog.Errorw("生成预签名URL失败", "error", err, "filename", req.Filename, "category", req.Category)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成预签名URL失败"})
		return
	}

	metrics.IncOperation("learn", "generate_presigned_url", "success")
	zlog.Infow("生成预签名URL成功", "material_id", response.MaterialID, "filename", req.Filename)
	c.JSON(http.StatusOK, response)
}

// LearnUpdateMaterial 更新学习资源信息
// @Summary 更新学习资源信息
// @Tags learn
// @Accept json
// @Produce json
// @Param id path int true "资源ID"
// @Param request body dto.UpdateMaterialRequest true "更新请求"
// @Success 200 {object} map[string]string
// @Router /learn/{id} [put]
func (h *LearnHandler) LearnUpdateMaterial(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "update", time.Since(start).Seconds())
	}()

	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		metrics.IncOperation("learn", "update", "error")
		zlog.Errorw("无效的资源ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资源ID"})
		return
	}

	// 绑定请求参数
	var req dto.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("learn", "update", "error")
		zlog.Errorw("无效的请求参数", "error", err, "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	if err := h.learnService.UpdateMaterial(c.Request.Context(), id, &req); err != nil {
		metrics.IncOperation("learn", "update", "error")
		zlog.Errorw("更新学习资源失败", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新学习资源失败"})
		return
	}

	metrics.IncOperation("learn", "update", "success")
	zlog.Infow("更新学习资源成功", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
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
