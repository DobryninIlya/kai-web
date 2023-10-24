package model

type ChatGPTAnswer struct {
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Message      struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Created int    `json:"created"`
	Id      string `json:"id"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error string `json:"error"`
}

type MessagesChatGPT struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTNewsParams struct {
	Header      string `json:"header"`
	Description string `json:"description"`
}
