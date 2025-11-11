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
			"agent_card_list": []*AgentResponse{mockAgentCard()},
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
	resp, total, err := GetAgentList(pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list", "Internal Error", c)
		return
	}

	SuccessResp(gin.H{
		"agent_card_list": resp,
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

	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}
	agentCard, err := GetCardResponse(agentUID)
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
			"agent_card_list": []*AgentResponse{mockAgentCard()},
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

	agents, total, err := GetAgentListByTrustModel(pageInt, pageSizeInt, trustModel)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by trust model", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
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
			"agent_list": []*AgentResponse{mockAgentCard()},
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

	agents, total, err := SearchAgentListBySkill(skill, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by skill", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

func GetAgentCardsSearchByNameHandler(c *gin.Context) {

	name := c.Query("name")
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
	agents, total, err := SearchAgentListByName(name, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by name", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

// func GetFeedbackListHandler(c *gin.Context) {
// 	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
// 	if err != nil {
// 		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
// 		return
// 	}
// 	page := c.Query("page")
// 	pageSize := c.Query("page_size")
// 	pageInt, err := strconv.Atoi(page)
// 	if err != nil {
// 		ErrResp(nil, "fail to get page", "Invalid Request", c)
// 		return
// 	}
// 	if pageInt <= 0 {
// 		pageInt = 1
// 	}
// 	pageSizeInt, err := strconv.Atoi(pageSize)
// 	if err != nil {
// 		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
// 		return
// 	}
// 	if pageSizeInt <= 0 {
// 		pageSizeInt = 10
// 	}

// 	feedbacks, total, err := model.GetFeedbackList(agentUID, pageInt, pageSizeInt)
// 	if err != nil {
// 		ErrResp(nil, "fail to get feedback list", "Internal Error", c)
// 		return
// 	}
// 	SuccessResp(gin.H{
// 		"feedbacks": feedbacks,
// 		"total":     total,
// 	}, c)
// }

func GetAgentCommentsHandler(c *gin.Context) {
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}

	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
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
	comments, total, err := model.GetCommentsByAgentUID(agentUID, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent comments", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"comments": comments,
		"total":    total,
	}, c)
}

func GetAgentFeedbacksHandler(c *gin.Context) {
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
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

	feedbacks, total, err := model.GetFeedbacksByAgentUID(agentUID, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent feedbacks", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"feedbacks": feedbacks,
		"total":     total,
	}, c)

}

func mockAgentCard() *AgentResponse {
	return &AgentResponse{
		UID:          1,
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

func mockFeedback() []*model.FeedbackResp {
	return []*model.FeedbackResp{
		{
			UID:                1,
			AgentUID:           1,
			ChainID:            "97",
			AgentID:            "1",
			ReputationRegistry: "0x0000000000000000000000000000000000000000",
			ClientAddress:      "0x0000000000000000000000000000000000000000",
			FeedbackIndex:      1,
			Score:              10,
			Tag1:               "1",
			Tag2:               "1",
			FeedbackURI:        "1",
			FeedbackHash:       "1",
			TxHash:             "1",
			Timestamps:         1,
			Name:               "1",
			Avatar:             "https://bnbattest.s3.ap-southeast-1.amazonaws.com/logo/bas.png",
			Passport:           true,
		},
		{
			UID:                2,
			AgentUID:           2,
			ChainID:            "97",
			AgentID:            "2",
			ReputationRegistry: "0x0000000000000000000000000000000000000000",
			ClientAddress:      "0x0000000000000000000000000000000000000000",
			FeedbackIndex:      2,
			Score:              10,
			Tag1:               "2",
			Tag2:               "2",
			FeedbackURI:        "2",
			FeedbackHash:       "2",
			TxHash:             "2",
			Timestamps:         2,
			Name:               "2",
			Avatar:             "https://bnbattest.s3.ap-southeast-1.amazonaws.com/logo/bas.png",
			Passport:           true,
		},
	}
}

func mockComment() []*model.Comment {
	return []*model.Comment{
		{
			CommentAttestationID: "1",
			Commenter:            "0x0000000000000000000000000000000000000000",
			AgentUID:             1,
			CommentText:          "1",
			Score:                10,
			Timestamps:           1,
			Name:                 "1",
			Avatar:               "https://bnbattest.s3.ap-southeast-1.amazonaws.com/logo/bas.png",
			Passport:             true,
		},
		{
			CommentAttestationID: "2",
			Commenter:            "0x0000000000000000000000000000000000000000",
			AgentUID:             2,
			CommentText:          "2",
			Score:                10,
			Timestamps:           2,
			Name:                 "2",
			Avatar:               "https://bnbattest.s3.ap-southeast-1.amazonaws.com/logo/bas.png",
			Passport:             true,
		},
	}
}
