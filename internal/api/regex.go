/*
RegexValidate      	// 判断是否为有效的正则表达式
RegexRecognize		// 查看是否正则表达式是否匹配字符串
*/
package api

import (
	"net/http"
	"time"

	"backend/internal/domain/dto"
	"backend/internal/metrics"
	re "backend/internal/service/regex_s"

	"github.com/gin-gonic/gin"
)

// 使用dto包中的结构体替代本地定义

// 限定所使用的符号为后端实际功能实现所支持的：仅包含：[a-zA-Z0-9]、*、+、？、|、（）
func RegexValidate(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "validate", time.Since(start).Seconds())
	}()
	var req dto.RegexValidateRequest

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.RegexValidateResponse{
			Msg:   "Invalid request format: " + err.Error(),
			Valid: false,
		})
		metrics.IncOperation("regex", "validate", "failure: parameter parsing error")
		return
	}

	_, err := re.RegexValidate(req.Pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.RegexValidateResponse{
			Msg:   "invalid regular expression",
			Valid: false,
			Error: err.Error(),
		})
		metrics.IncOperation("regex", "validate", "failure: invalid regex")
		return
	}

	// 验证通过
	metrics.IncOperation("regex", "validate", "success")
	c.JSON(200, dto.RegexValidateResponse{
		Valid:   true,
		Msg:     "Regular expression is valid",
		Pattern: req.Pattern,
	})
}

func RegexRecognize(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "recognize", time.Since(start).Seconds())
	}()
	var req dto.RegexRecognizeRequest
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.RegexRecognizeResponse{
			Msg:     "Invalid request format: " + err.Error(),
			Matched: false,
		})
		metrics.IncOperation("regex", "recognize", "failure: parameter parsing error")
		return
	}

	_, err := re.RegexValidate(req.Pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.RegexRecognizeResponse{
			Msg:     "invalid regular expression",
			Matched: false,
			Error:   err.Error(),
		})
		metrics.IncOperation("regex", "recognize", "failure: invalid regex")
		return
	}

	_, err = re.RegexValidString(req.Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.RegexRecognizeResponse{
			Msg:     "invalid input string",
			Matched: false,
			Error:   err.Error(),
		})
		metrics.IncOperation("regex", "recognize", "failure: invalid string")
		return
	}

	ok, err := re.RegexRecognize(req.Pattern, req.Str)
	if !ok {
		c.JSON(http.StatusOK, dto.RegexRecognizeResponse{
			Msg:     "recognition failed: " + err.Error(),
			Matched: false,
		})
		metrics.IncOperation("regex", "recognize", "failure: recognition failed")
		return
	}
	metrics.IncOperation("regex", "recognize", "success")
	c.JSON(http.StatusOK, dto.RegexRecognizeResponse{
		Matched: true,
		Msg:     "pattern matches the string",
	})
}
