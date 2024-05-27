package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Request to chatGPT API
type chatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response from chatGPT API
type chatCompletionResponse struct {
	Id                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []choice `json:"choices"`
	Usage             usage    `json:"usage"`
}

type choice struct {
	Index        int     `json:"index"`
	Message      message `json:"message"`
	Logprobs     bool    `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Endpoint
const API_URL = "https://api.openai.com/v1/chat/completions"

func getGeneratedResponse(prompt string) (string, error) {
	//Set gpt model and received prompt
	ccReq := &chatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: prompt},
		},
	}

	//Marshal Go struct into Json
	jsonData, err := json.Marshal(ccReq)
	if err != nil {
		log.Printf("Failed to Marshal: %v", err)
		return "", err
	}

	//Create Http request struct with request method, endpoint and request body
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", API_URL, bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("Failed to create http request struct: %v", err)
		return "", err
	}

	//Add necessary headers, including the API key for authorization
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("OPENAI_API_KEY environment variable is not set")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	//Execute http request to chatGPT and get response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get http response: %v", err)
		return "", err
	}

	//Check if http status code is ok
	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %v", res.Status)
		return "", err
	}

	//Read http response body
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		return "", err
	}

	//Unmarshal Json response into Go struct
	ccRes := &chatCompletionResponse{}
	err = json.Unmarshal(body, ccRes)
	if err != nil {
		log.Printf("Failed to Unmarshal: %v", err)
		return "", err
	}

	if len(ccRes.Choices) == 0 {
		log.Println("No choices returned from chatGPT")
		return "", err
	}

	//Return generated text from chatGPT
	return ccRes.Choices[0].Message.Content, nil
}

func main() {
	words := [3]string{"nonchalant", "reckon", "appalled"}
	prompt := fmt.Sprintf("Please create an English example sentence using following words: %s, %s, %s",
		words[0], words[1], words[2])

	fmt.Println("")
	fmt.Println("")

	fmt.Println("++++++ Prompt ++++++")
	fmt.Println(prompt)

	fmt.Println("")
	fmt.Println("")

	fmt.Println("++++++ Generated response ++++++")
	response, err := getGeneratedResponse(prompt)
	if err != nil {
		log.Fatalf("Failed to get generated response from ChatGPT API: %v", err)
	}
	fmt.Println(response)

	fmt.Println("")
	fmt.Println("")
}
