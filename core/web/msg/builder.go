package msg

import (
	"Kinux/tools/translator"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"strings"
)

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
	// 校验器翻译
	errs, ok := data.(validator.ValidationErrors)
	if ok {
		errsMap := errs.Translate(translator.Trans)
		var errsStringSlice = make([]string, 0, len(errsMap))
		for _, v := range errsMap {
			errsStringSlice = append(errsStringSlice, v)
		}
		data = strings.Join(errsStringSlice, " & ")
	}

	// FIX 修复错误类型无法正确输出
	if _err, ok := data.(error); ok {
		data = _err.Error()
	}

	return Build(CodeFailed, data, opts...)
}

// 快速构建成功信息
func BuildSuccess(data interface{}, opts ...BuildOption) (res Result) {
	return Build(CodeSuccess, data, opts...)
}
