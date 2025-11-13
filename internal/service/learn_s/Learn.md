你已经为 **学习资源模块（Learn Module）** 设计了清晰的存储策略：**资源文件存 OSS/S3，数据库只存元数据**。接下来，我们需要：

1. **设计数据库模型（PostgreSQL 表结构）**
2. **实现对应的 API 路由与 Handler**
3. **集成预签名 URL 机制（上传/下载）**

---

## ✅ 一、数据库模型设计（`learn_materials` 表）

### 表名：`learn_materials`

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | `BIGSERIAL PRIMARY KEY` | 主键 |
| `title` | `VARCHAR(255) NOT NULL` | 资源标题，如“DFA 构造详解” |
| `description` | `TEXT` | 可选描述 |
| `category` | `VARCHAR(50) NOT NULL` | 分类：`grammar` / `automaton` / `regex` / `general` |
| `file_key` | `VARCHAR(512) NOT NULL UNIQUE` | OSS 中的 object key，如 `materials/automaton/dfa_intro_v2.pdf` |
| `file_name` | `VARCHAR(255) NOT NULL` | 原始文件名（用于下载时显示） |
| `mime_type` | `VARCHAR(100)` | 如 `application/pdf`, `video/mp4` |
| `size_bytes` | `BIGINT` | 文件大小（字节） |
| `created_at` | `TIMESTAMP WITH TIME ZONE DEFAULT NOW()` | 创建时间 |
| `updated_at` | `TIMESTAMP WITH TIME ZONE DEFAULT NOW()` | 更新时间 |

> 💡 **为什么不用 `url` 字段？**  
> 因为 URL = `https://<bucket>.oss-cn-beijing.aliyuncs.com/` + `file_key`，可动态拼接，避免硬编码域名。

---

## ✅ 二、API 路由设计（补充到 `setupAuthRoutes`）

在 `/gdesign/learn` 下增加：

```go
learn := r.Group("/learn")
{
	learn.GET("/", LearnList)           // 获取某分类下的学习资料列表
	learn.GET("/:id", LearnGetByID)     // 获取单个资源详情（含下载链接）
	learn.POST("/presigned-url", LearnGeneratePresignedURL) // 生成上传预签名 URL（仅管理员）
}
```

> ⚠️ 注意：
> - 普通用户只能 **读取/下载**
> - **上传权限应限制为管理员**（可通过角色判断，或暂时只开放内部接口）

---

## ✅ 三、Handler 实现示例

### 1. 列出学习资料（按分类）

```go
// LearnList godoc
// @Summary 获取学习资料列表
// @Tags learn
// @Param category query string false "分类: grammar/automaton/regex/general"
// @Success 200 {array} dto.LearnMaterialResponse
// @Router /learn [get]
func LearnList(c *gin.Context) {
	category := c.DefaultQuery("category", "general")

	materials, err := learnService.ListMaterials(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch materials"})
		return
	}

	c.JSON(http.StatusOK, materials)
}
```

### 2. 获取单个资源 + 生成**下载预签名 URL**

```go
func LearnGetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	material, err := learnService.GetMaterialByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
		return
	}

	// 生成 5 分钟有效的下载链接
	downloadURL, err := ossClient.GeneratePresignedURL(material.FileKey, time.Minute*5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          material.ID,
		"title":       material.Title,
		"description": material.Description,
		"file_name":   material.FileName,
		"mime_type":   material.MimeType,
		"download_url": downloadURL, // 前端直接跳转或 <a href>
	})
}
```

### 3. （可选）生成上传预签名 URL（仅限管理员）

```go
type PresignedURLRequest struct {
	Filename string `json:"filename" binding:"required"`
	Category string `json:"category" binding:"required,oneof=grammar automaton regex general"`
}

func LearnGeneratePresignedURL(c *gin.Context) {
	var req PresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 验证用户是否为管理员
	// if !isAdmin(c) { c.AbortWithStatus(http.StatusForbidden); return }

	// 生成唯一 file_key，例如：materials/automaton/dfa_intro_20251111123456.pdf
	ext := filepath.Ext(req.Filename)
	key := fmt.Sprintf("materials/%s/%s_%d%s",
		req.Category,
		strings.TrimSuffix(req.Filename, ext),
		time.Now().Unix(),
		ext)

	uploadURL, err := ossClient.GenerateUploadPresignedURL(key, time.Minute*10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate upload URL"})
		return
	}

	// 先创建数据库记录（状态为 pending），等前端上传成功后再更新 size/mime_type
	// 或者让前端上传后回调通知后端
	materialID, err := learnService.CreatePendingMaterial(c.Request.Context(), &model.LearnMaterial{
		Title:     req.Filename,
		Category:  req.Category,
		FileKey:   key,
		FileName:  req.Filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"material_id": materialID,
	})
}
```

---

## ✅ 四、OSS 客户端封装（建议放在 `pkg/oss`）

```go
// pkg/oss/client.go
package oss

import (
	"time"
)

type Client interface {
	GeneratePresignedURL(key string, expire time.Duration) (string, error)
	GenerateUploadPresignedURL(key string, expire time.Duration) (string, error)
}

// 使用阿里云 OSS SDK 实现
type AliyunOSSClient struct {
	bucket *oss.Bucket
}

func (c *AliyunOSSClient) GeneratePresignedURL(key string, expire time.Duration) (string, error) {
	return c.bucket.SignURL(key, oss.HTTPGet, int64(expire.Seconds()))
}

func (c *AliyunOSSClient) GenerateUploadPresignedURL(key string, expire time.Duration) (string, error) {
	return c.bucket.SignURL(key, oss.HTTPPut, int64(expire.Seconds()))
}
```

在 `main.go` 中初始化并注入到 `LearnService`。

---

## ✅ 五、前端使用流程

### 下载：
1. 用户点击“下载课件”
2. 前端调用 `GET /gdesign/learn/123`
3. 后端返回 `{ download_url: "https://xxx?Expires=...&OSSAccessKeyId=...&Signature=..." }`
4. 前端执行：`window.location.href = download_url`

### 上传（管理员）：
1. 管理员选择文件 → 前端调用 `/presigned-url`
2. 后端返回 `upload_url` 和 `material_id`
3. 前端执行：`fetch(upload_url, { method: 'PUT', body: file })`
4. （可选）上传成功后，调用回调接口更新文件元信息（如 size/mime_type）

---

## ✅ 六、安全与扩展建议

- **权限控制**：普通用户只能读，上传需 RBAC（如 `role=admin`）
- **文件类型白名单**：只允许 `.pdf`, `.mp4`, `.pptx`, `.png` 等
- **病毒扫描**：集成云安全服务（如阿里云内容安全）
- **CDN 加速**：绑定自定义域名 + HTTPS + CDN
- **日志审计**：记录谁下载了什么资源（用于教学分析）

---

## ✅ 总结

| 步骤 | 内容 |
|------|------|
| 📦 数据库 | `learn_materials` 表，存元数据 |
| 🌐 API | `GET /learn`, `GET /learn/:id`, `POST /learn/presigned-url` |
| 🔒 安全 | 预签名 URL + JWT 认证 + 类型白名单 |
| ☁️ 存储 | 文件直传/直下 OSS，不经过应用服务器 |
| 🧩 扩展 | 支持视频转码、缩略图、审核等 |

---

需要我帮你生成：
- **GORM 模型代码**？
- **PostgreSQL migration 脚本**？
- **完整的 `LearnService` 接口与实现**？

欢迎继续提问！