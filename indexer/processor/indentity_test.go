package processor

import (
	"agent_identity/logger"
	"agent_identity/model"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
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

	ethClient, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		panic(err)
	}

	idx := NewCreateAgentProcessor(config.Identity.Addr, ethClient, config.Identity.FetchBlockInterval, config.Identity.StartBlock, _logger)
	idx.Process()

	// idx.setAgentCardInserted()
}

type Config struct {
	RpcURL     string `yaml:"rpc_url"`
	Dns        string `yaml:"dns"`
	Reputation struct {
		Addr               string `yaml:"addr"`
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
		Run                bool   `yaml:"run"`
	} `yaml:"reputation"`

	Identity struct {
		Addr               string `yaml:"addr"`
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
		Run                bool   `yaml:"run"`
	} `yaml:"identity"`
	Comment struct {
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
		Limit              int    `yaml:"limit"`
		CommentSchemaID    string `yaml:"comment_schema_id"`
		Run                bool   `yaml:"run"`
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
