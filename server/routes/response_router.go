package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterResponseRoutes 注册答卷相关路由
func RegisterResponseRoutes(router *gin.Engine) {
	group := router.Group("/api/survey")
	{
		group.POST("/:SurveyID", controllers.GetOptionCount)
		group.POST("/:SurveyID", controllers.GetTextFillinData)
		group.POST("/:SurveyID", controllers.GetNumFillinData)
		group.GET("/:SurveyID", controllers.GetSurveyResponsesHandler)
	}
}
