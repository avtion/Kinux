//  Client 用于与目标POD中的Container建立连接
package k8s

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

type clientReqOption func(r *rest.Request)

// PtyHandler 定义remoteCommand所需要的方法集
type PtyHandler interface {
	io.ReadWriter
	remotecommand.TerminalSizeQueue // 调整终端大小
	Done()                          // 终止
}

// Pod建立连接
func ConnectToPod(ctx context.Context, p *coreV1.Pod, container string, pty PtyHandler, cmd []string,
	options ...clientReqOption) (err error) {
	// 关闭pty连接
	defer func() {
		logrus.Debug("pty被释放")
		pty.Done()
	}()

	// 设置默认容器
	if len(container) == 0 {
		if len(p.Spec.Containers) > 0 {
			container = p.Spec.Containers[0].Name
		} else {
			return errors.New("pod has no container")
		}
	}

	// 设置默认执行的命令
	if len(cmd) == 0 {
		cmd = []string{"/bin/sh"}
	}

	// 创建请求
	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(p.GetNamespace()).
		SubResource("exec").
		Name(p.Name).
		VersionedParams(&coreV1.PodExecOptions{
			Container: container,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Command:   cmd,
		}, scheme.ParameterCodec)
	for _, opt := range options {
		opt(req)
	}

	// 创建SPDY执行器
	executor, err := remotecommand.NewSPDYExecutor(k8sConfig, http.MethodPost, req.URL())
	if err != nil {
		return
	}
	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
}
