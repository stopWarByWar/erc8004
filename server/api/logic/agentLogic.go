package logic

import (
	"agent_identity/config"
	"agent_identity/helper"
	"agent_identity/model"
	"agent_identity/server/api/types"
	serverTypes "agent_identity/server/api/types"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func GetAgentListByFilter(page, pageSize int, name *string, trustModel, chains, skills *[]string, x402Support, active, haveFeedback *bool) ([]*serverTypes.AgentResponse, int64, error) {
	agents, total, err := model.GetAgentsByFilter(name, page, pageSize, trustModel, chains, skills, x402Support, active, haveFeedback)
	if err != nil {
		return nil, 0, err
	}
	resp, err := formatAgentResponse(agents)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}

func GetCardResponse(agentUID uint64) (*serverTypes.AgentResponse, error) {
	agent, err := model.GetAgentByUID(agentUID)
	if err != nil {
		return nil, fmt.Errorf("fail to get agent by uid: %v", err)
	}

	if agent == nil || agent.AgentID == "" {
		return nil, fmt.Errorf("agent not found")
	}

	skills, err := model.GetOASFSkillsByAgentUID(agentUID)
	if err != nil {
		return nil, fmt.Errorf("fail to get skills by agent uid: %v", err)
	}

	provider, err := model.GetA2AProviderByAgentUID(agentUID)
	if err != nil {
		return nil, fmt.Errorf("fail to get provider by agent uid: %v", err)
	}

	trustModels, err := model.GetTrustModelsByAgentUID(agentUID)
	if err != nil {
		return nil, fmt.Errorf("fail to get trust models by agent uid: %v", err)
	}

	metadataRaw, err := model.GetMetadata(agent.ChainID, agent.IdentityRegistry, agent.AgentID)
	if err != nil {
		return nil, fmt.Errorf("fail to get metadata by agent uid: %v", err)
	}

	var metadataResponse = make([]serverTypes.MetadataResponse, 0)
	for _, metadata := range metadataRaw {
		var value string
		data, ok := decodeMetadataValue(metadata.Value)
		if !ok {
			value = metadata.Value
		} else {
			if len(data) == (20) {
				value = common.BytesToAddress(data).String()
			} else if len(data) == (32) {
				value = common.BytesToHash(data).String()
			} else {
				value = string(data)
			}
		}

		metadataResponse = append(metadataResponse, serverTypes.MetadataResponse{Key: metadata.Key, Value: value})
	}

	var skillTagsResponse = make([]serverTypes.SkillTagResponse, 0)
	for _, skill := range skills {
		skillTagsResponse = append(skillTagsResponse, serverTypes.SkillTagResponse{
			ID:          skill.SkillName,
			Name:        skill.SkillName,
			Description: "OASF Skill",
			Tags:        []string{skill.SkillName},
		})
	}

	var trustModelsResponse = make([]string, 0)
	for _, trustModel := range trustModels {
		trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
	}

	var providerResponse serverTypes.ProviderResponse
	if provider != nil {
		providerResponse = serverTypes.ProviderResponse{
			Organization: provider.Organization,
			URL:          provider.URL,
		}
	} else {
		providerResponse = serverTypes.ProviderResponse{}
	}

	services, err := model.GetServicesByAgentUID(agent.UID)
	if err != nil {
		return nil, fmt.Errorf("fail to get services by agent uid: %v", err)
	}
	var uri string
	var mcpEndpoint string
	var oasfEndpoint string
	var a2aEndpoint string
	var version string

	var endpoints = make([]serverTypes.EndpointResponse, 0)

	for i, _service := range services {
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
		if len(strings.TrimSpace(_service.ServiceName)) > 0 {
			endpoints = append(endpoints, serverTypes.EndpointResponse{
				No:       i + 1,
				Name:     strings.ToLower(_service.ServiceName),
				Endpoint: _service.Endpoint,
			})
		}
	}

	chainInfo, _ := config.GetChainInfo(agent.ChainID)

	deployerInfo := config.GetContractsDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

	var status string
	if agent.Active {
		status = "active"
	} else {
		status = "inactive"
	}

	resp := serverTypes.AgentResponse{
		UID:                agent.UID,
		AgentID:            agent.AgentID,
		A2AEndpoint:        a2aEndpoint,
		WalletAddress:      agent.AgentWallet,
		Owner:              agent.Owner,
		ChainID:            agent.ChainID,
		ChainName:          chainInfo.ChainName,
		ChainLogo:          chainInfo.ChainLogo,
		Name:               agent.Name,
		Description:        agent.Description,
		URL:                uri,
		Provider:           providerResponse,
		IconURL:            agent.Image,
		Version:            version,
		DocumentationURL:   agent.A2ADocumentationURL,
		Skills:             skillTagsResponse,
		TrustModels:        trustModelsResponse,
		IdentityRegistry:   agent.IdentityRegistry,
		Metadata:           metadataResponse,
		TokenURL:           agent.AgentURI,
		Deployer:           deployerInfo.Deployer,
		DeployerLogo:       deployerInfo.LogoURL,
		MCPEndpoint:        mcpEndpoint,
		OASFEndpoint:       oasfEndpoint,
		ReputationRegistry: deployerInfo.ReputationAddress,
		Status:             status,
		X402Support:        agent.X402Support,
		Endpoints:          endpoints,
	}
	return &resp, nil
}

func FilterSearchAgentListBySemantic(desc string, limit int, threshold float64, trustModelIDs *[]string, chainIDs *[]string, skills *[]string, x402Support *bool, active *bool, haveFeedback *bool) ([]*serverTypes.AgentResponse, error) {
	filters := &model.VectorSearchFilters{
		TrustModel:   trustModelIDs,
		ChainID:      chainIDs,
		Skills:       skills,
		X402Support:  x402Support,
		Active:       active,
		HaveFeedback: haveFeedback,
	}
	agentUIDs, err := model.SearchSimilarVectors(desc, limit, threshold, filters)
	if err != nil {
		return nil, err
	}
	agents, err := model.GetAgentsByUIDs(agentUIDs)
	if err != nil {
		return nil, err
	}

	formattedAgents, err := formatAgentResponse(agents)
	if err != nil {
		return nil, err
	}
	return formattedAgents, nil
}

func UploadAgentProfile(request serverTypes.UploadAgentProfileRequest, logoData []byte) (tokenURI string, err error) {
	logoURI, err := helper.GetHelper().UploadLogoToS3(request.ChainID, common.HexToAddress(request.IdentityRegistry).String(), request.AgentID, logoData)
	if err != nil {
		return "", fmt.Errorf("fail to upload logo to s3: %w", err)
	}

	agentID, err := strconv.ParseUint(request.AgentID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("fail to parse agent id: %w", err)
	}

	agentProfile := &types.AgentProfile{
		Type:        "https://eips.ethereum.org/EIPS/eip-8004#registration-v1",
		Name:        request.Name,
		Description: request.Description,
		Image:       logoURI,
		Endpoints:   request.Endpoints,
		Registrations: []types.Registration{
			{
				AgentId:       int64(agentID),
				AgentRegistry: fmt.Sprintf("eip155:%s:%s", request.ChainID, common.HexToAddress(request.IdentityRegistry).String()),
			},
		},
		SupportedTrust: request.SupportedTrust,
	}

	agentProfileData, err := json.Marshal(agentProfile)
	if err != nil {
		return "", fmt.Errorf("fail to marshal agent profile: %w", err)
	}

	tokenURI, err = helper.GetHelper().UploadAgentProfileToS3(request.ChainID, common.HexToAddress(request.IdentityRegistry).String(), request.AgentID, agentProfileData)
	if err != nil {
		return "", fmt.Errorf("fail to upload agent profile to s3: %w", err)
	}
	return tokenURI, nil
}
