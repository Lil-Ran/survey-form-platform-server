package services

import (
	"errors"
	"server/common"
)

// GetOptionCount 获取选项被选择的数量
func GetOptionCount(surveyID, optionID string) (int64, error) {
	// 验证问卷是否存在
	var survey common.Survey
	if err := common.DB.Where("SurveyID = ?", surveyID).First(&survey).Error; err != nil {
		return 0, errors.New("survey not found")
	}

	// 在 ResponseOption 表中统计指定 OptionID 的选择次数
	var count int64
	if err := common.DB.Model(&common.ResponseOption{}).
		Where("SurveyID = ? AND OptionID = ? AND IsSelect = true", surveyID, optionID).
		Count(&count).Error; err != nil {
		return 0, errors.New("failed to count option selections")
	}

	return count, nil
}

// GetTextFillinData 获取指定填空题的所有回答
func GetTextFillinData(surveyID, textFillinID string) ([]string, error) {
	// 验证问卷是否存在
	var survey common.Survey
	if err := common.DB.Where("SurveyID = ?", surveyID).First(&survey).Error; err != nil {
		return nil, errors.New("survey not found")
	}

	// 验证文本填空题是否属于该问卷
	var questionTextFillIn common.QuestionTextFillIn
	if err := common.DB.Where("TextFillInID = ? AND QuestionID IN (SELECT QuestionID FROM Questions WHERE SurveyID = ?)", textFillinID, surveyID).
		First(&questionTextFillIn).Error; err != nil {
		return nil, errors.New("text fill-in question not found for the specified survey")
	}

	// 查询 ResponseTextFillIn 表中的所有回答
	var responses []common.ResponseTextFillIn
	if err := common.DB.Where("SurveyID = ? AND TextFillInID = ?", surveyID, textFillinID).
		Find(&responses).Error; err != nil {
		return nil, errors.New("failed to retrieve responses for the text fill-in question")
	}

	// 提取回答
	var answers []string
	for _, response := range responses {
		answers = append(answers, response.TextContent)
	}

	return answers, nil
}

// GetNumFillinData 获取指定数字填空题的所有数值回答
func GetNumFillinData(surveyID, numFillInID string) ([]int, error) {
	// 验证问卷是否存在
	var survey common.Survey
	if err := common.DB.Where("SurveyID = ?", surveyID).First(&survey).Error; err != nil {
		return nil, errors.New("survey not found")
	}

	// 验证数字填空题是否属于该问卷
	var questionNumFillIn common.QuestionNumFillIn
	if err := common.DB.Where("NumFillInID = ? AND QuestionID IN (SELECT QuestionID FROM Questions WHERE SurveyID = ?)", numFillInID, surveyID).
		First(&questionNumFillIn).Error; err != nil {
		return nil, errors.New("number fill-in question not found for the specified survey")
	}

	// 查询 ResponseNumFillIn 表中的所有回答
	var responses []common.ResponseNumFillIn
	if err := common.DB.Where("SurveyID = ? AND NumFillInID = ?", surveyID, numFillInID).
		Find(&responses).Error; err != nil {
		return nil, errors.New("failed to retrieve responses for the number fill-in question")
	}

	// 提取数值回答
	var numbers []int
	for _, response := range responses {
		numbers = append(numbers, response.NumContent)
	}

	return numbers, nil
}

// ResponseDetailModel 答卷详情返回模型
type ResponseDetailModel struct {
	ResponseID string           `json:"ResponseID"`
	SurveyID   string           `json:"SurveyID"`
	Questions  []QuestionDetail `json:"QuestionResponse"`
}

type QuestionDetail struct {
	ResponseID   string                   `json:"ResponseID"`
	QuestionID   string                   `json:"QuestionID"`
	Title        string                   `json:"Title"`
	Description  string                   `json:"Description"`
	QuestionType string                   `json:"QuestionType"` // 表示题目类型：选择、文本填空或数字填空
	Options      []OptionDetail           `json:"Options"`
	TextFillIns  []ResponseTextFillInData `json:"TextFillIns"`
	NumFillIns   []ResponseNumFillInData  `json:"NumFillIns"`
}

type OptionDetail struct {
	ResponseID    string `json:"ResponseID"`
	OptionID      string `json:"OptionID"`
	OptionContent string `json:"OptionContent"`
	QuestionID    string `json:"QuestionID"`
	IsSelect      bool   `json:"IsSelect"`
}

