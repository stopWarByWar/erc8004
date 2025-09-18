package agentcard

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"trpc.group/trpc-go/trpc-a2a-go/server"
)

func GetAgentCard(address, agentID, domain string) ([]*AgentCard, error) {
	agentCard, err := GetAgentCardFromDomain(domain)
	if err != nil {
		return nil, err
	}
	return ValidateAgentCard(agentID, address, agentCard), nil
}

func ValidateAgentCard(agentID string, address string, agentCard *AgentCard) []*AgentCard {
	var agentCards []*AgentCard
	for _, agentRegistration := range agentCard.Registrations {
		address, chainID, namespace, err := formatAddress(agentRegistration.AgentAddress)
		if err != nil {
			continue
		}
		if agentRegistration.AgentID == agentID && agentRegistration.AgentAddress == address {
			if !validateSignature(agentRegistration.Signature, address) {
				continue
			}
			newCard := formatAgentCard(agentCard)
			newCard.ChainID = chainID
			newCard.Namespace = namespace
			newCard.Domain = agentCard.Domain
			newCard.Signature = agentRegistration.Signature
			newCard.AgentID = agentRegistration.AgentID
			newCard.AgentAddress = agentRegistration.AgentAddress
			agentCards = append(agentCards, newCard)
		}
	}
	return agentCards
}

func GetAgentCardFromDomain(domain string) (*AgentCard, error) {
	response, err := http.Get(domain)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalAgentCard(body)
}

func unmarshalAgentCard(body []byte) (*AgentCard, error) {
	var agentCard *AgentCard
	err := json.Unmarshal(body, &agentCard)
	if err != nil {
		return nil, err
	}
	return agentCard, nil
}

var errInvalidAddressFormat = errors.New("invalid address format")

func formatAddress(address string) (string, string, string, error) {
	addressSlice := strings.Split(address, ":")
	if len(addressSlice) != 3 {
		return "", "", "", errInvalidAddressFormat
	}
	return addressSlice[0], addressSlice[1], addressSlice[2], nil
}

// todo: validate signature
func validateSignature(signature string, address string) bool {
	return true
}

func formatAgentCard(agentCard *AgentCard) *AgentCard {
	if agentCard == nil {
		return nil
	}

	// Deep copy embedded server.AgentCard via JSON round-trip to handle maps/slices/pointers safely.
	var embeddedCopy server.AgentCard
	if data, err := json.Marshal(agentCard.AgentCard); err == nil {
		if err := json.Unmarshal(data, &embeddedCopy); err != nil {
			// Fallback to shallow copy if JSON unmarshal fails.
			embeddedCopy = agentCard.AgentCard
		}
	} else {
		// Fallback to shallow copy if JSON marshal fails.
		embeddedCopy = agentCard.AgentCard
	}

	// Deep copy simple slices/structs defined in this package.
	var trustModelsCopy []string
	if len(agentCard.TrustModels) > 0 {
		trustModelsCopy = append(make([]string, 0, len(agentCard.TrustModels)), agentCard.TrustModels...)
	}

	var registrationsCopy []Registration
	if len(agentCard.Registrations) > 0 {
		registrationsCopy = make([]Registration, len(agentCard.Registrations))
		copy(registrationsCopy, agentCard.Registrations)
	}

	return &AgentCard{
		AgentCard:       embeddedCopy,
		TrustModels:     trustModelsCopy,
		Registrations:   registrationsCopy,
		FeedbackDataURI: agentCard.FeedbackDataURI,
		ChainID:         agentCard.ChainID,
		Namespace:       agentCard.Namespace,
	}
}
