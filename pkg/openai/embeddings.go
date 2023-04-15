package openai

import (
	"errors"
	"log"
)

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// EmbeddingResponse is the response from a Create embeddings request.
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// SendEmbeddings returns an EmbeddingResponse
func (gpt ChatGPT) Embeddings(input string) (resp *EmbeddingResponse, err error) {
	requestBody := EmbeddingRequest{
		Input: input,
		Model: "text-embedding-ada-002",
	}

	resp = &EmbeddingResponse{}
	err = gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/embeddings", "POST",
		jsonBody,
		requestBody, resp)
	if err != nil {
		log.Println("embeddings err:", err)
		err = errors.New("openai 请求失败")
	}

	return resp, err
}
