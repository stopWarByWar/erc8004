package agentcard

import (
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

type AgentCard struct {
	server.AgentCard
	TrustModels     []string
	Registrations   []Registration
	FeedbackDataURI string
	ChainID         string
	Namespace       string
	Domain          string
	Signature       string
	AgentID         string
	AgentAddress    string
}

type Registration struct {
	AgentID      string
	AgentAddress string
	Signature    string
}

const TrustModelFeedback = "feedback"
const TrustModelInferenceValidation = "inference-validation"
const TrustModelTeeAttestation = "tee-attestation"
