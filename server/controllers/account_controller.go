package controllers

import (
	"net/http"
	"server/services"

	"github.com/gin-gonic/gin"
)

// RegisterUser 处理用户注册请求
func RegisterUser(c *gin.Context) {
	// 定义请求结构体，包含用户名、密码、邮箱和邮箱验证码
	var request struct {
		UserName  string `json:"userName" binding:"required,min=4,max=64"` // 用户名，必填，长度4-64
		Password  string `json:"password" binding:"required,min=8"`        // 密码，必填，最小长度8
		Email     string `json:"email" binding:"required,email"`           // 邮箱，必填，必须是有效邮箱格式
		EmailCode string `json:"emailCode" binding:"required,len=6"`       // 邮箱验证码，必填，长度6，必须是数字
	}

	// 绑定JSON请求体到结构体，并验证参数
	if err := c.ShouldBindJSON(&request); err != nil {
		// 如果参数无效，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request parameters", "code": 400})
		return
	}

	// 调用服务层的RegisterUser函数进行注册
	if err := services.RegisterUser(request.UserName, request.Password, request.Email, request.EmailCode); err != nil {
		// 如果注册失败，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	// 注册成功，返回200状态码
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "code": 200})
}

// RequestEmailCode 处理请求邮箱验证码
func RequestEmailCode(c *gin.Context) {
	// 定义请求结构体，包含邮箱
	var request struct {
		Email string `json:"email" binding:"required,email"` // 邮箱，必填，必须是有效邮箱格式
	}

	// 绑定JSON请求体到结构体，并验证参数
	if err := c.ShouldBindJSON(&request); err != nil {
		// 如果参数无效，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request parameters", "code": 400})
		return
	}

	// 调用服务层的SendEmailCode函数发送验证码
	if err := services.SendEmailCode(request.Email); err != nil {
		// 如果发送失败，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	// 发送成功，返回200状态码
	c.JSON(http.StatusOK, gin.H{"message": "Email code sent", "code": 200})
}

// LoginUser 用户登录
func LoginUser(c *gin.Context) {
	// 定义请求结构体，包含用户名或邮箱和密码
	var request struct {
		UserNameOrEmail string `json:"userName" binding:"required,min=4,max=64"` // 用户名或邮箱，必填，长度4-64
		Password        string `json:"password" binding:"required,min=8"`        // 密码，必填，最小长度8
	}

	// 绑定JSON请求体到结构体，并验证参数
	if err := c.ShouldBindJSON(&request); err != nil {
		// 如果参数无效，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request parameters", "code": 400})
		return
	}

	// 调用服务层的LoginUser函数进行登录，并获取用户信息
	user, err := services.LoginUser(request.UserNameOrEmail, request.Password)
	if err != nil {
		// 如果登录失败，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	// 将userId存储在cookie中
	c.SetCookie("token", user.UserID, 3600, "/", "localhost", false, true)

	// 登录成功，返回200状态码
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "code": 200})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 从Cookie中获取token
	token, err := c.Cookie("token")
	if err != nil {
		// 如果没有token，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is required", "code": 400})
		return
	}

	// 调用服务层的GetUserInfoByToken函数获取用户信息
	user, err := services.GetUserInfoByToken(token)
	if err != nil {
		// 如果获取失败，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	// 获取成功，返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"userId":   user.UserID,
		"userName": user.UserName,
		"email":    user.Email,
	})
}

// LogoutUser 用户退出
func LogoutUser(c *gin.Context) {
	// 获取token参数
	_, err := c.Cookie("token")
	if err != nil {
		// 如果没有token，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is required", "code": 400})
		return
	}

	// 清除cookie
	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	// 返回登出成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "用户已登出"})
}

// ResetPassword 处理用户重置密码请求
func ResetPassword(c *gin.Context) {
	// 定义请求结构体，包含密码、邮箱和邮箱验证码
	var request struct {
		Password  string `json:"password" binding:"required,min=8"`                    // 密码，必填，最小长度8
		Email     string `json:"email" binding:"required,email"`                       // 邮箱，必填，必须是有效邮箱格式
		EmailCode string `json:"emailCode" binding:"required,len=6,regexp=^[0-9]{6}$"` // 邮箱验证码，必填，长度6，必须是数字
	}

	// 绑定JSON请求体到结构体，并验证参数
	if err := c.ShouldBindJSON(&request); err != nil {
		// 如果参数无效，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request parameters", "code": 400})
		return
	}

	// 调用服务层的ResetPassword函数进行密码重置
	if err := services.ResetPassword(request.Email, request.Password, request.EmailCode); err != nil {
		// 如果重置失败，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	// 重置成功，返回200状态码
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful", "code": 200})
}
