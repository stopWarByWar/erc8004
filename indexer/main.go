package main

import (
	"agent_identity/indexer/processor"
	"agent_identity/logger"
	"agent_identity/model"
	"flag"
	"os"

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

	config, err := initConf(*configFile)
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns)

	// createAgentIdx := processor.NewCreateAgentProcessor(config.IdentityAddr, config.ReputationAddr, config.ValidationAddr, config.CommentAddr, config.RpcURL, config.FetchBlockInterval, config.StartBlock, _logger)
	// go createAgentIdx.Process()

	commentIdx := processor.NewCommentProcessor(config.StartBlock, config.FetchCommentLimit, config.FetchCommentInterval, config.CommentAttestor, _logger)
	commentIdx.Process()
}

type Config struct {
	IdentityAddr         string `yaml:"identity_addr"`
	ReputationAddr       string `yaml:"reputation_addr"`
	ValidationAddr       string `yaml:"validation_addr"`
	CommentAddr          string `yaml:"comment_addr"`
	CommentAttestor      string `yaml:"comment_attestor"`
	RpcURL               string `yaml:"rpc_url"`
	FetchBlockInterval   int64  `yaml:"fetch_block_interval"`
	FetchCommentInterval int64  `yaml:"fetch_comment_interval"`
	FetchCommentLimit    int    `yaml:"fetch_comment_limit"`
	Dns                  string `yaml:"dns"`
	StartBlock           uint64 `yaml:"start_block"`
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
