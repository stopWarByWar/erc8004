package processor

import (
	"agent_identity/logger"
	"agent_identity/model"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func TestReputationProcessor(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in -short")
	}
	if v := os.Getenv("RUN_INTEGRATION"); v != "1" {
		t.Skip("set RUN_INTEGRATION=1 to run integration test")
	}

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

	model.InitDB(config.Dns, config.OpenaiAPIKey)

	ethClient, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		panic(err)
	}

	identityExecBlockChan := make(chan uint64, 10)
	identityExecBlockChan <- 1000000000000000
	processor := NewReputationProcessor(config.Reputation.Addr, config.Identity.Addr, ethClient, config.Reputation.FetchBlockInterval, config.Reputation.StartBlock, _logger, identityExecBlockChan)
	processor.Process()
}
