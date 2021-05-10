package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// ListMissionCheckpoints 获取实验相关的检查点
func ListMissionCheckpoints(c *gin.Context) {
	params := &struct {
		Page, Size int
		MissionID  uint
		Containers []string
	}{
		Page:       cast.ToInt(c.DefaultQuery("page", "1")),
		Size:       cast.ToInt(c.DefaultQuery("size", "10")),
		MissionID:  cast.ToUint(c.DefaultQuery("mission", "0")),
		Containers: c.QueryArray("containers[]"),
	}
	data, err := models.ListMissionCheckpoints(c, params.MissionID, params.Containers,
		models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	type resStruct struct {
		ID              uint   `json:"id"`
		Percent         uint   `json:"percent"`
		MissionID       uint   `json:"mission_id"`
		CheckpointID    uint   `json:"checkpoint_id"`
		Priority        int    `json:"priority"`
		Mission         string `json:"mission"`
		Checkpoint      string `json:"checkpoint"`
		TargetContainer string `json:"target_container"`
	}

	if len(data) == 0 {
		c.JSON(http.StatusOK, msg.BuildSuccess([]*resStruct{}))
		return
	}

	missionIDs := make([]uint, 0, len(data))
	cpIDs := make([]uint, 0, len(data))
	for _, v := range data {
		missionIDs = append(missionIDs, v.Mission)
		cpIDs = append(cpIDs, v.CheckPoint)
	}
	missionsMapper, err := models.GetMissionsNameMapper(c, missionIDs...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	cpMapper, err := models.GetCheckpointsNameMapper(c, cpIDs...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	var res = make([]*resStruct, 0, len(data))
	for _, v := range data {
		res = append(res, &resStruct{
			ID:              v.ID,
			Percent:         v.Percent,
			MissionID:       v.Mission,
			CheckpointID:    v.CheckPoint,
			Priority:        v.Priority,
			Mission:         missionsMapper[v.Mission],
			Checkpoint:      cpMapper[v.CheckPoint],
			TargetContainer: v.TargetContainer,
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// CountMissionCheckpoints 统计实验相关检查点
func CountMissionCheckpoints(c *gin.Context) {
	params := &struct {
		MissionID  uint
		Containers []string
	}{
		MissionID:  cast.ToUint(c.DefaultQuery("mission", "0")),
		Containers: c.QueryArray("containers[]"),
	}
	res, err := models.CountMissionCheckpoints(c, params.MissionID, params.Containers)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// GetMissionCheckpointsPercent 获取实验相关检查点已占比例
func GetMissionCheckpointsPercent(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	res, err := models.CountMissionCheckpointPercent(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// AddMissionCheckpoint 新增实验检查点
func AddMissionCheckpoint(c *gin.Context) {
	params := &struct {
		Mission         uint   `json:"mission" binding:"required" `
		CheckPoint      uint   `json:"check_point" binding:"required"`
		Percent         uint   `json:"percent" binding:"required"`
		Priority        int    `json:"priority"`
		TargetContainer string `json:"target_container" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.AddMissionCheckpoint(c, &models.MissionCheckpoints{
		Mission:         params.Mission,
		CheckPoint:      params.CheckPoint,
		Percent:         params.Percent,
		Priority:        params.Priority,
		TargetContainer: params.TargetContainer,
	}); err != nil {
		if strings.Contains(err.Error(), "mission_checkpoint") {
			err = fmt.Errorf("检查点(id: %d)已经挂载容器(id: %s)", params.CheckPoint, params.TargetContainer)
		}
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("监测点添加成功"))
}

// EditMissionCheckpoint 修改实验检查点
func EditMissionCheckpoint(c *gin.Context) {
	params := &struct {
		ID              uint   `json:"id"  binding:"required"`
		Mission         uint   `json:"mission" binding:"required"`
		CheckPoint      uint   `json:"check_point" binding:"required"`
		Percent         uint   `json:"percent" binding:"required"`
		Priority        int    `json:"priority"`
		TargetContainer string `json:"target_container" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.EditMissionCheckpoint(c, &models.MissionCheckpoints{
		Model: gorm.Model{
			ID: params.ID,
		},
		Mission:         params.Mission,
		CheckPoint:      params.CheckPoint,
		Percent:         params.Percent,
		Priority:        params.Priority,
		TargetContainer: params.TargetContainer,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("监测点修改成功"))
}

// DeleteMissionCheckpoint 删除实验监测点
func DeleteMissionCheckpoint(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if err := models.DeleteMissionCheckpoint(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("监测点删除成功"))
}

// MissionCheckpointWithRaw 实验考点数据
type MissionCheckpointWithRaw struct {
	ID              uint   `json:"id"`
	Mission         uint   `json:"mission"`
	CheckPoint      uint   `json:"check_point"`
	Percent         uint   `json:"percent"`
	Priority        int    `json:"priority"`
	TargetContainer string `json:"target_container"`
	CheckpointID    uint   `json:"checkpoint_id"`
	CpName          string `json:"cp_name"`
	CpDesc          string `json:"cp_desc"`
	CpCommand       string `json:"cp_command"`
	CpMethod        uint   `json:"cp_method"`
	IsDone          bool   `json:"is_done"`
}

// GetCheckpoints 获取实验的检查点
func GetCheckpoints(c *gin.Context) {
	params := &struct {
		Exam    uint `json:"exam"`
		Mission uint `json:"mission"`
		Lesson  uint `json:"lesson"`
	}{
		Exam:    cast.ToUint(c.DefaultQuery("exam", "0")),
		Mission: cast.ToUint(c.DefaultQuery("mission", "0")),
		Lesson:  cast.ToUint(c.DefaultQuery("lesson", "0")),
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// TODO 支持考试自定义考点查询
	// 首先获取全部的检查点
	mcs, err := models.FindAllMissionCheckpoints(c, params.Mission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 检查点信息
	cpIDs := make([]uint, 0, len(mcs))
	for _, v := range mcs {
		cpIDs = append(cpIDs, v.CheckPoint)
	}
	cps, err := models.FindCheckpoints(c, cpIDs...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	var cpsMapping = make(map[uint]*models.Checkpoint, len(cps))
	for k := range cps {
		cpsMapping[cps[k].ID] = cps[k]
	}

	// 找到需要需要完成的检查点
	mcp, err := models.FindAllTodoMissionCheckpointsV2(c, ac.ID, params.Lesson, params.Exam, params.Mission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	// map[容器名]map[检查点ID]struct{}
	var mcpMapping = make(map[string]map[uint]struct{})
	for _, v := range mcp {
		if _, isExist := mcpMapping[v.TargetContainer]; !isExist {
			mcpMapping[v.TargetContainer] = make(map[uint]struct{})
		}
		mcpMapping[v.TargetContainer][v.CheckPoint] = struct{}{}
	}

	// 结果
	type resType struct {
		ContainerName string                      `json:"container_name"`
		Data          []*MissionCheckpointWithRaw `json:"data"`
	}
	var res = make([]*resType, 0)
	var resMapping = make(map[string]*resType, 0)
	for _, v := range mcs {
		cp, isExist := cpsMapping[v.CheckPoint]
		if !isExist {
			continue
		}
		var _res *resType
		_res, isExist = resMapping[v.TargetContainer]
		if !isExist {
			_res = &resType{
				ContainerName: v.TargetContainer,
				Data:          make([]*MissionCheckpointWithRaw, 0),
			}
			res = append(res, _res)
			resMapping[v.TargetContainer] = _res
		}

		_res.Data = append(_res.Data, &MissionCheckpointWithRaw{
			ID:              v.ID,
			Mission:         v.Mission,
			CheckPoint:      v.CheckPoint,
			Percent:         v.Percent,
			Priority:        v.Priority,
			TargetContainer: v.TargetContainer,
			CheckpointID:    cp.ID,
			CpName:          cp.Name,
			CpDesc:          cp.Desc,
			CpCommand: func() string {
				switch cp.Method {
				case models.MethodStdout, models.MethodTargetPort:
					return cp.Out
				default:
					return cp.In
				}
			}(),
			CpMethod: cp.Method,
			// TODO 应该反着过来查询
			IsDone: func() bool {
				if _, _isExist := mcpMapping[v.TargetContainer]; _isExist {
					if _, ok := mcpMapping[v.TargetContainer][v.CheckPoint]; ok {
						return false
					}
				}
				return true
			}(),
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
