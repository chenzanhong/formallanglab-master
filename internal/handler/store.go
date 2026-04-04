package handler

import (
	"net/http"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/domain/storage"
	storeSvc "github.com/chenzanhong/formallanglab-master/internal/service/store_s"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService storeSvc.StoreService
}

func NewStoreHandler(storeService storeSvc.StoreService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

func (h *StoreHandler) CreateAutomaton(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.CreateAutomatonResponse{Msg: "用户未认证", Result: false})
		return
	}

	var req dto.CreateAutomatonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateAutomatonResponse{Msg: "参数错误：" + err.Error(), Result: false})
		return
	}

	record := &storage.AutomatonRecord{
		Name:      req.Automaton.Name,
		Username:  username.(string),
		Automaton: req.Automaton.Automaton,
	}

	if err := h.storeService.CreateAutomaton(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateAutomatonResponse{Msg: "创建自动机失败：" + err.Error(), Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.CreateAutomatonResponse{Msg: "自动机存储成功", Result: true})
}

func (h *StoreHandler) FindAutomata(c *gin.Context) {
	username, _ := c.Get("username")
	lastID := getUintQueryParam(c, "last_id", 0)
	limit := getUintQueryParam(c, "limit", 100)
	if limit <= 0 {
		limit = 100
	}

	resp, err := h.storeService.FindAutomata(c, username.(string), lastID, int(limit))
	if err != nil {
		c.JSON(500, dto.FindAutomatonResponse{Msg: "服务器错误", Result: false})
		return
	}

	resp.Msg = "查询成功"
	resp.Result = true
	c.JSON(200, resp)
}

func (h *StoreHandler) CreateGrammar(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.CreateAutomatonResponse{Msg: "用户未认证", Result: false})
		return
	}

	var req dto.CreateGrammarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateGrammarResponse{Msg: "参数错误解析错误，请检查输入", Result: false})
		return
	}

	record := &storage.GrammarRecord{
		Name:     req.Grammar.Name,
		Username: username.(string),
		Grammar:  req.Grammar.Grammar,
	}

	if err := h.storeService.CreateGrammar(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateGrammarResponse{Msg: "创建文法失败：" + err.Error(), Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.CreateGrammarResponse{Msg: "文法创建成功", Result: true})
}

func (h *StoreHandler) FindGrammars(c *gin.Context) {
	username, _ := c.Get("username")
	lastID := getUintQueryParam(c, "last_id", 0)
	limit := getUintQueryParam(c, "limit", 100)
	if limit <= 0 {
		limit = 100
	}

	resp, err := h.storeService.FindGrammars(c, username.(string), lastID, int(limit))
	if err != nil {
		c.JSON(500, dto.FAToGrammarResponse{Msg: "服务器错误", Result: false})
		return
	}

	resp.Msg = "查询成功"
	resp.Result = true
	c.JSON(200, resp)
}

func (h *StoreHandler) CreateRegex(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.CreateRegexResponse{Msg: "用户未认证"})
		return
	}

	var req dto.CreateRegexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateRegexResponse{Msg: "参数错误解析失败，请检查输入", Result: false})
		return
	}

	record := &storage.RegexRecord{
		Name:     req.Pattern.Name,
		Username: username.(string),
		Pattern:  req.Pattern.Pattern,
	}

	if err := h.storeService.CreateRegex(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateRegexResponse{Msg: "创建正则表达式失败：" + err.Error(), Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.CreateRegexResponse{Msg: "正则表达式创建成功", Result: true})
}

func (h *StoreHandler) FindRegexes(c *gin.Context) {
	username, _ := c.Get("username")
	lastID := getUintQueryParam(c, "last_id", 0)
	limit := getUintQueryParam(c, "limit", 100)
	if limit <= 0 {
		limit = 100
	}

	resp, err := h.storeService.FindRegexes(c, username.(string), lastID, int(limit))
	if err != nil {
		c.JSON(500, dto.FAToRegexResponse{Msg: "服务器错误", Result: false})
		return
	}

	resp.Msg = "查询成功"
	resp.Result = true
	c.JSON(200, resp)
}

func (h *StoreHandler) DeleteAutomaton(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.DeleteAutomatonResponse{Msg: "用户未认证", Result: false})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.DeleteAutomatonResponse{Msg: "缺少必要参数", Result: false})
		return
	}

	if err := h.storeService.DeleteAutomaton(c.Request.Context(), id, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.DeleteAutomatonResponse{Msg: "删除自动机失败", Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteAutomatonResponse{Msg: "删除成功", Result: true})
}

func (h *StoreHandler) DeleteGrammar(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.DeleteGrammarResponse{Msg: "用户未认证", Result: false})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.DeleteGrammarResponse{Msg: "缺少必要参数", Result: false})
		return
	}

	if err := h.storeService.DeleteGrammar(c.Request.Context(), id, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.DeleteGrammarResponse{Msg: "删除文法失败", Result: false})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteGrammarResponse{Msg: "删除成功", Result: true})
}

func (h *StoreHandler) DeleteRegex(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.DeleteRegexResponse{Msg: "用户未认证", Result: false})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.DeleteRegexResponse{Msg: "缺少必要参数", Result: false})
		return
	}

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
