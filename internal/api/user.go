package api

import (
	"backend/internal/domain/dto"
	myErrors "backend/internal/errors"
	"backend/internal/metrics"
	userSvc "backend/internal/service/user_s"
	"backend/logs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService userSvc.UserService
}

func NewUserHandler(userService userSvc.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("user", "register", time.Since(start).Seconds())
	}()
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Code:    400,
			Message: "请求数据格式错误",
		})
		metrics.IncOperation("user", "register", "failure: parameter parsing error")
		return
	}

	newUser, err := h.userService.Register(c.Request.Context(), req.Name, req.Email, req.Password, req.Token)
	if err != nil {
		switch err {
		case myErrors.ErrPasswordHashFailed:
			c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
				Code:    500,
				Message: "密码加密失败",
			})
			metrics.IncOperation("user", "register", "failure: password encryption error")
			return
		case myErrors.ErrUserCreationFailed:
			c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
				Code:    500,
				Message: "用户创建失败",
			})
			metrics.IncOperation("user", "register", "failure: user creation error")
			return
		default:
			c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
				Code:    500,
				Message: "注册失败",
				Error:   err.Error(),
			})
			metrics.IncOperation("user", "register", "failure: unknown error")
			return
		}
	}

	metrics.IncOperation("user", "register", "success")
	c.JSON(http.StatusOK, dto.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		ID:      newUser.ID,
		Name:    newUser.Name,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("user", "login", time.Since(start).Seconds())
	}()
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.LoginResponse{
			Code:    400,
			Message: "登录数据解析失败",
			Error:   err.Error(),
		})
		metrics.IncOperation("user", "login", "failure: parameter parsing error")
		return
	}

	// 调用服务层登录逻辑
	tokenString, user, err := h.userService.Login(c.Request.Context(), req.Name, req.Password)
	if err != nil {
		switch err {
		case myErrors.ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, dto.LoginResponse{
				Code:    401,
				Message: "用户名或密码错误",
			})
			metrics.IncOperation("user", "login", "failure: user not found")
			return
		case myErrors.ErrUserCreationFailed:
			c.JSON(http.StatusUnauthorized, dto.LoginResponse{
				Code:    401,
				Message: "用户名或密码错误",
			})
			metrics.IncOperation("user", "login", "failure: invalid password")
			return
		default:
			c.JSON(http.StatusInternalServerError, dto.LoginResponse{
				Code:    500,
				Message: "登录失败",
				Error:   err.Error(),
			})
			metrics.IncOperation("user", "login", "failure: unknown error")
			return
		}
	}

	metrics.IncOperation("user", "login", "success")
	c.JSON(http.StatusOK, dto.LoginResponse{
		Code:    200,
		Message: "登录成功",
		Token:   tokenString,
		Name:    user.Name,
		ID:      user.ID,
	})
}

// 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("email", "reset_password", time.Since(start).Seconds())
	}()
	// 实现重置密码的逻辑
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BindJSON(&request); err != nil {
		logs.Sugar.Errorw("重置密码", "detail", "解析请求数据失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		metrics.IncOperation("email", "reset_password", "failure: parameter parsing error")
		return
	}

	if request.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "新密码不能为空"})
		metrics.IncOperation("email", "reset_password", "failure: empty password")
		return
	}

	// 验证 token
	err := h.userService.ResetPassword(c.Request.Context(), request.Token, request.NewPassword)
	if err != nil {
		switch err {
		case myErrors.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"message": "验证码错误或已过期"})
			metrics.IncOperation("email", "reset_password", "failure: invalid token")
			return
		case myErrors.ErrPasswordHashFailed:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败", "code": 500})
			metrics.IncOperation("email", "reset_password", "failure: password encryption error")
			return
		case myErrors.ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			metrics.IncOperation("email", "reset_password", "failure: user not found")
			return
		default:
			logs.Sugar.Errorw("重置密码", "detail", "重置密码失败")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "重置密码失败", "code": 500})
			metrics.IncOperation("email", "reset_password", "failure: reset password error")
			return
		}
	}

	logs.Sugar.Infow("重置密码", "detail", "重置密码成功。")
	metrics.IncOperation("email", "reset_password", "success")
	c.JSON(http.StatusOK, gin.H{
		"message": "重置密码成功",
	})
}
