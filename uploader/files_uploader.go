package uploader

import (
	"chat-gpt-service/helper"
	"errors"
	"fmt"
	"strings"
)

func FileUploader(fileData []byte, contentType string) (string, error) {
	// Specify the environment (e.g., "production" or "local")
	appEnvironment := helper.GetEnvWithDefault("APP_ENV", "development")
	// Specify the image URL

	extension, err := helper.GetExtensionFromContentType(contentType)
	if err != nil {
		fmt.Println("Invalid content type:", contentType)
		return "", err
	}

	fileKey := fmt.Sprintf("chat-gpt/%s.%s", helper.GenerateUniqueFileName(), extension)

	// Determine storage destination based on the environment
	switch strings.ToLower(appEnvironment) {
	case "production":
		// Call the function to save to AWS S3
		url, err := helper.SaveFileToS3(fileData, fileKey)
		if err != nil {
			fmt.Println("Error saving to AWS S3:", err)
			return "", err
		}
		fmt.Println("File saved to AWS S3!")
		return url, nil
	case "development":
		// Call the function to save locally
		filePath, err := helper.SaveFileToLocal(fileData, fileKey)
		if err != nil {
			fmt.Println("Error saving locally:", err)
			return "", err
		}
		fmt.Println("File saved locally. Path:", filePath)
		return filePath, nil
	default:
		fmt.Println("Invalid environment specified.")
	}

	return "", errors.New("Invalid environment specified")
}
