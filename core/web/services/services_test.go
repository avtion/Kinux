package services

import (
	"Kinux/core/k8s"
	"Kinux/tools/cfg"
	"github.com/spf13/viper"
)

func init() {
	cfg.InitConfig(func(v *viper.Viper) {
		v.AddConfigPath("../../../")
	})
	cfg.DefaultConfig.Kubernetes.KubeConfigPath = "../../../kubeConfig"
	k8s.InitKubernetes()
}
