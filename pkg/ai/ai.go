package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Service struct {
	Token string
	Request
}

type Response struct {
	Id      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created uint                     `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]int           `json:"usage"`
	Error   map[string]string        `json:"error"`
}

type Request struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float32 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             float32 `json:"top_p"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
	PresencePenalty  float32 `json:"presence_penalty"`
}

// CreateAiService gets token string
// and name of file with configuration of OpenAI API
func NewAiService(token string) (*Service, error) {
	configFilename := os.Getenv("OPENAI_CONFIG_FILE")
	if configFilename == "" {
		return nil, fmt.Errorf("OPENAI_CONFIG_FILE env is empty")
	}
	rawData, err := os.ReadFile(configFilename)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file, err: %w", err)
	}
	var requestConfig Request
	if err = json.Unmarshal(rawData, &requestConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshall request config, err: %w", err)
	}

	return &Service{
		Token:   token,
		Request: requestConfig,
	}, nil
}

// MakeRequest add text string to prompt of AiService struct
// and send request to OpenAI API for completion creating
// then parse response and return text of answer
func (s *Service) MakeRequest(text string) (answer string, err error) {
	s.Prompt = text
	body, err := json.Marshal(s.Request)
	if err != nil {
		return "", fmt.Errorf("cannot convert aiService in body format, err: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("cannot create request, err: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot get response, err: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read response body, err: %w", err)
	}
	var responseAI Response
	err = json.Unmarshal(data, &responseAI)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal reponse body, err: %w", err)
	}
	if len(responseAI.Choices) == 0 {
		return "", fmt.Errorf(responseAI.Error["message"])
	}
	return responseAI.Choices[0]["text"].(string), nil
}
