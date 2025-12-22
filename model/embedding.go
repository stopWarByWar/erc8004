package model

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

var openAIClient *openai.Client
var ctx = context.Background()

func textToEmbedding(text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}

	resp, err := openAIClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call embedding API: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("embedding result is empty")
	}

	return resp.Data[0].Embedding, nil
}
