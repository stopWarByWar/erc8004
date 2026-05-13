package logic

import (
	"agent_identity/model"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

// initOrchestratorTest is a lightweight DB initializer (no S3, no config dir
// requirements) for the orchestrator-level tests in this file. We only need
// the DB connection here; the broader initTest() in logic_test.go pulls in
// AWS + config files that may not exist in test envs.
func initOrchestratorTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run DB-backed orchestrator tests")
	}
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		t.Skipf("skipping: cannot read ./config.yaml (%v)", err)
	}
	var cfg struct {
		Dns          string `yaml:"dns"`
		OpenaiAPIKey string `yaml:"openai_api_key"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse config: %v", err)
	}
	if cfg.Dns == "" {
		t.Fatal("empty dns")
	}
	model.InitDB(cfg.Dns, cfg.OpenaiAPIKey)
}

func nukeAgent(t *testing.T, agentUID uint64) {
	t.Helper()
	if err := model.DeleteFeedbacksForTest(agentUID); err != nil {
		t.Fatalf("nuke feedbacks: %v", err)
	}
	if err := model.DeleteFeedbackCreditScoreForTest(agentUID); err != nil {
		t.Fatalf("nuke credit: %v", err)
	}
}

func feedFor(t *testing.T, aid uint64, tag, client string, value float64, ts uint64) {
	t.Helper()
	fb := &model.Feedback{
		ChainID: "TEST", AgentUID: aid, AgentID: "0",
		IdentityRegistry: "0xIDREG", ReputationRegistry: "0xREPREG",
		ClientAddress: client, FeedbackIndex: ts,
		FormatValue: value, Value: "1", ValueDecimals: 0,
		Tag1: tag, Timestamps: ts, BlockNumber: ts, TxHash: "0xTX",
	}
	if err := model.CreateFeedback(fb); err != nil {
		t.Fatalf("CreateFeedback: %v", err)
	}
}

// seedReviewerWeights pre-populates the reviewer weight cache so tests don't
// degenerate into the "all reviewers unknown, R below threshold" case. We use
// 0.8 — slightly below "saturated commerce history" (1.0) but well above the
// nil-cache floor (0.1) — to model "established reviewers".
func seedReviewerWeights(t *testing.T, addrs ...string) {
	t.Helper()
	for _, addr := range addrs {
		row := &model.FeedbackReviewerWeight{
			ClientAddress: addr,
			Weight:        0.8,
			Source:        "test",
			LastComputed:  time.Now().Unix(),
		}
		if err := model.UpsertReviewerWeight(row); err != nil {
			t.Fatalf("seed reviewer weight: %v", err)
		}
		t.Cleanup(func() {
			model.DeleteReviewerWeightForTest(addr)
		})
	}
}

// ─── Cold start ────────────────────────────────────────────────────────────

func TestGetFeedbackCredit_ColdStart_NoData(t *testing.T) {
	initOrchestratorTest(t)
	const aid = uint64(9100001)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	resp, err := GetFeedbackCredit(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCredit: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response even on cold start")
	}
	if resp.AgentUID != aid {
		t.Errorf("AgentUID=%d want %d", resp.AgentUID, aid)
	}
	if resp.CreditScore != nil {
		t.Errorf("CreditScore=%v want nil for cold start", *resp.CreditScore)
	}
	if resp.TagCount != 0 {
		t.Errorf("TagCount=%d want 0", resp.TagCount)
	}
	if len(resp.Tags) != 0 {
		t.Errorf("Tags=%d want 0", len(resp.Tags))
	}
}

// ─── Pure behavioral (no sentiment-eligible tags) ─────────────────────────

func TestGetFeedbackCredit_NonSentimentTag_AlphaIs1(t *testing.T) {
	initOrchestratorTest(t)
	const aid = uint64(9100002)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	// Latency-style values: out of [-1,1], not sentiment.
	seedReviewerWeights(t, "0xC1", "0xC2", "0xC3", "0xC4", "0xC5")
	now := uint64(time.Now().Unix())
	day := uint64(86400)
	feedFor(t, aid, "latency", "0xC1", 200, now-60*day)
	feedFor(t, aid, "latency", "0xC2", 350, now-50*day)
	feedFor(t, aid, "latency", "0xC3", 500, now-40*day)
	feedFor(t, aid, "latency", "0xC4", 800, now-30*day)
	feedFor(t, aid, "latency", "0xC5", 1200, now-10*day)

	resp, err := GetFeedbackCredit(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCredit: %v", err)
	}
	if resp.TagCount != 1 {
		t.Errorf("TagCount=%d want 1", resp.TagCount)
	}
	if resp.SentimentTagCount != 0 {
		t.Errorf("SentimentTagCount=%d want 0 (latency values are not sentiment-like)", resp.SentimentTagCount)
	}
	if resp.Alpha != 1.0 {
		t.Errorf("Alpha=%.2f want 1.0 (no sentiment → pure behavioral)", resp.Alpha)
	}
	if resp.SentimentScore != nil {
		t.Errorf("SentimentScore=%v want nil", *resp.SentimentScore)
	}
	if resp.CreditScore == nil {
		t.Fatal("expected CreditScore not nil (we have feedbacks)")
	}
	if *resp.CreditScore != resp.BehavioralScore {
		t.Errorf("CreditScore=%.4f BehavioralScore=%.4f, expected equal when α=1",
			*resp.CreditScore, resp.BehavioralScore)
	}
	if len(resp.Tags) != 1 {
		t.Fatalf("Tags=%d want 1", len(resp.Tags))
	}
	tg := resp.Tags[0]
	if tg.IsSentiment {
		t.Errorf("tag.IsSentiment=true want false")
	}
	if tg.ValueDistribution == nil {
		t.Errorf("expected ValueDistribution for non-sentiment tag")
	}
	if tg.UniqueReviewerCount != 5 {
		t.Errorf("UniqueReviewerCount=%d want 5", tg.UniqueReviewerCount)
	}
}

// ─── Sentiment-only tag ────────────────────────────────────────────────────

func TestGetFeedbackCredit_SentimentTag_AlphaIs04(t *testing.T) {
	initOrchestratorTest(t)
	const aid = uint64(9100003)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	// All positive sentiment, varied reviewers, spread across 60+ days so
	// active_days >= 30 → tag_authority > 0 → fusion engages.
	seedReviewerWeights(t, "0xC1", "0xC2", "0xC3", "0xC4", "0xC5")
	now := uint64(time.Now().Unix())
	day := uint64(86400)
	feedFor(t, aid, "accuracy", "0xC1", 0.9, now-60*day)
	feedFor(t, aid, "accuracy", "0xC2", 0.8, now-50*day)
	feedFor(t, aid, "accuracy", "0xC3", 0.7, now-40*day)
	feedFor(t, aid, "accuracy", "0xC4", 0.95, now-30*day)
	feedFor(t, aid, "accuracy", "0xC5", 0.85, now-10*day)

	resp, err := GetFeedbackCredit(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCredit: %v", err)
	}
	if resp.SentimentTagCount != 1 {
		t.Errorf("SentimentTagCount=%d want 1", resp.SentimentTagCount)
	}
	if resp.Alpha != 0.4 {
		t.Errorf("Alpha=%.2f want 0.4", resp.Alpha)
	}
	if resp.SentimentScore == nil {
		t.Fatal("SentimentScore is nil; expected sentiment in [0,1]")
	}
	if *resp.SentimentScore < 0.7 {
		t.Errorf("SentimentScore=%.4f, expected ≥ 0.7 (all positive 0.7–0.95)", *resp.SentimentScore)
	}
	if len(resp.Tags) != 1 || !resp.Tags[0].IsSentiment {
		t.Errorf("expected single sentiment tag, got %+v", resp.Tags)
	}
}

// ─── Mixed: one sentiment + one parameter-style tag ───────────────────────

func TestGetFeedbackCredit_MixedTags(t *testing.T) {
	initOrchestratorTest(t)
	const aid = uint64(9100004)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	seedReviewerWeights(t, "0xA1", "0xA2", "0xA3", "0xB1", "0xB2", "0xB3")
	now := uint64(time.Now().Unix())
	day := uint64(86400)
	// sentiment tag (spread to reach activeDays > 30)
	feedFor(t, aid, "rating", "0xA1", 0.8, now-60*day)
	feedFor(t, aid, "rating", "0xA2", 0.7, now-30*day)
	feedFor(t, aid, "rating", "0xA3", 0.6, now-10*day)
	// parameter tag
	feedFor(t, aid, "latency", "0xB1", 200, now-60*day)
	feedFor(t, aid, "latency", "0xB2", 500, now-30*day)
	feedFor(t, aid, "latency", "0xB3", 850, now-10*day)

	resp, err := GetFeedbackCredit(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCredit: %v", err)
	}
	if resp.TagCount != 2 {
		t.Errorf("TagCount=%d want 2", resp.TagCount)
	}
	if resp.SentimentTagCount != 1 {
		t.Errorf("SentimentTagCount=%d want 1", resp.SentimentTagCount)
	}
	// Confirm tag flags are correctly set
	byTag := map[string]FeedbackTagCreditResp{}
	for _, tg := range resp.Tags {
		byTag[tg.Tag] = tg
	}
	if !byTag["rating"].IsSentiment {
		t.Errorf("rating: IsSentiment=false want true")
	}
	if byTag["latency"].IsSentiment {
		t.Errorf("latency: IsSentiment=true want false")
	}
	if byTag["latency"].ValueDistribution == nil {
		t.Errorf("latency: expected ValueDistribution")
	}
}

// ─── Revoked exclusion ────────────────────────────────────────────────────

func TestGetFeedbackCredit_RevokedExcludedFromActive(t *testing.T) {
	initOrchestratorTest(t)
	const aid = uint64(9100005)
	nukeAgent(t, aid)
	t.Cleanup(func() { nukeAgent(t, aid) })

	seedReviewerWeights(t, "0xR1", "0xR2", "0xR3")
	now := uint64(time.Now().Unix())
	day := uint64(86400)
	feedFor(t, aid, "quality", "0xR1", 0.9, now-60*day)
	feedFor(t, aid, "quality", "0xR2", 0.8, now-30*day)
	feedFor(t, aid, "quality", "0xR3", -0.5, now-10*day)

	if err := model.RevokeFeedbackForTest(aid, "0xR1"); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	resp, err := GetFeedbackCredit(aid)
	if err != nil {
		t.Fatalf("GetFeedbackCredit: %v", err)
	}
	if len(resp.Tags) != 1 {
		t.Fatalf("Tags=%d want 1", len(resp.Tags))
	}
	tg := resp.Tags[0]
	if tg.FeedbackCount != 3 {
		t.Errorf("FeedbackCount=%d want 3 (total)", tg.FeedbackCount)
	}
	if tg.UniqueReviewerCount != 2 {
		t.Errorf("UniqueReviewerCount=%d want 2 (active only)", tg.UniqueReviewerCount)
	}
}
