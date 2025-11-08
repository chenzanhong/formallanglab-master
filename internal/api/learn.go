package api

import (
	"backend/internal/metrics"
	"backend/logs"
	"time"

	"github.com/gin-gonic/gin"
)

func LearnGet(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("learn", "get", time.Since(start).Seconds())
	}()

	metrics.IncOperation("learn", "get", "success")
	logs.Sugar.Infow("学习资源获取成功")
}
