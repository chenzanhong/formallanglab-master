// test/fixtures/mock_user_service.go
package fixtures

import (
	"backend/internal/domain/model"
	myErrors "backend/internal/errors"
	"context"
)



// MockUserService 模拟用户服务

type MockUserService struct {
	RegisterFunc func(ctx context.Context, name, email, password, token string) (*model.User, error)
	LoginFunc    func(ctx context.Context, name, password string) (string, *model.User, error)
}

// Register 模拟注册方法
func (m *MockUserService) Register(ctx context.Context, name, email, password, token string) (*model.User, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, name, email, password, token)
	}
	// 默认实现
	if name == "" || email == "" || password == "" || token == "" {
		return nil, myErrors.ErrUserCreationFailed
	}
	return &model.User{
		ID:    1,
		Name:  name,
		Email: email,
	},
		nil
}

// Login 模拟登录方法
func (m *MockUserService) Login(ctx context.Context, name, password string) (string, *model.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, name, password)
	}
	// 默认实现
	if name == ValidName && password == ValidUserPassword {
		return "mock_token_123", &model.User{
			ID:   1,
			Name: name,
		},
		nil
	}
	return "", nil, myErrors.ErrUserNotFound
}

// GetUserByID 模拟根据ID获取用户（实现接口需要）
func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return &model.User{
		ID:   id,
		Name: "user_" + string(rune(id+48)),
	},
		nil
}

// UpdateUser 模拟更新用户（实现接口需要）
func (m *MockUserService) UpdateUser(ctx context.Context, user *model.User) error {
	return nil
}

// DeleteUser 模拟删除用户（实现接口需要）
func (m *MockUserService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

// ResetPassword 模拟重置密码方法（实现接口需要）
func (m *MockUserService) ResetPassword(ctx context.Context, name, newPassword string) error {
	// 简单的默认实现
	if name == ValidName && newPassword != "" {
		return nil
	}
	return myErrors.ErrUserNotFound
}