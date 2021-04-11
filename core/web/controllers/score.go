package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
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
