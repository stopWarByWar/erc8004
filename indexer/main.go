package main

import (
	"agent_identity/indexer/processor"
	"agent_identity/logger"
	"agent_identity/model"
	"context"
	"flag"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var configFile = flag.String("f", "./conf.yaml", "the config file")

func main() {
	flag.Parse()

	logConf := &logger.Config{
		Level:        logrus.DebugLevel,
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

	ethClient, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		panic(err)
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns)

	// createAgentIdx := processor.NewCreateAgentProcessor(config.Identity.Addr, ethClient, config.Identity.FetchBlockInterval, config.Identity.StartBlock, _logger)
	// go createAgentIdx.Process()

	// reputationIdx := processor.NewReputationProcessor(config.Reputation.Addr, ethClient, config.Reputation.FetchBlockInterval, config.Reputation.StartBlock, _logger)
	// go reputationIdx.Process()

	commentIdx := processor.NewCommentProcessor(chainID.String(), config.Comment.CommentSchemaID, config.Comment.StartBlock, config.Comment.Limit, config.Comment.FetchBlockInterval, _logger)
	commentIdx.Process()
}

type Config struct {
	RpcURL     string `yaml:"rpc_url"`
	Dns        string `yaml:"dns"`
	Reputation struct {
		Addr               string `yaml:"addr"`
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
	} `yaml:"reputation"`

	Identity struct {
		Addr               string `yaml:"addr"`
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
	} `yaml:"identity"`
	Comment struct {
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
		Limit              int    `yaml:"limit"`
		CommentSchemaID    string `yaml:"comment_schema_id"`
	} `yaml:"comment"`
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
