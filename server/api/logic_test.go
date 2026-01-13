package api

import (
	"agent_identity/logger"
	"agent_identity/model"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func initTest() {
	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	config, err := initConf("./config.yaml")
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns, config.OpenaiAPIKey)
}

func TestGetCardResponse(t *testing.T) {
	initTest()
	CardResponse, err := GetCardResponse(5593)
	if err != nil {
		t.Errorf("GetCardResponse error: %v", err)
	}
	fmt.Println(CardResponse)
}

func TestGetCardList(t *testing.T) {
	initTest()
	agents, total, err := GetAgentList(1, 10)
	if err != nil {
		t.Errorf("GetAgentList error: %v", err)
	}
	for _, agent := range agents {
		fmt.Println(agent)
	}
	fmt.Println(total)
}

func TestGetCardResponseByTrustModel(t *testing.T) {
	initTest()
	agents, total, err := GetAgentListByFilter(1, 10, []string{"feedback", "inference-validation"}, []string{"97"})
	if err != nil {
		t.Errorf("GetAgentListByFilter error: %v", err)
	}
	for _, agent := range agents {
		fmt.Println(agent)
	}
	fmt.Println(total)
}

func TestSearchCardResponseBySkill(t *testing.T) {
	initTest()
	agents, total, err := SearchAgentListBySkill("traffic", 1, 10)
	if err != nil {
		t.Errorf("SearchAgentListBySkill error: %v", err)
	}
	for _, agent := range agents {
		fmt.Println(agent)
	}
	fmt.Println(total)
}

type Config struct {
	Dns                string `yaml:"dns"`
	Port               string `yaml:"port"`
	Mock               bool   `yaml:"mock"`
	S3Region           string `yaml:"s3_region"`
	S3BucketName       string `yaml:"s3_bucket_name"`
	AwsAccessKeyId     string `yaml:"aws_access_key_id"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"`
	OpenaiAPIKey       string `yaml:"openai_api_key"`
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
