package controllers

import (
	"net/http"
	"server/common"
	"server/services"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// 问卷创建
func GreateSurvey(c *gin.Context) {
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

// 通过surveyId获取问卷
func GetSurvey(c *gin.Context) {
	surveyId := c.Param("surveyId")
	survey, err := services.GetSurveyByID(surveyId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get survey")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Survey retrieved successfully", survey)
}
