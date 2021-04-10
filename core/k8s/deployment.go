// Package k8s Deployment 用于根据用户ID和JobID创建对应的Deployment
// 总体流程：根据AccountID检查用户已经有正在使用的Deployment，销毁其他正在工作的Deployment之后根据JobID获取对应的
// Deployment信息来创建新的Deployment
package k8s

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	appV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/yaml"
)

type DeploymentOption func(d *appV1.Deployment) (err error)

// 删除规则：直至Deployment所有资源被释放后才释放Deployment
var deletePolicy = metaV1.DeletePropagationForeground

// NewDeployment 创建新的Deployment
func NewDeployment(ctx context.Context, dp *appV1.Deployment, s labels.Set, opts ...DeploymentOption) (err error) {
	if dp == nil {
		return errors.New("deployment config is nil")
	}
	// 应用缺省方法
	if len(opts) > 0 {
		for _, opt := range opts {
			if err = opt(dp); err != nil {
				return
			}
		}
	}

	// 解决命名空间冲突的问题
	ns := namespace
	if ns != dp.GetNamespace() {
		ns = dp.GetNamespace()
	}

	// 确保label被正确挂载
	if dpLabels := labels.Set(dp.GetLabels()); !dpLabels.Has(s.String()) {
		s = labels.Merge(dpLabels, s)
		dp.SetLabels(s)

		var podLabels = labels.Set(dp.Spec.Template.GetLabels())
		for k, v := range s {
			podLabels[k] = v
		}
		dp.Spec.Template.SetLabels(podLabels)
		dp.Spec.Selector.MatchLabels = podLabels
	}

	// 创建新的Deployment
	dp, err = clientSet.AppsV1().Deployments(ns).Create(ctx, dp, metaV1.CreateOptions{})
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.WithFields(map[string]interface{}{
		"dp": dp.Name,
	}).Info("Deployment创建成功")
	return
}

// 解析Deployment的配置文件
func ParseDeploymentConfig(fileRaw []byte, strict bool) (dp *appV1.Deployment, err error) {
	dp = new(appV1.Deployment)
	if strict {
		err = yaml.UnmarshalStrict(fileRaw, dp)
	} else {
		err = yaml.Unmarshal(fileRaw, dp)
	}
	if err != nil {
		return
	}
	return
}

// 删除Deployment
func DeleteDeployment(ctx context.Context, ns string, dpName string) (err error) {
	return clientSet.AppsV1().Deployments(ns).Delete(ctx, dpName, metaV1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// 批量查询Deployment
func ListDeployments(ctx context.Context, ns string, s labels.Set) (dps *appV1.DeploymentList, err error) {
	if s == nil || len(s) == 0 {
		err = errors.New("批量查询Deployment失败: 选择器为空")
		return
	}

	// 修复因命名空间导致的错误
	if ns == "" {
		ns = namespace
	}

	dps, err = clientSet.AppsV1().Deployments(ns).List(ctx, metaV1.ListOptions{
		LabelSelector: s.String(),
	})
	if err != nil {
		return
	}
	return
}

// WatchDeploymentsToReady 监听Deployment
func WatchDeploymentsToReady(ctx context.Context, ns string, s labels.Set, errCh chan<- error) {
	w, err := clientSet.AppsV1().Deployments(ns).Watch(ctx, metaV1.ListOptions{
		LabelSelector: s.String(),
	})
	if err != nil {
		return
	}
	go func() {
		defer w.Stop()
		for {
			select {
			case <-ctx.Done():
				logrus.Trace("WatchDeploymentsToReady 上下文结束")
				errCh <- ctx.Err()
				return
			case event, ok := <-w.ResultChan():
				if !ok {
					return
				}
				logrus.WithField("event", event).Trace("WatchDeploymentsToReady 监听到事件")
				switch event.Type {
				case watch.Modified:
					obj, ok := event.Object.(*appV1.Deployment)
					if !ok {
						continue
					}
					for k := range obj.Status.Conditions {
						if obj.Status.Conditions[k].Type != appV1.DeploymentAvailable {
							continue
						}
						if obj.Status.Conditions[k].Status == v1.ConditionTrue {
							errCh <- nil
						}
					}
				default:
				}
			}
		}
	}()
}
