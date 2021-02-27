package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/tools/bytesconv"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

// 获取Deployment
func ListDeployment(c *gin.Context) {
	params := &struct {
		Name       string
		Page, Size int
	}{
		Name: c.DefaultQuery("name", ""),
		Page: cast.ToInt(c.DefaultQuery("page", "1")),
		Size: cast.ToInt(c.DefaultQuery("size", "10")),
	}
	data, err := models.ListDeployment(c, params.Name, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译结果
	type resType struct {
		ID        uint   ` json:"id"`
		Name      string `json:"name"`
		Raw       string `json:"raw"`
		CreatedAt string ` json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	res := make([]*resType, 0, len(data))
	for _, v := range data {
		res = append(res, &resType{
			ID:        v.ID,
			Name:      v.Name,
			Raw:       bytesconv.BytesToString(v.Raw),
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 添加Deployment
func AddDeployment(c *gin.Context) {
	params := &struct {
		Name string `json:"name"`
		Raw  string `json:"raw"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.AddDeployment(c, params.Name, bytesconv.StringToBytes(params.Raw)); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("Deployment创建成功"))
}

// 修改Deployment
func EditDeployment(c *gin.Context) {
	params := &struct {
		ID  uint   `json:"id"`
		Raw string `json:"raw"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.EditDeployment(c, params.ID, bytesconv.StringToBytes(params.Raw)); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("Deployment修改成功"))
}

// 删除Deployment
func DeleteDeployment(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if err := models.DeleteDeployment(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("Deployment删除成功"))
}

// 快速获取Deployment
func QuickListDeployment(c *gin.Context) {
	params := &struct {
		Name string
	}{
		Name: c.DefaultQuery("name", ""),
	}
	res, err := models.QuickListDeployment(c, params.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 统计Deployment
func CountDeployment(c *gin.Context) {
	params := &struct {
		Name string
	}{
		Name: c.DefaultQuery("name", ""),
	}
	res, err := models.CountDeployment(c, params.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
