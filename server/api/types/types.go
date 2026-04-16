package types

import (
	"encoding/json"
	"time"
)

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

	// CommerceScore provides global commerce score summary organized by role.
	// It is optional and may be empty when the DB has no commerce tables or no data.
	CommerceScore map[string]CommerceScore `json:"commerce_score,omitempty"`
}

// CommerceScore is the UI-facing commerce score payload for a specific role.
// It mirrors fields returned by GET agent/commerce/scores.
type CommerceScore struct {
	Role                      string  `json:"role"`
	CompletedCount            int     `json:"completed_count"`
	RejectedCount             int     `json:"rejected_count"`
	ExpiredResponsibleCount   int     `json:"expired_responsible_count"`
	SuccessRate               float64 `json:"success_rate"`
	WeightedScore             float64 `json:"weighted_score"`
	TotalVolumeUSD            float64 `json:"total_volume_usd"`
	WeightedScoreUSD          float64 `json:"weighted_score_usd"`
	CreatedCount              int     `json:"created_count"`
	FundedCount               int     `json:"funded_count"`
	FundedRate                float64 `json:"funded_rate"`
	CompletionRate            float64 `json:"completion_rate"`
	EvaluatedCount            int     `json:"evaluated_count"`
	ExpiredFromSubmittedCount int     `json:"expired_from_submitted_count"`
	Responsiveness            float64 `json:"responsiveness"`
	TotalJobs                 int     `json:"total_jobs"`
	TotalVolume               float64 `json:"total_volume"`
	UniqueCounterparties      int     `json:"unique_counterparties"`
	Confidence                float64 `json:"confidence"`
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
	JobCreatedAmount       int64
	JobCompletedAmount     int64
	ClientAmount           int64
	PaymentVolumeUSD       float64
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

// ─────────────── Commerce Stats ───────────────

// CommerceStats is the top-level stats holder for GET /agent/commerce/stats.
type CommerceStats struct {
	ActionBreakdown    ActionBreakdown    `json:"action_breakdown"`
	TimeSeries         TimeSeriesStats    `json:"time_series"`
	BudgetDistribution BudgetDistribution `json:"budget_distribution"`
}

// ActionBreakdown maps role → action counts. Empty roles are omitted.
type ActionBreakdown map[string]ActionCounts

// ActionCounts holds counts for the 6 core commerce actions.
type ActionCounts struct {
	JobCreated   int `json:"job_created"`
	JobFunded    int `json:"job_funded"`
	JobSubmitted int `json:"job_submitted"`
	JobCompleted int `json:"job_completed"`
	JobRejected  int `json:"job_rejected"`
	JobExpired   int `json:"job_expired"`
}

// TimeSeriesStats holds three time windows.
type TimeSeriesStats struct {
	Hours24 []TimeBucket `json:"24h"`
	Days7   []TimeBucket `json:"7d"`
	Days30  []TimeBucket `json:"30d"`
}

// TimeBucket represents one time-series bucket.
type TimeBucket struct {
	Bucket         int64   `json:"bucket"`
	CompletedCount int     `json:"completed_count"`
	RejectedCount  int     `json:"rejected_count"`
	SuccessRate    float64 `json:"success_rate"`
}

// BudgetDistribution maps role → contract address → stats.
type BudgetDistribution map[string]map[string]ContractBudgetStats

// ContractBudgetStats holds total count and bucketed stats for one contract.
type ContractBudgetStats struct {
	TotalCount int           `json:"total_count"`
	Buckets    BudgetBuckets `json:"buckets"`
}

// BudgetBuckets holds small/medium/large buckets.
// When TotalCount < 3, Buckets is a flat map[string]interface{}{"count": int, "max_amount": string}.
type BudgetBuckets struct {
	Small  *BudgetBucketStat `json:"small,omitempty"`
	Medium *BudgetBucketStat `json:"medium,omitempty"`
	Large  *BudgetBucketStat `json:"large,omitempty"`
	// Flat is used when sample size < 3.
	Flat any `json:"-"` // map[string]any{"count": int, "max_amount": string}
}

// MarshalJSON implements custom marshaling to handle flat bucket case.
func (b BudgetBuckets) MarshalJSON() ([]byte, error) {
	if b.Flat != nil {
		return json.Marshal(b.Flat)
	}
	type BB BudgetBuckets
	return json.Marshal(BB(b))
}

// BudgetBucketStat holds count and max_amount for one bucket.
type BudgetBucketStat struct {
	Count     int    `json:"count"`
	MaxAmount string `json:"max_amount"`
}

// ─────────────── Commerce Jobs (Job Browser) ───────────────

// CommerceJobDTO is a UI-facing job snapshot payload (commerce_jobs).
// NOTE: we currently reuse model.CommerceJob in handlers; this DTO is reserved for future decoupling.
type CommerceJobDTO struct {
	ChainID          string `json:"chain_id"`
	CommerceContract string `json:"commerce_contract"`
	JobID            uint64 `json:"job_id"`

	Status    string `json:"status"`
	UpdatedAt uint64 `json:"updated_at"`

	Client    string `json:"client"`
	Provider  string `json:"provider"`
	Evaluator string `json:"evaluator"`

	Description string `json:"description,omitempty"`
	HookAddress string `json:"hook_address,omitempty"`

	ExpiredAt   uint64 `json:"expired_at,omitempty"`
	SubmittedAt uint64 `json:"submitted_at,omitempty"`
	CompletedAt uint64 `json:"completed_at,omitempty"`

	PaymentToken    string `json:"payment_token"`
	PaymentDecimals uint   `json:"payment_decimals"`
	TokenSymbol     string `json:"token_symbol"`

	Budget        string  `json:"budget"`
	BudgetUSD     float64 `json:"budget_usd"`
	PaidAmount    string  `json:"paid_amount"`
	PaidAmountUSD float64 `json:"paid_amount_usd"`

	PlatformFeeAmount  string  `json:"platform_fee_amount"`
	PlatformFeeUSD     float64 `json:"platform_fee_usd"`
	EvaluatorFeeAmount string  `json:"evaluator_fee_amount"`
	EvaluatorFeeUSD    float64 `json:"evaluator_fee_usd"`

	LatestBlockNumber uint64 `json:"latest_block_number,omitempty"`
	LatestTxHash      string `json:"latest_tx_hash,omitempty"`
	LatestActionUID   uint64 `json:"latest_action_uid,omitempty"`
}

type CommerceJobsListResp struct {
	Jobs  []CommerceJobDTO `json:"jobs"`
	Total int64            `json:"total"`
}

type CommerceActionDTO struct {
	ChainID          string `json:"chain_id"`
	CommerceContract string `json:"commerce_contract"`
	JobID            uint64 `json:"job_id"`

	AgentUID     uint64 `json:"agent_uid"`
	AgentAddress string `json:"agent_address"`
	Role         string `json:"role"`
	Action       string `json:"action"`

	SignalPolarity  string  `json:"signal_polarity"`
	SignalWeight    float64 `json:"signal_weight"`
	SignalCertainty string  `json:"signal_certainty"`

	JobBudget string  `json:"job_budget"`
	BudgetUSD float64 `json:"budget_usd"`

	PaymentToken    string `json:"payment_token"`
	PaymentDecimals uint   `json:"payment_decimals"`
	TokenSymbol     string `json:"token_symbol"`

	Counterparty   string `json:"counterparty"`
	Reason         string `json:"reason"`
	Deliverable    string `json:"deliverable"`
	PreviousStatus string `json:"previous_status"`
	HookAddress    string `json:"hook_address"`

	BlockNumber    uint64 `json:"block_number"`
	TxHash         string `json:"tx_hash"`
	LogIndex       uint   `json:"log_index"`
	BlockTimestamp uint64 `json:"block_timestamp"`
}

type CommerceJobDetailEvidence struct {
	TimelineSource    string   `json:"timeline_source"`
	EventsTableSource string   `json:"events_table_source"`
	SettlementEvents  []string `json:"settlement_events"`
}

type CommerceJobDetailResp struct {
	Job      CommerceJobDTO            `json:"job"`
	Evidence CommerceJobDetailEvidence `json:"evidence"`
}

type CommerceJobsGeneralSummaryResp struct {
	JobsCount       int64            `json:"jobs_count"`
	PaidVolumeUSD   float64          `json:"paid_volume_usd"`
	BudgetVolumeUSD float64          `json:"budget_volume_usd"`
	OutcomeMix      map[string]int64 `json:"outcome_mix"`
	ActiveMix       map[string]int64 `json:"active_mix"`
	LastUpdated     uint64           `json:"last_updated"`
}

type CommerceJobsGeneralDistributionsResp struct {
	Token         any `json:"token"`
	ChainContract any `json:"chain_contract"`
	Fees          any `json:"fees"`
}

type CommerceJobsGeneralResp struct {
	Summary       CommerceJobsGeneralSummaryResp       `json:"summary"`
	Distributions CommerceJobsGeneralDistributionsResp `json:"distributions"`
}

type CommerceJobsChartsResp struct {
	Charts any `json:"charts"`
}

// ─────────────── Commerce Job Actions (Job Detail) ───────────────

// CommerceJobActionDTO is a UI-facing action row for Job Detail.
// It is derived from commerce_actions (event log) with minimal fields for rendering a table.
type CommerceJobActionDTO struct {
	ChainID          string `json:"chain_id"`
	CommerceContract string `json:"commerce_contract"`
	JobID            uint64 `json:"job_id"`

	// ActionType is a UI-friendly event name (e.g., JobCreated, BudgetSet).
	ActionType string `json:"action_type"`

	BlockTimestamp uint64 `json:"block_timestamp"`
	BlockNumber    uint64 `json:"block_number"`
	TxHash         string `json:"tx_hash"`
	LogIndex       uint   `json:"log_index"`

	Actor string `json:"actor,omitempty"`
	Role  string `json:"role,omitempty"` // client|provider|evaluator|platform|unknown

	PaymentToken    string `json:"payment_token,omitempty"`
	PaymentDecimals uint   `json:"payment_decimals,omitempty"`
	TokenSymbol     string `json:"token_symbol,omitempty"`

	// Amount/AmountUSD are optional. For now, we expose budget-like amounts from commerce_actions.
	Amount    *float64 `json:"amount,omitempty"`
	AmountUSD *float64 `json:"amount_usd,omitempty"`
}

type CommerceJobActionsListResp struct {
	Actions []CommerceJobActionDTO `json:"actions"`
	Total   int64                  `json:"total"`
}
