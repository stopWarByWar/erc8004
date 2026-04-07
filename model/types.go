package model

type Agent struct {
	UID              uint64 `gorm:"column:uid;type:bigint;primaryKey"`
	ChainID          string `gorm:"column:chain_id;type:varchar(255)"`
	IdentityRegistry string `gorm:"column:identity_registry;type:varchar(255)"`
	AgentID          string `gorm:"column:agent_id"`
	Owner            string `gorm:"column:owner"`

	AgentURI    string `gorm:"column:agent_uri;type:text"`
	AgentWallet string `gorm:"column:agent_wallet"`

	Type        string
	Name        string `gorm:"column:name;type:varchar(255)"`
	Description string `gorm:"column:description;type:text"`

	Image       string `gorm:"column:image;type:varchar(255)"`
	X402Support bool   `gorm:"column:x402_support;type:boolean"`
	Active      bool   `gorm:"column:active;type:boolean"`

	FeedbackCount uint64 `gorm:"column:feedback_count;type:bigint"`

	BlockNumber uint64 `gorm:"column:block_number;type:bigint"`
	Index       uint64 `gorm:"column:index;type:bigint"`
	TxHash      string `gorm:"column:tx_hash;type:char(66)"`
	Timestamps  uint64 `gorm:"column:timestamps;type:bigint"`

	A2AURI              string `gorm:"column:a2a_uri;type:text"`
	A2AVersion          string `gorm:"column:a2a_version;type:varchar(255)"`
	A2ADocumentationURL string `gorm:"column:a2a_documentation_url;type:text"`

	Inserted bool `gorm:"column:inserted;type:boolean"`
}

func (Agent) TableName() string { return "agents" }

type Service struct {
	AgentUID    uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	ServiceName string `gorm:"column:service_name;type:varchar(255)"`
	Endpoint    string `gorm:"column:endpoint;type:text"`
	Version     string `gorm:"column:version;type:varchar(255)"`
}

func (Service) TableName() string { return "services" }

type OASFSkill struct {
	AgentUID  uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	SkillName string `gorm:"column:skill_name;type:varchar(255)"`
	Version   string `gorm:"column:version;type:varchar(255)"`
}

func (OASFSkill) TableName() string { return "oasf_skills" }

type OASFDomain struct {
	AgentUID uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Domain   string `gorm:"column:domain;type:varchar(255)"`
	Version  string `gorm:"column:version;type:varchar(255)"`
}

func (OASFDomain) TableName() string { return "oasf_domains" }

type TrustModel struct {
	AgentUID   uint64 `gorm:"column:agent_uid;type:bigint;primaryKey;index"`
	TrustModel string `gorm:"column:trust_model;type:varchar(255);primaryKey"`
}

func (TrustModel) TableName() string { return "trust_models" }

type AgentVector struct {
	UID              uint64 `gorm:"column:uid;type:bigint;primaryKey"`
	AgentUID         uint64 `gorm:"column:agent_uid;type:bigint;not null;index"`
	IdentityRegistry string `gorm:"column:identity_registry;type:varchar(255);not null;index"`
	ChainID          string `gorm:"column:chain_id;type:varchar(255);not null;index"`
	CreateTimestamp  uint64 `gorm:"column:create_timestamp;type:bigint;not null;index"`
	Embedding        string `gorm:"column:embedding;type:vector(1536);not null"` // 使用string类型存储，实际使用时需要转换为[]float32
	Content          string `gorm:"column:content;type:text"`
	Metadata         string `gorm:"column:metadata;type:jsonb"` // 使用string存储JSON，也可以使用自定义类型
	CreatedAt        int64  `gorm:"column:created_at;type:timestamp"`
}

func (AgentVector) TableName() string { return "agent_vectors" }

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

//-------------A2A-------------

type Capability struct {
	AgentUID               uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Streaming              bool   `gorm:"column:streaming;type:boolean"`
	PushNotifications      bool   `gorm:"column:push_notifications;type:boolean"`
	StateTransitionHistory bool   `gorm:"column:state_transition_history;type:boolean"`
}

func (Capability) TableName() string { return "capabilities" }

type A2ASkill struct {
	AgentUID    uint64 `gorm:"column:agent_uid;primaryKey"`
	ID          string `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name;type:varchar(255)"`
	Description string `gorm:"column:description;type:text"`
}

func (A2ASkill) TableName() string { return "a2a_skills" }

type A2ASkillTag struct {
	AgentUID uint64 `gorm:"column:agent_uid;type:bigint;primaryKey;"`
	ID       string `gorm:"column:id;type:varchar(255);primaryKey"`
	Tag      string `gorm:"column:tag;type:varchar(255);primaryKey"`
}

func (A2ASkillTag) TableName() string { return "a2a_skill_tags" }

