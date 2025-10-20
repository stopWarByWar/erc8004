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
	Score            uint64 `gorm:"column:score;type:bigint"`
	CommentCount     uint64 `gorm:"column:comment_count;type:bigint"`
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

type AgentComment struct {
	CommentAttestationID string `gorm:"column:comment_attestation_id;type:bytes32;primaryKey"`
	AgentClientID        string `gorm:"column:agent_client_id;type:varchar(255)"`
	AgentServerID        string `gorm:"column:agent_server_id;type:varchar(255)"`
	CommentText          string `gorm:"column:comment_text;type:text"`
	Score                uint16 `gorm:"column:score;type:uint16"`
	Commenter            string `gorm:"column:commenter;type:varchar(255)"`
	Timestamps           uint64 `gorm:"column:timestamps;type:bigint"`
	NewestComment        bool   `gorm:"column:newest_comment;type:boolean"`
	Block                uint64 `gorm:"column:block;type:bigint"`
	Index                uint64 `gorm:"column:index;type:bigint"`
	TxHash               string `gorm:"column:tx_hash;type:varchar(255)"`
	IsAuthorized         bool   `gorm:"column:is_authorized;type:boolean"`
}

func (AgentComment) TableName() string { return "agent_comments" }

type Attestation struct {
	UID                   string `gorm:"column:uid;primaryKey"`
	Type                  string
	SchemaUID             string `gorm:"column:schema_uid"`
	Attestor              string
	Recipient             string `gorm:"column:recipient"`
	Timestamps            int64
	Expiration            int64
	Revoked               bool
	Revocable             bool   `gorm:"column:revocable"`
	RevokeTxHash          string `gorm:"column:revoke_tx_hash"`
	RevocationTime        int64
	TransactionId         string `gorm:"column:transaction_id"`
	ReferencedAttestation string `gorm:"column:referenced_attestation"`
	RawData               []byte `gorm:"column:raw_data"`
	Public                bool

	Block int64
	Index int
}

func (Attestation) TableName() string {
	return "attestation"
}

type AuthFeedback struct {
	AuthFeedbackID     string `gorm:"column:auth_feedback_id;type:varchar(255);primaryKey"`
	AgentClientID      string `gorm:"column:agent_client_id;type:varchar(255)"`
	AgentServerID      string `gorm:"column:agent_server_id;type:varchar(255)"`
	AgentClientAddress string `gorm:"column:agent_client_address;type:varchar(255)"`
	AgentServerAddress string `gorm:"column:agent_server_address;type:varchar(255)"`
	BlockNumber        int64  `gorm:"column:block_number;type:bigint"`
	Index              int64  `gorm:"column:index;type:bigint"`
	TxHash             string `gorm:"column:tx_hash;type:varchar(255)"`
}

func (AuthFeedback) TableName() string { return "auth_feedbacks" }

type Comment struct {
	Commenter     string `gorm:"column:commenter;type:varchar(255)"`
	AgentClientID string `gorm:"column:agent_client_id;type:varchar(255)"`
	CommentText   string `gorm:"column:comment_text;type:text"`
	Score         uint16 `gorm:"column:score;type:uint16"`
	Timestamps    uint64 `gorm:"column:timestamps;type:bigint"`

	Logo         string `gorm:"column:logo;type:varchar(255)"`
	Name         string `gorm:"column:name;type:varchar(255)"`
	Passport     bool
	IsAuthorized bool `gorm:"column:is_authorized;type:boolean"`
}
