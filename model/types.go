package model

type Agent struct {
	UID              uint64 `gorm:"column:uid;type:bigint;primaryKey"`
	AgentID          string `gorm:"column:agent_id"`
	Owner            string `gorm:"column:owner"`
	TokenURL         string `gorm:"column:token_url;type:text"`
	AgentWallet      string `gorm:"column:agent_wallet"`
	Namespace        string `gorm:"column:namespace;type:varchar(255)"`
	IdentityRegistry string `gorm:"column:identity_registry;type:varchar(255)"`
	ChainID          string `gorm:"column:chain_id;type:varchar(255)"`
	A2AEndpoint      string `gorm:"column:a2a_endpoint;type:text"`

	Type        string
	Name        string `gorm:"column:name;type:varchar(255)"`
	Description string `gorm:"column:description;type:text"`
	URL         string `gorm:"column:url;type:text"`

	Image            string `gorm:"column:image;type:varchar(255)"`
	Version          string `gorm:"column:version;type:varchar(255)"`
	DocumentationURL string `gorm:"column:documentation_url;type:text"`

	Score        uint64 `gorm:"column:score;type:bigint"`
	CommentCount uint64 `gorm:"column:comment_count;type:bigint"`
	Timestamps   uint64 `gorm:"column:timestamps;type:bigint"`

	UserInterfaceURL string `gorm:"column:user_interface_url;type:varchar(500)"`
}

func (Agent) TableName() string { return "agents" }

type Capability struct {
	AgentUID               uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Streaming              bool   `gorm:"column:streaming;type:boolean"`
	PushNotifications      bool   `gorm:"column:push_notifications;type:boolean"`
	StateTransitionHistory bool   `gorm:"column:state_transition_history;type:boolean"`
}

func (Capability) TableName() string { return "capabilities" }

type Skill struct {
	AgentUID    uint64 `gorm:"column:agent_uid;primaryKey"`
	ID          string `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name;type:varchar(255)"`
	Description string `gorm:"column:description;type:text"`
}

func (Skill) TableName() string { return "skills" }

type SkillTags struct {
	AgentUID uint64 `gorm:"column:agent_uid;type:bigint;primaryKey;"`
	SkillID  string `gorm:"column:id;type:varchar(255);primaryKey"`
	Tag      string `gorm:"column:tag;type:varchar(255);primaryKey"`
}

func (SkillTags) TableName() string { return "skill_tags" }

type Provider struct {
	AgentUID     uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Organization string `gorm:"column:organization;type:varchar(255);primaryKey"`
	URL          string `gorm:"column:url;type:varchar(255)"`
}

func (Provider) TableName() string { return "providers" }

type TrustModel struct {
	AgentUID   uint64 `gorm:"column:agent_uid;type:bigint;primaryKey;index"`
	TrustModel string `gorm:"column:trust_model;type:varchar(255);primaryKey"`
}

func (TrustModel) TableName() string { return "trust_models" }

type Extension struct {
	AgentUID    uint64 `gorm:"column:agent_uid;type:bigint;primaryKey;index"`
	URI         string `gorm:"column:uri;type:text;primaryKey"`
	Required    bool   `gorm:"column:required;type:boolean"`
	Description string `gorm:"column:description;type:text"`
}

func (Extension) TableName() string { return "extensions" }

type AgentRegistry struct {
	AgentID          string `gorm:"column:agent_id;type:varchar(255)"`
	IdentityRegistry string `gorm:"column:identity_registry;type:char(42)"`
	Owner            string `gorm:"column:owner;type:varchar(255)"`
	TokenURL         string `gorm:"column:token_url;type:text"`
	ChainID          string `gorm:"column:chain_id;type:char(42)"`
	BlockNumber      uint64 `gorm:"column:block_number;type:bigint"`
	Index            uint64 `gorm:"column:index;type:bigint"`
	TxHash           string `gorm:"column:tx_hash;type:char(66);primaryKey"`
	Timestamps       uint64 `gorm:"column:timestamps;type:bigint"`
	Inserted         bool   `gorm:"column:inserted;type:boolean"`
}

func (AgentRegistry) TableName() string { return "agent_registries" }

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

type Feedback struct {
	UID                uint64 `gorm:"column:uid;type:bigint;primaryKey"`
	AgentUID           uint64 `gorm:"column:agent_uid"`
	ChainID            string `gorm:"column:chain_id"`
	AgentID            string `gorm:"column:agent_id"`
	ReputationRegistry string `gorm:"column:reputation_registry"`
	ClientAddress      string `gorm:"column:client_address"`
	FeedbackIndex      uint64 `gorm:"column:feedback_index"`
	Score              uint8
	Tag1               string
	Tag2               string
	FeedbackURI        string
	FeedbackHash       string

	Fetched     bool `gorm:"column:fetched;type:boolean"`
	BlockNumber uint64
	Index       uint64
	TxHash      string
	Revoked     bool   `gorm:"column:revoked;type:boolean"`
	Timestamps  uint64 `gorm:"column:timestamps;type:bigint"`
}

