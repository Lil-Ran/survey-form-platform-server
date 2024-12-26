package controllers

import (
	"encoding/json"
	"net/http"
	"server/common"
	"server/services"

	"github.com/gorilla/mux"
)

// GetSurveyMetaHandler 获取问卷元数据
func GetSurveyMetaHandler(w http.ResponseWriter, r *http.Request) {
	surveyID := mux.Vars(r)["surveyId"]
	survey, err := services.GetSurveyMeta(surveyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey)
}

// GetSurveyQuestionsHandler 获取问卷的所有题目
func GetSurveyQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	surveyID := mux.Vars(r)["surveyId"]
	questions, err := services.GetSurveyQuestions(surveyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// AddSurveyQuestionHandler 新增问卷题目
func AddSurveyQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var question common.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	surveyID := mux.Vars(r)["surveyId"]
	question.SurveyID = surveyID
	if err := services.AddSurveyQuestion(question); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// EditSurveyQuestionHandler 修改问卷题目
func EditSurveyQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var question common.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := services.EditSurveyQuestion(question); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteSurveyQuestionHandler 删除问卷题目
func DeleteSurveyQuestionHandler(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["questionId"]
	surveyID := mux.Vars(r)["surveyId"]
	if err := services.DeleteSurveyQuestion(questionID, surveyID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// MoveSurveyQuestionHandler 移动问卷题目顺序
func MoveSurveyQuestionHandler(w http.ResponseWriter, r *http.Request) {
	surveyID := mux.Vars(r)["surveyId"]
	var requestBody struct {
		QuestionID string `json:"questionId"`
		NewIndex   int    `json:"newIndex"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := services.MoveSurveyQuestion(requestBody.QuestionID, surveyID, requestBody.NewIndex); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
