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
	"backend/logs"

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
		metrics.IncOperation("regex", "validate", "failure: parameter parsing error")
		logs.Sugar.Warnw("正则表达式验证失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(400, dto.RegexValidateResponse{
			Msg:    "Invalid request format: " + err.Error(),
			Valid:  false,
			Result: false,
		})
		return
	}

	if err := re.RegexValidate(req.Pattern); err != nil {
		metrics.IncOperation("regex", "validate", "failure: invalid regex")
		logs.Sugar.Warnw("正则表达式验证失败", "detail", "正则表达式格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexValidateResponse{
			Msg:    "invalid regular expression",
			Valid:  false,
			Error:  err.Error(),
			Result: false,
		})
		return
	}

	// 验证通过
	metrics.IncOperation("regex", "validate", "success")
	logs.Sugar.Infow("正则表达式验证成功")
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
		metrics.IncOperation("regex", "recognize", "failure: parameter parsing error")
		logs.Sugar.Warnw("正则表达式匹配失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(400, dto.RegexRecognizeResponse{
			Msg:     "Invalid request format: " + err.Error(),
			Matched: false,
			Result:  false,
		})
		return
	}

	if err := re.RegexValidate(req.Pattern); err != nil {
		metrics.IncOperation("regex", "recognize", "failure: invalid regex")
		logs.Sugar.Warnw("正则表达式匹配失败", "detail", "正则表达式格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexRecognizeResponse{
			Msg:     "invalid regular expression",
			Matched: false,
			Error:   err.Error(),
			Result:  false,
		})
		return
	}

	_, err := re.RegexValidString(req.Str)
	if err != nil {
		metrics.IncOperation("regex", "recognize", "failure: invalid string")
		logs.Sugar.Warnw("正则表达式匹配失败", "detail", "输入字符串无效")
		c.JSON(http.StatusBadRequest, dto.RegexRecognizeResponse{
			Msg:     "invalid input string",
			Matched: false,
			Error:   err.Error(),
			Result:  false,
		})
		return
	}

	ok, err := re.RegexRecognize(req.Pattern, req.Str)
	if !ok {
		metrics.IncOperation("regex", "recognize", "failure: recognition failed")
		logs.Sugar.Warnw("正则表达式匹配失败", "detail", "正则表达式无法匹配给定的字符串")
		c.JSON(http.StatusOK, dto.RegexRecognizeResponse{
			Msg:     "recognition failed: " + err.Error(),
			Matched: false,
			Result:  false,
		})
		return
	}
	metrics.IncOperation("regex", "recognize", "success")
	logs.Sugar.Infow("正则表达式匹配成功")
	c.JSON(http.StatusOK, dto.RegexRecognizeResponse{
		Matched: true,
		Msg:     "pattern matches the string",
		Result:  true,
	})
}
