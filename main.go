package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func getGeneratedResponse(prompt string) string {
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
		log.Fatalf("Failed to marshal json: %v", err)
	}

	//Create Http request struct with request method, endpoint and request body
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", API_URL, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Failed to create http request struct: %v", err)
	}

	//Add necessary headers, including the API key for authorization
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	//Execute http request to chatGPT
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get http response: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %v", res.Status)
	}

	//Decode http response body into Go struct
	defer res.Body.Close()
	ccRes := &chatCompletionResponse{}
	err = json.NewDecoder(res.Body).Decode(ccRes)
	if err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	if len(ccRes.Choices) == 0 {
		log.Fatal("No choices returned from chatGPT")
	}

	//Return generated text from chatGPT
	return ccRes.Choices[0].Message.Content
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
	fmt.Println(getGeneratedResponse(prompt))

	fmt.Println("")
	fmt.Println("")
}
