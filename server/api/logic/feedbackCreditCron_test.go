package logic

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"agent_identity/model"
)

// ─── RefreshFeedbackCreditCache (DB-backed) ───────────────────────────────

func TestRefreshFeedbackCreditCache_PopulatesCacheForAgent(t *testing.T) {
	initOrchestratorTest(t)
	const aid = uint64(9300001)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	seedReviewerWeights(t, "0xCR1", "0xCR2", "0xCR3", "0xCR4", "0xCR5")
	now := uint64(time.Now().Unix())
	day := uint64(86400)
	feedFor(t, aid, "accuracy", "0xCR1", 0.9, now-60*day)
	feedFor(t, aid, "accuracy", "0xCR2", 0.85, now-50*day)
	feedFor(t, aid, "accuracy", "0xCR3", 0.7, now-40*day)
	feedFor(t, aid, "accuracy", "0xCR4", 0.95, now-30*day)
	feedFor(t, aid, "accuracy", "0xCR5", 0.8, now-10*day)

	// Make sure no stale cache row biases the assertion.
	if err := model.DeleteFeedbackCreditScoreForTest(aid); err != nil {
		t.Fatalf("clear cache: %v", err)
	}

	refreshed, err := RefreshFeedbackCreditCache(context.Background())
	if err != nil {
		t.Fatalf("RefreshFeedbackCreditCache: %v", err)
	}
	if refreshed < 1 {
		t.Errorf("refreshed=%d want ≥ 1", refreshed)
	}

	// The cache row for our agent must now exist with values matching what
	// GetFeedbackCredit returns when called directly.
	live, err := GetFeedbackCredit(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCredit: %v", err)
	}
	cached, err := model.GetFeedbackCreditScore(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCreditScore: %v", err)
	}
	if cached == nil {
		t.Fatalf("expected cache row for aid=%d after refresh, got nil", aid)
	}
	// CreditScore should match within rounding (numeric(6,4) on DB side).
	if live.CreditScore == nil {
		t.Fatal("live CreditScore unexpectedly nil")
	}
	if cached.CreditScore == nil {
		t.Fatal("cached CreditScore unexpectedly nil")
	}
	if diff := *live.CreditScore - *cached.CreditScore; diff > 1e-3 || diff < -1e-3 {
		t.Errorf("CreditScore mismatch: live=%.6f cached=%.6f", *live.CreditScore, *cached.CreditScore)
	}
	if cached.TagCount != live.TagCount {
		t.Errorf("TagCount mismatch: cached=%d live=%d", cached.TagCount, live.TagCount)
	}
	if cached.SentimentTagCount != live.SentimentTagCount {
		t.Errorf("SentimentTagCount mismatch: cached=%d live=%d", cached.SentimentTagCount, live.SentimentTagCount)
	}
}

func TestRefreshFeedbackCreditCache_ContinuesOnPerAgentError(t *testing.T) {
	// Use injectable hook: orchestrator function returns an error for one
	// specific uid; the loop must continue and refresh the others.
	const (
		good  = uint64(9300100)
		bad   = uint64(9300101)
		good2 = uint64(9300102)
	)
	origList := listAgentUIDsForRefreshFn
	origGet := getFeedbackCreditForRefreshFn
	origUpsert := upsertFeedbackCreditScoreFn
	t.Cleanup(func() {
		listAgentUIDsForRefreshFn = origList
		getFeedbackCreditForRefreshFn = origGet
		upsertFeedbackCreditScoreFn = origUpsert
	})

	listAgentUIDsForRefreshFn = func() ([]uint64, error) {
		return []uint64{good, bad, good2}, nil
	}
	getFeedbackCreditForRefreshFn = func(uid uint64) (*FeedbackCreditResp, error) {
		if uid == bad {
			return nil, &errFake{msg: "boom"}
		}
		s := 0.5
		return &FeedbackCreditResp{
			AgentUID:    uid,
			CreditScore: &s,
			TagCount:    1,
		}, nil
	}

	var upserted []uint64
	upsertFeedbackCreditScoreFn = func(row *model.FeedbackCreditScore) error {
		upserted = append(upserted, row.AgentUID)
		return nil
	}

	refreshed, err := RefreshFeedbackCreditCache(context.Background())
	if err != nil {
		t.Fatalf("RefreshFeedbackCreditCache: %v", err)
	}
	if refreshed != 2 {
		t.Errorf("refreshed=%d want 2 (good + good2)", refreshed)
	}
	if len(upserted) != 2 || upserted[0] != good || upserted[1] != good2 {
		t.Errorf("upserted=%v want [%d %d]", upserted, good, good2)
	}
}

