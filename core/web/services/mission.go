package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

func init() {
	RegisterWebsocketOperation(wsOpNewPty, missionPtyRegister())
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

// 批量获取任务信息
func ListMissions(ctx context.Context, u *models.Account, name string, ns []string, page, size int) (res []*Mission, err error) {
	if u == nil {
		err = errors.New("用户信息不存在")
		return
	}

	// 获取用户资料后获取班级信息，用于确定命名空间的访问
	d, err := u.GetDepartment(ctx)
	if err != nil {
		return
	}
	departmentNS := d.GetNS()
	if len(departmentNS) == 0 {
		return
	}

	if len(ns) > 0 {
		// 如果请求指定命名空间，则需要判断是否合法
		for _, v := range ns {
			// 防止分隔符导致的逻辑错误
			if strings.ContainsRune(v, ';') {
				err = errors.New("输入的命名空间包括分隔符")
				return
			}
			if !strings.Contains(d.Namespace, v) {
				err = errors.New("命名空间合法访问")
				return
			}
		}
	} else {
		// 直接访问班级命名空间
		ns = departmentNS
	}

	// 从数据库中查询对应命名空间的任务集合
	ms, err := models.ListMissions(ctx, name, ns, models.NewPageBuilder(page, size))
	if err != nil {
		return
	}

	// TODO 获取已完成的任务
	dpStatusMapper, err := getDeploymentStatusForMission(ctx, "", NewLabelMarker().WithAccount(u.ID))
	if err != nil {
		return
	}

	// 遍历构造对应的响应结果
	for _, mission := range ms {
		status, isExist := dpStatusMapper[mission.ID]
		if !isExist {
			status = MissionStatusStop
		}
		res = append(res, &Mission{
			ID:     mission.ID,
			Name:   mission.Name,
			Desc:   mission.Desc,
			Guide:  mission.Guide,
			Status: status,
		})
	}

	return
}

// 根据Deployment的状态获取对应任务的状态
func getDeploymentStatusForMission(ctx context.Context, namespace string, l *labelMaker) (
	res map[uint]MissionStatus, err error) {
	if l == nil {
		l = NewLabelMarker()
	}

	// 从K8S调度模块查询Deployment的情况
	dps, err := k8s.ListDeployments(ctx, namespace, l.Do())
	if err != nil {
		return
	}

	// 遍历
	res = make(map[uint]MissionStatus, len(dps.Items))
	if len(dps.Items) > 0 {
		for _, item := range dps.Items {
			// 如果可用的副本等于要求的副本数量
			missionID := cast.ToUint(labels.Set(item.GetLabels()).Get(missionLabel))

			// FIX 修复mission标签为空的情况
			if missionID == 0 {
				continue
			}

			// 检查可用副本数量
			if item.Status.AvailableReplicas == *item.Spec.Replicas {
				res[missionID] = MissionStatusWorking
			} else {
				res[missionID] = MissionStatusPending
			}
		}
	}

	return
}

// 任务操作
type MissionOperation = uint

const (
	_             MissionOperation = iota
	MissionCreate                  // 创建
	MissionDelete                  // 删除
)

// 用户账号与任务绑定操作
func AccountMissionOpera(ctx context.Context, ac *models.Account,
	targetMission uint, operation MissionOperation) (err error) {
	defer func() {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"account":   ac.ID,
				"mission":   targetMission,
				"operation": operation,
			}).Error("任务操作失败")
		}
	}()
	if ac == nil {
		return errors.New("用户信息为空")
	}

	ms, err := models.GetMission(ctx, targetMission)
	if err != nil {
		return
	}

	// 校验任务的命名空间是否被允许访问
	d, err := ac.GetDepartment(ctx)
	if err != nil {
		return
	}
	if err = d.IsNamespaceAllowed(ms.Namespace); err != nil {
		return
	}

	// 创建任务控制器
	switch operation {
	case MissionCreate:
		err = NewMissionController(ctx).SetAc(ac).SetMission(ms).NewDeployment()
	case MissionDelete:
		err = NewMissionController(ctx).SetAc(ac).SetMission(ms).DestroyDeployment()
	default:
		return errors.New("unknown mission operation")
	}
	return
}

// 创建任务终端
func missionPtyRegister() WsOperationHandler {
	return func(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
		// 从 gin.Context 上下文获取用户信息
		ac, err := GetAccountFromCtx(ws.Context)
		if err != nil {
			return
		}

		// 获取任务信息
		missionRaw := &struct {
			ID        uint   `json:"id"`
			Container string `json:"container"`
		}{}
		any.Get("data").ToVal(missionRaw)
		if missionRaw.ID == 0 {
			return errors.New("目标任务不存在")
		}
		mission, err := models.GetMission(ws.Context, missionRaw.ID)
		if err != nil {
			return
		}

		// 校验命名空间
		d, err := ac.GetDepartment(ws.Context)
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
		pods, err := NewMissionController(ws.Context).SetAc(ac).SetMission(mission).GetPods("")
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

		// 使用装饰器维护websocket链接，并终止守护协程
		// TODO 并发情况下会导致通道重复关闭并panic
		ws.StopDaemon()

		return k8s.ConnectToPod(ws.Context, &pod, c.Name, ws.InitPtyWrapper(), mission.GetCommand())
	}
}
