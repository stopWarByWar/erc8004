package agentcard

import (
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

type TokenURLResponse struct {
	Type             string         `json:"type"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Image            string         `json:"image"`
	Endpoints        []Endpoint     `json:"endpoints"`
	Registrations    []Registration `json:"registrations"`
	SupportedTrust   []string       `json:"supportedTrust"`
	UserInterfaceURL string         `json:"userInterface"`
}
type Registration struct {
	AgentID       uint64 `json:"agentId"`
	AgentRegistry string `json:"agentRegistry"`
}

type Endpoint struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}

const TrustModelFeedback = "feedback"
const TrustModelInferenceValidation = "inference-validation"
const TrustModelTeeAttestation = "tee-attestation"

type Agent struct {
	Owner            string            `json:"owner"`
	Type             string            `json:"type"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Image            string            `json:"image"`
	Endpoint         string            `json:"endpoint"`
	IdentityRegistry string            `json:"registry"`
	SupportedTrust   []string          `json:"supportedTrust"`
	AgentID          string            `json:"agentId"`
	TokenURL         string            `json:"tokenUrl"`
	AgentCard        *server.AgentCard `json:"agentCard"`
	Namespace        string            `json:"namespace"`
	ChainID          string            `json:"chainId"`
	AgentWallet      string            `json:"agentWallet"`
	Timestamps       uint64            `json:"timestamps"`
	UserInterfaceURL string            `json:"userInterfaceURL"`
}
