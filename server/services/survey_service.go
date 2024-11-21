package services

import (
	"errors"
	"server/common"

	"gorm.io/gorm"
)

func CreateSurvey(survey *common.Survey) error {
	var user common.User
	result := common.DB.Where("UserID = ?", survey.UserID).First(&user)
	if result.Error != nil {
		return errors.New("User does not exist")
	}

	survey.UserID = user.UserID

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

// ListSurveys 获取问卷列表
func ListSurveys(count, skip int) ([]common.Survey, int64, error) {
	var surveys []common.Survey
	var total int64

	// 统计总数
	if err := common.DB.Model(&common.Survey{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := common.DB.Offset(skip).Limit(count).Find(&surveys).Error; err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
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
	return common.DB.Model(&common.Survey{}).Where("survey_id = ?", surveyID).Update("status", status).Error
}

// CopySurvey 复制问卷
func CopySurvey(surveyID string) (*common.Survey, error) {
	var survey common.Survey
	if err := common.DB.Where("survey_id = ?", surveyID).First(&survey).Error; err != nil {
		return nil, errors.New("survey not found")
	}

	// 创建新的问卷
	newSurvey := survey
	newSurvey.SurveyID = "" // 让 GORM 自动生成新 ID
	if err := common.DB.Create(&newSurvey).Error; err != nil {
		return nil, err
	}
	return &newSurvey, nil
}
