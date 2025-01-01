package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterResponseRoutes 注册答卷相关路由
func RegisterRespondentRoutes(apiGroup *gin.RouterGroup) {
	responseRoutes := apiGroup.Group("/respondent")
	{
		responseRoutes.GET("/:surveyId/questions", controllers.GetRespondentQuestionsController)
		responseRoutes.POST("/:surveyId/submit", controllers.SubmitSurveyResponseController)
	}
}
