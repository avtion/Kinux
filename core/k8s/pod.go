package k8s

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 根据Account和Job获取对应的POD
func GetPods(ctx context.Context, account, job string) (p *coreV1.PodList, err error) {
	return clientSet.CoreV1().
		Pods(namespace).
		List(ctx,
			metaV1.ListOptions{LabelSelector: selectorFactory(account, selectorJobOpt(job)).String()})
}
