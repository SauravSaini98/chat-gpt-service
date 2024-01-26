package uploader

import (
	"chat-gpt-service/db"
	"chat-gpt-service/helper"
	"chat-gpt-service/model"

	"errors"
	"fmt"
	"hash/fnv"

	"strings"

	"github.com/jinzhu/gorm"
)

func CheckAndSaveImage(fileUrl string) (model.UploadedFile, error) {
	uploadedFileData, contentType, err := helper.DownloadFileFromUrl(fileUrl)

	if err != nil {
		return model.UploadedFile{}, err
	}

	if !strings.HasPrefix(contentType, "image/") {
		return model.UploadedFile{}, fmt.Errorf("invalid content type: %v", contentType)
	}

	// Calculate hash of uploadedFile data
	hashValue := hash(uploadedFileData)

	// Check if hash already exists in the database
	uploadedFile, err := uploadedFileExists(hashValue)

	if err != nil {
		// Save uploadedFile to the database
		url, err := FileUploader(uploadedFileData, contentType)
		if err != nil {
			return model.UploadedFile{}, err
		}
		uploadedFile, err = model.CreateUploadedFile(url, contentType, hashValue)
		if err != nil {
			return model.UploadedFile{}, err
		}
		return uploadedFile, nil
	} else {
		return uploadedFile, nil
	}
}

func uploadedFileExists(hashValue string) (model.UploadedFile, error) {
	// Check if the hash value already exists in the database
	var uploadedFile model.UploadedFile
	result := db.DB.Where("hash = ?", hashValue).First(&uploadedFile)

	if result.Error == nil {
		return uploadedFile, nil
	} else if gorm.IsRecordNotFoundError(result.Error) {
		return model.UploadedFile{}, errors.New("file does not exist")
	} else {
		return model.UploadedFile{}, result.Error
	}
}

func hash(data []byte) string {
	// Calculate hash of the data using FNV-1a
	hasher := fnv.New32a()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
