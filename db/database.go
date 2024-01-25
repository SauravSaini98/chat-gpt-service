// db/database.go
package db

import (
	"chat-gpt-service/config"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable search_path=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSchema)

	fmt.Println("dbURI", dbURI)
	var errDB error
	DB, errDB = gorm.Open("postgres", dbURI)
	if errDB != nil {
		fmt.Println("Failed to connect to database")
		return errDB
	}

	return nil
}