func TestRefreshFeedbackCreditCache_RespectsContextCancel(t *testing.T) {
	origList := listAgentUIDsForRefreshFn
	origGet := getFeedbackCreditForRefreshFn
	origUpsert := upsertFeedbackCreditScoreFn
	t.Cleanup(func() {
		listAgentUIDsForRefreshFn = origList
		getFeedbackCreditForRefreshFn = origGet
		upsertFeedbackCreditScoreFn = origUpsert
	})

	listAgentUIDsForRefreshFn = func() ([]uint64, error) {
		uids := make([]uint64, 100)
		for i := range uids {
			uids[i] = uint64(9400000 + i)
		}
		return uids, nil
	}
	var processed atomic.Int32
	getFeedbackCreditForRefreshFn = func(uid uint64) (*FeedbackCreditResp, error) {
		processed.Add(1)
		s := 0.5
		return &FeedbackCreditResp{AgentUID: uid, CreditScore: &s}, nil
	}
	upsertFeedbackCreditScoreFn = func(row *model.FeedbackCreditScore) error { return nil }

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediate cancel

	_, err := RefreshFeedbackCreditCache(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	// We should have stopped early — definitely not processed all 100.
	if processed.Load() >= 100 {
		t.Errorf("expected early stop, processed=%d/100", processed.Load())
	}
}

// ─── StartFeedbackCreditCron ──────────────────────────────────────────────

func TestStartFeedbackCreditCron_RunsAndStops(t *testing.T) {
	resetCronStartedFlag()
	t.Cleanup(resetCronStartedFlag)

	origRefresh := refreshFeedbackCreditCacheFn
	t.Cleanup(func() { refreshFeedbackCreditCacheFn = origRefresh })

	var runs atomic.Int32
	refreshFeedbackCreditCacheFn = func(ctx context.Context) (int, error) {
		runs.Add(1)
		return 0, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	StartFeedbackCreditCron(ctx, 20*time.Millisecond)

	time.Sleep(80 * time.Millisecond)
	cancel()
	time.Sleep(20 * time.Millisecond)

	got := runs.Load()
	if got < 2 {
		t.Errorf("expected ≥ 2 runs (initial + at least 1 ticker), got %d", got)
	}
	if got > 12 {
		t.Errorf("expected ≤ 12 runs (sanity bound), got %d", got)
	}
}

func TestStartFeedbackCreditCron_IdempotentStart(t *testing.T) {
	resetCronStartedFlag()
	t.Cleanup(resetCronStartedFlag)

	origRefresh := refreshFeedbackCreditCacheFn
	t.Cleanup(func() { refreshFeedbackCreditCacheFn = origRefresh })

	var runs atomic.Int32
	refreshFeedbackCreditCacheFn = func(ctx context.Context) (int, error) {
		runs.Add(1)
		return 0, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartFeedbackCreditCron(ctx, 50*time.Millisecond)
	StartFeedbackCreditCron(ctx, 50*time.Millisecond) // second call must be a no-op
	StartFeedbackCreditCron(ctx, 50*time.Millisecond) // ditto

	time.Sleep(30 * time.Millisecond) // long enough for the initial run only
	cancel()
	time.Sleep(20 * time.Millisecond)

	// Only the FIRST Start should have spawned a goroutine. Initial-run on
	// that one goroutine → exactly 1 invocation in this short window.
	if got := runs.Load(); got != 1 {
		t.Errorf("expected 1 run (idempotent), got %d", got)
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────

type errFake struct{ msg string }

func (e *errFake) Error() string { return e.msg }
