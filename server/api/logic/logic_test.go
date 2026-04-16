package logic

import (
	"agent_identity/config"
	"agent_identity/helper"
	"agent_identity/logger"
	"agent_identity/model"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func initTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run DB-backed tests")
	}
	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	_config, err := initConf("./config.yaml")
	if err != nil {
		panic(err)
	}
	err = config.Init("../../../config/config.yaml")
	if err != nil {
		panic(err)
	}

	helper.InitHelper(_config.S3Region, _config.S3BucketName, _config.AWSAccessKeyId, _config.AWSSecretAccessKey)

	model.InitDB(_config.Dns, _config.OpenaiAPIKey)
	if err := model.EnsureCommerceActionsSchema(); err != nil {
		panic(err)
	}
	if err := model.EnsureCommerceScoresSchema(); err != nil {
		panic(err)
	}
}

type Config struct {
	Dns                string `yaml:"dns"`
	Port               string `yaml:"port"`
	Mock               bool   `yaml:"mock"`
	S3Region           string `yaml:"s3_region"`
	S3BucketName       string `yaml:"s3_bucket_name"`
	AWSAccessKeyId     string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
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
