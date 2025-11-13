package dto

// LearnMaterialResponse 学习资源响应DTO
type LearnMaterialResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
}

// LearnMaterialDetailResponse 学习资源详情响应DTO（包含下载链接）
type LearnMaterialDetailResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
	DownloadURL string `json:"download_url"`
}

// PresignedURLRequest 生成上传预签名URL请求DTO
type PresignedURLRequest struct {
	Filename string `json:"filename" binding:"required"`
	Category string `json:"category" binding:"required,oneof=grammar automaton regex general"`
}

// PresignedURLResponse 生成上传预签名URL响应DTO
type PresignedURLResponse struct {
	UploadURL  string `json:"upload_url"`
	MaterialID int64  `json:"material_id"`
}

// UpdateMaterialRequest 更新学习资源请求DTO
type UpdateMaterialRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
}
