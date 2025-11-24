package agentcard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"trpc.group/trpc-go/trpc-a2a-go/server"
)

const defaultIPFSGateway = "https://ipfs.io/ipfs/"

func GetAgentCardFromTokenURL(owner, tokenId, tokenURL, chainID, identityRegistryAddr string, timestamps uint64) (*Agent, bool, []error) {
	var errors []error
	var inserted bool = true
	body, err := fetchTokenURLBody(tokenURL)
	if err != nil {
		errors = append(errors, err)
		return nil, false, errors
	}

	var tokenURLResponse TokenURLResponse
	err = json.Unmarshal(body, &tokenURLResponse)
	if err != nil {
		errors = append(errors, err)
		return nil, false, errors
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
				errors = append(errors, err)
				inserted = false
			}
			agent.AgentCard = agentCard
			agent.Endpoint = endpoint.Endpoint
		}

		if endpoint.Name == "agentWallet" {
			namespace, _chainID, agentWallet, err := formatAddress(endpoint.Endpoint)
			if err != nil {
				errors = append(errors, err)
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
				agent.AgentID = tokenId
				namespace, _chainID, registryAddr, err := formatAddress(registration.AgentRegistry)
				if err != nil {
					errors = append(errors, err)
					inserted = false
				}
				if namespace == "eip155" && _chainID == chainID && registryAddr == identityRegistryAddr {
					agent.IdentityRegistry = registryAddr
				}
			}
		}
	}
	if len(agent.AgentID) == 0 || len(agent.IdentityRegistry) == 0 || len(agent.Namespace) == 0 || len(agent.AgentWallet) == 0 {
		errors = append(errors, fmt.Errorf("invalid agent agent id:%s, identity registry:%s, namespace:%s, agent wallet:%s", tokenId, identityRegistryAddr, agent.Namespace, agent.AgentWallet))
		inserted = false
	}

	return agent, inserted, errors
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

func fetchTokenURLBody(tokenURL string) ([]byte, error) {
	resolvedURL, err := resolveTokenURL(tokenURL)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(resolvedURL)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to get token URL response: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed, status: %s, status code: %d", response.Status, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to read token URL response body: %v", err)
	}
	if len(body) == 0 {
		return nil, errors.New("empty token URL response body")
	}
	return body, nil
}

func resolveTokenURL(tokenURL string) (string, error) {
	if strings.HasPrefix(tokenURL, "ipfs://") {
		cidPath := strings.TrimPrefix(tokenURL, "ipfs://")
		if len(cidPath) == 0 {
			return "", fmt.Errorf("invalid ipfs token URL: %s", tokenURL)
		}
		gateway := os.Getenv("IPFS_GATEWAY_URL")
		if len(gateway) == 0 {
			gateway = defaultIPFSGateway
		} else if !strings.HasSuffix(gateway, "/") {
			gateway += "/"
		}
		return fmt.Sprintf("%s%s", gateway, cidPath), nil
	}
	if strings.HasPrefix(tokenURL, "https://") {
		return tokenURL, nil
	}
	return "", fmt.Errorf("unsupported token URL scheme: %s", tokenURL)
}
