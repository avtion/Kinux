package services

import (
	"Kinux/core/web/models"
	"context"
	"errors"
	"math"
	"sort"
)

type scoreListData struct {
	ID                uint    `json:"id"`
	Account           uint    `json:"account"`
	DepartmentID      uint    `json:"department_id"`
	Username          string  `json:"username"`
	Department        string  `json:"department"`
	Score             float64 `json:"score"`
	FinishCheckpoints uint    `json:"finish_checkpoints"`
}

type ScoresListResult struct {
	Total       uint             `json:"total"`
	CpsNum      int              `json:"cps_num"`
	FinishCount int              `json:"finish_count"`
	Data        []*scoreListData `json:"data"`
}

func ListMissionScores(ctx context.Context, mission uint, department uint) (res *ScoresListResult, err error) {
	if department == 0 || mission == 0 {
		err = errors.New("非法参数")
		return
	}
	// 获取班级和实验对象
	dp, err := models.GetDepartmentByID(ctx, department)
	if err != nil {
		return
	}
	ms, err := models.GetMission(ctx, mission)
	if err != nil {
		return
	}

	// 获取用户列表
	acs, err := models.ListAccountsWithProfiles(ctx, nil,
		models.AccountRoleFilter(models.RoleNormalAccount),
		models.AccountDepartmentFilter(int(dp.ID)),
	)
	if err != nil {
		return
	}
	res = &ScoresListResult{Total: ms.Total}
	if len(acs) == 0 {
		res.Data = make([]*scoreListData, 0)
		return
	}

	// 获取用户ID和初始化结果
	acsIDs := make([]uint, 0, len(acs))
	res.Data = make([]*scoreListData, 0, len(acs))
	acMapper := make(map[uint]*scoreListData, len(acs)) // id -> *ScoresListResult
	for _, v := range acs {
		acsIDs = append(acsIDs, v.ID)
		_res := &scoreListData{
			ID:                0,
			Account:           v.ID,
			DepartmentID:      v.DepartmentId,
			Username:          v.Username,
			Department:        v.Department,
			Score:             0,
			FinishCheckpoints: 0,
		}
		res.Data = append(res.Data, _res)
		acMapper[v.ID] = _res
	}

	// 获取检查点
	cps, err := models.FindAllMissionCheckpoints(ctx, mission)
	if err != nil {
		return
	}
	res.CpsNum = len(cps)
	// 构建二维匹配哈希 x=目标容器 y=目标检查点 --> 成绩比例
	var cpsMapper = func() (mapper map[string]map[uint]uint) {
		mapper = make(map[string]map[uint]uint)
		for _, cp := range cps {
			if _, isExist := mapper[cp.TargetContainer]; !isExist {
				mapper[cp.TargetContainer] = make(map[uint]uint)
			}
			mapper[cp.TargetContainer][cp.CheckPoint] = cp.Percent
		}
		return
	}()

	mcs, err := models.FindAccountsFinishScore(ctx, acsIDs, mission)
	mcsCounterMapper := make(map[uint]int, len(acs))
	if err != nil {
		return
	}
	for _, mc := range mcs {
		// 防止空指针
		if _, isExist := cpsMapper[mc.Container]; !isExist {
			continue
		}
		if _, isExist := acMapper[mc.Account]; !isExist {
			continue
		}
		percent, isExist := cpsMapper[mc.Container][mc.Checkpoint]
		if !isExist {
			continue
		}
		acMapper[mc.Account].FinishCheckpoints++
		score := float64(ms.Total) * (float64(percent) * 0.01)
		acMapper[mc.Account].Score += math.RoundToEven(score)
		if _, isExist = mcsCounterMapper[mc.Account]; isExist {
			mcsCounterMapper[mc.Account]++
		} else {
			mcsCounterMapper[mc.Account] = 1
		}
		if mcsCounterMapper[mc.Account] == res.CpsNum {
			res.FinishCount++
		}
	}

	sort.Slice(res.Data, func(i, j int) bool {
		return res.Data[i].Score > res.Data[j].Score
	})
	for k, v := range res.Data {
		v.ID = uint(k)
	}
	return
}
