// model/chat_gpt_response.go
package model

import (
	"github.com/jinzhu/gorm"
)

type ChatGPTResponse struct {
	gorm.Model
	Engine         string       `json:"engine"`
	Prompt         string       `json:"prompt"`
	Answer         string       `gorm:"type:text" json:"answer"`
	Success        bool         `json:"success"`
	UploadedFileID uint         `json:"uploaded_file_id"`
	UploadedFile   UploadedFile `gorm:"foreignkey:UploadedFileID"`
}

// TableName sets the table name for the ChatGPTResponse model.
func (ChatGPTResponse) TableName() string {
	return "chat_gpt_responses"
}
