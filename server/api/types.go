package api

type AgentResponse struct {
	UID                  uint64 // uid
	AgentID              string `json:"agentId"` // agent id in contract json file
	A2AEndpoint          string `json:"a2aEndpoint"`
	WalletAddress        string `json:"walletAddress"`        //wallet address
	WalletAddressScanURL string `json:"walletAddressScanURL"` //wallet address scan url
	Owner                string `json:"owner"`
	ChainID              string
	ChainName            string `json:"chainName"`
	ChainLogo            string `json:"chainLogo"`
	Namespace            string
	IdentityRegistry     string             `json:"identityRegistry"`
	Name                 string             `json:"name"`
	Description          string             `json:"description"`
	URL                  string             `json:"url"`
	Provider             ProviderResponse   `json:"provider"`
	IconURL              string             `json:"iconUrl,omitempty"`
	Version              string             `json:"version"`
	DocumentationURL     string             `json:"documentationUrl,omitempty"`
	Skills               []SkillTagResponse `json:"skills"`
	TrustModels          []string           `json:"trustModels"`
	UserInterface        string             `json:"userInterface"`
	Score                float64            `json:"score"`
	Metadata             []MetadataResponse `json:"metadata"`
	TokenURL             string             `json:"tokenUrl"`
	Deployer             string             `json:"deployer"`
	DeployerLogo         string             `json:"deployerLogo"`
	MACEndpoint          string             `json:"mcpEndpoint,omitempty"`
	OASFEndpoint         string             `json:"oasfEndpoint,omitempty"`

	WalletAddressExpirationTime uint64 `json:"walletAddressExpirationTime"`
	ReputationRegistry          string `json:"reputationRegistry"`
}

type SkillTagResponse struct {
	ID          string
	Name        string
	Description string
	Tags        []string
}

type ProviderResponse struct {
	Organization string
	URL          string
}
type ProofOfPayment struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	ChainId     string `json:"chainId"`
	TxHash      string `json:"txHash"`
}
type UploadFeedbackRequest struct {
	UID            uint64          `json:"uid"`
	ClientAddress  string          `json:"clientAddress"`
	IndexLimit     uint64          `json:"indexLimit"`
	Score          int             `json:"score"`
	Tag1           *string         `json:"tag1,omitempty"`
	Tag2           *string         `json:"tag2,omitempty"`
	Skill          *string         `json:"skill,omitempty"`
	Context        *string         `json:"context,omitempty"`
	Task           *string         `json:"task,omitempty"`
	Capability     *string         `json:"capability,omitempty"`
	Endpoint       *string         `json:"endpoint,omitempty"`
	Domain         *string         `json:"domain,omitempty"`
	Name           *string         `json:"name,omitempty"`
	ProofOfPayment *ProofOfPayment `json:"proofOfPayment,omitempty"`
}
type UploadAgentProfileRequest struct {
	AgentID            string   `json:"agentId" form:"agentId" binding:"required"`
	ChainID            string   `json:"chainId" form:"chainId" binding:"required"`
	Name               string   `json:"name" form:"name" binding:"required"`
	Description        string   `json:"description" form:"description" binding:"required"`
	A2AEndpoint        string   `json:"a2aEndpoint" form:"a2aEndpoint" binding:"required"`
	IdentityRegistry   string   `json:"identityRegistry" form:"identityRegistry" binding:"required"`
	SupportedTrust     []string `json:"supportedTrust" form:"supportedTrust"`
	AgentWalletAddress string   `json:"agentWallet" form:"agentWallet" binding:"required"`
}

type MetadataResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Validation struct {
	AgentUID           uint64 `json:"agent_uid"`
	ChainID            string `json:"chain_id"`
	AgentID            string `json:"agent_id"`
	ValidationRegistry string `json:"validation_registry"`

	ValidatorAddress string `json:"validator_address"`

	RequestURI  string `json:"request_uri"`
	RequestHash string `json:"request_hash"`

	Response     int    `json:"response"`
	ResponseURI  string `json:"response_uri"`
	ResponseHash string `json:"response_hash"`
	Tag1         string `json:"tag1"`
	Timestamps   uint64 `json:"timestamps"`
}
