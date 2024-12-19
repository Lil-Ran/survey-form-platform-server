package controllers

import (
	"net/http"
	"server/services"

	"github.com/gin-gonic/gin"
)

// 设置 Cookie
func SetCookie(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id is required"})
		return
	}
	services.SetCookie(c, userID)
}

// 获取 Cookie
func GetCookie(c *gin.Context) {
	claims, err := services.GetCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

// 删除 Cookie
func DeleteCookie(c *gin.Context) {
	services.DeleteCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Cookie has been deleted"})
}
