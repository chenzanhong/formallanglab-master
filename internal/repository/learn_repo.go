package repository

import (
	"context"
	"errors"

	"backend/internal/domain/model"

	"gorm.io/gorm"
)

// LearnRepository 学习资源仓库接口
type LearnRepository interface {
	// ListMaterials 获取学习资源列表
	ListMaterials(ctx context.Context, category string) ([]*model.LearnMaterial, error)
	// GetMaterialByID 根据ID获取学习资源
	GetMaterialByID(ctx context.Context, id int64) (*model.LearnMaterial, error)
	// CreateMaterial 创建学习资源
	CreateMaterial(ctx context.Context, material *model.LearnMaterial) (int64, error)
	// DeleteMaterial 删除学习资源
	DeleteMaterial(ctx context.Context, id int64) error
}

// LearnRepositoryImpl 学习资源仓库实现
type LearnRepositoryImpl struct {
	DB *gorm.DB
}

// NewLearnRepository 创建学习资源仓库
func NewLearnRepository(db *gorm.DB) LearnRepository {
	return &LearnRepositoryImpl{DB: db}
}

// ListMaterials 获取学习资源列表
func (r *LearnRepositoryImpl) ListMaterials(ctx context.Context, category string) ([]*model.LearnMaterial, error) {
	var materials []*model.LearnMaterial
	query := r.DB.WithContext(ctx)

	// 如果指定了分类，则按分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 按创建时间倒序排列
	if err := query.Order("created_at DESC").Find(&materials).Error; err != nil {
		return nil, err
	}

	return materials, nil
}

// GetMaterialByID 根据ID获取学习资源
func (r *LearnRepositoryImpl) GetMaterialByID(ctx context.Context, id int64) (*model.LearnMaterial, error) {
	var material model.LearnMaterial
	if err := r.DB.WithContext(ctx).First(&material, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回nil表示资源不存在
		}
		return nil, err
	}
	return &material, nil
}

// CreateMaterial 创建学习资源
func (r *LearnRepositoryImpl) CreateMaterial(ctx context.Context, material *model.LearnMaterial) (int64, error) {
	if err := r.DB.WithContext(ctx).Create(material).Error; err != nil {
		return 0, err
	}
	return material.ID, nil
}

// DeleteMaterial 删除学习资源
func (r *LearnRepositoryImpl) DeleteMaterial(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.LearnMaterial{}, id).Error
}
