package model

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

type fcTestConfig struct {
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
}

// initFeedbackCreditTest sets up the DB connection and validates the v2 schema
// is in place. Skips when RUN_DB_TESTS != "1".
func initFeedbackCreditTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run DB-backed tests")
	}
	if db != nil {
		return // already initialized by a previous test
	}
	data, err := os.ReadFile("../server/api/logic/config.yaml")
	if err != nil {
		t.Skipf("skipping: cannot read ../server/api/logic/config.yaml (%v)", err)
	}
	var cfg fcTestConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}
	if cfg.Dns == "" {
		t.Fatal("empty dns in ../server/api/logic/config.yaml")
	}
	InitDB(cfg.Dns, cfg.OpenaiAPIKey)
}

// nukeAgent removes all feedback + scored data for a test agent uid so each test
// starts from a clean slate without colliding with real data.
func nukeAgent(t *testing.T, agentUID uint64) {
	t.Helper()
	if err := db.Where("agent_uid = ?", agentUID).Delete(&Feedback{}).Error; err != nil {
		t.Fatalf("cleanup feedbacks: %v", err)
	}
	if err := db.Exec("DELETE FROM feedback_tag_scores_v2 WHERE agent_uid = ?", agentUID).Error; err != nil {
		t.Fatalf("cleanup feedback_tag_scores_v2: %v", err)
	}
	if err := db.Exec("DELETE FROM feedback_credit_scores WHERE agent_uid = ?", agentUID).Error; err != nil {
		t.Fatalf("cleanup feedback_credit_scores: %v", err)
	}
}

// insertFeedback is a test helper that constructs a minimal valid Feedback
// row. It returns the row so callers can read back its UID.
func insertFeedback(t *testing.T, agentUID uint64, tag1, client string, value float64, ts uint64) *Feedback {
	t.Helper()
	fb := &Feedback{
		ChainID:            "TEST",
		AgentUID:           agentUID,
		AgentID:            "0",
		IdentityRegistry:   "0xIDREG",
		ReputationRegistry: "0xREPREG",
		ClientAddress:      client,
		FeedbackIndex:      uint64(ts), // unique per row in tests
		FormatValue:        value,
		Value:              "1",
		ValueDecimals:      0,
		Tag1:               tag1,
		Timestamps:         ts,
		BlockNumber:        ts,
		Index:              0,
		TxHash:             "0xTESTTX",
	}
	if err := CreateFeedback(fb); err != nil {
		t.Fatalf("insertFeedback: %v", err)
	}
	return fb
}

// ─── Trigger: INSERT path ──────────────────────────────────────────────────

func TestFeedbackCredit_TriggerInsert_BasicAggregation(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000001)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	insertFeedback(t, aid, "accuracy", "0xCLIENT_A", 0.8, 1_700_000_001)
	insertFeedback(t, aid, "accuracy", "0xCLIENT_B", 0.5, 1_700_000_002)
	insertFeedback(t, aid, "accuracy", "0xCLIENT_A", 0.9, 1_700_000_003) // dup reviewer

	rows, err := GetFeedbackTagScoresV2(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 tag row, got %d", len(rows))
	}
	r := rows[0]
	if r.Tag != "accuracy" {
		t.Errorf("tag=%q want accuracy", r.Tag)
	}
	if r.FeedbackCount != 3 {
		t.Errorf("feedback_count=%d want 3", r.FeedbackCount)
	}
	if r.ActiveCount != 3 {
		t.Errorf("active_count=%d want 3", r.ActiveCount)
	}
	if r.RevokedCount != 0 {
		t.Errorf("revoked_count=%d want 0", r.RevokedCount)
	}
	if r.UniqueReviewerCount != 2 {
		t.Errorf("unique_reviewer_count=%d want 2 (CLIENT_A,CLIENT_B)", r.UniqueReviewerCount)
	}
	if r.FirstActiveTs != 1_700_000_001 {
		t.Errorf("first_active_ts=%d want 1_700_000_001", r.FirstActiveTs)
	}
	if r.LastActiveTs != 1_700_000_003 {
		t.Errorf("last_active_ts=%d want 1_700_000_003", r.LastActiveTs)
	}
}

func TestFeedbackCredit_TriggerInsert_ValueDistribution(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000002)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	values := []float64{-1.0, -0.5, 0.0, 0.5, 1.0}
	for i, v := range values {
		insertFeedback(t, aid, "rating", "0xC"+string(rune('A'+i)), v, uint64(1_700_000_100+i))
	}

	rows, err := GetFeedbackTagScoresV2(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	r := rows[0]
	if r.ValueMin == nil || *r.ValueMin != -1.0 {
		t.Errorf("value_min=%v want -1.0", r.ValueMin)
	}
	if r.ValueMax == nil || *r.ValueMax != 1.0 {
		t.Errorf("value_max=%v want 1.0", r.ValueMax)
	}
	if r.ValueP50 == nil || *r.ValueP50 != 0.0 {
		t.Errorf("value_p50=%v want 0.0", r.ValueP50)
	}
}

func TestFeedbackCredit_TriggerInsert_NoTag_NoRow(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000003)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	// Empty tag1 — trigger must skip.
	insertFeedback(t, aid, "", "0xC1", 0.5, 1_700_000_200)

	rows, err := GetFeedbackTagScoresV2(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty tag, got %d", len(rows))
	}
}

