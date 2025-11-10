// test/e2e/auth_flow_test.go
package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthFlow 测试完整的认证流程
// 使用模拟服务器模拟完整的注册->登录->访问受保护资源流程
func TestAuthFlow(t *testing.T) {
	// 创建模拟服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/user/register":
			// 模拟注册接口
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// 模拟注册成功响应
			resp := map[string]interface{}{
				"result": true,
				"msg":    "注册成功",
				"id":     uint(1),
				"name":   "testuser",
			}
			json.NewEncoder(w).Encode(resp)
			return

		case "/api/user/login":
			// 模拟登录接口
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// 模拟登录成功响应
			resp := map[string]interface{}{
				"result": true,
				"msg":    "登录成功",
				"token":  "mock_jwt_token",
				"name":   "testuser",
				"id":     uint(1),
			}
			json.NewEncoder(w).Encode(resp)
			return

		case "/api/automaton/list":
			// 模拟受保护资源接口
			if r.Method != "GET" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// 验证认证token
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer mock_jwt_token" {
				w.WriteHeader(http.StatusUnauthorized)
				resp := map[string]interface{}{
					"result": false,
					"msg":    "未授权访问",
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			// 模拟受保护资源响应
			resp := map[string]interface{}{
				"result": true,
				"data": []map[string]interface{}{
					{
						"id":   uint(1),
						"name": "测试自动机",
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return

		default:
			w.WriteHeader(http.StatusNotFound)
			resp := map[string]interface{}{
				"result": false,
				"msg":    "接口不存在",
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer mockServer.Close()

	// 创建HTTP客户端
	client := &http.Client{}

	// 步骤1: 用户注册
	registerReq, err := http.NewRequest("POST", mockServer.URL+"/api/user/register", strings.NewReader(`{
		"name": "testuser",
		"email": "test@example.com",
		"password": "Password123",
		"token": "valid_token"
	}`))
	assert.NoError(t, err)
	registerReq.Header.Set("Content-Type", "application/json")

	registerResp, err := client.Do(registerReq)
	assert.NoError(t, err)
	defer registerResp.Body.Close()

	assert.Equal(t, http.StatusOK, registerResp.StatusCode)

	// 解析注册响应
	var registerData map[string]interface{}
	err = json.NewDecoder(registerResp.Body).Decode(&registerData)
	assert.NoError(t, err)
	assert.True(t, registerData["result"].(bool))
	assert.Equal(t, "注册成功", registerData["msg"].(string))

	// 步骤2: 用户登录
	loginReq, err := http.NewRequest("POST", mockServer.URL+"/api/user/login", strings.NewReader(`{
		"name": "testuser",
		"password": "Password123"
	}`))
	assert.NoError(t, err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := client.Do(loginReq)
	assert.NoError(t, err)
	defer loginResp.Body.Close()

	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// 解析登录响应
	var loginData map[string]interface{}
	err = json.NewDecoder(loginResp.Body).Decode(&loginData)
	assert.NoError(t, err)
	assert.True(t, loginData["result"].(bool))
	token, ok := loginData["token"].(string)
	assert.True(t, ok, "登录响应中应包含token")
	assert.NotEmpty(t, token)

	// 步骤3: 访问受保护资源
	protectedReq, err := http.NewRequest("GET", mockServer.URL+"/api/automaton/list", nil)
	assert.NoError(t, err)
	protectedReq.Header.Set("Authorization", "Bearer "+token)
	protectedReq.Header.Set("Content-Type", "application/json")

	protectedResp, err := client.Do(protectedReq)
	assert.NoError(t, err)
	defer protectedResp.Body.Close()

	assert.Equal(t, http.StatusOK, protectedResp.StatusCode)

	// 解析受保护资源响应
	var protectedData map[string]interface{}
	err = json.NewDecoder(protectedResp.Body).Decode(&protectedData)
	assert.NoError(t, err)
	assert.True(t, protectedData["result"].(bool))
	assert.Contains(t, protectedData, "data")

	// 测试无效token访问受保护资源
	invalidTokenReq, _ := http.NewRequest("GET", mockServer.URL+"/api/automaton/list", nil)
	invalidTokenReq.Header.Set("Authorization", "Bearer invalid_token")
	invalidTokenReq.Header.Set("Content-Type", "application/json")

	invalidTokenResp, _ := client.Do(invalidTokenReq)
	defer invalidTokenResp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, invalidTokenResp.StatusCode)
}

// MockServerTest 测试使用模拟服务器的集成流程
func TestMockServerTest(t *testing.T) {
	// 创建模拟服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/api/user/login" {
			// 模拟登录成功
			resp := map[string]interface{}{
				"result": true,
				"msg":    "登录成功",
				"token":  "mock_jwt_token",
				"name":   "testuser",
				"id":     1,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		if r.URL.Path == "/api/automaton/recognize" {
			// 检查认证token
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer mock_jwt_token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// 模拟自动机识别
			resp := map[string]interface{}{
				"result": true,
				"data": map[string]interface{}{
					"accepted": true,
					"message":  "字符串被接受",
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// 测试登录
	loginReq, _ := http.NewRequest("POST", mockServer.URL+"/api/user/login", nil)
	loginReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	loginResp, err := client.Do(loginReq)
	assert.NoError(t, err)
	defer loginResp.Body.Close()

	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// 解析登录响应
	var loginData map[string]interface{}
	err = json.NewDecoder(loginResp.Body).Decode(&loginData)
	assert.NoError(t, err)
	assert.True(t, loginData["result"].(bool))
	token := loginData["token"].(string)

	// 测试访问受保护资源
	automatonReq, _ := http.NewRequest("POST", mockServer.URL+"/api/automaton/recognize", nil)
	automatonReq.Header.Set("Authorization", "Bearer "+token)
	automatonReq.Header.Set("Content-Type", "application/json")

	automatonResp, err := client.Do(automatonReq)
	assert.NoError(t, err)
	defer automatonResp.Body.Close()

	assert.Equal(t, http.StatusOK, automatonResp.StatusCode)
}
