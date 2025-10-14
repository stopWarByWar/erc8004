package processor

import (
	"agent_identity/logger"
	"agent_identity/model"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func TestProcessor(t *testing.T) {
	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_logger, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	config, err := initConf("./config.yaml")
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns)

	idx := NewCreateAgentProcessor(config.IdentityAddr, config.ReputationAddr, config.ValidationAddr, config.CommentAddr, config.RpcURL, config.FetchBlockInterval, config.StartBlock, _logger)
	idx.Process()
}

type Config struct {
	IdentityAddr       string `yaml:"identity_addr"`
	ReputationAddr     string `yaml:"reputation_addr"`
	ValidationAddr     string `yaml:"validation_addr"`
	CommentAddr        string `yaml:"comment_addr"`
	RpcURL             string `yaml:"rpc_url"`
	FetchBlockInterval int64  `yaml:"fetch_block_interval"`
	Dns                string `yaml:"dns"`
	StartBlock         uint64 `yaml:"start_block"`
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
