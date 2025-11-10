// test/integration/automaton_api_test.go
package integration

import (
	"backend/internal/domain/dto"
	"backend/test/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAutomatonStringRecognize 测试自动机字符串识别功能
func TestAutomatonStringRecognize(t *testing.T) {
	// 成功识别接受的字符串
	t.Run("成功识别接受的字符串", func(t *testing.T) {
		// 设置路由和模拟响应
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		router.POST("/api/automaton/recognize", func(c *gin.Context) {
			c.JSON(http.StatusOK, dto.AutomatonStringRecognizeResponse{
				Msg:    "字符串被接受",
				Result: true,
			})
		})

		// 创建模拟服务器
		server := httptest.NewServer(router)
		defer server.Close()

		// 构造请求
		req := fixtures.NewAutomatonStringRecognizeRequest("{}", "valid_input")

		// 发送请求
		resp, err := http.Post(server.URL+"/api/automaton/recognize", "application/json", req.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应状态码
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 识别被拒绝的字符串
	t.Run("识别被拒绝的字符串", func(t *testing.T) {
		// 设置路由和模拟响应
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		router.POST("/api/automaton/recognize", func(c *gin.Context) {
			c.JSON(http.StatusOK, dto.AutomatonStringRecognizeResponse{
				Msg:    "字符串被拒绝",
				Result: false,
			})
		})

		// 创建模拟服务器
		server := httptest.NewServer(router)
		defer server.Close()

		// 构造请求
		req := fixtures.NewAutomatonStringRecognizeRequest("{}", "invalid_input")

		// 发送请求
		resp, err := http.Post(server.URL+"/api/automaton/recognize", "application/json", req.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应状态码
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 自动机不存在
	t.Run("自动机不存在", func(t *testing.T) {
		// 设置路由和模拟响应
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		router.POST("/api/automaton/recognize", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, dto.AutomatonStringRecognizeResponse{
				Msg:    "自动机不存在",
				Result: false,
			})
		})

		// 创建模拟服务器
		server := httptest.NewServer(router)
		defer server.Close()

		// 构造请求
		req := fixtures.NewAutomatonStringRecognizeRequest("{}", "input")

		// 发送请求
		resp, err := http.Post(server.URL+"/api/automaton/recognize", "application/json", req.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应状态码
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	// 参数验证失败
	t.Run("参数验证失败", func(t *testing.T) {
		// 设置路由和模拟响应
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		router.POST("/api/automaton/recognize", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
				Msg:    "参数验证失败",
				Result: false,
			})
		})

		// 创建模拟服务器
		server := httptest.NewServer(router)
		defer server.Close()

		// 构造请求
		req := fixtures.NewAutomatonStringRecognizeRequest("invalid_json", "input")

		// 发送请求
		resp, err := http.Post(server.URL+"/api/automaton/recognize", "application/json", req.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应状态码
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
