package agentcard

import "encoding/json"

// Services 支持 JSON 中 services 为 string 数组或 Service 对象数组
// 例如: ["trading","defi"] 或 [{"name":"trading","endpoint":"..."}]
type Services []Service

func (s *Services) UnmarshalJSON(data []byte) error {
	var strSlice []string
	if err := json.Unmarshal(data, &strSlice); err == nil {
		*s = make([]Service, len(strSlice))
		for i, name := range strSlice {
			(*s)[i] = Service{Name: name}
		}
		return nil
	}
	var objSlice []Service
	if err := json.Unmarshal(data, &objSlice); err != nil {
		return err
	}
	*s = objSlice
	return nil
}

type TokenURLResponse struct {
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Image          string         `json:"image"`
	Services       Services       `json:"services"`
	X402Support    bool           `json:"x402Support"`
	Active         bool           `json:"active"`
	Registrations  []Registration `json:"registrations"`
	SupportedTrust []string       `json:"supportedTrust"`
}

type Registration struct {
	AgentID       any    `json:"agentId"`
	AgentRegistry string `json:"agentRegistry"`
}

type Service struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Version  *string `json:"version,omitempty"`
	// Capabilities map[string]interface{} `json:"capabilities,omitempty"`
	Skills  []string `json:"skills,omitempty"`
	Domains []string `json:"domains,omitempty"`
}

const TrustModelReputation = "reputation"
const TrustModelCryptoEconomicValidation = "crypto-economic-validation"
const TrustModelTeeAttestation = "tee-attestation"
