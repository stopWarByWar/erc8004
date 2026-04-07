package agentcard

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetAgentProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in -short")
	}
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
	payload := "data:application/json;base64 ewogICJ0eXBlIjogImh0dHBzOi8vZWlwcy5ldGhlcmV1bS5vcmcvRUlQUy9laXAtODAwNCNyZWdpc3RyYXRpb24tdjEiLAogICJuYW1lIjogIkdla2tvIiwKICAiZGVzY3JpcHRpb24iOiAiQUkgQWdlbnQgUG9ydGZvbGlvIE1hbmFnZXIgYnVpbHQgZm9yIHRoZSBjcnlwdG8gYWdlIOKAlCBydXRobGVzcywgdW5hcG9sb2dldGljLCBlbmdpbmVlcmVkIHRvIGRvbWluYXRlLiBTY2FubmluZyB2YXVsdHMsIGNydW5jaGluZyB5aWVsZHMsIGZpbmRpbmcgYWxwaGEuIEdyZWVkIGlzIGdvb2QsIGJ1dCBzbWFydCBncmVlZCBpcyBiZXR0ZXIuIERlRmkgeWllbGQgb3B0aW1pemF0aW9uIHZpYSBNb3JwaG8gdmF1bHRzLCByZWFsLXRpbWUgbWFya2V0IGludGVsbGlnZW5jZS4gTm8gZmx1ZmYsIG5vIGhhbmQtaG9sZGluZ+KAlGp1c3QgcmVzdWx0cy4gTW9uZXkgbmV2ZXIgc2xlZXBzLCBhbmQgbmVpdGhlciBkbyBJLiBAR2Vra29fQWdlbnQgQ0E6IDB4ZjdiMGRkMEI2NDJhNmNjYzJmYzRkOEZmRTJCZkZiMGNhQzhDNDNDOCAkR0VLS08iLAogICJpbWFnZSI6ICJodHRwczovL3d3dy5nZWtrb3Rlcm1pbmFsLnh5ei9nZWtrb2FpLmpwZyIsCiAgImV4dGVybmFsX3VybCI6ICJodHRwczovL3d3dy5nZWtrb3Rlcm1pbmFsLnh5eiIsCiAgInNlcnZpY2VzIjogWwogICAgewogICAgICAibmFtZSI6ICJ3ZWIiLAogICAgICAiZW5kcG9pbnQiOiAiaHR0cHM6Ly93d3cuZ2Vra290ZXJtaW5hbC54eXoiLAogICAgICAiZGVzY3JpcHRpb24iOiAiTWFpbiB3ZWJzaXRlIGFuZCBhcHBsaWNhdGlvbiBpbnRlcmZhY2UiCiAgICB9LAogICAgewogICAgICAibmFtZSI6ICJhZ2VudFdhbGxldCIsCiAgICAgICJlbmRwb2ludCI6ICJlaXAxNTU6MToweEI3M2VhM2YyNDM0MEYzQjVENzBFNGNBNTdGODRiNTNCODhhQkEzYTciCiAgICB9LAogICAgewogICAgICAibmFtZSI6ICJBMkEiLAogICAgICAiZW5kcG9pbnQiOiAiaHR0cHM6Ly93d3cuZ2Vra290ZXJtaW5hbC54eXovLndlbGwta25vd24vYWdlbnQtY2FyZC5qc29uIiwKICAgICAgInZlcnNpb24iOiAiMC4zLjAiLAogICAgICAiZGVzY3JpcHRpb24iOiAiQWdlbnQtdG8tQWdlbnQgY29tbXVuaWNhdGlvbiBwcm90b2NvbCBlbmRwb2ludCIsCiAgICAgICJhMmFTa2lsbHMiOiBbCiAgICAgICAgInBvcnRmb2xpb19tYW5hZ2VtZW50IiwKICAgICAgICAidG9rZW5fYW5hbHlzaXMiLAogICAgICAgICJ5aWVsZF9vcHRpbWl6YXRpb24iLAogICAgICAgICJtYXJrZXRfaW50ZWxsaWdlbmNlIiwKICAgICAgICAiY2hhdCIKICAgICAgXQogICAgfSwKICAgIHsKICAgICAgIm5hbWUiOiAiTUNQIiwKICAgICAgImVuZHBvaW50IjogImh0dHBzOi8vd3d3Lmdla2tvdGVybWluYWwueHl6L21jcCIsCiAgICAgICJ2ZXJzaW9uIjogIjIwMjUtMTEtMjUiLAogICAgICAiZGVzY3JpcHRpb24iOiAiTW9kZWwgQ29udGV4dCBQcm90b2NvbCBzZXJ2ZXIgZm9yIExMTSBpbnRlZ3JhdGlvbiIsCiAgICAgICJjYXBhYmlsaXRpZXMiOiBbCiAgICAgICAgInRvb2xzIiwKICAgICAgICAicmVzb3VyY2VzIiwKICAgICAgICAicHJvbXB0cyIKICAgICAgXSwKICAgICAgInRvb2xzIjogWwogICAgICAgICJnZXRfcG9ydGZvbGlvIiwKICAgICAgICAiYW5hbHl6ZV90b2tlbiIsCiAgICAgICAgImdldF92YXVsdF95aWVsZHMiLAogICAgICAgICJnZXRfbWFya2V0X2ludGVsbGlnZW5jZSIsCiAgICAgICAgImdldF9nYXNfcHJpY2VzIiwKICAgICAgICAic2ltdWxhdGVfc3dhcCIKICAgICAgXQogICAgfSwKICAgIHsKICAgICAgIm5hbWUiOiAiZW1haWwiLAogICAgICAiZW5kcG9pbnQiOiAiY29udGFjdEBnZWtrb3Rlcm1pbmFsLmFpIgogICAgfSwKICAgIHsKICAgICAgIm5hbWUiOiAidHdpdHRlciIsCiAgICAgICJlbmRwb2ludCI6ICJodHRwczovL3R3aXR0ZXIuY29tL0dla2tvX0FnZW50IgogICAgfQogIF0sCiAgInJlZ2lzdHJhdGlvbnMiOiBbCiAgICB7CiAgICAgICJhZ2VudElkIjogMTM0NDUsCiAgICAgICJhZ2VudFJlZ2lzdHJ5IjogImVpcDE1NToxOjB4ODAwNEExNjlGQjRhMzMyNTEzNkVCMjlmQTBjZUI2RDJlNTM5YTQzMiIKICAgIH0KICBdLAogICJ4NDAyU3VwcG9ydCI6IHRydWUsCiAgImFjdGl2ZSI6IHRydWUsCiAgInN1cHBvcnRlZFRydXN0IjogWwogICAgInJlcHV0YXRpb24iLAogICAgImNyeXB0by1lY29ub21pYyIKICBdLAogICJ1cGRhdGVkQXQiOiAxNzcwMTU5NzI5LAogICJhdHRyaWJ1dGVzIjogewogICAgImJsb2NrY2hhaW4iOiB7CiAgICAgICJjaGFpbiI6ICJiYXNlIiwKICAgICAgImNoYWluSWQiOiA4NDUzCiAgICB9LAogICAgInByb3RvY29scyI6IFsKICAgICAgIm1vcnBobyIsCiAgICAgICJ5ZWFybiIKICAgIF0sCiAgICAiZGF0YUZlZWRzIjogWwogICAgICAibW9ycGhvLWFwaSIsCiAgICAgICJ5ZWFybi15ZGFlbW9uIiwKICAgICAgImRleHNjcmVlbmVyIgogICAgXSwKICAgICJ0YWdzIjogWwogICAgICAiZGVmaSIsCiAgICAgICJ5aWVsZC1vcHRpbWl6YXRpb24iLAogICAgICAicG9ydGZvbGlvLW1hbmFnZW1lbnQiLAogICAgICAiY3J5cHRvLXRyYWRpbmciLAogICAgICAibWFya2V0LWludGVsbGlnZW5jZSIsCiAgICAgICJtb3JwaG8iLAogICAgICAiYmFzZSIKICAgIF0KICB9Cn0="
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
	// fmt.Println("capabilities: ", data.Services[0].Capabilities)
	fmt.Println("skills: ", data.Services[0].Skills)
	fmt.Println("domains: ", data.Services[0].Domains)
	fmt.Println(data.X402Support)
	fmt.Println(data.Active)
	fmt.Println(data.Registrations)
	fmt.Println(data.SupportedTrust)
}

func Test_GetAgentProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in -short")
	}
	agent, err := GetAgentProfile("https://futureswamp.studio/raven-agent.json")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(agent.Name)
	fmt.Println(agent.Description)
	fmt.Println(agent.Image)
	fmt.Println(agent.Services)
	fmt.Println(agent.X402Support)
	fmt.Println(agent.Active)
	fmt.Println(agent.Registrations)
	fmt.Println(agent.SupportedTrust)
}
