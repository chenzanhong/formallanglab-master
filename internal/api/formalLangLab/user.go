package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/metrics"
	r_init "backend/internal/repository"
	email "backend/internal/service/email_s"
	"backend/internal/utils"
	"backend/pkg/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register 用户注册
func Register(c *gin.Context) {
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

	// 检查用户是否已存在
	var existingUser model.User
	if err := r_init.DB.Where("name = ?", req.Name).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, dto.RegisterResponse{
			Code:    409,
			Message: "用户名已存在",
		})
		metrics.IncOperation("user", "register", "failure: username exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Code:    500,
			Message: "数据库查询失败",
			Error:   err.Error(),
		})
		metrics.IncOperation("user", "register", "failure: db query error")
		return
	}

	// 检查邮箱是否已经被注册
	if err := r_init.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, dto.RegisterResponse{
			Code:    409,
			Message: "邮箱已被注册",
		})
		metrics.IncOperation("user", "register", "failure: email exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Code:    500,
			Message: "数据库查询失败",
			Error:   err.Error(),
		})
		metrics.IncOperation("user", "register", "failure: db query error")
		return
	}

	// 校验token
	if !email.ValidateRegisterVerificationToken(req.Email, req.Token) {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Code:    400,
			Message: "验证码错误或已过期",
		})
		metrics.IncOperation("user", "register", "failure: invalid token")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Code:    500,
			Message: "密码加密失败",
			Error:   err.Error(),
		})
		metrics.IncOperation("user", "register", "failure: password encryption error")
		return
	}

	// 创建用户
	newUser := model.User{
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := r_init.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Code:    500,
			Message: "用户创建失败",
			Error:   err.Error(),
		})
		metrics.IncOperation("user", "register", "failure: user creation error")
		return
	}

	// 注册成功后删除token，不处理 error
	email.DeleteRegisterVerificationToken(req.Email, req.Token)

	metrics.IncOperation("user", "register", "success")
	c.JSON(http.StatusOK, dto.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		ID:      newUser.ID,
		Name:    newUser.Name,
	})
}

func Login(c *gin.Context) {
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

	var user model.User
	if err := r_init.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, dto.LoginResponse{
				Code:    401,
				Message: "用户名或密码错误",
			})
			metrics.IncOperation("user", "login", "failure: user not found")
		} else {
			c.JSON(http.StatusInternalServerError, dto.LoginResponse{
				Code:    500,
				Message: "数据库查询错误",
				Error:   err.Error(),
			})
			metrics.IncOperation("user", "login", "failure: db query error")
		}
		return
	}

	// 使用 bcrypt 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, dto.LoginResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
		metrics.IncOperation("user", "login", "failure: invalid password")
		return
	}

	// 生成 JWT
	tokenString, err := middleware.GenerateToken(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.LoginResponse{
			Code:    500,
			Message: "生成Token失败",
			Error:   err.Error(),
		})
		metrics.IncOperation("user", "login", "failure: token generation error")
		return
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
