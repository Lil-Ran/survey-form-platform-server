package common

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户结构体
type User struct {
	gorm.Model
	UserID       string    `gorm:"column:UserID;primaryKey"`
	UserName     string    `gorm:"column:UserName"`
	Email        string    `gorm:"column:Email"`
	Password     string    `gorm:"column:Password"`
	RegisterDate time.Time `gorm:"column:RegisterDate"`
	Surveys      []Survey  `gorm:"foreignKey:UserID"`
}

// Survey 问卷结构体
type Survey struct {
	gorm.Model
	SurveyID          string     `gorm:"column:SurveyID;primaryKey"`
	AccessID          string     `gorm:"column:AccessID"`
	UserID            string     `gorm:"column:UserID;index"`
	Title             string     `gorm:"column:Title"`
	Description       string     `gorm:"column:Description"`
	CreateTime        time.Time  `gorm:"column:CreateTime"`
	ExpireTime        time.Time  `gorm:"column:ExpireTime"`
	LastUpdateTime    time.Time  `gorm:"column:LastUpdateTime"`
	LastUpdateUser    string     `gorm:"column:LastUpdateUser"`
	Status            int        `gorm:"column:Status"`
	ResponseCount     int        `gorm:"column:ResponseCount"`
	ThemeColor        int        `gorm:"column:ThemeColor"`
	TextColor         int        `gorm:"column:TextColor"`
	PCBackgroundImage string     `gorm:"column:PCBackgroundImage"`
	PCBannerImage     string     `gorm:"column:PCBannerImage"`
	Footer            string     `gorm:"column:Footer"`
	DisplayStyle      int        `gorm:"column:DisplayStyle"`
	ButtonText        []string   `gorm:"column:ButtonText"`
	StartTime         time.Time  `gorm:"column:StartTime"`
	EndTime           time.Time  `gorm:"column:EndTime"`
	DayStartTime      time.Time  `gorm:"column:DayStartTime"`
	DayEndTime        time.Time  `gorm:"column:DayEndTime"`
	PasswordStrategy  int        `gorm:"column:PasswordStrategy"`
	Password          []string   `gorm:"column:Password"`
	MaxResponseCount  int        `gorm:"column:MaxResponseCount"`
	BrowserLimit      bool       `gorm:"column:BrowserLimit"`
	IPLimit           bool       `gorm:"column:IPLimit"`
	KeepContent       bool       `gorm:"column:KeepContent"`
	FailMessage       string     `gorm:"column:FailMessage"`
	ShowAfterSubmit   int        `gorm:"column:ShowAfterSubmit"`
	ShowContent       string     `gorm:"column:ShowContent"`
	Questions         []Question `gorm:"foreignKey:SurveyID"`
	Responses         []Response `gorm:"foreignKey:SurveyID"`
}

// Question 问题结构体
type Question struct {
	gorm.Model
	QuestionID       string           `gorm:"column:QuestionID;primaryKey"`
	SurveyID         string           `gorm:"column:SurveyID;index"`
	QuestionIndex    int              `gorm:"column:QuestionIndex"`
	Title            string           `gorm:"column:Title"`
	Description      string           `gorm:"column:Description"`
	IsRequired       bool             `gorm:"column:IsRequired"`
	QuestionType     int              `gorm:"column:QuestionType"`
	DisplayCondition string           `gorm:"column:DisplayCondition"`
	Options          []QuestionOption `gorm:"foreignKey:QuestionID"`
}

// QuestionOption 问题选项结构体
type QuestionOption struct {
	gorm.Model
	OptionID         string `gorm:"column:OptionID;primaryKey"`
	QuestionID       string `gorm:"column:QuestionID;index"`
	SurveyID         string `gorm:"column:SurveyID;index"`
	OptionIndex      int    `gorm:"column:OptionIndex"`
	Title            string `gorm:"column:Title"`
	DisplayCondition string `gorm:"column:DisplayCondition"`
	Attribute        int    `gorm:"column:Attribute"`
}

// TextFillIn 文本填空结构体
type TextFillIn struct {
	gorm.Model
	TextFillInID string `gorm:"column:TextFillInID;primaryKey"`
	QuestionID   string `gorm:"column:QuestionID;index"`
	TextContent  string `gorm:"column:TextContent"`
}

// NumFillIn 数字填空结构体
type NumFillIn struct {
	gorm.Model
	NumFillInID string `gorm:"column:NumFillInID;primaryKey"`
	QuestionID  string `gorm:"column:QuestionID;index"`
	NumContent  int    `gorm:"column:NumContent"`
	MaxNum      int    `gorm:"column:MaxNum"`
	LeastNum    int    `gorm:"column:LeastNum"`
	Precision   int    `gorm:"column:Precision"`
}

// Response 答卷结构体
type Response struct {
	gorm.Model
	ResponseID    string         `gorm:"column:ResponseID;primaryKey"`
	SurveyID      string         `gorm:"column:SurveyID;index"`
	ResponseIndex int            `gorm:"column:ResponseIndex"`
	StartTime     time.Time      `gorm:"column:StartTime"`
	SubmitTime    time.Time      `gorm:"column:SubmitTime"`
	Duration      int            `gorm:"column:Duration"`
	Source        string         `gorm:"column:Source"`
	IP            string         `gorm:"column:IP"`
	IsStar        bool           `gorm:"column:IsStar"`
	IsInvalid     bool           `gorm:"column:IsInvalid"`
	Data          []ResponseData `gorm:"foreignKey:ResponseID"`
}

// ResponseData 全部答卷内容结构体
type ResponseData struct {
	gorm.Model
	ResponseID string `gorm:"column:ResponseID;index"`
	SurveyID   string `gorm:"column:SurveyID;index"`
	QuestionID string `gorm:"column:QuestionID;index"`
	OptionID   string `gorm:"column:OptionID;index"`
	Content    string `gorm:"column:Content"`
}

func InitDb() {
	// 数据库连接信息
	dsn := "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移模式
	db.AutoMigrate(&User{}, &Survey{}, &Question{}, &QuestionOption{}, &TextFillIn{}, &NumFillIn{}, &Response{}, &ResponseData{})

}
