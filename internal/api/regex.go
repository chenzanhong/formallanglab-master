/*
RegexValidate      	// 判断是否为有效的正则表达式
RegexRecognize		// 查看是否正则表达式是否匹配字符串
*/
package api

import (
	"net/http"

	re "backend/internal/service/regex_s"

	"github.com/gin-gonic/gin"
)

// RegexValidateRequest 正则表达式验证请求结构
type RegexValidateRequest struct {
	Pattern string `json:"pattern" binding:"required"`
}

type RegexRecognizeRequest struct {
	Pattern string `json:"pattern" binding:"required"`
	Str     string `json:"str" binding:"required"`
}

// 限定所使用的符号为后端实际功能实现所支持的：仅包含：[a-zA-Z0-9]、*、+、？、|、（）
func RegexValidate(c *gin.Context) {
	var req RegexValidateRequest

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request format: " + err.Error()})
		return
	}

	_, err := re.RegexValidate(req.Pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid regular expression", "valid": false, "error": err})
	}

	// 验证通过
	c.JSON(200, gin.H{
		"valid":   true,
		"msg":     "Regular expression is valid",
		"pattern": req.Pattern,
	})
}

func RegexRecognize(c *gin.Context) {
	var req RegexRecognizeRequest
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request format: " + err.Error()})
		return
	}

	_, err := re.RegexValidate(req.Pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid regular expression", "matched": false, "error": err})
	}

	_, err = re.RegexValidString(req.Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid input string", "matched": false, "error": err.Error()})
		return
	}

	ok, err := re.RegexRecognize(req.Pattern, req.Str)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"msg": "recognition failed: " + err.Error(), "matched": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"matched": true, "msg": "pattern matches the string"})
}
