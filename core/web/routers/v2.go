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

	}
}
