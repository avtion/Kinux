package routers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	addRouters(debugRouters())
}

// Debug调试专用路由
// TODO 注释
func debugRouters() initFunc {
	return func(r *gin.Engine) {
		var defaultSuccess gin.HandlerFunc = func(c *gin.Context) {
			c.JSON(http.StatusOK, msg.BuildSuccess("ok"))
		}
		debug := r.Group("/debug")
		{
			debug.GET("/login", func(c *gin.Context) {
				token, _, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
					Username: "systemTestAccount",
					Role:     models.RoleAdmin,
				})
				if err != nil {
					c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("生成测试密钥失败"))
					return
				}
				c.JSON(http.StatusOK, map[string]string{
					"msg":   "系统测试密钥生成成功",
					"token": token,
				})
			})
			debug.Use(middlewares.JsonWebTokenAuth).GET("/auth_test", defaultSuccess)
		}
	}
}
