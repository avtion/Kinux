package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"Kinux/tools/bytesconv"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

func init() {
	RegisterWebsocketOperation(wsOpNewPty, missionPtyRegister)
	RegisterWebsocketOperation(wsOpResize, missionPtyResizeRegister)
	RegisterWebsocketOperation(wsOpStdin, missionPtyStdinRegister)
}

// 创建任务终端
func missionPtyRegister(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	// 从 WebsocketSchedule 获取用户信息
	if ws.Account == nil {
		return errors.New("user info not exist")
	}

	// 获取任务信息
	missionRaw := &struct {
		ID        string `json:"id"`
		Container string `json:"container"`
	}{}
	any.Get("data").ToVal(missionRaw)
	if cast.ToInt(missionRaw.ID) == 0 {
		return errors.New("目标任务为空")
	}
	mission, err := models.GetMission(ws.Context, cast.ToUint(missionRaw.ID))
	if err != nil {
		return
	}

	// 校验命名空间
	d, err := ws.Account.GetDepartment(ws.Context)
	if err != nil {
		return
	}
	if err = d.IsNamespaceAllowed(mission.Namespace); err != nil {
		return
	}

	// 校验容器
	if missionRaw.Container == "" {
		if mission.ExecContainer == "" {
			return errors.New("目标任务未制定容器")
		}
		missionRaw.Container = mission.ExecContainer
	}
	if !mission.IsContainerAllowed(missionRaw.Container) {
		err = errors.New("container不合法")
		return
	}

	// 确定目标容器
	pods, err := NewMissionController(ws.Context).SetAc(ws.Account).SetMission(mission).GetPods("")
	if err != nil {
		return
	}
	if len(pods.Items) == 0 {
		return errors.New("目标任务无可用节点")
	}

	// 从k8s调度器中获取目标任务的POD运行状态
	var pod v1.Pod
	for _, v := range pods.Items {
		if v.Status.Phase == v1.PodRunning {
			pod = v
			break
		}
	}
	if pod.Status.Phase != v1.PodRunning || len(pod.Spec.Containers) == 0 {
		return errors.New("目标任务的节点未准备就绪或无可用容器")
	}
	var c = pod.Spec.Containers[0]
	for _, v := range pod.Spec.Containers {
		if v.Name == missionRaw.Container {
			c = v
		}
	}
	if c.Name == "" {
		return errors.New("目标任务无可用容器")
	}

	// TODO 移除监听者测试
	stdinListenerOpt, stdinListener := NewPtyWrapperListenerOpt(ListenStdin)
	stdinListener.DebugPrint(ws.Context)
	stdoutListenerOpt, stdoutListener := NewPtyWrapperListenerOpt(ListenStdout)
	stdoutListener.DebugPrint(ws.Context)

	// 挂载检查点
	cps, err := models.FindAllTodoCheckpoints(ws.Context, ws.Account.ID, mission.ID, c.Name)
	if err != nil {
		return
	}
	checkpointListenerOpt := NewWrapperForCheckpointCallback(ws.Context, ws.Account, nil, mission, c.Name, cps...)

	go func() {
		if _err := k8s.ConnectToPod(ws.Context, &pod, c.Name, ws.InitPtyWrapper(
			stdinListenerOpt, stdoutListenerOpt, checkpointListenerOpt), mission.GetCommand()); _err != nil {
			logrus.Error("创建POD终端失败", err)
		}
	}()

	return
}

// 调整终端窗体大小处理函数
func missionPtyResizeRegister(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	// 调整窗口大小
	var size = &struct {
		Rows uint16 `json:"rows"`
		Cols uint16 `json:"cols"`
	}{}
	any.Get("data").ToVal(size)

	// 防止 Read 接口发生阻塞
	select {
	case ws.sizeChan <- remotecommand.TerminalSize{Width: size.Cols, Height: size.Rows}:
	default:
	}

	return
}

// 终端写入处理函数
func missionPtyStdinRegister(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	if ws.PtyStdin == nil {
		return errors.New("没有可用的终端")
	}
	_, err = ws.PtyStdin.Write(bytesconv.StringToBytes(any.Get("data").ToString()))
	return nil
}
