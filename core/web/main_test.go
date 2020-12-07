package web

import (
	"Kinux/core/k8s"
	"Kinux/tools/cfg"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func init() {
	cfg.InitConfig(func(v *viper.Viper) {
		v.AddConfigPath("../../")
	})
	cfg.DefaultConfig.Kubernetes.KubeConfigPath = "../../kubeConfig"
	k8s.InitKubernetes()
}

func TestInitWebService(t *testing.T) {
	tests := []struct {
		name     string
		noFinish bool
	}{
		{
			name:     "test",
			noFinish: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitWebService()
			if tt.noFinish {
				select {}
			} else {
				<-time.After(10 * time.Second)
			}
		})
	}
}
