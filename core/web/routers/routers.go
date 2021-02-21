package routers

import (
	"Kinux/tools/translator"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	zhTranslatorPkg "github.com/go-playground/validator/v10/translations/zh"
	"github.com/sirupsen/logrus"
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

	// 修改校验器语言
	if err := zhTranslatorPkg.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate),
		translator.Trans); err != nil {
		err = fmt.Errorf("gin校验器翻译初始化失败: %v", err)
		logrus.Panic(err)
		panic(err)
	}

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
			http.MethodDelete,
		},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
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
