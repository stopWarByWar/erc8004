package model

type AgentCard struct {
	AgentID          string
	AgentDomain      string
	AgentAddress     string
	Name             string
	Description      string
	URL              string
	IconURL          string
	Version          string
	DocumentationURL string
	FeedbackDataURI  string
	ChainID          string
	Namespace        string
	Domain           string
	Signature        string
}

func (a AgentCard) AgentCard() string {
	return "agent_cards"
}

type Capability struct {
	AgentID                string
	Streaming              bool
	PushNotifications      bool
	StateTransitionHistory bool
}

func (a Capability) Capabilities() string {
	return "capabilities"
}

type Skill struct {
	AgentID     string
	Name        string
	Description string
	Tags        []string
	InputModes  []string
	OutputModes []string
}

func (a Skill) Skills() string {
	return "skills"
}

type Provider struct {
	AgentID      string
	Organization string
	URL          string
}

func (a Provider) Provider() string {
	return "providers"
}

type TrustModel struct {
	AgentID    string
	TrustModel []string
}

func (a TrustModel) TrustModels() string {
	return "trust_models"
}

type Extension struct {
	AgentID     string
	URI         string
	Required    bool
	Description string
}

func (a Extension) Extensions() string {
	return "extensions"
}

type AgentRegistry struct {
	AgentID      string
	AgentAddress string
	AgentDomain  string
	BlockNumber  uint64
	Index        uint64
	TxHash       string
	Timestamps   uint64
}

func (a AgentRegistry) AgentRegistries() string {
	return "agent_registries"
}