type ResponseTextFillInData struct {
	ResponseID   string `json:"ResponseID"`
	TextFillInID string `json:"TextFillInID"`
	QuestionID   string `json:"QuestionID"`
	TextContent  string `json:"TextContent"`
}

type ResponseNumFillInData struct {
	ResponseID  string `json:"ResponseID"`
	NumFillInID string `json:"NumFillInID"`
	QuestionID  string `json:"QuestionID"`
	NumContent  int    `json:"NumContent"`
}

// GetSurveyResponses 获取指定问卷的所有答卷内容
func GetSurveyResponses(surveyID string) ([]ResponseDetailModel, error) {
	// 验证问卷是否存在
	var survey common.Survey
	if err := common.DB.Where("SurveyID = ?", surveyID).First(&survey).Error; err != nil {
		return nil, errors.New("survey not found")
	}

	// 获取该问卷下的所有答卷
	var responses []common.SurveyResponse
	if err := common.DB.Where("SurveyID = ?", surveyID).Find(&responses).Error; err != nil {
		return nil, errors.New("failed to retrieve responses")
	}

	// 构建返回结果
	var responseDetails []ResponseDetailModel
	for _, response := range responses {
		// 获取该问卷下的所有问题
		var questions []common.Question
		if err := common.DB.Where("SurveyID = ?", surveyID).Find(&questions).Error; err != nil {
			return nil, errors.New("failed to retrieve questions")
		}

		// 构建问题详情
		var questionDetails []QuestionDetail
		for _, question := range questions {
			questionDetail := QuestionDetail{
				ResponseID:   response.ResponseID,
				QuestionID:   question.QuestionID,
				Title:        question.Title,
				Description:  question.Description,
				QuestionType: question.QuestionType,
			}

			switch question.QuestionType {
			case "SingleChoice", "MultiChoice": // 单选/多选
				var options []common.ResponseOption
				if err := common.DB.Where("QuestionID = ? AND ResponseID = ?", question.QuestionID, response.ResponseID).Find(&options).Error; err == nil {
					for _, option := range options {
						questionDetail.Options = append(questionDetail.Options, OptionDetail{
							ResponseID:    option.ResponseID,
							OptionID:      option.OptionID,
							OptionContent: option.OptionContent,
							QuestionID:    option.QuestionID,
							IsSelect:      option.IsSelect,
						})
					}
				}
				questionDetail.TextFillIns = []ResponseTextFillInData{}
				questionDetail.NumFillIns = []ResponseNumFillInData{}
			case "SingleTextFillIn", "MultiTextFillIn": // 单文本填空/多文本填空
				var textFillIns []common.ResponseTextFillIn
				if err := common.DB.Where("QuestionID = ? AND ResponseID = ?", question.QuestionID, response.ResponseID).Find(&textFillIns).Error; err == nil {
					for _, text := range textFillIns {
						questionDetail.TextFillIns = append(questionDetail.TextFillIns, ResponseTextFillInData{
							ResponseID:   text.ResponseID,
							TextFillInID: text.TextFillInID,
							QuestionID:   text.QuestionID,
							TextContent:  text.TextContent,
						})
					}
				}
				questionDetail.Options = []OptionDetail{}
				questionDetail.NumFillIns = []ResponseNumFillInData{}
			case "SingleNumFillIn", "MultiNumFillIn": // 单数字填空/多数字填空
				var numFillIns []common.ResponseNumFillIn
				if err := common.DB.Where("QuestionID = ? AND ResponseID = ?", question.QuestionID, response.ResponseID).Find(&numFillIns).Error; err == nil {
					for _, num := range numFillIns {
						questionDetail.NumFillIns = append(questionDetail.NumFillIns, ResponseNumFillInData{
							ResponseID:  num.ResponseID,
							NumFillInID: num.NumFillInID,
							QuestionID:  num.QuestionID,
							NumContent:  num.NumContent,
						})
					}
				}
				questionDetail.Options = []OptionDetail{}
				questionDetail.TextFillIns = []ResponseTextFillInData{}
			default:
				questionDetail.Options = []OptionDetail{}
				questionDetail.TextFillIns = []ResponseTextFillInData{}
				questionDetail.NumFillIns = []ResponseNumFillInData{}
			}

			questionDetails = append(questionDetails, questionDetail)
		}

		// 构建每个答卷的模型
		responseDetails = append(responseDetails, ResponseDetailModel{
			ResponseID: response.ResponseID,
			SurveyID:   surveyID,
			Questions:  questionDetails,
		})
	}

	return responseDetails, nil
}
