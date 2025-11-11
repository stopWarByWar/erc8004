package api

import (
	"agent_identity/model"
	"errors"
	"math"
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

	var skillTagsResponse []SkillTagResponse
	for _, skill := range skills {
		var tags []string
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

	var trustModelsResponse []string
	for _, trustModel := range trustModels {
		trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
	}

	score := 0.0
	if agent.CommentCount > 0 {
		score = math.Round(float64(agent.Score)/float64(agent.CommentCount)*10) / 10
	}

	return &AgentResponse{
		UID:          agent.UID,
		AgentID:      agent.AgentID,
		AgentDomain:  agent.A2AEndpoint,
		AgentAddress: agent.AgentWallet,
		ChainID:      agent.ChainID,
		Namespace:    agent.Namespace,
		Name:         agent.Name,
		Description:  agent.Description,
		URL:          agent.URL,
		Provider: ProviderResponse{
			Organization: provider.Organization,
			URL:          provider.URL,
		},
		IconURL:          agent.Image,
		Version:          agent.Version,
		DocumentationURL: agent.DocumentationURL,
		Skills:           skillTagsResponse,
		TrustModels:      trustModelsResponse,
		Score:            score,
		UserInterface:    agent.UserInterfaceURL,
		IdentityRegistry: agent.IdentityRegistry,
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

func GetAgentListByTrustModel(page, pageSize int, trustModel []string) ([]*AgentResponse, int64, error) {
	agents, total, err := model.GetAgentsByTrustModel(page, pageSize, trustModel)
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
		var skillTagsResponse []SkillTagResponse
		for _, skill := range skills[agent.UID] {
			var tags []string
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

		var trustModelsResponse []string
		for _, trustModel := range trustModels[agent.UID] {
			trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
		}

		score := 0.0
		if agent.CommentCount > 0 {
			score = math.Round(float64(agent.Score)/float64(agent.CommentCount)*10) / 10
		}

		resp = append(resp, &AgentResponse{
			UID:          agent.UID,
			AgentID:      agent.AgentID,
			AgentDomain:  agent.A2AEndpoint,
			AgentAddress: agent.AgentWallet,
			ChainID:      agent.ChainID,
			Namespace:    agent.Namespace,
			Name:         agent.Name,
			Description:  agent.Description,
			URL:          agent.URL,
			Provider: ProviderResponse{
				Organization: providers[agent.UID].Organization,
				URL:          providers[agent.UID].URL,
			},
			IconURL:          agent.Image,
			Version:          agent.Version,
			DocumentationURL: agent.DocumentationURL,
			Skills:           skillTagsResponse,
			TrustModels:      trustModelsResponse,
			Score:            score,
			UserInterface:    agent.UserInterfaceURL,
			IdentityRegistry: agent.IdentityRegistry,
		})
	}

	return resp, nil
}
