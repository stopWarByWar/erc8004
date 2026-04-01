package types

import "time"

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
	// Namespace            string
	IdentityRegistry string             `json:"identityRegistry"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	URL              string             `json:"url"`
	Provider         ProviderResponse   `json:"provider"`
	IconURL          string             `json:"iconUrl,omitempty"`
	Version          string             `json:"version"`
	DocumentationURL string             `json:"documentationUrl,omitempty"`
	Skills           []SkillTagResponse `json:"skills"`
	TrustModels      []string           `json:"trustModels"`
	UserInterface    string             `json:"userInterface"`
	// Score              float64            `json:"score"`
	Metadata           []MetadataResponse `json:"metadata"`
	TokenURL           string             `json:"tokenUrl"`
	Deployer           string             `json:"deployer"`
	DeployerLogo       string             `json:"deployerLogo"`
	MCPEndpoint        string             `json:"mcpEndpoint,omitempty"`
	OASFEndpoint       string             `json:"oasfEndpoint,omitempty"`
	ReputationRegistry string             `json:"reputationRegistry"`
	Status             string             `json:"status"`
	X402Support        bool               `json:"x402Support"`
	Endpoints          []EndpointResponse `json:"endpoints"`
}

type EndpointResponse struct {
	No       int    `json:"no"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
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
	AgentID          string     `json:"agentId" form:"agentId" binding:"required"`
	ChainID          string     `json:"chainId" form:"chainId" binding:"required"`
	Name             string     `json:"name" form:"name" binding:"required"`
	Description      string     `json:"description" form:"description" binding:"required"`
	IdentityRegistry string     `json:"identityRegistry" form:"identityRegistry" binding:"required"`
	SupportedTrust   []string   `json:"supportedTrust" form:"supportedTrust"`
	Endpoints        []Endpoint `json:"endpoints" form:"endpoints" binding:"required"`
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
	ValidatorLogo    string `json:"validator_logo"`

	RequestURI  string `json:"request_uri"`
	RequestHash string `json:"request_hash"`

	Response     int    `json:"response"`
	ResponseURI  string `json:"response_uri"`
	ResponseHash string `json:"response_hash"`
	Status       string `json:"status"`
	Tag1         string `json:"tag1"`
	Timestamps   uint64 `json:"timestamps"`
}

type ValidatorValidation struct {
	ChainName        string `json:"chain_name"`
	ChainLogo        string `json:"chain_logo"`
	ContractDeployer string `json:"contract_deployer"`

	AgentUID  uint64 `json:"agent_uid"`
	ChainID   string `json:"chain_id"`
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`

	ValidationRegistry string `json:"validation_registry"`

	ValidatorAddress string `json:"validator_address"`

	RequestURI  string `json:"request_uri"`
	RequestHash string `json:"request_hash"`

	Response     int    `json:"response"`
	ResponseURI  string `json:"response_uri"`
	ResponseHash string `json:"response_hash"`
	Tag1         string `json:"tag1"`
	Timestamps   uint64 `json:"timestamps"`
	Status       string `json:"status"`
}

type Endpoint struct {
	Name         string                  `json:"name"`
	Endpoint     string                  `json:"endpoint"`
	Version      *string                 `json:"version,omitempty"`
	Capabilities *map[string]interface{} `json:"capabilities,omitempty"`
	Skills       []string                `json:"skills,omitempty"`
	Domains      []string                `json:"domains,omitempty"`
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
	Value         int    `json:"value"`
	ValueDecimals int    `json:"valueDecimals"`
	// MAY fields
	Tag1     *string `json:"tag1,omitempty"`
	Tag2     *string `json:"tag2,omitempty"`
	Endpoint *string `json:"endpoint,omitempty"`
	Context  *string `json:"context,omitempty"`

	MCP            *map[string]interface{} `json:"mcp,omitempty"`
	A2A            *A2A                    `json:"a2a,omitempty"`
	OASF           *OASF                   `json:"oasf,omitempty"`
	ProofOfPayment *ProofOfPayment         `json:"proofOfPayment,omitempty"`
}

type A2A struct {
	Skills    []string `json:"skills,omitempty"`
	ContextId string   `json:"contextId,omitempty"`
	TaskId    string   `json:"taskId,omitempty"`
}

type OASF struct {
	Skills  []string `json:"skills,omitempty"`
	Domains []string `json:"domains,omitempty"`
}

type ProofOfPayment struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	ChainId     string `json:"chainId"`
	TxHash      string `json:"txHash"`
}

type NetworkResponse struct {
	ChainId     string `json:"chainId"`
	ChainName   string `json:"chainName"`
	ChainLogo   string `json:"chainLogo"`
	AgentAmount uint64 `json:"agentAmount"`

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

type ValidatorListInfo struct {
	Rank           uint64 `json:"rank"`
	Address        string `json:"address"`
	PendingAmount  uint64 `json:"pendingAmount"`
	FinishedAmount uint64 `json:"finishedAmount"`
}

type LeaderboardInfo struct {
	AgentAmount            int64
	FeedbackAmount         int64
	AgentAmountWithIn7Days int64
	NetworkAmount          int64
	Leaderboard            []LeaderboardAgentInfo
}

type LeaderboardAgentInfo struct {
	Name string
	Key  string
	Data []SimpleAgentInfo
}

type SimpleAgentInfo struct {
	UID              uint64
	AgentID          string
	AgentName        string
	AgentImage       string
	AgentDescription string
	ChainID          string
	ChainName        string
	ChainLogo        string
}

type AgentValidationEvalReportResponse struct {
	ID           uint      `json:"id"`
	AgentUID     uint64    `json:"agent_uid"`
	AgentTokenID int       `json:"agent_token_id"`
	ChainID      int       `json:"chainid"`
	Reporter     string    `json:"reporter"`
	Version      string    `json:"version"`
	Score        float64   `json:"score"`
	ReportURL    string    `json:"report_url"`
	Desc         string    `json:"desc"`
	ValidatedAt  time.Time `json:"validated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type AgentValidationEvalDimensionResponse struct {
	Dimension string    `json:"dimension"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type AgentValidationEvalLatestResponse struct {
	Report     *AgentValidationEvalReportResponse     `json:"report"`
	Dimensions []AgentValidationEvalDimensionResponse `json:"dimensions"`
}
