package logic

import (
	"agent_identity/config"
	"agent_identity/server/api/types"
	"fmt"
)

func GetNetworkList() ([]types.NetworkResponse, error) {
	networkList := make([]types.NetworkResponse, 0)
	for _, chain := range config.ChainList {
		network := types.NetworkResponse{
			ChainId:   chain.ChainId,
			ChainName: chain.ChainName,
			ChainLogo: chain.ChainLogo,
		}

		deployers := make([]types.ContractInfo, 0)
		for _, register := range config.RegisterMap[chain.ChainId] {
			deployers = append(deployers, types.ContractInfo{
				IdentityAddress:       register.IdentityAddress,
				IdentityContractURL:   fmt.Sprintf("%s/address/%s", chain.ScanPrefix, register.IdentityAddress),
				ReputationAddress:     register.ReputationAddress,
				ReputationContractURL: fmt.Sprintf("%s/address/%s", chain.ScanPrefix, register.ReputationAddress),
				ValidationAddress:     register.ValidationAddress,
				ValidationContractURL: fmt.Sprintf("%s/address/%s", chain.ScanPrefix, register.ValidationAddress),
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
