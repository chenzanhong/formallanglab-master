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
		metrics.IncOperation("user", "register", "failure: parameter parsing error")
		logs.Sugar.Warnw("请求数据格式错误", "detail", err.Error())
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Result: false,
			Msg:    "请求数据格式错误",
		})
		return
	}

	newUser, err := h.userService.Register(c.Request.Context(), req.Name, req.Email, req.Password, req.Token)
	if err != nil {
		switch err {
		case myErrors.ErrPasswordHashFailed:
			metrics.IncOperation("user", "register", "failure: password encryption error")
			logs.Sugar.Errorw("密码加密失败")
			c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
				Result: false,
				Msg:    "密码加密失败",
			})
			return
		case myErrors.ErrUserCreationFailed:
			metrics.IncOperation("user", "register", "failure: user creation error")
			logs.Sugar.Errorw("用户创建失败")
			c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
				Result: false,
				Msg:    "用户创建失败",
			})
			return
		default:
			metrics.IncOperation("user", "register", "failure: unknown error")
			logs.Sugar.Warnw("注册失败", "detail", err.Error())
			c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
				Result: false,
				Msg:    "注册失败",
				Error:  err.Error(),
			})
			return
		}
	}

	metrics.IncOperation("user", "register", "success")
	logs.Sugar.Infow("注册成功", "user_id", newUser.ID, "username", newUser.Name)
	c.JSON(http.StatusOK, dto.RegisterResponse{
		Result: true,
		Msg:    "注册成功",
		ID:     newUser.ID,
		Name:   newUser.Name,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("user", "login", time.Since(start).Seconds())
	}()
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("user", "login", "failure: parameter parsing error")
		logs.Sugar.Warnw("登录数据解析失败", "detail", err.Error())
		c.JSON(http.StatusBadRequest, dto.LoginResponse{
			Result: false,
			Msg:    "登录数据解析失败",
			Error:  err.Error(),
		})
		return
	}

	// 调用服务层登录逻辑
	tokenString, user, err := h.userService.Login(c.Request.Context(), req.Name, req.Password)
	if err != nil {
		switch err {
		case myErrors.ErrUserNotFound:
			metrics.IncOperation("user", "login", "failure: user not found")
			logs.Sugar.Warnw("用户名或密码错误", "username", req.Name)
			c.JSON(http.StatusUnauthorized, dto.LoginResponse{
				Result: false,
				Msg:    "用户名或密码错误",
			})
			return
		case myErrors.ErrUserCreationFailed:
			metrics.IncOperation("user", "login", "failure: invalid password")
			logs.Sugar.Warnw("用户名或密码错误", "username", req.Name)
			c.JSON(http.StatusUnauthorized, dto.LoginResponse{
				Result: false,
				Msg:    "用户名或密码错误",
			})
			return
		default:
			metrics.IncOperation("user", "login", "failure: unknown error")
			logs.Sugar.Warnw("登录失败", "detail", err.Error())
			c.JSON(http.StatusInternalServerError, dto.LoginResponse{
				Result: false,
				Msg:    "登录失败",
				Error:  err.Error(),
			})
			return
		}
	}

	metrics.IncOperation("user", "login", "success")
	logs.Sugar.Infow("登录成功", "user_id", user.ID, "username", user.Name)
	c.JSON(http.StatusOK, dto.LoginResponse{
		Result: true,
		Msg:    "登录成功",
		Token:  tokenString,
		Name:   user.Name,
		ID:     user.ID,
	})
}

// 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("email", "reset_password", time.Since(start).Seconds())
	}()
	// 实现重置密码的逻辑
	var request dto.ResetPasswordRequest

	if err := c.BindJSON(&request); err != nil {
		metrics.IncOperation("email", "reset_password", "failure: parameter parsing error")
		logs.Sugar.Warnw("重置密码失败", "detail", "解析请求数据失败")
		c.JSON(http.StatusBadRequest, dto.ResetPasswordResponse{Msg: "请求数据格式错误", Result: false})
		return
	}

	if request.NewPassword == "" {
		metrics.IncOperation("email", "reset_password", "failure: empty password")
		logs.Sugar.Warnw("重置密码失败", "detail", "新密码为空")
		c.JSON(http.StatusBadRequest, dto.ResetPasswordResponse{Msg: "新密码不能为空", Result: false})
		return
	}

	err := h.userService.ResetPassword(c.Request.Context(), request.Token, request.NewPassword)
	if err != nil {
		switch err {
		case myErrors.ErrInvalidToken:
			metrics.IncOperation("email", "reset_password", "failure: invalid token")
			logs.Sugar.Warnw("重置密码失败", "detail", "验证码错误或已过期")
			c.JSON(http.StatusUnauthorized, dto.ResetPasswordResponse{Msg: "验证码错误或已过期", Result: false})
			return
		case myErrors.ErrPasswordHashFailed:
			metrics.IncOperation("email", "reset_password", "failure: password encryption error")
			logs.Sugar.Errorw("重置密码失败", "detail", "密码加密失败")
			c.JSON(http.StatusInternalServerError, dto.ResetPasswordResponse{Msg: "密码加密失败", Result: false})
			return
		case myErrors.ErrUserNotFound:
			metrics.IncOperation("email", "reset_password", "failure: user not found")
			logs.Sugar.Warnw("重置密码失败", "detail", "用户不存在")
			c.JSON(http.StatusUnauthorized, dto.ResetPasswordResponse{Msg: "用户不存在", Result: false})
			return
		default:
			metrics.IncOperation("email", "reset_password", "failure: reset password error")
			logs.Sugar.Warnw("重置密码失败", "detail", err.Error())
			c.JSON(http.StatusInternalServerError, dto.ResetPasswordResponse{Msg: "重置密码失败", Result: false})
			return
		}
	}

	metrics.IncOperation("email", "reset_password", "success")
	logs.Sugar.Infow("重置密码成功")
	c.JSON(http.StatusOK, dto.ResetPasswordResponse{
		Msg:    "重置密码成功",
		Result: true,
	})
}

func (h *UserHandler) CheckMe(c *gin.Context) {
	// 不做任何处理，只是借助JWT判断token是否还有效
	c.JSON(http.StatusOK, gin.H{})
}
