package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	r_init "backend/internal/repository"
	"backend/internal/utils"
	"backend/logs"
	"backend/pkg/middleware"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Code:    400,
			Message: "请求数据格式错误",
			Error:   err.Error(),
		})
		return
	}

	// 检查用户是否已存在
	var existingUser model.User
	if err := r_init.DB.Where("name = ?", req.Name).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, dto.RegisterResponse{
			Code:    409,
			Message: "用户名已存在",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Code:    500,
			Message: "数据库查询失败",
			Error:   err.Error(),
		})
		return
	}

	// 检查邮箱是否已经被注册
	if err := r_init.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, dto.RegisterResponse{
			Code:    409,
			Message: "邮箱已被注册",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Code:    500,
			Message: "数据库查询失败",
			Error:   err.Error(),
		})
		return
	}

	// 校验token
	if !ValidateVerificationToken(req.Email, req.Token) {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Code:    400,
			Message: "验证码错误或已过期",
		})
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
		return
	}

	// 注册成功后删除token
	DeleteVerificationToken(req.Email)

	c.JSON(http.StatusOK, dto.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		ID:      newUser.ID,
		Name:    newUser.Name,
	})
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.LoginResponse{
			Code:    400,
			Message: "登录数据解析失败",
			Error:   err.Error(),
		})
		return
	}

	var user model.User
	if err := r_init.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, dto.LoginResponse{
				Code:    401,
				Message: "用户名或密码错误",
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.LoginResponse{
				Code:    500,
				Message: "数据库查询错误",
				Error:   err.Error(),
			})
		}
		return
	}

	// 使用 bcrypt 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, dto.LoginResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
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
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Code:    200,
		Message: "登录成功",
		Token:   tokenString,
		Name:    user.Name,
		ID:      user.ID,
	})
}

// 找回密码，发送验证码

// 发送注册验证码
func SendVerificationCode(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		logs.Sugar.Errorw("发送注册验证码", "detail", "解析请求失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误", "code": 400})
		return
	}

	// 检查邮箱格式
	if request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入邮箱地址", "code": 400})
		return
	}

	// 检查邮箱是否已经被注册
	var existingUser model.User
	if err := r_init.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "邮箱已被注册", "code": 409})
		return
	} else if err != gorm.ErrRecordNotFound {
		logs.Sugar.Errorw("发送注册验证码", "detail", "数据库查询失败")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询失败", "code": 500})
		return
	}

	// 生成验证码
	verificationCode := GenerateRandomToken(6)
	logs.Sugar.Infow("发送注册验证码", "email", request.Email, "code", verificationCode)

	// 保存token到Redis，过期时间10分钟
	if err := SaveVerificationToken(request.Email, verificationCode); err != nil {
		logs.Sugar.Errorw("发送注册验证码", "detail", "保存验证码失败")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存验证码失败", "code": 500})
		return
	}

	// 发送验证码邮件
	err := SendVerificationEmail(request.Email, verificationCode)
	if err != nil {
		logs.Sugar.Errorw("发送注册验证码", "detail", "发送验证码邮件失败")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "发送验证码失败", "code": 500})
		return
	}

	// TODO: 在实际应用中，应该将验证码存储在Redis或数据库中，设置过期时间
	// 这里暂时返回成功消息
	logs.Sugar.Infow("发送泣册验证码", "detail", "验证码发送成功")
	c.JSON(http.StatusOK, gin.H{
		"message": "验证码已发送，请检查邮箱",
		"code":    200,
	})
}

const charset = "0123456789"

func GenerateRandomToken(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[r.Intn(len(charset))]
	}

	return string(token)
}

