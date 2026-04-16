package main

import (
	"agent_identity/config"
	"agent_identity/indexer/processor"
	"agent_identity/logger"
	"agent_identity/model"
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

	model.InitDB(config.Dns, config.OpenaiAPIKey)

	var reputationIdx *processor.ReputationProcessor
	var validationIdx *processor.ValidationRegistryProcessor
	var identityIdx *processor.IdentityProcessor

	// 创建channel用于identityIdx和reputationIdx、validationIdx之间的通信
	// 使用Fan-out模式：identityIdx发送到一个channel，然后分发到多个接收者
	var identityExecBlockChan chan uint64
	var reputationExecBlockChan chan uint64
	var validationExecBlockChan chan uint64

	// 如果identity和至少一个processor运行，创建主channel
	if config.Identity.Run && (config.Reputation.Run || config.Validation.Run) {
		identityExecBlockChan = make(chan uint64, 10)
		if config.Reputation.Run {
			reputationExecBlockChan = make(chan uint64, 10)
		}
		if config.Validation.Run {
			validationExecBlockChan = make(chan uint64, 10)
		}

		// Fan-out goroutine: 从identityIdx接收，同时发送给reputationIdx和validationIdx
		go func() {
			for execBlock := range identityExecBlockChan {
				// 同时发送给所有需要接收的processor
				if reputationExecBlockChan != nil {
					select {
					case reputationExecBlockChan <- execBlock:
					default:
						// channel已满，跳过本次发送（非阻塞）
					}
				}
				if validationExecBlockChan != nil {
					select {
					case validationExecBlockChan <- execBlock:
					default:
						// channel已满，跳过本次发送（非阻塞）
					}
				}
			}
		}()
	}

	if config.Reputation.Run {
		reputationIdx = processor.NewReputationProcessor(config.Reputation.Addr, config.Identity.Addr, ethClient, config.Reputation.FetchBlockInterval, config.Reputation.StartBlock, _logger, reputationExecBlockChan)
		go reputationIdx.Process()
	}

	if config.Validation.Run {
		validationIdx = processor.NewValidationRegistryProcessor(config.Validation.Addr, config.Identity.Addr, ethClient, config.Validation.FetchBlockInterval, config.Validation.StartBlock, _logger, validationExecBlockChan)
		go validationIdx.Process()
	}

	if config.Identity.Run {
		var sendChan chan<- uint64
		if identityExecBlockChan != nil {
			sendChan = identityExecBlockChan
		}
		identityIdx = processor.NewCreateAgentProcessor(config.Identity.Addr, ethClient, config.Identity.FetchBlockInterval, config.Identity.StartBlock, _logger, sendChan)
		go identityIdx.Process()
	}

	if config.Commerce.Run {
		commerceIdx := processor.NewCommerceProcessor(config.Commerce.Addr, ethClient, config.Commerce.FetchBlockInterval, config.Commerce.StartBlock, _logger)
		go commerceIdx.Process()
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
	config.Validation.Addr = common.HexToAddress(config.Validation.Addr).String()
	config.Commerce.Addr = common.HexToAddress(config.Commerce.Addr).String()
	config.Comment.CommentSchemaID = common.HexToHash(config.Comment.CommentSchemaID).String()

	processor.SetCoingeckoAPIKey(config.CoingeckoAPIKey)

	return config, nil
}
