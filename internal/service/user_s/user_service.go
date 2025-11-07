// internal/service/user_s/user_service.go
package user_s

import (
	"backend/internal/domain/model"
	"backend/internal/repository"
	"context"
)

type UserService interface {
	Register(ctx context.Context, name, email, password, token string) (*model.User, error)
	Login(ctx context.Context, name, password string) (string, *model.User, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type UserServiceImpl struct {
	userRepo  repository.UserRepository
	emailRepo repository.EmailRepository
}

func NewUserService(userRepo repository.UserRepository, emailRepo repository.EmailRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo, emailRepo: emailRepo}
}
