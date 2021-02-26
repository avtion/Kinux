package routers

import (
	"Kinux/core/web/controllers"
	"Kinux/core/web/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	addRouters(v1Routers())
}

func v1Routers() initFunc {
	return func(r *gin.Engine) {
		v1 := r.Group("v1")

		// 用户登陆接口不使用鉴权中间件
		v1.POST("/account/login", controllers.LoginAccount)

		// WebSocket需要手动进行鉴权
		ws := v1.Group("/ws")
		{
			ws.GET("/", controllers.WebSocketHandlerV1)
		}

		// 挂载JWT鉴权中间件
		v1 = v1.Group("", middlewares.JsonWebTokenAuth)

		// 任务相关
		ms := v1.Group("/mission")
		{
			ms.GET("/", controllers.QueryMissions)                                // 批量查询任务
			ms.POST("/op/:id/", controllers.NewMission)                           // 创建任务
			ms.DELETE("/op/:id/", controllers.DeleteMission)                      // 删除任务
			ms.GET("/guide/:id/", controllers.GetMissionGuide)                    // 获取任务的实验文档
			ms.GET("/cnames/:id/", controllers.ListMissionAllowedContainersNames) // 获取任务允许的容器名列表
		}

		// 用户账号相关
		ac := v1.Group("/account")
		{
			ac.PUT("/avatar", controllers.UpdateAccountAvatarSeed) // 更新用户头像种子
			ac.POST("/pw", controllers.UpdatePassword)             // 更新用户密码
			ac.GET("/", controllers.ListAccounts)                  // list
			ac.GET("/count/", controllers.CountAccounts)           // quick
			ac.POST("/", controllers.AddAccount)                   // add
			ac.PUT("/", controllers.EditAccount)                   // edit
			ac.DELETE("/:id/", controllers.DeleteAccount)          // delete
		}

		// 部门（班级）相关
		department := v1.Group("/department")
		{
			department.GET("/", controllers.ListDepartments)            // list
			department.POST("/", controllers.AddDepartment)             // add
			department.PUT("/", controllers.EditDepartment)             // edit
			department.DELETE("/:id/", controllers.DeleteDepartment)    // delete
			department.GET("/count/", controllers.CountDepartments)     // count
			department.GET("/quick/", controllers.QuickListDepartments) // options quick
		}
	}
}
