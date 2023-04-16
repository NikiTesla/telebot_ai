package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type AiService struct {
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
func CreateAiService(token string, confFilename string) (*AiService, error) {
	rawData, err := os.ReadFile(confFilename)
	if err != nil {
		logrus.Printf("cannot configure OpenAI API, error: &s", err.Error())
		return nil, err
	}

	var req Request
	err = json.Unmarshal(rawData, &req)
	if err != nil {
		logrus.Printf("cannot unmarshall configuration of AI API, err: %s", err.Error())
		return nil, err
	}

	return &AiService{token, req}, nil
}

// MakeRequest add text string to prompt of AiService struct
// and send request to OpenAI API for completion creating
// then parse response and return text of answer
func (aS *AiService) MakeRequest(text string) (answer string, err error) {
	aS.Prompt = text
	body, err := json.Marshal(aS.Request)
	if err != nil {
		logrus.Printf("can't convert aiService in body format, err: %s", err.Error())
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewReader(body))
	if err != nil {
		logrus.Printf("Cannot create request %s", err.Error())
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+aS.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Printf("Can't get response, err: %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var responseAI Response
	err = json.Unmarshal(data, &responseAI)
	if err != nil {
		logrus.Printf("Cannot unmarshall response body %s", err.Error())
		return "", err
	}

	if len(responseAI.Choices) == 0 {
		return "", errors.New(responseAI.Error["message"])
	}

	return responseAI.Choices[0]["text"].(string), nil
}
