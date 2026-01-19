package processor

import (
	"agent_identity/logger"
	"agent_identity/model"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func TestValidationProcessor(t *testing.T) {
	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_logger, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	config, err := initConf("../../config/testnet/eth_sepolica_foundation.yaml")
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns, config.OpenaiAPIKey)

	ethClient, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		panic(err)
	}

	identityExecBlockChan := make(chan uint64, 10)
	processor := NewValidationRegistryProcessor(config.Validation.Addr, config.Identity.Addr, ethClient, config.Validation.FetchBlockInterval, config.Validation.StartBlock, _logger, identityExecBlockChan)
	go func() {
		identityExecBlockChan <- 1000000000000000
	}()
	processor.Process()
}
