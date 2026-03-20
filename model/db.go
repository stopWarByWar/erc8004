package model

import (
	"os"

	openai "github.com/sashabaranov/go-openai"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(dns, openaiAPIKey string) {
	var err error
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		cfg := openai.DefaultConfig(openaiAPIKey)
		cfg.BaseURL = baseURL
		openAIClient = openai.NewClientWithConfig(cfg)
		return
	}

	openAIClient = openai.NewClient(openaiAPIKey)

}