func TestFeedbackCredit_TriggerInsert_MultipleTags(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000004)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	insertFeedback(t, aid, "accuracy", "0xC1", 0.8, 1_700_000_300)
	insertFeedback(t, aid, "latency", "0xC1", 0.6, 1_700_000_301)
	insertFeedback(t, aid, "latency", "0xC2", 0.7, 1_700_000_302)

	rows, err := GetFeedbackTagScoresV2(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 tag rows, got %d", len(rows))
	}
	tags := map[string]FeedbackTagScoreV2{}
	for _, r := range rows {
		tags[r.Tag] = r
	}
	if tags["accuracy"].FeedbackCount != 1 || tags["accuracy"].UniqueReviewerCount != 1 {
		t.Errorf("accuracy: got count=%d uniq=%d", tags["accuracy"].FeedbackCount, tags["accuracy"].UniqueReviewerCount)
	}
	if tags["latency"].FeedbackCount != 2 || tags["latency"].UniqueReviewerCount != 2 {
		t.Errorf("latency: got count=%d uniq=%d", tags["latency"].FeedbackCount, tags["latency"].UniqueReviewerCount)
	}
}

// ─── Trigger: REVOKE path ──────────────────────────────────────────────────

func TestFeedbackCredit_TriggerRevoke_DecrementsActive(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000005)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	insertFeedback(t, aid, "quality", "0xC1", 0.9, 1_700_000_400)
	insertFeedback(t, aid, "quality", "0xC2", 0.8, 1_700_000_401)
	insertFeedback(t, aid, "quality", "0xC3", 0.7, 1_700_000_402)

	// Revoke one
	if err := db.Model(&Feedback{}).
		Where("agent_uid = ? AND client_address = ?", aid, "0xC1").
		Update("revoked", true).Error; err != nil {
		t.Fatalf("revoke: %v", err)
	}

	rows, err := GetFeedbackTagScoresV2(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	r := rows[0]
	if r.FeedbackCount != 3 {
		t.Errorf("feedback_count=%d want 3 (revoked still counts)", r.FeedbackCount)
	}
	if r.ActiveCount != 2 {
		t.Errorf("active_count=%d want 2", r.ActiveCount)
	}
	if r.RevokedCount != 1 {
		t.Errorf("revoked_count=%d want 1", r.RevokedCount)
	}
	if r.UniqueReviewerCount != 2 {
		t.Errorf("unique_reviewer_count=%d want 2 (C2, C3)", r.UniqueReviewerCount)
	}
}

func TestFeedbackCredit_TriggerRevoke_AllRevoked_KeepsRowWithZeroActive(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000006)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	insertFeedback(t, aid, "x", "0xC1", 0.5, 1_700_000_500)

	if err := db.Model(&Feedback{}).
		Where("agent_uid = ?", aid).
		Update("revoked", true).Error; err != nil {
		t.Fatalf("revoke: %v", err)
	}

	rows, err := GetFeedbackTagScoresV2(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].ActiveCount != 0 || rows[0].RevokedCount != 1 {
		t.Errorf("got active=%d revoked=%d want 0/1", rows[0].ActiveCount, rows[0].RevokedCount)
	}
	// Value distribution should be NULL since no active rows
	if rows[0].ValueP50 != nil {
		t.Errorf("expected value_p50=NULL after all revoked, got %v", *rows[0].ValueP50)
	}
}

// ─── Cache table queries ───────────────────────────────────────────────────

func TestFeedbackCredit_GetFeedbackCreditScore_Missing(t *testing.T) {
	initFeedbackCreditTest(t)
	got, err := GetFeedbackCreditScore(uint64(99999999))
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing agent, got %+v", got)
	}
}

func TestFeedbackCredit_UpsertAndGetCreditScore(t *testing.T) {
	initFeedbackCreditTest(t)
	const aid = uint64(9000007)
	t.Cleanup(func() {
		db.Exec("DELETE FROM feedback_credit_scores WHERE agent_uid = ?", aid)
	})

	credit := 0.7234
	conf := 0.62
	row := &FeedbackCreditScore{
		AgentUID:          aid,
		BehavioralScore:   ptrFloat(0.81),
		SentimentScore:    ptrFloat(0.65),
		CreditScore:       &credit,
		Alpha:             ptrFloat(0.4),
		Confidence:        &conf,
		TagCount:          7,
		SentimentTagCount: 4,
		EffectiveN:        ptrFloat(18.3),
		LastComputedAt:    ptrInt64(1_778_481_371),
		LastUpdated:       ptrInt64(1_778_481_371),
	}
	if err := UpsertFeedbackCreditScore(row); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	got, err := GetFeedbackCreditScore(aid)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if got == nil {
		t.Fatal("expected row, got nil")
	}
	if got.CreditScore == nil || *got.CreditScore != credit {
		t.Errorf("credit_score=%v want %.4f", got.CreditScore, credit)
	}
	if got.TagCount != 7 {
		t.Errorf("tag_count=%d want 7", got.TagCount)
	}
}

func TestFeedbackCredit_ReviewerWeightCache(t *testing.T) {
	initFeedbackCreditTest(t)
	const addr = "0xREVIEWER_TEST_001"
	t.Cleanup(func() {
		db.Exec("DELETE FROM feedback_reviewer_weights WHERE client_address = ?", addr)
	})

	// not yet cached
	got, err := GetReviewerWeight(addr)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}

	// Upsert
	row := &FeedbackReviewerWeight{
		ClientAddress: addr,
		Weight:        0.85,
		Source:        "commerce",
		LastComputed:  1_778_481_371,
	}
	if err := UpsertReviewerWeight(row); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	got, err = GetReviewerWeight(addr)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if got == nil || got.Weight != 0.85 || got.Source != "commerce" {
		t.Errorf("got %+v want weight=0.85 source=commerce", got)
	}
}

// ─── tiny helpers ─────────────────────────────────────────────────────────

func ptrFloat(f float64) *float64 { return &f }
func ptrInt64(i int64) *int64     { return &i }
