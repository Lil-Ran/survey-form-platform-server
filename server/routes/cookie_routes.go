package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCookieRoutes(router *gin.RouterGroup) {
	cookieGroup := router.Group("/cookie")
	{
		cookieGroup.GET("/set", controllers.SetCookie)       // 设置 Cookie
		cookieGroup.GET("/get", controllers.GetCookie)       // 获取 Cookie
		cookieGroup.GET("/delete", controllers.DeleteCookie) // 删除 Cookie
	}
}
