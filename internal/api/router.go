package api

import (
	"backend/internal/middleware"
	aiSvc "backend/internal/service/ai_s"
	emailSvc "backend/internal/service/email_s"
	userSvc "backend/internal/service/user_s"
	"net/http"
	"net/http/pprof"

	mtr "backend/internal/metrics"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userService userSvc.UserService, emailService emailSvc.EmailService, aiService aiSvc.AIService, enable_pprof string) *gin.Engine {

	userHandler := NewUserHandler(userService)
	emailHandler := NewEmailHandler(emailService)
	aiHandler := NewAIHandler(aiService)

	router := gin.New()
	// 1. 恢复中间件 - 最先使用，捕获所有panic
	router.Use(gin.Recovery())
	router.GET("/gdesign/ai/ws", aiHandler.AIChatWS) // WebSocket聊天接口，不经过JWT中间件

	// 2. 请求ID中间件 - 尽早设置，让后续中间件都能使用
	router.Use(middleware.RequestID())
	// 3. 全局速率限制 - 在处理请求初期进行限制，避免资源浪费，
	// 但为了与UserRateLimitMiddleware不重复，只在后面的公共路由组添加
	// router.Use(middleware.GlobalRateLimitMiddleware())
	// 4. CORS中间件 - 尽早处理跨域请求，避免不必要的后续处理
	router.Use(middleware.CORSMiddleware())
	// 5. 日志中间件 - 在业务逻辑前记录请求，在业务逻辑后记录响应。日志中间件，但是感觉有点笨重，暂时不使用
	// router.Use(middleware.Logging(middleware.DefaultLoggingConfig))
	// 6. 指标收集 - 收集所有处理过程的指标
	router.Use(mtr.HTTPMiddleware())
	if enable_pprof == "true" {
		setupPprof(router)
	}

	setupPublicRoutes(router, userHandler, emailHandler) // 注册公开路由

	setupAuthRoutes(router, userHandler, aiHandler) // 注册需要认证的路由

	return router
}

func setupPprof(router *gin.Engine) {
	pprofGroup := router.Group("/debug/pprof")
	pprofGroup.Use(func(c *gin.Context) {
		if c.ClientIP() != "127.0.0.1" && c.ClientIP() != "::1" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func setupPublicRoutes(router *gin.Engine, userHandler *UserHandler, emailHandler *EmailHandler) {
	router.GET("/gdesign/metrics", mtr.MetricsHandler())                                                                      // 不需要限速                                                                        // prometheus.yml中加上 metrics_path: /gdesign/metrics
	router.POST("/gdesign/register", middleware.GlobalRateLimitMiddleware(), userHandler.Register)                            // 注册
	router.POST("/gdesign/login", middleware.GlobalRateLimitMiddleware(), userHandler.Login)                                  // 登录
	router.POST("/gdesign/register/code", middleware.GlobalRateLimitMiddleware(), emailHandler.SendRegisterVerificationCode)  // 发送验证码（注册用）
	router.POST("/gdesign/reset-pwd/code", middleware.GlobalRateLimitMiddleware(), emailHandler.SendResetPwdVerificationCode) // 请求重置密码
	router.POST("/gdesign/reset-pwd", middleware.GlobalRateLimitMiddleware(), userHandler.ResetPassword)                      // 重置密码
}

func setupAuthRoutes(router *gin.Engine, userHandler *UserHandler, aiHandler *AIhandler) {
	// 使用 JWT、Rate 中间件保护这些路由
	router.GET("/gdesign/user/me", middleware.JWTAuthMiddleware(), userHandler.CheckMe)
	r := router.Group("/gdesign", middleware.JWTAuthMiddleware(), middleware.UserRateLimitMiddleware())

	// 文法相关接口
	grammar := r.Group("/grammar")
	{
		grammar.POST("/validate", GrammarValidate)            // 文法校验——是否有效
		grammar.POST("/type", GrammarTypeDetermine)           // 判断所给文法的类型
		grammar.POST("/simplify", GrammarSimplify)            // 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
		grammar.POST("/ambiguity", GrammarAmbiguityCheck)     // 正则文法的二义性判断
		grammar.POST("/equivalence", GrammarEquivalenceCheck) // 判断所给的两个正则文法是否等价
		grammar.POST("/recognize", GrammarStringRecognize)    // 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
		grammar.POST("/first", GrammarFirstSet)               // 计算文法的First集
		grammar.POST("/follow", GrammarFollowSet)             // 计算文法的Follow集
	}

	// 正则表达式相关接口
	regEx := r.Group("/regex")
	{
		regEx.POST("/validate", RegexValidate) // 判断是否为有效的正则表达式
	}

	// 自动机相关接口
	automaton := r.Group("/automaton")
	{
		automaton.POST("/validate", AutomatonValidate)         // 是否有效
		automaton.POST("/recognize", AutomatonStringRecognize) // 字符串识别
		automaton.POST("/cleanup", AutomatonCleanup)           // 去无效符号、不可达符号
		automaton.POST("/minimize", DFAMinimize)               // DFA 最小化
		automaton.POST("/nfatodfa", NFADeterminization)        // NFA 转 DFA，NFA确定化
	}

	// 文法、自动机间的转换
	convert := r.Group("/convert")
	{
		convert.POST("/grammar-to-nfa", GrammarToNFA) // 右线性文法转成 NFA
		convert.POST("/fa-to-grammar", FAToGrammar)   // FA 转成文法
		convert.POST("/regex-to-nfa", RegexToNFA)     // 正则表达式转为NFA
		convert.POST("/fa-to-regex", FAToRegex)       // FA转为正则表达式
	}

	// 知识学习
	learn := r.Group("/learn")
	{
		learn.GET("/", LearnGet) // 获取学习资料或信息
	}

	// AI 相关接口
	ai := r.Group("/ai")
	{
		ai.POST("/sse", aiHandler.AIChatSSE) // 流式AI聊天接口，SSE
	}
}
