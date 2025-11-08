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

func (h *AIhandler) AIChatSSE(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("ai", "chat_sse", time.Since(start).Seconds())
	}()

	var req dto.AIChatRequest
	if err := c.BindJSON(&req); err != nil {
		metrics.IncOperation("ai", "chat_sse", "failure: request body required")
		logs.Sugar.Errorw("AI对话请求失败", "detail", "请求体不能为空")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "request body is required", "result": false})
		return
	}

	// JWT 中间件中Set的
	usernameVal, _ := c.Get("username")
	username := usernameVal.(string)
	stream, session, err := h.aiService.StreamChat(c.Request.Context(), username, &req)
	if err != nil {
		metrics.IncOperation("ai", "chat_sse", "failure: service error")
		logs.Sugar.Errorw("AI对话请求失败", "detail", "AI服务调用失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "result": false})
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

	ticker := time.NewTicker(64 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if buffer.Len() > 0 {
					// 写入客户端
					// 不使用c.SSE，ai响应本身就是流式，无需再SSE
					c.Writer.Write([]byte(buffer.String()))
					c.Writer.Flush()
					// 同步到完整响应记录
					aiResp.WriteString(buffer.String())
					buffer.Reset()
				}
			case <-done:
				return
			}
		}
	}()

	for stream.Next() {
		buffer.WriteString(stream.Current().Choices[0].Delta.Content)
	}

	// 停止ticker并关闭done通道
	ticker.Stop()
	close(done)

	// 最终 flush 剩余内容
	if buffer.Len() > 0 {
		c.Writer.Write([]byte(buffer.String()))
		c.Writer.Flush()
		aiResp.WriteString(buffer.String())
		buffer.Reset()
	}

	if err := stream.Err(); err != nil {
		metrics.IncOperation("ai", "chat_sse", "failure: stream error")
		logs.Sugar.Errorw("AI流式传输失败", "detail", "流式传输过程中发生错误")
	}

	go func() {
		ctx := context.Background()
		session.RecentTurns = append(session.RecentTurns, model.QAPair{
			User: req.Question,
			AI:   aiResp.String(),
		})
		session.Trim()
		if err := h.aiService.SaveSession(ctx, username, session); err != nil {
			metrics.IncOperation("ai", "chat_sse", "failure: save session error")
			logs.Sugar.Errorw("AI会话保存失败", "detail", "无法保存用户会话信息")
		}
	}()

	metrics.IncOperation("ai", "chat_sse", "success")
	logs.Sugar.Infow("AI对话请求成功")
}
