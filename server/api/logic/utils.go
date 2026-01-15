package logic

import (
	"agent_identity/config"
	"agent_identity/model"
	"fmt"
	"math"

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

	skills, err := model.GetSkillsByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, err
	}

	skillTags, err := model.GetSkillTagsByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, err
	}

	providers, err := model.GetProvidersByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, err
	}

	trustModels, err := model.GetTrustModelsByAgentUIDs(agentUIDs)
	if err != nil {
		return nil, err
	}

	for _, agent := range agents {
		var skillTagsResponse = make([]serverTypes.SkillTagResponse, 0)
		for _, skill := range skills[agent.UID] {
			var tags = make([]string, 0)
			for _, skillTag := range skillTags[agent.UID][skill.ID] {
				tags = append(tags, skillTag.Tag)
			}
			skillTagsResponse = append(skillTagsResponse, serverTypes.SkillTagResponse{
				ID:          skill.ID,
				Name:        skill.Name,
				Description: skill.Description,
				Tags:        tags,
			})
		}

		var trustModelsResponse = make([]string, 0)
		for _, trustModel := range trustModels[agent.UID] {
			trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
		}

		score := 0.0
		if agent.CommentCount > 0 {
			score = math.Round(float64(agent.Score)/float64(agent.CommentCount)*10) / 10
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

		resp = append(resp, &serverTypes.AgentResponse{
			UID:              agent.UID,
			AgentID:          agent.AgentID,
			A2AEndpoint:      agent.A2AEndpoint,
			WalletAddress:    agent.AgentWallet,
			Owner:            agent.Owner,
			ChainID:          agent.ChainID,
			ChainName:        chainInfo.ChainName,
			ChainLogo:        chainInfo.ChainLogo,
			Namespace:        agent.Namespace,
			Name:             agent.Name,
			Description:      agent.Description,
			URL:              agent.URL,
			Provider:         providerResponse,
			IconURL:          agent.Image,
			Version:          agent.Version,
			DocumentationURL: agent.DocumentationURL,
			Skills:           skillTagsResponse,
			TrustModels:      trustModelsResponse,
			Score:            score,
			UserInterface:    agent.UserInterfaceURL,
			IdentityRegistry: agent.IdentityRegistry,
			Deployer:         deployerInfo.Deployer,
			DeployerLogo:     deployerInfo.LogoURL,

			WalletAddressScanURL:        fmt.Sprintf("%s/address/%s", chainInfo.ScanPrefix, agent.AgentWallet),
			WalletAddressExpirationTime: agent.AgentWalletExpirationTime,
			ReputationRegistry:          deployerInfo.ReputationAddress,
		})
	}

	return resp, nil
}
