package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 快速获取角色选项
func QuickListRoles(c *gin.Context) {
	type resType struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	res := make([]*resType, 0, len(models.RoleArray))
	for _, v := range models.RoleArray {
		if v == models.RoleAnonymous || v == models.RoleAdmin {
			continue
		}
		res = append(res, &resType{
			ID:   v,
			Name: models.RoleTranslator(v),
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
