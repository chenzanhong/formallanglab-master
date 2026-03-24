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
	// AddMaterial 添加学习资源（管理员直接在 OSS 上传后调用）
	AddMaterial(ctx context.Context, req *dto.AddMaterialRequest) error
	// SyncOSSFiles 同步OSS文件到数据库（自动识别新增文件）
	SyncOSSFiles(ctx context.Context) (*dto.SyncOSSFilesResponse, error)
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
		"grammar":   true,
		"automaton": true,
		"regex":     true,
		"general":   true,
		"":          true, // 空字符串表示获取所有分类
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

// AddMaterial 添加学习资源（管理员直接在 OSS 上传后调用）
func (s *LearnServiceImpl) AddMaterial(ctx context.Context, req *dto.AddMaterialRequest) error {
	// 验证文件是否在 OSS 中存在
	exists, err := s.ossClient.CheckObjectExists(req.FileKey)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("file does not exist in OSS")
	}

	// 创建学习资源记录
	material := &model.LearnMaterial{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		FileKey:     req.FileKey,
		FileName:    req.FileName,
		MimeType:    req.MimeType,
		SizeBytes:   req.SizeBytes,
	}

	// 保存到数据库
	_, err = s.learnRepo.CreateMaterial(ctx, material)
	if err != nil {
		return fmt.Errorf("failed to create material: %w", err)
	}

	return nil
}

// SyncOSSFiles 同步OSS文件到数据库（自动识别新增文件）
func (s *LearnServiceImpl) SyncOSSFiles(ctx context.Context) (*dto.SyncOSSFilesResponse, error) {
	// 定义要扫描的分类目录
	categories := []string{"grammar", "automaton", "regex", "general"}

	// 获取数据库中已存在的所有文件
	existingMaterials, err := s.learnRepo.ListMaterials(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list existing materials: %w", err)
	}

	// 创建已存在文件的映射（file_key -> material）
	existingFileKeys := make(map[string]bool)
	for _, material := range existingMaterials {
		existingFileKeys[material.FileKey] = true
	}

	response := &dto.SyncOSSFilesResponse{
		AddedFiles: []string{},
	}

	// 扫描每个分类目录
	for _, category := range categories {
		prefix := fmt.Sprintf("materials/%s/", category)

		// 列出OSS中的文件
		ossObjects, err := s.ossClient.ListObjects(prefix)
		if err != nil {
			return nil, fmt.Errorf("failed to list OSS objects for prefix %s: %w", prefix, err)
		}

		response.TotalFiles += int64(len(ossObjects))

		// 对比找出新增的文件
		for _, obj := range ossObjects {
			if existingFileKeys[obj.Key] {
				response.ExistingFiles++
				continue
			}

			// 新增文件，添加到数据库
			fileName := extractFileName(obj.Key)
			mimeType := inferMimeType(fileName)

			material := &model.LearnMaterial{
				Title:       fileName,
				Description: "",
				Category:    category,
				FileKey:     obj.Key,
				FileName:    fileName,
				MimeType:    mimeType,
				SizeBytes:   obj.Size,
			}

			_, err := s.learnRepo.CreateMaterial(ctx, material)
			if err != nil {
				return nil, fmt.Errorf("failed to create material for %s: %w", obj.Key, err)
			}

			response.NewFiles++
			response.AddedFiles = append(response.AddedFiles, obj.Key)
		}
	}

	return response, nil
}

// extractFileName 从 file_key 中提取文件名
func extractFileName(fileKey string) string {
	parts := strings.Split(fileKey, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return fileKey
}

// inferMimeType 根据文件扩展名推断 MIME 类型
func inferMimeType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".ppt", ".pptx":
		return "application/vnd.ms-powerpoint"
	case ".doc", ".docx":
		return "application/msword"
	case ".txt":
		return "text/plain"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
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
