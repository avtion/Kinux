package k8s

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// 获取POD列表
func GetPods(ctx context.Context, ns string, s labels.Set) (p *coreV1.PodList, err error) {
	if ns == "" {
		ns = namespace
	}
	return clientSet.CoreV1().
		Pods(ns).
		List(ctx,
			metaV1.ListOptions{LabelSelector: s.String()})
}
