// test/fixtures/user_fixtures.go
package fixtures

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// 创建用户注册请求
func NewRegisterRequest(t *testing.T, name, email, password, token string) *http.Request {
	jsonData := `{"name":"` + name + `","email":"` + email + `","password":"` + password + `","token":"` + token + `"}`
	return httptest.NewRequest("POST", "/api/user/register", strings.NewReader(jsonData))
}

// 创建用户登录请求
func NewLoginRequest(t *testing.T, name, password string) *http.Request {
	jsonData := `{"name":"` + name + `","password":"` + password + `"}`
	return httptest.NewRequest("POST", "/api/user/login", strings.NewReader(jsonData))
}

// 创建Gin测试上下文
func NewTestContext(t *testing.T, req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request = req
	return c, r
}

// 有效的用户测试数据
var (
	ValidName     = "testuser"
	ValidUserEmail    = "test@example.com"
	ValidUserPassword = "Password123"
	ValidUserToken    = "valid_token"
)

// 无效的用户测试数据
var (
	InvalidName     = ""
	InvalidUserPassword = ""
	InvalidUserEmail    = "invalid-email"
)
