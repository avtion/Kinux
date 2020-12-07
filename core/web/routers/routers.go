package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type initFunc func(r *gin.Engine)

var initQueue = make([]initFunc, 0, 1<<5)

// 初始化路由
func NewRouters() (r *gin.Engine) {
	return newRouters(
		initQueue...,
	)
}

// 内部实现
func newRouters(fns ...initFunc) (r *gin.Engine) {
	r = gin.Default()

	// CORS跨域中间件
	r.Use(cors.New(cors.Config{
		// TODO 修改跨域请求源地址
		AllowAllOrigins: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodOptions,
			http.MethodPost,
			http.MethodPut,
			http.MethodHead,
		},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	for _, fn := range fns {
		fn(r)
	}
	return
}

// 追加路由
func addRouters(fns ...initFunc) {
	initQueue = append(initQueue, fns...)
}
