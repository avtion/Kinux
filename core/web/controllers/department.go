package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

// 获取班级列表
func ListDepartments(c *gin.Context) {
	params := &struct {
		Page       int    `json:"page"`
		Size       int    `json:"size"`
		NameFilter string `json:"name_filter"`
	}{
		Page:       cast.ToInt(c.DefaultQuery("page", "1")),
		Size:       cast.ToInt(c.DefaultQuery("size", "10")),
		NameFilter: c.DefaultQuery("name_filter", ""),
	}
	dps, err := models.ListDepartments(c, params.NameFilter, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译一下结果
	type resType struct {
		ID        uint     `json:"id"`
		Name      string   `json:"name"`
		CreatAt   string   `json:"creat_at"`
		UpdatedAt string   `json:"updated_at"`
		Namespace []string `json:"namespace"`
	}
	res := make([]resType, 0, len(dps))
	for _, v := range dps {
		res = append(res, resType{
			ID:        v.ID,
			Name:      v.Name,
			CreatAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Namespace: v.GetNS(),
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 统计班级总数
func CountDepartments(c *gin.Context) {
	params := &struct {
		NameFilter string `json:"name_filter"`
	}{
		NameFilter: c.DefaultQuery("name_filter", ""),
	}
	res, err := models.CountDepartments(c, params.NameFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 新增班级
func AddDepartment(c *gin.Context) {
	params := new(struct {
		Name      string   `json:"name" binding:"required"`
		Namespace []string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if _, err := models.NewDepartment(c, params.Name, models.DepartmentNsOpt(params.Namespace...)); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("班级创建成功"))
}

// 删除班级
func DeleteDepartment(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("id为空"))
		return
	}
	if err := models.DeleteDepartment(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("班级删除成功"))
}

// 修改班级
func EditDepartment(c *gin.Context) {
	params := new(struct {
		ID        int      `json:"id" binding:"required"`
		Name      string   `json:"name" binding:"required"`
		Namespace []string `json:"namespace"`
	})
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.UpdateDepartment(c, params.ID,
		models.DepartmentNameOpt(params.Name),
		models.DepartmentNsOpt(params.Namespace...),
	); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("班级更新成功"))
}

// 快速返回班级数据
func QuickListDepartments(c *gin.Context) {
	res, err := models.QuickListDepartment(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
