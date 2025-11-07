package dto

import (
	"backend/internal/domain/model"
	"backend/pkg/binding"
)

type AIChatRequest struct {
	Question  string           `json:"question" binding:"required"`
	Page      binding.PageType `json:"page" binding:"pageValid"`
	Automaton *model.Automaton `json:"automaton,omitempty"` // ← 指针
	Grammar   *model.Grammar   `json:"grammar,omitempty"`
	Regex     *model.Regex     `json:"regex,omitempty"`
}
