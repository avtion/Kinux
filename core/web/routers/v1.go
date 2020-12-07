package routers

import (
	"Kinux/core/web/controllers"
	"github.com/gin-gonic/gin"
)

func init() {
	addRouters(v1Routers())
}

func v1Routers() initFunc {
	return func(r *gin.Engine) {
		v1 := r.Group("v1")

		// WebSocket
		v1.GET("/ws", controllers.WebSocketHandlerV1)
	}
}
