// Package logic — feedback credit score background refresher.
//
// RefreshFeedbackCreditCache walks every agent that has feedback data,
// computes their credit score via GetFeedbackCredit, and persists the result
// into the feedback_credit_scores cache table. StartFeedbackCreditCron is a
// fire-and-forget starter that runs the refresher on a ticker.
//
// Design: docs/designs/202605120000_feedback-credit-score.html — Phase 4+
// extension. The orchestrator still computes scores on demand; this cron is
// purely a preheater so cache-backed consumers (leaderboards, agent bank)
// can read fresh aggregates without paying the per-agent computation cost.
package logic

import (
	"context"
	"sync/atomic"
	"time"

	"agent_identity/model"

	"github.com/sirupsen/logrus"
)

// ─── Tunables ──────────────────────────────────────────────────────────────

const (
	feedbackCreditDefaultRefreshInterval = 5 * time.Minute
)

// ─── Pluggable dependencies (overridable in tests) ────────────────────────

var (
	// Inputs to RefreshFeedbackCreditCache, overridable for unit tests.
	listAgentUIDsForRefreshFn     = model.ListAgentUIDsWithFeedback
	getFeedbackCreditForRefreshFn = GetFeedbackCredit
	upsertFeedbackCreditScoreFn   = model.UpsertFeedbackCreditScore

	// Indirection so the cron loop test can stub the work function.
	refreshFeedbackCreditCacheFn = RefreshFeedbackCreditCache

	// Single-start guard.
	feedbackCreditCronStarted atomic.Bool
)

// resetCronStartedFlag is a test-only hook so successive cron tests do not
// share state across goroutines.
func resetCronStartedFlag() {
	feedbackCreditCronStarted.Store(false)
}

// ─── RefreshFeedbackCreditCache ───────────────────────────────────────────

// RefreshFeedbackCreditCache iterates every agent with at least one feedback
// tag stat row, computes their credit response, and upserts the rolled-up
// result into feedback_credit_scores. It returns the number of agents
// successfully refreshed.
//
// Per-agent errors are logged but do not abort the loop. Context cancellation
// is honored between agents and short-circuits with ctx.Err().
func RefreshFeedbackCreditCache(ctx context.Context) (int, error) {
	uids, err := listAgentUIDsForRefreshFn()
	if err != nil {
		return 0, err
	}

	var refreshed int
	for _, uid := range uids {
		// Context check before each agent — cheap and lets callers exit early.
		if err := ctx.Err(); err != nil {
			return refreshed, err
		}

		resp, err := getFeedbackCreditForRefreshFn(uid)
		if err != nil {
			logRefreshErr("compute", uid, err)
			continue
		}
		row := buildCreditScoreRow(resp)
		if err := upsertFeedbackCreditScoreFn(row); err != nil {
			logRefreshErr("upsert", uid, err)
			continue
		}
		refreshed++
	}
	return refreshed, nil
}

// buildCreditScoreRow converts an in-memory FeedbackCreditResp to the
// persisted FeedbackCreditScore model row. Pointer fields propagate nullness
// so the JSON contract round-trips cleanly.
func buildCreditScoreRow(resp *FeedbackCreditResp) *model.FeedbackCreditScore {
	row := &model.FeedbackCreditScore{
		AgentUID:          resp.AgentUID,
		CreditScore:       resp.CreditScore,
		SentimentScore:    resp.SentimentScore,
		TagCount:          resp.TagCount,
		SentimentTagCount: resp.SentimentTagCount,
	}
	// Non-nullable fields → wrap into pointer for the table.
	beh := resp.BehavioralScore
	row.BehavioralScore = &beh
	conf := resp.Confidence
	row.Confidence = &conf
	alpha := resp.Alpha
	row.Alpha = &alpha
	effN := resp.EffectiveN
	row.EffectiveN = &effN

	nowUnix := nowUnixFn()
	row.LastComputedAt = &nowUnix
	if resp.LastUpdated > 0 {
		lu := int64(resp.LastUpdated)
		row.LastUpdated = &lu
	}
	return row
}

// ─── StartFeedbackCreditCron ──────────────────────────────────────────────

// StartFeedbackCreditCron spawns a single background goroutine that calls
// RefreshFeedbackCreditCache once immediately and then on every tick of the
// given interval, until the context is cancelled.
//
// Idempotent: subsequent calls after the first are no-ops. interval ≤ 0
// falls back to feedbackCreditDefaultRefreshInterval.
func StartFeedbackCreditCron(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = feedbackCreditDefaultRefreshInterval
	}
	if !feedbackCreditCronStarted.CompareAndSwap(false, true) {
		return
	}

	go func() {
		// Best-effort initial refresh so cache is usable shortly after boot.
		_, _ = refreshFeedbackCreditCacheFn(ctx)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = refreshFeedbackCreditCacheFn(ctx)
			}
		}
	}()
}

// ─── helpers ──────────────────────────────────────────────────────────────

func logRefreshErr(stage string, uid uint64, err error) {
	logrus.WithFields(logrus.Fields{
		"stage":     stage,
		"agent_uid": uid,
		"error":     err,
	}).Warn("feedback credit cache refresh: per-agent error; continuing")
}
