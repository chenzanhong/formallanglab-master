package user_s

import (
	"backend/internal/domain/model"
	"backend/internal/errors"

	// "backend/internal/utils"
	"backend/internal/middleware"
	"backend/pkg/cryptoutil"
	"context"
)

func (s *UserServiceImpl) Login(ctx context.Context, name, password string) (string, *model.User, error) {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByName(ctx, name)
	if err != nil {
		return "", nil, errors.ErrInternal
	}
	if !exists {
		return "", nil, errors.ErrUserNotFound
	}
	// 检查密码是否正确
	user, err := s.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return "", nil, errors.ErrInternal
	}

	if !cryptoutil.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.ErrInvalidCredentials
	}

	// 生成JWT
	token, err := middleware.GenerateToken(name)
	if err != nil {
		return "", nil, errors.ErrTokenGenerationFailed
	}
	return token, user, nil
}
