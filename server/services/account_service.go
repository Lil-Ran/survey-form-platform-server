package services

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
	"server/common"
	"server/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterUser 注册用户
func RegisterUser(userName, password, email, emailCode string) error {
	// 验证用户名格式
	if len(userName) < 4 || len(userName) > 64 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !isValidUserName(userName) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	// 验证密码格式
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !isValidPassword(password) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	// 验证邮箱验证码
	var emailVerification struct {
		Email  string
		Code   string
		Expiry time.Time
	}
	if err := common.DB.Raw("SELECT email, code, expiry FROM email_verifications WHERE email = ? AND code = ?", email, emailCode).Scan(&emailVerification).Error; err != nil {
		return errors.New("invalid email code")
	}

	// 检查验证码是否过期
	if emailVerification.Expiry.Before(time.Now()) {
		return errors.New("email code expired")
	}

	// 创建新用户到数据库
	userID := uuid.New().String()
	sql := "INSERT INTO users (UserID, UserName, Email, Password, RegisterDate) VALUES (?, ?, ?, ?, ?)"
	if err := common.DB.Exec(sql, userID, userName, email, password, time.Now()).Error; err != nil {
		return errors.New("failed to register user")
	}

	return nil
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(email string) error {
	// 生成随机验证码
	code := generateRandomCode(6)

	// 保存验证码到数据库
	sql := "INSERT INTO email_verifications (Email, Code, Expiry) VALUES (?, ?, ?)"
	if err := common.DB.Exec(sql, email, code, time.Now().Add(10*time.Minute)).Error; err != nil {
		return errors.New("failed to save email code")
	}

	// 发送邮件
	if err := sendEmail(email, code); err != nil {
		if strings.Contains(err.Error(), "short response") {
			return nil
		}
		// fmt.Printf("SendMail returned error: %v\n", err)
		return errors.New("failed to send email")
	}

	return nil
}

// generateRandomCode 生成指定长度的随机数字验证码
func generateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < length; i++ {
		code += string(rune('0' + r.Intn(10)))
	}
	return code
}

// sendEmail 发送邮件
func sendEmail(to, code string) error {
	from := config.Config.SMTP.From
	password := config.Config.SMTP.Password
	smtpHost := config.Config.SMTP.Host
	smtpPort := config.Config.SMTP.Port

	// 设置邮件内容
	subject := "Subject: Your Verification Code\n\n"
	fromHeader := fmt.Sprintf("From: SURVEY—FORM <%s>\n", from)
	toHeader := fmt.Sprintf("To: %s\n", to)
	body := fmt.Sprintf("Your verification code is: %s\n", code)

	// 构建完整的邮件内容
	message := []byte(fromHeader + toHeader + subject + "\n" + body)

	// 认证信息
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 发送邮件
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)

	if err != nil {
		return err
	}

	return nil
}

// isValidUserName 验证用户名格式
func isValidUserName(userName string) bool {
	// 使用正则表达式验证用户名格式
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(userName)
}

// isValidPassword 验证密码格式
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	// hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	return hasLower && hasUpper && hasDigit
	// return true
}

// LoginUser 用户登录
func LoginUser(userNameOrEmail, password string) (*common.User, error) {
	var user common.User
	// 查询用户信息，用户名或邮箱匹配
	if err := common.DB.Where("UserName = ? OR Email = ?", userNameOrEmail, userNameOrEmail).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 验证密码
	if user.Password != password {
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}

// GetUserInfoByToken 根据token获取用户信息
func GetUserInfoByToken(c *gin.Context) (*common.User, error) {
	// 通过 GetCookie 提取并验证 Cookie 中的 JWT
	claims, err := GetCookie(c)
	if err != nil {
		return nil, errors.New("unauthorized: invalid or missing token")
	}

	// 从 claims 中提取 userID
	userID, ok := claims["userID"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// 根据 userID 查询用户信息
	var user common.User
	if err := common.DB.Where("UserID = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// ResetPassword 重置用户密码
func ResetPassword(email, newPassword, emailCode string) error {
	// 验证邮箱验证码
	var emailVerification common.EmailVerification
	if err := common.DB.Where("Email = ? AND Code = ?", email, emailCode).First(&emailVerification).Error; err != nil {
		return errors.New("invalid email code")
	}

	// 检查验证码是否过期
	if emailVerification.Expiry.Before(time.Now()) {
		return errors.New("email code expired")
	}

	// 更新用户密码
	if err := common.DB.Model(&common.User{}).Where("Email = ?", email).Update("Password", newPassword).Error; err != nil {
		return errors.New("failed to reset password")
	}

	return nil
}
