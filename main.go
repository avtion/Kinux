package main

import (
	"Kinux/core/k8s"
	"Kinux/core/web"
	"Kinux/tools/cfg"
	"Kinux/tools/health"
	"github.com/sirupsen/logrus"
	"time"
	_ "time/tzdata"
)

func main() {
	// 修复时区问题
	var err error
	time.Local, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Panic(err)
	}
	cfg.InitConfig()
	health.InitHealCheck() // 活性探针
	k8s.InitKubernetes()
	web.InitWebService() // 启动Web服务
}
