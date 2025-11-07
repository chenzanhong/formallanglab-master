// backend/internal/api/email/email.go
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/metrics"
	email "backend/internal/service/email_s"
	"backend/logs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService email.EmailService
}

func NewEmailHandler(emailService email.EmailService) *EmailHandler {
	return &EmailHandler{emailService: emailService}
}

// 注册，发送验证码
func (h *EmailHandler) SendRegisterVerificationCode(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("email", "send_register_code", time.Since(start).Seconds())
	}()
	var request dto.EmailRequest

	if err := c.BindJSON(&request); err != nil {
		logs.Sugar.Errorw("发送注册验证码", "detail", "解析请求失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误", "code": 400})
		metrics.IncOperation("email", "send_register_code", "failure: parameter parsing error")
		return
	}

	// 检查邮箱格式
	if request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入邮箱地址", "code": 400})
		metrics.IncOperation("email", "send_register_code", "failure: empty email")
		return
	}

	err := h.emailService.SendRegisterVerificationCode(c.Request.Context(), request.Email)
	if err != nil {
		logs.Sugar.Errorw("发送注册验证码", "detail", "验证码发送失败")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "验证码发送失败", "code": 500})
		metrics.IncOperation("email", "send_register_code", "failure: send code error")
		return
	}

	// 这里暂时返回成功消息
	logs.Sugar.Infow("发送注册验证码", "detail", "验证码发送成功")
	metrics.IncOperation("email", "send_register_code", "success")
	c.JSON(http.StatusOK, gin.H{
		"message": "验证码已发送，请检查邮箱",
		"code":    200,
	})
}

// 重置密码的验证码
func (h *EmailHandler) SendResetPwdVerificationCode(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("email", "send_reset_password_code", time.Since(start).Seconds())
	}()
	// 实现请求重置密码的逻辑
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		logs.Sugar.Errorw("重置密码请求", "detail", "解析请求失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		metrics.IncOperation("email", "send_reset_password_code", "failure: parameter parsing error")
		return
	}

	// 检查邮箱格式
	if request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入邮箱地址", "code": 400})
		metrics.IncOperation("email", "send_reset_password_code", "failure: empty email")
		return
	}

	err := h.emailService.SendResetPwdVerificationCode(c.Request.Context(), request.Email)
	if err != nil {
		logs.Sugar.Errorw("重置密码请求", "detail", "重置密码请求失败")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "重置密码请求失败", "code": 500})
		metrics.IncOperation("email", "send_reset_password_code", "failure: send code error")
		return
	}

	logs.Sugar.Infow("重置密码请求", "detail", "重置密码请求成功。")
	metrics.IncOperation("email", "send_reset_password_code", "success")
	c.JSON(http.StatusOK, gin.H{
		"message": "重置密码请求成功",
	})
}
