package logic

import (
	"agent_identity/config"
	"agent_identity/server/api/types"
	"fmt"
)

func GetNetworkList() ([]types.NetworkResponse, error) {
	networkList := make([]types.NetworkResponse, 0)
	chainInfoMap := config.GetChainInfoMap()
	for chainId, info := range chainInfoMap {
		network := types.NetworkResponse{
			ChainId:     chainId,
			ChainName:   info.ChainName,
			ChainLogo:   info.ChainLogo,
			AgentAmount: info.AgentAmount,
		}

		deployers := make([]types.ContractInfo, 0)
		for _, register := range config.RegisterMap[chainId] {
			deployers = append(deployers, types.ContractInfo{
				IdentityAddress:       register.IdentityAddress,
				IdentityContractURL:   fmt.Sprintf("%s/address/%s", info.ScanPrefix, register.IdentityAddress),
				ReputationAddress:     register.ReputationAddress,
				ReputationContractURL: fmt.Sprintf("%s/address/%s", info.ScanPrefix, register.ReputationAddress),
				ValidationAddress:     register.ValidationAddress,
				ValidationContractURL: fmt.Sprintf("%s/address/%s", info.ScanPrefix, register.ValidationAddress),
				Deployer:              register.Deployer,
				Description:           register.Description,
				LogoURL:               register.LogoURL,
			})
		}
		network.ContractInfo = deployers
		networkList = append(networkList, network)
	}
	return networkList, nil
}

func GetAgentAmountForEachChain() ([]config.ChainInfo, uint64) {
	chainInfoMap := config.GetChainInfoMap()
	chainInfos := make([]config.ChainInfo, 0)
	total := uint64(0)
	for _, info := range chainInfoMap {
		chainInfos = append(chainInfos, info)
		total += info.AgentAmount
	}

	return chainInfos, total
}
