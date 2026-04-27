package main

import (
	"agent_identity/config"
	"agent_identity/helper"
	"agent_identity/logger"
	"agent_identity/model"
	"agent_identity/server/api"
	apiUtils "agent_identity/server/api/utils"
	"flag"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var configFile = flag.String("f", "./conf.yaml", "the config file")

func main() {
	flag.Parse()

	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_logger, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	_config, err := initConf(*configFile)
	if err != nil {
		panic(err)
	}

	// config.yaml is at project_root/config/config.yaml
	// *configFile is project_root/server/api/logic/config.yaml
	// Go up 3 levels: server/api/logic -> server/api -> server -> project_root
	config.Init(filepath.Join(filepath.Dir(*configFile), "..", "config", "config.yaml"))

	helper.InitHelper(_config.S3Region, _config.S3BucketName, _config.S3AccessKey, _config.S3SecretKey)

	model.InitDB(_config.Dns, _config.OpenaiAPIKey)
	api.InitRouter(_logger, apiUtils.RuntimeConfig{Mock: _config.Mock, FeedbackMock: _config.FeedbackMock})

	api.Run([]string{"*"}, _config.Port)
}

type Config struct {
	Dns          string `yaml:"dns"`
	Port         string `yaml:"port"`
	Mock         bool   `yaml:"mock"`
	FeedbackMock bool   `yaml:"feedback_mock"`
	S3Region     string `yaml:"s3_region"`
	S3BucketName string `yaml:"s3_bucket_name"`
	S3AccessKey  string `yaml:"aws_access_key_id"`
	S3SecretKey  string `yaml:"aws_secret_access_key"`
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
