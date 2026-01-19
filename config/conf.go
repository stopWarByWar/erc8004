package config

import (
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v2"
)

type ChainInfo struct {
	ChainId     string `json:"chain_id" yaml:"chain_id"`
	ChainName   string `json:"chain_name" yaml:"chain_name"`
	ChainLogo   string `json:"chain_logo" yaml:"chain_logo"`
	ScanPrefix  string `json:"scan_prefix" yaml:"scan_prefix"`
	AgentAmount uint64 `json:"agent_amount" yaml:"agent_amount"`
}

type ContractInfo struct {
	ChainId           string `json:"chain_id" yaml:"chain_id"`
	IdentityAddress   string `json:"identity_address" yaml:"identity_address"`
	ReputationAddress string `json:"reputation_address" yaml:"reputation_address"`
	ValidationAddress string `json:"validation_address" yaml:"validation_address"`
	Deployer          string `json:"name" yaml:"name"`
	Description       string `json:"description" yaml:"description"`
	LogoURL           string `json:"logo_url" yaml:"logo_url"`
}

type Config struct {
	ChainList    []ChainInfo    `yaml:"chain_list"`
	RegisterList []ContractInfo `yaml:"register_list"`
}

var ChainList = []ChainInfo{}
var RegisterList = []ContractInfo{}
var RegisterMap = make(map[string]map[string]ContractInfo)
var ChainMap = make(map[string]ChainInfo)

// Init 从配置文件加载 ChainList 和 RegisterList
// 如果 configPath 为空，将尝试从多个常见路径加载
func Init(configPath string) error {
	if configPath == "" {
		return errors.New("config path is empty")
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return errors.New("config file not found")
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return errors.New("config file is invalid")
	}

	ChainList = config.ChainList
	RegisterList = config.RegisterList

	for _, register := range RegisterList {
		register.IdentityAddress = common.HexToAddress(register.IdentityAddress).String()
		register.ReputationAddress = common.HexToAddress(register.ReputationAddress).String()
		register.ValidationAddress = common.HexToAddress(register.ValidationAddress).String()
	}

	for _, chain := range ChainList {
		ChainMap[chain.ChainId] = chain
	}

	for _, register := range RegisterList {
		if RegisterMap[register.ChainId] == nil {
			RegisterMap[register.ChainId] = make(map[string]ContractInfo)
		}
		RegisterMap[register.ChainId][register.IdentityAddress] = register
	}
	return nil
}

func GetChainInfo(chainId string) (ChainInfo, bool) {
	chain, ok := ChainMap[chainId]
	if !ok {
		return ChainInfo{}, false
	}
	return chain, true
}

func SetChainAgentAmount(chainId string, amount int64) {
	chain, ok := ChainMap[chainId]
	if !ok {
		return
	}
	chain.AgentAmount = uint64(amount)
}
func GetContractsDeployerInfo(ChainID string, RegistryAddress string) ContractInfo {
	register, ok := RegisterMap[ChainID][RegistryAddress]
	if !ok {
		return ContractInfo{}
	}
	return register
}

type IndexerConfig struct {
	Name         string `yaml:"name"`
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
		Limit              int    `yaml:"lim	it"`
		CommentSchemaID    string `yaml:"comment_schema_id"`
		Run                bool   `yaml:"run"`
	} `yaml:"comment"`
	Validation struct {
		Addr               string `yaml:"addr"`
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
		Run                bool   `yaml:"run"`
	} `yaml:"validation"`
}
