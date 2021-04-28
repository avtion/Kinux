/*
	Kinux Web
	使用Gin框架拉起一个Web服务
*/
package web

import (
	"Kinux/core/web/routers"
	"Kinux/core/web/services"
	"Kinux/tools/cfg"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// InitWebService 初始化Web服务
func InitWebService() {
	logrus.Info(cfg.DefaultConfig.Web.Enable)
	if !cfg.DefaultConfig.Web.Enable {
		return
	}
	logrus.WithField("端口", cfg.DefaultConfig.Web.Port).Trace("Web服务正在启动")
	gin.SetMode(cfg.DefaultConfig.Web.Mode)

	// 初始化监考
	services.InitExamWatcher(context.Background())

	if err := initWebService(cfg.DefaultConfig.Web.Port); err != nil {
		logrus.Error(err)
		return
	}
	return
}

// 初始化Web服务内部实现
func initWebService(port string) (err error) {
	svr := &http.Server{
		Addr:              ":" + port,
		Handler:           routers.NewRouters(),
		ReadHeaderTimeout: 1 * time.Minute,
		WriteTimeout:      10 * time.Minute,
	}
	return svr.ListenAndServe()
}
