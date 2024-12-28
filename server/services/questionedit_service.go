package services

import (
	"errors"
	"server/common"

	"gorm.io/gorm"
)

// GetSurveyMeta 获取问卷元数据
func GetSurveyMeta(surveyID string) (common.Survey, error) {
	var survey common.Survey
	if err := common.DB.Where("survey_id = ?", surveyID).First(&survey).Error; err != nil {
		return survey, errors.New("survey not found")
	}
	return survey, nil
}

// GetSurveyQuestions 获取问卷的所有题目
func GetSurveyQuestions(surveyID string) ([]common.Question, error) {
	var questions []common.Question
	if err := common.DB.Where("survey_id = ?", surveyID).Order("question_index ASC").Find(&questions).Error; err != nil {
		return nil, errors.New("failed to fetch questions")
	}
	return questions, nil
}

// AddSurveyQuestion 新增问卷题目
func AddSurveyQuestion(question common.Question) error {
	var maxIndex int
	common.DB.Model(&common.Question{}).Where("survey_id = ?", question.SurveyID).Select("MAX(question_index)").Scan(&maxIndex)
	question.QuestionIndex = maxIndex + 1 // 设置新题目的索引
	return common.DB.Create(&question).Error
}

// EditSurveyQuestion 修改问卷题目
func EditSurveyQuestion(question common.Question) error {
	result := common.DB.Model(&common.Question{}).Where("question_id = ?", question.QuestionID).Updates(map[string]interface{}{
		"title":             question.Title,
		"description":       question.Description,
		"is_required":       question.IsRequired,
		"question_type":     question.QuestionType,
		"display_condition": question.DisplayCondition,
	})
	if result.RowsAffected == 0 {
		return errors.New("question not found")
	}
	return result.Error
}

// DeleteSurveyQuestion 删除问卷题目
func DeleteSurveyQuestion(questionID string, surveyID string) error {
	return common.DB.Transaction(func(tx *gorm.DB) error {
		var question common.Question
		if err := tx.Where("question_id = ? AND survey_id = ?", questionID, surveyID).First(&question).Error; err != nil {
			return errors.New("question not found")
		}

		// 删除题目
		if err := tx.Delete(&question).Error; err != nil {
			return err
		}

		// 更新剩余题目的索引
		if err := tx.Exec(`
			UPDATE questions 
			SET question_index = question_index - 1 
			WHERE survey_id = ? AND question_index > ?
		`, surveyID, question.QuestionIndex).Error; err != nil {
			return err
		}
		return nil
	})
}

// MoveSurveyQuestion 移动问卷题目顺序
func MoveSurveyQuestion(questionID string, surveyID string, newIndex int) error {
	return common.DB.Transaction(func(tx *gorm.DB) error {
		var question common.Question
		if err := tx.Where("question_id = ? AND survey_id = ?", questionID, surveyID).First(&question).Error; err != nil {
			return errors.New("question not found")
		}

		oldIndex := question.QuestionIndex
		if oldIndex == newIndex {
			return nil // 无需移动
		}

		// 更新中间题目索引
		if oldIndex < newIndex {
			// 向下移动
			if err := tx.Exec(`
				UPDATE questions
				SET question_index = question_index - 1
				WHERE survey_id = ? AND question_index > ? AND question_index <= ?
			`, surveyID, oldIndex, newIndex).Error; err != nil {
				return err
			}
		} else {
			// 向上移动
			if err := tx.Exec(`
				UPDATE questions
				SET question_index = question_index + 1
				WHERE survey_id = ? AND question_index >= ? AND question_index < ?
			`, surveyID, newIndex, oldIndex).Error; err != nil {
				return err
			}
		}

		// 更新当前题目的索引
		return tx.Model(&common.Question{}).Where("question_id = ?", questionID).Update("question_index", newIndex).Error
	})
}
