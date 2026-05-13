package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommerceLeaderboardStats struct {
	JobCreatedAmount   int64
	JobCompletedAmount int64
	ClientAmount       int64
	PaymentVolumeUSD   float64
}

// GetCommerceLeaderboardStats returns global commerce aggregates for leaderboard display.
// It uses commerce_jobs snapshot table for efficiency.
func GetCommerceLeaderboardStats() (*CommerceLeaderboardStats, error) {
	type row struct {
		CreatedJobs   int64
		CompletedJobs int64
		Clients       int64
		PaidVolumeUSD float64
	}
	var r row
	err := db.Model(&CommerceJob{}).
		Select(`
			COUNT(*) as created_jobs,
			COUNT(*) FILTER (WHERE status = ?) as completed_jobs,
			COUNT(DISTINCT client) FILTER (WHERE client <> '') as clients,
			COALESCE(SUM(paid_amount_usd), 0) as paid_volume_usd
		`, StatusCompleted).
		Scan(&r).Error
	if err != nil {
		return nil, err
	}
	return &CommerceLeaderboardStats{
		JobCreatedAmount:   r.CreatedJobs,
		JobCompletedAmount: r.CompletedJobs,
		ClientAmount:       r.Clients,
		PaymentVolumeUSD:   r.PaidVolumeUSD,
	}, nil
}

// ─────────────── Commerce Jobs Filters (distinct options) ───────────────

func GetCommerceJobsDistinctChainIDs() ([]string, error) {
	var ids []string
	if err := db.Model(&CommerceJob{}).
		Distinct("chain_id").
		Where("chain_id <> ''").
		Pluck("chain_id", &ids).Error; err != nil {
		return nil, err
	}
	ids = uniqueNonEmptyStrings(ids)
	sort.Strings(ids)
	return ids, nil
}

func GetCommerceJobsDistinctCommerceContracts() ([]string, error) {
	var contracts []string
	if err := db.Model(&CommerceJob{}).
		Distinct("commerce_contract").
		Where("commerce_contract <> ''").
		Pluck("commerce_contract", &contracts).Error; err != nil {
		return nil, err
	}
	contracts = uniqueNonEmptyStrings(contracts)
	sort.Strings(contracts)
	return contracts, nil
}

// CommerceJobPaymentTokenFacet is one distinct (chain_id, payment_token) row from commerce_jobs
// with an aggregated token_symbol for filter UI.
type CommerceJobPaymentTokenFacet struct {
	ChainID      string `gorm:"column:chain_id"`
	PaymentToken string `gorm:"column:payment_token"`
	TokenSymbol  string `gorm:"column:token_symbol"`
}

