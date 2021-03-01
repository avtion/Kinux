package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

// 任务查询
type MissionQuery struct {
	Name      string   `form:"name" `
	Namespace []string `form:"namespace"`
	Page      uint     `form:"page"`
	Size      uint     `form:"size"`
}

// 查询任务
func QueryMissions(c *gin.Context) {
	var query = new(MissionQuery)
	if err := c.ShouldBind(query); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ms, err := services.ListMissions(c, ac, query.Name, query.Namespace, int(query.Page), int(query.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(ms))
	return
}

// 创建新的任务
// TODO 前端测试
func NewMissionDeployment(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("任务id为空"))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err = services.AccountMissionOpera(c, ac, id, services.MissionCreate); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("任务创建成功"))
}

// 删除正在进行的任务
// TODO 前端测试
func DeleteMissionDeployment(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("任务id为空"))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// TODO 任务点检查

	if err = services.AccountMissionOpera(c, ac, id, services.MissionDelete); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("任务删除成功"))
}

// 获取任务允许的容器名列表
func ListMissionAllowedContainersNames(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("任务id为空"))
		return
	}
	names, err := services.ListMissionAllowedContainersNames(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(names))
}

// 获取任务的实验文档
func GetMissionGuide(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("任务id为空"))
		return
	}
	res, err := services.GetMissionGuide(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 获取任务数据
func ListMissions(c *gin.Context) {
	params := &struct {
		Page, Size int
		Name       string
		Namespace  []string
	}{
		Page:      cast.ToInt(c.DefaultQuery("page", "1")),
		Size:      cast.ToInt(c.DefaultQuery("size", "10")),
		Name:      c.DefaultQuery("name", ""),
		Namespace: c.QueryArray("ns[]"),
	}
	data, err := models.ListMissions(c, params.Name, params.Namespace, models.NewPageBuilder(params.Page, params.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译
	type resType struct {
		ID         uint     `json:"id"`
		Total      uint     `json:"total"`
		Name       string   `json:"name"`
		Desc       string   `json:"desc"`
		Namespace  string   `json:"namespace"`
		Containers []string `json:"containers"`

		Deployment    uint   `json:"deployment"`
		ExecContainer string `json:"exec_container"`
		Command       string `json:"command"`
	}
	var res = make([]*resType, 0, len(data))
	for _, v := range data {
		res = append(res, &resType{
			ID:        v.ID,
			Total:     v.Total,
			Name:      v.Name,
			Desc:      v.Desc,
			Namespace: v.Namespace,
			// 修复参数为空的情况下会返回[""]的情况
			Containers: func(cs string) []string {
				if cs != "" {
					return strings.Split(v.WhiteListContainer, ";")
				} else {
					return []string{}
				}
			}(v.WhiteListContainer),
			Deployment:    v.Deployment,
			ExecContainer: v.ExecContainer,
			Command:       v.Command,
		})
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 统计任务数量
func CountMissions(c *gin.Context) {
	params := &struct {
		Name      string
		Namespace []string
	}{
		Name:      c.DefaultQuery("name", ""),
		Namespace: c.QueryArray("ns"),
	}
	res, err := models.CountMissions(c, params.Name, params.Namespace)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 删除任务
func DeleteMission(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if err := models.DeleteMission(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("删除实验成功"))
}

// 获取所有任务的命名空间
func ListMissionNamespaces(c *gin.Context) {
	res, err := models.ListMissionNamespaces(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 编辑任务
func EditMission(c *gin.Context) {

}
