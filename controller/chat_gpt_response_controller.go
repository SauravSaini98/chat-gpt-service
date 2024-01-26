// controller/chat_gpt_response_controller.go
package controller

import (
	"chat-gpt-service/db"
	"chat-gpt-service/helper"
	"chat-gpt-service/model"
	"chat-gpt-service/uploader"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetChatGPTResponseHandler handles the creation or retrieval of a ChatGPTResponse.
func GetChatGPTResponseHandler(c *gin.Context) {
	var requestParams struct {
		Engine       string `json:"engine" binding:"required"`
		Prompt       string `json:"prompt" binding:"required"`
		ResponseType string `json:"response_type"`
	}

	if err := c.ShouldBindJSON(&requestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors := make(map[string]string)
	allowedResponseTypes := []string{"json", "string"}

	engine := trimAndValidate(requestParams.Engine, "engine", &validationErrors)
	prompt := trimAndValidate(requestParams.Prompt, "prompt", &validationErrors)
	responseType := validateEnumVal(requestParams.ResponseType, "response_type", &validationErrors, allowedResponseTypes)

	// Check if there are validation errors
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Check the local database

	var chatGPTNullResponse model.ChatGPTResponse
	nullQuery := db.DB.Where("prompt = ? and engine = ? and success is null and created_at > (NOW() - INTERVAL '30 minutes')", prompt, engine)

	if err := nullQuery.Last(&chatGPTNullResponse).Error; err == nil {
		// Data found in the database
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Request already send please wait and resend it again"})
		return
	}

	var chatGPTResponse model.ChatGPTResponse
	query := db.DB.Where("prompt = ? and success is true and engine = ?", prompt, engine)

	if err := query.First(&chatGPTResponse).Error; err == nil {
		// Data found in the database
		jsonResponseData := helper.SetJSONResponse(chatGPTResponse.Answer, responseType)
		c.JSON(http.StatusOK, gin.H{"data": jsonResponseData})
		return
	}

	newChatGPTResponse := model.ChatGPTResponse{
		Engine: engine,
		Prompt: prompt,
	}

	if err := db.DB.Create(&newChatGPTResponse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If data not found, make a request to a third-party API (simulated here)
	response, err := helper.GetChatCompleteResponse(prompt, requestParams.Engine, 1000)

	if err != nil {
		newChatGPTResponse.Success = false
		if err := db.DB.Save(&newChatGPTResponse).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newChatGPTResponse.Answer = response
	newChatGPTResponse.Success = true

	if err := db.DB.Save(&newChatGPTResponse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonResponseData := helper.SetJSONResponse(response, responseType)

	// c.JSON(http.StatusOK, gin.H{"data": newChatGPTResponse.Answer})
	c.JSON(http.StatusOK, gin.H{"data": jsonResponseData})
}

// GetChatGPTResponseHandler handles the creation or retrieval of a ChatGPTResponse.
func GetChatGPTVisionResponseHandler(c *gin.Context) {
	engine := "gpt-4-vision-preview"
	var requestParams struct {
		ImageUrl     string `json:"image_url" binding:"required"`
		Prompt       string `json:"prompt" binding:"required"`
		ResponseType string `json:"response_type"`
	}

	if err := c.ShouldBindJSON(&requestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors := make(map[string]string)
	allowedResponseTypes := []string{"json", "string"}

	imageUrl := trimAndValidate(requestParams.ImageUrl, "image_url", &validationErrors)
	prompt := trimAndValidate(requestParams.Prompt, "prompt", &validationErrors)
	responseType := validateEnumVal(requestParams.ResponseType, "response_type", &validationErrors, allowedResponseTypes)

	// Check if there are validation errors
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	uploadedFile, err := uploader.CheckAndSaveImage(imageUrl)

	fmt.Println("UPloaded File OBJECT", uploadedFile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var chatGPTNullResponse model.ChatGPTResponse
	nullQuery := db.DB.Where("prompt = ? and engine = ? and uploaded_file_id = ? and success IS null AND created_at > (NOW() - INTERVAL '5 minutes')", prompt, engine, uploadedFile.ID)
	if err := nullQuery.Last(&chatGPTNullResponse).Error; err == nil {
		// Data found in the database
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Request already send please wait and resend it again"})
		return
	}

	// Check the local database
	var chatGPTResponse model.ChatGPTResponse
	query := db.DB.Where("prompt = ? and success is true and engine = ? and uploaded_file_id = ?", prompt, engine, uploadedFile.ID)

	if err := query.Last(&chatGPTResponse).Error; err == nil {
		// Data found in the database
		jsonResponseData := helper.SetJSONResponse(chatGPTResponse.Answer, responseType)
		c.JSON(http.StatusOK, gin.H{"data": jsonResponseData})
		return
	}

	// imageFileUrl := uploadedFile.FileURL
	newChatGPTResponse := model.ChatGPTResponse{
		Engine:         "gpt-4-vision-preview",
		Prompt:         prompt,
		UploadedFileID: uploadedFile.ID,
	}

	if err := db.DB.Create(&newChatGPTResponse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uploadedFileUrl := uploadedFile.BuildFileUrl()

	fmt.Println("UPloaded File URL", uploadedFileUrl)

	// If data not found, make a request to a third-party API (simulated here)
	response, err := helper.GetChatGptVisionResponse(prompt, uploadedFileUrl, 1000)

	if err != nil {
		newChatGPTResponse.Success = false
		if err := db.DB.Save(&newChatGPTResponse).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newChatGPTResponse.Answer = response
	newChatGPTResponse.Success = true

	if err := db.DB.Save(&newChatGPTResponse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonResponseData := helper.SetJSONResponse(response, responseType)

	// c.JSON(http.StatusOK, gin.H{"data": newChatGPTResponse.Answer})
	c.JSON(http.StatusOK, gin.H{"data": jsonResponseData})
}

// helper functions
func trimAndValidate(param string, paramName string, errors *map[string]string) string {
	trimmedParam := strings.TrimSpace(param)
	if trimmedParam == "" {
		errorMessage := paramName + " cannot be empty"
		(*errors)[paramName] = errorMessage
		return ""
	}

	return trimmedParam
}

// helper functions
func validateEnumVal(param string, paramName string, errors *map[string]string, allowedValues []string) string {
	trimmedParam := trimAndValidate(param, paramName, errors)
	if trimmedParam == "" {
		return ""
	}

	// Check if trimmedParam is in the allowed values
	validValue := false
	for _, allowed := range allowedValues {
		if trimmedParam == allowed {
			validValue = true
			break
		}
	}

	if !validValue {
		errorMessage := paramName + " should be one of " + strings.Join(allowedValues, ", ")
		(*errors)[paramName] = errorMessage
		return ""
	}

	return trimmedParam
}
