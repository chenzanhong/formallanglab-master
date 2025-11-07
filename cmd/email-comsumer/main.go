// backend/cmd/email-worker/main.go
/*
	邮件消费者
*/
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	elogs "backend/logs/email"

	eConfig "backend/configs/email"
	kafka_s "backend/internal/service/kafka_s"

	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
)

func main() {
	elogs.InitZapSugarDefault()
	eConfig.SetEmailEnvVariables()
	
	brokers := strings.Split(strings.TrimSpace(os.Getenv("KAFKA_BROKERS")), ",") // 从配置读取
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: "email-service-group",
	})
	fmt.Println(os.Getenv("KAFKA_TOPIC"))

	// 优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	elogs.EmailSugar.Info("Email worker started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			elogs.EmailSugar.Info("Shutting down email worker...")
			reader.Close()
			return
		default:
			// fmt.Println(1)
			msg, err := reader.ReadMessage(ctx) // 阻塞读取
			// fmt.Println(2)
			if err != nil {
				elogs.EmailSugar.Infof("Error reading message: %v", err)
				continue
			}

			var event kafka_s.KafkaEmailEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				elogs.EmailSugar.Infof("Failed to unmarshal email event: %v", err)
				continue
			}

			// 重试机制（可选：指数退避）
			maxRetries := 3
			for attempt := 1; attempt <= maxRetries; attempt++ {
				err := sendEmailSync(event.To, event.Subject, event.ContentType, event.Body)
				if err == nil {
					elogs.EmailSugar.Infow("Email sent successfully", "to", event.To, "subject", event.Subject)
					break
				}

				elogs.EmailSugar.Warnw("Failed to send email, retrying...", "attempt", attempt, "error", err)
				if attempt < maxRetries {
					time.Sleep(time.Duration(1<<uint(attempt-1)) * time.Second)
				} else {
					elogs.EmailSugar.Errorw("Email sending failed after retries", "to", event.To, "error", err)
					// 可选：发送到死信队列（DLQ），不过验证码一分钟后过期，所以这里不放死信队列，前端用户手动重试
				}
			}
		}
	}
}

// sendEmailSync 是实际的同步发送逻辑
func sendEmailSync(email, subject, contextType, body string) error {
	myEmail := os.Getenv("EMAIL_NAME")
	myPassword := os.Getenv("EMAIL_PASSWORD")
	smtpServerHost := os.Getenv("SMTP_SERVER_HOST")
	smtpServerPortStr := os.Getenv("SMTP_SERVER_PORT")

	if myEmail == "" || myPassword == "" || smtpServerHost == "" || smtpServerPortStr == "" {
		return errors.New("环境变量未正确设置")
	}

	smtpServerPort, err := strconv.Atoi(smtpServerPortStr)
	if err != nil {
		return fmt.Errorf("SMTP 端口解析失败: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", myEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody(contextType, body)

	d := gomail.NewDialer(smtpServerHost, smtpServerPort, myEmail, myPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 生产环境建议使用有效证书

	if err := d.DialAndSend(m); err != nil {
		if strings.Contains(err.Error(), "535") {
			return errors.New("SMTP 身份验证失败，请检查邮箱账号或授权码")
		} else if strings.Contains(err.Error(), "connection refused") {
			return errors.New("无法连接 SMTP 服务器，请检查网络或服务器地址")
		}
		return fmt.Errorf("邮件发送失败: %w", err)
	}

	elogs.EmailSugar.Infow("邮件发送成功", "to", email, "subject", subject)
	return nil
}
