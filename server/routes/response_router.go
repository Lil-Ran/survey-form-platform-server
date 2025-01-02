package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterResponseRoutes 注册答卷相关路由
func RegisterResponseRoutes(group *gin.RouterGroup) {
	surveyGroup := group.Group("/survey")
	surveyGroup.POST("/:SurveyID/GetOption", controllers.GetOptionCount)
	surveyGroup.POST("/:SurveyID/GetText", controllers.GetTextFillinData)
	surveyGroup.POST("/:SurveyID/GetNum", controllers.GetNumFillinData)
	surveyGroup.GET("/:SurveyID", controllers.GetSurveyResponsesHandler)
}
