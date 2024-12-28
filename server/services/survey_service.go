package services

import (
	"errors"
	// "fmt"
	"server/common"
	"time"

	"github.com/google/uuid"

	// "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateSurvey(userID string, survey *common.Survey) error {

	var user common.User
	result := common.DB.Where("UserID = ?", userID).First(&user)
	if result.Error != nil {
		return result.Error
	}

	survey.UserID = userID

	result = common.DB.Create(survey)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetSurveyByID(surveyID string) (*common.Survey, error) {
	var survey common.Survey
	if err := common.DB.Preload("Questions.Options").Where("survey_id = ?", surveyID).First(&survey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey not found")
		}
		return nil, err
	}
	return &survey, nil
}

// ListSurveys 获取用户的问卷列表
type SurveyResponse struct {
	SurveyID           string     `json:"surveyId"`
	AccessID           *string    `json:"accessId"`
	Title              string     `json:"title"`
	Status             string     `json:"status"`
	ResponseCount      int        `json:"responseCount"`
	OwnerID            string     `json:"ownerId"`
	OwnerName          string     `json:"ownerName"`
	CreateTime         *time.Time `json:"createTime"`
	LastUpdateTime     *time.Time `json:"lastUpdateTime"`
	LastUpdateUserID   string     `json:"lastUpdataUserID"`
	LastUpdateUserName string     `json:"lastUpdateUserName"`
}

func ListSurveys(userID string, count, skip int) ([]SurveyResponse, int64, int64, error) {
	var surveys []common.Survey
	var total int64

	// 统计用户的问卷总数
	if err := common.DB.Model(&common.Survey{}).Where("UserID = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 分页查询用户的问卷
	if err := common.DB.Model(&common.Survey{}).Where("UserID = ?", userID).Offset(skip).Limit(count).Find(&surveys).Error; err != nil {
		return nil, 0, 0, err
	}
	// 调试日志
	// fmt.Printf("Surveys retrieved: %+v\n", surveys)
	// 构建响应结构
	var responses []SurveyResponse
	for _, survey := range surveys {
		var owner common.User

		// 获取问卷拥有者信息
		if err := common.DB.Where("UserID = ?", survey.UserID).First(&owner).Error; err != nil {
			return nil, 0, 0, err
		}

		// 获取最后更新用户信息
		var lastUpdateUserName string
		var lastUpdateUser common.User
		if survey.LastUpdateUser != "" {
			if err := common.DB.Where("UserID = ?", survey.LastUpdateUser).First(&lastUpdateUser).Error; err == nil {
				lastUpdateUserName = lastUpdateUser.UserName
			}
		}

		// 构建单个问卷响应
		responses = append(responses, SurveyResponse{
			SurveyID:           survey.SurveyID,
			AccessID:           survey.AccessID,
			Title:              survey.Title,
			Status:             survey.Status,
			ResponseCount:      survey.ResponseCount,
			OwnerID:            survey.UserID,
			OwnerName:          owner.UserName,
			CreateTime:         survey.CreateTime,
			LastUpdateTime:     survey.LastUpdateTime,
			LastUpdateUserID:   survey.LastUpdateUser,
			LastUpdateUserName: lastUpdateUserName,
		})
	}

	return responses, total, int64(len(responses)), nil
}

// UpdateSurveyStatus 更新问卷状态
func UpdateSurveyStatus(surveyID, status string) error {
	// 确保状态合法性
	validStatuses := map[string]bool{
		"Ongoing": true, "Suspended": true, "Deleted": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// 更新状态
	return common.DB.Model(&common.Survey{}).Where("SurveyID = ?", surveyID).Update("status", status).Error
}

// CopySurvey 复制问卷
func CopySurvey(surveyID string) error {
	var survey common.Survey
	if err := common.DB.Where("SurveyID = ?", surveyID).First(&survey).Error; err != nil {
		return errors.New("survey not found")
	}

	// 创建新的问卷
	newSurvey := survey
	newSurvey.SurveyID = uuid.New().String() // 让 GORM 自动生成新 ID
	if err := common.DB.Create(&newSurvey).Error; err != nil {
		return err
	}
	return nil
}
