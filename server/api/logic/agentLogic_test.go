package logic

import (
	"fmt"
	"testing"

	agentcard "agent_identity/agentCard"
	"agent_identity/server/api/types"

	"agent_identity/config"
)

func TestUpdateFilterInfo(t *testing.T) {
	initTest()
	config.UpdateGeneralInfo()
	fmt.Printf("generalInfo: %+v\n", config.GetGeneralInfo())

}

var name = "Agent"
var trustModelIDs = []string{}
var chainIDs = []string{"11155111"}
var skills = []string{}
var x402Support = false
var active = false
var haveFeedback = false

func TestFilterSearchAgentListByFilter(t *testing.T) {
	initTest()
	page := 1
	pageSize := 10
	agents, total, err := GetAgentListByFilter(page, pageSize, &name, &trustModelIDs, &chainIDs, &skills, &x402Support, &active, &haveFeedback)
	if err != nil {
		t.Errorf("GetAgentListByFilter error: %v", err)
	}
	fmt.Println(agents, total)
}

func TestGetAgentListByFilter(t *testing.T) {
	initTest()
	page := 1
	pageSize := 10

	agents, total, err := GetAgentListByFilter(page, pageSize, &name, &trustModelIDs, &chainIDs, &skills, &x402Support, &active, &haveFeedback)
	if err != nil {
		t.Errorf("GetAgentListByFilter error: %v", err)
	}
	if len(agents) > 1 {
		fmt.Printf("agents: %+v, total: %d\n", agents[0], total)
	} else {
		t.Errorf("agents: %+v, total: %d\n", agents, total)
	}
}
func TestGetCardResponse(t *testing.T) {
	initTest()
	agentUID := uint64(1)
	agentCard, err := GetCardResponse(agentUID)
	if err != nil {
		t.Errorf("GetCardResponse error: %v", err)
	}
	fmt.Printf("agentCard: %+v\n", agentCard)
}

func TestFilterSearchAgentListBySemantic(t *testing.T) {
	initTest()
	desc := "test"
	limit := 10
	threshold := 0.5
	agents, err := FilterSearchAgentListBySemantic(desc, limit, threshold, &trustModelIDs, &chainIDs, &skills, &x402Support, &active, &haveFeedback)
	if err != nil {
		t.Errorf("FilterSearchAgentListBySemantic error: %v", err)
	}
	fmt.Printf("agents: %+v\n", agents)
}

func TestUploadAgentProfile(t *testing.T) {
	initTest()
	versionString := "1.0.0"
	capabilities := map[string]interface{}{
		"tools":     []string{"test"},
		"prompts":   []string{"test"},
		"resources": []string{"test"},
	}
	request := types.UploadAgentProfileRequest{
		AgentID:          "1",
		ChainID:          "11155111",
		Name:             "test",
		Description:      "test",
		IdentityRegistry: "0x8004AA63c570c570eBF15376c0dB199918BFe9Fb",
		SupportedTrust:   []string{agentcard.TrustModelTeeAttestation},
		Endpoints: []types.Endpoint{
			{
				Name:     "A2A",
				Endpoint: "https://api.test.com/a2a",
			},
			{
				Name:     "OASF",
				Endpoint: "https://api.test.com/oasf",
				Skills:   []string{"test"},
			},
			{
				Name:         "MCP",
				Endpoint:     "https://api.test.com/mcp",
				Version:      &versionString,
				Capabilities: &capabilities,
			},
		},
	}
	logoData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	tokenURI, err := UploadAgentProfile(request, logoData)
	if err != nil {
		t.Errorf("UploadAgentProfile error: %v", err)
		return
	}
	fmt.Printf("tokenURI: %s\n", tokenURI)
}
