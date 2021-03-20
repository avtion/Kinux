package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取课程
func ListLessons(c *gin.Context) {
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
