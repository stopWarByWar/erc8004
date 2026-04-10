package config

import (
	"agent_identity/model"
	"errors"
	"fmt"
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

type FilterInfoAmount struct {
	Name   string `json:"name" yaml:"name"`
	Amount int64  `json:"amount" yaml:"amount"`
}

type FilterStatusInfo struct {
	Active       int64 `json:"active" yaml:"active"`
	HaveFeedback int64 `json:"have_feedback" yaml:"have_feedback"`
}

type FilterInfo struct {
	Networks    []ChainInfo        `yaml:"networks"`
	TrustModels []FilterInfoAmount `yaml:"trust_models"`
	Status      FilterStatusInfo   `yaml:"status"`
	X402Support int64              `yaml:"x402_support"`
	Skills      []FilterInfoAmount `yaml:"skills"`
}

var filterInfo = FilterInfo{}

var RegisterMap = make(map[string]map[string]ContractInfo)
var ChainMap = make(map[string]ChainInfo)
var ValidationRegistryMap = make(map[string]map[string]ContractInfo)

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

	for _, register := range config.RegisterList {
		register.IdentityAddress = common.HexToAddress(register.IdentityAddress).String()
		register.ReputationAddress = common.HexToAddress(register.ReputationAddress).String()
		register.ValidationAddress = common.HexToAddress(register.ValidationAddress).String()
	}

	for _, chain := range config.ChainList {
		ChainMap[chain.ChainId] = chain
	}

	for _, register := range config.RegisterList {
		if RegisterMap[register.ChainId] == nil {
			RegisterMap[register.ChainId] = make(map[string]ContractInfo)
		}
		RegisterMap[register.ChainId][register.IdentityAddress] = register
		if ValidationRegistryMap[register.ChainId] == nil {
			ValidationRegistryMap[register.ChainId] = make(map[string]ContractInfo)
		}
		ValidationRegistryMap[register.ChainId][register.ValidationAddress] = register
	}
	return nil
}

func GetChainInfoMap() map[string]ChainInfo {
	return ChainMap
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
	ChainMap[chainId] = chain
}
func GetContractsDeployerInfo(ChainID string, RegistryAddress string) ContractInfo {
	register, ok := RegisterMap[ChainID][RegistryAddress]
	if !ok {
		return ContractInfo{}
	}
	return register
}
func GetIdentityAddressByValidationAddress(chainID, ValidationAddress string) string {
	registry, ok := ValidationRegistryMap[chainID][ValidationAddress]
	if !ok {
		return ""
	}
	return registry.IdentityAddress
}

func UpdateFilterInfo() error {
	var newFilterInfo = FilterInfo{}
	var chainIds []string
	for chainId := range ChainMap {
		chainIds = append(chainIds, chainId)
	}
	agentAmounts, err := model.GetAgentAmountForEachChain(chainIds)
	if err != nil {
		return fmt.Errorf("failed to get agent amount for each chain: %w", err)
	}

	var chainInfos []ChainInfo
	for chainId, agentAmount := range agentAmounts {
		chainInfo, ok := ChainMap[chainId]
		if !ok {
			continue
		}
		chainInfo.AgentAmount = uint64(agentAmount)
		chainInfos = append(chainInfos, chainInfo)
		SetChainAgentAmount(chainId, int64(agentAmount))
	}
	newFilterInfo.Networks = chainInfos

	trustModelAmounts, err := model.GetAgentAmountForEachTrustModel()
	if err != nil {
		return fmt.Errorf("failed to get agent amount for each trust model: %w", err)
	}
	var trustModelAmountInfos []FilterInfoAmount
	for _, trustModelAmount := range trustModelAmounts {
		trustModelAmountInfos = append(trustModelAmountInfos, FilterInfoAmount{
			Name:   trustModelAmount.Name,
			Amount: trustModelAmount.Amount,
		})
	}
	newFilterInfo.TrustModels = trustModelAmountInfos

	activeAgentAmount, err := model.GetActiveAgentAmount()
	if err != nil {
		fmt.Println("failed to get active agent amount", err)
		return fmt.Errorf("failed to get active agent amount: %w", err)
	}
	newFilterInfo.Status.Active = activeAgentAmount

	feedbackCount, err := model.GetAgentAmountWithFeedback()
	if err != nil {
		fmt.Println("failed to get agent amount with feedback", err)
		return fmt.Errorf("failed to get agent amount with feedback: %w", err)
	}

	newFilterInfo.Status.HaveFeedback = feedbackCount

	x402SupportAgentAmount, err := model.GetX402SupportAgentAmount()
	if err != nil {
		fmt.Println("failed to get x402 support agent amount", err)
		return fmt.Errorf("failed to get x402 support agent amount: %w", err)
	}
	newFilterInfo.X402Support = x402SupportAgentAmount

	skillAmounts, err := model.GetAgentAmountForEachSkill(50)
	if err != nil {
		fmt.Println("failed to get agent amount for each skill", err)
		return fmt.Errorf("failed to get agent amount for each skill: %w", err)
	}
	var skillInfos []FilterInfoAmount
	for _, skillAmount := range skillAmounts {
		skillInfos = append(skillInfos, FilterInfoAmount{
			Name:   skillAmount.Name,
			Amount: skillAmount.Amount,
		})
	}
	newFilterInfo.Skills = skillInfos
	filterInfo = newFilterInfo
	return nil
}

type IndexerConfig struct {
	Name            string `yaml:"name"`
	RpcURL          string `yaml:"rpc_url"`
	Dns             string `yaml:"dns"`
	OpenaiAPIKey    string `yaml:"openai_api_key"`
	CoingeckoAPIKey string `yaml:"coingecko_api_key"`
	Reputation      struct {
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
	Commerce struct {
		Addr               string `yaml:"addr"`
		FetchBlockInterval int64  `yaml:"fetch_block_interval"`
		StartBlock         uint64 `yaml:"start_block"`
		Run                bool   `yaml:"run"`
	} `yaml:"commerce"`
}

func GetFilterInfo() FilterInfo {
	return filterInfo
}
