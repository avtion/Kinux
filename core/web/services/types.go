package services

// 终止标识符EOT
const EndOfTransmission = "\u0004"

type Operation = string

// WebsocketMessage 消息协议中的规定的操作
var Operations = struct {
	Stdin         Operation // 输入
	Stdout        Operation // 输出
	Resize        Operation // 重新调整窗体大小
	Ping          Operation // 保持活跃
	Msg           Operation // 向客户端发送通知
	ResourceApply Operation // 资源申请
}{
	Stdin:         "stdin",
	Stdout:        "stdout",
	Resize:        "resize",
	Ping:          "ping",
	Msg:           "msg",
	ResourceApply: "apply",
}

// WebsocketMessage C/S消息协议
type WebsocketMessage struct {
	Operation Operation   `json:"operation"`
	Data      interface{} `json:"data"`
	// 窗体大小 - 对应 remotecommand.TerminalSize 结构体
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}

// 任务状态
type MissionStatus = int

const (
	_                    MissionStatus = iota
	MissionStatusStop                  // 未启动
	MissionStatusPending               // 正在启动
	MissionStatusWorking               // 运行中
	MissionStatusDone                  // 已经完成
)

var _ = [...]MissionStatus{MissionStatusStop, MissionStatusPending, MissionStatusWorking, MissionStatusDone}

// 业务层的任务结构体, 用于响应
type Mission struct {
	ID     uint
	Name   string
	Desc   string
	Guide  string
	Status MissionStatus
}