type A2AProvider struct {
	AgentUID     uint64 `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Organization string `gorm:"column:organization;type:varchar(255);primaryKey"`
	URL          string `gorm:"column:url;type:varchar(255)"`
}

func (A2AProvider) TableName() string { return "a2a_providers" }

type A2AExtension struct {
	AgentUID    uint64 `gorm:"column:agent_uid;type:bigint;primaryKey;index"`
	URI         string `gorm:"column:uri;type:text;primaryKey"`
	Required    bool   `gorm:"column:required;type:boolean"`
	Description string `gorm:"column:description;type:text"`
}

func (A2AExtension) TableName() string { return "a2a_extensions" }

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

//-------------Feedback-------------

type Feedback struct {
	UID                uint64  `gorm:"column:uid;type:bigint;primaryKey"`
	AgentUID           uint64  `gorm:"column:agent_uid"`
	ChainID            string  `gorm:"column:chain_id"`
	AgentID            string  `gorm:"column:agent_id"`
	IdentityRegistry   string  `gorm:"column:identity_registry"`
	ReputationRegistry string  `gorm:"column:reputation_registry"`
	ClientAddress      string  `gorm:"column:client_address"`
	FeedbackIndex      uint64  `gorm:"column:feedback_index"`
	FormatValue        float64 `gorm:"column:format_value"`
	Value              string  `gorm:"column:value"`
	ValueDecimals      uint    `gorm:"column:value_decimals"`
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
	Endpoint    string `gorm:"column:endpoint;type:text"`
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

type FeedbackResp struct {
	Feedback
	Name     string
	Avatar   string
	Passport bool
}

// -------------Validation-------------
type Validation struct {
	AgentUID           uint64 `gorm:"column:agent_uid;type:bigint;not null"`
	ChainID            string `gorm:"column:chain_id;type:varchar(255);not null"`
	AgentID            string `gorm:"column:agent_id;type:varchar(255);not null"`
	ValidationRegistry string `gorm:"column:validation_registry;type:char(42);not null"`
	ValidatorAddress   string `gorm:"column:validator_address;type:char(42);not null"`
	RequestHash        string `gorm:"column:request_hash;type:char(66);not null"`
	ResponseURI        string `gorm:"column:response_uri;type:varchar(255);not null"`
	ResponseHash       string `gorm:"column:response_hash;type:char(66);not null"`
	Response           int    `gorm:"column:response;type:int;not null"`
	Tag1               string `gorm:"column:tag1;type:varchar(255);not null"`
	BlockNumber        uint64 `gorm:"column:block_number;type:bigint;not null"`
	Index              uint64 `gorm:"column:index;type:bigint;not null"`
	Timestamps         uint64 `gorm:"column:timestamps;type:bigint;not null"`
	RequestTxHash      string `gorm:"column:request_tx_hash;type:char(66);not null;primaryKey"`
	ResponseTxHash     string `gorm:"column:response_tx_hash;type:char(66);not null"`
	RequestURI         string `gorm:"column:request_uri;type:varchar(255);not null"`
}

func (Validation) TableName() string { return "validations" }

// Validator 验证者统计表结构体
type Validator struct {
	Address        string `gorm:"column:address;type:char(42);primaryKey"`
	PendingAmount  uint64 `gorm:"column:pending_amount;type:bigint;not null;default:0"`
	FinishedAmount uint64 `gorm:"column:finished_amount;type:bigint;not null;default:0"`
}

func (Validator) TableName() string { return "validators" }

type FeedbackTagScore struct {
	AgentUID                         uint64  `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Tag                              string  `gorm:"column:tag;type:varchar(255);primaryKey"`
	Score                            float64 `gorm:"column:score;type:numeric(36, 8)"`
	UpdatedAt                        uint64  `gorm:"column:updated_at;type:bigint"`
	FeedbackCount                    uint64  `gorm:"column:feedback_count;type:bigint"`
	UniqueFeedbackClientAddressCount uint64  `gorm:"column:unique_feedback_client_address_count;type:bigint"`
}

func (FeedbackTagScore) TableName() string { return "feedback_tag_scores" }

// ─────────────── ERC-8183 Commerce Reputation ───────────────

