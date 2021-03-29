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
			ws.GET("/list/", controllers.ListLiveWebsocket)
			ws.POST("/msg/", controllers.SendMessageToTargetWs)
		}

		// 挂载JWT鉴权中间件
		v1 = v1.Group("", middlewares.JsonWebTokenAuth)

		// 任务相关
		ms := v1.Group("/mission")
		{
			ms.GET("/", controllers.QueryMissions)                                // 批量查询任务
			ms.POST("/op/:id/", controllers.NewMissionDeployment)                 // 创建任务
			ms.DELETE("/op/:id/", controllers.DeleteMissionDeployment)            // 删除任务
			ms.GET("/guide/:id/", controllers.GetMissionGuide)                    // 获取任务的实验文档
			ms.PUT("/guide/", controllers.UpdateMissionGuide)                     // 更新任务实验文档
			ms.GET("/cnames/:id/", controllers.ListMissionAllowedContainersNames) // 获取任务允许的容器名列表
			ms.GET("/list/", controllers.ListMissions)                            // list
			ms.GET("/count/", controllers.CountMissions)                          // count
			ms.POST("/model/", controllers.AddMission)                            // add
			ms.PUT("/model/", controllers.EditMission)                            // edit
			ms.DELETE("/delete/:id/", controllers.DeleteMission)                  // delete
			ms.GET("/ns/", controllers.ListMissionNamespaces)                     // 获取所有任务的命名空间
		}

		// 用户账号相关
		ac := v1.Group("/account")
		{
			ac.PUT("/avatar", controllers.UpdateAccountAvatarSeed) // 更新用户头像种子
			ac.POST("/pw", controllers.UpdatePassword)             // 更新用户密码
			ac.GET("/", controllers.ListAccounts)                  // list
			ac.GET("/count/", controllers.CountAccounts)           // count
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

		// 配置相关
		deployment := v1.Group("/deployment")
		{
			deployment.GET("/", controllers.ListDeployment)                          // list
			deployment.GET("/count/", controllers.CountDeployment)                   // count
			deployment.GET("/quick/", controllers.QuickListDeployment)               // option quick
			deployment.PUT("/", controllers.EditDeployment)                          // edit
			deployment.POST("/", controllers.AddDeployment)                          // add
			deployment.DELETE("/:id/", controllers.DeleteDeployment)                 // delete
			deployment.GET("/containers/:id/", controllers.ListDeploymentContainers) // 获取配置的所有容器名
		}

		// 成绩相关
		//score := v1.Group("/score")
		{
			//score.GET("/:type/", controllers.ListScores)         // list
			//score.DELETE("/:type/:id/", controllers.DeleteScore) // delete
			//score.GET("/ms/", controllers.ListMissionScore) // list
		}

		// 角色相关
		role := v1.Group("/role")
		{
			role.GET("/quick/", controllers.QuickListRoles) // quick
		}

		// 检查点相关
		checkpoint := v1.Group("/cp")
		{
			checkpoint.GET("/", controllers.ListCheckpoints)            // list
			checkpoint.POST("/", controllers.AddCheckpoint)             // add
			checkpoint.PUT("/", controllers.EditCheckpoint)             // edit
			checkpoint.DELETE("/:id/", controllers.DeleteCheckpoint)    // delete
			checkpoint.GET("/count/", controllers.CountCheckpoints)     // count
			checkpoint.GET("/quick/", controllers.QuickListCheckpoints) // options quick
		}

		// 实验监测点相关
		missionCheckpoint := v1.Group("/mcp")
		{
			missionCheckpoint.GET("/", controllers.ListMissionCheckpoints)                   // list
			missionCheckpoint.GET("/count/", controllers.CountMissionCheckpoints)            // count
			missionCheckpoint.GET("/percent/:id/", controllers.GetMissionCheckpointsPercent) // option percent
			missionCheckpoint.POST("/", controllers.AddMissionCheckpoint)                    // add
			missionCheckpoint.PUT("/", controllers.EditMissionCheckpoint)                    // edit
			missionCheckpoint.DELETE("/:id/", controllers.DeleteMissionCheckpoint)           // delete
		}

		// 考试相关
		exam := v1.Group("/exam")
		{
			exam.GET("/list/", controllers.ListExams)          // list
			exam.GET("/count/", controllers.CountExams)        // count
			exam.DELETE("/:id/", controllers.DeleteExam)       // delete
			exam.POST("/", controllers.AddExam)                // add
			exam.PUT("/", controllers.EditExam)                // edit
			exam.GET("/dp/", controllers.ListExamByDepartment) // list For user
		}

		// 考试实验相关
		examMission := v1.Group("/em")
		{
			examMission.GET("/", controllers.ListExamMissions)                      // list
			examMission.GET("/count/", controllers.CountExamMissions)               // count
			examMission.DELETE("/:id/", controllers.DeleteExamMission)              // delete
			examMission.POST("/", controllers.AddExamMission)                       // add
			examMission.PUT("/", controllers.EditExamMission)                       // edit
			examMission.GET("/percent/:id/", controllers.GetExamMissionUsedPercent) // option percent
		}
	}
}
