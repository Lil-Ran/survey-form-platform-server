package services

import (
	"errors"
	"server/common"
	"strings"
)

// SurveyMetaModel 定义符合 API 文档的响应结构
type SurveyMetaModel struct {
	SurveyID       string `json:"surveyId"`
	AccessID       string `json:"accessId"`
	UserID         string `json:"userId"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CreateTime     string `json:"createTime"`
	LastUpdateTime string `json:"lastUpdateTime"`
	LastUpdateUser string `json:"lastUpdateUser"`
	Status         string `json:"status"`
}

// GetSurveyMetaService 获取问卷元数据
func GetSurveyMetaService(surveyId string) (*SurveyMetaModel, error) {
	var survey common.Survey

	// 查询数据库获取问卷
	err := common.DB.Where("SurveyID = ?", surveyId).First(&survey).Error
	if err != nil {
		return nil, errors.New("survey not found")
	}

	// 构造元数据响应
	meta := &SurveyMetaModel{
		SurveyID:       survey.SurveyID,
		AccessID:       survey.AccessID,
		UserID:         survey.UserID,
		Title:          survey.Title,
		Description:    survey.Description,
		CreateTime:     survey.CreateTime.String(),
		LastUpdateTime: survey.LastUpdateTime.String(),
		LastUpdateUser: "admin", // 假设使用固定管理员作为示例
		Status:         survey.Status,
	}

	return meta, nil
}

//====================================================================================================

// Option 定义符合 API 的选项结构
type Option struct {
	OptionID      string `json:"OptionID"`
	OptionContent string `json:"OptionContent"`
	QuestionID    string `json:"QuestionID"`
	SurveyID      string `json:"SurveyID"`
}

// NumFillIn 定义符合 API 的数字填空结构
type NumFillIn struct {
	NumFillInID string `json:"NumFillInID"`
	QuestionID  string `json:"QuestionID"`
	SurveyID    string `json:"SurveyID"`
}

// TextFillIn 定义符合 API 的文本填空结构
type TextFillIn struct {
	TextFillInID string `json:"TextFillInID"`
	QuestionID   string `json:"QuestionID"`
	SurveyID     string `json:"SurveyID"`
}

// // GetSurveyQuestionsService 获取问卷题目信息
// func GetSurveyQuestionsService(surveyId string) (*SurveyModel, error) {
// 	var survey common.Survey

// 	// 查询问卷信息
// 	err := common.DB.Where("SurveyID = ?", surveyId).First(&survey).Error
// 	if err != nil {
// 		return nil, errors.New("survey not found")
// 	}

// 	// 查询题目信息
// 	var questions []common.Question
// 	err = common.DB.Where("SurveyID = ?", surveyId).Find(&questions).Error
// 	if err != nil {
// 		return nil, errors.New("questions not found")
// 	}

// 	questionModels := make([]QuestionModel, 0)
// 	for _, question := range questions {
// 		// 查询选项
// 		var options []common.QuestionOption
// 		err := common.DB.Where("QuestionID = ? AND SurveyID = ?", question.QuestionID, surveyId).Find(&options).Error
// 		if err != nil {
// 			return nil, errors.New("Options not found for question " + question.QuestionID)
// 		}

// 		// 查询数字填空
// 		var numFillIns []common.QuestionNumFillIn
// 		err = common.DB.Where("QuestionID = ? AND SurveyID = ?", question.QuestionID, surveyId).Find(&numFillIns).Error
// 		if err != nil {
// 			return nil, errors.New("NumFillIns not found for question " + question.QuestionID)
// 		}

// 		// 查询文本填空
// 		var textFillIns []common.QuestionTextFillIn
// 		err = common.DB.Where("QuestionID = ? AND SurveyID = ?", question.QuestionID, surveyId).Find(&textFillIns).Error
// 		if err != nil {
// 			return nil, errors.New("TextFillIns not found for question " + question.QuestionID)
// 		}

// 		// 构造题目信息
// 		questionModels = append(questionModels, QuestionModel{
// 			Type:        question.QuestionType,
// 			Label:       question.QuestionLabel,
// 			QuestionID:  question.QuestionID,
// 			Title:       question.Title,
// 			Description: question.Description,
// 			LeastChoice: question.LeastChoice,
// 			MaxChoice:   question.MaxChoice,
// 			SurveyID:    question.SurveyID,
// 			Options:     options,     // 直接使用查询结果，无需再构建
// 			NumFillIns:  numFillIns,  // 直接使用查询结果，无需再构建
// 			TextFillIns: textFillIns, // 直接使用查询结果，无需再构建
// 		})
// 	}

// 	// 构造问卷信息
// 	return &SurveyModel{
// 		ID:        survey.SurveyID,
// 		Title:     survey.Title,
// 		IsOpening: survey.Status == "Ongoing",
// 		Questions: questionModels,
// 	}, nil
// }

// GetRespondentQuestionsController 获取问卷及问题
func GetSurveyQuestionsService(surveyId string) (*SurveyModel, error) {
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
		if question.OptionIDs != "" {
			optionIDs := strings.Split(question.OptionIDs, ",")
			for _, optionID := range optionIDs {
				var option common.QuestionOption
				err := common.DB.Where("OptionID = ? AND SurveyID = ?", optionID, surveyId).First(&option).Error
				if err != nil {
					return nil, errors.New("Option not found for optionID: " + optionID)
				}
				options = append(options, option)
			}
		}

		// 查询数字填空
		var numFillIns []common.QuestionNumFillIn
		if question.NumFillInIDs != "" {
			numFillInIDs := strings.Split(question.NumFillInIDs, ",")
			for _, numFillInID := range numFillInIDs {
				var numFillIn common.QuestionNumFillIn
				err := common.DB.Where("NumFillInID = ? AND SurveyID = ?", numFillInID, surveyId).First(&numFillIn).Error
				if err != nil {
					return nil, errors.New("NumFillIn not found for NumFillInID: " + numFillInID)
				}
				numFillIns = append(numFillIns, numFillIn)
			}
		}

		// 查询文本填空
		var textFillIns []common.QuestionTextFillIn
		if question.TextFillInIDs != "" {
			textFillInIDs := strings.Split(question.TextFillInIDs, ",")
			for _, textFillInID := range textFillInIDs {
				var textFillIn common.QuestionTextFillIn
				err := common.DB.Where("TextFillInID = ? AND SurveyID = ?", textFillInID, surveyId).First(&textFillIn).Error
				if err != nil {
					return nil, errors.New("TextFillIn not found for TextFillInID: " + textFillInID)
				}
				textFillIns = append(textFillIns, textFillIn)
			}
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

// SaveSurveyEditService 处理问卷编辑保存的服务层逻辑
func SaveSurveyEditService(surveyId string, surveyData *SurveyModel) error {
	// 检查问卷是否存在
	var survey common.Survey
	err := common.DB.Where("SurveyID = ?", surveyId).First(&survey).Error
	if err != nil {
		return errors.New("survey not found")
	}

	// 更新问卷信息
	survey.Title = surveyData.Title
	survey.Status = "Ongoing"

	err = common.DB.Save(&survey).Error
	if err != nil {
		return errors.New("failed to update survey")
	}

	// 删除旧问题及其相关数据
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.Question{}).Error
	if err != nil {
		return errors.New("failed to delete old questions")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionOption{}).Error
	if err != nil {
		return errors.New("failed to delete old options")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionTextFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete old text fill-ins")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionNumFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete old num fill-ins")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionResponse{}).Error
	if err != nil {
		return errors.New("failed to delete old responsequestion fill-ins")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.ResponseOption{}).Error
	if err != nil {
		return errors.New("failed to delete old responseoption fill-ins")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.ResponseTextFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete old responsetext fill-ins")
	}
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.ResponseNumFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete old responsenumfill fill-ins")
	}
	err = common.DB.Debug().Where("SurveyID = ?", surveyId).Delete(&common.SurveyResponse{}).Error
	if err != nil {
		return errors.New("failed to delete old response fill-ins")
	}
	//清空 Survey 中的 ResponseCount 字段
	err = common.DB.Debug().Model(&survey).Where("SurveyID = ?", surveyId).Update("ResponseCount", 0).Error
	if err != nil {
		return errors.New("failed to reset response count")
	}

	// if err := common.DB.Model(&survey).Where("SurveyID = ?", surveyId).Update("ResponseCount", survey.ResponseCount+1).Error; err != nil {
	// 	return errors.New("failed to update survey response count: " + err.Error())
	// }

	// 存储问题 ID 的列表
	questionIDs := []string{}

	// 保存新问题及其相关数据
	for _, question := range surveyData.Questions {
		// 收集选项 IDs
		optionIDs := []string{}
		for _, option := range question.Options {
			optionIDs = append(optionIDs, option.OptionID)
		}

		// 收集文本填空 IDs
		textFillInIDs := []string{}
		for _, textFillIn := range question.TextFillIns {
			textFillInIDs = append(textFillInIDs, textFillIn.TextFillInID)
		}

		// 收集数字填空 IDs
		numFillInIDs := []string{}
		for _, numFillIn := range question.NumFillIns {
			numFillInIDs = append(numFillInIDs, numFillIn.NumFillInID)
		}

		// 收集问题 ID
		questionIDs = append(questionIDs, question.QuestionID)

		// 构造问题数据
		newQuestion := common.Question{
			QuestionID:    question.QuestionID,
			SurveyID:      surveyId,
			Title:         question.Title,
			Description:   question.Description,
			LeastChoice:   int(question.LeastChoice),
			MaxChoice:     int(question.MaxChoice),
			QuestionType:  question.Type,
			QuestionLabel: question.Label,
			OptionIDs:     strings.Join(optionIDs, ","),
			TextFillInIDs: strings.Join(textFillInIDs, ","),
			NumFillInIDs:  strings.Join(numFillInIDs, ","),
		}

		// 插入新问题
		err = common.DB.Create(&newQuestion).Error
		if err != nil {
			return errors.New("Failed to save question: " + question.QuestionID)
		}

		// 保存问题选项
		for _, option := range question.Options {
			newOption := common.QuestionOption{
				OptionID:      option.OptionID,
				QuestionID:    question.QuestionID,
				SurveyID:      surveyId,
				OptionContent: option.OptionContent,
			}
			err = common.DB.Create(&newOption).Error
			if err != nil {
				return errors.New("Failed to save option: " + option.OptionID)
			}
		}

		// 保存文本填空
		for _, textFillIn := range question.TextFillIns {
			newTextFillIn := common.QuestionTextFillIn{
				TextFillInID: textFillIn.TextFillInID,
				QuestionID:   question.QuestionID,
				SurveyID:     surveyId,
			}
			err = common.DB.Create(&newTextFillIn).Error
			if err != nil {
				return errors.New("Failed to save text fill-in: " + textFillIn.TextFillInID)
			}
		}

		// 保存数字填空
		for _, numFillIn := range question.NumFillIns {
			newNumFillIn := common.QuestionNumFillIn{
				NumFillInID: numFillIn.NumFillInID,
				QuestionID:  question.QuestionID,
				SurveyID:    surveyId,
			}
			err = common.DB.Create(&newNumFillIn).Error
			if err != nil {
				return errors.New("Failed to save num fill-in: " + numFillIn.NumFillInID)
			}
		}
	}

	// 将问题 ID 列表保存为以逗号分隔的字符串
	survey.QuestionIDs = strings.Join(questionIDs, ",")
	err = common.DB.Save(&survey).Error
	if err != nil {
		return errors.New("failed to save survey question IDs")
	}

	return nil
}

// DeleteSurveyService 处理问卷删除逻辑
func DeleteSurveyService(surveyId string) error {
	// 检查问卷是否存在
	var survey common.Survey
	err := common.DB.Where("SurveyID = ?", surveyId).First(&survey).Error
	if err != nil {
		return errors.New("survey not found")
	}

	// 删除问卷相关联的问题
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.Question{}).Error
	if err != nil {
		return errors.New("failed to delete questions related to the survey")
	}

	// 删除问卷相关联的问题选项
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionOption{}).Error
	if err != nil {
		return errors.New("failed to delete question options related to the survey")
	}

	// 删除问卷相关联的文本填空
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionTextFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete text fill-ins related to the survey")
	}

	// 删除问卷相关联的数字填空
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.QuestionNumFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete num fill-ins related to the survey")
	}

	// 删除问卷本身
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.Survey{}).Error
	if err != nil {
		return errors.New("failed to delete survey")
	}

	// 删除问卷相关联的答卷选项
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.ResponseOption{}).Error
	if err != nil {
		return errors.New("failed to delete response options related to the survey")
	}

	// 删除问卷相关联的文本填空答卷
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.ResponseTextFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete response text fill-ins related to the survey")
	}

	// 删除问卷相关联的答卷
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.SurveyResponse{}).Error
	if err != nil {
		return errors.New("failed to delete survey responses related to the survey")
	}

	// 删除问卷相关联的数字填空答卷
	err = common.DB.Where("SurveyID = ?", surveyId).Delete(&common.ResponseNumFillIn{}).Error
	if err != nil {
		return errors.New("failed to delete response num fill-ins related to the survey")
	}

	return nil
}
