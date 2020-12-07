// K8S集群容器探针
package health

import (
	"Kinux/tools/cfg"
	"github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
	"net/http"
)

var health = healthcheck.NewHandler()

// 初始化活性探针
func InitHealCheck() {
	var port = "8700"
	if cfg.DefaultConfig.Live.Port != "" {
		port = cfg.DefaultConfig.Live.Port
	}
	go func() {
		logrus.Infof("Kinux活性探针启动(端口: %s)", port)
		if err := http.ListenAndServe(":"+port, health); err != nil {
			logrus.Error(err)
			return
		}
	}()
}
