package api

import (
	"agent_identity/config"
	"agent_identity/helper"
	"agent_identity/model"
	"agent_identity/types"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func GetCardResponse(agentUID uint64) (*AgentResponse, error) {
	agent, err := model.GetAgentByUID(agentUID)
	if err != nil {
		return nil, err
	}

	if agent == nil || agent.AgentID == "" {
		return nil, errors.New("agent not found")
	}

	skills, err := model.GetSkillsByAgentUID(agentUID)
	if err != nil {
		return nil, err
	}

	skillTags, err := model.GetSkillTagsByAgentUID(agentUID)
	if err != nil {
		return nil, err
	}

	provider, err := model.GetProviderByAgentUID(agentUID)
	if err != nil {
		return nil, err
	}

	trustModels, err := model.GetTrustModelsByAgentUID(agentUID)
	if err != nil {
		return nil, err
	}

	metadataRaw, err := model.GetMetadata(agent.ChainID, agent.IdentityRegistry, agent.AgentID)
	if err != nil {
		return nil, err
	}

	tokenURL, err := model.GetTokenURL(agent.ChainID, agent.IdentityRegistry, agent.AgentID)
	if err != nil {
		return nil, err
	}

	var metadataResponse = make([]MetadataResponse, 0)
	for _, metadata := range metadataRaw {

		data, err := hex.DecodeString(metadata.Value)
		if err != nil {
			return nil, err
		}

		var value string
		if len(data) == (20) {
			value = common.BytesToAddress(data).String()
		} else if len(data) == (32) {
			value = common.BytesToHash(data).String()
		} else {
			value = string(data)
		}

		metadataResponse = append(metadataResponse, MetadataResponse{Key: metadata.Key, Value: value})
	}

	var skillTagsResponse = make([]SkillTagResponse, 0)
	for _, skill := range skills {
		var tags = make([]string, 0)
		for _, skillTag := range skillTags[skill.ID] {
			tags = append(tags, skillTag.Tag)
		}
		skillTagsResponse = append(skillTagsResponse, SkillTagResponse{
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

	var providerResponse ProviderResponse
	if provider != nil {
		providerResponse = ProviderResponse{
			Organization: provider.Organization,
			URL:          provider.URL,
		}
	} else {
		providerResponse = ProviderResponse{}
	}

	chainInfo := config.GetChainInfo(agent.ChainID)

	deployerInfo := config.GetContractsDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

	resp := AgentResponse{
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
		return nil, err
	}
	oasfEndpoint, err := model.GetOASFEndpointByAgentUID(agent.UID)
	if err != nil {
		return nil, err
	}
	if mcpEndpoint != nil {
		resp.MACEndpoint = mcpEndpoint.Endpoint
	}
	if oasfEndpoint != nil {
		resp.OASFEndpoint = oasfEndpoint.Endpoint
	}

	return &resp, nil
}

func GetAgentList(page, pageSize int) ([]*AgentResponse, int64, error) {
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

func GetAgentListByFilter(page, pageSize int, trustModel []string, chains []string) ([]*AgentResponse, int64, error) {
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

func SearchAgentListBySkill(skill string, page, pageSize int) ([]*AgentResponse, int, error) {
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

func SearchAgentListByName(name string, page, pageSize int) ([]*AgentResponse, int, error) {
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

func FilterSearchAgentListByFilter(name string, page, pageSize int, trustModelIDs, chainIDs []string) ([]*AgentResponse, int64, error) {
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

func FilterSearchAgentListBySemantic(desc string, limit int, threshold float64, trustModelIDs, chainIDs []string) ([]*AgentResponse, error) {
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

func formatAgentResponse(agents []*model.Agent) ([]*AgentResponse, error) {
	if len(agents) == 0 {
		return []*AgentResponse{}, nil
	}

	agentUIDs := make([]uint64, 0, len(agents))
	for _, agent := range agents {
		agentUIDs = append(agentUIDs, agent.UID)
	}

	var resp []*AgentResponse

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
		var skillTagsResponse = make([]SkillTagResponse, 0)
		for _, skill := range skills[agent.UID] {
			var tags = make([]string, 0)
			for _, skillTag := range skillTags[agent.UID][skill.ID] {
				tags = append(tags, skillTag.Tag)
			}
			skillTagsResponse = append(skillTagsResponse, SkillTagResponse{
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

		var providerResponse ProviderResponse

		if providers[agent.UID] != nil {
			providerResponse = ProviderResponse{
				Organization: providers[agent.UID].Organization,
				URL:          providers[agent.UID].URL,
			}
		} else {
			providerResponse = ProviderResponse{}
		}

		chainInfo := config.GetChainInfo(agent.ChainID)

		deployerInfo := config.GetContractsDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

		resp = append(resp, &AgentResponse{
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

func SetFeedback(request UploadFeedbackRequest) (string, string, error) {
	agent, err := model.GetAgentByUID(request.UID)
	if err != nil {
		return "", "", fmt.Errorf("fail to get agent: %w", err)
	}

	agentID, err := strconv.ParseUint(agent.AgentID, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("fail to parse agent id: %w", err)
	}

	agentRegistry := common.HexToAddress(agent.IdentityRegistry).String()
	clientAddress := common.HexToAddress(request.ClientAddress).String()

	feedback := &types.Feedback{
		AgentRegistry: fmt.Sprintf("eip155:%s:%s", agent.ChainID, agentRegistry),
		AgentId:       int64(agentID),
		ClientAddress: fmt.Sprintf("eip155:%s:%s", agent.ChainID, clientAddress),
		CreatedAt:     strconv.FormatInt(time.Now().Unix(), 10),
		Score:         request.Score,
		Tag1:          request.Tag1,
		Tag2:          request.Tag2,
		Context:       request.Context,
		Task:          request.Task,
		Capability:    request.Capability,
		Endpoint:      request.Endpoint,
		Domain:        request.Domain,
		Name:          request.Name,
	}

	if request.ProofOfPayment != nil {
		feedback.ProofOfPayment = &types.ProofOfPayment{
			FromAddress: common.HexToAddress(request.ProofOfPayment.FromAddress).String(),
			ToAddress:   common.HexToAddress(request.ProofOfPayment.ToAddress).String(),
			ChainId:     agent.ChainID,
			TxHash:      request.ProofOfPayment.TxHash,
		}
	}

	feedbackData, err := json.Marshal(feedback)
	if err != nil {
		return "", "", fmt.Errorf("fail to marshal feedback: %w", err)
	}

	feedbackURI, err := helper.GetHelper().UploadFeedbackToS3(agent.ChainID, agentRegistry, agent.AgentID, clientAddress, request.IndexLimit, feedbackData)
	if err != nil {
		return "", "", fmt.Errorf("fail to upload feedback to s3: %w", err)
	}

	feedbackHash := common.BytesToHash(sha256.New().Sum(feedbackData)).String()

	return feedbackURI, feedbackHash, nil
}

func getNetworkList() ([]types.NetworkResponse, error) {
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
