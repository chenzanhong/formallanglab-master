package dto

import (
	"time"

	"backend/internal/domain/model"
)

// ———————— Item DTOs（用于列表响应，含完整数据） ————————

type AutomatonItem struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	Automaton model.Automaton `json:"automaton"`
	CreatedAt time.Time       `json:"created_at"`
}

type GrammarItem struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Grammar   model.Grammar `json:"grammar"`
	CreatedAt time.Time     `json:"created_at"`
}

type RegexItem struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Pattern   model.Regex `json:"pattern"`
	CreatedAt time.Time   `json:"created_at"`
}

// ================== 创建 ===================
// CreateAutomatonRequest 创建自动机请求DTO
type CreateAutomatonRequest struct {
	Automaton AutomatonItem `json:"automaton" binding:"required"`
}

type CreateAutomatonResponse = BaseResponse

// CreateGrammarRequest 创建文法请求DTO
type CreateGrammarRequest struct {
	Grammar GrammarItem `json:"grammar" binding:"required"`
}

// CreateGrammarResponse 创建文法响应DTO
type CreateGrammarResponse = BaseResponse

// CreateRegexRequest 创建正则表达式请求DTO
type CreateRegexRequest struct {
	Pattern RegexItem `json:"pattern" binding:"required"`
}

// CreateRegexResponse 创建正则表达式响应DTO
type CreateRegexResponse = BaseResponse

// ================== 删除 ===================
// DeleteGrammarRequest 删除文法请求DTO
type DeleteGrammarRequest struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteGrammarResponse 删除文法响应DTO
type DeleteGrammarResponse = BaseResponse

// DeleteRegexRequest 删除正则表达式请求DTO
type DeleteRegexRequest struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteRegexResponse 删除正则表达式响应DTO
type DeleteRegexResponse = BaseResponse

// DeleteAutomatonRequest 删除自动机请求DTO
type DeleteAutomatonRequest struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteAutomatonResponse 删除自动机响应DTO
type DeleteAutomatonResponse = BaseResponse

// ================== 分页查找 ===================
type FindAutomatonResponse struct {
	Msg        string          `json:"msg"`
	Result     bool            `json:"result"`
	Automatons []AutomatonItem `json:"automatons"`
	HasMore    bool            `json:"has_more"`
	NextCursor uint            `json:"next_cursor,omitempty"`
}

type FindGrammarResponse struct {
	Msg        string        `json:"msg"`
	Result     bool          `json:"result"`
	Grammars   []GrammarItem `json:"grammars"`
	HasMore    bool          `json:"has_more"`
	NextCursor uint          `json:"next_cursor,omitempty"`
}

type FindRegexResponse struct {
	Msg        string      `json:"msg"`
	Result     bool        `json:"result"`
	Patterns   []RegexItem `json:"patterns"`
	HasMore    bool        `json:"has_more"`
	NextCursor uint        `json:"next_cursor,omitempty"`
}
