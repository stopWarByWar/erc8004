package agentcard

import (
	"testing"
)

func TestGetAgentProfile(t *testing.T) {
	agent, err := GetAgentProfile("ipfs://QmNfTaioEdMfdnSRTHXJy6juAw4xFNGeCnAnfkqckjZHth")
	if err != nil {
		panic(err)
	}
	t.Log(agent.Type)
	t.Log(agent.Name)
	t.Log(agent.Description)
	t.Log(agent.Image)
	t.Log(agent.Services)
	t.Log(agent.X402Support)
	t.Log(agent.Active)
	t.Log(agent.Registrations)
	t.Log(agent.SupportedTrust)
}
