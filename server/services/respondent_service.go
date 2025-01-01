package services

import (
	"errors"
	"server/common"
)

type QuestionModel struct {
	Type        string                      `json:"QuestionType"`
	Label       string                      `json:"QuestionLabel"`
	QuestionID  string                      `json:"QuestionID"`
	Title       string                      `json:"title"`
	Description string                      `json:"Explanation"`
	LeastChoice int                         `json:"leastChoice"`
	MaxChoice   int                         `json:"maxChoice"`
	SurveyID    string                      `json:"surveyId"`
	Options     []common.QuestionOption     `json:"Options"`
	NumFillIns  []common.QuestionNumFillIn  `json:"NumFillIns"`
	TextFillIns []common.QuestionTextFillIn `json:"TextFillIns"`
}

type SurveyModel struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	IsOpening bool            `json:"isopening"`
	Questions []QuestionModel `json:"questions"`
}

type ResponseModel struct {
	ResponseID        string                  `json:"responseid"`
	SurveyID          string                  `json:"surveyid"`
	QuestionsResponse []QuestionResponseModel `json:"questionsResponse"`
}

type QuestionResponseModel struct {
	ResponseID string `json:"responseid"`
	QID        string `json:"qid"`
	Type       string `json:"type"`
}

// GetRespondentQuestionsController 获取问卷及问题
func GetRespondentQuestionsController(surveyId string) (*SurveyModel, error) {
	var survey common.Survey

	// 查询 Survey
	err := common.DB.Where("SurveyID = ?", surveyId).First(&survey).Error
	if err != nil {
		return nil, errors.New("survey not found")
	}

	// 从 QuestionIDs 中逐个查询 Question 表
	questions := make([]QuestionModel, 0)
	for _, questionID := range survey.QuestionIDs {
		var question common.Question
		err := common.DB.Where("QuestionID = ?", questionID).First(&question).Error
		if err != nil {
			return nil, errors.New("Failed to find question: " + questionID)
		}

		// 将 Question 转换为 QuestionModel
		questions = append(questions, QuestionModel{
			Type:        question.QuestionType,
			Label:       question.QuestionLabel,
			QuestionID:  question.QuestionID,
			Title:       question.Title,
			Description: question.Description,
			LeastChoice: question.LeastChoice,
			MaxChoice:   question.MaxChoice,
			SurveyID:    question.SurveyID,
		})
	}

	// 构造响应数据
	return &SurveyModel{
		ID:        survey.SurveyID,
		Title:     survey.Title,
		IsOpening: survey.Status == "open",
		Questions: questions,
	}, nil
}

// SubmitSurveyResponseService 保存答卷
func SubmitSurveyResponseService(response ResponseModel) error {
	// 检查问卷是否存在
	var survey common.Survey
	err := common.DB.Where("SurveyID = ?", response.SurveyID).First(&survey).Error
	if err != nil {
		return errors.New("survey not found")
	}

	// 保存答卷
	var surveyResponse common.SurveyResponse
	surveyResponse.ResponseID = response.ResponseID
	surveyResponse.SurveyID = response.SurveyID
	err = common.DB.Create(&surveyResponse).Error
	if err != nil {
		return errors.New("failed to save response")
	}

	// 保存问题答卷
	for _, question := range response.QuestionsResponse {
		var questionResponse common.QuestionResponse
		questionResponse.ResponseID = response.ResponseID
		questionResponse.QuestionID = question.QID
		questionResponse.SurveyID = response.SurveyID
		err := common.DB.Create(&questionResponse).Error
		if err != nil {
			return errors.New("failed to save question response")
		}
	}

	return nil
}
