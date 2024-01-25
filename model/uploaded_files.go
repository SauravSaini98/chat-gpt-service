// model/chat_gpt_response.go
package model

import (
	"github.com/jinzhu/gorm"
)

type UploadedFile struct {
	gorm.Model
	Hash        string `gorm:"uniqueIndex" json:"hash"`
	ContentData []byte `json:"content_data"`
	ContentType string `gorm:"uniqueIndex" json:"content_type"`
}

// TableName sets the table name for the Image model.
func (UploadedFile) TableName() string {
	return "uploaded_files"
}
