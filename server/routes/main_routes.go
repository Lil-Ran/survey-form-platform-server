package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		RegisterAccountRoutes(apiGroup) // 注册账户相关路由
		RegisterCookieRoutes(apiGroup)  // 注册 Cookie 相关路由
		RegisterSurveyRoutes(apiGroup)  // 注册问卷相关路由
		// RegisterQuestionEditRoutes(apiGroup) // 注册问题路由
	}
}
