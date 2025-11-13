package main

import (
	cf "backend/configs"
	"backend/internal/api"
	"backend/internal/middleware"
	rep "backend/internal/repository"
	aiSvc "backend/internal/service/ai_s"
	"backend/logs"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mtr "backend/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func init() {
	mtr.PrometheusRegister() // 初始化Prometheus
	cf.SetEnvVariables()     // 初始化配置以及环境变量设置
	logs.InitLoggerFromEnv()
}

// 独立的WebSocket服务器
func main() {
	// 1. 初始化数据库和Redis连接
	repo, err := rep.Init()
	if err != nil {
		logs.Sugar.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. 初始化OpenAI客户端
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
		option.WithBaseURL(baseURL),
	)

	// 3. 初始化AI服务
	aiRepo := rep.NewAIRepository(repo.Redis)
	aiService := aiSvc.NewAIService(&aiClient, aiRepo)

	// 4. 创建AI处理器
	qaCache := aiSvc.NewQACache()
	// 加载预置高频问题缓存
	if err := qaCache.LoadCache("./knowledge/qa/qa.json"); err != nil {
		logs.Sugar.Fatalf("Failed to load QACache", "detail", err.Error())
	}
	aiHandler := api.NewAIHandler(aiService, qaCache)

	// 5. 创建Gin引擎
	router := gin.New()
	// 只添加必要的中间件
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware()) // 需要导入middleware包

	// 6. 只注册WebSocket路由
	router.GET("/gdesign/ai/ws", aiHandler.AIChatWS)

	// 7. 创建HTTP服务实例
	wsPort := os.Getenv("WS_PORT")
	if wsPort == "" {
		wsPort = "8081" // 默认端口
	}
	srv := &http.Server{
		Addr:    ":" + wsPort,
		Handler: router,
	}

	// 8. 创建context监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 9. 启动HTTP服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Sugar.Fatalf("WebSocket server ListenAndServe error: %v", err)
		}
	}()

	logs.Sugar.Infof("WebSocket server started on :%s", wsPort)

	// 10. 等待中断信号
	<-ctx.Done()

	logs.Sugar.Info("Shutting down WebSocket server...")

	// 11. 创建一个超时context控制优雅关闭时间
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 12. 停止HTTP服务
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logs.Sugar.Errorf("WebSocket server Shutdown error: %v", err)
	} else {
		logs.Sugar.Info("WebSocket server gracefully stopped")
	}

	// 13. 关闭数据库连接
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

	// 14. 关闭Redis连接
	if err := repo.Redis.Close(); err != nil {
		logs.Sugar.Errorf("Redis Close error: %v", err)
	} else {
		logs.Sugar.Info("Redis connection closed")
	}

	logs.Sugar.Info("WebSocket server exited")
}
