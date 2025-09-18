package server

import "agent_identity/model"

type CardResponse struct {
	AgentID      string
	AgentDomain  string
	AgentAddress string

	Name             string              `json:"name"`
	Description      string              `json:"description"`
	URL              string              `json:"url"`
	Provider         *model.Provider     `json:"provider"`
	IconURL          string              `json:"iconUrl,omitempty"`
	Version          string              `json:"version"`
	DocumentationURL string              `json:"documentationUrl,omitempty"`
	Skills           []*model.Skill      `json:"skills"`
	TrustModels      []*model.TrustModel `json:"trustModels"`
}
