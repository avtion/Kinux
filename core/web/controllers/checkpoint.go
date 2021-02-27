package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
)

// 获取检查点数据
func ListCheckpoints(c *gin.Context) {
	params := &struct {
		Page, Size int
		Name       string
		Method     uint
	}{
		Page:   cast.ToInt(c.DefaultQuery("page", "1")),
		Size:   cast.ToInt(c.DefaultQuery("size", "10")),
		Name:   c.DefaultQuery("name", ""),
		Method: cast.ToUint(c.DefaultQuery("method", "0")),
	}
	data, err := models.ListCheckpoints(c, params.Name,
		params.Method, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译
	type resType struct {
		ID        uint   `json:"id"`
		Method    uint   `json:"method"`
		Name      string `json:"name"`
		Desc      string `json:"desc"`
		In        string ` json:"in"`
		Out       string `json:"out"`
		CreatedAt string ` json:"created_at"`
	}

	var res = make([]*resType, 0, len(data))
	for _, v := range data {
		res = append(res, &resType{
			ID:        v.ID,
			Method:    v.Method,
			Name:      v.Name,
			Desc:      v.Desc,
			In:        v.In,
			Out:       v.Out,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 统计检查点
func CountCheckpoints(c *gin.Context) {
	params := &struct {
		Name   string
		Method uint
	}{
		Name:   c.DefaultQuery("name", ""),
		Method: cast.ToUint(c.DefaultQuery("method", "0")),
	}
	res, err := models.CountCheckpoints(c, params.Name, params.Method)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 快速获取检查点选项
func QuickListCheckpoints(c *gin.Context) {
	params := &struct {
		Name   string
		Method uint
	}{
		Name:   c.DefaultQuery("name", ""),
		Method: cast.ToUint(c.DefaultQuery("method", "0")),
	}
	res, err := models.QuickListCheckpoint(c, params.Name, params.Method)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 创建检查点
func AddCheckpoint(c *gin.Context) {
	params := &struct {
		Method uint   `json:"method"`
		Name   string `json:"name"`
		Desc   string `json:"desc"`
		In     string `json:"in"`
		Out    string `json:"out"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.Checkpoint{
		Name:   params.Name,
		Desc:   params.Desc,
		In:     params.In,
		Out:    params.Out,
		Method: params.Method,
	}).Create(c); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("创建检查点成功"))
}

// 修改检查点
func EditCheckpoint(c *gin.Context) {
	params := &struct {
		ID     uint   `json:"id"`
		Method uint   `json:"method"`
		Name   string `json:"name"`
		Desc   string `json:"desc"`
		In     string `json:"in"`
		Out    string `json:"out"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.Checkpoint{
		Model: gorm.Model{
			ID: params.ID,
		},
		Name:   params.Name,
		Desc:   params.Desc,
		In:     params.In,
		Out:    params.Out,
		Method: params.Method,
	}).Edit(c); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("修改检查点"))
}

// 删除检查点
func DeleteCheckpoint(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if err := models.DeleteCheckpoint(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("检查点删除成功"))
}
