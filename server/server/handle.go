package server

import (
	agentcard "agent_identity/agentCard"
	"strconv"

	"github.com/gin-gonic/gin"
)

var mock = false

func GetAgentCardListHandler(c *gin.Context) {
	if mock {
		SuccessResp(gin.H{
			"agent_card_list": []*CardResponse{mockAgentCard()},
			"total":           1,
			"current_page":    1,
		}, c)
		return
	}

	page := c.Query("page")
	pageSize := c.Query("page_size")
	limit := c.Query("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ErrResp(nil, "fail to get limit", "Invalid Request", c)
	}
	agentCardList, total, err := GetCardList(pageInt, pageSizeInt, limitInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list", "Internal Error", c)
	}

	SuccessResp(gin.H{
		"agent_card_list": agentCardList,
		"total":           total,
		"current_page":    pageInt,
	}, c)
}

func GetAgentCardDetailHandler(c *gin.Context) {
	if mock {
		SuccessResp(gin.H{
			"agent_card": mockAgentCard(),
		}, c)
		return
	}

	agentID := c.Query("agent_id")
	agentCard, err := GetCardResponse(agentID)
	if err != nil {
		ErrResp(nil, "fail to get agent card detail", "Internal Error", c)
	}
	SuccessResp(gin.H{
		"agent_card": agentCard,
	}, c)
}

func GetAgentCardListByTrustModelHandler(c *gin.Context) {
	if mock {
		SuccessResp(gin.H{
			"agent_card_list": []*CardResponse{mockAgentCard()},
			"total":           1,
			"current_page":    1,
		}, c)
		return
	}

	trustModel := c.QueryArray("trust_model")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	limit := c.Query("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ErrResp(nil, "fail to get limit", "Invalid Request", c)
	}
	agentCardList, total, err := GetCardResponseByTrustModel(pageInt, pageSizeInt, limitInt, trustModel)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by trust model", "Internal Error", c)
	}
	SuccessResp(gin.H{
		"agent_card_list": agentCardList,
		"total":           total,
		"current_page":    pageInt,
	}, c)
}

func GetTrustModelListHandler(c *gin.Context) {
	SuccessResp(gin.H{
		"trust_model_list": []string{agentcard.TrustModelFeedback, agentcard.TrustModelInferenceValidation, agentcard.TrustModelTeeAttestation},
	}, c)
}

func GetAgentCardsSearchBySkillHandler(c *gin.Context) {
	if mock {
		SuccessResp(gin.H{
			"agent_card_list": []*CardResponse{mockAgentCard()},
		}, c)
		return
	}

	skill := c.Query("skill")
	agentCardList, err := SearchCardResponseBySkill(skill)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by skill", "Internal Error", c)
	}
	SuccessResp(gin.H{
		"agent_card_list": agentCardList,
	}, c)
}

func mockAgentCard() *CardResponse {
	return &CardResponse{
		AgentID:      "1",
		AgentDomain:  "passport.bnbattest.io",
		AgentAddress: "0x0000000000000000000000000000000000000000",
		ChainID:      "97",
		Namespace:    "eip:155",

		Name:        "Test-Agent",
		Description: "This is a test for the agent card. It is used to test the agent card API. This is a test for the agent card. This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.",
		URL:         "https://passport.bnbattest.io",
		Provider: ProviderResponse{
			Organization: "BAS",
			URL:          "https://passport.bnbattest.io",
		},
		IconURL:          "https://bnbattest.s3.ap-southeast-1.amazonaws.com/logo/bas.png",
		Version:          "1.0.0",
		DocumentationURL: "https://passport.bnbattest.io",
		Skills: []SkillTagResponse{
			{
				ID:          "route-optimizer-traffic",
				Name:        "Traffic-Aware Route Optimizer",
				Description: "Calculates the optimal driving route between two or more locations, taking into account real-time traffic conditions, road closures, and user preferences (e.g., avoid tolls, prefer highways).",
				Tags:        []string{"maps", "routing", "navigation", "directions", "traffic"},
			},
			{
				ID:          "custom-map-generator",
				Name:        "Personalized Map Generator",
				Description: "Creates custom map images or interactive map views based on user-defined points of interest, routes, and style preferences. Can overlay data layers.",
				Tags:        []string{"maps", "customization", "visualization", "cartography"},
			},
		},
		TrustModels: []string{agentcard.TrustModelFeedback, agentcard.TrustModelInferenceValidation, agentcard.TrustModelTeeAttestation},
	}
}
