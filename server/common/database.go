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
	gorm.Model
	UserID       string    `gorm:"column:UserID;primaryKey"` // 用户ID
	UserName     string    `gorm:"column:UserName"`          // 用户名
	Email        string    `gorm:"column:Email"`             // 邮箱
	Password     string    `gorm:"column:Password"`          // 密码
	RegisterDate time.Time `gorm:"column:RegisterDate"`      // 注册日期
	Surveys      []Survey  `gorm:"foreignKey:UserID"`        // 用户的问卷列表
}

// Survey 问卷结构体
type Survey struct {
	gorm.Model
	SurveyID          string     `gorm:"column:SurveyID;primaryKey"` // 问卷ID
	AccessID          string     `gorm:"column:AccessID"`            // 访问ID
	UserID            string     `gorm:"column:UserID;index"`        // 用户ID
	Title             string     `gorm:"column:Title"`               // 问卷标题
	Description       string     `gorm:"column:Description"`         // 问卷描述
	CreateTime        time.Time  `gorm:"column:CreateTime"`          // 创建时间
	ExpireTime        time.Time  `gorm:"column:ExpireTime"`          // 过期时间
	LastUpdateTime    time.Time  `gorm:"column:LastUpdateTime"`      // 最后更新时间
	LastUpdateUser    string     `gorm:"column:LastUpdateUser"`      // 最后更新用户
	Status            int        `gorm:"column:Status"`              // 问卷状态
	ResponseCount     int        `gorm:"column:ResponseCount"`       // 响应数量
	ThemeColor        int        `gorm:"column:ThemeColor"`          // 主题颜色
	TextColor         int        `gorm:"column:TextColor"`           // 文字颜色
	PCBackgroundImage string     `gorm:"column:PCBackgroundImage"`   // PC背景图片
	PCBannerImage     string     `gorm:"column:PCBannerImage"`       // PC横幅图片
	Footer            string     `gorm:"column:Footer"`              // 页脚
	DisplayStyle      int        `gorm:"column:DisplayStyle"`        // 显示样式
	ButtonText        string     `gorm:"type:json"`                  // JSON 存储
	StartTime         time.Time  `gorm:"column:StartTime"`           // 开始时间
	EndTime           time.Time  `gorm:"column:EndTime"`             // 结束时间
	DayStartTime      time.Time  `gorm:"column:DayStartTime"`        // 每日开始时间
	DayEndTime        time.Time  `gorm:"column:DayEndTime"`          // 每日结束时间
	PasswordStrategy  int        `gorm:"column:PasswordStrategy"`    // 密码策略
	Password          string     `gorm:"type:json"`                  // JSON 存储
	MaxResponseCount  int        `gorm:"column:MaxResponseCount"`    // 最大响应数量
	BrowserLimit      bool       `gorm:"column:BrowserLimit"`        // 浏览器限制
	IPLimit           bool       `gorm:"column:IPLimit"`             // IP限制
	KeepContent       bool       `gorm:"column:KeepContent"`         // 保留内容
	FailMessage       string     `gorm:"column:FailMessage"`         // 失败消息
	ShowAfterSubmit   int        `gorm:"column:ShowAfterSubmit"`     // 提交后显示
	ShowContent       string     `gorm:"column:ShowContent"`         // 显示内容
	Questions         []Question `gorm:"foreignKey:SurveyID"`        // 问卷中的问题列表
	Responses         []Response `gorm:"foreignKey:SurveyID"`        // 问卷的响应列表
}

// Question 问题结构体
type Question struct {
	gorm.Model
	QuestionID       string           `gorm:"column:QuestionID;primaryKey"` // 问题ID
	SurveyID         string           `gorm:"column:SurveyID;index"`        // 问卷ID
	QuestionIndex    int              `gorm:"column:QuestionIndex"`         // 问题索引
	Title            string           `gorm:"column:Title"`                 // 问题标题
	Description      string           `gorm:"column:Description"`           // 问题描述
	IsRequired       bool             `gorm:"column:IsRequired"`            // 是否必填
	QuestionType     int              `gorm:"column:QuestionType"`          // 问题类型
	DisplayCondition string           `gorm:"column:DisplayCondition"`      // 显示条件
	Options          []QuestionOption `gorm:"foreignKey:QuestionID"`        // 问题选项列表
}

// QuestionOption 问题选项结构体
type QuestionOption struct {
	gorm.Model
	OptionID         string `gorm:"column:OptionID;primaryKey"` // 选项ID
	QuestionID       string `gorm:"column:QuestionID;index"`    // 问题ID
	SurveyID         string `gorm:"column:SurveyID;index"`      // 问卷ID
	OptionIndex      int    `gorm:"column:OptionIndex"`         // 选项索引
	Title            string `gorm:"column:Title"`               // 选项标题
	DisplayCondition string `gorm:"column:DisplayCondition"`    // 显示条件
	Attribute        int    `gorm:"column:Attribute"`           // 选项属性
}

// TextFillIn 文本填空结构体
type TextFillIn struct {
	gorm.Model
	TextFillInID string `gorm:"column:TextFillInID;primaryKey"` // 文本填空ID
	QuestionID   string `gorm:"column:QuestionID;index"`        // 问题ID
	TextContent  string `gorm:"column:TextContent"`             // 文本内容
}

// NumFillIn 数字填空结构体
type NumFillIn struct {
	gorm.Model
	NumFillInID string `gorm:"column:NumFillInID;primaryKey"` // 数字填空ID
	QuestionID  string `gorm:"column:QuestionID;index"`       // 问题ID
	NumContent  int    `gorm:"column:NumContent"`             // 数字内容
	MaxNum      int    `gorm:"column:MaxNum"`                 // 最大值
	LeastNum    int    `gorm:"column:LeastNum"`               // 最小值
	Precision   int    `gorm:"column:Precision"`              // 精度
}

// Response 答卷结构体
type Response struct {
	gorm.Model
	ResponseID    string         `gorm:"column:ResponseID;primaryKey"` // 答卷ID
	SurveyID      string         `gorm:"column:SurveyID;index"`        // 问卷ID
	ResponseIndex int            `gorm:"column:ResponseIndex"`         // 答卷索引
	StartTime     time.Time      `gorm:"column:StartTime"`             // 开始时间
	SubmitTime    time.Time      `gorm:"column:SubmitTime"`            // 提交时间
	Duration      int            `gorm:"column:Duration"`              // 持续时间
	Source        string         `gorm:"column:Source"`                // 来源
	IP            string         `gorm:"column:IP"`                    // IP地址
	IsStar        bool           `gorm:"column:IsStar"`                // 是否加星
	IsInvalid     bool           `gorm:"column:IsInvalid"`             // 是否无效
	Data          []ResponseData `gorm:"foreignKey:ResponseID"`        // 答卷数据列表
}

// ResponseData 全部答卷内容结构体
type ResponseData struct {
	gorm.Model
	ResponseID string `gorm:"column:ResponseID;index"` // 答卷ID
	SurveyID   string `gorm:"column:SurveyID;index"`   // 问卷ID
	QuestionID string `gorm:"column:QuestionID;index"` // 问题ID
	OptionID   string `gorm:"column:OptionID;index"`   // 选项ID
	Content    string `gorm:"column:Content"`          // 内容
}

// EmailVerification 邮箱验证码结构体
type EmailVerification struct {
	gorm.Model
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
	err = DB.AutoMigrate(&User{}, &Survey{}, &Question{}, &QuestionOption{}, &TextFillIn{}, &NumFillIn{}, &Response{}, &ResponseData{}, &EmailVerification{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	fmt.Println("Database migration completed!")
}
