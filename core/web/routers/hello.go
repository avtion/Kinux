package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	addRouters(func(r *gin.Engine) {
		r.GET("hello", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, map[string]string{
				"msg": "hello kinux",
			})
		})
	})
}
