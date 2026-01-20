package logic

import (
	"agent_identity/config"
	"agent_identity/helper"
	"agent_identity/model"
	"agent_identity/server/api/types"
	serverTypes "agent_identity/server/api/types"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

func FilterSearchAgentListByFilter(name string, page, pageSize int, trustModelIDs, chainIDs []string) ([]*serverTypes.AgentResponse, int64, error) {
	agents, total, err := model.FilterSearchAgentsByName(name, page, pageSize, trustModelIDs, chainIDs)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatAgentResponse(agents)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func GetAgentListByFilter(page, pageSize int, trustModel []string, chains []string) ([]*serverTypes.AgentResponse, int64, error) {
	agents, total, err := model.GetAgentsByFilter(page, pageSize, trustModel, chains)
	if err != nil {
		return nil, 0, err
	}
	resp, err := formatAgentResponse(agents)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}

func GetAgentList(page, pageSize int) ([]*serverTypes.AgentResponse, int64, error) {
	agents, total, err := model.GetAgentList(page, pageSize)
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

	skills, err := model.GetSkillsByAgentUID(agentUID)
	if err != nil {
		return nil, fmt.Errorf("fail to get skills by agent uid: %v", err)
	}

	skillTags, err := model.GetSkillTagsByAgentUID(agentUID)
	if err != nil {
		return nil, fmt.Errorf("fail to get skill tags by agent uid: %v", err)
	}

	provider, err := model.GetProviderByAgentUID(agentUID)
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

	tokenURL, err := model.GetTokenURL(agent.ChainID, agent.IdentityRegistry, agent.AgentID)
	if err != nil {
		return nil, fmt.Errorf("fail to get token url by agent uid: %v", err)
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
		var tags = make([]string, 0)
		for _, skillTag := range skillTags[skill.ID] {
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
	for _, trustModel := range trustModels {
		trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
	}

	score := 0.0
	if agent.CommentCount > 0 {
		score = math.Round(float64(agent.Score)/float64(agent.CommentCount)*10) / 10
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

	chainInfo, _ := config.GetChainInfo(agent.ChainID)

	deployerInfo := config.GetContractsDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

	resp := serverTypes.AgentResponse{
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
		Metadata:         metadataResponse,
		TokenURL:         tokenURL,
		Deployer:         deployerInfo.Deployer,
		DeployerLogo:     deployerInfo.LogoURL,

		WalletAddressScanURL:        fmt.Sprintf("%s/address/%s", chainInfo.ScanPrefix, agent.AgentWallet),
		WalletAddressExpirationTime: agent.AgentWalletExpirationTime,
		ReputationRegistry:          deployerInfo.ReputationAddress,
	}

	mcpEndpoint, err := model.GetMCPEndpointByAgentUID(agent.UID)
	if err != nil {
		return nil, fmt.Errorf("fail to get mcp endpoint by agent uid: %v", err)
	}
	oasfEndpoint, err := model.GetOASFEndpointByAgentUID(agent.UID)
	if err != nil {
		return nil, fmt.Errorf("fail to get oasf endpoint by agent uid: %v", err)
	}
	if mcpEndpoint != nil {
		resp.MACEndpoint = mcpEndpoint.Endpoint
	}
	if oasfEndpoint != nil {
		resp.OASFEndpoint = oasfEndpoint.Endpoint
	}

	return &resp, nil
}

func SearchAgentListBySkill(skill string, page, pageSize int) ([]*serverTypes.AgentResponse, int, error) {
	agents, total, err := model.SearchAgentsBySkill(skill, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatAgentResponse(agents)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func SearchAgentListByName(name string, page, pageSize int) ([]*serverTypes.AgentResponse, int, error) {
	agents, total, err := model.SearchAgentsByName(name, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatAgentResponse(agents)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func FilterSearchAgentListBySemantic(desc string, limit int, threshold float64, trustModelIDs, chainIDs []string) ([]*serverTypes.AgentResponse, error) {
	filters := &model.VectorSearchFilters{
		TrustModel: trustModelIDs,
		ChainID:    chainIDs,
	}
	agentUIDs, err := model.SearchSimilarVectors(desc, limit, threshold, filters)
	if err != nil {
		return nil, err
	}
	agents, err := model.GetAgentsByUIDs(agentUIDs)
	if err != nil {
		return nil, err
	}

	cards, err := formatAgentResponse(agents)
	if err != nil {
		return nil, err
	}
	return cards, nil
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
