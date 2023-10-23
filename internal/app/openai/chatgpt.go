package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/model"
	"net/http"
	"time"
)

const (
	apiEndpoint    = "https://neuroapi.host"
	apiCompletions = apiEndpoint + "/v1/chat/completions"
	NEWS_PROMPT    = `Представь, что ты новостной редактор. 
Ниже я предоставлю текст, тебе необходимо будет выделить из него короткий заголовок до 100 символов и небольшое текстовоепревью (короткое описание) до 210 символов. 
Формат ответа должен быть в json. Ключи json header и description Текст:
`
)

type ChatGPT struct {
	model       string  // model name
	temperature float64 // 0.0 - 2.0
	role        string
	prompt      string
	log         *logrus.Logger
	ctx         context.Context
}

func NewChatGPT(ctx context.Context, logger *logrus.Logger, model string, temperature float64, role string) *ChatGPT {
	return &ChatGPT{
		model:       model,
		temperature: temperature,
		role:        role,
		ctx:         ctx,
		log:         logger,
	}
}

func (c *ChatGPT) WithPrompt(prompt string) {
	c.prompt = prompt
}

func (c *ChatGPT) getPayload(message string) ([]byte, error) {
	var answer struct {
		Model       string                  `json:"model"`
		Temperature float64                 `json:"temperature"`
		Messages    []model.MessagesChatGPT `json:"messages"`
	}
	answer.Model = c.model
	answer.Temperature = c.temperature
	answer.Messages = append(answer.Messages, model.MessagesChatGPT{Role: c.role, Content: c.prompt + message})
	result, err := json.Marshal(answer)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ChatGPT) GenerateAnswer(message string) (model.ChatGPTNewsParams, error) {
	ctx, closeCtx := context.WithTimeout(c.ctx, time.Second*10)
	defer closeCtx()
	payload, err := c.getPayload(message)
	if err != nil {
		c.log.Logf(
			logrus.ErrorLevel,
			"Error while generating payload: %s",
			err.Error(),
		)
		return model.ChatGPTNewsParams{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiCompletions, bytes.NewBuffer(payload))
	if err != nil {
		c.log.Logf(
			logrus.ErrorLevel,
			"Error while creating request: %s",
			err.Error(),
		)
		return model.ChatGPTNewsParams{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 15}
	resp, err := client.Do(req)
	if err != nil {
		c.log.Logf(
			logrus.ErrorLevel,
			"Error while sending request: %s",
			err.Error(),
		)
		return model.ChatGPTNewsParams{}, err
	}
	defer resp.Body.Close()
	var answer model.ChatGPTAnswer
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &answer)
	if err != nil {
		c.log.Logf(
			logrus.ErrorLevel,
			"Error while reading response: %s",
			err.Error(),
		)
		return model.ChatGPTNewsParams{}, err
	}
	content := answer.Choices[0].Message.Content
	var result model.ChatGPTNewsParams
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		c.log.Logf(
			logrus.ErrorLevel,
			"Error while reading response: %s",
			err.Error(),
		)
		return model.ChatGPTNewsParams{}, err
	}

	return result, nil

}
