/*
RegexValidate      			// 判断是否为有效的正则表达式
RegexRecognize				// 查看是否正则表达式是否匹配字符串
RegexEquivalenceCheck 		// 检查两个正则表达式是否等价
RegexGenerateExampleString 	// 生成正则表达式可匹配的字符串示例
*/
package handler

import (
	"net/http"
	"time"

	"github.com/chenzanhong/zlog"
	"github.com/gin-gonic/gin"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/middleware/metrics"
	re "github.com/chenzanhong/formallanglab-master/internal/service/regex_s"
)

func RegexValidate(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "validate", time.Since(start).Seconds())
	}()

	var req dto.RegexValidateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("regex", "validate", "failure: parameter parsing error")
		zlog.Warnw("正则表达式验证失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.RegexValidateResponse{
			Msg:    "无效的请求格式：" + err.Error(),
			Valid:  false,
			Result: false,
		})

		return
	}

	if err := re.RegexValidate(req.Pattern); err != nil {
		metrics.IncOperation("regex", "validate", "failure: invalid regex")
		zlog.Warnw("正则表达式验证失败", "detail", "正则表达式格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexValidateResponse{
			Msg:    "无效的正则表达式",
			Valid:  false,
			Result: false,
		})

		return
	}

	metrics.IncOperation("regex", "validate", "success")
	zlog.Infow("正则表达式验证成功")
	c.JSON(200, dto.RegexValidateResponse{
		Valid:   true,
		Msg:     "正则表达式有效",
		Pattern: req.Pattern,
	})
}

func RegexRecognize(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "recognize", time.Since(start).Seconds())
	}()

	var req dto.RegexRecognizeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("regex", "recognize", "failure: parameter parsing error")
		zlog.Warnw("正则表达式匹配失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(400, dto.RegexRecognizeResponse{
			Msg:     "无效的请求格式：" + err.Error(),
			Matched: false,
			Result:  false,
		})

		return
	}

	if err := re.RegexValidate(req.Pattern); err != nil {
		metrics.IncOperation("regex", "recognize", "failure: invalid regex")
		zlog.Warnw("正则表达式匹配失败", "detail", "正则表达式格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexRecognizeResponse{
			Msg:     "无效的正则表达式",
			Matched: false,
			Result:  false,
		})

		return
	}

	_, err := re.RegexValidString(req.Str)
	if err != nil {
		metrics.IncOperation("regex", "recognize", "failure: invalid string")
		zlog.Warnw("正则表达式匹配失败", "detail", "输入字符串无效")
		c.JSON(http.StatusBadRequest, dto.RegexRecognizeResponse{
			Msg:     "无效的输入字符串",
			Matched: false,
			Result:  false,
		})

		return
	}

	ok, err := re.RegexRecognize(req.Pattern, req.Str)
	if err != nil {
		// 内部计算错误
		metrics.IncOperation("regex", "recognize", "failure: recognition error")
		zlog.Warnw("正则表达式匹配失败", "detail", "匹配过程中发生错误")
		c.JSON(http.StatusInternalServerError, dto.RegexRecognizeResponse{
			Msg:     "匹配过程中发生错误：" + err.Error(),
			Matched: false,
			Result:  false,
		})

		return
	}

	if !ok {
		// 匹配结果为否（正常的逻辑结果）
		metrics.IncOperation("regex", "recognize", "success: not matched")
		zlog.Infow("正则表达式匹配成功", "detail", "正则表达式无法匹配给定的字符串")
		c.JSON(http.StatusOK, dto.RegexRecognizeResponse{
			Msg:     "字符串不匹配",
			Matched: false,
			Result:  true,
		})

		return
	}

	metrics.IncOperation("regex", "recognize", "success")
	zlog.Infow("正则表达式匹配成功")
	c.JSON(http.StatusOK, dto.RegexRecognizeResponse{
		Matched: true,
		Msg:     "模式匹配字符串",
		Result:  true,
	})
}

