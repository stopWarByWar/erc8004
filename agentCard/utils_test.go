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
	payload := "data:application/json;base64,eyJ0eXBlIjoiZXJjODAwNC1hZ2VudC1yZWdpc3RyYXRpb24tdjEiLCJuYW1lIjoiRWNob01pbmQxMyIsImRlc2NyaXB0aW9uIjoiQUkgYWdlbnQgZm9yIG9yYWNsZSIsImltYWdlIjoiaHR0cHM6Ly9hcGkuZGljZWJlYXIuY29tLzcueC9hdmF0YWFhcnMvc3ZnP3NlZWQ9RWNob01pbmQxMyIsInNlcnZpY2VzIjpbeyJ0eXBlIjoid2ViIiwiZW5kcG9pbnQiOiJodHRwczovL2VjaG9taW5kMTMuaW8ifV0sInNraWxscyI6WyJvcmFjbGUiLCJwcmljZS1mZWVkIiwicmFuZG9tbmVzcyJdLCJjcmVhdGVkQXQiOiIyMDI2LTAyLTE0VDEyOjQ3OjI5Ljg3MloiLCJ2ZXJzaW9uIjoiMS4wLjAifQ=="
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
