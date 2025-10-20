package api

import (
	agentcard "agent_identity/agentCard"
	"agent_identity/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	agentCardList, total, err := GetCardList(pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list", "Internal Error", c)
		return
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
		return
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
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	_logger.WithFields(logrus.Fields{
		"page":        pageInt,
		"page_size":   pageSizeInt,
		"trust_model": trustModel,
	}).Info("GetCardResponseByTrustModel")

	agentCardList, total, err := GetCardResponseByTrustModel(pageInt, pageSizeInt, trustModel)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by trust model", "Internal Error", c)
		return
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
	page := c.Query("page")
	pageSize := c.Query("page_size")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	agentCardList, total, err := SearchCardResponseBySkill(skill, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by skill", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_card_list": agentCardList,
		"total":           total,
		"current_page":    pageInt,
	}, c)
}

func CheckAuthFeedbackExistsHandler(c *gin.Context) {
	clientAddress := c.Query("address")
	agentServerID := c.Query("agent_id")
	authed, agentClientID, err := model.CheckAuthFeedbackExists(clientAddress, agentServerID)
	if err != nil {
		ErrResp(nil, "fail to check auth feedback exists", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"authed":          authed,
		"agent_client_id": agentClientID,
	}, c)
}

func GetCommentListHandler(c *gin.Context) {
	agentID := c.Query("agent_id")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	isAuthorized := c.Query("type") == "authorized"
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	if pageInt <= 0 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	if pageSizeInt <= 0 {
		pageSizeInt = 10
	}

	comments, total, err := model.GetCommentList(agentID, pageInt, pageSizeInt, isAuthorized)
	if err != nil {
		ErrResp(nil, "fail to get comment list", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"comments": comments,
		"total":    total,
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
		TrustModels:   []string{agentcard.TrustModelFeedback, agentcard.TrustModelInferenceValidation, agentcard.TrustModelTeeAttestation},
		UserInterface: "https://passport.bnbattest.io",
		Score:         10.0,
	}
}
