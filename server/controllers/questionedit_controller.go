package controllers

import (
	"net/http"
	"server/services"

	"github.com/gin-gonic/gin"
)

// GetSurveyMetaController 获取问卷元数据
func GetSurveyMetaController(c *gin.Context) {
	// 从路径参数中获取 surveyId
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "surveyId is required",
			"code":    400,
		})
		return
	}

	// 调用服务层获取元数据
	meta, err := services.GetSurveyMetaService(surveyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// 返回元数据
	c.JSON(http.StatusOK, meta)
}

// GetSurveyQuestionsController 获取问卷题目信息
func GetSurveyQuestionsController(c *gin.Context) {
	// 从路径参数中获取 surveyId
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "surveyId is required",
			"code":    400,
		})
		return
	}

	// 调用服务层获取数据
	survey, err := services.GetSurveyQuestionsService(surveyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, survey)
}

// SaveSurveyEditController 处理问卷编辑保存的控制器
func SaveSurveyEditController(c *gin.Context) {
	// 获取路径参数中的 surveyId
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "surveyId is required",
			"code":    400,
		})
		return
	}

	// 解析请求体
	var surveyData services.SurveyModel
	if err := c.ShouldBindJSON(&surveyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    400,
		})
		return
	}

	// 调用服务层逻辑保存问卷
	err := services.SaveSurveyEditService(surveyId, &surveyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Survey updated successfully",
		"code":    200,
	})
}

// DeleteSurveyController 问卷删除控制器
func DeleteSurveyController(c *gin.Context) {
	// 获取路径参数中的 surveyId
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "surveyId is required",
			"code":    400,
		})
		return
	}

	// 调用服务层逻辑删除问卷
	err := services.DeleteSurveyService(surveyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Survey deleted successfully",
		"code":    200,
	})
}
