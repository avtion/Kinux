package main

import (
	"Kinux/core/k8s"
	"Kinux/core/web"
	"Kinux/tools/cfg"
	"Kinux/tools/health"
)

func main() {
	cfg.InitConfig()
	health.InitHealCheck() // 活性探针
	k8s.InitKubernetes()
	web.InitWebService() // 启动Web服务
}
