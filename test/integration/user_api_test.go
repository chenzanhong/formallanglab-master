// test/integration/user_api_test.go
package integration

import (
	"backend/internal/api"
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	myErrors "backend/internal/errors"
	"backend/logs"
	"backend/test/fixtures"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化日志系统以避免nil pointer错误
	logs.InitLoggerDefault()
}

// TestUserRegister 测试用户注册功能
func TestUserRegister(t *testing.T) {
	t.Run("成功注册", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{
			RegisterFunc: func(ctx context.Context, name, email, password, token string) (*model.User, error) {
				return &model.User{
						ID:    1,
						Name:  name,
						Email: email,
					},
					nil
			},
		}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行
		req := fixtures.NewRegisterRequest(t, fixtures.ValidName, fixtures.ValidUserEmail, fixtures.ValidUserPassword, fixtures.ValidUserToken)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Register(c)

		// 验证
		assert.Equal(t, http.StatusOK, r.Code)

		var resp dto.RegisterResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Result)
		assert.Equal(t, "注册成功", resp.Msg)
		assert.Equal(t, uint(1), resp.ID)
		assert.Equal(t, fixtures.ValidName, resp.Name)
	})

	t.Run("参数验证失败", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行 - 缺少名称
		req := fixtures.NewRegisterRequest(t, fixtures.InvalidName, fixtures.ValidUserEmail, fixtures.ValidUserPassword, fixtures.ValidUserToken)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Register(c)

		// 验证
		assert.Equal(t, http.StatusBadRequest, r.Code)

		var resp dto.RegisterResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Result)
		assert.Equal(t, "请求数据格式错误", resp.Msg)
	})

	t.Run("密码加密失败", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{
			RegisterFunc: func(ctx context.Context, name, email, password, token string) (*model.User, error) {
				return nil, myErrors.ErrPasswordHashFailed
			},
		}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行
		req := fixtures.NewRegisterRequest(t, fixtures.ValidName, fixtures.ValidUserEmail, fixtures.ValidUserPassword, fixtures.ValidUserToken)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Register(c)

		// 验证
		assert.Equal(t, http.StatusInternalServerError, r.Code)

		var resp dto.RegisterResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Result)
		assert.Equal(t, "密码加密失败", resp.Msg)
	})

	t.Run("用户创建失败", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{
			RegisterFunc: func(ctx context.Context, name, email, password, token string) (*model.User, error) {
				return nil, myErrors.ErrUserCreationFailed
			},
		}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行
		req := fixtures.NewRegisterRequest(t, fixtures.ValidName, fixtures.ValidUserEmail, fixtures.ValidUserPassword, fixtures.ValidUserToken)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Register(c)

		// 验证
		assert.Equal(t, http.StatusInternalServerError, r.Code)

		var resp dto.RegisterResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Result)
		assert.Equal(t, "用户创建失败", resp.Msg)
	})
}

// TestUserLogin 测试用户登录功能
func TestUserLogin(t *testing.T) {
	t.Run("成功登录", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{
			LoginFunc: func(ctx context.Context, name, password string) (string, *model.User, error) {
				return "jwt_token_123", &model.User{
						ID:   1,
						Name: name,
					},
					nil
			},
		}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行
		req := fixtures.NewLoginRequest(t, fixtures.ValidName, fixtures.ValidUserPassword)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Login(c)

		// 验证
		assert.Equal(t, http.StatusOK, r.Code)

		var resp dto.LoginResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Result)
		assert.Equal(t, "登录成功", resp.Msg)
		assert.Equal(t, "jwt_token_123", resp.Token)
		assert.Equal(t, fixtures.ValidName, resp.Name)
	})

	t.Run("参数验证失败", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行 - 缺少密码
		req := fixtures.NewLoginRequest(t, fixtures.ValidName, fixtures.InvalidUserPassword)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Login(c)

		// 验证
		assert.Equal(t, http.StatusBadRequest, r.Code)

		var resp dto.LoginResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Result)
		assert.Equal(t, "登录数据解析失败", resp.Msg)
	})

	t.Run("用户不存在", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{
			LoginFunc: func(ctx context.Context, name, password string) (string, *model.User, error) {
				return "", nil, myErrors.ErrUserNotFound
			},
		}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行
		req := fixtures.NewLoginRequest(t, "nonexistent", fixtures.ValidUserPassword)
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Login(c)

		// 验证
		assert.Equal(t, http.StatusUnauthorized, r.Code)

		var resp dto.LoginResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Result)
		assert.Equal(t, "用户名或密码错误", resp.Msg)
	})

	t.Run("密码错误", func(t *testing.T) {
		// 设置
		mockUserService := &fixtures.MockUserService{
			LoginFunc: func(ctx context.Context, name, password string) (string, *model.User, error) {
				return "", nil, myErrors.ErrUserCreationFailed // 注意：实际代码中使用了这个错误表示密码错误
			},
		}
		userHandler := api.NewUserHandler(mockUserService)

		// 执行
		req := fixtures.NewLoginRequest(t, fixtures.ValidName, "WrongPassword")
		c, r := fixtures.NewTestContext(t, req)
		userHandler.Login(c)

		// 验证
		assert.Equal(t, http.StatusUnauthorized, r.Code)

		var resp dto.LoginResponse
		err := json.Unmarshal(r.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Result)
		assert.Equal(t, "用户名或密码错误", resp.Msg)
	})
}
