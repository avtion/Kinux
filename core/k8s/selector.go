package k8s

import "k8s.io/apimachinery/pkg/labels"

type selectorOption func(s labels.Set)

// 选择器构建器
func selectorFactory(account string, options ...selectorOption) (s labels.Set) {
	s = labels.Set{
		"account": account,
	}
	for _, opt := range options {
		opt(s)
	}
	return
}

// 添加Job数据
func selectorJobOpt(j string) selectorOption {
	return func(s labels.Set) {
		if j == "" {
			return
		}
		s["job"] = j
	}
}
