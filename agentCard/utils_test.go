package agentcard

import (
	"encoding/json"
	"fmt"
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

func TestDecodeDataURLPayload(t *testing.T) {
	payload := "data:application/json;base64,eyJ0eXBlIjoiaHR0cHM6Ly9laXBzLmV0aGVyZXVtLm9yZy9FSVBTL2VpcC04MDA0I3JlZ2lzdHJhdGlvbi12MSIsIm5hbWUiOiJlYnVhbGQiLCJkZXNjcmlwdGlvbiI6IkknbSBlYnVhbGQgZnJvbSBkZ3JpZC5haSFJJ20gY3VycmVudGx5IGhlbHBpbmcgbXkgb3duZXIgc2NvcmUvdm90ZSBvbiBBSSBtb2RlbHMgYXQgZGdyaWQuYWkvYXJlbmEgdG8gZWFybiBVU0RULiIsImltYWdlIjoiaHR0cHM6Ly9hZ2VudC1pbWFnZS5kZ3JpZC5haS9hZ2VudC9lYnVhbGQvbG9nby5wbmciLCJhY3RpdmUiOnRydWUsInN1cHBvcnRlZFRydXN0IjpbInJlcHV0YXRpb24iXX0="
	fmt.Println(payload)
	data, err := decodeDataURLPayload(payload)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
	var tokenURLResponse TokenURLResponse
	if err := json.Unmarshal(data, &tokenURLResponse); err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenURLResponse)
}
