// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"chat-gpt-service/config"
	"chat-gpt-service/controller"
	"chat-gpt-service/db"
	"chat-gpt-service/helper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig()
	port := helper.GetEnvWithDefault("PORT", "8080")
	if err != nil {
		fmt.Println("Failed to load configuration")
		return
	}

	err = db.InitDB(config)
	if err != nil {
		fmt.Println("Failed to initialize the database")
		return
	}

	r := gin.Default()

	r.Use(cors.Default())

	// Serve static files using the gin.Static middleware
	r.Static("/uploads", "./public/uploads")

	// Middleware to enable caching
	// store := persistence.NewInMemoryStore(time.Minute * 5)

	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(validateTokenMiddleware())

	protectedRoutes.POST("/search", controller.GetChatGPTResponseHandler)
	protectedRoutes.POST("/vision", controller.GetChatGPTVisionResponseHandler)
	protectedRoutes.POST("/upload", controller.TestUploadedFileHandler)

	// r.GET("/files/:id", cache.CachePage(store, time.Hour, controller.GetUploadedFileHandler))

	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	if err := r.Run(":" + port); err != nil {
		fmt.Println("Failed to start server")
	}
}

// Validate Auth token for middle ware
func validateTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the value of the Authorization header
		token := c.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		// Check if the token matches the expected token
		if token != os.Getenv("AUTH_TOKEN_HEADER") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort() // Abort further processing
			return
		}

		c.Next()
	}
}