func (Feedback) TableName() string { return "feedbacks" }

type Response struct {
	UID                uint64 `gorm:"column:uid;primaryKey"`
	AgentUID           uint64 `gorm:"column:agent_uid"`
	FeedbackUID        uint64 `gorm:"column:feedback_uid"`
	ChainID            string `gorm:"column:chain_id"`
	AgentID            string `gorm:"column:agent_id"`
	ReputationRegistry string `gorm:"column:reputation_registry"`
	ClientAddress      string `gorm:"column:client_address"`
	FeedbackIndex      uint64 `gorm:"column:feedback_index"`
	Responder          string
	ResponseURI        string
	ResponseHash       string

	Fetched     bool `gorm:"column:fetched;type:boolean"`
	BlockNumber uint64
	Index       uint64
	TxHash      string
	Timestamps  uint64 `gorm:"column:timestamps;type:bigint"`
}

func (Response) TableName() string { return "responses" }

// agent comment schema
// identity_registry char(42) NOT NULL,
// agent_id varchar(255) NOT NULL,
// commenter char(42) NOT NULL,
// comment_text text NOT NULL,

type AgentComment struct {
	CommentAttestationID string `gorm:"column:comment_attestation_id;type:bytes32;primaryKey"`
	AgentUID             uint64 `gorm:"column:agent_uid;type:bigint"`
	ChainID              string `gorm:"column:chain_id;type:char(42)"`
	IdentityRegistry     string `gorm:"column:identity_registry;type:char(42)"`
	AgentID              string `gorm:"column:agent_id;type:varchar(255)"`
	Commenter            string `gorm:"column:commenter;type:char(42)"`
	CommentText          string `gorm:"column:comment_text;type:text"`
	Score                uint16 `gorm:"column:score;type:uint16"`
	Timestamps           uint64 `gorm:"column:timestamps;type:bigint"`
	Block                uint64 `gorm:"column:block;type:bigint"`
	Index                uint64 `gorm:"column:index;type:bigint"`
	TxHash               string `gorm:"column:tx_hash;type:char(66)"`
	Revoked              bool   `gorm:"column:revoked;type:boolean"`
}

func (AgentComment) TableName() string { return "agent_comments" }

type Comment struct {
	CommentAttestationID string `gorm:"column:comment_attestation_id"`
	Commenter            string `gorm:"column:commenter;type:varchar(255)"`
	AgentUID             uint64 `gorm:"column:agent_uid;type:bigint"`
	CommentText          string `gorm:"column:comment_text;type:text"`
	Score                uint16 `gorm:"column:score;type:uint16"`
	Timestamps           uint64 `gorm:"column:timestamps;type:bigint"`

	Avatar   string `gorm:"column:avatar;type:varchar(255)"`
	Name     string `gorm:"column:name;type:varchar(255)"`
	Passport bool
}

type FeedbackResp struct {
	UID                uint64 `gorm:"column:uid;type:bigint;primaryKey"`
	AgentUID           uint64 `gorm:"column:agent_uid"`
	ChainID            string `gorm:"column:chain_id"`
	AgentID            string `gorm:"column:agent_id"`
	ReputationRegistry string `gorm:"column:reputation_registry"`
	ClientAddress      string `gorm:"column:client_address"`
	FeedbackIndex      uint64 `gorm:"column:feedback_index"`
	Score              uint8
	Tag1               string
	Tag2               string
	FeedbackURI        string
	FeedbackHash       string

	TxHash     string
	Timestamps uint64 `gorm:"column:timestamps;type:bigint"`

	Name     string
	Avatar   string
	Passport bool
}

type Metadata struct {
	ChainID          string `gorm:"column:chain_id;type:char(42);primaryKey"`
	IdentityRegistry string `gorm:"column:identity_registry;type:char(42);primaryKey"`
	AgentID          string `gorm:"column:agent_id;type:varchar(255);primaryKey"`
	Key              string `gorm:"column:key;type:varchar(255);primaryKey"`
	Value            string `gorm:"column:value;type:text"`
	Block            uint64 `gorm:"column:block;type:bigint"`
	Index            uint64 `gorm:"column:index;type:bigint"`
	TxHash           string `gorm:"column:tx_hash;type:char(66)"`
}

func (Metadata) TableName() string { return "agent_metadatas" }
