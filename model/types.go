package model

type AgentCard struct {
	AgentID          string `gorm:"column:agent_id;type:varchar(255);primaryKey"`
	AgentDomain      string `gorm:"column:agent_domain;type:text"`
	AgentAddress     string `gorm:"column:agent_address;type:text"`
	Name             string `gorm:"column:name;type:varchar(255)"`
	Description      string `gorm:"column:description;type:text"`
	URL              string `gorm:"column:url;type:text"`
	IconURL          string `gorm:"column:icon_url;type:varchar(255)"`
	Version          string `gorm:"column:version;type:varchar(255)"`
	DocumentationURL string `gorm:"column:documentation_url;type:text"`
	FeedbackDataURI  string `gorm:"column:feedback_data_uri;type:text"`
	ChainID          string `gorm:"column:chain_id;type:varchar(255)"`
	Namespace        string `gorm:"column:namespace;type:varchar(255)"`
	Signature        string `gorm:"column:signature;type:varchar(255)"`
	UserInterface    string `gorm:"column:user_interface;type:text"`
}

func (AgentCard) TableName() string { return "agent_cards" }

type Capability struct {
	AgentID                string `gorm:"column:agent_id;type:varchar(255);primaryKey"`
	Streaming              bool   `gorm:"column:streaming;type:boolean"`
	PushNotifications      bool   `gorm:"column:push_notifications;type:boolean"`
	StateTransitionHistory bool   `gorm:"column:state_transition_history;type:boolean"`
}

func (Capability) TableName() string { return "capabilities" }

type Skill struct {
	AgentID     string `gorm:"column:agent_id;type:varchar(255);primaryKey;"`
	ID          string `gorm:"column:id;type:varchar(255);primaryKey"`
	Name        string `gorm:"column:name;type:varchar(255)"`
	Description string `gorm:"column:description;type:text"`
}

func (Skill) TableName() string { return "skills" }

type SkillTags struct {
	AgentID string `gorm:"column:agent_id;type:varchar(255);primaryKey;"`
	SkillID string `gorm:"column:id;type:varchar(255);primaryKey"`
	Tag     string `gorm:"column:tag;type:varchar(255);primaryKey"`
}

func (SkillTags) TableName() string { return "skill_tags" }

type Provider struct {
	AgentID      string `gorm:"column:agent_id;type:varchar(255);primaryKey"`
	Organization string `gorm:"column:organization;type:varchar(255);primaryKey"`
	URL          string `gorm:"column:url;type:varchar(255)"`
}

func (Provider) TableName() string { return "providers" }

type TrustModel struct {
	AgentID    string `gorm:"column:agent_id;type:varchar(255);primaryKey;index"`
	TrustModel string `gorm:"column:trust_model;type:varchar(255);primaryKey"`
}

func (TrustModel) TableName() string { return "trust_models" }

type Extension struct {
	AgentID     string `gorm:"column:agent_id;type:varchar(255);primaryKey;index"`
	URI         string `gorm:"column:uri;type:text;primaryKey"`
	Required    bool   `gorm:"column:required;type:boolean"`
	Description string `gorm:"column:description;type:text"`
}

func (Extension) TableName() string { return "extensions" }

type AgentRegistry struct {
	AgentID      string `gorm:"column:agent_id;type:varchar(255);index"`
	AgentAddress string `gorm:"column:agent_address;type:varchar(255);index"`
	AgentDomain  string `gorm:"column:agent_domain;type:varchar(255)"`
	BlockNumber  uint64 `gorm:"column:block_number;type:bigint;primaryKey"`
	Index        uint64 `gorm:"column:index;type:bigint;primaryKey"`
	TxHash       string `gorm:"column:tx_hash;type:varchar(255)"`
	Timestamps   uint64 `gorm:"column:timestamps;type:bigint"`
	Inserted     bool   `gorm:"column:inserted;type:boolean"`
}

func (AgentRegistry) TableName() string { return "agent_registries" }
