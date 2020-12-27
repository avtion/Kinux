// Deployment 用于根据用户ID和JobID创建对应的Deployment
// 总体流程：根据AccountID检查用户已经有正在使用的Deployment，销毁其他正在工作的Deployment之后根据JobID获取对应的
// Deployment信息来创建新的Deployment
package k8s

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	appV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/yaml"
)

type DeploymentOption func(d *appV1.Deployment) (err error)

// 删除规则：直至Deployment所有资源被释放后才释放Deployment
var deletePolicy = metaV1.DeletePropagationForeground

// 创建新的Deployment
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
func DeleteDeployments(ctx context.Context, ns string, s labels.Set) (err error) {
	if s == nil || len(s) == 0 {
		return errors.New("删除Deployment失败: 选择器为空")
	}

	// 修复因命名空间导致的错误
	if ns == "" {
		ns = namespace
	}

	// 检查用户是否有正在使用的Deployment
	dps, err := clientSet.AppsV1().Deployments(ns).List(ctx, metaV1.ListOptions{
		LabelSelector: s.String(),
	})
	if err != nil {
		return
	}

	// 逐个将原本使用的Deployment删除
	if len(dps.Items) > 0 {
		logrus.WithField("选择器", s).Debug("发现用户拥有正在使用的Deployment")
		// TODO 防止误操作
		for _, deployment := range dps.Items {
			if err = clientSet.AppsV1().Deployments(deployment.Namespace).Delete(ctx, deployment.Name, metaV1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				logrus.Error(err)
				return
			}
			logrus.WithField("deployment", deployment.Name).Debug("删除用户正在使用的Deployment")
		}
	}
	return
}
