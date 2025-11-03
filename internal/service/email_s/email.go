// backend/internal/api/email/email.go
package email

import (
	"backend/internal/domain/model"
	r_init "backend/internal/repository"
	kafka_s "backend/internal/service/kafka_s"
	"backend/internal/utils"
	"backend/logs"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ====== 频率限制与异步发送逻辑 ======
const (
	EmailTypeRegister = "register"
	EmailTypeResetPwd = "reset"

	// ⏱️ 请求间隔限制：1分钟内只能发1次（防刷）
	EmailSendRateLimitTTL = 1 * time.Minute
	// 验证码有效期保持1分钟
	VerificationTokenTTL = 1 * time.Minute
)

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

// ======================= 服务 =======================

// 服务：注册账号，发送验证码
func SendRegisterVerificationCodeService(email string) error {
	// 限制频率，验证码有效期一分钟，不能重复发送
	if hasRegisterVerificationToken(email) {
		logs.Sugar.Errorw("发送注册验证码", "detail", "操作太频繁，请稍后重试")
		return errors.New("操作太频繁，请稍后重试")
	}

	// 检查邮箱是否存在
	var existingUser model.User
	if err := r_init.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		logs.Sugar.Errorw("发送注册验证码", "detail", "邮箱已存在")
		return errors.New("邮箱已存在")
	} else if err != gorm.ErrRecordNotFound {
		logs.Sugar.Errorw("发送注册验证码", "detail", "数据库查询失败")
		return errors.New("数据库查询失败")
	}

	// 生成验证码
	verificationCode := GenerateRandomToken(6)
	logs.Sugar.Infow("发送注册验证码", "email", email, "code", verificationCode)

	// 保存token到Redis，过期时间1分钟
	if err := saveRegisterVerificationToken(email, verificationCode); err != nil {
		return errors.New("保存验证码失败")
	}

	// （内部）异步
	return SendRegisterEmail(email, verificationCode)
}

// 服务：处理重置密码请求
func SendResetPwdVerificationCodeService(email string) error {
	// 限制频率，验证码有效期一分钟，不能重复发送
	if hasResetPwdToken(email) {
		logs.Sugar.Errorw("发送重置密码验证码", "detail", "操作太频繁，请稍后重试")
		return errors.New("操作太频繁，请稍后重试")
	}

	// 查找用户
	var user model.User
	err := r_init.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logs.Sugar.Errorw("重置密码请求", "detail", "用户未找到。")
			return errors.New("用户未找到")
		} else {
			logs.Sugar.Errorw("重置密码请求", "detail", "数据库查询失败。")
			return errors.New("数据库查询失败")
		}
	}

	// 生成 6 位数字 token
	token := GenerateRandomToken(6)
	logs.Sugar.Infow("生成找回密码 token", "email", email, "token", token)

	// 保存到 Redis，1 分钟过期
	if err := saveResetPwdToken(token, email); err != nil {
		logs.Sugar.Errorw("保存找回密码 token 失败", "error", err)
		return errors.New("系统繁忙，请稍后重试")
	}

	// 发送重置密码邮件（异步）
	SendResetPwdEmail(email, token)

	logs.Sugar.Infow("重置密码请求", "detail", "重置密码请求成功。")
	return nil
}

// 服务：重置密码
func ResetPassword(token, newPassword string) error {
	email, err := getEmailByResetPwdToken(token)
	if err != nil {
		logs.Sugar.Errorw("无效或过期的重置 token", "token", token)
		return errors.New("无效或过期的重置链接")
	}

	// 根据 email 查用户
	var user model.User
	if err = r_init.DB.Where("email = ?", email).First(&user).Error; err != nil {
		logs.Sugar.Errorw("根据 email 查不到用户", "email", email)
		return errors.New("用户异常")
	}

	// 更新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logs.Sugar.Errorw("密码加密失败", "error", err)
		return errors.New("密码加密失败")
	}
	if err := r_init.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		logs.Sugar.Errorw("密码更新失败", "error", err)
		return errors.New("密码更新失败")
	}

	logs.Sugar.Infow("重置密码", "detail", "重置密码成功。")
	return nil
}

// ======================= 邮件生产者=======================

