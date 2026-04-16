package model

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestTextToEmbedding(t *testing.T) {
	if os.Getenv("RUN_OPENAI_INTEGRATION_TESTS") != "1" {
		t.Skip("skipping OpenAI integration test; set RUN_OPENAI_INTEGRATION_TESTS=1 to enable")
	}
	config, err := initConf("./config.yaml")
	if err != nil {
		panic(err)
	}
	InitDB(config.Dns, config.OpenaiAPIKey)
	embedding, err := textToEmbedding("Hello, world!")
	if err != nil {
		t.Fatalf("failed to convert content to embedding: %v", err)
	}
	t.Logf("embedding: %v", embedding)
}

type Config struct {
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
}

func initConf(confPath string) (*Config, error) {
	config := &Config{}
	dataBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(dataBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func TestSearchSimilarVectors(t *testing.T) {
	if os.Getenv("RUN_OPENAI_INTEGRATION_TESTS") != "1" {
		t.Skip("skipping OpenAI integration test; set RUN_OPENAI_INTEGRATION_TESTS=1 to enable")
	}
	config, err := initConf("./config.yaml")
	if err != nil {
		panic(err)
	}
	InitDB(config.Dns, config.OpenaiAPIKey)
	TrustModel := []string{}
	IdentityRegistry := []string{}
	ChainID := []string{}
	Skills := []string{}
	X402Support := false
	Active := true
	HaveFeedback := false
	limit := 10
	threshold := 0.5
	filters := &VectorSearchFilters{
		TrustModel:       &TrustModel,
		IdentityRegistry: &IdentityRegistry,
		ChainID:          &ChainID,
		Skills:           &Skills,
		X402Support:      &X402Support,
		Active:           &Active,
		HaveFeedback:     &HaveFeedback,
	}

	desc := "transfer token from one address to another"
	vectors, err := SearchSimilarVectors(desc, limit, threshold, filters)
	if err != nil {
		t.Fatalf("failed to search similar vectors: %v", err)
	}
	t.Logf("vectors: %v", vectors)
}
