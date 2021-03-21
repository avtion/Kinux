package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
)

// 获取课程
func GetLessonsOptions(c *gin.Context) {
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	dp, err := ac.GetDepartment(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	lessons, err := dp.GetLessons(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	type resType struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	var res = make([]*resType, 0, len(lessons))
	for _, v := range lessons {
		res = append(res, &resType{
			ID:   v.ID,
			Name: v.Name,
			Desc: v.Desc,
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 增加课程
func AddLesson(c *gin.Context) {
	params := &struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.GetGlobalDB().WithContext(c).Create(&models.Lesson{
		Name: params.Name,
		Desc: params.Desc,
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("课程创建成功"))
}

// 获取课程
func ListLessons(c *gin.Context) {
	params := &struct {
		Page, Size int
	}{
		Page: cast.ToInt(c.DefaultQuery("page", "1")),
		Size: cast.ToInt(c.DefaultQuery("size", "10")),
	}
	selectRes := make([]*models.Lesson, 0, params.Size)
	if err := models.GetGlobalDB().WithContext(c).Model(new(models.Lesson)).Scopes(
		models.NewPageBuilder(params.Page, params.Size).Build,
		func(db *gorm.DB) *gorm.DB {
			if params.Name == "" {
				return db
			}
			return db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", params.Name))
		},
	).Find(&selectRes).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	type resType struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	var res = make([]*resType, 0, len(selectRes))
	for _, v := range selectRes {
		res = append(res, &resType{
			ID:   v.ID,
			Name: v.Name,
			Desc: v.Desc,
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 修改课程
func EditLesson(c *gin.Context) {
	params := &struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.GetGlobalDB().WithContext(c).Save(&models.Lesson{
		Model: gorm.Model{
			ID: params.ID,
		},
		Name: params.Name,
		Desc: params.Desc,
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("课程修改成功"))
}

// 删除课程
func DeleteLesson(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("id为空"))
		return
	}
	if err := models.GetGlobalDB().WithContext(c).Unscoped().Delete(&models.Lesson{
		Model: gorm.Model{ID: id},
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("课程删除成功"))
}

// 获取课程数量
func CountLessons(c *gin.Context) {
	var res int64
	if err := models.GetGlobalDB().WithContext(c).Model(new(models.Lesson)).Count(&res).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 增加课程实验
func AddLessonMission(c *gin.Context) {
	params := &struct {
		Lesson   uint `json:"lesson" binding:"gt=0"`
		Mission  uint `json:"mission" binding:"gt=0"`
		Priority uint `json:"priority"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.GetGlobalDB().WithContext(c).Create(&models.LessonMission{
		Model:    gorm.Model{},
		Lesson:   params.Lesson,
		Mission:  params.Mission,
		Priority: params.Priority,
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("课程实验创建成功"))
}

// 修改课程实验
func EditLessonMission(c *gin.Context) {
	params := &struct {
		ID       uint `json:"id"  binding:"gt=0"`
		Lesson   uint `json:"lesson"  binding:"gt=0"`
		Mission  uint `json:"mission"  binding:"gt=0"`
		Priority uint `json:"priority"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err := models.GetGlobalDB().WithContext(c).Save(&models.LessonMission{
		Model: gorm.Model{
			ID: params.ID,
		},
		Lesson:   params.Lesson,
		Mission:  params.Mission,
		Priority: params.Priority,
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("课程实验修改成功"))
}

// 获取课程实验
func ListLessonMission(c *gin.Context) {
	params := &struct {
		Lesson     uint
		Page, Size int
	}{
		Page:   cast.ToInt(c.DefaultQuery("page", "1")),
		Size:   cast.ToInt(c.DefaultQuery("size", "10")),
		Lesson: cast.ToUint(c.DefaultQuery("lesson", "0")),
	}
	if params.Lesson == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("课程不能为空"))
		return
	}
	type resType struct {
		ID          uint   `json:"id"`
		MissionID   uint   `json:"mission_id"`
		MissionName string `json:"mission_name"`
		MissionDesc string `json:"mission_desc"`
		Priority    uint   `json:"priority"`
	}
	var res = make([]*resType, 0, params.Size)
	if err := models.GetGlobalDB().WithContext(c).Model(new(models.LessonMission)).Joins(
		"left join missions ON lesson_missions.mission = missions.id").Select(
		"lesson_missions.id AS id, missions.id AS mission_id, "+
			"missions.desc AS mission_desc, missions.name AS mission_name, priority",
	).Scopes(
		models.NewPageBuilder(params.Page, params.Size).Build,
		models.ScopeLessonMissionOrderByOrder(),
	).Where("lesson_missions.lesson = ?", params.Lesson).Scan(&res).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 统计课程实验
func CountLessonMission(c *gin.Context) {
	params := &struct {
		Lesson uint
	}{
		Lesson: cast.ToUint(c.DefaultQuery("lesson", "0")),
	}
	if params.Lesson == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("课程不能为空"))
		return
	}
	var res int64
	if err := models.GetGlobalDB().WithContext(c).Model(
		new(models.LessonMission)).Where("lesson = ?", params.Lesson).Count(&res).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 删除课程实验
func DeleteLessonMission(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("id为空"))
		return
	}
	if err := models.GetGlobalDB().WithContext(c).Unscoped().Delete(&models.LessonMission{
		Model: gorm.Model{ID: id},
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("实验删除成功"))
}
