package agentcard

import (
	"fmt"
	"os"
	"testing"
)

func Test_unmarshalAgentCard(t *testing.T) {
	data, err := os.ReadFile("agentCard_example.json")
	if err != nil {
		t.Fatalf("failed to read agent card example: %v", err)
	}
	agentCard, err := unmarshalAgentCard(data)
	if err != nil {
		t.Fatalf("failed to unmarshal agent card: %v", err)
	}
	fmt.Printf("agent card: %v", agentCard)
}

func Test_ValidateAgentCard(t *testing.T) {
	data, err := os.ReadFile("agentCard_example.json")
	if err != nil {
		t.Fatalf("failed to read agent card example: %v", err)
	}
	agentCard, err := unmarshalAgentCard(data)
	if err != nil {
		t.Fatalf("failed to unmarshal agent card: %v", err)
	}

	agentCardList := ValidateAgentCard("agent_id", "0x1234567890123456789012345678901234567890", agentCard)
	for _, card := range agentCardList {
		fmt.Printf("agent card: %v", card)
	}
}
