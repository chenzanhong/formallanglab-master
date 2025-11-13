// 用于获取学习资源
package learn_s

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/repository"
	"backend/pkg/oss"
)

// LearnService 学习资源服务接口
type LearnService interface {
	// ListMaterials 获取学习资源列表
	ListMaterials(ctx context.Context, category string) ([]*dto.LearnMaterialResponse, error)
	// GetMaterialByID 根据ID获取学习资源详情
	GetMaterialByID(ctx context.Context, id int64) (*dto.LearnMaterialDetailResponse, error)
	// GeneratePresignedURL 生成上传预签名URL
	GeneratePresignedURL(ctx context.Context, req *dto.PresignedURLRequest) (*dto.PresignedURLResponse, error)
	// UpdateMaterial 更新学习资源信息
	UpdateMaterial(ctx context.Context, id int64, req *dto.UpdateMaterialRequest) error
	// DeleteMaterial 删除学习资源
	DeleteMaterial(ctx context.Context, id int64) error
}

// LearnServiceImpl 学习资源服务实现
type LearnServiceImpl struct {
	learnRepo repository.LearnRepository
	ossClient oss.Client
}

// NewLearnService 创建学习资源服务
func NewLearnService(learnRepo repository.LearnRepository, ossClient oss.Client) LearnService {
	return &LearnServiceImpl{
		learnRepo: learnRepo,
		ossClient: ossClient,
	}
}

// ListMaterials 获取学习资源列表
func (s *LearnServiceImpl) ListMaterials(ctx context.Context, category string) ([]*dto.LearnMaterialResponse, error) {
	// 验证分类参数
	validCategories := map[string]bool{
		"grammar":  true,
		"automaton": true,
		"regex":    true,
		"general":  true,
		"":         true, // 空字符串表示获取所有分类
	}
	if !validCategories[category] {
		category = "general" // 默认为general分类
	}

	// 从仓库获取数据
	materials, err := s.learnRepo.ListMaterials(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials: %w", err)
	}

	// 转换为DTO
	responses := make([]*dto.LearnMaterialResponse, len(materials))
	for i, material := range materials {
		responses[i] = &dto.LearnMaterialResponse{
			ID:          material.ID,
			Title:       material.Title,
			Description: material.Description,
			Category:    material.Category,
			FileName:    material.FileName,
			MimeType:    material.MimeType,
			SizeBytes:   material.SizeBytes,
		}
	}

	return responses, nil
}

// GetMaterialByID 根据ID获取学习资源详情
func (s *LearnServiceImpl) GetMaterialByID(ctx context.Context, id int64) (*dto.LearnMaterialDetailResponse, error) {
	// 从仓库获取数据
	material, err := s.learnRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get material: %w", err)
	}
	if material == nil {
		return nil, nil // 返回nil表示资源不存在
	}

	// 生成5分钟有效的下载链接
	downloadURL, err := s.ossClient.GeneratePresignedURL(material.FileKey, time.Minute*5)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

	// 构造响应
	response := &dto.LearnMaterialDetailResponse{
		ID:          material.ID,
		Title:       material.Title,
		Description: material.Description,
		Category:    material.Category,
		FileName:    material.FileName,
		MimeType:    material.MimeType,
		SizeBytes:   material.SizeBytes,
		DownloadURL: downloadURL,
	}

	return response, nil
}

// GeneratePresignedURL 生成上传预签名URL
func (s *LearnServiceImpl) GeneratePresignedURL(ctx context.Context, req *dto.PresignedURLRequest) (*dto.PresignedURLResponse, error) {
	// 生成唯一的file_key
	ext := filepath.Ext(req.Filename)
	fileNameWithoutExt := strings.TrimSuffix(req.Filename, ext)
	fileKey := fmt.Sprintf("materials/%s/%s_%d%s",
		req.Category,
		fileNameWithoutExt,
		time.Now().Unix(),
		ext)

	// 生成10分钟有效的上传预签名URL
	uploadURL, err := s.ossClient.GenerateUploadPresignedURL(fileKey, time.Minute*10)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	// 创建待处理的数据库记录
	material := &model.LearnMaterial{
		Title:     req.Filename,
		Category:  req.Category,
		FileKey:   fileKey,
		FileName:  req.Filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	materialID, err := s.learnRepo.CreatePendingMaterial(ctx, material)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending material: %w", err)
	}

	return &dto.PresignedURLResponse{
		UploadURL:  uploadURL,
		MaterialID: materialID,
	}, nil
}

// UpdateMaterial 更新学习资源信息
func (s *LearnServiceImpl) UpdateMaterial(ctx context.Context, id int64, req *dto.UpdateMaterialRequest) error {
	// 获取现有资源
	material, err := s.learnRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get material: %w", err)
	}
	if material == nil {
		return fmt.Errorf("material not found")
	}

	// 更新字段
	material.Title = req.Title
	material.Description = req.Description
	material.MimeType = req.MimeType
	material.SizeBytes = req.SizeBytes
	material.UpdatedAt = time.Now()

	// 保存更新
	if err := s.learnRepo.UpdateMaterial(ctx, material); err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}

	return nil
}

// DeleteMaterial 删除学习资源
func (s *LearnServiceImpl) DeleteMaterial(ctx context.Context, id int64) error {
	// 检查资源是否存在
	material, err := s.learnRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get material: %w", err)
	}
	if material == nil {
		return fmt.Errorf("material not found")
	}

	// 删除资源
	if err := s.learnRepo.DeleteMaterial(ctx, id); err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}

	return nil
}