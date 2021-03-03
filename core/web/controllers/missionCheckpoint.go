package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// 获取实验相关的检查点
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

// 统计实验相关检查点
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

// 获取实验相关检查点已占比例
func GetMissionCheckpointsPercent(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	res, err := models.CountMissionCheckpointPercent(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 新增实验检查点
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

// 修改实验检查点
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

// 删除实验监测点
func DeleteMissionCheckpoint(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if err := models.DeleteMissionCheckpoint(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("监测点删除成功"))
}