// GetCommerceJobsDistinctPaymentTokenFacets returns distinct (chain_id, payment_token) pairs
// with a deterministic non-empty symbol when any row has token_symbol set.
func GetCommerceJobsDistinctPaymentTokenFacets() ([]CommerceJobPaymentTokenFacet, error) {
	var rows []CommerceJobPaymentTokenFacet
	err := db.Model(&CommerceJob{}).
		Select("chain_id, payment_token, MAX(NULLIF(TRIM(token_symbol), '')) AS token_symbol").
		Where("payment_token <> ?", "").
		Where("chain_id <> ?", "").
		Group("chain_id, payment_token").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]CommerceJobPaymentTokenFacet, 0, len(rows))
	for _, r := range rows {
		r.ChainID = strings.TrimSpace(r.ChainID)
		r.PaymentToken = strings.TrimSpace(r.PaymentToken)
		if r.ChainID == "" || r.PaymentToken == "" {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

func uniqueNonEmptyStrings(in []string) []string {
	if len(in) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func CreateCommerceAction(action *CommerceAction) error {
	if action == nil {
		return errors.New("nil action")
	}
	// Rely on DB unique constraint (see migrations) to guarantee idempotency.
	return db.Omit("uid").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(action).Error
}

// GetLatestCommerceActionTimestampByUID returns the latest block_timestamp for an agent.
func GetLatestCommerceActionTimestampByUID(uid uint64) (uint64, error) {
	var action CommerceAction
	err := db.Where("agent_uid = ?", uid).
		Order("block_timestamp DESC").
		First(&action).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return action.BlockTimestamp, nil
}

func GetLatestCommerceAction(chainID, commerceContract string) (uint64, uint64, error) {
	var action CommerceAction
	err := db.Where("chain_id = ? AND commerce_contract = ?", chainID, commerceContract).
		Order("block_number DESC, log_index DESC").
		First(&action).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	return action.BlockNumber, uint64(action.LogIndex), nil
}

type CommerceActionsQuery struct {
	UID          uint64
	Role         string
	Actions      []string // if empty, no filter
	Certainty    string
	Polarities   []string // if empty, no filter
	ChainID      string
	Contract     string
	Counterparty string

	// HookAddress supports exact match.
	// HasHook supports "only has hook" or "only no hook" by treating empty-string as no-hook.
	HookAddress string
	HasHook     *bool

	MinBudget *float64
	MaxBudget *float64

	PaymentToken string
	TokenSymbol  string
	MinBudgetUSD *float64
	MaxBudgetUSD *float64

	StartTime *uint64
	EndTime   *uint64

	// SortBy supports: "timestamp" (default), "budget", "signal_weight"
	SortBy string

	Page     int
	PageSize int
}

func (q CommerceActionsQuery) validate() error {
	if q.UID == 0 {
		return errors.New("invalid uid")
	}
	if q.Page <= 0 || q.PageSize <= 0 {
		return errors.New("invalid page or pageSize")
	}
	if q.SortBy == "" {
		return nil
	}
	switch strings.ToLower(q.SortBy) {
	case "timestamp", "budget", "signal_weight":
		return nil
	default:
		return fmt.Errorf("invalid sort_by: %s", q.SortBy)
	}
}

func GetCommerceActionsByQuery(q CommerceActionsQuery) ([]CommerceAction, int64, error) {
	if err := q.validate(); err != nil {
		return nil, 0, err
	}

	query := db.Model(&CommerceAction{}).Where("agent_uid = ?", q.UID)
	if q.Role != "" {
		query = query.Where("role = ?", q.Role)
	}
	if len(q.Actions) > 0 {
		query = query.Where("action IN ?", q.Actions)
	}
	if q.Certainty != "" {
		query = query.Where("signal_certainty = ?", q.Certainty)
	}
	if len(q.Polarities) > 0 {
		query = query.Where("signal_polarity IN ?", q.Polarities)
	}
	if q.ChainID != "" {
		query = query.Where("chain_id = ?", q.ChainID)
	}
	if q.Contract != "" {
		query = query.Where("commerce_contract = ?", q.Contract)
	}
	if q.Counterparty != "" {
		query = query.Where("counterparty = ?", q.Counterparty)
	}
	if q.HookAddress != "" {
		query = query.Where("hook_address = ?", q.HookAddress)
	}
	if q.HasHook != nil {
		if *q.HasHook {
			query = query.Where("hook_address <> ''")
		} else {
			query = query.Where("hook_address = ''")
		}
	}
	if q.MinBudget != nil {
		query = query.Where("job_budget >= ?", *q.MinBudget)
	}
	if q.MaxBudget != nil {
		query = query.Where("job_budget <= ?", *q.MaxBudget)
	}
	if q.PaymentToken != "" {
		query = query.Where("payment_token = ?", q.PaymentToken)
	}
	if q.TokenSymbol != "" {
		query = query.Where("token_symbol = ?", q.TokenSymbol)
	}
	if q.MinBudgetUSD != nil {
		query = query.Where("(budget_usd IS NULL OR budget_usd >= ?)", *q.MinBudgetUSD)
	}
	if q.MaxBudgetUSD != nil {
		query = query.Where("(budget_usd IS NULL OR budget_usd <= ?)", *q.MaxBudgetUSD)
	}
	if q.StartTime != nil {
		query = query.Where("block_timestamp >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		query = query.Where("block_timestamp <= ?", *q.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []CommerceAction{}, 0, nil
	}

	order := "block_timestamp DESC"
	switch strings.ToLower(q.SortBy) {
	case "", "timestamp":
		order = "block_timestamp DESC"
	case "budget":
		order = "job_budget DESC, block_timestamp DESC"
	case "signal_weight":
		order = "signal_weight DESC, block_timestamp DESC"
	}

	var actions []CommerceAction
	if err := query.
		Order(order).
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&actions).Error; err != nil {
		return nil, 0, err
	}
	return actions, total, nil
}

// Backward compatible wrapper (existing call sites).
func GetCommerceActionsByAgentUID(uid uint64, role, action, certainty string, page, pageSize int) ([]CommerceAction, int64, error) {
	var actions []string
	if action != "" {
		actions = []string{action}
	}
	return GetCommerceActionsByQuery(CommerceActionsQuery{
		UID:       uid,
		Role:      role,
		Actions:   actions,
		Certainty: certainty,
		Page:      page,
		PageSize:  pageSize,
	})
}

func GetCommerceActionsByJobID(chainID, contract string, jobID uint64) ([]CommerceAction, error) {
	var actions []CommerceAction
	err := db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).
		Order("block_number ASC, log_index ASC").
		Find(&actions).Error
	return actions, err
}

func GetSettlementActionsByJobID(chainID, contract string, jobID uint64) ([]CommerceAction, error) {
	var actions []CommerceAction
	err := db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ? AND action IN ?",
		chainID, contract, jobID,
		[]string{ActionPaymentReleased, ActionPlatformFeePaid, ActionEvaluatorFeePaid}).
		Order("block_number ASC, log_index ASC").
		Find(&actions).Error
	return actions, err
}

// ─────────────── Commerce Job Actions (Job Detail) ───────────────

type CommerceJobActionsQuery struct {
	ChainID  string
	Contract string
	JobID    uint64

	// ActionTypes filters commerce_actions.action by a list of raw action strings (snake_case).
	// If empty, no filter.
	ActionTypes []string

	StartTime *uint64 // block_timestamp >= start
	EndTime   *uint64 // block_timestamp <= end

	Page     int
	PageSize int
}

func (q CommerceJobActionsQuery) validate() error {
	if q.ChainID == "" || q.Contract == "" || q.JobID == 0 {
		return errors.New("missing chain_id/commerce_contract/job_id")
	}
	if q.Page <= 0 || q.PageSize <= 0 {
		return errors.New("invalid page or pageSize")
	}
	return nil
}

func GetCommerceJobActions(q CommerceJobActionsQuery) ([]CommerceAction, int64, error) {
	if err := q.validate(); err != nil {
		return nil, 0, err
	}

	query := db.Model(&CommerceAction{}).
		Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", q.ChainID, q.Contract, q.JobID)

	if len(q.ActionTypes) > 0 {
		query = query.Where("action IN ?", q.ActionTypes)
	}
	if q.StartTime != nil {
		query = query.Where("block_timestamp >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		query = query.Where("block_timestamp <= ?", *q.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []CommerceAction{}, 0, nil
	}

	var actions []CommerceAction
	if err := query.
		Order("block_timestamp DESC, log_index DESC").
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&actions).Error; err != nil {
		return nil, 0, err
	}
	return actions, total, nil
}

type TokenVolumeItem struct {
	TokenSymbol    string  `json:"token_symbol"`
	TotalVolume    float64 `json:"total_volume"`
	TotalVolumeUSD float64 `json:"total_volume_usd"`
	JobCount       int     `json:"job_count"`
}

func GetTokenVolumeBreakdown(uid uint64, role, chainID, contract string) ([]TokenVolumeItem, error) {
	q := db.Model(&CommerceAction{}).Where("agent_uid = ? AND role = ?", uid, role)
	if chainID != "" {
		q = q.Where("chain_id = ?", chainID)
	}
	if contract != "" {
		q = q.Where("commerce_contract = ?", contract)
	}

	rows, err := q.Select("token_symbol, SUM(job_budget) as total_volume, SUM(budget_usd) as total_volume_usd, COUNT(*) as job_count").
		Where("token_symbol != ''").
		Group("token_symbol").
		Order("total_volume_usd DESC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TokenVolumeItem
	for rows.Next() {
		var item TokenVolumeItem
		if err := rows.Scan(&item.TokenSymbol, &item.TotalVolume, &item.TotalVolumeUSD, &item.JobCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func GetCommerceScoresByAgentUID(uid uint64) ([]CommerceScore, error) {
	var scores []CommerceScore
	err := db.Where("agent_uid = ?", uid).Find(&scores).Error
	return scores, err
}

func GetCommerceScoreGlobal(uid uint64) ([]CommerceScoreGlobal, error) {
	var scores []CommerceScoreGlobal
	err := db.Where("agent_uid = ?", uid).Find(&scores).Error
	return scores, err
}

// GetAgentReputationStats returns the average feedback score and total feedback count for an agent.
func GetAgentReputationStats(uid uint64) (avgScore float64, count int, err error) {
	type row struct {
		AvgScore float64
		Count    int
	}
	var r row
	err = db.Model(&FeedbackTagScore{}).
		Select("COALESCE(AVG(score), 0) as avg_score, COUNT(*) as count").
		Where("agent_uid = ?", uid).
		Scan(&r).Error
	if err != nil {
		return 0, 0, err
	}
	return r.AvgScore, r.Count, nil
}

// GetCommerceJobsByProviderAddress returns commerce jobs where the agent is the provider.
func GetCommerceJobsByProviderAddress(providerAddress string, chainID string) ([]CommerceJob, error) {
	var jobs []CommerceJob
	err := db.Where("provider = ? AND chain_id = ?", providerAddress, chainID).Find(&jobs).Error
	return jobs, err
}

func GetCommerceScoresByAgentUIDAndContract(uid uint64, chainID, contract string) ([]CommerceScore, error) {
	var scores []CommerceScore
	err := db.Where("agent_uid = ? AND chain_id = ? AND commerce_contract = ?", uid, chainID, contract).
		Find(&scores).Error
	return scores, err
}

// EnsureCommerceActionsSchema creates the minimal commerce_actions schema if migrations
// haven't been applied yet. This is primarily used by tests and local dev setups.
func EnsureCommerceActionsSchema() error {
	if db == nil {
		return errors.New("db not initialized")
	}
	return db.Exec(`
CREATE TABLE IF NOT EXISTS commerce_actions (
  uid serial PRIMARY KEY,
  chain_id varchar(255) NOT NULL,
  commerce_contract varchar(255) NOT NULL,
  job_id bigint NOT NULL,
  agent_uid bigint NOT NULL,
  agent_address varchar(255) NOT NULL,
  role varchar(32) NOT NULL,
  action varchar(64) NOT NULL,
  signal_polarity varchar(16) NOT NULL,
  signal_weight numeric(4, 2) NOT NULL DEFAULT 0,
  signal_certainty varchar(16) NOT NULL,
  job_budget numeric(36, 8) DEFAULT 0,
  counterparty varchar(255) DEFAULT '',
  reason varchar(255) DEFAULT '',
  deliverable varchar(255) DEFAULT '',
  previous_status varchar(32) DEFAULT '',
  hook_address varchar(255) DEFAULT '',
  payment_token varchar(255) DEFAULT '',
  payment_decimals smallint DEFAULT 18,
  token_symbol varchar(32) DEFAULT '',
  budget_usd numeric(36,8) DEFAULT 0,
  block_number bigint NOT NULL,
  tx_hash varchar(255) NOT NULL,
  log_index integer NOT NULL,
  block_timestamp bigint NOT NULL,
  created_at timestamp DEFAULT NOW(),
  CONSTRAINT uniq_commerce_action UNIQUE (chain_id, commerce_contract, job_id, role, action, block_number, log_index)
);
CREATE INDEX IF NOT EXISTS idx_ca_job ON commerce_actions(chain_id, commerce_contract, job_id);
CREATE INDEX IF NOT EXISTS idx_ca_agent_timestamp ON commerce_actions(agent_uid, block_timestamp);
CREATE INDEX IF NOT EXISTS idx_ca_agent_certainty ON commerce_actions(agent_uid, signal_certainty);
CREATE INDEX IF NOT EXISTS idx_ca_hook ON commerce_actions(hook_address);
`).Error
}

// EnsureCommerceScoresSchema creates minimal commerce_scores(_global) tables if missing.
// Triggers are managed by migrations; this helper is primarily for tests.
func EnsureCommerceScoresSchema() error {
	if db == nil {
		return errors.New("db not initialized")
	}
	return db.Exec(`
CREATE TABLE IF NOT EXISTS commerce_scores (
  agent_uid bigint NOT NULL,
  role varchar(32) NOT NULL,
  chain_id varchar(255) NOT NULL,
  commerce_contract varchar(255) NOT NULL,
  completed_count integer NOT NULL DEFAULT 0,
  rejected_count integer NOT NULL DEFAULT 0,
  expired_responsible_count integer NOT NULL DEFAULT 0,
  success_rate numeric(6, 4) NOT NULL DEFAULT 0,
  weighted_score numeric(6, 4) NOT NULL DEFAULT 0,
  weighted_volume_sum numeric(36, 8) NOT NULL DEFAULT 0,
  created_count integer NOT NULL DEFAULT 0,
  funded_count integer NOT NULL DEFAULT 0,
  funded_rate numeric(6, 4) NOT NULL DEFAULT 0,
  completion_rate numeric(6, 4) NOT NULL DEFAULT 0,
  evaluated_count integer NOT NULL DEFAULT 0,
  expired_from_submitted_count integer NOT NULL DEFAULT 0,
  responsiveness numeric(6, 4) NOT NULL DEFAULT 0,
  total_jobs integer NOT NULL DEFAULT 0,
  total_volume numeric(36, 8) NOT NULL DEFAULT 0,
  unique_counterparties integer NOT NULL DEFAULT 0,
  confidence numeric(4, 2) NOT NULL DEFAULT 0,
  updated_at timestamp DEFAULT NOW(),
  PRIMARY KEY (agent_uid, role, chain_id, commerce_contract)
);
CREATE TABLE IF NOT EXISTS commerce_scores_global (
  agent_uid bigint NOT NULL,
  role varchar(32) NOT NULL,
  completed_count integer NOT NULL DEFAULT 0,
  rejected_count integer NOT NULL DEFAULT 0,
  expired_responsible_count integer NOT NULL DEFAULT 0,
  success_rate numeric(6, 4) NOT NULL DEFAULT 0,
  weighted_score numeric(6, 4) NOT NULL DEFAULT 0,
  weighted_volume_sum numeric(36, 8) NOT NULL DEFAULT 0,
  created_count integer NOT NULL DEFAULT 0,
  funded_count integer NOT NULL DEFAULT 0,
  funded_rate numeric(6, 4) NOT NULL DEFAULT 0,
  completion_rate numeric(6, 4) NOT NULL DEFAULT 0,
  evaluated_count integer NOT NULL DEFAULT 0,
  expired_from_submitted_count integer NOT NULL DEFAULT 0,
  responsiveness numeric(6, 4) NOT NULL DEFAULT 0,
  total_jobs integer NOT NULL DEFAULT 0,
  total_volume numeric(36, 8) NOT NULL DEFAULT 0,
  unique_counterparties integer NOT NULL DEFAULT 0,
  confidence numeric(4, 2) NOT NULL DEFAULT 0,
  updated_at timestamp DEFAULT NOW(),
  PRIMARY KEY (agent_uid, role)
);
`).Error
}

// EnsureCommerceJobsSchema creates a minimal commerce_jobs table if missing.
// This is primarily used by tests and local dev setups.
func EnsureCommerceJobsSchema() error {
	if db == nil {
		return errors.New("db not initialized")
	}
	return db.Exec(`
CREATE TABLE IF NOT EXISTS commerce_jobs (
  uid bigserial PRIMARY KEY,
  chain_id varchar(255) NOT NULL,
  commerce_contract varchar(255) NOT NULL,
  job_id bigint NOT NULL,

  client varchar(255) NOT NULL DEFAULT '',
  provider varchar(255) NOT NULL DEFAULT '',
  evaluator varchar(255) NOT NULL DEFAULT '',
  description text DEFAULT '',

  budget numeric(36,8) DEFAULT 0,
  budget_usd numeric(36,8) DEFAULT 0,
  paid_amount numeric(36,8) DEFAULT 0,
  paid_amount_usd numeric(36,8) DEFAULT 0,
  platform_fee_amount numeric(36,8) DEFAULT 0,
  platform_fee_usd numeric(36,8) DEFAULT 0,
  evaluator_fee_amount numeric(36,8) DEFAULT 0,
  evaluator_fee_usd numeric(36,8) DEFAULT 0,

  payment_token varchar(255) DEFAULT '',
  payment_decimals smallint DEFAULT 18,
  token_symbol varchar(32) DEFAULT '',

  status varchar(32) NOT NULL DEFAULT 'open',
  hook_address varchar(255) DEFAULT '',

  expired_at bigint DEFAULT 0,
  submitted_at bigint DEFAULT 0,
  completed_at bigint DEFAULT 0,

  latest_action_uid bigint DEFAULT 0,
  latest_block_number bigint DEFAULT 0,
  latest_tx_hash varchar(255) DEFAULT '',
  updated_at bigint NOT NULL DEFAULT 0,

  CONSTRAINT uniq_commerce_job UNIQUE (chain_id, commerce_contract, job_id)
);
CREATE INDEX IF NOT EXISTS idx_cj_chain_contract ON commerce_jobs(chain_id, commerce_contract);
CREATE INDEX IF NOT EXISTS idx_cj_status ON commerce_jobs(status);
CREATE INDEX IF NOT EXISTS idx_cj_updated ON commerce_jobs(updated_at);
`).Error
}

// ─────────────── Commerce Stats ───────────────

// TimeBucketAgg represents one time-series bucket aggregation.
type TimeBucketAgg struct {
	BucketTime     uint64
	CompletedCount int
	RejectedCount  int
	SuccessRate    float64
}

// GetCommerceActionCountsByRole returns nested map: role → action → count.
func GetCommerceActionCountsByRole(uid uint64) (map[string]map[string]int, error) {
	type row struct {
		Role   string
		Action string
		Count  int
	}
	var rows []row
	err := db.Model(&CommerceAction{}).
		Select("role, action, COUNT(*) as count").
		Where("agent_uid = ?", uid).
		Group("role, action").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]int)
	for _, r := range rows {
		if result[r.Role] == nil {
			result[r.Role] = make(map[string]int)
		}
		result[r.Role][r.Action] = r.Count
	}
	return result, nil
}

// GetCommerceTimeSeries returns bucketed time-series data.
func GetCommerceTimeSeries(uid uint64, windowStart, windowEnd, bucketDuration uint64) ([]TimeBucketAgg, error) {
	if windowEnd <= windowStart {
		return []TimeBucketAgg{}, nil
	}
	var results []TimeBucketAgg
	// bucket_key = (block_timestamp / bucketDuration) * bucketDuration
	err := db.Model(&CommerceAction{}).
		Select(`(block_timestamp / ?) * ? as bucket_time,
			SUM(CASE WHEN action = 'job_completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN action = 'job_rejected' THEN 1 ELSE 0 END) as rejected_count`,
			bucketDuration, bucketDuration).
		Where("agent_uid = ? AND block_timestamp >= ? AND block_timestamp < ?", uid, windowStart, windowEnd).
		Group("bucket_time").
		Order("bucket_time ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	// Compute success_rate for each bucket
	for i := range results {
		total := results[i].CompletedCount + results[i].RejectedCount
		if total > 0 {
			results[i].SuccessRate = float64(results[i].CompletedCount) / float64(total)
		}
	}
	return results, nil
}

// RawBudgetItem represents a single terminal action budget row for stats computation.
type RawBudgetItem struct {
	Role     string
	Contract string
	Budget   float64
}

// GetCommerceBudgetItems returns all terminal action budgets for a uid.
// Logic layer uses this to compute percentile-based bucketing.
func GetCommerceBudgetItems(uid uint64) ([]RawBudgetItem, error) {
	var items []RawBudgetItem
	err := db.Model(&CommerceAction{}).
		Select("role, commerce_contract as contract, job_budget as budget").
		Where("agent_uid = ? AND action IN ('job_completed', 'job_rejected')", uid).
		Order("role ASC, commerce_contract ASC, job_budget ASC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// ─────────────── Commerce Jobs ───────────────

// UpsertCommerceJob inserts or updates a CommerceJob record.
func UpsertCommerceJob(job *CommerceJob) error {
	if job == nil {
		return errors.New("nil job")
	}
	// UpdatedAt must be an on-chain timestamp (block_timestamp).
	// Do NOT fall back to wall time, otherwise analytics becomes inconsistent.
	if job.UpdatedAt == 0 {
		return errors.New("commerce_job.updated_at must be set to on-chain block_timestamp")
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "commerce_contract"},
			{Name: "job_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"client", "provider", "evaluator", "description",
			"budget", "budget_usd",
			"paid_amount", "paid_amount_usd",
			"platform_fee_amount", "platform_fee_usd",
			"evaluator_fee_amount", "evaluator_fee_usd",
			"payment_token", "payment_decimals", "token_symbol",
			"status", "hook_address",
			"expired_at", "submitted_at", "completed_at",
			"latest_action_uid", "latest_block_number", "latest_tx_hash", "updated_at",
		}),
	}).Create(job).Error
}

// FetchOrCreateCommerceJob finds an existing job or creates a new one.
func FetchOrCreateCommerceJob(chainID, contract string, jobID uint64) (*CommerceJob, error) {
	var job CommerceJob
	err := db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).First(&job)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return &CommerceJob{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            jobID,
			Status:           StatusOpen,
		}, nil
	}
	if err.Error != nil {
		return nil, fmt.Errorf("fetch commerce job: %w", err.Error)
	}
	return &job, nil
}

// CommerceJobsQuery defines filter parameters for job queries.
type CommerceJobsQuery struct {
	ChainID  string
	Contract string
	Status   string
	Role     string // "client" | "provider" | "evaluator" (optional view)
	// AgentAddress filters by the address field specified by Role.
	// When Role is empty, AgentAddress is ignored.
	AgentAddress string
	// Counterparty filters the opposite party based on Role:
	// - role=client -> counterparty matches provider
	// - role=provider -> counterparty matches client
	// - role=evaluator -> counterparty matches (client or provider) (P1; current behavior: no-op unless specified later)
	Counterparty string

	PaymentToken string
	TokenSymbol  string

	MinBudget     *float64
	MaxBudget     *float64
	MinPaidAmount *float64
	MaxPaidAmount *float64
	MinBudgetUSD  *float64
	MaxBudgetUSD  *float64

	StartTime *uint64
	EndTime   *uint64

	SortBy    string // "updated_at" (default), "created_at", "budget", "budget_usd", "paid_amount_usd"
	SortOrder string // "desc" (default) | "asc"

	Page     int
	PageSize int
}

func (q CommerceJobsQuery) validate() error {
	if q.Page <= 0 || q.PageSize <= 0 {
		return errors.New("invalid page or pageSize")
	}
	return nil
}

// GetCommerceJobs queries jobs with filters and pagination.
func GetCommerceJobs(q CommerceJobsQuery) ([]CommerceJob, int64, error) {
	if err := q.validate(); err != nil {
		return nil, 0, err
	}

	query := db.Model(&CommerceJob{})

	if q.ChainID != "" {
		query = query.Where("chain_id = ?", q.ChainID)
	}
	if q.Contract != "" {
		query = query.Where("commerce_contract = ?", q.Contract)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.Role != "" {
		switch strings.ToLower(q.Role) {
		case "client":
			if q.AgentAddress != "" {
				query = query.Where("client = ?", q.AgentAddress)
			} else {
				query = query.Where("client <> ''")
			}
			if q.Counterparty != "" {
				query = query.Where("provider = ?", q.Counterparty)
			}
		case "provider":
			if q.AgentAddress != "" {
				query = query.Where("provider = ?", q.AgentAddress)
			} else {
				query = query.Where("provider <> ''")
			}
			if q.Counterparty != "" {
				query = query.Where("client = ?", q.Counterparty)
			}
		case "evaluator":
			if q.AgentAddress != "" {
				query = query.Where("evaluator = ?", q.AgentAddress)
			} else {
				query = query.Where("evaluator <> ''")
			}
			// Counterparty for evaluator is ambiguous (client/provider). Keep for future extension.
		}
	}
	if q.PaymentToken != "" {
		query = query.Where("payment_token = ?", q.PaymentToken)
	}
	if q.TokenSymbol != "" {
		query = query.Where("token_symbol = ?", q.TokenSymbol)
	}
	if q.MinBudget != nil {
		query = query.Where("budget >= ?", *q.MinBudget)
	}
	if q.MaxBudget != nil {
		query = query.Where("budget <= ?", *q.MaxBudget)
	}
	if q.MinPaidAmount != nil {
		query = query.Where("paid_amount >= ?", *q.MinPaidAmount)
	}
	if q.MaxPaidAmount != nil {
		query = query.Where("paid_amount <= ?", *q.MaxPaidAmount)
	}
	if q.MinBudgetUSD != nil {
		query = query.Where("budget_usd >= ?", *q.MinBudgetUSD)
	}
	if q.MaxBudgetUSD != nil {
		query = query.Where("budget_usd <= ?", *q.MaxBudgetUSD)
	}
	if q.StartTime != nil {
		query = query.Where("updated_at >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		query = query.Where("updated_at <= ?", *q.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []CommerceJob{}, 0, nil
	}

	sortBy := strings.ToLower(strings.TrimSpace(q.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(q.SortOrder))
	if sortOrder == "" {
		sortOrder = "desc"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	order := "updated_at DESC, job_id DESC"
	switch sortBy {
	case "", "updated_at":
		order = "updated_at " + sortOrder + ", job_id " + sortOrder
	case "budget":
		order = "budget " + sortOrder + ", updated_at DESC, job_id DESC"
	case "budget_usd":
		order = "budget_usd " + sortOrder + ", updated_at DESC, job_id DESC"
	case "paid_amount_usd":
		order = "paid_amount_usd " + sortOrder + ", updated_at DESC, job_id DESC"
	}

	var jobs []CommerceJob
	if err := query.Order(order).Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
}

// ─────────────── Commerce Jobs General / Charts ───────────────

type CommerceJobsGeneralSummary struct {
	JobsCount       int64   `json:"jobs_count"`
	PaidVolumeUSD   float64 `json:"paid_volume_usd"`
	BudgetVolumeUSD float64 `json:"budget_volume_usd"`
	LastUpdated     uint64  `json:"last_updated"`

	OutcomeCompleted int64 `json:"-"`
	OutcomeRejected  int64 `json:"-"`
	OutcomeExpired   int64 `json:"-"`

	ActiveOpen      int64 `json:"-"`
	ActiveFunded    int64 `json:"-"`
	ActiveSubmitted int64 `json:"-"`
}

type CommerceJobsTokenDistributionItem struct {
	TokenSymbol     string         `json:"token_symbol"`
	JobsCount       int64          `json:"jobs_count"`
	BudgetVolumeUSD float64        `json:"budget_volume_usd"`
	PaidVolumeUSD   float64        `json:"paid_volume_usd"`
	Drilldown       map[string]any `json:"drilldown,omitempty"`
}

type CommerceJobsChainContractDistributionItem struct {
	ChainID          string         `json:"chain_id"`
	ChainName        string         `json:"chain_name,omitempty"`
	ChainLogo        string         `json:"chain_logo,omitempty"`
	CommerceContract string         `json:"commerce_contract"`
	JobsCount        int64          `json:"jobs_count"`
	BudgetVolumeUSD  float64        `json:"budget_volume_usd"`
	PaidVolumeUSD    float64        `json:"paid_volume_usd"`
	Drilldown        map[string]any `json:"drilldown,omitempty"`
}

type CommerceJobsChainContractItem struct {
	CommerceContract string         `json:"commerce_contract"`
	JobsCount        int64          `json:"jobs_count"`
	BudgetVolumeUSD  float64        `json:"budget_volume_usd"`
	PaidVolumeUSD    float64        `json:"paid_volume_usd"`
	Drilldown        map[string]any `json:"drilldown,omitempty"`
}

type CommerceJobsChainContractsItem struct {
	ChainID          string                          `json:"chain_id"`
	ChainName        string                          `json:"chain_name,omitempty"`
	ChainLogo        string                          `json:"chain_logo,omitempty"`
	ERC8183Contracts []CommerceJobsChainContractItem `json:"erc8183_contracts"`
}

type CommerceJobsFeesDistribution struct {
	PlatformFeeUSD  float64 `json:"platform_fee_usd"`
	EvaluatorFeeUSD float64 `json:"evaluator_fee_usd"`
}

type CommerceJobsChartsBucket struct {
	Bucket         uint64         `json:"bucket"`
	Count          int64          `json:"count,omitempty"`
	ValueUSD       float64        `json:"value_usd,omitempty"`
	PaidVolumeUSD  float64        `json:"paid_volume_usd,omitempty"`
	PlatformFeeUSD float64        `json:"platform_fee_usd,omitempty"`
	Drilldown      map[string]any `json:"drilldown"`
}

type CommerceJobsChartOutcomeItem struct {
	Status    string         `json:"status"`
	Count     int64          `json:"count"`
	Drilldown map[string]any `json:"drilldown"`
}

type CommerceJobsChartFunnelItem struct {
	Stage     string         `json:"stage"`
	Count     int64          `json:"count"`
	Drilldown map[string]any `json:"drilldown"`
}

type CommerceJobsChartFeesItem struct {
	Type      string         `json:"type"`
	Value     float64        `json:"value"`
	Drilldown map[string]any `json:"drilldown"`
}

type CommerceJobsChartHistogramItem struct {
	Range struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"range"`
	Count     int64          `json:"count"`
	Drilldown map[string]any `json:"drilldown"`
}

type CommerceJobsChartTopEventItem struct {
	Kind             string         `json:"kind"`
	ChainID          string         `json:"chain_id"`
	CommerceContract string         `json:"commerce_contract"`
	JobID            uint64         `json:"job_id"`
	ValueUSD         float64        `json:"value_usd"`
	Drilldown        map[string]any `json:"drilldown"`
}

type CommerceJobsCharts struct {
	Activity struct {
		ActiveJobsOverTime      []CommerceJobsChartsBucket `json:"active_jobs_over_time"`
		CreatedJobsOverTime     []CommerceJobsChartsBucket `json:"created_jobs_over_time"`
		UniqueClientsOverTime   []CommerceJobsChartsBucket `json:"unique_clients_over_time"`
		UniqueProvidersOverTime []CommerceJobsChartsBucket `json:"unique_providers_over_time"`
	} `json:"activity"`
	Health struct {
		OutcomeMix []CommerceJobsChartOutcomeItem `json:"outcome_mix"`
		Funnel     []CommerceJobsChartFunnelItem  `json:"funnel"`
	} `json:"health"`
	Volume struct {
		PaidVolumeUSDOverTime  []CommerceJobsChartsBucket  `json:"paid_volume_usd_over_time"`
		PlatformFeeUSDOverTime []CommerceJobsChartsBucket  `json:"platform_fee_usd_over_time"`
		FeesBreakdown          []CommerceJobsChartFeesItem `json:"fees_breakdown"`
	} `json:"volume"`
	Distribution struct {
		TokenDistribution  []CommerceJobsTokenDistributionItem `json:"token_distribution"`
		ChainContracts     []CommerceJobsChainContractsItem    `json:"chain_contracts"`
		BudgetHistogramUSD []CommerceJobsChartHistogramItem    `json:"budget_histogram_usd"`
	} `json:"distribution"`
	Evidence struct {
		TopEvents []CommerceJobsChartTopEventItem `json:"top_events"`
	} `json:"evidence"`
}

// GetCommerceJobsCharts returns aggregated chart data for Job Browser charts.
// NOTE: current implementation is jobs-snapshot based (commerce_jobs). Action-level charts can be added later.
func GetCommerceJobsCharts(q CommerceJobsQuery, windowStart, windowEnd uint64, bucketSeconds uint64) (*CommerceJobsCharts, error) {
	if bucketSeconds == 0 {
		return nil, errors.New("invalid bucket_seconds")
	}

	// Use the same filtered base query as general.
	summary, tokenDist, ccDist, fees, err := GetCommerceJobsGeneral(q, true)
	if err != nil {
		return nil, err
	}

	out := &CommerceJobsCharts{}
	out.Distribution.TokenDistribution = tokenDist
	out.Distribution.ChainContracts = groupChainContracts(ccDist)
	for i := range out.Distribution.TokenDistribution {
		sym := out.Distribution.TokenDistribution[i].TokenSymbol
		if sym != "" {
			out.Distribution.TokenDistribution[i].Drilldown = map[string]any{"token_symbol": sym}
		}
	}
	out.Volume.FeesBreakdown = []CommerceJobsChartFeesItem{
		{Type: "platform_fee_usd", Value: fees.PlatformFeeUSD, Drilldown: map[string]any{"status": StatusCompleted}},
		{Type: "evaluator_fee_usd", Value: fees.EvaluatorFeeUSD, Drilldown: map[string]any{"status": StatusCompleted}},
	}

	// Outcome mix and funnel from status counts already included in summary.
	out.Health.OutcomeMix = []CommerceJobsChartOutcomeItem{
		{Status: StatusCompleted, Count: summary.OutcomeCompleted, Drilldown: map[string]any{"status": StatusCompleted}},
		{Status: StatusRejected, Count: summary.OutcomeRejected, Drilldown: map[string]any{"status": StatusRejected}},
		{Status: StatusExpired, Count: summary.OutcomeExpired, Drilldown: map[string]any{"status": StatusExpired}},
	}
	out.Health.Funnel = []CommerceJobsChartFunnelItem{
		{Stage: StatusOpen, Count: summary.ActiveOpen, Drilldown: map[string]any{"status": StatusOpen}},
		{Stage: StatusFunded, Count: summary.ActiveFunded, Drilldown: map[string]any{"status": StatusFunded}},
		{Stage: StatusSubmitted, Count: summary.ActiveSubmitted, Drilldown: map[string]any{"status": StatusSubmitted}},
		{Stage: StatusCompleted, Count: summary.OutcomeCompleted, Drilldown: map[string]any{"status": StatusCompleted}},
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if q.ChainID != "" {
			tx = tx.Where("chain_id = ?", q.ChainID)
		}
		if q.Contract != "" {
			tx = tx.Where("commerce_contract = ?", q.Contract)
		}
		if q.Status != "" {
			tx = tx.Where("status = ?", q.Status)
		}
		if q.Role != "" {
			switch strings.ToLower(q.Role) {
			case "client":
				if q.AgentAddress != "" {
					tx = tx.Where("client = ?", q.AgentAddress)
				} else {
					tx = tx.Where("client <> ''")
				}
				if q.Counterparty != "" {
					tx = tx.Where("provider = ?", q.Counterparty)
				}
			case "provider":
				if q.AgentAddress != "" {
					tx = tx.Where("provider = ?", q.AgentAddress)
				} else {
					tx = tx.Where("provider <> ''")
				}
				if q.Counterparty != "" {
					tx = tx.Where("client = ?", q.Counterparty)
				}
			case "evaluator":
				if q.AgentAddress != "" {
					tx = tx.Where("evaluator = ?", q.AgentAddress)
				} else {
					tx = tx.Where("evaluator <> ''")
				}
				// Counterparty for evaluator is ambiguous (client/provider). Keep for future extension.
			}
		} else if q.AgentAddress != "" {
			// When role is empty, match any party if agent_address is provided.
			tx = tx.Where("(client = ? OR provider = ? OR evaluator = ?)", q.AgentAddress, q.AgentAddress, q.AgentAddress)
		}
		if q.Role == "" && q.Counterparty != "" {
			// Minimal behavior: match any party.
			tx = tx.Where("(client = ? OR provider = ? OR evaluator = ?)", q.Counterparty, q.Counterparty, q.Counterparty)
		}
		if q.PaymentToken != "" {
			tx = tx.Where("payment_token = ?", q.PaymentToken)
		}
		if q.TokenSymbol != "" {
			tx = tx.Where("token_symbol = ?", q.TokenSymbol)
		}
		if q.MinBudget != nil {
			tx = tx.Where("budget >= ?", *q.MinBudget)
		}
		if q.MaxBudget != nil {
			tx = tx.Where("budget <= ?", *q.MaxBudget)
		}
		if q.MinBudgetUSD != nil {
			tx = tx.Where("budget_usd >= ?", *q.MinBudgetUSD)
		}
		if q.MaxBudgetUSD != nil {
			tx = tx.Where("budget_usd <= ?", *q.MaxBudgetUSD)
		}
		if windowStart > 0 {
			tx = tx.Where("updated_at >= ?", windowStart)
		}
		if windowEnd > 0 {
			tx = tx.Where("updated_at <= ?", windowEnd)
		}
		return tx
	}

	type bucketRow struct {
		Bucket uint64
		Cnt    int64
		SumUSD float64
	}
	var bucketRows []bucketRow
	if err := applyFilters(db.Model(&CommerceJob{})).
		Select("(updated_at / ?) * ? as bucket, COUNT(*) as cnt, COALESCE(SUM(paid_amount_usd),0) as sum_usd", bucketSeconds, bucketSeconds).
		Group("bucket").
		Order("bucket ASC").
		Scan(&bucketRows).Error; err != nil {
		return nil, err
	}
	for _, r := range bucketRows {
		start := r.Bucket
		end := r.Bucket + bucketSeconds
		out.Activity.ActiveJobsOverTime = append(out.Activity.ActiveJobsOverTime, CommerceJobsChartsBucket{
			Bucket: start,
			Count:  r.Cnt,
			Drilldown: map[string]any{
				"start_time": start,
				"end_time":   end,
			},
		})
		out.Volume.PaidVolumeUSDOverTime = append(out.Volume.PaidVolumeUSDOverTime, CommerceJobsChartsBucket{
			Bucket:        start,
			ValueUSD:      r.SumUSD,
			PaidVolumeUSD: r.SumUSD,
			Drilldown: map[string]any{
				"status":     StatusCompleted,
				"start_time": start,
				"end_time":   end,
			},
		})
	}

	// Created jobs over time (JobCreated actions) — join commerce_jobs to reuse filters.
	type actionBucketRow struct {
		Bucket uint64
		Cnt    int64
		SumUSD float64
	}
	applyActionFilters := func(tx *gorm.DB) *gorm.DB {
		// Join on job key so we can apply the same job filters (status/token/budget/time).
		tx = tx.Joins("JOIN commerce_jobs cj ON cj.chain_id = commerce_actions.chain_id AND cj.commerce_contract = commerce_actions.commerce_contract AND cj.job_id = commerce_actions.job_id")
		if q.ChainID != "" {
			tx = tx.Where("cj.chain_id = ?", q.ChainID)
		}
		if q.Contract != "" {
			tx = tx.Where("cj.commerce_contract = ?", q.Contract)
		}
		if q.Status != "" {
			tx = tx.Where("cj.status = ?", q.Status)
		}
		if q.Role != "" {
			switch strings.ToLower(q.Role) {
			case "client":
				if q.AgentAddress != "" {
					tx = tx.Where("cj.client = ?", q.AgentAddress)
				} else {
					tx = tx.Where("cj.client <> ''")
				}
				if q.Counterparty != "" {
					tx = tx.Where("cj.provider = ?", q.Counterparty)
				}
			case "provider":
				if q.AgentAddress != "" {
					tx = tx.Where("cj.provider = ?", q.AgentAddress)
				} else {
					tx = tx.Where("cj.provider <> ''")
				}
				if q.Counterparty != "" {
					tx = tx.Where("cj.client = ?", q.Counterparty)
				}
			case "evaluator":
				if q.AgentAddress != "" {
					tx = tx.Where("cj.evaluator = ?", q.AgentAddress)
				} else {
					tx = tx.Where("cj.evaluator <> ''")
				}
				// Counterparty for evaluator is ambiguous (client/provider). Keep for future extension.
			}
		} else if q.AgentAddress != "" {
			tx = tx.Where("(cj.client = ? OR cj.provider = ? OR cj.evaluator = ?)", q.AgentAddress, q.AgentAddress, q.AgentAddress)
		}
		if q.Role == "" && q.Counterparty != "" {
			tx = tx.Where("(cj.client = ? OR cj.provider = ? OR cj.evaluator = ?)", q.Counterparty, q.Counterparty, q.Counterparty)
		}
		if q.PaymentToken != "" {
			tx = tx.Where("cj.payment_token = ?", q.PaymentToken)
		}
		if q.TokenSymbol != "" {
			tx = tx.Where("cj.token_symbol = ?", q.TokenSymbol)
		}
		if q.MinBudget != nil {
			tx = tx.Where("cj.budget >= ?", *q.MinBudget)
		}
		if q.MaxBudget != nil {
			tx = tx.Where("cj.budget <= ?", *q.MaxBudget)
		}
		if q.MinBudgetUSD != nil {
			tx = tx.Where("cj.budget_usd >= ?", *q.MinBudgetUSD)
		}
		if q.MaxBudgetUSD != nil {
			tx = tx.Where("cj.budget_usd <= ?", *q.MaxBudgetUSD)
		}
		// Window uses action time (block_timestamp) for action-based series.
		if windowStart > 0 {
			tx = tx.Where("commerce_actions.block_timestamp >= ?", windowStart)
		}
		if windowEnd > 0 {
			tx = tx.Where("commerce_actions.block_timestamp <= ?", windowEnd)
		}
		return tx
	}

	var createdRows []actionBucketRow
	if err := applyActionFilters(db.Model(&CommerceAction{})).
		Where("commerce_actions.action = ?", ActionJobCreated).
		Select("(commerce_actions.block_timestamp / ?) * ? as bucket, COUNT(*) as cnt, 0 as sum_usd", bucketSeconds, bucketSeconds).
		Group("bucket").
		Order("bucket ASC").
		Scan(&createdRows).Error; err != nil {
		return nil, err
	}
	for _, r := range createdRows {
		start := r.Bucket
		end := r.Bucket + bucketSeconds
		out.Activity.CreatedJobsOverTime = append(out.Activity.CreatedJobsOverTime, CommerceJobsChartsBucket{
			Bucket: start,
			Count:  r.Cnt,
			Drilldown: map[string]any{
				"start_time": start,
				"end_time":   end,
			},
		})
	}

	// Unique clients/providers over time (distinct addresses on commerce_jobs, bucketed by updated_at)
	type uniqBucketRow struct {
		Bucket uint64
		Cnt    int64
	}
	var uniqClients []uniqBucketRow
	if err := applyFilters(db.Model(&CommerceJob{})).
		Where("client <> ''").
		Select("(updated_at / ?) * ? as bucket, COUNT(DISTINCT client) as cnt", bucketSeconds, bucketSeconds).
		Group("bucket").
		Order("bucket ASC").
		Scan(&uniqClients).Error; err != nil {
		return nil, err
	}
	for _, r := range uniqClients {
		start := r.Bucket
		end := r.Bucket + bucketSeconds
		out.Activity.UniqueClientsOverTime = append(out.Activity.UniqueClientsOverTime, CommerceJobsChartsBucket{
			Bucket: start,
			Count:  r.Cnt,
			Drilldown: map[string]any{
				"start_time": start,
				"end_time":   end,
			},
		})
	}
	var uniqProviders []uniqBucketRow
	if err := applyFilters(db.Model(&CommerceJob{})).
		Where("provider <> ''").
		Select("(updated_at / ?) * ? as bucket, COUNT(DISTINCT provider) as cnt", bucketSeconds, bucketSeconds).
		Group("bucket").
		Order("bucket ASC").
		Scan(&uniqProviders).Error; err != nil {
		return nil, err
	}
	for _, r := range uniqProviders {
		start := r.Bucket
		end := r.Bucket + bucketSeconds
		out.Activity.UniqueProvidersOverTime = append(out.Activity.UniqueProvidersOverTime, CommerceJobsChartsBucket{
			Bucket: start,
			Count:  r.Cnt,
			Drilldown: map[string]any{
				"start_time": start,
				"end_time":   end,
			},
		})
	}

	// Platform fee USD over time — sum the per-event USD amount from actions (stored in budget_usd).
	var feeRows []actionBucketRow
	if err := applyActionFilters(db.Model(&CommerceAction{})).
		Where("commerce_actions.action = ?", "platform_fee_paid").
		Select("(commerce_actions.block_timestamp / ?) * ? as bucket, 0 as cnt, COALESCE(SUM(commerce_actions.budget_usd),0) as sum_usd", bucketSeconds, bucketSeconds).
		Group("bucket").
		Order("bucket ASC").
		Scan(&feeRows).Error; err != nil {
		return nil, err
	}
	for _, r := range feeRows {
		start := r.Bucket
		end := r.Bucket + bucketSeconds
		out.Volume.PlatformFeeUSDOverTime = append(out.Volume.PlatformFeeUSDOverTime, CommerceJobsChartsBucket{
			Bucket:         start,
			PlatformFeeUSD: r.SumUSD,
			Drilldown: map[string]any{
				"status":     StatusCompleted,
				"start_time": start,
				"end_time":   end,
			},
		})
	}

	// Budget histogram (USD) with fixed ranges (P0). Tunable later.
	ranges := [][2]float64{{0, 10}, {10, 100}, {100, 1000}, {1000, 10000}, {10000, 100000}, {100000, 1e18}}
	for _, rg := range ranges {
		minV, maxV := rg[0], rg[1]
		var cnt int64
		if err := applyFilters(db.Model(&CommerceJob{})).
			Where("budget_usd >= ? AND budget_usd < ?", minV, maxV).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt == 0 {
			continue
		}
		item := CommerceJobsChartHistogramItem{Count: cnt, Drilldown: map[string]any{"min_budget_usd": minV, "max_budget_usd": maxV}}
		item.Range.Min = minV
		item.Range.Max = maxV
		out.Distribution.BudgetHistogramUSD = append(out.Distribution.BudgetHistogramUSD, item)
	}

	// Top events (evidence-like) from commerce_jobs snapshot: max_paid, max_budget.
	var maxPaid CommerceJob
	_ = applyFilters(db.Model(&CommerceJob{})).Order("paid_amount_usd DESC").First(&maxPaid).Error
	if maxPaid.JobID != 0 && maxPaid.PaidAmountUSD > 0 {
		out.Evidence.TopEvents = append(out.Evidence.TopEvents, CommerceJobsChartTopEventItem{
			Kind:             "max_paid",
			ChainID:          maxPaid.ChainID,
			CommerceContract: maxPaid.CommerceContract,
			JobID:            maxPaid.JobID,
			ValueUSD:         maxPaid.PaidAmountUSD,
			Drilldown: map[string]any{
				"chain_id":          maxPaid.ChainID,
				"commerce_contract": maxPaid.CommerceContract,
				"job_id":            maxPaid.JobID,
			},
		})
	}
	var maxBudget CommerceJob
	_ = applyFilters(db.Model(&CommerceJob{})).Order("budget_usd DESC").First(&maxBudget).Error
	if maxBudget.JobID != 0 && maxBudget.BudgetUSD > 0 {
		out.Evidence.TopEvents = append(out.Evidence.TopEvents, CommerceJobsChartTopEventItem{
			Kind:             "max_budget",
			ChainID:          maxBudget.ChainID,
			CommerceContract: maxBudget.CommerceContract,
			JobID:            maxBudget.JobID,
			ValueUSD:         maxBudget.BudgetUSD,
			Drilldown: map[string]any{
				"chain_id":          maxBudget.ChainID,
				"commerce_contract": maxBudget.CommerceContract,
				"job_id":            maxBudget.JobID,
			},
		})
	}

	return out, nil
}

func groupChainContracts(flat []CommerceJobsChainContractDistributionItem) []CommerceJobsChainContractsItem {
	if len(flat) == 0 {
		return []CommerceJobsChainContractsItem{}
	}

	byChain := make(map[string]*CommerceJobsChainContractsItem)
	order := make([]string, 0, 8)

	for _, it := range flat {
		cid := strings.TrimSpace(it.ChainID)
		if cid == "" {
			continue
		}
		g, ok := byChain[cid]
		if !ok {
			g = &CommerceJobsChainContractsItem{
				ChainID:          cid,
				ChainName:        it.ChainName,
				ChainLogo:        it.ChainLogo,
				ERC8183Contracts: []CommerceJobsChainContractItem{},
			}
			byChain[cid] = g
			order = append(order, cid)
		} else {
			// best-effort fill
			if g.ChainName == "" && it.ChainName != "" {
				g.ChainName = it.ChainName
			}
			if g.ChainLogo == "" && it.ChainLogo != "" {
				g.ChainLogo = it.ChainLogo
			}
		}

		contract := strings.TrimSpace(it.CommerceContract)
		if contract == "" {
			continue
		}
		g.ERC8183Contracts = append(g.ERC8183Contracts, CommerceJobsChainContractItem{
			CommerceContract: contract,
			JobsCount:        it.JobsCount,
			BudgetVolumeUSD:  it.BudgetVolumeUSD,
			PaidVolumeUSD:    it.PaidVolumeUSD,
			Drilldown: map[string]any{
				"chain_id":          cid,
				"commerce_contract": contract,
			},
		})
	}

	out := make([]CommerceJobsChainContractsItem, 0, len(order))
	for _, cid := range order {
		g := byChain[cid]
		if g == nil || len(g.ERC8183Contracts) == 0 {
			continue
		}
		out = append(out, *g)
	}
	return out
}

// GetCommerceJobsGeneral aggregates summary + distributions for a given jobs query filter.
func GetCommerceJobsGeneral(q CommerceJobsQuery, includeDistributions bool) (*CommerceJobsGeneralSummary, []CommerceJobsTokenDistributionItem, []CommerceJobsChainContractDistributionItem, *CommerceJobsFeesDistribution, error) {
	// Reuse the same filter logic as GetCommerceJobs by applying it to a base query.
	// We intentionally ignore pagination for aggregates.
	q.Page = 1
	q.PageSize = 1
	if err := q.validate(); err != nil {
		return nil, nil, nil, nil, err
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		// Duplicate filter logic from GetCommerceJobs; safe for aggregation queries.
		if q.ChainID != "" {
			tx = tx.Where("chain_id = ?", q.ChainID)
		}
		if q.Contract != "" {
			tx = tx.Where("commerce_contract = ?", q.Contract)
		}
		if q.Status != "" {
			tx = tx.Where("status = ?", q.Status)
		}
		if q.Role != "" {
			switch strings.ToLower(q.Role) {
			case "client":
				if q.AgentAddress != "" {
					tx = tx.Where("client = ?", q.AgentAddress)
				} else {
					tx = tx.Where("client <> ''")
				}
				if q.Counterparty != "" {
					tx = tx.Where("provider = ?", q.Counterparty)
				}
			case "provider":
				if q.AgentAddress != "" {
					tx = tx.Where("provider = ?", q.AgentAddress)
				} else {
					tx = tx.Where("provider <> ''")
				}
				if q.Counterparty != "" {
					tx = tx.Where("client = ?", q.Counterparty)
				}
			case "evaluator":
				if q.AgentAddress != "" {
					tx = tx.Where("evaluator = ?", q.AgentAddress)
				} else {
					tx = tx.Where("evaluator <> ''")
				}
			}
		}
		if q.PaymentToken != "" {
			tx = tx.Where("payment_token = ?", q.PaymentToken)
		}
		if q.TokenSymbol != "" {
			tx = tx.Where("token_symbol = ?", q.TokenSymbol)
		}
		if q.MinBudget != nil {
			tx = tx.Where("budget >= ?", *q.MinBudget)
		}
		if q.MaxBudget != nil {
			tx = tx.Where("budget <= ?", *q.MaxBudget)
		}
		if q.MinPaidAmount != nil {
			tx = tx.Where("paid_amount >= ?", *q.MinPaidAmount)
		}
		if q.MaxPaidAmount != nil {
			tx = tx.Where("paid_amount <= ?", *q.MaxPaidAmount)
		}
		if q.MinBudgetUSD != nil {
			tx = tx.Where("budget_usd >= ?", *q.MinBudgetUSD)
		}
		if q.MaxBudgetUSD != nil {
			tx = tx.Where("budget_usd <= ?", *q.MaxBudgetUSD)
		}
		if q.StartTime != nil {
			tx = tx.Where("updated_at >= ?", *q.StartTime)
		}
		if q.EndTime != nil {
			tx = tx.Where("updated_at <= ?", *q.EndTime)
		}
		return tx
	}

	type summaryRow struct {
		JobsCount       int64
		PaidVolumeUSD   float64
		BudgetVolumeUSD float64
		LastUpdated     uint64
	}
	var sr summaryRow
	if err := applyFilters(db.Model(&CommerceJob{})).
		Select("COUNT(*) as jobs_count, COALESCE(SUM(paid_amount_usd),0) as paid_volume_usd, COALESCE(SUM(budget_usd),0) as budget_volume_usd, COALESCE(MAX(updated_at),0) as last_updated").
		Scan(&sr).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	out := &CommerceJobsGeneralSummary{
		JobsCount:       sr.JobsCount,
		PaidVolumeUSD:   sr.PaidVolumeUSD,
		BudgetVolumeUSD: sr.BudgetVolumeUSD,
		LastUpdated:     sr.LastUpdated,
	}

	// outcome_mix / active_mix counts
	type statusCount struct {
		Status string
		Cnt    int64
	}
	var rows []statusCount
	if err := applyFilters(db.Model(&CommerceJob{})).
		Select("status, COUNT(*) as cnt").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	for _, r := range rows {
		switch r.Status {
		case StatusCompleted:
			out.OutcomeCompleted = r.Cnt
		case StatusRejected:
			out.OutcomeRejected = r.Cnt
		case StatusExpired:
			out.OutcomeExpired = r.Cnt
		case StatusOpen:
			out.ActiveOpen = r.Cnt
		case StatusFunded:
			out.ActiveFunded = r.Cnt
		case StatusSubmitted:
			out.ActiveSubmitted = r.Cnt
		}
	}

	if !includeDistributions {
		return out, nil, nil, &CommerceJobsFeesDistribution{}, nil
	}

	var tokenDist []CommerceJobsTokenDistributionItem
	if err := applyFilters(db.Model(&CommerceJob{})).
		Select("token_symbol, COUNT(*) as jobs_count, COALESCE(SUM(budget_usd),0) as budget_volume_usd, COALESCE(SUM(paid_amount_usd),0) as paid_volume_usd").
		Where("token_symbol <> ''").
		Group("token_symbol").
		Order("paid_volume_usd DESC").
		Scan(&tokenDist).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var ccDist []CommerceJobsChainContractDistributionItem
	if err := applyFilters(db.Model(&CommerceJob{})).
		Select("chain_id, commerce_contract, COUNT(*) as jobs_count, COALESCE(SUM(budget_usd),0) as budget_volume_usd, COALESCE(SUM(paid_amount_usd),0) as paid_volume_usd").
		Group("chain_id, commerce_contract").
		Order("paid_volume_usd DESC").
		Scan(&ccDist).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	type feeRow struct {
		PlatformFeeUSD  float64
		EvaluatorFeeUSD float64
	}
	var fr feeRow
	if err := applyFilters(db.Model(&CommerceJob{})).
		Select("COALESCE(SUM(platform_fee_usd),0) as platform_fee_usd, COALESCE(SUM(evaluator_fee_usd),0) as evaluator_fee_usd").
		Scan(&fr).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	fees := &CommerceJobsFeesDistribution{PlatformFeeUSD: fr.PlatformFeeUSD, EvaluatorFeeUSD: fr.EvaluatorFeeUSD}

	return out, tokenDist, ccDist, fees, nil
}

// GetCommerceJobByID retrieves a single job by its chain+contract+jobID.
func GetCommerceJobByID(chainID, contract string, jobID uint64) (*CommerceJob, error) {
	var job CommerceJob
	err := db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).First(&job).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get commerce job: %w", err)
	}
	return &job, nil
}

func GetAgentUIDsByAddress(chainIDs []string, address string, limit int) ([]uint64, error) {
	var results1 []uint64

	query1 := db.Model(&CommerceJob{}).
		Select("DISTINCT agent_uid").
		Where("client = ?", address)

	if len(chainIDs) > 0 {
		query1 = query1.Where("chain_id IN ?", chainIDs)
	}

	if limit > 0 {
		query1 = query1.Limit(limit)
	}

	if err := query1.Scan(&results1).Error; err != nil {
		return nil, err
	}

	var results2 []uint64

	query2 := db.Model(&Agent{}).
		Select("DISTINCT agent_uid").
		Where("agent_wallet = ?", address)

	if len(chainIDs) > 0 {
		query2 = query2.Where("chain_id IN ?", chainIDs)
	}

	if limit > 0 {
		query2 = query2.Limit(limit)
	}
	if err := query2.Scan(&results2).Error; err != nil {
		return nil, err
	}

	//去重合并 result1，result2，返回
	results1 = append(results1, results2...)
	results1 = uniqueUint64s(results1)
	return results1, nil
}

func uniqueUint64s(in []uint64) []uint64 {
	if len(in) == 0 {
		return []uint64{}
	}
	seen := make(map[uint64]struct{}, len(in))
	out := make([]uint64, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
