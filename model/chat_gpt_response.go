// model/chat_gpt_response.go
package model

import (
	"github.com/jinzhu/gorm"
)

type ChatGPTResponse struct {
	gorm.Model
	Engine   string `json:"engine"`
	Prompt   string `json:"prompt"`
	Answer   string `gorm:"type:text" json:"answer"`
	ImageURL string `json:"image_url"`
	Success  bool   `json:"success"`
}

// TableName sets the table name for the ChatGPTResponse model.
func (ChatGPTResponse) TableName() string {
	return "chat_gpt_responses"
}
