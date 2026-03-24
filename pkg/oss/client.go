package oss

import (
	"fmt"
	"os"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Client OSS客户端接口
type Client interface {
	// GeneratePresignedURL 生成下载预签名URL
	GeneratePresignedURL(key string, expire time.Duration) (string, error)
	// GenerateUploadPresignedURL 生成上传预签名URL
	GenerateUploadPresignedURL(key string, expire time.Duration) (string, error)
	// CheckObjectExists 检查对象是否存在
	CheckObjectExists(key string) (bool, error)
	// ListObjects 列出指定前缀的所有对象
	ListObjects(prefix string) ([]OSSObjectInfo, error)
}

// OSSObjectInfo OSS对象信息
type OSSObjectInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
	ETag         string
}

// MockOSSClient 模拟OSS客户端实现（用于开发测试）
type MockOSSClient struct {
	BaseURL string
}

// NewMockOSSClient 创建模拟OSS客户端
func NewMockOSSClient(baseURL string) *MockOSSClient {
	return &MockOSSClient{BaseURL: baseURL}
}

// GeneratePresignedURL 生成模拟的下载预签名URL
func (c *MockOSSClient) GeneratePresignedURL(key string, expire time.Duration) (string, error) {
	// 模拟生成预签名URL
	return c.BaseURL + "/" + key + "?expire=" + time.Now().Add(expire).Format(time.RFC3339), nil
}

// GenerateUploadPresignedURL 生成模拟的上传预签名URL
func (c *MockOSSClient) GenerateUploadPresignedURL(key string, expire time.Duration) (string, error) {
	// 模拟生成上传预签名URL
	return c.BaseURL + "/upload/" + key + "?expire=" + time.Now().Add(expire).Format(time.RFC3339), nil
}

// CheckObjectExists 模拟检查对象是否存在
func (c *MockOSSClient) CheckObjectExists(key string) (bool, error) {
	// 简单模拟，始终返回存在
	// 在实际使用中，可以根据需要调整逻辑
	return true, nil
}

// ListObjects 模拟列出对象
func (c *MockOSSClient) ListObjects(prefix string) ([]OSSObjectInfo, error) {
	// 模拟返回一些测试数据
	return []OSSObjectInfo{}, nil
}

// AliyunOSSClient 阿里云OSS客户端实现
type AliyunOSSClient struct {
	client     *oss.Client
	service    *oss.Bucket
	baseURL    string
	bucketName string
}

// NewAliyunOSSClient 创建阿里云OSS客户端
func NewAliyunOSSClient() (*AliyunOSSClient, error) {
	// 从环境变量中读取配置
	endpoint := os.Getenv("OSS_ENDPOINT")
	accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	bucketName := os.Getenv("OSS_BUCKET_NAME")
	baseURL := os.Getenv("OSS_BASE_URL")

	// 验证必要的配置
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return nil, fmt.Errorf("阿里云OSS配置不完整，请检查环境变量")
	}

	// 创建OSS客户端
	client, err := oss.New(endpoint,
		accessKeyID,
		accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	// 获取存储空间
	service, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("获取存储空间失败: %w", err)
	}

	return &AliyunOSSClient{
		client:     client,
		service:    service,
		baseURL:    baseURL,
		bucketName: bucketName,
	}, nil
}

// GeneratePresignedURL 生成阿里云OSS下载预签名URL
func (c *AliyunOSSClient) GeneratePresignedURL(key string, expire time.Duration) (string, error) {
	// 生成预签名URL用于下载
	u, err := c.service.SignURL(key, oss.HTTPGet, int64(expire.Seconds()))
	if err != nil {
		return "", fmt.Errorf("生成下载预签名URL失败: %w", err)
	}

	return u, nil
}

// GenerateUploadPresignedURL 生成阿里云OSS上传预签名URL
func (c *AliyunOSSClient) GenerateUploadPresignedURL(key string, expire time.Duration) (string, error) {
	// 生成预签名URL用于上传
	u, err := c.service.SignURL(key, oss.HTTPPut, int64(expire.Seconds()))
	if err != nil {
		return "", fmt.Errorf("生成上传预签名URL失败: %w", err)
	}

	return u, nil
}

// CheckObjectExists 检查阿里云OSS对象是否存在
func (c *AliyunOSSClient) CheckObjectExists(key string) (bool, error) {
	// 调用阿里云SDK检查对象是否存在
	exists, err := c.service.IsObjectExist(key)
	if err != nil {
		return false, fmt.Errorf("检查对象是否存在失败: %w", err)
	}

	return exists, nil
}

// ListObjects 列出阿里云OSS指定前缀的所有对象
func (c *AliyunOSSClient) ListObjects(prefix string) ([]OSSObjectInfo, error) {
	// 调用阿里云SDK列出对象
	lor, err := c.service.ListObjects(oss.Prefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("列出对象失败: %w", err)
	}

	// 转换为OSSObjectInfo数组
	objects := make([]OSSObjectInfo, 0, len(lor.Objects))
	for _, obj := range lor.Objects {
		objects = append(objects, OSSObjectInfo{
			Key:          obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
			ETag:         obj.ETag,
		})
	}

	return objects, nil
}
