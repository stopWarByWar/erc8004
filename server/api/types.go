package api

type AgentResponse struct {
	UID              uint64 // uid
	AgentID          string `json:"agentId"` // agent id in contract json file
	AgentDomain      string //a2a endpoint
	AgentAddress     string //wallet address
	Owner            string `json:"owner"`
	ChainID          string
	ChainName        string `json:"chainName"`
	ChainLogo        string `json:"chainLogo"`
	Namespace        string
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
	Score            float64            `json:"score"`
	Metadata         []MetadataResponse `json:"metadata"`
	TokenURL         string             `json:"tokenUrl"`
	Deployer         string             `json:"deployer"`
	DeployerLogo     string             `json:"deployerLogo"`
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
	UID                   uint64       `json:"uid"`
	Score                 int          `json:"score"`
	Tag1                  string       `json:"tag1"`
	Tag2                  string       `json:"tag2"`
	Skill                 string       `json:"skill,omitempty"`
	Context               string       `json:"context,omitempty"`
	Task                  string       `json:"task,omitempty"`
	Capability            string       `json:"capability,omitempty"`
	FeedbackAuth          FeedbackAuth `json:"feedbackAuth"`
	FeedbackAuthSignature string       `json:"feedbackAuthSignature"`
}

type FeedbackAuth struct {
	AgentId          string `json:"agentId"`
	ClientAddress    string `json:"clientAddress"`
	IndexLimit       uint64 `json:"indexLimit"`
	Expiry           uint64 `json:"expiry"`
	ChainId          string `json:"chainId"`
	IdentityRegistry string `json:"identityRegistry"`
	SignerAddress    string `json:"signerAddress"`
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
