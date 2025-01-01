package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterSurveyRoutes 注册问卷编辑相关路由
func RegisterQuestionEditRoutes(router *gin.RouterGroup) {
	editGroup := router.Group("/edit")
	{
		editGroup.GET("/:surveyId/meta", controllers.GetSurveyMetaController)
		editGroup.GET("/:surveyId/questions", controllers.GetSurveyQuestionsController)
		editGroup.POST("/:surveyId/qedit", controllers.SaveSurveyEditController)
		editGroup.DELETE("/:surveyId/delete", controllers.DeleteSurveyController)
	}
}
