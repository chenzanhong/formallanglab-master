package api

import (
	"backend/internal/middleware"

	// aiSvc "backend/internal/service/ai"

	mtr "backend/internal/metrics"

	"github.com/chenzanhong/goutil/jwtx"
	"github.com/gin-gonic/gin"
)

func SetupRouter(storeHandler *StoreHandler, learnHandler *LearnHandler) *gin.Engine {

	router := gin.Default()
	// 1. 恢复中间件 - 最先使用，捕获所有panic
	// router.Use(gin.Recovery())
	// router.GET("/master/ai/ws", aiHandler.AIChatWS) // WebSocket聊天接口，不经过JWT中间件

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

	router.GET("/gdesign/master/metrics", mtr.MetricsHandler())
	router.GET("/gdesign/master/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.HEAD("/gdesign/master/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	setupAuthRoutes(router, learnHandler, storeHandler) // 注册需要认证的路由
	return router
}

func setupAuthRoutes(router *gin.Engine, learnHandler *LearnHandler, storeHandler *StoreHandler) {
	// 使用 JWT、Rate 中间件保护这些路由
	r := router.Group("/gdesign/master", jwtx.GinJWTAuthMiddleware(), middleware.UserRateLimitMiddleware())

	// 文法相关接口
	grammar := r.Group("/grammar")
	{
		grammar.POST("/validate", GrammarValidate)              // 文法校验——是否有效
		grammar.POST("/type", GrammarTypeDetermine)             // 判断所给文法的类型
		grammar.POST("/simplify", GrammarSimplify)              // 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
		grammar.POST("/ambiguity", GrammarAmbiguityCheck)       // 正则文法的二义性判断
		grammar.POST("/equivalence", GrammarEquivalenceCheck)   // 判断所给的两个正则文法是否等价
		grammar.POST("/recognize", GrammarStringRecognize)      // 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
		grammar.POST("/generate", GrammarGenerateExampleString) // 生成可推导和不可推导字符串
		grammar.POST("/first", GrammarFirstSet)                 // 计算文法的First集
		grammar.POST("/follow", GrammarFollowSet)               // 计算文法的Follow集
	}

	// 自动机相关接口
	automaton := r.Group("/automaton")
	{
		automaton.POST("/validate", AutomatonValidate)              // 是否有效
		automaton.POST("/recognize", AutomatonStringRecognize)      // 字符串识别
		automaton.POST("/cleanup", AutomatonCleanup)                // 去无效符号、不可达符号
		automaton.POST("/minimize", DFAMinimize)                    // DFA 最小化
		automaton.POST("/nfatodfa", NFADeterminization)             // NFA 转 DFA，NFA确定化
		automaton.POST("/equivalence", AutomatonEquivalenceCheck)   // 判断所给的两个自动机是否等价
		automaton.POST("/generate", AutomatonGenerateExampleString) // 生成可接受和不可接受字符串
	}

	// 正则表达式相关接口
	regEx := r.Group("/regex")
	{
		regEx.POST("/validate", RegexValidate) // 判断是否为有效的正则表达式
		regEx.POST("/recognize", RegexRecognize)
		regEx.POST("/equivalence", RegexEquivalenceCheck)   // 判断所给的两个正则表达式是否等价
		regEx.POST("/generate", RegexGenerateExampleString) // 生成可匹配和不可匹配字符串
	}

	// 文法、自动机间的转换
	convert := r.Group("/convert")
	{
		convert.POST("/grammar-to-nfa", GrammarToFA) // 右线性文法转成 NFA
		convert.POST("/fa-to-grammar", FAToGrammar)  // FA 转成文法
		convert.POST("/regex-to-nfa", RegexToFA)     // 正则表达式转为NFA
		convert.POST("/fa-to-regex", FAToRegex)      // FA转为正则表达式
	}

	// 知识学习
	learn := r.Group("/learn")
	{
		learn.GET("/", learnHandler.LearnList)                 // 获取学习资料列表
		learn.POST("/", learnHandler.LearnAddMaterial)         // 添加学习资源（管理员在OSS上传后调用）
		learn.POST("/sync", learnHandler.LearnSyncOSSFiles)    // 同步OSS文件到数据库（自动识别新增文件）
		learn.GET("/:id", learnHandler.LearnGetByID)           // 获取单个资源详情
		learn.DELETE("/:id", learnHandler.LearnDeleteMaterial) // 删除学习资源
	}

	// 存储模块接口
	store := r.Group("/store")
	{
		// 自动机存储
		store.POST("/automaton", storeHandler.CreateAutomaton)
		store.GET("/automatons", storeHandler.FindAutomata)
		store.DELETE("/automaton/:id", storeHandler.DeleteAutomaton)
		// 文法存储
		store.POST("/grammar", storeHandler.CreateGrammar)
		store.GET("/grammars", storeHandler.FindGrammars)
		store.DELETE("/grammar/:id", storeHandler.DeleteGrammar)
		// 正则表达式存储
		store.POST("/regex", storeHandler.CreateRegex)
		store.GET("/regexes", storeHandler.FindRegexes)
		store.DELETE("/regex/:id", storeHandler.DeleteRegex)
	}
}
