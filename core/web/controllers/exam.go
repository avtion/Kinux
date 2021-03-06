package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// list
func ListExams(c *gin.Context) {
	params := &struct {
		Page, Size int
		Ns         string
	}{
		Page: cast.ToInt(c.DefaultQuery("page", "1")),
		Size: cast.ToInt(c.DefaultQuery("size", "10")),
		Ns:   c.DefaultQuery("ns", ""),
	}
	data, err := models.ListExams(c, params.Ns, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译
	type resType struct {
		ID         uint          `json:"id"`
		Name       string        `json:"name"`
		Desc       string        `json:"desc"`
		Namespace  string        `json:"namespace"`
		Total      uint          `json:"total"`
		ForceOrder bool          `json:"force_order"`
		BeginAt    string        `json:"begin_at"`
		EndAt      string        `json:"end_at"`
		CreatedAt  string        `json:"created_at"`
		TimeLimit  time.Duration `json:"time_limit"`
	}

	res := make([]*resType, 0, len(data))
	for _, v := range data {
		res = append(res, &resType{
			ID:         v.ID,
			Name:       v.Name,
			Desc:       v.Desc,
			Namespace:  v.Name,
			Total:      v.Total,
			ForceOrder: v.ForceOrder,
			BeginAt:    v.BeginAt.Format("2006-01-02 15:04:05"),
			EndAt:      v.EndAt.Format("2006-01-02 15:04:05"),
			CreatedAt:  v.CreatedAt.Format("2006-01-02 15:04:05"),
			TimeLimit:  v.TimeLimit / time.Minute,
		})
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// count
func CountExams(c *gin.Context) {
	params := &struct {
		Ns string
	}{
		Ns: c.DefaultQuery("ns", ""),
	}
	res, err := models.CountExams(c, params.Ns)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// delete
func DeleteExam(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if err := models.DeleteExam(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("实验删除成功"))
}

// add
func AddExam(c *gin.Context) {
	params := &struct {
		Name       string `json:"name" binding:"required"`
		Desc       string `json:"desc"`
		Namespace  string `json:"namespace" binding:"required"`
		Total      uint   `json:"total" binding:"required"`
		BeginAt    int64  `json:"begin_at" binding:"required"`
		EndAt      int64  `json:"end_at" binding:"required"`
		ForceOrder bool   `json:"force_order"`
		TimeLimit  uint   `json:"time_limit"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.Exam{
		Name:       params.Name,
		Desc:       params.Desc,
		Namespace:  params.Namespace,
		Total:      params.Total,
		BeginAt:    time.Unix(params.BeginAt, 0),
		EndAt:      time.Unix(params.EndAt, 0),
		ForceOrder: params.ForceOrder,
		TimeLimit:  time.Duration(params.TimeLimit) * time.Minute,
	}).Create(c); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
}

// edit
func EditExam(c *gin.Context) {
	params := &struct {
		ID         uint   `json:"id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		Desc       string `json:"desc"`
		Namespace  string `json:"namespace" binding:"required"`
		Total      uint   `json:"total" binding:"required"`
		BeginAt    int64  `json:"begin_at" binding:"required"`
		EndAt      int64  `json:"end_at" binding:"required"`
		ForceOrder bool   `json:"force_order"`
		TimeLimit  uint   `json:"time_limit"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.Exam{
		Model:      gorm.Model{ID: params.ID},
		Name:       params.Name,
		Desc:       params.Desc,
		Namespace:  params.Namespace,
		Total:      params.Total,
		BeginAt:    time.Unix(params.BeginAt, 0),
		EndAt:      time.Unix(params.EndAt, 0),
		ForceOrder: params.ForceOrder,
		TimeLimit:  time.Duration(params.TimeLimit) * time.Minute,
	}).Save(c); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
}

// 获取已有考试的命名空间列表
func GetExamsNamespaces(c *gin.Context) {
	res, err := models.GetExamsNamespaces(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
