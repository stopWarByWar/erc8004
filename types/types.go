package types

type Endpoint struct {
	Name         string                 `json:"name"`
	Endpoint     string                 `json:"endpoint"`
	Version      string                 `json:"version,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
}

type Registration struct {
	AgentId       int64  `json:"agentId"`
	AgentRegistry string `json:"agentRegistry"`
}

type AgentProfile struct {
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Image          string         `json:"image"`
	Endpoints      []Endpoint     `json:"endpoints"`
	Registrations  []Registration `json:"registrations"`
	SupportedTrust []string       `json:"supportedTrust"`
}

type Feedback struct {
	// MUST fields
	AgentRegistry string `json:"agentRegistry"`
	AgentId       int64  `json:"agentId"`
	ClientAddress string `json:"clientAddress"`
	CreatedAt     string `json:"createdAt"`
	Score         int    `json:"score"`

	// MAY fields
	Tag1           *string         `json:"tag1,omitempty"`
	Tag2           *string         `json:"tag2,omitempty"`
	Endpoint       *string         `json:"endpoint,omitempty"`
	Domain         *string         `json:"domain,omitempty"`
	Skill          *string         `json:"skill,omitempty"`
	Context        *string         `json:"context,omitempty"`
	Task           *string         `json:"task,omitempty"`
	Capability     *string         `json:"capability,omitempty"`
	Name           *string         `json:"name,omitempty"`
	ProofOfPayment *ProofOfPayment `json:"proofOfPayment,omitempty"`
}

type ProofOfPayment struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	ChainId     string `json:"chainId"`
	TxHash      string `json:"txHash"`
}

type NetworkResponse struct {
	ChainId   string `json:"chainId"`
	ChainName string `json:"chainName"`
	ChainLogo string `json:"chainLogo"`

	ContractInfo []ContractInfo `json:"contractInfo"`
}

type ContractInfo struct {
	IdentityAddress       string `json:"identityAddress"`
	IdentityContractURL   string `json:"identityContractURL"`
	ReputationAddress     string `json:"reputationAddress"`
	ReputationContractURL string `json:"reputationContractURL"`
	ValidationAddress     string `json:"validationAddress"`
	ValidationContractURL string `json:"validationContractURL"`
	Deployer              string `json:"deployer"`
	Description           string `json:"description"`
	LogoURL               string `json:"logoURL"`
}
