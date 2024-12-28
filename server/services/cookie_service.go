package services

import (
	"fmt"
	"server/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte
var tokenExpiry time.Duration

// 初始化服务配置
func InitAuthConfig() {
	jwtSecret = []byte(config.Config.Auth.JWTSecret)
	// 解析 token_expiry 配置
	expiry, err := time.ParseDuration(config.Config.Auth.TokenExpiry)
	if err != nil {
		fmt.Println("Error:", err)
		panic("Invalid token_expiry format in configuration")
	}
	tokenExpiry = expiry
}

// 生成 JWT
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(tokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 验证 JWT
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// 设置 Cookie
func SetCookie(c *gin.Context, userID string) error {
	token, err := GenerateJWT(userID)
	if err != nil {
		return err
	}
	c.SetCookie("token", token, int(tokenExpiry.Seconds()), "/", "localhost", false, true)
	return nil
}

// 获取 Cookie
func GetCookie(c *gin.Context) (jwt.MapClaims, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	return ValidateJWT(token)
}

// 删除 Cookie
func DeleteCookie(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", true, true)
	// c.JSON(http.StatusOK, gin.H{"message": "Cookie has been deleted"})
}
