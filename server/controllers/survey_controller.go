package controllers

import (
	"math/rand"
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	accessID := rng.Intn(900000) + 100000 // 生成六位数字（100000到999999）
	// 调用服务层创建问卷
	// 构造问卷对象
	survey := common.Survey{
		SurveyID:          uuid.New().String(),
		AccessID:          strconv.Itoa(accessID),
		Title:             request.Title,
		Description:       "", // 默认描述为空字符串，可以根据需要调整
		CreateTime:        now,
		ExpireTime:        now.AddDate(0, 1, 0), // 默认过期时间为 1 个月后
		LastUpdateTime:    now,
		Status:            "Ongoing",
		ResponseCount:     0,                                                                            // 初始响应数量为 0
		ThemeColor:        0,                                                                            // 默认主题颜色
		TextColor:         0,                                                                            // 默认文字颜色
		PCBackgroundImage: "",                                                                           // 默认背景图片为空
		PCBannerImage:     "",                                                                           // 默认横幅图片为空
		Footer:            nil,                                                                          // 页脚默认值为空
		DisplayStyle:      0,                                                                            // 默认显示样式
		ButtonText:        nil,                                                                          // 按钮文字默认值为空
		StartTime:         now,                                                                          // 默认开始时间为当前时间
		EndTime:           now.AddDate(0, 1, 0),                                                         // 默认结束时间为 1 个月后
		DayStartTime:      time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),    // 默认每日开始时间为 00:00
		DayEndTime:        time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()), // 默认每日结束时间为 23:59
		PasswordStrategy:  0,                                                                            // 默认密码策略
		Password:          "{}",                                                                         // 默认无密码
		MaxResponseCount:  0,                                                                            // 默认无限制
		BrowserLimit:      false,                                                                        // 默认不限制浏览器
		IPLimit:           false,                                                                        // 默认不限制 IP
		KeepContent:       false,                                                                        // 默认不保留内容
		FailMessage:       "",                                                                           // 默认失败消息为空
		ShowAfterSubmit:   0,                                                                            // 默认提交后不显示内容
		ShowContent:       "",                                                                           // 默认显示内容为空
		QuestionIDs:       "{}",                                                                         // 默认无问题列表
		ResponseIDs:       nil,                                                                          // 默认无响应列表
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
	// 检查 surveys 是否为空
	if surveys == nil {
		surveys = []services.SurveyResponse{} // 将 nil 替换为空数组
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