// 处理重置密码请求
func RequestResetPassword(c *gin.Context) {
	// 实现请求重置密码的逻辑
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		logs.Sugar.Errorw("重置密码请求", "detail", "解析请求失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}

	// 查找用户
	var user model.User
	err := r_init.DB.Where("email = ?", request.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logs.Sugar.Errorw("重置密码请求", "detail", "用户未找到。")
			c.JSON(http.StatusNotFound, gin.H{"message": "用户未找到"})
		} else {
			logs.Sugar.Errorw("重置密码请求", "detail", "数据库查询失败。")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询失败"})
		}
		return
	}

	// 生成唯一的重置密码 token
	token := GenerateRandomToken(6) // 生成6位长度的token
	fmt.Println("密码找回时生成的token为：", token)
	// 在数据库中保存 token
	err = r_init.DB.Model(&user).Update("token", token).Error
	if err != nil {
		logs.Sugar.Errorw("重置密码请求", "detail", "保存 token 失败。")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存 token 失败"})
		return
	}

	// 发送重置密码邮件
	err = SendEmail(request.Email, token)
	if err != nil {
		logs.Sugar.Errorw("重置密码请求", "detail", "发送重置密码邮件失败。")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "发送重置密码邮件失败"})
		return
	}

	logs.Sugar.Infow("重置密码请求", "detail", "重置密码请求成功。")
	c.JSON(http.StatusOK, gin.H{
		"message": "重置密码请求成功",
	})
}

// 重置密码
func ResetPassword(c *gin.Context) {
	// 实现重置密码的逻辑
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BindJSON(&request); err != nil {
		logs.Sugar.Errorw("重置密码", "detail", "解析请求数据失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}
	var user model.User
	err := r_init.DB.Where("token = ?", request.Token).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logs.Sugar.Errorw("重置密码", "detail", "无效的重置密码 token。")
			c.JSON(http.StatusNotFound, gin.H{"message": "无效的重置密码 token"})
		} else {
			logs.Sugar.Errorw("重置密码", "detail", "数据库查询失败。")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询失败"})
		}
		return
	}

	err = r_init.DB.Model(&user).Update("password", request.NewPassword).Error
	if err != nil {
		logs.Sugar.Errorw("重置密码", "detail", "密码重置失败。")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码重置失败"})
		return
	}

	err = r_init.DB.Model(&user).Update("token", nil).Error
	if err != nil {
		logs.Sugar.Errorw("重置密码", "detail", "密码重置成功，但是 token 重置失败。")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码重置成功，但是 token 重置失败"})
		return
	}

	logs.Sugar.Infow("重置密码", "detail", "重置密码成功。")
	c.JSON(http.StatusOK, gin.H{
		"message": "重置密码成功",
	})
}

