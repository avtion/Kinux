package cfg

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type configInitFunc func(v *viper.Viper)

// 默认配置对象
var DefaultConfig = new(config)

// 初始化日志模块
func initLogModule() {
	// JSON格式输出
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 追踪
	logrus.SetReportCaller(true)
}

// 初始化配置
func InitConfig(fns ...configInitFunc) {
	initLogModule()

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath("./")
	for _, fn := range fns {
		fn(v)
	}
	if err := v.ReadInConfig(); err != nil {
		logrus.Panic(err)
		return
	}
	if err := v.Unmarshal(DefaultConfig); err != nil {
		logrus.Panic(err)
		return
	}

	// 设置日志等级
	lvl, _err := logrus.ParseLevel(DefaultConfig.LogLevel)
	if _err != nil {
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)

	logrus.WithField("配置信息", *DefaultConfig).Trace("配置加载成功")
	return
}
