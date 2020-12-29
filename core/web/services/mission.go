package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

// 批量获取任务信息
func ListMissions(c *gin.Context, u *models.Account, name string, ns []string, page, size int) (res []*Mission, err error) {
	if u == nil {
		err = errors.New("用户信息不存在")
		return
	}

	// 获取用户资料后获取班级信息，用于确定命名空间的访问
	profile, err := u.GetProfile(c)
	if err != nil {
		return
	}
	department, err := profile.GetDepartment(c)
	if err != nil {
		return
	}
	departmentNS := department.GetNS()
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
			if !strings.Contains(department.Namespace, v) {
				err = errors.New("命名空间合法访问")
				return
			}
		}
	} else {
		// 直接访问班级命名空间
		ns = departmentNS
	}

	// 从数据库中查询对应命名空间的任务集合
	ms, err := models.ListMissions(c, name, ns, models.NewPageBuilder(page, size))
	if err != nil {
		return
	}

	// TODO 获取已完成的任务

	// 从K8S调度模块查询Deployment的情况
	dps, err := k8s.ListDeployments(c, "", NewLabelMarker().WithAccount(u.ID).Do())
	if err != nil {
		return
	}
	var dpStatusMapper = make(map[uint]MissionStatus, len(dps.Items))
	if len(dps.Items) > 0 {
		for _, item := range dps.Items {
			// 如果可用的副本等于要求的副本数量
			missionID := cast.ToUint(labels.Set(item.GetLabels()).Get(missionLabel))
			if item.Status.AvailableReplicas == *item.Spec.Replicas {
				dpStatusMapper[missionID] = MissionStatusWorking
			} else {
				dpStatusMapper[missionID] = MissionStatusPending
			}
		}
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
