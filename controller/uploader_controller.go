package controller

import (
	"chat-gpt-service/db"
	"chat-gpt-service/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUploadedFileHandler(c *gin.Context) {
	// Assume you have an ID parameter in the URL.
	fileID := c.Param("id")

	fmt.Println("File Id", fileID)

	var file model.UploadedFile
	query := db.DB.Where("id = ?", fileID)
	if err := query.First(&file).Error; err == nil {
		// Data found in the database
		c.Data(http.StatusOK, file.ContentType, file.ContentData)
		return
	}

	c.Header("Cache-Control", "max-age="+strconv.Itoa(24*60*60))
	c.JSON(http.StatusNotFound, "file not found")
	// Serve the image data as a response.
}
