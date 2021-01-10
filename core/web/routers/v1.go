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
		v1 := r.Group("v1")

		// 用户登陆接口不使用鉴权中间件
		v1.POST("/account/login", controllers.LoginAccount)

		// 挂载JWT鉴权中间件
		v1 = v1.Group("", middlewares.JsonWebTokenAuth)

		// Account
		//ac := v1.Group("/account")
		{
			//ac.POST("/login", controllers.LoginAccount)
		}

		// WebSocket
		ws := v1.Group("/ws")
		{
			ws.GET("/", controllers.WebSocketHandlerV1)
		}

		// 任务相关
		ms := v1.Group("/mission")
		{
			ms.GET("/", controllers.QueryMissions)        // 批量查询任务
			ms.POST("/:id/", controllers.NewMission)      // 创建任务
			ms.DELETE("/:id/", controllers.DeleteMission) // 删除任务
		}

	}
}
