package processor

import (
	agentcard "agent_identity/agentCard"
	"agent_identity/logger"
	"agent_identity/model"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

	model.InitDB(config.Dns, config.OpenaiAPIKey)

	ethClient, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		panic(err)
	}

	sendChan := make(chan uint64, 10)
	idx := NewCreateAgentProcessor(common.HexToAddress(config.Identity.Addr).String(), ethClient, config.Identity.FetchBlockInterval, config.Identity.StartBlock, _logger, sendChan)
	idx.Process()

	// idx.setAgentCardInserted()
}

type Config struct {
	RpcURL       string `yaml:"rpc_url"`
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
	Reputation   struct {
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

func Test_GetTokenURL(t *testing.T) {
	tokenURL := "https://baspublic.s3.ap-southeast-1.amazonaws.com/erc8004/agent_profile/97-0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345-382.json"
	agent, provider, err := agentcard.GetAgentCardFromTokenURL("0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345", "382", tokenURL, "97", "0xA98A5542a1AaB336397d487e32021E0E48BEF717", 1731196800)
	if err != nil {
		panic(err)
	}
	t.Log(agent)
	t.Log(provider)
}