// 发送验证码
func SendEmail(email, token string) error {
	myEmail := os.Getenv("EMAIL_NAME")
	myPassword := os.Getenv("EMAIL_PASSWORD")
	// baseUrl := os.Getenv("BASE_URL")
	smtpServerHost := os.Getenv("SMTP_SERVER_HOST")
	smtpServerPortStr := os.Getenv("SMTP_SERVER_PORT")

	if myEmail == "" || myPassword == "" || smtpServerHost == "" || smtpServerPortStr == "" {
		logs.Sugar.Errorw("发送验证码邮件", "detail", "环境变量未正确设置")
		return errors.New("环境变量未正确设置")
	}

	smtpServerPort, err := strconv.Atoi(smtpServerPortStr)
	if err != nil {
		logs.Sugar.Errorw("发送验证码邮件", "detail", "将端口号转换为整数时出错: %v", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", myEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset Request")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>FormalLangLab</h1>
		<p>这是你的验证码：%s</p>
	`, token))

	d := gomail.NewDialer(smtpServerHost, smtpServerPort, myEmail, myPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 跳过证书验证，生产环境中应谨慎使用
	if err := d.DialAndSend(m); err != nil {
		if strings.Contains(err.Error(), "535") { // 例如，检查错误消息中是否包含 SMTP 身份验证失败的代码
			return errors.New("发送邮件失败可能是 SMTP 身份验证错误")
		} else if strings.Contains(err.Error(), "connection refused") {
			return errors.New("发送邮件失败：SMTP 服务器连接被拒绝")
		}
	}
	logs.Sugar.Infow("发送验证码邮件", "detail", "邮件发送成功")
	return nil
}

// 发送注册验证码邮件
func SendVerificationEmail(email, code string) error {
	myEmail := os.Getenv("EMAIL_NAME")
	myPassword := os.Getenv("EMAIL_PASSWORD")
	smtpServerHost := os.Getenv("SMTP_SERVER_HOST")
	smtpServerPortStr := os.Getenv("SMTP_SERVER_PORT")

	if myEmail == "" || myPassword == "" || smtpServerHost == "" || smtpServerPortStr == "" {
		logs.Sugar.Errorw("发送注册验证码邮件", "detail", "环境变量未正确设置")
		return errors.New("环境变量未正确设置")
	}

	smtpServerPort, err := strconv.Atoi(smtpServerPortStr)
	if err != nil {
		logs.Sugar.Errorw("发送注册验证码邮件", "detail", "将端口号转换为整数时出错: %v", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", myEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "GDesign 注册验证码")
	m.SetBody("text/html", fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f7f9fc;">
			<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; border-radius: 10px 10px 0 0; text-align: center;">
				<h1 style="margin: 0; font-size: 28px; font-weight: bold;">GDesign</h1>
				<p style="margin: 10px 0 0 0; font-size: 16px; opacity: 0.9;">形式语言与自动机学习系统</p>
			</div>
			<div style="background-color: white; padding: 40px; border-radius: 0 0 10px 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
				<h2 style="color: #333; margin-bottom: 20px; font-size: 24px;">注册验证码</h2>
				<p style="color: #666; font-size: 16px; line-height: 1.6; margin-bottom: 30px;">感谢您注册 GDesign 学习平台！请使用以下验证码完成注册：</p>
				<div style="background-color: #f8f9fa; border: 2px dashed #667eea; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0;">
					<span style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 6px; font-family: monospace;">%s</span>
				</div>
				<p style="color: #999; font-size: 14px; margin-bottom: 20px;">• 验证码有效期为 10 分钟</p>
				<p style="color: #999; font-size: 14px; margin-bottom: 20px;">• 请勿将验证码告知他人</p>
				<p style="color: #999; font-size: 14px;">如果您没有申请注册，请忽略此邮件。</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="color: #999; font-size: 12px; text-align: center;">此邮件由 GDesign 系统自动发送，请勿回复。</p>
			</div>
		</div>
	`, code))

	d := gomail.NewDialer(smtpServerHost, smtpServerPort, myEmail, myPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		if strings.Contains(err.Error(), "535") {
			return errors.New("发送邮件失败可能是 SMTP 身份验证错误")
		} else if strings.Contains(err.Error(), "connection refused") {
			return errors.New("发送邮件失败：SMTP 服务器连接被拒绝")
		}
		return err
	}
	logs.Sugar.Infow("发送注册验证码邮件", "detail", "邮件发送成功")
	return nil
}

// Redis Token 管理函数

// 保存验证码到Redis，过期时间10分钟
func SaveVerificationToken(email, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("verification_token:%s", email)

	err := r_init.RedisClient.Set(ctx, key, token, 10*time.Minute).Err()
	if err != nil {
		logs.Sugar.Errorw("保存验证码到Redis", "error", err)
		return err
	}

	logs.Sugar.Infow("验证码保存成功", "email", email, "key", key)
	return nil
}

// 校验验证码
func ValidateVerificationToken(email, token string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("verification_token:%s", email)

	storedToken, err := r_init.RedisClient.Get(ctx, key).Result()
	if err != nil {
		logs.Sugar.Errorw("获取验证码失败", "error", err, "email", email)
		return false
	}

	return storedToken == token
}

// 删除验证码
func DeleteVerificationToken(email string) {
	ctx := context.Background()
	key := fmt.Sprintf("verification_token:%s", email)

	err := r_init.RedisClient.Del(ctx, key).Err()
	if err != nil {
		logs.Sugar.Errorw("删除验证码失败", "error", err, "email", email)
	} else {
		logs.Sugar.Infow("验证码删除成功", "email", email)
	}
}
