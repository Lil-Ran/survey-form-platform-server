package controllers

import (
	"fmt"
	"net/http"
	"server/services"

	"github.com/gin-gonic/gin"
)

// GetRespondentQuestionsController 获取答卷问题信息
func GetRespondentQuestionsController(c *gin.Context) {
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "surveyId is required",
			"code":    400,
		})
		return
	}

	// 调用服务层获取问卷数据
	survey, err := services.GetRespondentQuestionsController(surveyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// 返回结果，严格符合API文档
	c.JSON(http.StatusOK, survey)
}

// SubmitSurveyResponseController 提交答卷
func SubmitSurveyResponseController(c *gin.Context) {
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "surveyId is required",
			"code":    400,
		})
		return
	}

	var responseModel services.ResponseModel
	if err := c.ShouldBindJSON(&responseModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    400,
		})
		return
	}

	fmt.Printf("ResponseModel: %+v\n", responseModel)

	// 校验 surveyId 是否一致
	if responseModel.SurveyID != surveyId {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Survey ID mismatch",
			"code":    400,
		})
		return
	}

	// 调用服务层保存答卷
	err := services.SubmitSurveyResponseService(responseModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    500,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Response submitted successfully",
	})
}
