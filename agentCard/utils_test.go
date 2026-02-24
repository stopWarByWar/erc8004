package agentcard

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetAgentProfile(t *testing.T) {
	agent, err := GetAgentProfile("ipfs://QmcLpNkoqchCXjJS5diuP61hAwnUDaJnrAGTDn4PzkHKzn")
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
	payload := "data:application/json;base64,eyJ0eXBlIjoiaHR0cHM6Ly9laXBzLmV0aGVyZXVtLm9yZy9FSVBTL2VpcC04MDA0I3JlZ2lzdHJhdGlvbi12MSIsIm5hbWUiOiJDcnlwdG9NYXN0ZXIgQWdlbnQgIzE2NyIsImRlc2NyaXB0aW9uIjoiQXV0b21hdGVkIHRyYWRpbmcgYWdlbnQgcG93ZXJlZCBieSBDcnlwdG9NYXN0ZXIiLCJzZXJ2aWNlcyI6WyJ0cmFkaW5nIiwiZGVmaSJdLCJhY3RpdmUiOnRydWUsInN1cHBvcnRlZFRydXN0IjpbInJlcHV0YXRpb24iXX0="
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

func TestDecodeAgentProfileData(t *testing.T) {
	tokenURLResponse := `{"type":"https://eips.ethereum.org/EIPS/eip-8004#registration-v1","name":"Vector Zero","description":"Strategic Technogist | Revenue & Runway Optimization","image":"https://i.ibb.co/CKN66Mg3/DB8-DC218-52-A2-4824-B793-C9942336641-A5.png","services":[]}`
	encoded, err := decodeAgentProfileData([]byte(tokenURLResponse))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(encoded.Type)
	fmt.Println(encoded.Name)
	fmt.Println(encoded.Description)
	fmt.Println(encoded.Image)
	fmt.Println(encoded.Services)
	fmt.Println(encoded.X402Support)
	fmt.Println(encoded.Active)
	fmt.Println(encoded.Registrations)
	fmt.Println(encoded.SupportedTrust)
}

func Test_getAgentProfileFromEncodedData(t *testing.T) {
	payload := "data:application/json;base64,eyJ0eXBlIjoiaHR0cHM6Ly9laXBzLmV0aGVyZXVtLm9yZy9FSVBTL2VpcC04MDA0I3JlZ2lzdHJhdGlvbi12MSIsIm5hbWUiOiJDcnlwdG9NYXN0ZXIgQWdlbnQgIzQ2IiwiZGVzY3JpcHRpb24iOiJBdXRvbWF0ZWQgdHJhZGluZyBhZ2VudCBwb3dlcmVkIGJ5IENyeXB0b01hc3RlciIsInNlcnZpY2VzIjpbInRyYWRpbmciLCJkZWZpIl0sImFjdGl2ZSI6dHJ1ZSwic3VwcG9ydGVkVHJ1c3QiOlsicmVwdXRhdGlvbiJdfQ=="
	fmt.Println(payload)
	data, err := getAgentProfileFromEncodedData(payload)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data.Type)
	fmt.Println(data.Name)
	fmt.Println(data.Description)
	fmt.Println(data.Image)
	fmt.Println("name: ", data.Services[0].Name)
	fmt.Println("endpoint: ", data.Services[0].Endpoint)
	fmt.Println("version: ", data.Services[0].Version)
	fmt.Println("capabilities: ", data.Services[0].Capabilities)
	fmt.Println("skills: ", data.Services[0].Skills)
	fmt.Println("domains: ", data.Services[0].Domains)
	fmt.Println(data.X402Support)
	fmt.Println(data.Active)
	fmt.Println(data.Registrations)
	fmt.Println(data.SupportedTrust)
}
