package repository

import (
	"backend/internal/domain/model"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error

	ExistsByID(ctx context.Context, id int64) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)

	UpdateUser(ctx context.Context, user *model.User) error
	UpdatePasswordByUsername(ctx context.Context, username, password string) error
	UpdatePasswordByEmail(ctx context.Context, email, password string) error

	DeleteUser(ctx context.Context, user *model.User) error
	DeleteUserByID(ctx context.Context, id int64) error
	DeleteUserByEmail(ctx context.Context, email string) error
	DeleteUserByName(ctx context.Context, name string) error
}

type UserRepositoryImpl struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &UserRepositoryImpl{DB: db, Redis: redis}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.User{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *UserRepositoryImpl) UpdatePasswordByUsername(ctx context.Context, username, password string) error {
	return r.DB.WithContext(ctx).Model(&model.User{}).Where("name = ?", username).Update("password", password).Error
}

func (r *UserRepositoryImpl) UpdatePasswordByEmail(ctx context.Context, email, password string) error {
	return r.DB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Update("password", password).Error
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Delete(user).Error
}

func (r *UserRepositoryImpl) DeleteUserByID(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *UserRepositoryImpl) DeleteUserByEmail(ctx context.Context, email string) error {
	return r.DB.WithContext(ctx).Delete(&model.User{}, "email = ?", email).Error
}

func (r *UserRepositoryImpl) DeleteUserByName(ctx context.Context, name string) error {
	return r.DB.WithContext(ctx).Delete(&model.User{}, "name = ?", name).Error
}
