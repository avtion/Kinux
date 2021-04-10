package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/labels"
)

// ExamListResult 考试列表查询结果
type ExamListResult struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Desc       string     `json:"desc"`
	Total      uint       `json:"total"`
	BeginAt    string     `json:"begin_at"`
	EndAt      string     `json:"end_at"`
	ForceOrder bool       `json:"force_order"`
	Missions   []*Mission `json:"missions"`
}

// ListExams 查询考试列表
func ListExams(ctx context.Context, ac *models.Account, page, size int) (res []*ExamListResult, err error) {
	exams, err := models.ListExams(ctx, models.NewPageBuilder(page, size).Build)
	if err != nil {
		return
	}

	// 初始化结果
	resMapper := make(map[uint]*ExamListResult, len(exams))
	res = make([]*ExamListResult, 0, len(exams))
	examIDs := make([]uint, 0, len(exams))
	if len(exams) == 0 {
		return
	}
	for _, exam := range exams {
		examIDs = append(examIDs, exam.ID)
		_res := &ExamListResult{
			ID:         exam.ID,
			Name:       exam.Name,
			Desc:       exam.Desc,
			Total:      exam.Total,
			BeginAt:    exam.BeginAt.Format("2006-01-02 15:04:05"),
			EndAt:      exam.EndAt.Format("2006-01-02 15:04:05"),
			ForceOrder: exam.ForceOrder,
			Missions:   make([]*Mission, 0),
		}
		resMapper[exam.ID] = _res
		res = append(res, _res)
	}

	// 查询考试关联的实验
	eMissions, err := models.GetExamMissions(ctx, examIDs...)
	if err != nil {
		return
	}
	if len(eMissions) == 0 {
		return
	}
	missionIDs := make([]uint, 0, len(eMissions))
	for _, mission := range eMissions {
		missionIDs = append(missionIDs, mission.ID)
	}

	// 查询实验的详细数据
	missions, err := models.GetMissions(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", missionIDs)
	})
	if err != nil {
		return
	}
	missionsMapper := make(map[uint]*models.Mission, len(missions))
	for k, v := range missions {
		missionsMapper[v.ID] = missions[k]
	}

	// 查询当前Deployment的状态
	dpStatusMapper, err := getDeploymentStatusForExam(ctx, "", NewLabelMarker().WithAccount(ac.ID))
	if err != nil {
		return
	}

	// 追加结果
	for _, v := range eMissions {
		_res, isExist := resMapper[v.Exam]
		if !isExist {
			continue
		}

		// TODO 检查是否有正在运行的Deployment
		// TODO 检查是否完成
		ms, isExist := missionsMapper[v.Mission]
		if !isExist {
			continue
		}

		status, isExist := dpStatusMapper[v.Exam]
		if !isExist {
			status = MissionStatusStop
		}

		_res.Missions = append(_res.Missions, &Mission{
			ID:     ms.ID,
			Name:   ms.Name,
			Desc:   ms.Desc,
			Guide:  ms.Guide,
			Status: status, // TODO 完成状态监测
		})
	}

	return
}

// 根据Deployment的状态获取对应考试的状态
func getDeploymentStatusForExam(ctx context.Context, namespace string, l *labelMaker) (
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
			examID := cast.ToUint(labels.Set(item.GetLabels()).Get(examLabel))

			// FIX 修复mission标签为空的情况
			if examID == 0 {
				continue
			}

			// 检查可用副本数量
			if item.Status.AvailableReplicas == *item.Spec.Replicas {
				res[examID] = MissionStatusWorking
			} else {
				res[examID] = MissionStatusPending
			}
		}
	}

	return
}

// ExamStatus 考试状态
type ExamStatus uint

const (
	_          ExamStatus = iota
	ESNotStart            // 考试未开始
	ESRunning             // 正在考试
	ESFinish              // 考试已经结束
)

// GetExamStatus 获取考试的状态
func GetExamStatus(ctx context.Context, ac uint, exam uint) (res ExamStatus) {
	_eWatcher, isExist := ExamWatchers.Load(ac)
	if isExist {
		eWatcher, _ := _eWatcher.(*ExamWatcher)
		if eWatcher.ELog.Exam == exam {
			return ESRunning
		}
	}
	eLog := &models.ExamLog{
		Account: ac,
		Exam:    exam,
	}
	if err := models.GetGlobalDB().WithContext(ctx).Where(
		"account = ? AND exam = ?", eLog.Account, eLog.Exam).First(eLog).Error; err != nil {
		return ESNotStart
	}
	if eLog.EndAt.IsZero() {
		return ESRunning
	} else {
		return ESFinish
	}
}

// GetAllTodoCheckpointsForExam 查找考试中需要完成的检查点
func GetAllTodoCheckpointsForExam(ctx context.Context, ac, lesson, exam, mission uint, containers ...string) (cps []*models.Checkpoint, err error) {
	// TODO 支持自定义考试考点
	// 默认为实验的考点
	cps, err = models.FindAllTodoMissionCheckpoints(ctx, ac, lesson, exam, mission, containers...)
	return
}
