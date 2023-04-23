package openai

import (
	"fmt"
	"lark-vkm/internal/initialization"
	"os"
	"strings"
	"testing"
)

func TestSendEmbeddings(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		wantEmbeddingResponse EmbeddingResponse
		wantErr               bool
	}{
		{
			name:                  "1",
			input:                 "alice is a beautiful girl",
			wantEmbeddingResponse: EmbeddingResponse{},
			wantErr:               false,
		},
	}

	client := NewChatGPT(initialization.Config{
		OpenaiApiKeys: strings.Split(os.Getenv("OPENAI_KEY"), ","),
		OpenaiApiUrl:  os.Getenv("API_URL"),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmbeddingResponse, err := client.Embeddings(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEmbeddings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 只要检测有字段 model 就可以了
			if gotEmbeddingResponse.Model == "" {
				t.Errorf("CreateEmbeddings() gotEmbeddingResponse = %v, want %v", gotEmbeddingResponse, tt.wantEmbeddingResponse)
			}
			fmt.Println(gotEmbeddingResponse)
		})
	}
}
