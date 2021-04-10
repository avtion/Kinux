package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

func init() {
	RegisterWebsocketOperation(wsOpResetContainers, missionResetRegister)
	RegisterWebsocketOperation(wsOpMissionApply, missionApply)
}

// MissionStatus 任务状态
type MissionStatus = int

const (
	_                    MissionStatus = iota
	MissionStatusStop                  // 未启动
	MissionStatusPending               // 正在启动
	MissionStatusWorking               // 运行中
	MissionStatusDone                  // 已经完成
	MissionStatusBlock                 // 无法进行
)

var _ = [...]MissionStatus{MissionStatusStop, MissionStatusPending, MissionStatusWorking, MissionStatusDone}

// Mission 业务层的任务结构体, 用于响应
type Mission struct {
	ID     uint          `json:"id"`
	Name   string        `json:"name"`
	Desc   string        `json:"desc"`
	Guide  string        `json:"guide"`
	Status MissionStatus `json:"status"`
}

// ListMissions 批量获取任务信息
// Deprecated: 删除命名空间
func ListMissions(_ context.Context, _ *models.Account, _ string, _ []string, _, _ int) (res []*Mission, err error) {
	return []*Mission{}, errors.New("deprecated: 删除命名空间")
}

// ListMissionsV2 获取任务信息
func ListMissionsV2(c *gin.Context, lessonID uint, page, size int) (res []*Mission, err error) {
	ac, err := GetAccountFromCtx(c)
	if err != nil {
		return
	}
	lesson, err := models.GetLesson(c, lessonID)
	if err != nil {
		return
	}

	missionIDs, err := models.GetMissionIDsByLessons(c, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(models.NewPageBuilder(page, size).Build).Where("lesson = ?", lesson.ID)
	})
	if err != nil {
		return
	}

	ms, err := models.GetMissions(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", missionIDs)
	})
	if err != nil {
		return
	}

	dpStatusMapper, err := GetDeploymentStatusForMission(c, "",
		NewLabelMarker().WithAccount(ac.ID).WithLesson(lesson.ID).WithExam(0))
	if err != nil {
		return
	}

	// 遍历构造对应的响应结果
	for _, mission := range ms {
		status, isExist := dpStatusMapper[mission.ID]
		if !isExist {
			status = MissionStatusStop
		}

		// 查询任务是否已经完成
		var cps []*models.Checkpoint
		cps, err = models.FindAllTodoMissionCheckpoints(c, ac.ID, lesson.ID, 0, mission.ID)
		if err != nil {
			return
		}
		if len(cps) == 0 {
			status = MissionStatusDone
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

// GetDeploymentStatusForMission 根据Deployment的状态获取对应任务的状态
func GetDeploymentStatusForMission(ctx context.Context, namespace string, l *labelMaker) (
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

// MissionOperation 任务操作
type MissionOperation = uint

const (
	_             MissionOperation = iota
	MissionCreate                  // 创建
	MissionDelete                  // 删除
)

// AccountMissionOpera 用户账号与任务绑定操作
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

// ListMissionAllowedContainersNames 获取任务允许的容器名列表
func ListMissionAllowedContainersNames(ctx context.Context, missionID uint) (
	res []string, err error) {
	mission, err := models.GetMission(ctx, missionID)
	if err != nil {
		return
	}
	containers, err := mission.ListAllowedContainers(ctx)
	if err != nil {
		return
	}
	for _, c := range containers {
		if c.Name == mission.ExecContainer {
			res = append([]string{c.Name}, res...)
			continue
		}
		res = append(res, c.Name)
	}
	return
}

// GetMissionGuide 获取任务的实验文档
func GetMissionGuide(ctx context.Context, missionID uint) (res string, err error) {
	mission, err := models.GetMission(ctx, missionID)
	if err != nil {
		return
	}
	return mission.Guide, nil
}

// 实验重置处理函数
func missionResetRegister(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	// 从 WebsocketSchedule 获取用户信息
	if ws.Account == nil {
		return errors.New("user info not exist")
	}

	// 获取任务信息
	missionRaw := &struct {
		ID string `json:"id"`
	}{}
	any.Get("data").ToVal(missionRaw)
	if cast.ToInt(missionRaw.ID) == 0 {
		return errors.New("目标任务为空")
	}
	mission, err := models.GetMission(ws.Context, cast.ToUint(missionRaw.ID))
	if err != nil {
		return
	}

	// 使用协程处理重置终端操作避免websocket连接发生阻塞
	go func() {
		ctx, cancel := context.WithTimeout(ws.Context, 2*time.Minute)
		defer cancel()

		mc := NewMissionController(ctx).SetAc(ws.Account).SetMission(mission)

		// 启动监听
		errCh := mc.WatchDeploymentToReady("")

		// 删除所有的可用POD
		if err = mc.ResetMission(""); err != nil {
			return
		}
		select {
		case _err := <-errCh:
			if _err != nil {
				logrus.WithField("重置容器失败", _err)
				return
			}
		case <-ctx.Done():
			logrus.WithField("重置容器超时", ctx.Err())
			return
		}
		// 这里让程序等1秒避免容器没准备好
		<-time.After(1 * time.Second)
		data, err := jsoniter.Marshal(&WebsocketMessage{
			Op:   wsOpContainersDone,
			Data: nil,
		})
		if err != nil {
			return
		}
		ws.SendData(data)
	}()
	return
}

// 启动任务处理函数
func missionApply(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	// 从 WebsocketSchedule 获取用户信息
	if ws.Account == nil {
		return errors.New("user info not exist")
	}

	// 获取任务信息
	missionRaw := &struct {
		ID     string `json:"id"`
		Exam   string `json:"exam"`
		Lesson string `json:"lesson"`
	}{}
	any.Get("data").ToVal(missionRaw)

	// 课程ID
	var lessonID = cast.ToUint(missionRaw.Lesson)

	mission, err := models.GetMission(ws.Context, cast.ToUint(missionRaw.ID))
	if err != nil {
		return
	}

	var exam *models.Exam
	examID := cast.ToUint(missionRaw.Exam)
	if examID != 0 {
		exam, err = models.GetExam(ws.Context, examID)
		if err != nil {
			return
		}
		lessonID = exam.Lesson
	}
	if lessonID == 0 {
		return errors.New("目标课程为空")
	}
	lesson, err := models.GetLesson(ws.Context, lessonID)
	if err != nil {
		return
	}

	// 初始化
	mc := NewMissionController(ws).SetAc(ws.Account).SetLesson(
		lesson).SetExam(exam).SetMission(mission)

	// 启动监听
	errCh := mc.WatchDeploymentToReady("")

	// 创建新的Deployment
	if err = mc.NewDeployment(); err != nil {
		return
	}
	// 等待dp状态更新
	if err = <-errCh; err != nil {
		return
	}

	// 防止K8S通知太快导致容器未完成部署
	<-time.After(1 * time.Second)
	data, err := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpContainersDone,
		Data: nil,
	})
	if err != nil {
		return
	}
	ws.SendData(data)
	return
}
