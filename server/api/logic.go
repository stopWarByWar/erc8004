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
		metadataResponse = append(metadataResponse, MetadataResponse{
			Key:   metadata.Key,
			Value: metadata.Value,
		})
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

	deployerInfo := config.GetDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

	return &AgentResponse{
		UID:              agent.UID,
		AgentID:          agent.AgentID,
		AgentDomain:      agent.A2AEndpoint,
		AgentAddress:     agent.AgentWallet,
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
	}, nil
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

func formatAgentResponse(agents []*model.Agent) ([]*AgentResponse, error) {
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

		deployerInfo := config.GetDeployerInfo(agent.ChainID, common.HexToAddress(agent.IdentityRegistry).String())

		resp = append(resp, &AgentResponse{
			UID:              agent.UID,
			AgentID:          agent.AgentID,
			AgentDomain:      agent.A2AEndpoint,
			AgentAddress:     agent.AgentWallet,
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
		})
	}

	return resp, nil
}

func SetFeedback(request UploadFeedbackRequest) (string, string, error) {
	agent, err := model.GetAgentByUID(request.UID)
	if err != nil {
		return "", "", fmt.Errorf("fail to get agent: %w", err)
	}

	if agent.AgentID != request.FeedbackAuth.AgentId {
		return "", "", fmt.Errorf("agent id mismatch")
	}

	if common.HexToAddress(agent.IdentityRegistry).String() != common.HexToAddress(request.FeedbackAuth.IdentityRegistry).String() {
		return "", "", fmt.Errorf("identity registry mismatch")
	}

	feedbackAuthData, err := json.Marshal(request.FeedbackAuth)
	if err != nil {
		return "", "", fmt.Errorf("fail to marshal feedback auth: %w", err)
	}

	agentID, err := strconv.ParseUint(request.FeedbackAuth.AgentId, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("fail to parse agent id: %w", err)
	}

	feedback := &types.Feedback{
		AgentRegistry: common.HexToAddress(request.FeedbackAuth.IdentityRegistry).String(),
		AgentId:       int64(agentID),
		ClientAddress: common.HexToAddress(request.FeedbackAuth.ClientAddress).String(),
		CreatedAt:     strconv.FormatInt(time.Now().Unix(), 10),
		FeedbackAuth:  hex.EncodeToString(feedbackAuthData),
		Score:         request.Score,
		Tag1:          &request.Tag1,
		Tag2:          &request.Tag2,
		Context:       &request.Context,
		Task:          &request.Task,
		Capability:    &request.Capability,
	}

	feedbackData, err := json.Marshal(feedback)
	if err != nil {
		return "", "", fmt.Errorf("fail to marshal feedback: %w", err)
	}

	feedbackURI, err := helper.GetHelper().UploadFeedbackToS3(request.FeedbackAuth.ChainId, common.HexToAddress(request.FeedbackAuth.IdentityRegistry).String(), request.FeedbackAuth.AgentId, common.HexToAddress(request.FeedbackAuth.ClientAddress).String(), request.FeedbackAuth.IndexLimit, feedbackData)
	if err != nil {
		return "", "", fmt.Errorf("fail to upload feedback to s3: %w", err)
	}

	feedbackHash := common.BytesToHash(sha256.New().Sum(feedbackData)).String()

	return feedbackURI, feedbackHash, nil
}
