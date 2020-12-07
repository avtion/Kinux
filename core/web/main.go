/*
	Kinux Web
	使用Gin框架拉起一个Web服务
*/
package web

import (
	"Kinux/core/web/routers"
	"Kinux/tools/cfg"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Web服务关闭通道
var shutDownNotify = make(chan struct{})

// 初始化Web服务
func InitWebService() {
	logrus.Info(cfg.DefaultConfig.Web.Enable)
	if !cfg.DefaultConfig.Web.Enable {
		return
	}
	logrus.Trace("Web服务正在启动")
	gin.SetMode(cfg.DefaultConfig.Web.Mode)
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
	svr.RegisterOnShutdown(func() {
		// TODO 通知系统
	})
	go func() {
		if err = svr.ListenAndServe(); err != nil {
			logrus.Error(err)
			return
		}
	}()

	go func() {
		// 需要处理两个停止信号，一个是系统，另外一个是程序
		osNotify := make(chan os.Signal)
		signal.Notify(osNotify, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-osNotify:
		case <-shutDownNotify:
		}
		if err := svr.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("Web服务关闭")
	}()
	return
}

// 关闭Web服务
func ShutDown() (err error) {
	shutDownNotify <- struct{}{}
	return
}
