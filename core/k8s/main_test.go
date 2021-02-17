package k8s

import (
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
	"time"
)

// 测试环境初始化
func init() {
	// 读取配置
	configRawData, err := os.ReadFile("../../kubeConfig")
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logrus.TraceLevel)
	if err := initKubernetes(configRawData, "kinux"); err != nil {
		logrus.Panic(err)
	}
}

// 测试k8s集群访问情况
// 通过当前目录下的kubeConfig文件读取集群信息并测试联通情况
func Test_AccessToCluster(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取节点信息
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), v1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range nodes.Items {
		t.Logf("节点名: %s | 节点信息: %v", node.Name, node.Status)
	}

	// 获取POD信息
	pods, err := clientSet.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pods.String())
}
