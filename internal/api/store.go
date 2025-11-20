package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/storage"
	storeSvc "backend/internal/service/store_s"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StoreHandler 存储模块处理器
type StoreHandler struct {
	storeService storeSvc.StoreService
}

// NewStoreHandler 创建存储处理器实例
func NewStoreHandler(storeService storeSvc.StoreService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

// CreateAutomaton 创建自动机存储记录
func (h *StoreHandler) CreateAutomaton(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.CreateAutomatonResponse{Msg: "用户未认证", Result: false})
		return
	}

	// 绑定请求参数
	var req dto.CreateAutomatonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateAutomatonResponse{Msg: "参数错误: " + err.Error(), Result: false})
		return
	}

	// 创建自动机记录对象
	record := &storage.AutomatonRecord{
		Name:      req.Automaton.Name,
		Username:  username.(string),
		Automaton: req.Automaton.Automaton,
	}

	// 调用服务创建自动机
	if err := h.storeService.CreateAutomaton(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateAutomatonResponse{Msg: "创建自动机失败: " + err.Error(), Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.CreateAutomatonResponse{Msg: "自动机存储成功", Result: true})
}

// FindAutomata 根据用户名分页查询自动机列表
func (h *StoreHandler) FindAutomata(c *gin.Context) {
	username, _ := c.Get("username")
	lastID := getUintQueryParam(c, "last_id", 0)
	limit := getUintQueryParam(c, "limit", 100)
	if limit <= 0 {
		limit = 100
	}

	resp, err := h.storeService.FindAutomata(c, username.(string), lastID, int(limit))
	if err != nil {
		c.JSON(500, dto.FindAutomatonResponse{Msg: "server error", Result: false})
		return
	}
	// 设置成功响应的Msg和Result字段
	resp.Msg = "查询成功"
	resp.Result = true
	c.JSON(200, resp)
}

// CreateGrammar 创建文法存储记录
func (h *StoreHandler) CreateGrammar(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 绑定请求参数
	var req dto.CreateGrammarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateGrammarResponse{Msg: "参数错误: " + err.Error(), Result: false})
		return
	}

	// 创建文法记录对象
	record := &storage.GrammarRecord{
		Name:     req.Grammar.Name,
		Username: username.(string),
		Grammar:  req.Grammar.Grammar,
	}

	// 调用服务创建文法
	if err := h.storeService.CreateGrammar(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateGrammarResponse{Msg: "创建文法失败: " + err.Error(), Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.CreateGrammarResponse{Msg: "文法创建成功", Result: true})
}

// FindGrammars 根据用户名分页查询文法列表
func (h *StoreHandler) FindGrammars(c *gin.Context) {
	username, _ := c.Get("username")
	lastID := getUintQueryParam(c, "last_id", 0)
	limit := getUintQueryParam(c, "limit", 100)
	if limit <= 0 {
		limit = 100
	}

	resp, err := h.storeService.FindGrammars(c, username.(string), lastID, int(limit))
	if err != nil {
		c.JSON(500, dto.FAToGrammarResponse{Msg: "server error", Result: false})
		return
	}
	// 设置成功响应的Msg和Result字段
	resp.Msg = "查询成功"
	resp.Result = true
	c.JSON(200, resp)
}

// CreateRegex 创建正则表达式存储记录
func (h *StoreHandler) CreateRegex(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 绑定请求参数
	var req dto.CreateRegexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateRegexResponse{Msg: "参数错误: " + err.Error(), Result: false})
		return
	}

	// 创建正则表达式记录对象
	record := &storage.RegexRecord{
		Name:     req.Pattern.Name,
		Username: username.(string),
		Pattern:  req.Pattern.Pattern,
	}

	// 调用服务创建正则表达式
	if err := h.storeService.CreateRegex(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateRegexResponse{Msg: "创建正则表达式失败: " + err.Error(), Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.CreateRegexResponse{Msg: "正则表达式创建成功", Result: true})
}

// FindRegexByUsername 根据用户名查询正则表达式列表
func (h *StoreHandler) FindRegexes(c *gin.Context) {
	username, _ := c.Get("username")
	lastID := getUintQueryParam(c, "last_id", 0)
	limit := getUintQueryParam(c, "limit", 100)
	if limit <= 0 {
		limit = 100
	}

	resp, err := h.storeService.FindRegexes(c, username.(string), lastID, int(limit))
	if err != nil {
		c.JSON(500, dto.FAToRegexResponse{Msg: "server error", Result: false})
		return
	}
	// 设置成功响应的Msg和Result字段
	resp.Msg = "查询成功"
	resp.Result = true
	c.JSON(200, resp)
}

// DeleteAutomaton 删除自动机记录
func (h *StoreHandler) DeleteAutomaton(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.DeleteAutomatonResponse{Msg: "用户未认证", Result: false})
		return
	}

	// 获取路径参数中的id
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.DeleteAutomatonResponse{Msg: "缺少必要参数", Result: false})
		return
	}

	// 调用服务删除自动机
	if err := h.storeService.DeleteAutomaton(c.Request.Context(), id, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.DeleteAutomatonResponse{Msg: "删除自动机失败", Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteAutomatonResponse{Msg: "删除成功", Result: true})
}

// DeleteGrammar 删除文法记录
func (h *StoreHandler) DeleteGrammar(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.DeleteGrammarResponse{Msg: "用户未认证", Result: false})
		return
	}

	// 获取路径参数中的id
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.DeleteGrammarResponse{Msg: "缺少必要参数", Result: false})
		return
	}

	// 调用服务删除文法
	if err := h.storeService.DeleteGrammar(c.Request.Context(), id, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.DeleteGrammarResponse{Msg: "删除文法失败", Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteGrammarResponse{Msg: "删除成功", Result: true})
}

// DeleteRegex 删除正则表达式记录
func (h *StoreHandler) DeleteRegex(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.DeleteRegexResponse{Msg: "用户未认证", Result: false})
		return
	}

	// 获取路径参数中的id
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.DeleteRegexResponse{Msg: "缺少必要参数", Result: false})
		return
	}

	// 调用服务删除正则表达式
	if err := h.storeService.DeleteRegex(c.Request.Context(), id, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.DeleteRegexResponse{Msg: "删除正则表达式失败", Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteRegexResponse{Msg: "删除成功", Result: true})
}

func getUintQueryParam(c *gin.Context, queryParam string, defaultValue uint) uint {
	v, exist := c.Get(queryParam)
	if !exist {
		return defaultValue
	}
	return v.(uint)
}
