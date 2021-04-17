package routers

import (
	"Kinux/core/web/controllers"
	"Kinux/core/web/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	addRouters(v2Router())
}

func v2Router() initFunc {
	return func(r *gin.Engine) {
		v2 := r.Group("v2")

		v2WithAuth := v2.Group("", middlewares.JsonWebTokenAuth)

		// 课程管理
		lesson := v2WithAuth.Group("/lesson")
		{
			lesson.GET("/", controllers.GetLessonsOptions)
			lesson.GET("/list", controllers.ListLessons)
			lesson.GET("/count", controllers.CountLessons)
			lesson.PUT("/", controllers.EditLesson)
			lesson.POST("/", controllers.AddLesson)
			lesson.DELETE("/:id/", controllers.DeleteLesson)
		}

		// 课程实验管理
		lm := v2WithAuth.Group("/lm")
		{
			lm.GET("/list", controllers.ListLessonMission)
			lm.GET("/count", controllers.CountLessonMission)
			lm.PUT("/", controllers.EditLessonMission)
			lm.POST("/", controllers.AddLessonMission)
			lm.DELETE("/:id/", controllers.DeleteLessonMission)
		}

		// 班级课程管理
		dl := v2WithAuth.Group("/dl")
		{
			dl.GET("/list", controllers.ListDepartmentLesson)
			dl.GET("/count", controllers.CountDepartmentLesson)
			dl.POST("/", controllers.AddDepartmentLesson)
			dl.DELETE("/:id/", controllers.DeleteDepartmentLesson)
		}

		ms := v2WithAuth.Group("/ms")
		{
			ms.GET("/", controllers.ListMissionsV2)
		}

		score := v2WithAuth.Group("/score")
		{
			score.GET("/mission/", controllers.GetMissionScore)                      // 查询实验成绩
			score.GET("/exam/", controllers.GetExamScore)                            // 查询考试成绩
			score.GET("/save/", controllers.SaveScoreForAdmin)                       // 成绩存档
			score.GET("/quick/", controllers.QuickScoreSaverForAdmin)                // 快速获取成绩存档选项
			score.GET("/save/:id/", controllers.GetScoreSaversForAdmin)              // 获取成绩存档
			score.GET("/excel/:dp/:lesson/:mode/:st/:target/", controllers.GetExcel) // 获取成绩Excel文件
		}

		counter := v2WithAuth.Group("/counter")
		{
			counter.GET("/", controllers.Counter) // 数据统计
		}
	}
}
