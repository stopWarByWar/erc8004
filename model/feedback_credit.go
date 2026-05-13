// Package model — feedback credit score query layer.
//
// Read-only Go interface to the three new tables created by the migration
// 202605120000_feedback_credit_score.psql. The triggers maintain
// feedback_tag_scores_v2 automatically; this file exposes Go-side
// upsert/read for the cache tables (feedback_credit_scores,
// feedback_reviewer_weights).
package model

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ─── feedback_tag_scores_v2 ────────────────────────────────────────────────

// GetFeedbackTagScoresV2 returns all raw-stat rows for the given agent,
// ordered by active_count DESC (most-evaluated tags first).
func GetFeedbackTagScoresV2(agentUID uint64) ([]FeedbackTagScoreV2, error) {
	var rows []FeedbackTagScoreV2
	err := db.Where("agent_uid = ?", agentUID).
		Order("active_count DESC, feedback_count DESC, tag ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ─── feedback_credit_scores ────────────────────────────────────────────────

// GetFeedbackCreditScore returns the cached credit score for an agent, or nil
// if no row exists (not yet computed).
func GetFeedbackCreditScore(agentUID uint64) (*FeedbackCreditScore, error) {
	var row FeedbackCreditScore
	err := db.Where("agent_uid = ?", agentUID).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// UpsertFeedbackCreditScore writes the computed credit score, updating any
// existing row by primary key (agent_uid).
func UpsertFeedbackCreditScore(row *FeedbackCreditScore) error {
	if row == nil {
		return errors.New("nil feedback credit score row")
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "agent_uid"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"behavioral_score", "sentiment_score", "credit_score",
			"alpha", "confidence",
			"tag_count", "sentiment_tag_count",
			"effective_n",
			"last_computed_at", "last_updated",
		}),
	}).Create(row).Error
}

// ─── feedback_reviewer_weights ─────────────────────────────────────────────

// GetReviewerWeight returns the cached weight for a reviewer address, or nil
// if uncached.
func GetReviewerWeight(clientAddress string) (*FeedbackReviewerWeight, error) {
	var row FeedbackReviewerWeight
	err := db.Where("client_address = ?", clientAddress).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// UpsertReviewerWeight inserts or updates the cached reviewer weight by
// client_address.
func UpsertReviewerWeight(row *FeedbackReviewerWeight) error {
	if row == nil {
		return errors.New("nil reviewer weight row")
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "client_address"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"weight", "source", "last_computed",
		}),
	}).Create(row).Error
}

// ─── agent enumeration for cron preheater ─────────────────────────────────

// ListAgentUIDsWithFeedback returns the distinct agent_uids that currently
// have at least one row in feedback_tag_scores_v2 — i.e. agents whose credit
// score is worth (re)computing. Used by the background refresh cron.
//
// Order is ascending agent_uid for deterministic batching across runs.
func ListAgentUIDsWithFeedback() ([]uint64, error) {
	var uids []uint64
	err := db.Model(&FeedbackTagScoreV2{}).
		Distinct("agent_uid").
		Order("agent_uid ASC").
		Pluck("agent_uid", &uids).Error
	if err != nil {
		return nil, err
	}
	return uids, nil
}

// ─── feedback iteration for credit computation ────────────────────────────

// GetActiveFeedbacksForCredit returns all non-revoked feedbacks for an agent,
// ordered by (tag1, timestamps ASC). Returned rows include only the fields
// needed by the credit-score orchestrator; large columns like FeedbackURI
// are still loaded by GORM but unused.
func GetActiveFeedbacksForCredit(agentUID uint64) ([]Feedback, error) {
	var rows []Feedback
	err := db.
		Where("agent_uid = ? AND revoked = ? AND tag1 IS NOT NULL AND tag1 <> ''", agentUID, false).
		Order("tag1 ASC, timestamps ASC, uid ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ─── test helpers (exported only for use by *_test.go in other packages) ──

// DeleteFeedbacksForTest removes all feedback + v2 score data for a test
// agent_uid. Test-only helper; do not call from production code.
func DeleteFeedbacksForTest(agentUID uint64) error {
	if err := db.Where("agent_uid = ?", agentUID).Delete(&Feedback{}).Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM feedback_tag_scores_v2 WHERE agent_uid = ?", agentUID).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFeedbackCreditScoreForTest clears the credit cache for an agent.
func DeleteFeedbackCreditScoreForTest(agentUID uint64) error {
	return db.Exec("DELETE FROM feedback_credit_scores WHERE agent_uid = ?", agentUID).Error
}

// RevokeFeedbackForTest marks all feedback rows for (agent_uid, client_address)
// as revoked. Test-only helper.
func RevokeFeedbackForTest(agentUID uint64, clientAddress string) error {
	return db.Model(&Feedback{}).
		Where("agent_uid = ? AND client_address = ?", agentUID, clientAddress).
		Update("revoked", true).Error
}

// DeleteReviewerWeightForTest removes a single reviewer weight cache row.
func DeleteReviewerWeightForTest(clientAddress string) error {
	return db.Exec("DELETE FROM feedback_reviewer_weights WHERE client_address = ?", clientAddress).Error
}
