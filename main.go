package main

import (
	"backend/internal/api"
	rep "backend/internal/repository"
	"backend/logs"
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	mtr "backend/internal/metrics"
)

func init() {
	mtr.PrometheusRegister()
}

func main() {
	// 初始化zap日志配置
	logs.InitZapSugarDefault()
	fmt.Println("Init zap")
	// 设置环境变量
	api.SetEnvVariables()

	// 初始化数据库
	if err := rep.InitDB(); err != nil {
		logs.Sugar.Fatalf("Failed to initialize database: %v", err)
	}

	// 注册路由
	r := api.SetupRouter()
	// r.Run(":8080")

	// 创建 HTTP 服务实例
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 创建 context 监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 启动 HTTP 服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Sugar.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	logs.Sugar.Info("Server started on :8080")

	// 等待中断信号
	<-ctx.Done()

	logs.Sugar.Info("Shutting down server...")

	// 创建一个超时 context 控制优雅关闭时间
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. 停止 HTTP 服务
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logs.Sugar.Errorf("HTTP server Shutdown error: %v", err)
	} else {
		logs.Sugar.Info("HTTP server gracefully stopped")
	}

	// 2. 关闭 pg 数据库连接 *gorm.DB
	sqlDB, gormErr := rep.DB.DB()
	if gormErr == nil {
		if err := sqlDB.Close(); err != nil {
			logs.Sugar.Errorf("PostgreSQL GORM DB Close error: %v", err)
		} else {
			logs.Sugar.Info("PostgreSQL GORM DB closed")
		}
	} else {
		logs.Sugar.Error("Failed to get underlying SQL DB from GORM")
	}

	logs.Sugar.Info("Server exited")
}
