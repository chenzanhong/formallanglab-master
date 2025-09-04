/* 
	AI模块
	调用大模型对当前的文法或自动机做相关问题发解答，再格式化返回给前端
	 将 automaton 从真实上下文中注入（不要硬编码）
✅ 前端记录 trace 路径并传给后端
✅ 支持更多功能：complete, equivalence
✅ 增加缓存（相同自动机 + 相同问题不重复调用）
*/
package api

import(
	"github.com/gin-gonic/gin"
)

func AIGet(c *gin.Context){

}