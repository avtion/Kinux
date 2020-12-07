package k8s

import (
	"Kinux/tools/cfg"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sConfig *rest.Config
	clientSet *kubernetes.Clientset
	namespace string
)

// 初始化集群配置信息
func InitKubernetes() {
	raw, err := ioutil.ReadFile(cfg.DefaultConfig.Kubernetes.KubeConfigPath)
	if err != nil {
		logrus.Panic(err)
		return
	}
	if err = initKubernetes(raw, cfg.DefaultConfig.Kubernetes.Namespace); err != nil {
		logrus.Panic(err)
		return
	}
}

// 初始化集群配置信息内部实现
func initKubernetes(kubeConfig []byte, ns string) (err error) {
	// 初始化k8s配置
	cc, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return
	}
	k8sConfig, err = cc.ClientConfig()
	if err != nil {
		return
	}

	// 初始化集群
	clientSet, err = kubernetes.NewForConfig(k8sConfig)

	// 初始化命名空间
	namespace = ns
	return
}
