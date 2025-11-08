package binding

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const ( // 用户当前停留的页面
	PageRegex     = "regex"
	PageGrammar   = "grammar"
	PageAutomaton = "automaton"
	PageLearn     = "learn"
	PageHome      = "home"
)

var validPages = map[string]bool{
	PageRegex:     true,
	PageGrammar:   true,
	PageAutomaton: true,
	PageLearn:     true,
	PageHome:      true,
}

type PageType string

func (p PageType) IsVaild() bool {
	return validPages[string(p)]
}

// 自定义校验函数
func PageValid(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}
	return validPages[fl.Field().String()]
}

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pageValid", PageValid)
	}
}
