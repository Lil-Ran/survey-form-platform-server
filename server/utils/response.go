package utils

import (
	"github.com/gin-gonic/gin"
)

func JSONResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := gin.H{
		"code":    statusCode,
		"message": message,
	}

	// 如果 data 不为 nil，则直接添加到顶层
	if data != nil {
		for key, value := range data.(gin.H) {
			response[key] = value
		}
	}

	c.JSON(statusCode, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	JSONResponse(c, statusCode, message, nil)
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	JSONResponse(c, statusCode, message, data)
}
