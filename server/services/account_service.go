package services

import (
	"errors"
	"math/rand"
	"net/smtp"
	"server/common"
	"time"
	"github.com/google/uuid"
	"regexp"
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
	from := "your-email@example.com"
	password := "your-email-password"

	// 设置SMTP服务器信息
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	// 设置邮件内容
	subject := "Subject: Your Verification Code\n"
	body := "Your verification code is: " + code
	message := []byte(subject + "\n" + body)

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
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
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
func GetUserInfoByToken(token string) (*common.User, error) {
	var user common.User
	//根据token查找用户
	if err := common.DB.Where("UserID = ?", token).First(&user).Error; err != nil {
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

