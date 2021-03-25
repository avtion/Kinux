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
		Lesson     uint
	}{
		Page:   cast.ToInt(c.DefaultQuery("page", "1")),
		Size:   cast.ToInt(c.DefaultQuery("size", "10")),
		Lesson: cast.ToUint(c.DefaultQuery("lesson", "0")),
	}
	data, err := models.ListExams(c, models.NewPageBuilder(params.Page, params.Size).Build, func(db *gorm.DB) *gorm.DB {
		return db.Where("lesson = ?", params.Lesson)
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译
	type resType struct {
		ID         uint          `json:"id"`
		Name       string        `json:"name"`
		Desc       string        `json:"desc"`
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
		Lesson uint
	}{
		Lesson: cast.ToUint(c.DefaultQuery("lesson", "0")),
	}
	res, err := models.CountExams(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("lesson = ?", params.Lesson)
	})
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
		Total      uint   `json:"total" binding:"required"`
		BeginAt    int64  `json:"begin_at" binding:"required"`
		EndAt      int64  `json:"end_at" binding:"required"`
		ForceOrder bool   `json:"force_order"`
		TimeLimit  uint   `json:"time_limit"`
		Lesson     uint   `json:"lesson" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.Exam{
		Name:       params.Name,
		Desc:       params.Desc,
		Total:      params.Total,
		BeginAt:    time.Unix(params.BeginAt, 0),
		EndAt:      time.Unix(params.EndAt, 0),
		ForceOrder: params.ForceOrder,
		TimeLimit:  time.Duration(params.TimeLimit) * time.Minute,
		Lesson:     params.Lesson,
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

// list
func ListExamMissions(c *gin.Context) {
	params := &struct {
		Page, Size int
		Exam       uint
		Mission    uint
	}{
		Page:    cast.ToInt(c.DefaultQuery("page", "1")),
		Size:    cast.ToInt(c.DefaultQuery("size", "10")),
		Exam:    cast.ToUint(c.DefaultQuery("exam", "")),
		Mission: cast.ToUint(c.DefaultQuery("mission", "")),
	}
	data, err := models.ListExamMissions(c, params.Exam, params.Mission, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译一下
	type resType struct {
		ID       uint `json:"id"`
		Exam     uint `json:"exam"`
		Mission  uint `json:"mission"`
		Percent  uint `json:"percent"`
		Priority int  `json:"priority"`
	}
	var res = make([]*resType, 0, len(data))

	for _, v := range data {
		res = append(res, &resType{
			ID:       v.ID,
			Exam:     v.Exam,
			Mission:  v.Mission,
			Percent:  v.Percent,
			Priority: v.Priority,
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// count
func CountExamMissions(c *gin.Context) {
	params := &struct {
		Exam    uint
		Mission uint
	}{
		Exam:    cast.ToUint(c.DefaultQuery("exam", "")),
		Mission: cast.ToUint(c.DefaultQuery("mission", "")),
	}
	res, err := models.CountExamMissions(c, params.Exam, params.Mission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// delete
func DeleteExamMission(c *gin.Context) {
	if err := models.DeleteExamMission(c, cast.ToUint(c.Param("id"))); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("考试实验删除成功"))
}

// add
func AddExamMission(c *gin.Context) {
	params := &struct {
		Exam     uint `json:"exam" binding:"required"`
		Mission  uint `json:"mission" binding:"required"`
		Percent  uint `json:"percent" binding:"required"`
		Priority int  `json:"priority" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.ExamMissions{
		Exam:     params.Exam,
		Mission:  params.Mission,
		Percent:  params.Percent,
		Priority: params.Priority,
	}).Create(c); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("创建成功"))
}

// edit
func EditExamMission(c *gin.Context) {
	params := &struct {
		ID       uint `json:"id" binding:"required"`
		Exam     uint `json:"exam" binding:"required"`
		Mission  uint `json:"mission" binding:"required"`
		Percent  uint `json:"percent" binding:"required"`
		Priority int  `json:"priority" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := (&models.ExamMissions{
		Model:    gorm.Model{ID: params.ID},
		Exam:     params.Exam,
		Mission:  params.Mission,
		Percent:  params.Percent,
		Priority: params.Priority,
	}).Save(c); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("修改成功"))
}

// 获取考试实验已经使用的比例
func GetExamMissionUsedPercent(c *gin.Context) {
	res, err := models.GetExamMissionsUsedPercent(c, cast.ToUint(c.Param("id")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
