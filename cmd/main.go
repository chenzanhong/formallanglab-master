package main

import (
	"backend/configs"
	cf "backend/configs"
	"backend/internal/api"
	"backend/internal/middleware"
	rep "backend/internal/repository"
	kafka_s "backend/internal/service/kafka_s"
	learnSvc "backend/internal/service/learn_s"
	storeSvc "backend/internal/service/store_s"
	"backend/pkg/binding"
	"backend/pkg/oss"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	mtr "backend/internal/metrics"

	"github.com/chenzanhong/goutil/jwtx"
	"github.com/chenzanhong/zlog"
)

func init() {
	mtr.PrometheusRegister()     // 初始化 Prometheus
	binding.RegisterValidation() // 注册自定义验证器
	// logs.InitLoggerFromEnv()
}

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败：%v", err.Error())
	}
	cf.SyncConfigToEnv(*config) // 环境变量设置
	jwtx.InitWithHS256(config.JWT.Key, &middleware.Claims{}, jwtx.WithAutoInject(true))

	// 日志
	zlog.InitLogger(config.Log)

	// 1. 初始化数据库
	repo, err := rep.Init()
	if err != nil {
		zlog.Fatalf("Failed to initialize database: %v", err)
	}
	// apiKey := os.Getenv("DASHSCOPE_API_KEY")
	// if apiKey == "" {
	// 	zlog.Fatal("DASHSCOPE_API_KEY is required")
	// }
	// baseURL := os.Getenv("DASHSCOPE_BASE_URL")
	// if baseURL == "" {
	// 	zlog.Fatal("DASHSCOPE_BASE_URL is required")
	// }
	// aiClient := openai.NewClient(
	// 	option.WithAPIKey(apiKey),
	// 	// 以下是北京地域base_url，如果使用新加坡地域的模型，需要将base_url替换为：https://dashscope-intl.aliyuncs.com/compatible-mode/v1
	// 	option.WithBaseURL(baseURL),
	// )

	// 2. 初始化OSS客户端
	ossClient, err := oss.NewAliyunOSSClient()
	if err != nil {
		zlog.Errorf("初始化阿里云OSS客户端失败: %v", err)
	}

	// 3. 组装服务
	// userRepo := rep.NewUserRepository(repo.DB, repo.Redis)
	// emailRepo := rep.NewEmailRepository(repo.DB, repo.Redis)
	storeRepo := rep.NewStoreRepository(repo.DB)
	// aiRepo := rep.NewAIRepository(repo.Redis)
	learnRepo := rep.NewLearnRepository(repo.DB)
	kafkaProducer := kafka_s.NewDefaultKafkaProducerService()
	// userService := userSvc.NewUserService(userRepo, emailRepo)
	// emailService := emailSvc.NewEmailService(emailRepo, userRepo, kafkaProducer)
	// aiService := aiSvc.NewAIService(&aiClient, aiRepo)
	storeService := storeSvc.NewStoreService(storeRepo)
	learnService := learnSvc.NewLearnService(learnRepo, ossClient)

	// 4. 初始化处理器
	// userHandler := api.NewUserHandler(userService)
	// emailHandler := api.NewEmailHandler(emailService)
	storeHandler := api.NewStoreHandler(storeService)
	learnHandler := api.NewLearnHandler(learnService)

	// 5. 注册路由
	// r := api.SetupRouter(userService, emailService, aiService, learnHandler, enable_pprof)
	r := api.SetupRouter(storeHandler, learnHandler)

	// 5. 创建 HTTP 服务实例
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: r,
	}

	// 6. 创建 context 监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 应用 trace
	go func() {
		if v, ok := os.LookupEnv("ENABLE_TRACE"); ok && v == "true" {
			f, _ := os.Create("server.trace")
			defer f.Close()
			trace.Start(f)
			// go tool trace server.trace
			defer trace.Stop()
		}
	}()

	// 启动pprof http服务
	go func() {
		if os.Getenv("PPROF_PORT") != "0" {
			zlog.Infow("Starting pprof on localhost:", config.Server.PprofPort)
			http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Server.PprofPort), nil)
		}
	}()

	// 7. 启动 HTTP 服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	zlog.Info("Server started on :8080")

	// 8. 等待中断信号
	<-ctx.Done()

	zlog.Info("Shutting down server...")

	// 9. 创建一个超时 context 控制优雅关闭时间
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10. 停止 HTTP 服务
	if err := srv.Shutdown(shutdownCtx); err != nil {
		zlog.Errorf("HTTP server Shutdown error: %v", err)
	} else {
		zlog.Info("HTTP server gracefully stopped")
	}

	// 11. 关闭 pg 数据库连接 *gorm.DB
	sqlDB, gormErr := repo.DB.DB()
	if gormErr == nil {
		if err := sqlDB.Close(); err != nil {
			zlog.Errorf("PostgreSQL GORM DB Close error: %v", err)
		} else {
			zlog.Info("PostgreSQL GORM DB closed")
		}
	} else {
		zlog.Error("Failed to get underlying SQL DB from GORM")
	}

	// 12. 关闭 Redis 连接
	if err := repo.Redis.Close(); err != nil {
		zlog.Errorf("Redis Close error: %v", err)
	} else {
		zlog.Info("Redis connection closed")
	}

	// 13. 关闭Kafka生产者
	kafkaProducer.Close()

	zlog.Info("Server exited")
}
