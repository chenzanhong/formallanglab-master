// test/fixtures/automaton_fixtures.go
package fixtures

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// 创建自动机识别请求
func NewAutomatonStringRecognizeRequest(automatonJSON, input string) *http.Request {
	jsonData := `{"automaton":` + automatonJSON + `,"str":"` + input + `"}`
	return httptest.NewRequest("POST", "/api/automaton/recognize", strings.NewReader(jsonData))
}

// 创建带认证的请求
func NewAuthenticatedRequest(method, url, body, token string) *http.Request {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return req
}

// 有效的自动机测试数据
var (
	ValidAutomatonID    = uint(1)
	ValidAutomatonInput = "abab"
	InvalidAutomatonID  = uint(999)
	RejectionInput      = "invalid_input"
)
