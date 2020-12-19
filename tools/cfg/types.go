package cfg

// 配置
type config struct {
	Kubernetes k8s
	Web        web
	Live       live
	Database   database
	Casbin     casbin
	LogLevel   string
}

// 集群配置
type k8s struct {
	IsInCluster    bool
	KubeConfigPath string
	Namespace      string
}

// Web配置
type web struct {
	Enable bool
	Port   string
	Mode   string
}

// 活性探针
type live struct {
	Port string
}

// 数据库
type database struct {
	Name string
	Dsn  string
}

// Casbin鉴权中间件
type casbin struct {
	Enable bool
}
