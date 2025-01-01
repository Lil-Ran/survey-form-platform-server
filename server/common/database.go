package common

import (
	"fmt"
	"server/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局 DB 变量，供其他包使用
var DB *gorm.DB

// User 用户结构体
type User struct {
	UserID       string    `gorm:"column:UserID;primaryKey;size:36;unique"` // UUID 格式
	UserName     string    `gorm:"column:UserName;size:100"`                // 限制长度
	Email        string    `gorm:"column:Email;size:150;unique"`            // 邮箱唯一
	Password     string    `gorm:"column:Password;size:255"`                // 密码长度
	RegisterDate time.Time `gorm:"column:RegisterDate"`                     // 可为空值
	Surveys      []Survey  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Survey 问卷结构体
type Survey struct {
	SurveyID          string    `gorm:"column:SurveyID;primaryKey;size:36"` // 问卷ID
	AccessID          string    `gorm:"column:AccessID"`                    // 访问ID
	UserID            string    `gorm:"column:UserID;size:36"`              // 用户ID
	Title             string    `gorm:"column:Title"`                       // 问卷标题
	Description       string    `gorm:"column:Description"`                 // 问卷描述
	CreateTime        time.Time `gorm:"column:CreateTime"`                  // 创建时间
	ExpireTime        time.Time `gorm:"column:ExpireTime"`                  // 过期时间
	LastUpdateTime    time.Time `gorm:"column:LastUpdateTime"`              // 最后更新时间
	Status            string    `gorm:"column:Status"`                      // 问卷状态
	ResponseCount     int       `gorm:"column:ResponseCount"`               // 响应数量
	ThemeColor        int       `gorm:"column:ThemeColor"`                  // 主题颜色
	TextColor         int       `gorm:"column:TextColor"`                   // 文字颜色
	PCBackgroundImage string    `gorm:"column:PCBackgroundImage"`           // PC背景图片
	PCBannerImage     string    `gorm:"column:PCBannerImage"`               // PC横幅图片
	Footer            *string   `gorm:"column:Footer"`                      // 页脚
	DisplayStyle      int       `gorm:"column:DisplayStyle"`                // 显示样式
	ButtonText        *string   `gorm:"type:json"`                          // JSON 存储
	StartTime         time.Time `gorm:"column:StartTime"`                   // 开始时间
	EndTime           time.Time `gorm:"column:EndTime"`                     // 结束时间
	DayStartTime      time.Time `gorm:"column:DayStartTime"`                // 每日开始时间
	DayEndTime        time.Time `gorm:"column:DayEndTime"`                  // 每日结束时间
	PasswordStrategy  int       `gorm:"column:PasswordStrategy"`            // 密码策略
	Password          string    `gorm:"type:json"`                          // JSON 存储
	MaxResponseCount  int       `gorm:"column:MaxResponseCount"`            // 最大响应数量
	BrowserLimit      bool      `gorm:"column:BrowserLimit"`                // 浏览器限制
	IPLimit           bool      `gorm:"column:IPLimit"`                     // IP限制
	KeepContent       bool      `gorm:"column:KeepContent"`                 // 保留内容
	FailMessage       string    `gorm:"column:FailMessage"`                 // 失败消息
	ShowAfterSubmit   int       `gorm:"column:ShowAfterSubmit"`             // 提交后显示
	ShowContent       string    `gorm:"column:ShowContent"`                 // 显示内容
	QuestionIDs       []string  `gorm:"type:json"`                          // 问卷中的问题列表
	ResponseIDs       []string  `gorm:"type:json"`                          // 问卷的响应列表
}

// Question 问题结构体
type Question struct {
	QuestionID    string `gorm:"column:QuestionID;primaryKey"` // 问题ID
	SurveyID      string `gorm:"column:SurveyID;index"`        // 问卷ID
	Title         string `gorm:"column:Title"`                 // 问题标题
	Description   string `gorm:"column:Description"`           // 问题描述
	LeastChoice   int    `gorm:"column:LeastChoice"`           //最少选择数
	MaxChoice     int    `gorm:"column:MaxChoice"`             //最多选择数
	QuestionType  string `gorm:"column:QuestionType"`          // 问题类型
	QuestionLabel string `gorm:"column:QuestionLabel"`         // 问题类型中文
	OptionIDs     string `gorm:"column:OptionIDs"`             // 问题选项列表
	TextFillInIDs string `gorm:"column:TextFillInIDs"`         // 文本填空框
	NumFillInIDs  string `gorm:"column:NumFillInIDs"`          // 数字填空类型
}

// QuestionOption 问题选项结构体
type QuestionOption struct {
	OptionID      string `gorm:"column:OptionID;primaryKey"` // 选项ID
	QuestionID    string `gorm:"column:QuestionID;index"`    // 问题ID
	SurveyID      string `gorm:"column:SurveyID;index"`      // 问卷ID
	OptionContent string `gorm:"column:OptionContent"`       // 选项内容
}

type QuestionTextFillIn struct {
	TextFillInID string `gorm:"column:TextFillInID;primaryKey"` // 文本填空ID
	QuestionID   string `gorm:"column:QuestionID;index"`        // 问题ID
	SurveyID     string `gorm:"column:SurveyID;index"`          // 问卷ID
}

type QuestionNumFillIn struct {
	NumFillInID string `gorm:"column:NumFillInID;primaryKey"` // 数字填空ID
	QuestionID  string `gorm:"column:QuestionID;index"`       // 问题ID
	SurveyID    string `gorm:"column:SurveyID;index"`         // 问卷ID
}

///====================================================///=============================///===========================================================================

// ResponseOption 问题选项结构体
type ResponseOption struct {
	ResponseID    string `gorm:"column:ResponseID;primaryKey"` // 联合主键之一
	OptionID      string `gorm:"column:OptionID;primaryKey"`   // 联合主键之一
	QuestionID    string `gorm:"column:QuestionID;index"`      // 问题ID
	SurveyID      string `gorm:"column:SurveyID;index"`        // 问卷ID
	OptionContent string `gorm:"column:OptionContent"`         // 选项内容
	IsSelect      bool   `gorm:"column:IsSelect"`              // 选项内容
}

// TextFillIn 文本填空结构体
type ResponseTextFillIn struct {
	ResponseID   string `gorm:"column:ResponseID;primaryKey"`   // 联合主键之一
	TextFillInID string `gorm:"column:TextFillInID;primaryKey"` // 文本填空ID
	QuestionID   string `gorm:"column:QuestionID;index"`        // 问题ID
	SurveyID     string `gorm:"column:SurveyID;index"`          // 问卷ID
	TextContent  string `gorm:"column:TextContent"`             // 文本内容
}

// NumFillIn 数字填空结构体
type ResponseNumFillIn struct {
	ResponseID  string `gorm:"column:ResponseID;primaryKey"`  // 联合主键之一
	NumFillInID string `gorm:"column:NumFillInID;primaryKey"` // 数字填空ID
	QuestionID  string `gorm:"column:QuestionID;index"`       // 问题ID
	SurveyID    string `gorm:"column:SurveyID;index"`         // 问卷ID
	NumContent  int    `gorm:"column:NumContent"`             // 数字内容
}

// QuestionResponse 问题结构体
type QuestionResponse struct {
	ResponseID            string   `gorm:"column:ResponseID;primaryKey"` // 答卷ID
	QuestionID            string   `gorm:"column:QuestionID;primaryKey"` // 问题ID
	SurveyID              string   `gorm:"column:SurveyID;index"`        // 问卷ID
	ResponseOptionIDs     []string `gorm:"type:json"`                    // 问题选项列表
	ResponseTextFillInIDs []string `gorm:"type:json"`                    // 文本填空框
	ResponseNumFillInIDs  []string `gorm:"type:json"`                    // 数字填空类型
}

// SurveyResponse 答卷结构体
type SurveyResponse struct {
	ResponseID string `gorm:"column:ResponseID;primaryKey"` // 答卷ID
	SurveyID   string `gorm:"column:SurveyID;index"`        // 问卷ID
	Source     string `gorm:"column:Source"`                // 来源
	IP         string `gorm:"column:IP"`                    // IP地址
	IsStar     bool   `gorm:"column:IsStar"`                // 是否加星
	IsInvalid  bool   `gorm:"column:IsInvalid"`             // 是否无效
}

// EmailVerification 邮箱验证码结构体
type EmailVerification struct {
	Email  string    `gorm:"column:Email;index"` // 邮箱
	Code   string    `gorm:"column:Code"`        // 验证码
	Expiry time.Time `gorm:"column:Expiry"`      // 过期时间
}

// InitDb 初始化数据库连接
func InitDb() {
	ds := config.Config.Datasource

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		ds.Username, ds.Password, ds.Host, ds.Port, ds.DBName, ds.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 打印日志，确认自动迁移开始
	fmt.Println("Starting AutoMigrate...")
	// 自动迁移所有表结构
	err = DB.AutoMigrate(
		&User{},               // 用户表
		&Survey{},             // 问卷表
		&Question{},           // 问题表
		&QuestionOption{},     // 问题选项表
		&QuestionTextFillIn{}, // 文本填空表
		&QuestionNumFillIn{},  // 数字填空表
		&ResponseOption{},     // 答卷选项表
		&ResponseTextFillIn{}, // 文本填空答卷表
		&ResponseNumFillIn{},  // 数字填空答卷表
		&QuestionResponse{},   // 问题答卷表
		&SurveyResponse{},     // 问卷答卷表
		&EmailVerification{},  // 邮箱验证表
	)
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	fmt.Println("Database migration completed!")
}