type CommerceAction struct {
	UID              uint64  `gorm:"column:uid;type:bigint;primaryKey"`
	ChainID          string  `gorm:"column:chain_id;type:varchar(255);not null"`
	CommerceContract string  `gorm:"column:commerce_contract;type:varchar(255);not null"`
	JobID            uint64  `gorm:"column:job_id;type:bigint;not null"`
	AgentUID         uint64  `gorm:"column:agent_uid;type:bigint;not null"`
	AgentAddress     string  `gorm:"column:agent_address;type:varchar(255);not null"`
	Role             string  `gorm:"column:role;type:varchar(32);not null"`
	Action           string  `gorm:"column:action;type:varchar(64);not null"`
	SignalPolarity   string  `gorm:"column:signal_polarity;type:varchar(16);not null"`
	SignalWeight     float64 `gorm:"column:signal_weight;type:numeric(4,2);not null;default:0"`
	SignalCertainty  string  `gorm:"column:signal_certainty;type:varchar(16);not null"`
	JobBudget        float64 `gorm:"column:job_budget;type:numeric(36,8);default:0"`
	Counterparty     string  `gorm:"column:counterparty;type:varchar(255);default:''"`
	Reason           string  `gorm:"column:reason;type:varchar(255);default:''"`
	Deliverable      string  `gorm:"column:deliverable;type:varchar(255);default:''"`
	PreviousStatus   string  `gorm:"column:previous_status;type:varchar(32);default:''"`
	HookAddress      string  `gorm:"column:hook_address;type:varchar(255);default:''"`
	BlockNumber      uint64  `gorm:"column:block_number;type:bigint;not null"`
	TxHash           string  `gorm:"column:tx_hash;type:varchar(255);not null"`
	LogIndex         uint    `gorm:"column:log_index;type:integer;not null"`
	BlockTimestamp   uint64  `gorm:"column:block_timestamp;type:bigint;not null"`
}

func (CommerceAction) TableName() string { return "commerce_actions" }

type CommerceScore struct {
	AgentUID                 uint64  `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Role                     string  `gorm:"column:role;type:varchar(32);primaryKey"`
	ChainID                  string  `gorm:"column:chain_id;type:varchar(255);primaryKey"`
	CommerceContract         string  `gorm:"column:commerce_contract;type:varchar(255);primaryKey"`
	CompletedCount           int     `gorm:"column:completed_count;default:0"`
	RejectedCount            int     `gorm:"column:rejected_count;default:0"`
	ExpiredResponsibleCount  int     `gorm:"column:expired_responsible_count;default:0"`
	SuccessRate              float64 `gorm:"column:success_rate;type:numeric(6,4);default:0"`
	WeightedScore            float64 `gorm:"column:weighted_score;type:numeric(6,4);default:0"`
	WeightedVolumeSum        float64 `gorm:"column:weighted_volume_sum;type:numeric(36,8);default:0"`
	CreatedCount             int     `gorm:"column:created_count;default:0"`
	FundedCount              int     `gorm:"column:funded_count;default:0"`
	FundedRate               float64 `gorm:"column:funded_rate;type:numeric(6,4);default:0"`
	CompletionRate           float64 `gorm:"column:completion_rate;type:numeric(6,4);default:0"`
	EvaluatedCount           int     `gorm:"column:evaluated_count;default:0"`
	ExpiredFromSubmittedCount int    `gorm:"column:expired_from_submitted_count;default:0"`
	Responsiveness           float64 `gorm:"column:responsiveness;type:numeric(6,4);default:0"`
	TotalJobs                int     `gorm:"column:total_jobs;default:0"`
	TotalVolume              float64 `gorm:"column:total_volume;type:numeric(36,8);default:0"`
	UniqueCounterparties     int     `gorm:"column:unique_counterparties;default:0"`
	Confidence               float64 `gorm:"column:confidence;type:numeric(4,2);default:0"`
}

func (CommerceScore) TableName() string { return "commerce_scores" }

type CommerceScoreGlobal struct {
	AgentUID                 uint64  `gorm:"column:agent_uid;type:bigint;primaryKey"`
	Role                     string  `gorm:"column:role;type:varchar(32);primaryKey"`
	CompletedCount           int     `gorm:"column:completed_count;default:0"`
	RejectedCount            int     `gorm:"column:rejected_count;default:0"`
	ExpiredResponsibleCount  int     `gorm:"column:expired_responsible_count;default:0"`
	SuccessRate              float64 `gorm:"column:success_rate;type:numeric(6,4);default:0"`
	WeightedScore            float64 `gorm:"column:weighted_score;type:numeric(6,4);default:0"`
	WeightedVolumeSum        float64 `gorm:"column:weighted_volume_sum;type:numeric(36,8);default:0"`
	CreatedCount             int     `gorm:"column:created_count;default:0"`
	FundedCount              int     `gorm:"column:funded_count;default:0"`
	FundedRate               float64 `gorm:"column:funded_rate;type:numeric(6,4);default:0"`
	CompletionRate           float64 `gorm:"column:completion_rate;type:numeric(6,4);default:0"`
	EvaluatedCount           int     `gorm:"column:evaluated_count;default:0"`
	ExpiredFromSubmittedCount int    `gorm:"column:expired_from_submitted_count;default:0"`
	Responsiveness           float64 `gorm:"column:responsiveness;type:numeric(6,4);default:0"`
	TotalJobs                int     `gorm:"column:total_jobs;default:0"`
	TotalVolume              float64 `gorm:"column:total_volume;type:numeric(36,8);default:0"`
	UniqueCounterparties     int     `gorm:"column:unique_counterparties;default:0"`
	Confidence               float64 `gorm:"column:confidence;type:numeric(4,2);default:0"`
}

func (CommerceScoreGlobal) TableName() string { return "commerce_scores_global" }
