package main

import (
	cf "backend/configs"
	"backend/internal/api"
	rep "backend/internal/repository"
	aiSvc "backend/internal/service/ai_s"
	emailSvc "backend/internal/service/email_s"
	kafka_s "backend/internal/service/kafka_s"
	userSvc "backend/internal/service/user_s"
	"backend/logs"
	"backend/pkg/binding"
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	mtr "backend/internal/metrics"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func init() {
	mtr.PrometheusRegister()     // 初m始化Prometheus
	binding.RegisterValidation() // 注册自定义验证器
	logs.InitZapSugarDefault()   // 初始化zap日志配置
	cf.SetEnvVariables()         // 设置环境变量
}

/*
启动该main后需启动 backend\cmd\email-worker\main.go 开启kafka消费者
*/
func main() {
	// 1. 初始化数据库和ai连接
	repo, err := rep.Init()
	if err != nil {
		logs.Sugar.Fatalf("Failed to initialize database: %v", err)
	}
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		logs.Sugar.Fatal("DASHSCOPE_API_KEY is required")
	}
	baseURL := os.Getenv("DASHSCOPE_BASE_URL")
	if baseURL == "" {
		logs.Sugar.Fatal("DASHSCOPE_BASE_URL is required")
	}
	aiClient := openai.NewClient(
		option.WithAPIKey(apiKey),
		// 以下是北京地域base_url，如果使用新加坡地域的模型，需要将base_url替换为：https://dashscope-intl.aliyuncs.com/compatible-mode/v1
		option.WithBaseURL(baseURL),
	)

	// 2. 组装服务
	userRepo := rep.NewUserRepository(repo.DB, repo.Redis)
	emailRepo := rep.NewEmailRepository(repo.DB, repo.Redis)
	aiRepo := rep.NewAIRepository(repo.Redis)
	kafkaProducer := kafka_s.NewDefaultKafkaProducerService()
	userService := userSvc.NewUserService(userRepo, emailRepo)
	emailService := emailSvc.NewEmailService(emailRepo, userRepo, kafkaProducer)
	aiService := aiSvc.NewAIService(&aiClient, aiRepo)

	enable_pprof := os.Getenv("ENABLE_PPROF")

	// 4. 注册路由
	r := api.SetupRouter(userService, emailService, aiService, enable_pprof)

	// 5. 创建 HTTP 服务实例
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 6. 创建 context 监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 应用 trace
	if os.Getenv("ENABLE_TRACE") == "true" {
		f, _ := os.Create("server.trace")
		defer f.Close()
		trace.Start(f)
		// go tool trace server.trace
		defer trace.Stop()
	}

	// 7. 启动 HTTP 服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Sugar.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	logs.Sugar.Info("Server started on :8080")

	// 8. 等待中断信号
	<-ctx.Done()

	logs.Sugar.Info("Shutting down server...")

	// 9. 创建一个超时 context 控制优雅关闭时间
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10. 停止 HTTP 服务
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logs.Sugar.Errorf("HTTP server Shutdown error: %v", err)
	} else {
		logs.Sugar.Info("HTTP server gracefully stopped")
	}

	// 11. 关闭 pg 数据库连接 *gorm.DB
	sqlDB, gormErr := repo.DB.DB()
	if gormErr == nil {
		if err := sqlDB.Close(); err != nil {
			logs.Sugar.Errorf("PostgreSQL GORM DB Close error: %v", err)
		} else {
			logs.Sugar.Info("PostgreSQL GORM DB closed")
		}
	} else {
		logs.Sugar.Error("Failed to get underlying SQL DB from GORM")
	}

	// 12. 关闭 Redis 连接
	if err := repo.Redis.Close(); err != nil {
		logs.Sugar.Errorf("Redis Close error: %v", err)
	} else {
		logs.Sugar.Info("Redis connection closed")
	}

	// 13. 关闭Kafka生产者
	kafkaProducer.Close()

	logs.Sugar.Info("Server exited")
}
