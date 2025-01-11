package services

import (
	"errors"
	"log"
	"server/common"
	"strings"
)

type QuestionModel struct {
	Type        string                      `json:"QuestionType"`
	Label       string                      `json:"QuestionLabel"`
	QuestionID  string                      `json:"QuestionID"`
	Title       string                      `json:"Title"`
	Description string                      `json:"Description"`
	LeastChoice int                         `json:"LeastChoice"`
	MaxChoice   int                         `json:"MaxChoice"`
	SurveyID    string                      `json:"SurveyID"`
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
	ResponseID        string                  `json:"ResponseID"`
	SurveyID          string                  `json:"SurveyID"`
	QuestionsResponse []QuestionResponseModel `json:"QuestionResponse"`
}

type QuestionResponseModel struct {
	ResponseID  string                      `json:"ResponseID"`
	QID         string                      `json:"QuestionID"`
	Type        string                      `json:"QuestionType"`
	Options     []common.ResponseOption     `json:"Options"`
	TextFillIns []common.ResponseTextFillIn `json:"TextFillIns"`
	NumFillIns  []common.ResponseNumFillIn  `json:"NumFillIns"`
}

// GetRespondentQuestionsController 获取问卷及问题
func GetRespondentQuestionsController(surveyId string) (*SurveyModel, error) {
	var survey common.Survey

	// 查询 Survey
	err := common.DB.Where("SurveyID = ?", surveyId).First(&survey).Error
	if err != nil {
		return nil, errors.New("survey not found")
	}

	// 将 QuestionIDs 转换为问题 ID 的数组
	questionIDArray := strings.Split(survey.QuestionIDs, ",")

	// 从 QuestionIDs 中逐个查询 Question 表
	questions := make([]QuestionModel, 0)
	for _, questionID := range questionIDArray {
		var question common.Question
		err := common.DB.Where("QuestionID = ?", questionID).First(&question).Error
		if err != nil {
			return nil, errors.New("Failed to find question: " + questionID)
		}

		// 查询选项
		var options []common.QuestionOption
		err = common.DB.Where("QuestionID = ? AND SurveyID = ?", question.QuestionID, surveyId).Find(&options).Error
		if err != nil {
			return nil, errors.New("Options not found for question " + question.QuestionID)
		}

		// 查询数字填空
		var numFillIns []common.QuestionNumFillIn
		err = common.DB.Where("QuestionID = ? AND SurveyID = ?", question.QuestionID, surveyId).Find(&numFillIns).Error
		if err != nil {
			return nil, errors.New("NumFillIns not found for question " + question.QuestionID)
		}

		// 查询文本填空
		var textFillIns []common.QuestionTextFillIn
		err = common.DB.Where("QuestionID = ? AND SurveyID = ?", question.QuestionID, surveyId).Find(&textFillIns).Error
		if err != nil {
			return nil, errors.New("TextFillIns not found for question " + question.QuestionID)
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
			Options:     options,     // 直接使用查询结果，无需再构建
			NumFillIns:  numFillIns,  // 直接使用查询结果，无需再构建
			TextFillIns: textFillIns, // 直接使用查询结果，无需再构建
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

func SubmitSurveyResponseService(response ResponseModel) error {
	// 检查问卷是否存在
	var survey common.Survey
	if err := common.DB.Where("SurveyID = ?", response.SurveyID).First(&survey).Error; err != nil {
		return errors.New("survey not found")
	}

	// 检查是否已存在答卷
	var existingResponse common.SurveyResponse
	if err := common.DB.Where("ResponseID = ?", response.ResponseID).First(&existingResponse).Error; err == nil {
		return errors.New("response already exists")
	}

	// 保存答卷
	surveyResponse := common.SurveyResponse{
		ResponseID: response.ResponseID,
		SurveyID:   response.SurveyID,
	}
	if err := common.DB.Create(&surveyResponse).Error; err != nil {
		return errors.New("failed to save survey response: " + err.Error())
	}

	// 保存问题答卷
	for _, question := range response.QuestionsResponse {
		// 初始化字段，避免 nil 数据
		if question.Options == nil {
			question.Options = []common.ResponseOption{}
		}
		if question.TextFillIns == nil {
			question.TextFillIns = []common.ResponseTextFillIn{}
		}
		if question.NumFillIns == nil {
			question.NumFillIns = []common.ResponseNumFillIn{}
		}

		switch question.Type {
		case "SingleChoice", "MultiChoice": // 单选/多选题
			for _, option := range question.Options {

				responseOption := common.ResponseOption{
					ResponseID:    response.ResponseID,
					OptionID:      option.OptionID,
					QuestionID:    question.QID,
					SurveyID:      response.SurveyID,
					OptionContent: option.OptionContent,
					IsSelect:      option.IsSelect,
				}
				if err := common.DB.Create(&responseOption).Error; err != nil {
					return errors.New("failed to save response option: " + err.Error())
				}
			}
		case "SingleTextFillIn", "MultiTextFillIn": // 单文本/多文本填空题
			for _, textFillIn := range question.TextFillIns {
				responseTextFillIn := common.ResponseTextFillIn{
					ResponseID:   response.ResponseID,
					TextFillInID: textFillIn.TextFillInID,
					QuestionID:   question.QID,
					SurveyID:     response.SurveyID,
					TextContent:  textFillIn.TextContent,
				}
				if err := common.DB.Create(&responseTextFillIn).Error; err != nil {
					return errors.New("failed to save text fill-in response: " + err.Error())
				}
			}
		case "SingleNumFillIn", "MultiNumFillIn": // 单数字/多数字填空题
			for _, numFillIn := range question.NumFillIns {
				responseNumFillIn := common.ResponseNumFillIn{
					ResponseID:  response.ResponseID,
					NumFillInID: numFillIn.NumFillInID,
					QuestionID:  question.QID,
					SurveyID:    response.SurveyID,
					NumContent:  numFillIn.NumContent,
				}
				if err := common.DB.Create(&responseNumFillIn).Error; err != nil {
					return errors.New("failed to save number fill-in response: " + err.Error())
				}
			}
		default:
			log.Printf("Unsupported question type: %s", question.Type)
			continue // 跳过未知类型
		}
	}

	// 更新问卷的 ResponseCount
	if err := common.DB.Model(&survey).Update("ResponseCount", survey.ResponseCount+1).Error; err != nil {
		return errors.New("failed to update survey response count: " + err.Error())
	}

	return nil
}
