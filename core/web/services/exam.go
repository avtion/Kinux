package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/labels"
)

// 考试列表查询结果
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

// 查询考试列表
func ListExams(ctx context.Context, ac *models.Account, namespace string, page, size int) (res []*ExamListResult, err error) {
	exams, err := models.ListExams(ctx, namespace, models.NewPageBuilder(page, size))
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
