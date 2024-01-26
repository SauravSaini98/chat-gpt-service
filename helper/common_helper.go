package helper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// check and set default ENV
func GetEnvWithDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// common function to build s3 url
func BuildFileUrl(fileUrl string) (string, error) {
	appEnvironment := GetEnvWithDefault("APP_ENV", "development")
	baseUrl := ""
	switch strings.ToLower(appEnvironment) {
	case "production":
		// Call the function to save to AWS S3
		baseUrl = GetEnvWithDefault("S3_HOST_URL", "http://localhost:8080/")
	case "development":
		// Call the function to save locally
		baseUrl = GetEnvWithDefault("BASE_APP_URL", "http://localhost:8080/")
	default:
		fmt.Println("Invalid environment specified.")
	}
	if baseUrl == "" {
		return "", errors.New("Invalid App Environment")
	}

	return fmt.Sprintf("%s/%s", baseUrl, fileUrl), nil
}

// Function to save to AWS S3
func SaveFileToS3(data []byte, key string) (string, error) {
	// Create a new AWS session
	awsRegion := GetEnvWithDefault("S3_AWS_REGION", "ap-south-1")
	bucketName := GetEnvWithDefault("S3_BUCKET_NAME", "test-app")
	key = fmt.Sprintf("uploads/%s", key)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return "", err
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Upload the file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"), // Optional: Set the ACL for the object
	})
	if err != nil {
		return "", err
	}

	// Construct the URL based on the S3 bucket and key
	return key, nil
}

// Function to download an file from a URL
func DownloadFileFromUrl(url string) ([]byte, string, error) {
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
	// Read image data
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	return imageData, contentType, nil
}

// Function to save locally and return the file path
func SaveFileToLocal(data []byte, fileName string) (string, error) {
	// Specify the directory where you want to save the files
	baseDirectory := "public/uploads/chat-gpt/"

	// Create the full file path
	filePath := filepath.Join(baseDirectory, fileName)

	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create or open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	filePathWithoutPublic := strings.TrimPrefix(filePath, "public/")

	return filePathWithoutPublic, nil
}

// Function to derive file extension from content type
func GetExtensionFromContentType(contentType string) (string, error) {
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		return "", fmt.Errorf("unable to determine file extension for content type %s", contentType)
	}
	return exts[0][1:], nil // Remove the dot from the extension
}

// Function to generate a unique file name (you can customize this based on your needs)
func GenerateUniqueFileName() string {
	// You can use a library or logic to generate a unique file name here
	// For simplicity, using a constant prefix and timestamp in this example
	return "file_" + currentTimeStamp()
}

// Function to get the current timestamp (you can use your preferred timestamp logic)
func currentTimeStamp() string {
	timestamp := time.Now().Unix()
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Sprintf("%s_%d", "file", timestamp)
	}
	uniqueCode := randomUUID.String()
	return fmt.Sprintf("%s_%d", uniqueCode, timestamp)
}
