package model

import (
	"errors"
	"fmt"
	"strings"

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
