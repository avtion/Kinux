// Deployment 用于根据用户ID和JobID创建对应的Deployment
// 总体流程：根据AccountID检查用户已经有正在使用的Deployment，销毁其他正在工作的Deployment之后根据JobID获取对应的
// Deployment信息来创建新的Deployment
package k8s

import (
	"Kinux/tools"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	appV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

// Deployment调度主体
type DeploymentJob struct {
	AccountId string
	JobId     string
	Status    uint
}

type DeploymentOption func(d *appV1.Deployment)

// 删除规则：直至Deployment所有资源被释放后才释放Deployment
var deletePolicy = metaV1.DeletePropagationForeground

// 创建新的Deployment
func NewDeployment(ctx context.Context, account, job string, options ...DeploymentOption) (dp *appV1.Deployment, err error) {
	// 检查用户是否有正在使用的Deployment
	dps, err := clientSet.AppsV1().Deployments(namespace).List(ctx, metaV1.ListOptions{
		LabelSelector: selectorFactory(account).String(),
	})
	if err != nil {
		return
	}

	// 逐个将原本使用的Deployment删除
	if len(dps.Items) > 0 {
		logrus.WithField("account", account).Debug("发现用户拥有正在使用的Deployment")
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

	// 标记label
	l := selectorFactory(account, selectorJobOpt(job))

	// 创建新的Deployment
	dp, err = clientSet.AppsV1().Deployments(namespace).Create(ctx, &appV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%s", account, job, tools.GetRandomString(6)),
			Namespace: namespace,
			Labels:    l,
		},
		Spec: appV1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metaV1.LabelSelector{MatchLabels: l},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: l,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "centos7",
							Image:           "centos:centos7",
							Command:         []string{"sh", "-c", "trap : TERM INT; sleep infinity &amp; wait"},
							ImagePullPolicy: v1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}, metaV1.CreateOptions{})
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.WithFields(map[string]interface{}{
		"dp":      dp,
		"account": account,
	}).Info("Deployment创建成功")
	return
}

// 解析Deployment的配置文件
func ParseDeploymentConfig(_ context.Context, fileRaw []byte, strict bool) (dp *appV1.Deployment, err error) {
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
