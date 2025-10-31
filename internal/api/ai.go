/*
AI模块
*/
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
流式对话
中文乱码	确保 Go 后端设置 Content-Type: text/event-stream; charset=utf-8
跨域（CORS）	在 Gin 中添加中间件：
r.Use(cors.Default())（需 github.com/gin-contrib/cors）
特殊字符（如 ?, &）	前端用 encodeURIComponent，后端自动解码
长连接超时	可在后端定期发送 : ping\n\n 保活（可选）
*/
func AIChatSSE(c *gin.Context) {
	msg := c.Query("message")
	if msg == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "message is required"})
		return
	}
}
