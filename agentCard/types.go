package agentcard

import (
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

type AgentCard struct {
	server.AgentCard
	TrustModels     []string       `json:"trustModels"`
	Registrations   []Registration `json:"registrations"`
	FeedbackDataURI string         `json:"FeedbackDataURI"`
	ChainID         string         `json:"chainId"`
	Namespace       string         `json:"namespace"`
	Domain          string         `json:"domain"`
	Signature       string         `json:"signature"`
	AgentID         string         `json:"agentId"`
	AgentAddress    string         `json:"agentAddress"`
	UserInterface   string         `json:"userInterface"`
}

type Registration struct {
	AgentID      string `json:"agentId"`
	AgentAddress string `json:"agentAddress"`
	Signature    string `json:"signature"`
}

const TrustModelFeedback = "feedback"
const TrustModelInferenceValidation = "inference-validation"
const TrustModelTeeAttestation = "tee-attestation"
