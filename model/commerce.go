package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

	PaymentToken  string
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

type TokenVolumeItem struct {
	TokenSymbol   string  `json:"token_symbol"`
	TotalVolume  float64 `json:"total_volume"`
	TotalVolumeUSD float64 `json:"total_volume_usd"`
	JobCount     int     `json:"job_count"`
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
	job.UpdatedAt = uint64(time.Now().Unix())
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "commerce_contract"},
			{Name: "job_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"client", "provider", "evaluator", "description",
			"budget", "paid_amount", "paid_amount_usd",
			"payment_token", "token_symbol",
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
	Role     string // "client" or "provider"

	MinBudget     *float64
	MaxBudget     *float64
	MinPaidAmount *float64
	MaxPaidAmount *float64

	StartTime *uint64
	EndTime   *uint64

	SortBy string // "updated_at" (default), "created_at", "budget"

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
	if q.Role == "client" {
		query = query.Where("client != ''")
	} else if q.Role == "provider" {
		query = query.Where("provider != ''")
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

	order := "updated_at DESC"
	switch strings.ToLower(q.SortBy) {
	case "created_at", "budget":
		order = strings.ToLower(q.SortBy) + " DESC, updated_at DESC"
	case "updated_at", "":
		order = "updated_at DESC"
	}

	var jobs []CommerceJob
	if err := query.Order(order).Offset((q.Page-1)*q.PageSize).Limit(q.PageSize).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
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
