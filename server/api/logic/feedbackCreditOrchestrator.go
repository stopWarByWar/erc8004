// Package logic — feedback credit score orchestrator.
//
// GetFeedbackCredit assembles per-agent credit scores from raw aggregates
// (maintained by Postgres triggers in feedback_tag_scores_v2) and active
// feedback rows. Time-dependent math (time_decay, sentiment detection,
// authority, fusion) is applied here, on demand, using the pure helpers in
// feedbackCreditLogic.go.
//
// Design: docs/designs/202605120000_feedback-credit-score.html §5–§8 + §12
package logic

import (
	"time"

	"agent_identity/model"
)

// ─── Response types ────────────────────────────────────────────────────────

// FeedbackCreditResp is the per-agent response, mapped to the JSON contract
// in the design doc §12.
type FeedbackCreditResp struct {
	AgentUID          uint64                  `json:"agent_uid"`
	CreditScore       *float64                `json:"credit_score"`
	Confidence        float64                 `json:"confidence"`
	BehavioralScore   float64                 `json:"behavioral_score"`
	SentimentScore    *float64                `json:"sentiment_score"`
	Alpha             float64                 `json:"alpha"`
	TagCount          int                     `json:"tag_count"`
	SentimentTagCount int                     `json:"sentiment_tag_count"`
	EffectiveN        float64                 `json:"effective_n"`
	LastUpdated       uint64                  `json:"last_updated"`
	// TotalFeedbacks is the cross-tag count of all (active + revoked) feedback
	// rows backing this credit. Computed live, not cached in the DB row.
	TotalFeedbacks uint64 `json:"total_feedbacks"`
	// RevokedCount is the subset of TotalFeedbacks that were revoked. Used by
	// consumers to surface "0 revoked / N total" credibility signal.
	RevokedCount uint64                  `json:"revoked_count"`
	Tags         []FeedbackTagCreditResp `json:"tags"`
}

// FeedbackTagCreditResp is the per-(agent, tag) breakdown.
type FeedbackTagCreditResp struct {
	Tag                 string             `json:"tag"`
	IsSentiment         bool               `json:"is_sentiment"`
	DetectedScale       string             `json:"detected_scale,omitempty"` // "percent" | "five_star" | "" (non-sentiment)
	SentimentScore      *float64           `json:"sentiment_score,omitempty"`
	Authority           float64            `json:"authority"`
	Confidence          float64            `json:"confidence"`
	FeedbackCount       uint64             `json:"feedback_count"`
	UniqueReviewerCount uint64             `json:"unique_reviewer_count"`
	EffectiveN          float64            `json:"effective_n"`
	// LastActiveTs is the most recent feedback timestamp for this (agent, tag),
	// exposed so UI can sort tags by recency. Unix seconds.
	LastActiveTs      uint64             `json:"last_active_ts"`
	ValueDistribution *ValueDistribution `json:"value_distribution,omitempty"`
}

// ValueDistribution is shown for non-sentiment tags so consumers can see the
// raw parameter shape (e.g., p95 latency = 850 ms).
type ValueDistribution struct {
	Min float64 `json:"min"`
	P50 float64 `json:"p50"`
	P95 float64 `json:"p95"`
	Max float64 `json:"max"`
}

// ─── Tunables ──────────────────────────────────────────────────────────────

const (
	creditHalfLifeDays      = 180.0
	creditConfidenceTarget  = 10.0 // weighted unique reviewers needed to saturate confidence
	creditMinEffectiveForOK = 1.0  // below this we treat agent as "cold" and return CreditScore=nil
)

// ─── Pluggable dependencies (overridable in tests) ────────────────────────

// We expose function variables so tests can substitute fakes without spinning
// up a DB. Defaults wire through the model package.
var (
	fetchTagStatsFn      = model.GetFeedbackTagScoresV2
	fetchActiveFeedbacks = model.GetActiveFeedbacksForCredit
	fetchReviewerWeight  = model.GetReviewerWeight
	nowUnixFn            = func() int64 { return time.Now().Unix() }
)

// ─── Main entry point ─────────────────────────────────────────────────────

