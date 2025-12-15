package agentcard

import (
	"testing"
)

func TestGetAgentCardFromTokenURL(t *testing.T) {
	agent, provider, err := GetAgentCardFromTokenURL("0xA7132182Cbc0ceA8bE148FDE88faaD3BB9410d48", "1", "ipfs://QmNfTaioEdMfdnSRTHXJy6juAw4xFNGeCnAnfkqckjZHth", "11155111", "0x8004a6090Cd10A7288092483047B097295Fb8847", 1760570592)
	if err != nil {
		panic(err)
	}
	t.Log(agent)
	t.Log(provider)
}
