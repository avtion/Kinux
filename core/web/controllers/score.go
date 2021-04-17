package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type missionScoreQuery struct {
	Department uint
	Lesson     uint
	Mission    uint
}

// GetMissionScore 获取实验成绩
func GetMissionScore(c *gin.Context) {
	params := &missionScoreQuery{
		Department: cast.ToUint(c.DefaultQuery("dp", "0")),
		Lesson:     cast.ToUint(c.DefaultQuery("lesson", "0")),
		Mission:    cast.ToUint(c.DefaultQuery("mission", "0")),
	}
	if params.Lesson == 0 || params.Mission == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(errors.New("目标课程和实验不能为空")))
		return
	}

	// 兼容管理员接口
	if params.Department != 0 {
		GetMissionScoreForAdmin(c, params)
		return
	}

	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 查询成绩
	res, err := services.GetMissionScore(c, ac.ID, params.Lesson, params.Mission, 0)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// GetMissionScoreForAdmin 管理员查询实验成绩
func GetMissionScoreForAdmin(c *gin.Context, params *missionScoreQuery) {
	res, err := services.GetMissionScoreForAdmin(c, params.Department, params.Lesson, params.Mission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

type examScoreParams struct {
	Department uint
	Lesson     uint
	Exam       uint
}

// GetExamScore 获取考试成绩
func GetExamScore(c *gin.Context) {
	params := &examScoreParams{
		Department: cast.ToUint(c.DefaultQuery("dp", "0")),
		Lesson:     cast.ToUint(c.DefaultQuery("lesson", "0")),
		Exam:       cast.ToUint(c.DefaultQuery("exam", "0")),
	}
	if params.Lesson == 0 || params.Exam == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(errors.New("目标课程和考试不能为空")))
		return
	}

	// 兼容管理员接口
	if params.Department != 0 {
		GetExamScoreForAdmin(c, params)
		return
	}

	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	res, err := services.GetExamScore(c, ac.ID, params.Lesson, params.Exam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// GetExamScoreForAdmin 管理员查询考试成绩
func GetExamScoreForAdmin(c *gin.Context, params *examScoreParams) {
	res, err := services.GetExamScoreForAdmin(c, params.Department, params.Lesson, params.Exam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// SaveScoreForAdmin 管理员存档成绩
func SaveScoreForAdmin(c *gin.Context) {
	params := &struct {
		Department uint
		Lesson     uint
		Target     uint
		ScoreType  models.ScoreSaverType
	}{
		Department: cast.ToUint(c.DefaultQuery("dp", "0")),
		Lesson:     cast.ToUint(c.DefaultQuery("lesson", "0")),
		Target:     cast.ToUint(c.DefaultQuery("target", "0")),
		ScoreType:  cast.ToUint(c.DefaultQuery("type", "0")),
	}
	if params.Department == 0 || params.Lesson == 0 || params.Target == 0 || params.ScoreType == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(errors.New("参数不能为空")))
		return
	}
	var rawData []byte

	switch params.ScoreType {
	case models.ScoreTypeMission:
		data, err := services.GetMissionScoreForAdmin(c, params.Department, params.Lesson, params.Target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
		rawData, err = jsoniter.Marshal(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
	case models.ScoreTypeExam:
		data, err := services.GetExamScoreForAdmin(c, params.Department, params.Lesson, params.Target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
		rawData, err = jsoniter.Marshal(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
	}
	if err := models.NewScoreSave(c, params.ScoreType, params.Target, rawData); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("成绩存档成功"))
}

// QuickScoreSaverForAdmin 获取选项
func QuickScoreSaverForAdmin(c *gin.Context) {
	params := &struct {
		Department uint
		Lesson     uint
		Page       int
		Size       int
		ScoreType  models.ScoreSaverType
		Name       string
	}{
		Department: cast.ToUint(c.DefaultQuery("dp", "0")),
		Lesson:     cast.ToUint(c.DefaultQuery("lesson", "0")),
		Page:       cast.ToInt(c.DefaultQuery("page", "0")),
		Size:       cast.ToInt(c.DefaultQuery("size", "0")),
		ScoreType:  cast.ToUint(c.DefaultQuery("type", "0")),
		Name:       c.DefaultQuery("name", ""),
	}
	if params.Department == 0 || params.Lesson == 0 || params.ScoreType == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(errors.New("参数不能为空")))
		return
	}
	data, err := models.ListScoreSave(c, params.ScoreType, models.NewPageBuilder(params.Page, params.Size).Build,
		func(db *gorm.DB) *gorm.DB {
			return db.Select("id, raw_id, raw_name, raw_created_at, created_at")
		},
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	type resType struct {
		ID           uint      `json:"id"`
		RawID        uint      `json:"raw_id"`
		RawName      string    `json:"raw_name"`
		RawCreatedAt time.Time `json:"raw_created_at"`
		CreatedAt    time.Time `json:"created_at"`
	}
	var res = make([]*resType, 0, len(data))
	for _, v := range data {
		res = append(res, &resType{
			ID:           v.ID,
			RawID:        v.RawID,
			RawName:      v.RawName,
			RawCreatedAt: v.RawCreatedAt,
			CreatedAt:    v.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// GetScoreSaversForAdmin 获取实验存档
func GetScoreSaversForAdmin(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(errors.New("参数不能为空")))
		return
	}
	var data = new(models.ScoresSaver)
	if err := models.GetGlobalDB().WithContext(c).Where("id = ?", id).First(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 反序列化数据
	var dataObj interface{}
	switch data.ScoreType {
	case models.ScoreTypeMission:
		dataObj = make([]*services.MissionScoreForAdmin, 0)
	case models.ScoreTypeExam:
		dataObj = make([]*services.ExamScoreForAdmin, 0)
	}
	if err := jsoniter.Unmarshal(data.Data, &dataObj); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	type resType struct {
		ID           uint                  `json:"id"`
		ScoreType    models.ScoreSaverType `json:"score_type"`
		RawID        uint                  `json:"raw_id"`         // 实验或考试的原ID
		RawName      string                `json:"raw_name"`       // 实验或者考试原名
		RawCreatedAt time.Time             `json:"raw_created_at"` // 实验或者考试的创建时间
		Data         interface{}           `json:"data"`
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(&resType{
		ID:           data.ID,
		ScoreType:    data.ScoreType,
		RawID:        data.RawID,
		RawName:      data.RawName,
		RawCreatedAt: data.RawCreatedAt,
		Data:         dataObj,
	}))
}

// GetExcel 获取成绩的Excel
func GetExcel(c *gin.Context) {
	params := &struct {
		Department uint
		Lesson     uint
		Mode       uint                  // 1-实时 2-存档
		ScoreType  models.ScoreSaverType // 1-实验 2-考试
		TargetID   uint                  // 目标ID
	}{
		Department: cast.ToUint(c.Param("dp")),
		Lesson:     cast.ToUint(c.Param("lesson")),
		Mode:       cast.ToUint(c.Param("mode")),
		ScoreType:  cast.ToUint(c.Param("st")),
		TargetID:   cast.ToUint(c.Param("target")),
	}

	if params.Department == 0 ||
		params.Lesson == 0 ||
		params.Mode == 0 ||
		params.ScoreType == 0 ||
		params.TargetID == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(errors.New("参数为空")))
		return
	}

	// 获取班级
	dp, err := models.GetDepartmentByID(c, params.Department)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 课程
	lesson, err := models.GetLesson(c, params.Lesson)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	var (
		f               *excelize.File
		fileNameBuilder = strings.Builder{}
	)

	// 生成文件名
	fileNameBuilder.WriteString(dp.Name)
	fileNameBuilder.WriteString("-" + lesson.Name)
	fileNameBuilder.WriteString("-" + cast.ToString(time.Now().Unix()))
	fileNameBuilder.WriteString(".xlsx")

	// 实时或者存档
	switch params.Mode {
	case 1:
		switch params.ScoreType {
		case models.ScoreTypeMission:
			res, err := services.GetMissionScoreForAdmin(c, params.Department, params.Lesson, params.TargetID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
				return
			}
			f, err = res.GetExcel()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
				return
			}
		case models.ScoreTypeExam:
			res, err := services.GetExamScoreForAdmin(c, params.Department, params.Lesson, params.TargetID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
				return
			}
			f, err = res.GetExcel()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
				return
			}
		}
	case 2:
		var data = new(models.ScoresSaver)
		if err := models.GetGlobalDB().WithContext(c).Where(
			"id = ?", params.TargetID).First(&data).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}

		// 反序列化数据
		switch data.ScoreType {
		case models.ScoreTypeMission:
			raw := make(services.MissionScoreForAdminSlice, 0)
			if err := jsoniter.Unmarshal(data.Data, &raw); err != nil {
				c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
				return
			}
			f, err = raw.GetExcel()
		case models.ScoreTypeExam:
			raw := make(services.ExamScoreForAdminSlice, 0)
			if err := jsoniter.Unmarshal(data.Data, &raw); err != nil {
				c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
				return
			}
			f, err = raw.GetExcel()
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
	}
	if f == nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("成绩Excel生成失败，清联系管理员"))
		return
	}

	// 创建缓冲流
	buff, err := f.WriteToBuffer()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	c.DataFromReader(http.StatusOK, int64(buff.Len()),
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		buff,
		map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`,
				fileNameBuilder.String()),
		},
	)
	return
}
