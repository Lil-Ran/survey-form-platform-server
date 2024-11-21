package controllers

import (
	"net/http"
	"server/common"
	"server/services"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 问卷创建
func CreateSurvey(c *gin.Context) {
	var survey common.Survey
	if err := c.ShouldBindJSON(&survey); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	if err := services.CreateSurvey(&survey); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create survey")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Survey created successfully", survey)
}

// // 通过surveyId获取问卷
// func GetSurvey(c *gin.Context) {
// 	surveyId := c.Param("surveyId")
// 	survey, err := services.GetSurveyByID(surveyId)
// 	if err != nil {
// 		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get survey")
// 		return
// 	}

// 	utils.SuccessResponse(c, http.StatusOK, "Survey retrieved successfully", survey)
// }

// ListSurveys 获取问卷列表
func ListSurveys(c *gin.Context) {
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))

	surveys, total, err := services.ListSurveys(count, skip)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve surveys")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Surveys retrieved successfully", gin.H{
		"data":  surveys,
		"total": total,
	})
}

// UpdateSurveyStatus 更新问卷状态（开始/暂停/删除）
func UpdateSurveyStatus(c *gin.Context) {
	var statusUpdate struct {
		SurveyID string `json:"surveyId"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	if err := services.UpdateSurveyStatus(statusUpdate.SurveyID, statusUpdate.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Survey status updated successfully", nil)
}

// CopySurvey 复制问卷
func CopySurvey(c *gin.Context) {
	var copyRequest struct {
		SurveyID string `json:"surveyId"`
	}

	if err := c.ShouldBindJSON(&copyRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	newSurvey, err := services.CopySurvey(copyRequest.SurveyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Survey copied successfully", newSurvey)
}
