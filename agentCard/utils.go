package agentcard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"trpc.group/trpc-go/trpc-a2a-go/server"
)

func GetAgentCardFromTokenURL(owner, tokenId, tokenURL, chainID, identityRegistryAddr string, timestamps uint64) (*Agent, error) {
	response, err := http.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to get token URL response: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed, status: %s, status code: %d", response.Status, response.StatusCode)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to read token URL response body: %v", err)
	}
	var tokenURLResponse TokenURLResponse
	err = json.Unmarshal(body, &tokenURLResponse)
	if err != nil {
		return nil, fmt.Errorf("data error: failed to unmarshal token URL response body: %v", err)
	}

	var agent = &Agent{
		Type:             tokenURLResponse.Type,
		Name:             tokenURLResponse.Name,
		Description:      tokenURLResponse.Description,
		Image:            tokenURLResponse.Image,
		SupportedTrust:   tokenURLResponse.SupportedTrust,
		AgentID:          tokenId,
		TokenURL:         tokenURL,
		ChainID:          chainID,
		Owner:            owner,
		Timestamps:       timestamps,
		UserInterfaceURL: tokenURLResponse.UserInterfaceURL,
	}

	var agentCard *server.AgentCard

	for _, endpoint := range tokenURLResponse.Endpoints {
		if endpoint.Name == "A2A" {
			agentCard, err = getAgentCardFromA2AEndpoint(endpoint.Endpoint)
			if err != nil {
				return nil, err
			}
			agent.AgentCard = agentCard
			agent.Endpoint = endpoint.Endpoint
		}

		if endpoint.Name == "agentWallet" {
			namespace, _chainID, agentWallet, err := formatAddress(endpoint.Endpoint)
			if err != nil {
				return nil, fmt.Errorf("data error: failed to format address: %v", err)
			}
			if namespace == "eip155" && _chainID == chainID {
				agent.Namespace = namespace
				agent.AgentWallet = agentWallet
			}
		}
	}

	if tokenURLResponse.Registrations != nil {
		for _, registration := range tokenURLResponse.Registrations {
			if strconv.FormatUint(registration.AgentID, 10) == tokenId {
				namespace, _chainID, registryAddr, err := formatAddress(registration.AgentRegistry)
				if err != nil {
					return nil, fmt.Errorf("data error: failed to format address: %v", err)
				}
				if namespace == "eip155" && _chainID == chainID && registryAddr == identityRegistryAddr {
					agent.IdentityRegistry = registryAddr
					agent.Namespace = namespace
				}
			}
		}
	}

	return agent, nil
}

func getAgentCardFromA2AEndpoint(endpoint string) (*server.AgentCard, error) {
	if !validateA2AEndpoint(endpoint) {
		return nil, fmt.Errorf("invalid A2A endpoint: %s", endpoint)
	}

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to get A2A endpoint response: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed, status: %s, status code: %d", response.Status, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to read A2A endpoint response body: %v", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body")
	}
	agentCard, err := unmarshalAgentCard(body)
	if err != nil {
		return nil, fmt.Errorf("data error: failed to unmarshal A2A endpoint response body: %v", err)
	}
	return agentCard, nil
}

func unmarshalAgentCard(body []byte) (*server.AgentCard, error) {
	var agentCard *server.AgentCard
	err := json.Unmarshal(body, &agentCard)
	if err != nil {
		return nil, fmt.Errorf("data error: failed to unmarshal A2A endpoint response body: %v", err)
	}
	return agentCard, nil
}

func formatAddress(address string) (string, string, string, error) {
	addressSlice := strings.Split(address, ":")
	if len(addressSlice) != 3 {
		return "", "", "", errors.New("invalid address format")
	}
	return addressSlice[0], addressSlice[1], addressSlice[2], nil
}
func validateA2AEndpoint(endpoint string) bool {
	//检查https://agent.example/.well-known/agent-card.json这个格式
	if !strings.HasSuffix(endpoint, "/.well-known/agent-card.json") || !strings.HasPrefix(endpoint, "https://") {
		return false
	}
	return true
}
