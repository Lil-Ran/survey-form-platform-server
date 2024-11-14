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
