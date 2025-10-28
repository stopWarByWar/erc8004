package api

import (
	"agent_identity/model"
	"math"
)

func GetCardResponse(agentID string) (*CardResponse, error) {
	agentCard, err := model.GetAgentCard(agentID)
	if err != nil {
		return nil, err
	}

	if agentCard == nil {
		return nil, nil
	}

	skills, err := model.GetSkillsByAgentID(agentID)
	if err != nil {
		return nil, err
	}

	skillTags, err := model.GetSkillTagsByAgentID(agentID)
	if err != nil {
		return nil, err
	}

	provider, err := model.GetProviderByAgentID(agentID)
	if err != nil {
		return nil, err
	}

	trustModels, err := model.GetTrustModelsByAgentID(agentID)
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
	if agentCard.CommentCount > 0 {
		score = math.Round(float64(agentCard.Score)/float64(agentCard.CommentCount)*10) / 10
	}

	return &CardResponse{
		AgentID:      agentCard.AgentID,
		AgentDomain:  agentCard.AgentDomain,
		AgentAddress: agentCard.AgentAddress,
		ChainID:      agentCard.ChainID,
		Namespace:    agentCard.Namespace,
		Name:         agentCard.Name,
		Description:  agentCard.Description,
		URL:          agentCard.URL,
		Provider: ProviderResponse{
			Organization: provider.Organization,
			URL:          provider.URL,
		},
		IconURL:          agentCard.IconURL,
		Version:          agentCard.Version,
		DocumentationURL: agentCard.DocumentationURL,
		Skills:           skillTagsResponse,
		TrustModels:      trustModelsResponse,
		Score:            score,
	}, nil
}

func GetCardList(page, pageSize int) ([]*CardResponse, int64, error) {
	agentCards, total, err := model.GetAgentList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatCardResponse(agentCards)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func GetCardResponseByTrustModel(page, pageSize int, trustModel []string) ([]*CardResponse, int64, error) {
	agentCards, total, err := model.GetAgentCardsByTrustModel(page, pageSize, trustModel)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatCardResponse(agentCards)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func SearchCardResponseBySkill(skill string, page, pageSize int) ([]*CardResponse, int, error) {
	agentCards, total, err := model.SearchSkillsAgentCards(skill, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatCardResponse(agentCards)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func SearchCardResponseByName(name string, page, pageSize int) ([]*CardResponse, int, error) {
	agentCards, total, err := model.SearchAgentCardByName(name, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	cards, err := formatCardResponse(agentCards)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func formatCardResponse(agentCards []*model.AgentCard) ([]*CardResponse, error) {
	agentIDs := make([]string, 0, len(agentCards))
	for _, agentCard := range agentCards {
		agentIDs = append(agentIDs, agentCard.AgentID)
	}

	var cards []*CardResponse

	skills, err := model.GetSkillsByAgentIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	skillTags, err := model.GetSkillTagsByAgentIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	provider, err := model.GetProvidersByAgentIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	trustModels, err := model.GetTrustModelsByAgentIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	for _, agentCard := range agentCards {
		var skillTagsResponse []SkillTagResponse
		for _, skill := range skills[agentCard.AgentID] {
			var tags []string
			for _, skillTag := range skillTags[agentCard.AgentID][skill.ID] {
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
		for _, trustModel := range trustModels[agentCard.AgentID] {
			trustModelsResponse = append(trustModelsResponse, trustModel.TrustModel)
		}

		score := 0.0
		if agentCard.CommentCount > 0 {
			score = math.Round(float64(agentCard.Score)/float64(agentCard.CommentCount)*10) / 10
		}

		cards = append(cards, &CardResponse{
			AgentID:      agentCard.AgentID,
			AgentDomain:  agentCard.AgentDomain,
			AgentAddress: agentCard.AgentAddress,
			ChainID:      agentCard.ChainID,
			Namespace:    agentCard.Namespace,
			Name:         agentCard.Name,
			Description:  agentCard.Description,
			URL:          agentCard.URL,
			Provider: ProviderResponse{
				Organization: provider[agentCard.AgentID].Organization,
				URL:          provider[agentCard.AgentID].URL,
			},
			IconURL:          agentCard.IconURL,
			Version:          agentCard.Version,
			DocumentationURL: agentCard.DocumentationURL,
			Skills:           skillTagsResponse,
			TrustModels:      trustModelsResponse,
			UserInterface:    agentCard.UserInterface,
			Score:            score,
		})
	}

	return cards, nil
}
