package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strings"
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
		ID            uint   `json:"id"`
		Name          string `json:"name"`
		Desc          string `json:"desc"`
		Total         uint   `json:"total"`
		ForceOrder    bool   `json:"force_order"`
		BeginAt       string `json:"begin_at"`
		EndAt         string `json:"end_at"`
		CreatedAt     string `json:"created_at"`
		TimeLimit     string `json:"time_limit"`
		BeginAtUnix   int64  `json:"begin_at_unix"`
		EndAtUnix     int64  `json:"end_at_unix"`
		CreatedAtUnix int64  `json:"created_at_unix"`
		TimeLimitUnix int64  `json:"time_limit_unix"`
	}

	res := make([]*resType, 0, len(data))
	for _, v := range data {
		res = append(res, &resType{
			ID:            v.ID,
			Name:          v.Name,
			Desc:          v.Desc,
			Total:         v.Total,
			ForceOrder:    v.ForceOrder,
			BeginAt:       v.BeginAt.Format("2006-01-02 15:04:05"),
			EndAt:         v.EndAt.Format("2006-01-02 15:04:05"),
			CreatedAt:     v.CreatedAt.Format("2006-01-02 15:04:05"),
			TimeLimit:     v.TimeLimit.String(),
			BeginAtUnix:   v.BeginAt.Unix(),
			EndAtUnix:     v.EndAt.Unix(),
			CreatedAtUnix: v.CreatedAt.Unix(),
			TimeLimitUnix: int64(v.TimeLimit / time.Second),
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
		TimeLimit  int64  `json:"time_limit"`
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
		BeginAt:    time.Unix(0, params.BeginAt*int64(time.Millisecond)),
		EndAt:      time.Unix(0, params.EndAt*int64(time.Millisecond)),
		ForceOrder: params.ForceOrder,
		TimeLimit:  time.Duration(params.TimeLimit) * time.Second,
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
		TimeLimit  int64  `json:"time_limit"`
		Lesson     uint   `json:"lesson"`
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
		BeginAt:    time.Unix(0, params.BeginAt*int64(time.Millisecond)),
		EndAt:      time.Unix(0, params.EndAt*int64(time.Millisecond)),
		ForceOrder: params.ForceOrder,
		TimeLimit:  time.Duration(params.TimeLimit) * time.Second,
		Lesson:     params.Lesson,
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
		ID        uint   `json:"id"`
		Exam      uint   `json:"exam"`
		Mission   string `json:"mission"`
		Percent   uint   `json:"percent"`
		Priority  int    `json:"priority"`
		MissionID uint   `json:"mission_id"`
	}
	var res = make([]*resType, 0, len(data))

	// 查找实验的名称
	var missionIDs = make([]uint, 0, len(data))
	for _, v := range data {
		missionIDs = append(missionIDs, v.ID)
	}
	nameMapper, _ := models.GetMissionsNameMapper(c, missionIDs...)

	for _, v := range data {
		missionName, _ := nameMapper[v.Mission]
		res = append(res, &resType{
			ID:        v.ID,
			Exam:      v.Exam,
			Mission:   missionName,
			Percent:   v.Percent,
			Priority:  v.Priority,
			MissionID: v.Mission,
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
		Priority int  `json:"priority"`
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
		if strings.Contains(err.Error(), "UNIQUE") {
			err = errors.New("该考试已经存在当前实验")
		}
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
		Priority int  `json:"priority"`
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

// 根据班级获取对应的考试
func ListExamByDepartment(c *gin.Context) {
	params := &struct {
		Dp uint
	}{
		Dp: cast.ToUint(c.DefaultQuery("dp", "0")),
	}

	// 查看班级的课程
	type LessonType struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	var lessons = make([]*LessonType, 0)
	if err := models.GetGlobalDB().WithContext(c).Model(new(models.LessonDepartment)).Joins(
		"left join lessons ON lesson_departments.lesson = lessons.id",
	).Select("lesson_departments.id as id, lessons.desc as desc, lessons.name as name").Where(
		"lesson_departments.department = ?", params.Dp).Scan(&lessons).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	var lessonIDs = make([]uint, 0, len(lessons))
	var lessonMapping = make(map[uint]*LessonType, len(lessons))
	for k, v := range lessons {
		lessonIDs = append(lessonIDs, v.ID)
		lessonMapping[v.ID] = lessons[k]
	}

	data, err := models.ListExams(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("lesson IN ?", lessonIDs)
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译
	type resType struct {
		ID            uint   `json:"id"`
		Name          string `json:"name"`
		Desc          string `json:"desc"`
		Total         uint   `json:"total"`
		ForceOrder    bool   `json:"force_order"`
		BeginAt       string `json:"begin_at"`
		EndAt         string `json:"end_at"`
		CreatedAt     string `json:"created_at"`
		TimeLimit     string `json:"time_limit"`
		BeginAtUnix   int64  `json:"begin_at_unix"`
		EndAtUnix     int64  `json:"end_at_unix"`
		CreatedAtUnix int64  `json:"created_at_unix"`
		TimeLimitUnix int64  `json:"time_limit_unix"`
		Lesson        uint   `json:"lesson"`
		LessonName    string `json:"lesson_name"`
		LessonDesc    string `json:"lesson_desc"`

		// 考试状态
		ExamStatus services.ExamStatus
	}

	// 为了获取考试状态
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	res := make([]*resType, 0, len(data))
	for _, v := range data {
		l, _ := lessonMapping[v.Lesson]
		res = append(res, &resType{
			ID:            v.ID,
			Name:          v.Name,
			Desc:          v.Desc,
			Total:         v.Total,
			ForceOrder:    v.ForceOrder,
			BeginAt:       v.BeginAt.Format("2006-01-02 15:04:05"),
			EndAt:         v.EndAt.Format("2006-01-02 15:04:05"),
			CreatedAt:     v.CreatedAt.Format("2006-01-02 15:04:05"),
			TimeLimit:     v.TimeLimit.String(),
			BeginAtUnix:   v.BeginAt.Unix(),
			EndAtUnix:     v.EndAt.Unix(),
			CreatedAtUnix: v.CreatedAt.Unix(),
			TimeLimitUnix: int64(v.TimeLimit / time.Second),
			Lesson:        v.Lesson,
			LessonName:    l.Name,
			LessonDesc:    l.Desc,

			ExamStatus: services.GetExamStatus(c, ac.ID, v.ID),
		})
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 确定考试状态（用户在开始考试之前先检查全局考试状态）
func CheckinExamStatus(c *gin.Context) {
	// 检查用户是否处于考试状态
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	_eWatcher, isExist := services.ExamWatchers.Load(ac)
	if isExist {
		ew := _eWatcher.(*services.ExamWatcher)
		c.JSON(http.StatusOK, msg.BuildSuccess(services.NewExamRunningInfo(ew)))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(nil))
}

// 开始考试
func StartExam(c *gin.Context) {
	params := &struct {
		ExamID uint
	}{
		ExamID: cast.ToUint(c.DefaultQuery("exam", "0")),
	}
	if params.ExamID == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("考试ID为空"))
		return
	}

	// 检查用户是否处于考试状态
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 检查考试状态
	switch examStatus := services.GetExamStatus(c, ac.ID, params.ExamID); examStatus {
	case services.ESNotStart:
		// 启动监控者
		if err = services.NewExamWatcher(c, ac.ID, params.ExamID); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
		c.JSON(http.StatusOK, msg.BuildFailed("考试开始成功"))
	case services.ESRunning:
		c.JSON(http.StatusOK, msg.BuildSuccess("考试已开始"))
	case services.ESFinish:
		c.JSON(http.StatusOK, msg.BuildFailed("考试已经结束"))
	}
	return
}
