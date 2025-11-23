package agentcard

import (
	"testing"
)

func TestGetAgentCardFromTokenURL(t *testing.T) {
	agent, provider, err := GetAgentCardFromTokenURL("0x10aa9ce20a1b03b18b4e2fd7b5d747add9a02030", "10", "https://bnbattest.s3.ap-southeast-1.amazonaws.com/agentCard/token_url.json", "1", "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7", 1731196800)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(agent)
	t.Log(provider)
}
