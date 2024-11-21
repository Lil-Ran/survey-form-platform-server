package utils

import (
	"github.com/gin-gonic/gin"
)

// JSONResponse 统一 JSON 响应
func JSONResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	JSONResponse(c, statusCode, message, nil)
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	JSONResponse(c, statusCode, message, data)
}
