/* 
	有限状态自动机
	FSMValidate                             // 是否有效
    FSMCleanup                              // 去无效符号、不可达符号
    DFAMinimize                             // DFA 最小化
	NFAToDFA                                // NFA 转 DFA
	FSMStringRecognize						// 字符串识别
*/
package api

import (
	"github.com/gin-gonic/gin"
)

func FSMValidate(c *gin.Context){ // 是否有效

}

func FSMCleanup(c *gin.Context){ // 去无效符号、不可达符号

}

func DFAMinimize(c *gin.Context){ // DFA 最小化

}

func NFAToDFA(c *gin.Context){ // NFA 转 DFA

}

func FSMStringRecognize(c *gin.Context){ // 字符串识别

}