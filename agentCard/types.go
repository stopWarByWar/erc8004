package agentcard

type TokenURLResponse struct {
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Image          string         `json:"image"`
	Services       []Service      `json:"services"`
	X402Support    bool           `json:"x402Support"`
	Active         bool           `json:"active"`
	Registrations  []Registration `json:"registrations"`
	SupportedTrust []string       `json:"supportedTrust"`
}

type Registration struct {
	AgentID       uint64 `json:"agentId"`
	AgentRegistry string `json:"agentRegistry"`
}

type Service struct {
	Name         string                 `json:"name"`
	Endpoint     string                 `json:"endpoint"`
	Version      *string                `json:"version,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
	Skills       []string               `json:"skills,omitempty"`
	Domains      []string               `json:"domains,omitempty"`
}

const TrustModelReputation = "reputation"
const TrustModelCryptoEconomicValidation = "crypto-economic-validation"
const TrustModelTeeAttestation = "tee-attestation"