func RegexEquivalenceCheck(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "equivalence_check", time.Since(start).Seconds())
	}()

	var req dto.RegexEquivalenceCheckRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("regex", "equivalence_check", "failure: parameter parsing error")
		zlog.Warnw("正则表达式等价检查失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(400, dto.RegexEquivalenceCheckResponse{
			Msg:    "无效的请求格式：" + err.Error(),
			Result: false,
		})

		return
	}

	if err := re.RegexValidate(req.Pattern1); err != nil {
		metrics.IncOperation("regex", "equivalence_check", "failure: invalid regex1")
		zlog.Warnw("正则表达式等价检查失败", "detail", "正则表达式 1 格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexEquivalenceCheckResponse{
			Msg:    "无效的正则表达式",
			Result: false,
		})

		return
	}

	if err := re.RegexValidate(req.Pattern2); err != nil {
		metrics.IncOperation("regex", "equivalence_check", "failure: invalid regex2")
		zlog.Warnw("正则表达式等价检查失败", "detail", "正则表达式 2 格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexEquivalenceCheckResponse{
			Msg:    "无效的正则表达式",
			Result: false,
		})

		return
	}

	ok, err := re.RegexEquivalenceCheck(req.Pattern1, req.Pattern2)
	if err != nil {
		// 内部计算错误
		metrics.IncOperation("regex", "equivalence_check", "failure: equivalence check error")
		zlog.Warnw("正则表达式等价检查失败", "detail", "等价性判断过程中发生错误")
		c.JSON(http.StatusInternalServerError, dto.RegexEquivalenceCheckResponse{
			Msg:          "等价性判断出错：" + err.Error(),
			Result:       false,
			IsEquivalent: false,
		})

		return
	}

	if !ok {
		// 检查结果为不等价（这是成功的检查结果）
		metrics.IncOperation("regex", "equivalence_check", "success: not equivalent")
		zlog.Infow("正则表达式等价检查成功", "detail", "正则表达式 1 和 2 不等价")
		c.JSON(http.StatusOK, dto.RegexEquivalenceCheckResponse{
			Msg:          "这两个正则表达式不是等价的",
			Result:       true,
			IsEquivalent: false,
		})

		return
	}

	metrics.IncOperation("regex", "equivalence_check", "success")
	zlog.Infow("正则表达式等价检查成功")
	c.JSON(http.StatusOK, dto.RegexEquivalenceCheckResponse{
		Msg:          "等价性检查通过",
		Result:       true,
		IsEquivalent: true,
	})
}

func RegexGenerateExampleString(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "generate_example", time.Since(start).Seconds())
	}()

	var req dto.RegexGenerateExampleStringRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("regex", "generate_example", "failure: parameter parsing error")
		zlog.Warnw("正则表达式示例生成失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(400, dto.RegexGenerateExampleStringResponse{
			Msg:    "无效的请求格式：" + err.Error(),
			Result: false,
		})

		return
	}

	if err := re.RegexValidate(req.Pattern); err != nil {
		metrics.IncOperation("regex", "generate_example", "failure: invalid regex")
		zlog.Warnw("正则表达式示例生成失败", "detail", "正则表达式格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexGenerateExampleStringResponse{
			Msg:    "无效的正则表达式",
			Result: false,
		})

		return
	}

	accept, reject := re.RegexGenerateExampleString(req.Pattern)

	metrics.IncOperation("regex", "generate_example", "success")
	zlog.Infow("正则表达式示例生成成功")
	c.JSON(http.StatusOK, dto.RegexGenerateExampleStringResponse{
		Msg:            "示例字符串生成成功",
		AcceptExamples: accept,
		RejectExamples: reject,
		Result:         true,
	})
}

func RegexSimplify(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("regex", "simplify", time.Since(start).Seconds())
	}()

	var req dto.RegexSimplifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("regex", "simplify", "failure: parameter parsing error")
		zlog.Warnw("正则表达式化简失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(400, dto.RegexSimplifyResponse{
			Msg:    "无效的请求格式：" + err.Error(),
			Result: false,
		})

		return
	}

	if err := re.RegexValidate(req.Pattern); err != nil {
		metrics.IncOperation("regex", "simplify", "failure: invalid regex")
		zlog.Warnw("正则表达式化简失败", "detail", "正则表达式格式无效")
		c.JSON(http.StatusBadRequest, dto.RegexSimplifyResponse{
			Msg:    "无效的正则表达式",
			Result: false,
		})

		return
	}

	simplified, isEmptyLanguage := re.SimplifyRegex(string(req.Pattern))

	metrics.IncOperation("regex", "simplify", "success")
	zlog.Infow("正则表达式化简成功")
	c.JSON(http.StatusOK, dto.RegexSimplifyResponse{
		Msg:             "正则表达式化简成功",
		Result:          true,
		Original:        string(req.Pattern),
		Simplified:      simplified,
		IsEmptyLanguage: isEmptyLanguage,
	})
}
