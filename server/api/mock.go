package api

import (
	agentcard "agent_identity/agentCard"
	"agent_identity/model"
)

func mockAgent() *AgentResponse {
	return &AgentResponse{
		UID:              1,
		AgentID:          "1",
		AgentDomain:      "passport.bnbattest.io",
		AgentAddress:     "0x0000000000000000000000000000000000000000",
		ChainID:          "97",
		Namespace:        "eip:155",
		IdentityRegistry: "0xa98a5542a1aab336397d487e32021e0e48bef717",
		Name:             "Test-Agent",
		Description:      "This is a test for the agent card. It is used to test the agent card API. This is a test for the agent card. This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.This is a test for the agent card.",
		URL:              "https://passport.bnbattest.io",
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

func mockFeedbacks() []*model.FeedbackResp {
	return []*model.FeedbackResp{
		{
			UID:                1,
			AgentUID:           9001,
			ChainID:            "97",
			AgentID:            "224",
			ReputationRegistry: "0x5f4b38ae73d1f17c39f93a07c11d36f5a8a9c011",
			ClientAddress:      "0xc2b4d2a115d5c9c82ba4f54aafaef47de7b8a2a5",
			FeedbackIndex:      1,
			Score:              9,
			Tag1:               "响应速度",
			Tag2:               "服务态度",
			FeedbackURI:        "https://feedback.bitagent.test/agents/224/entries/1001",
			FeedbackHash:       "0x41f5dd7f78b498c4d84f9ef3c3e94d89f3c0a907bfae27c34dbe61308bb94f01",
			TxHash:             "0x8c40ce54f17f7b8b9f71d0b7dd93e9f66d37cbf3f6a4a5dd0f934e2b8700e501",
			Timestamps:         1728105600,
			Name:               "Alice Chen",
			Avatar:             "https://cdn.bitagent.test/avatars/alice.png",
			Passport:           true,
		},
		{
			UID:                1,
			AgentUID:           9001,
			ChainID:            "97",
			AgentID:            "224",
			ReputationRegistry: "0x5f4b38ae73d1f17c39f93a07c11d36f5a8a9c011",
			ClientAddress:      "0x5a07c82bc1dd3642c1f13e2c21b40c60d2837a60",
			FeedbackIndex:      2,
			Score:              8,
			Tag1:               "功能完整",
			Tag2:               "易用性",
			FeedbackURI:        "https://feedback.bitagent.test/agents/224/entries/1002",
			FeedbackHash:       "0x59e3cfb6c1a329be461b7d98a4f4cb6ccf5dce1db02dbe1cb05eae8cf7d246d8",
			TxHash:             "0x6be49fdf792d6f0c0b095499f917f6579c42b9d0e751fa035d13dbfba3a5d7c4",
			Timestamps:         1728019200,
			Name:               "Marco Li",
			Avatar:             "https://cdn.bitagent.test/avatars/marco.png",
			Passport:           false,
		},
		{
			UID:                1,
			AgentUID:           9001,
			ChainID:            "97",
			AgentID:            "224",
			ReputationRegistry: "0x5f4b38ae73d1f17c39f93a07c11d36f5a8a9c011",
			ClientAddress:      "0x5a07c82bc1dd3642c1f13e2c21b40c60d2837a60",
			FeedbackIndex:      2,
			Score:              8,
			Tag1:               "功能完整",
			Tag2:               "易用性",
			FeedbackURI:        "https://feedback.bitagent.test/agents/224/entries/1002",
			FeedbackHash:       "0x59e3cfb6c1a329be461b7d98a4f4cb6ccf5dce1db02dbe1cb05eae8cf7d246d8",
			TxHash:             "0x6be49fdf792d6f0c0b095499f917f6579c42b9d0e751fa035d13dbfba3a5d7c4",
			Timestamps:         1728019200,
			Name:               "Marco Li",
			Avatar:             "https://cdn.bitagent.test/avatars/marco.png",
			Passport:           false,
		},
		{
			UID:                1,
			AgentUID:           9001,
			ChainID:            "97",
			AgentID:            "224",
			ReputationRegistry: "0x5f4b38ae73d1f17c39f93a07c11d36f5a8a9c011",
			ClientAddress:      "0x5a07c82bc1dd3642c1f13e2c21b40c60d2837a60",
			FeedbackIndex:      2,
			Score:              8,
			Tag1:               "功能完整",
			Tag2:               "易用性",
			FeedbackURI:        "https://feedback.bitagent.test/agents/224/entries/1002",
			FeedbackHash:       "0x59e3cfb6c1a329be461b7d98a4f4cb6ccf5dce1db02dbe1cb05eae8cf7d246d8",
			TxHash:             "0x6be49fdf792d6f0c0b095499f917f6579c42b9d0e751fa035d13dbfba3a5d7c4",
			Timestamps:         1728019200,
			Name:               "Marco Li",
			Avatar:             "https://cdn.bitagent.test/avatars/marco.png",
			Passport:           false,
		},
		{
			UID:                1,
			AgentUID:           9001,
			ChainID:            "97",
			AgentID:            "224",
			ReputationRegistry: "0x5f4b38ae73d1f17c39f93a07c11d36f5a8a9c011",
			ClientAddress:      "0x5a07c82bc1dd3642c1f13e2c21b40c60d2837a60",
			FeedbackIndex:      2,
			Score:              8,
			Tag1:               "功能完整",
			Tag2:               "易用性",
			FeedbackURI:        "https://feedback.bitagent.test/agents/224/entries/1002",
			FeedbackHash:       "0x59e3cfb6c1a329be461b7d98a4f4cb6ccf5dce1db02dbe1cb05eae8cf7d246d8",
			TxHash:             "0x6be49fdf792d6f0c0b095499f917f6579c42b9d0e751fa035d13dbfba3a5d7c4",
			Timestamps:         1728019200,
			Name:               "Marco Li",
			Avatar:             "https://cdn.bitagent.test/avatars/marco.png",
			Passport:           false,
		},
	}
}

func mockComments() []*model.Comment {
	return []*model.Comment{
		{
			CommentAttestationID: "0x9a4ddf4c0d4d4c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f60718293a4b5c6d7e8f9",
			Commenter:            "0xc2b4d2a115d5c9c82ba4f54aafaef47de7b8a2a5",
			AgentUID:             9001,
			CommentText:          "客服响应很快，问题当日就解决了。",
			Score:                5,
			Timestamps:           1728109200,
			Name:                 "Alice Chen",
			Avatar:               "https://cdn.bitagent.test/avatars/alice.png",
			Passport:             true,
		},
		{
			CommentAttestationID: "0x7be42a9f52f4bc0ad6c1cddf38f0a7f2bb1932c4d0f1a2b3c4d5e6f708192a3b",
			Commenter:            "0x5a07c82bc1dd3642c1f13e2c21b40c60d2837a60",
			AgentUID:             9001,
			CommentText:          "功能覆盖了我们的大部分需求，希望文档再细致一些。",
			Score:                4,
			Timestamps:           1728022800,
			Name:                 "Marco Li",
			Avatar:               "https://cdn.bitagent.test/avatars/marco.png",
			Passport:             false,
		},
	}
}

func mockAgents(pageSizeInt int) []*AgentResponse {
	agents := make([]*AgentResponse, 0)
	for i := 0; i < pageSizeInt; i++ {
		agent := mockAgent()
		agent.UID = uint64(i + 1)
		agents = append(agents, agent)
	}
	return agents
}
