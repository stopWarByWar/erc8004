package api

type CardResponse struct {
	AgentID      string
	AgentDomain  string
	AgentAddress string
	ChainID      string
	Namespace    string

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
