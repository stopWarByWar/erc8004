package server

import "agent_identity/model"

func GetCardResponse(agentID string) (*CardResponse, error) {
	agentCard, err := model.GetAgentCard(agentID)
	if err != nil {
		return nil, err
	}

	skills, err := model.GetSkillsByAgentID(agentID)
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

	return &CardResponse{
		AgentID:          agentCard.AgentID,
		AgentDomain:      agentCard.AgentDomain,
		AgentAddress:     agentCard.AgentAddress,
		Name:             agentCard.Name,
		Description:      agentCard.Description,
		URL:              agentCard.URL,
		Provider:         provider,
		IconURL:          agentCard.IconURL,
		Version:          agentCard.Version,
		DocumentationURL: agentCard.DocumentationURL,
		Skills:           skills,
		TrustModels:      trustModels,
	}, nil
}

func GetCardList(page, pageSize, limit int) ([]*CardResponse, error) {
	agentCards, err := model.GetAgentList(page, pageSize, limit)
	if err != nil {
		return nil, err
	}
	return formatCardResponse(agentCards)
}

func GetCardResponseByTrustModel(trustModel []string) ([]*CardResponse, error) {
	agentCards, err := model.GetAgentCardsByTrustModel(trustModel)
	if err != nil {
		return nil, err
	}
	return formatCardResponse(agentCards)
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

	provider, err := model.GetProvidersByAgentIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	trustModels, err := model.GetTrustModelsByAgentIDs(agentIDs)
	if err != nil {
		return nil, err
	}

	for _, agentCard := range agentCards {
		cards = append(cards, &CardResponse{
			AgentID:          agentCard.AgentID,
			AgentDomain:      agentCard.AgentDomain,
			AgentAddress:     agentCard.AgentAddress,
			Name:             agentCard.Name,
			Description:      agentCard.Description,
			URL:              agentCard.URL,
			Provider:         provider[agentCard.AgentID],
			IconURL:          agentCard.IconURL,
			Version:          agentCard.Version,
			DocumentationURL: agentCard.DocumentationURL,
			Skills:           skills[agentCard.AgentID],
			TrustModels:      trustModels[agentCard.AgentID],
		}, nil)
	}

	return cards, nil
}
