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

// AddMaterialRequest 添加学习资源请求DTO
type AddMaterialRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	FileKey     string `json:"file_key" binding:"required"`
	FileName    string `json:"file_name" binding:"required"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
	Category    string `json:"category" binding:"required,oneof=grammar automaton regex general"`
}

// SyncOSSFilesResponse 同步OSS文件响应DTO
type SyncOSSFilesResponse struct {
	TotalFiles    int64    `json:"total_files"`    // OSS中的文件总数
	NewFiles      int64    `json:"new_files"`      // 新增的文件数
	ExistingFiles int64    `json:"existing_files"` // 已存在的文件数
	AddedFiles    []string `json:"added_files"`    // 新增的文件列表
}
