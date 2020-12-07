package msg

import "github.com/sirupsen/logrus"

type Result = map[string]interface{}

type BuildOption func(r Result)

// 构造器
func Build(code int, data interface{}, opts ...BuildOption) (res Result) {
	if code == 0 {
		code = CodeFailed
	}
	res = Result{
		"code": code,
		"data": data,
	}
	for _, opt := range opts {
		opt(res)
	}
	return
}

// 日志输出
func WithLogPrint(l ...logrus.Level) BuildOption {
	return func(r Result) {
		logLv := logrus.DebugLevel
		if len(l) > 0 {
			logLv = l[0]
		}
		logrus.StandardLogger().Log(logLv, r)
	}
}

// 快速构建失败信息
func BuildFailed(data interface{}, opts ...BuildOption) (res Result) {
	return Build(CodeFailed, data, opts...)
}

// 快速构建成功信息
func BuildSuccess(data interface{}, opts ...BuildOption) (res Result) {
	return Build(CodeSuccess, data, opts...)
}
