package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterResponseRoutes 注册答卷相关路由
func RegisterResponseRoutes(router *gin.Engine) {
	group := router.Group("/api/survey")
	{
		group.POST("/:SurveyID/GetOption", controllers.GetOptionCount)
		group.POST("/:SurveyID/GetText", controllers.GetTextFillinData)
		group.POST("/:SurveyID/GetNum", controllers.GetNumFillinData)
		group.GET("/:SurveyID", controllers.GetSurveyResponsesHandler)
	}
}
