package middlewares

import (
	"Kinux/tools/cfg"
	"github.com/spf13/viper"
)

// 初始化配置文件
func init() {
	cfg.InitConfig(func(v *viper.Viper) {
		v.AddConfigPath("../../../")
	})
}
