// model/chat_gpt_response.go
package model

import (
	"chat-gpt-service/db"
	"chat-gpt-service/helper"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UploadedFile struct {
	gorm.Model
	Hash        string `gorm:"uniqueIndex" json:"hash"`
	FileURL     string `json:"file_url"`
	ContentType string `gorm:"uniqueIndex" json:"content_type"`
}

// TableName sets the table name for the Image model.
func (UploadedFile) TableName() string {
	return "uploaded_files"
}

func (a *UploadedFile) BuildFileUrl() string {
	// Access fields of the specific record using 'a'
	url, err := helper.BuildFileUrl(a.FileURL)
	fmt.Println("url", url, "err", err)

	if err != nil {
		return a.FileURL
	}
	return url
}

func CreateUploadedFile(fileUrl string, contentType string, hashValue string) (UploadedFile, error) {
	uploadedFile := UploadedFile{
		FileURL:     fileUrl,
		ContentType: contentType,
		Hash:        hashValue,
	}

	result := db.DB.Create(&uploadedFile)
	if result.Error != nil {
		return UploadedFile{}, result.Error
	}

	return uploadedFile, nil
}
