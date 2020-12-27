package routers

import (
	"Kinux/core/web/controllers"
	"Kinux/core/web/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	addRouters(v1Routers())
}

func v1Routers() initFunc {
	return func(r *gin.Engine) {
		v1 := r.Group("v1", middlewares.JsonWebTokenAuth)

		// WebSocket
		v1.GET("/ws", controllers.WebSocketHandlerV1)

		// Account
		ac := v1.Group("/account")
		{
			ac.POST("/login", controllers.LoginAccount)
		}

		// 任务相关
		ms := v1.Group("/mission")
		{
			// 批量查询任务
			ms.GET("/", controllers.QueryMissions)
		}

	}
}
