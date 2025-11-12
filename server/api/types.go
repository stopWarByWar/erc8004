package api

type AgentResponse struct {
	UID              uint64 // uid
	AgentID          string `json:"agentId"` // agent id in contract json file
	AgentDomain      string //a2a endpoint
	AgentAddress     string //wallet address
	Owner            string `json:"owner"`
	ChainID          string
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
