package routers

import (
	"Kinux/core/web/controllers"
	"Kinux/core/web/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	addRouters(v2Router())
}

func v2Router() initFunc {
	return func(r *gin.Engine) {
		v2 := r.Group("v2")

		v2WithAuth := v2.Group("", middlewares.JsonWebTokenAuth)

		lesson := v2WithAuth.Group("/lesson")
		{
			lesson.GET("/", controllers.ListLessons)
		}

		ms := v2WithAuth.Group("/ms")
		{
			ms.GET("/", controllers.ListMissionsV2)
		}
	}
}
