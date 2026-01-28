package logic

import (
	"agent_identity/config"
	"agent_identity/model"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	serverTypes "agent_identity/server/api/types"

	"github.com/ethereum/go-ethereum/common"
)

func formatAgentResponse(agents []*model.Agent) ([]*serverTypes.AgentResponse, error) {
	if len(agents) == 0 {
		return []*serverTypes.AgentResponse{}, nil
	}

	agentUIDs := make([]uint64, 0, len(agents))
	for _, agent := range agents {
		agentUIDs = append(agentUIDs, agent.UID)
	}

	var resp []*serverTypes.AgentResponse

	skills, err := model.GetOASFSkillsByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, fmt.Errorf("fail to get skills by agent uids: %v", err)
	}

	providers, err := model.GetA2AProvidersByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, fmt.Errorf("fail to get providers by agent uids: %v", err)
	}

	trustModels, err := model.GetTrustModelsByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, fmt.Errorf("fail to get trust models by agent uids: %v", err)
	}

	services, err := model.GetServicesByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, fmt.Errorf("fail to get services by agent uids: %v", err)
	}

	for _, agent := range agents {
		_services := services[agent.UID]
		var uri string
		var mcpEndpoint string
		var oasfEndpoint string
		var a2aEndpoint string
		var version string

		for _, _service := range _services {
			if _service.ServiceName == "a2a" {
				a2aEndpoint = _service.Endpoint
				version = _service.Version
			}
			if _service.ServiceName == "mcp" {
				mcpEndpoint = _service.Endpoint
			}
			if _service.ServiceName == "oasf" {
				oasfEndpoint = _service.Endpoint
			}
			if _service.ServiceName == "web" {
				uri = _service.Endpoint
			}
		}

		var skillTagsResponse = make([]serverTypes.SkillTagResponse, 0)
		for i, skill := range skills[agent.UID] {
			skillTagsResponse = append(skillTagsResponse, serverTypes.SkillTagResponse{
				ID:          strconv.Itoa(i),
				Name:        skill.SkillName,
				Description: "OASF Skill",
				Tags:        []string{skill.SkillName},
			})
		}

		var trustModelsResponse = make([]string, 0)
		for _, trustModel := range trustModels[agent.UID] {
			trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
		}

		var providerResponse serverTypes.ProviderResponse

		if providers[agent.UID] != nil {
			providerResponse = serverTypes.ProviderResponse{
				Organization: providers[agent.UID].Organization,
				URL:          providers[agent.UID].URL,
			}
		} else {
			providerResponse = serverTypes.ProviderResponse{}
		}

		chainInfo, _ := config.GetChainInfo(agent.ChainID)

		deployerInfo := config.GetContractsDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

		var status string
		if agent.Active {
			status = "active"
		} else {
			status = "inactive"
		}

		resp = append(resp, &serverTypes.AgentResponse{
			UID:              agent.UID,
			AgentID:          agent.AgentID,
			A2AEndpoint:      a2aEndpoint,
			WalletAddress:    agent.AgentWallet,
			Owner:            agent.Owner,
			ChainID:          agent.ChainID,
			ChainName:        chainInfo.ChainName,
			ChainLogo:        chainInfo.ChainLogo,
			Name:             agent.Name,
			Description:      agent.Description,
			URL:              uri,
			Provider:         providerResponse,
			IconURL:          agent.Image,
			Skills:           skillTagsResponse,
			TrustModels:      trustModelsResponse,
			IdentityRegistry: agent.IdentityRegistry,
			Deployer:         deployerInfo.Deployer,
			DeployerLogo:     deployerInfo.LogoURL,
			MCPEndpoint:      mcpEndpoint,
			OASFEndpoint:     oasfEndpoint,
			Version:          version,

			WalletAddressScanURL: fmt.Sprintf("%s/address/%s", chainInfo.ScanPrefix, agent.AgentWallet),
			ReputationRegistry:   deployerInfo.ReputationAddress,
			Status:               status,
		})
	}

	return resp, nil
}

// 检测value是不是hex编码，如果是则解码，否则返回原值
func decodeMetadataValue(value string) ([]byte, bool) {
	data, err := hex.DecodeString(strings.TrimPrefix(value, "0x"))
	if err != nil {
		return nil, false
	}
	return data, true
}
