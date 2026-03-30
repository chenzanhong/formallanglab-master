# 学习资料上传指南

## 功能说明

本系统支持学习资料的下载功能，资料存储在阿里云 OSS 中，通过数据库管理元数据。

## 已完成的功能

✅ 后端 API 接口已实现：
- `GET /master/learn` - 获取学习资料列表
- `GET /master/learn/{id}` - 获取单个资料详情（包含下载链接）
- `POST /master/learn` - 添加学习资源（需要先上传到 OSS）
- `DELETE /master/learn/{id}` - 删除学习资源
- `POST /master/learn/sync` - 同步 OSS 文件到数据库

✅ 数据库表已创建：`learn_materials`

✅ OSS 客户端已实现（支持阿里云 OSS 和模拟模式）

## 简单上传方式

### 方案：直接放在 materials/ 目录下

1. **在阿里云 OSS 控制台创建 `materials/` 目录**
   - 进入你的 Bucket
   - 新建文件夹，命名为 `materials`

2. **直接上传文件到 materials/ 目录**
   ```
   materials/
   ├── DFA教程.pptx
   ├── NFA指南.pptx
   ├── 正则表达式基础.pptx
   └── 文法转换讲解.pptx
   ```

3. **同步到数据库**
   
   调用同步接口自动识别新增文件：
   ```bash
   curl -X POST http://localhost:8081/master/learn/sync
   ```

4. **验证上传结果**
   ```bash
   curl http://localhost:8081/master/learn
   ```

## OSS 文件命名规范

```
materials/{filename}.{ext}
```

支持的格式：
- PDF: `.pdf`
- PPT: `.ppt`, `.pptx`
- DOC: `.doc`, `.docx`
- TXT: `.txt`

示例：
- `materials/DFA教程.pptx`
- `materials/NFA指南.pptx`
- `materials/正则表达式基础.pdf`

## 分类说明

⚠️ **注意：系统暂时不使用分类功能**

- 所有上传的文件都会被归类为 `general`
- 分类字段保留是为了未来可能的扩展需求
- 前端只显示"全部"按钮，不提供分类筛选

## 下载流程

1. 用户调用 `GET /master/learn` 获取资料列表
2. 用户点击下载按钮
3. 后端生成 5 分钟有效的预签名下载链接
4. 用户通过预签名链接直接从 OSS 下载文件

## 环境变量配置

确保 `.env` 文件中配置了阿里云 OSS 信息：

```bash
OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
OSS_ACCESS_KEY_ID=your_access_key_id
OSS_ACCESS_KEY_SECRET=your_access_key_secret
OSS_BUCKET_NAME=your-bucket-name
OSS_BASE_URL=https://your-bucket-name.oss-cn-hangzhou.aliyuncs.com
```

## 开发测试

如果还没有阿里云 OSS 账号，可以使用模拟模式进行开发测试：

```go
// 在 main.go 中
ossClient := oss.NewMockOSSClient("http://localhost:9000")
```

## 常见问题

### Q: 上传的文件在哪里查看？
A: 上传后需要调用 `/master/learn/sync` 接口同步到数据库，然后通过 `/master/learn` 接口查看。

### Q: 下载链接有效期是多久？
A: 预签名下载链接有效期为 5 分钟。

### Q: 可以直接下载 OSS 文件吗？
A: 可以，但需要生成预签名 URL。直接暴露 OSS URL 不安全。

### Q: 如何删除学习资料？
A: 调用 `DELETE /master/learn/{id}` 接口，只会删除数据库记录，不会删除 OSS 中的文件。如需删除 OSS 文件，需要手动在 OSS 控制台操作。

### Q: 文件名需要特殊格式吗？
A: 不需要，系统会自动识别分类。但建议使用有意义的文件名，方便管理。

### Q: 为什么没有分类功能？
A: 系统暂时不使用分类功能，所有文件统一归类为 `general`。分类字段保留是为了未来可能的扩展需求。

## 下一步

1. 在阿里云 OSS 创建 `materials/` 目录
2. 上传你的 PPT 文件到该目录
3. 调用 sync 接口同步
4. 在前端页面查看和下载资料

