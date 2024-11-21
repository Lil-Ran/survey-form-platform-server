package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSurveyRoutes(router *gin.RouterGroup) {
	surveyGroup := router.Group("/survey")
	{
		surveyGroup.POST("/create", controllers.CreateSurvey)       // 问卷创建
		surveyGroup.POST("/switch", controllers.UpdateSurveyStatus) // 修改问卷状态
		surveyGroup.POST("/copy", controllers.CopySurvey)           // 问卷复制
		surveyGroup.GET("", controllers.ListSurveys)                // 获取问卷列表
	}
}
