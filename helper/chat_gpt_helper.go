package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

// GPTConfig represents the configuration for a GPT model.
type GPTConfig struct {
	APIKey    string
	Endpoint  string
	Engine    string
	MaxTokens int
}

func GetChatCompleteResponse(question string, engine string, maxTokens int) (string, error) {
	if maxTokens <= 0 {
		maxTokens = 1000
	}

	if engine == "" {
		engine = "gpt-4-1106-preview"
	}

	payload := map[string]interface{}{
		"model": engine,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": question,
			},
		},
		"max_tokens": maxTokens,
	}

	responseMessage, err := chatGptApiResponse(payload)
	if err != nil {
		return "", err
	}
	return responseMessage, nil
}

func GetChatGptVisionResponse(question string, imageUrl string, maxTokens int) (string, error) {

	if maxTokens <= 0 {
		maxTokens = 1000
	}

	payload := map[string]interface{}{
		"model": "gpt-4-vision-preview",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": question,
					},
					{
						"type": "image",
						"image_url": map[string]string{
							"url": imageUrl,
						},
					},
				},
			},
		},
		"max_tokens": maxTokens,
	}

	responseMessage, err := chatGptApiResponse(payload)
	if err != nil {
		return "", err
	}
	return responseMessage, nil
}

func chatGptApiResponse(payload map[string]interface{}) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	endpoint := "https://api.openai.com/v1/chat/completions"
	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode JSON response
	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", err
	}

	// Check if the response object is "chat.completion"
	if object, ok := responseData["object"].(string); !ok || object != "chat.completion" {
		// Update model_object (you need to replace this with your actual update logic)
		// model_object.update(success: false)

		// Raise an error with the error message from the response
		if errorMessage, ok := responseData["error"].(map[string]interface{})["message"].(string); ok {
			return "", errors.New(errorMessage)
		} else {
			return "", errors.New("time out error")
		}
	}

	// Extract the desired information from the decoded JSON
	var messages []interface{}
	if choices, ok := responseData["choices"].([]interface{}); ok {
		for _, choice := range choices {
			if message, ok := choice.(map[string]interface{})["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					messages = append(messages, content)
				}
			}
		}
	}

	// Join the extracted messages
	responseMessage := ""
	if len(messages) > 0 {
		responseMessage = fmt.Sprintf("%v", messages)
	} else {
		return "", errors.New("somthing went wrong")
	}

	return responseMessage, nil
}

type ModelObject struct {
	// Add fields as needed
	JSONResponse map[string]interface{}
}

// SetResponse sets the response based on the message and response type.
func SetJSONResponse(message string, responseType string) interface{} {
	if responseType == "json" {
		jsonResponse := setMultipleJSON(message)
		if jsonResponse == message {
			jsonResponse = setJSONData(message)
		}
		return jsonResponse
	}

	return message
}

// SetJSONData parses the input data and returns a map.
func setJSONData(data string) map[string]interface{} {
	var hash map[string]interface{}

	// Try to parse JSON
	err := json.Unmarshal([]byte(data), &hash)
	if err != nil {
		hash = make(map[string]interface{})

		// If JSON parsing fails, attempt to extract key-value pairs using regex
		re := regexp.MustCompile(`"([^"]+)":\s*"([^"]+)"|\{\s*"([^"]+)":\s*"([^"]+)"\s*\}|"([^"]+)":\s*([\d\.]+)`)
		matches := re.FindAllStringSubmatch(data, -1)

		for _, match := range matches {
			for i := 1; i < len(match); i += 2 {
				if match[i] != "" {
					hash[match[i]] = match[i+1]
				}
			}
		}
	}

	return hash
}

// SetMultipleJSON extracts JSON data enclosed in triple backticks and parses it.
func setMultipleJSON(data string) interface{} {
	answer := extractJSONFromBackticks(data)
	if answer != "" {
		var result interface{}
		err := json.Unmarshal([]byte(answer), &result)
		if err == nil {
			return result
		}
	}

	return data
}

// extractJSONFromBackticks extracts JSON data enclosed in triple backticks.
func extractJSONFromBackticks(data string) string {
	re := regexp.MustCompile("```json\\n(.*?)\\n```")
	matches := re.FindStringSubmatch(data)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
