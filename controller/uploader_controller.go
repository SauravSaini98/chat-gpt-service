package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestUploadedFileHandler(c *gin.Context) {
	// Assume you have an ID parameter in the URL.
	var requestParams struct {
		ImageUrl string `json:"image_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// url, err := uploader.FileUploader(requestParams.ImageUrl)

	// if err != nil {
	// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// return
	// }

	c.JSON(http.StatusOK, gin.H{"data": "saved"})
	// Serve the image data as a response.
}
