package controllers

import (
	"net/http"
	"server/services"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// GetResponseDetails 获取答卷详细信息
func GetOptionCount(c *gin.Context) {
	// 从路径中获取 SurveyID
	surveyID := c.Param("SurveyID")
	if surveyID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "SurveyID is required")
		return
	}

	// 从请求体中获取 OptionID
	var request struct {
		OptionID string `json:"OptionID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	// 调用服务层逻辑
	count, err := services.GetOptionCount(surveyID, request.OptionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回响应
	utils.SuccessResponse(c, http.StatusOK, "Option count retrieved successfully", gin.H{
		"num": count,
	})
}

// GetTextFillinData 获取指定 TextFillinID 的文本数据
func GetTextFillinData(c *gin.Context) {
	// 从路径参数中获取 SurveyID
	surveyID := c.Param("SurveyID")
	if surveyID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "SurveyID is required")
		return
	}

	// 从请求体中获取 TextFillinID
	var request struct {
		TextFillinID string `json:"TextFillinID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	// 调用服务层逻辑
	texts, err := services.GetTextFillinData(surveyID, request.TextFillinID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, http.StatusOK, "Text data retrieved successfully", gin.H{
		"stringarray": texts,
	})
}

// GetNumFillinData 获取指定数字填空题的数值回答
func GetNumFillinData(c *gin.Context) {
	// 从路径参数中获取 SurveyID
	surveyID := c.Param("SurveyID")
	if surveyID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "SurveyID is required")
		return
	}

	// 从请求体中获取 NumFillInID
	var request struct {
		NumFillInID string `json:"NumFillInID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	// 调用服务层逻辑
	numbers, err := services.GetNumFillinData(surveyID, request.NumFillInID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回响应
	utils.SuccessResponse(c, http.StatusOK, "Number data retrieved successfully", gin.H{
		"numarray": numbers,
	})
}

// GetSurveyResponsesHandler 获取指定问卷的所有答卷内容
func GetSurveyResponsesHandler(c *gin.Context) {
	// 获取 SurveyID 参数
	surveyID := c.Param("SurveyID")
	if surveyID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "SurveyID is required")
		return
	}

	// 调用服务层逻辑
	responses, err := services.GetSurveyResponses(surveyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, responses)
}
