package routers

import (
	"Kinux/tools/cfg"
	"github.com/spf13/viper"
)

func init() {
	cfg.InitConfig(func(v *viper.Viper) {
		v.AddConfigPath("../../../")
	})
}
