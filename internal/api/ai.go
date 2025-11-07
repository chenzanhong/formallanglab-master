/*
AI模块 OpenAI Go SDK版本不低于 v2.4.0
*/
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/metrics"
	aiSvc "backend/internal/service/ai_s"
	"backend/logs"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AIhandler struct {
	aiService aiSvc.AIService
}

func NewAIHandler(aiService aiSvc.AIService) *AIhandler {
	return &AIhandler{aiService: aiService}
}

/*
流式对话
中文乱码	确保 Go 后端设置 Content-Type: text/event-stream; charset=utf-8
跨域（CORS）	在 Gin 中添加中间件：已添加
特殊字符（如 ?, &）	前端用 encodeURIComponent，后端自动解码
长连接超时	可在后端定期发送 : ping\n\n 保活（可选）
*/

func (h *AIhandler) AIChatSSE(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("ai", "chat_sse", time.Since(start).Seconds())
	}()

	var req dto.AIChatRequest // 给我5行四字成语，一行一个
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "message is required"})
		metrics.IncOperation("ai", "chat_sse", "failure: message required")
		return
	}

	// JWT 中间件中Set的
	usernameVal, _ := c.Get("username")
	username := usernameVal.(string)
	stream, session, err := h.aiService.StreamChat(c.Request.Context(), username, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		metrics.IncOperation("ai", "chat_sse", "failure: service error")
		return
	}

	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Accel-Buffering", "no") // Disable buffering for nginx

	// 用于记录完整 AI 响应（用于保存会话）
	var aiResp strings.Builder
	// 用于缓冲待发送的内容（提升性能）
	var buffer strings.Builder

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case <-c.Request.Context().Done():
				return
			case <-ticker.C:
				if buffer.Len() > 0 {
					// 写入客户端
					c.Writer.Write([]byte(buffer.String()))
					c.Writer.Flush()
					// 同步到完整响应记录
					aiResp.WriteString(buffer.String())
					buffer.Reset()
				}
			}
		}
	}()

	for stream.Next() {
		buffer.WriteString(stream.Current().Choices[0].Delta.Content)
	}

	// 最终 flush 剩余内容
	if buffer.Len() > 0 {
		c.Writer.Write([]byte(buffer.String()))
		c.Writer.Flush()
		aiResp.WriteString(buffer.String())
	}
	<-done

	if err := stream.Err(); err != nil {
		logs.Sugar.Errorf("AI stream error: %v", err)
	}

	go func() {
		ctx := context.Background()
		session.RecentTurns = append(session.RecentTurns, model.QAPair{
			User: req.Question,
			AI:   aiResp.String(),
		})
		session.Trim()
		if err := h.aiService.SaveSession(ctx, username, session); err != nil {
			logs.Sugar.Errorf("Failed to save AI session for user %s: %v", username, err)
		}
	}()

	metrics.IncOperation("ai", "chat_sse", "success")
}
