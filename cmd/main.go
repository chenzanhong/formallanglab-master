package main

import (
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

	"github.com/chenzanhong/goutil/jwtx"
	"github.com/chenzanhong/zlog"
	"github.com/gin-gonic/gin"

	"github.com/chenzanhong/formallanglab-master/configs"
	cf "github.com/chenzanhong/formallanglab-master/configs"
	"github.com/chenzanhong/formallanglab-master/internal/handler"
	mtr "github.com/chenzanhong/formallanglab-master/internal/metrics"
	"github.com/chenzanhong/formallanglab-master/internal/middleware"
	rep "github.com/chenzanhong/formallanglab-master/internal/repository"
	learnSvc "github.com/chenzanhong/formallanglab-master/internal/service/learn_s"
	storeSvc "github.com/chenzanhong/formallanglab-master/internal/service/store_s"
	"github.com/chenzanhong/formallanglab-master/pkg/binding"
	"github.com/chenzanhong/formallanglab-master/pkg/oss"
)

func init() {
	mtr.PrometheusRegister()     // 初始化 Prometheus
	binding.RegisterValidation() // 注册自定义验证器
}

func main() {
	// 1. 加载配置
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败：%v", err.Error())
	}
	// 2. 设置环境变量
	cf.SyncConfigToEnv(*config)
	// 3. 初始化JWT
	jwtx.InitWithHS256(config.JWT.Key, &middleware.Claims{}, jwtx.WithAutoInject(true))

	// 4. 初始化日志
	zlog.InitLogger(config.Log)

	// 5. 初始化数据库
	repo, err := rep.Init()
	if err != nil {
		zlog.Fatalf("Failed to initialize database: %v", err)
	}

	// 6. 初始化OSS客户端
	ossClient, err := oss.NewAliyunOSSClient()
	if err != nil {
		zlog.Errorf("初始化阿里云OSS客户端失败: %v", err)
	}

	// 7. 组装服务
	storeRepo := rep.NewStoreRepository(repo.DB)
	learnRepo := rep.NewLearnRepository(repo.DB)
	storeService := storeSvc.NewStoreService(storeRepo)
	learnService := learnSvc.NewLearnService(learnRepo, ossClient)

	// 8. 初始化处理器
	storeHandler := handler.NewStoreHandler(storeService)
	learnHandler := handler.NewLearnHandler(learnService)

	// 9. 注册路由
	r := handler.SetupRouter(storeHandler, learnHandler)

	// 10. 创建 HTTP 服务实例
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: r,
	}

	// 11. 创建 context 监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 应用 trace（通过 ENABLE_TRACE 环境变量控制）
	go func() {
		if v, ok := os.LookupEnv("ENABLE_TRACE"); ok && v == "true" {
			f, _ := os.Create("server.trace")
			defer f.Close()
			trace.Start(f)
			// go tool trace server.trace
			defer trace.Stop()
		}
	}()

	// 12. 启动 pprof http 服务（通过 PPROF_PORT 环境变量控制，默认为 6060）
	go func() {
		if pprofPort := os.Getenv("PPROF_PORT"); pprofPort != "0" && pprofPort != "" {
			zlog.Info("Starting pprof on :"+pprofPort)
			if err := http.ListenAndServe(":"+pprofPort, nil); err != nil {
				zlog.Errorf("pprof server error: %v", err)
			}
		}
	}()

	// 13. 启动独立的 metrics 服务（通过 METRICS_PORT 环境变量控制）
	go func() {
		if metricsPort := os.Getenv("METRICS_PORT"); metricsPort != "0" && metricsPort != "" {
			r := gin.New()
			r.Use(gin.Recovery())
			zlog.Infow("Starting metrics on localhost:" + metricsPort)
			r.GET("/gdesign/master/metrics", mtr.MetricsHandler())
			r.Run(fmt.Sprintf(":%s", metricsPort))
		}
	}()

	// 14. 启动主 HTTP 服务
	// 启动前自动同步 OSS 文件
	go func() {
		zlog.Info("开始同步 OSS 文件...")
		ctx := context.Background()
		_, err := learnService.SyncOSSFiles(ctx)
		if err != nil {
			zlog.Errorf("同步 OSS 文件失败: %v", err)
		} else {
			zlog.Info("OSS 文件同步完成")
		}
	}()

	// 启动主 HTTP 服务（通过 SERVER_PORT 环境变量控制，默认为 8081）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
		zlog.Info("server exited")
	}()

	zlog.Info("Server started on :8081")

	// 等待中断信号（SIGINT, SIGTERM）
	<-ctx.Done()

	zlog.Info("Shutting down server...")

	// 15. 优雅关闭 HTTP 服务
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zlog.Errorf("HTTP server Shutdown error: %v", err)
	} else {
		zlog.Info("HTTP server gracefully stopped")
	}

	// 16. 关闭 PostgreSQL 数据库连接
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

	zlog.Info("Server exited")
}