// 发送注册验证码邮件
func SendRegisterEmail(email, code string) error {
	subject := "FormalLangLab 注册验证码"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f7f9fc;">
			<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; border-radius: 10px 10px 0 0; text-align: center;">
				<h1 style="margin: 0; font-size: 28px; font-weight: bold;">FormalLangLab</h1>
				<p style="margin: 10px 0 0 0; font-size: 16px; opacity: 0.9;">形式语言与自动机学习系统</p>
			</div>
			<div style="background-color: white; padding: 40px; border-radius: 0 0 10px 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
				<h2 style="color: #333; margin-bottom: 20px; font-size: 24px;">注册验证码</h2>
				<p style="color: #666; font-size: 16px; line-height: 1.6; margin-bottom: 30px;">感谢您注册 FormalLangLab 学习平台！请使用以下验证码完成注册：</p>
				<div style="background-color: #f8f9fa; border: 2px dashed #667eea; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0;">
					<span style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 6px; font-family: monospace;">%s</span>
				</div>
				<p style="color: #999; font-size: 14px; margin-bottom: 20px;">• 验证码有效期为 1 分钟</p>
				<p style="color: #999; font-size: 14px; margin-bottom: 20px;">• 请勿将验证码告知他人</p>
				<p style="color: #999; font-size: 14px;">如果您没有申请注册，请忽略此邮件。</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="color: #999; font-size: 12px; text-align: center;">此邮件由 FormalLangLab 系统自动发送，请勿回复。</p>
			</div>
		</div>
	`, code)

	return SendEmailViaKafka(email, subject, "text/html", body)
}

// 发送重置密码的验证码邮件
func SendResetPwdEmail(email, token string) error {
	subject := "FormalLangLab 重置密码"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f7f9fc;">
			<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; border-radius: 10px 10px 0 0; text-align: center;">
				<h1 style="margin: 0; font-size: 28px; font-weight: bold;">FormalLangLab</h1>
				<p style="margin: 10px 0 0 0; font-size: 16px; opacity: 0.9;">形式语言与自动机学习系统</p>
			</div>
			<div style="background-color: white; padding: 40px; border-radius: 0 0 10px 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
				<h2 style="color: #333; margin-bottom: 20px; font-size: 24px;">重置密码</h2>
				<p style="color: #666; font-size: 16px; line-height: 1.6; margin-bottom: 30px;">您正在请求重置 FormalLangLab 学习平台的密码，请使用以下验证码完成操作：</p>
				<div style="background-color: #f8f9fa; border: 2px dashed #667eea; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0;">
					<span style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 6px; font-family: monospace;">%s</span>
				</div>
				<p style="color: #999; font-size: 14px; margin-bottom: 20px;">• 验证码有效期为 1 分钟</p>
				<p style="color: #999; font-size: 14px; margin-bottom: 20px;">• 请勿将验证码告知他人</p>
				<p style="color: #999; font-size: 14px;">如果您没有请求重置密码，请忽略此邮件。</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="color: #999; font-size: 12px; text-align: center;">此邮件由 FormalLangLab 系统自动发送，请勿回复。</p>
			</div>
		</div>
	`, token)

	return SendEmailViaKafka(email, subject, "text/html", body)
}

// SendEmailViaKafka 发送邮件事件到 Kafka
func SendEmailViaKafka(email, subject, contextType, body string) error {
	event := &kafka_s.EmailEvent{
		To:          email,
		Subject:     subject,
		ContentType: contextType,
		Body:        body,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := kafka_s.SendEmailEvent(ctx, event); err != nil {
		return fmt.Errorf("发送邮件事件到 Kafka 失败: %w", err)
	}

	return nil
}

// ======================= Redis Token 管理函数 =======================

// ============= 注册相关 ==============

// 保存注册验证码到Redis，过期时间10分钟
func saveRegisterVerificationToken(email, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("register_verification_token:%s", email)

	err := r_init.RedisClient.Set(ctx, key, token, VerificationTokenTTL).Err()
	if err != nil {
		logs.Sugar.Errorw("保存注册验证码到Redis", "detail", "保存验证码失败", "error", err,"email", email, "token", token)
		return err
	}

	logs.Sugar.Infow("注册验证码保存成功", "detail", "保存注册验证码成功", "email", email, "token", token)
	return nil
}

func hasRegisterVerificationToken(email string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("register_verification_token:%s", email)

	_, err := r_init.RedisClient.Get(ctx, key).Result()
	return err == nil
}

// 校验注册验证码
func ValidateRegisterVerificationToken(email, token string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("register_verification_token:%s", email)

	storedToken, err := r_init.RedisClient.Get(ctx, key).Result()
	if err != nil {
		logs.Sugar.Errorw("获取注册验证码失败", "detail", "获取注册验证码失败", "error", err, "email", email, "token", token)
		return false
	}

	return storedToken == token
}

// 删除注册验证码
func DeleteRegisterVerificationToken(email, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("register_verification_token:%s", email)

	err := r_init.RedisClient.Del(ctx, key).Err()
	if err != nil {
		logs.Sugar.Errorw("删除注册验证码失败", "detail", "删除注册验证码失败", "error", err, "email", email, "token", token)
		return err
	}
	logs.Sugar.Infow("注册验证码删除成功", "detail", "删除注册验证码成功", "email", email, "token", token)
	return nil
}

// ============= 找回密码相关 ==============

// saveResetPwdToken 保存找回密码 reset_pwd_token:token -> email
func saveResetPwdToken(token, email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("reset_pwd_token:%s", token)
	return r_init.RedisClient.Set(ctx, key, email, VerificationTokenTTL).Err()
}

func hasResetPwdToken(token string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("reset_pwd_token:%s", token)

	_, err := r_init.RedisClient.Get(ctx, key).Result()
	return err == nil
}

// 通过 token 获取 email
func getEmailByResetPwdToken(token string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("reset_pwd_token:%s", token)
	email, err := r_init.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("token 不存在或已过期")
	}
	return email, err
}

// 删除 token（可选，非必须，因为会自动过期）
func deleteResetPwdToken(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("reset_pwd_token:%s", token)
	return r_init.RedisClient.Del(ctx, key).Err()
}
