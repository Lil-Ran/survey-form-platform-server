package controllers

import (
	"net/http"
	"server/common"
	"server/services"
	"server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 问卷创建
func CreateSurvey(c *gin.Context) {

	claims, err := services.GetCookie(c)
	if err != nil {
		return
	}

	// 从 claims 中提取 userID
	userID, ok := claims["userID"].(string)
	if !ok {
		return
	}
	var request struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}
	now := time.Now()
	// 调用服务层创建问卷
	survey := common.Survey{
		SurveyID: uuid.New().String(),

		Status:         "Ongoing",
		Title:          request.Title,
		CreateTime:     &now,
		LastUpdateTime: &now,
		LastUpdateUser: userID,
	}
	if err := services.CreateSurvey(userID, &survey); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create survey")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Survey created successfully", nil)
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

// ListSurveys 获取用户问卷列表
func ListSurveys(c *gin.Context) {
	// 从 Cookie 中提取用户 ID
	claims, err := services.GetCookie(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized: invalid or missing token")
		return
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid token claims")
		return
	}

	// 获取分页参数
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))

	// 调用服务层获取问卷列表
	surveys, total, length, err := services.ListSurveys(userID, count, skip)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve surveys")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Surveys retrieved successfully", gin.H{
		"data":   surveys, // surveys 直接作为数组返回
		"total":  total,   // 总数
		"length": length,  // 当前分页的数量
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

	err := services.CopySurvey(copyRequest.SurveyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Survey copied successfully", nil)
}
