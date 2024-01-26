package helper

import (
	"chat-gpt-service/db"
	"chat-gpt-service/model"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"strings"
)

func CheckAndSaveImage(url string) (string, error) {
	uploadedFileData, contentType, err := downloadImage(url)
	if err != nil {
		return "", err
	}

	// Calculate hash of uploadedFile data
	hashValue := hash(uploadedFileData)

	// Check if hash already exists in the database
	url, err = uploadedFileExists(hashValue)

	if err != nil {
		// Save uploadedFile to the database
		url, err = saveImageToDB(uploadedFileData, contentType, hashValue)
		if err != nil {
			return "", err
		}
		return url, nil
	} else {
		return url, nil
	}
}

func downloadImage(url string) ([]byte, string, error) {
	// Download the image from the provided URL
	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP request failed with status: %v", response.Status)
	}

	// Check content type
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", fmt.Errorf("invalid content type: %v", contentType)
	}

	// Read image data
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	return imageData, contentType, nil
}

func uploadedFileExists(hashValue string) (string, error) {
	// Check if the hash value already exists in the database
	var uploadedFile model.UploadedFile
	result := db.DB.Where("hash = ?", hashValue).First(&uploadedFile)

	if result.Error == nil {
		baseAppUrl := os.Getenv("BASE_APP_URL")
		fileUrl := fmt.Sprintf("%s/files/%d", baseAppUrl, uploadedFile.ID)
		return fileUrl, nil
	} else {
		return "", errors.New("file does not exist")
	}
}

func saveImageToDB(uploadedFileData []byte, contentType string, hashValue string) (string, error) {
	uploadedFile := model.UploadedFile{
		ContentData: uploadedFileData,
		ContentType: contentType,
		Hash:        hashValue,
	}

	result := db.DB.Create(&uploadedFile)
	if result.Error != nil {
		return "", result.Error
	}

	baseAppUrl := os.Getenv("BASE_APP_URL")
	fileUrl := fmt.Sprintf("%s/files/%d", baseAppUrl, uploadedFile.ID)
	return fileUrl, nil
}

func hash(data []byte) string {
	// Calculate hash of the data using FNV-1a
	hasher := fnv.New32a()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