// GetFeedbackCredit returns the rolled-up credit response for the given agent.
// Cold-start (no feedback rows) returns a response with CreditScore=nil and
// empty Tags so consumers can distinguish "unknown" from "zero score".
func GetFeedbackCredit(agentUID uint64) (*FeedbackCreditResp, error) {
	tagStats, err := fetchTagStatsFn(agentUID)
	if err != nil {
		return nil, err
	}
	if len(tagStats) == 0 {
		return &FeedbackCreditResp{AgentUID: agentUID, Tags: []FeedbackTagCreditResp{}}, nil
	}

	feedbacks, err := fetchActiveFeedbacks(agentUID)
	if err != nil {
		return nil, err
	}

	now := uint64(nowUnixFn())

	// Group active feedbacks by tag and resolve once.
	feedsByTag := make(map[string][]model.Feedback, len(tagStats))
	for _, fb := range feedbacks {
		feedsByTag[fb.Tag1] = append(feedsByTag[fb.Tag1], fb)
	}

	// reviewer-weight memoization per call (avoid repeated DB hits).
	reviewerCache := make(map[string]float64)
	lookupReviewer := func(addr string) float64 {
		if w, ok := reviewerCache[addr]; ok {
			return w
		}
		row, err := fetchReviewerWeight(addr)
		var w float64
		if err == nil && row != nil {
			w = row.Weight
		} else {
			w = reviewerWeight(nil) // floor (0.1) when uncached
		}
		reviewerCache[addr] = w
		return w
	}

	// Tally cross-tag aggregates for Layer A as we walk tags.
	var (
		totalFeedback           uint64
		totalRevoked            uint64
		reviewerWeightSum       = make(map[string]float64) // per-reviewer sum of w_i across tags
		tagResponses            = make([]FeedbackTagCreditResp, 0, len(tagStats))
		sentimentNumeratorByAth float64 // for Layer B fusion
		sentimentDenomByAth     float64
		sentimentTagCount       int
		latestUpdated           uint64
		totalEffectiveN         float64
	)

	for _, ts := range tagStats {
		totalFeedback += ts.FeedbackCount
		totalRevoked += ts.RevokedCount
		if ts.LastUpdated > latestUpdated {
			latestUpdated = ts.LastUpdated
		}

		active := feedsByTag[ts.Tag]
		// Per-reviewer dedupe ordinal (sorted by timestamps ASC from query)
		dedupeK := make(map[string]int)
		// Per-feedback weight and per-reviewer aggregate weight
		var (
			values             = make([]float64, 0, len(active))
			perReviewerWeight  = make(map[string]float64)
			weightedSentSum    float64
			weightSum          float64
		)
		for _, fb := range active {
			dedupeK[fb.ClientAddress]++
			k := dedupeK[fb.ClientAddress]
			rw := lookupReviewer(fb.ClientAddress)
			td := timeDecay(fb.Timestamps, now, creditHalfLifeDays)
			dd := dedupeFactor(k)
			w := rw * td * dd
			perReviewerWeight[fb.ClientAddress] += w
			weightSum += w
			values = append(values, fb.FormatValue)
		}
		// effective_n = weighted unique reviewer count
		var effectiveN float64
		for _, w := range perReviewerWeight {
			effectiveN += w
		}
		totalEffectiveN += effectiveN
		// fold into cross-tag reviewer-weight map for HHI
		for addr, w := range perReviewerWeight {
			reviewerWeightSum[addr] += w
		}

		// Sentiment detection + tag sentiment score
		isSent, scale := isSentimentEligible(values)
		var tagSent *float64
		var detectedScale string
		if scale != nil {
			detectedScale = scale.Name
		}
		if isSent && weightSum > 0 {
			// Re-walk feedbacks with weighted normalized+asymmetric values.
			// Identical weights to the first pass; dedupe ordinal recomputed
			// from scratch to keep this branch self-contained.
			dedupeRound2 := make(map[string]int, len(perReviewerWeight))
			for _, fb := range active {
				dedupeRound2[fb.ClientAddress]++
				k := dedupeRound2[fb.ClientAddress]
				rw := lookupReviewer(fb.ClientAddress)
				td := timeDecay(fb.Timestamps, now, creditHalfLifeDays)
				dd := dedupeFactor(k)
				w := rw * td * dd
				n := normalizeSentiment(fb.FormatValue, scale)
				n = applyNegativeAsymmetry(n)
				weightedSentSum += n * w
			}
			s := clamp01(weightedSentSum / weightSum)
			tagSent = &s
		}

		// authority + confidence
		activeDays := 0
		if ts.LastActiveTs > ts.FirstActiveTs {
			activeDays = int((ts.LastActiveTs - ts.FirstActiveTs) / 86400)
		}
		auth := tagAuthority(int(ts.UniqueReviewerCount), activeDays, 0 /* newcomer dominance: v1 placeholder */)
		conf := effectiveN / creditConfidenceTarget
		if conf > 1 {
			conf = 1
		}

		tr := FeedbackTagCreditResp{
			Tag:                 ts.Tag,
			IsSentiment:         isSent,
			DetectedScale:       detectedScale,
			SentimentScore:      tagSent,
			Authority:           auth,
			Confidence:          conf,
			FeedbackCount:       ts.FeedbackCount,
			UniqueReviewerCount: ts.UniqueReviewerCount,
			EffectiveN:          effectiveN,
			LastActiveTs:        ts.LastActiveTs,
		}
		if !isSent {
			tr.ValueDistribution = makeValueDistribution(ts)
		}
		tagResponses = append(tagResponses, tr)

		if isSent && tagSent != nil {
			sentimentTagCount++
			w := auth * conf
			sentimentNumeratorByAth += *tagSent * w
			sentimentDenomByAth += w
		}
	}

	// ─ Layer A ─
	var R float64
	reviewerWeights := make([]float64, 0, len(reviewerWeightSum))
	for _, w := range reviewerWeightSum {
		R += w
		reviewerWeights = append(reviewerWeights, w)
	}
	revocationRate := 0.0
	if totalFeedback > 0 {
		revocationRate = float64(totalRevoked) / float64(totalFeedback)
	}
	diversity := 1 - hhi(reviewerWeights)
	if diversity < 0 {
		diversity = 0
	}
	behavioral := behavioralScore(R, revocationRate, diversity)

	// ─ Layer B ─
	var sentimentPtr *float64
	if sentimentDenomByAth > 0 {
		b := clamp01(sentimentNumeratorByAth / sentimentDenomByAth)
		sentimentPtr = &b
	}

	// ─ Fuse ─
	hasSentiment := sentimentPtr != nil
	var bForFuse float64
	if hasSentiment {
		bForFuse = *sentimentPtr
	}
	creditScore, alpha := fuseCreditScore(behavioral, bForFuse, hasSentiment)

	// Confidence at agent-level: same saturation curve as per-tag.
	overallConf := R / creditConfidenceTarget
	if overallConf > 1 {
		overallConf = 1
	}

	resp := &FeedbackCreditResp{
		AgentUID:          agentUID,
		Confidence:        overallConf,
		BehavioralScore:   behavioral,
		SentimentScore:    sentimentPtr,
		Alpha:             alpha,
		TagCount:          len(tagStats),
		SentimentTagCount: sentimentTagCount,
		EffectiveN:        totalEffectiveN,
		LastUpdated:       latestUpdated,
		TotalFeedbacks:    totalFeedback,
		RevokedCount:      totalRevoked,
		Tags:              tagResponses,
	}
	if R >= creditMinEffectiveForOK {
		resp.CreditScore = &creditScore
	}
	return resp, nil
}

func makeValueDistribution(ts model.FeedbackTagScoreV2) *ValueDistribution {
	if ts.ValueMin == nil && ts.ValueP50 == nil && ts.ValueP95 == nil && ts.ValueMax == nil {
		return nil
	}
	d := &ValueDistribution{}
	if ts.ValueMin != nil {
		d.Min = *ts.ValueMin
	}
	if ts.ValueP50 != nil {
		d.P50 = *ts.ValueP50
	}
	if ts.ValueP95 != nil {
		d.P95 = *ts.ValueP95
	}
	if ts.ValueMax != nil {
		d.Max = *ts.ValueMax
	}
	return d
}
