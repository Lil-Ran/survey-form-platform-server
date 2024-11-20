package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

// 注册账户相关路由
func RegisterAccountRoutes(router *gin.Engine) {
	accountGroup := router.Group("/account")
	{
		// 用户注册路由
		accountGroup.POST("/register", controllers.RegisterUser)
		// 请求邮箱验证码路由
		accountGroup.POST("/request-email-code", controllers.RequestEmailCode)
		// 用户登录路由
		accountGroup.POST("/login", controllers.LoginUser)
		// 用户退出路由
		accountGroup.POST("/logout", controllers.LogoutUser)
		// 用户重置密码路由
		accountGroup.POST("/reset-password", controllers.ResetPassword)
	}
}
