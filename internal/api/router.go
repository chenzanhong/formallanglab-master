package api

import (
	"backend/configs"
	"backend/logs"
	"backend/pkg/middleware"
	"strconv"

	"os"

	mtr "backend/internal/metrics"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.Use(mtr.HTTPMiddleware(), middleware.CORSMiddleware())

	setupPublicRoutes(router) // 注册公开路由
	setupAuthRoutes(router)   // 注册需要认证的路由

	return router
}

func SetEnvVariables() {
	config, err := configs.LoadConfig()
	if err != nil {
		logs.Sugar.Fatalf("加载配置失败：%v", err.Error())
	}
	os.Setenv("DB_USER", config.PG.User)
	os.Setenv("DB_PASSWORD", config.PG.Password)
	os.Setenv("DB_HOST", config.PG.Host)
	os.Setenv("DB_PORT", config.PG.Port)
	os.Setenv("DB_NAME", config.PG.Name)
	os.Setenv("REDIS_HOST", config.Redis.Host)
	os.Setenv("REDIS_PORT", config.Redis.Port)
	os.Setenv("REDIS_PASSWORD", config.Redis.Password)
	os.Setenv("REDIS_DB", strconv.Itoa(config.Redis.DB))
	os.Setenv("EMAIL_NAME", config.Email.Name)
	os.Setenv("EMAIL_PASSWORD", config.Email.Password)
	os.Setenv("SMTP_SERVER_HOST", config.SMTPServer.Host)
	os.Setenv("SMTP_SERVER_PORT", config.SMTPServer.Port)
}

func setupPublicRoutes(router *gin.Engine) {
	router.GET("/gdesign/metrics", mtr.MetricsHandler())                 // prometheus.yml中加上 metrics_path: /gdesign/metrics
	router.POST("/gdesign/register", Register)                           // 注册
	router.POST("/gdesign/login", Login)                                 // 登录
	router.POST("/gdesign/send_verification_code", SendVerificationCode) // 发送验证码（注册用）
	router.POST("/gdesign/req_resetpassword", RequestResetPassword)      // 请求重置密码
	router.POST("/gdesign/resetpassword", ResetPassword)                 // 重置密码
}

func setupAuthRoutes(router *gin.Engine) {
	// 使用 JWT 中间件保护这些路由
	r := router.Group("/gdesign", middleware.JWTAuthMiddleware())

	// 文法相关接口
	grammar := r.Group("/grammar")
	{
		grammar.POST("/validate", GrammarValidate)            // 文法校验——是否有效
		grammar.POST("/type", GrammarTypeDetermine)           // 判断所给文法的类型
		grammar.POST("/simplify", GrammarSimplify)            // 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
		grammar.POST("/ambiguity", GrammarAmbiguityCheck)     // 正则文法的二义性判断
		grammar.POST("/equivalence", GrammarEquivalenceCheck) // 判断所给的两个正则文法是否等价
		grammar.POST("/recognize", GrammarStringRecognize)    // 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
	}

	// 正则表达式相关接口
	regEx := r.Group("/regex")
	{
		regEx.POST("/validate", RegexValidate) // 判断是否为有效的正则表达式
	}

	// 自动机相关接口
	fsm := r.Group("/fsm")
	{
		fsm.POST("/validate", FSMValidate)               // 是否有效
		fsm.POST("/recognize", FSMStringRecognize)       // 字符串识别
		fsm.POST("/cleanup", FSMCleanup)                 // 去无效符号、不可达符号
		fsm.POST("/minimize", DFAMinimize)               // DFA 最小化
		fsm.POST("/determinization", NFADeterminization) // NFA 转 DFA，NFA确定化
	}

	// 文法、自动机间的转换
	convert := r.Group("/convert")
	{
		convert.POST("/grammartonfa", GrammarToNFA) // 文法转成 NFA
		convert.POST("/fatogrammar", FAToGrammar)   // FA 转成文法
		convert.POST("/regextonfa", RegexToNFA)      // 正则表达式转为NFA
		convert.POST("/fatoregex", FAToRegex)        // FA转为正则表达式
	}

	// 知识学习
	learn := r.Group("/learn")
	{
		learn.GET("/", LearnGet) // 获取学习资料或信息
	}

	// AI 相关接口
	ai := r.Group("/ai")
	{
		ai.GET("/", AIGet) // 获取AI相关信息或执行特定操作
	}
}
