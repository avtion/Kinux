package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

// 获取成绩
func ListScores(c *gin.Context) {
	scoreType := cast.ToInt(c.Param("type"))
	sc := models.NewScore(scoreType)
	if sc == nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("scoreType错误"))
		return
	}

	params := &struct {
		Page, Size int
	}{
		Page: cast.ToInt(c.DefaultQuery("page", "1")),
		Size: cast.ToInt(c.DefaultQuery("size", "10")),
	}
	res, err := sc.ListScores(c, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("scoreType错误"))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 删除成绩
func DeleteScore(c *gin.Context) {
	scoreType := cast.ToInt(c.Param("type"))
	id := cast.ToUint(c.Param("id"))

	sc := models.NewScore(scoreType)
	if sc == nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("scoreType错误"))
		return
	}
	if err := sc.DeleteScore(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("成绩删除成功"))
}
