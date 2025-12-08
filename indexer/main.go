package main

import (
	"agent_identity/config"
	"agent_identity/indexer/processor"
	"agent_identity/logger"
	"agent_identity/model"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var configFile = flag.String("f", "./config/testnet/base_sepolia.yaml", "the config file")

func main() {
	flag.Parse()

	config, err := initConf(*configFile)
	if err != nil {
		panic(err)
	}

	ethClient, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		panic(err)
	}

	logConf := &logger.Config{
		Level:        logrus.DebugLevel,
		ReportCaller: true,
		FilePath:     fmt.Sprintf("./allLogger/indexer/%s", config.Name),
	}

	_logger, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns)

	if config.Reputation.Run {
		reputationIdx := processor.NewReputationProcessor(config.Reputation.Addr, ethClient, config.Reputation.FetchBlockInterval, config.Reputation.StartBlock, _logger)
		go reputationIdx.Process()
	}

	if config.Identity.Run {
		identityIdx := processor.NewCreateAgentProcessor(config.Identity.Addr, ethClient, config.Identity.FetchBlockInterval, config.Identity.StartBlock, _logger)
		go identityIdx.Process()
	}

	if config.Comment.Run {
		commentIdx := processor.NewCommentProcessor(chainID.String(), config.Comment.CommentSchemaID, config.Comment.StartBlock, config.Comment.Limit, config.Comment.FetchBlockInterval, _logger)
		go commentIdx.Process()
	}

	select {}
}

func initConf(confPath string) (*config.IndexerConfig, error) {
	config := &config.IndexerConfig{}
	dataBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dataBytes, config)
	if err != nil {
		return nil, err
	}

	config.Reputation.Addr = common.HexToAddress(config.Reputation.Addr).String()
	config.Identity.Addr = common.HexToAddress(config.Identity.Addr).String()
	config.Comment.CommentSchemaID = common.HexToHash(config.Comment.CommentSchemaID).String()
	return config, nil
}
