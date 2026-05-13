// Package logic — feedback credit score pure-function building blocks.
//
// Design reference: docs/designs/202605120000_feedback-credit-score.html
//
// This file holds value-agnostic, side-effect-free math used by both the
// online orchestrator (GetFeedbackCredit) and the offline trigger fallback
// (Postgres functions can mirror these formulas for sanity testing).
package logic

import (
	"math"
	"sort"

	"agent_identity/model"
)

// ─── Constants ──────────────────────────────────────────────────────────────

const (
	// reviewerWeight tuning
	reviewerWeightFloor          = 0.1  // assigned when global commerce history is nil
	reviewerWeightBaselineFloor  = 0.3  // assigned even with zero counts
	reviewerWeightSpan           = 0.7  // span above baseline floor up to 1.0
	reviewerWeightSaturatedAt    = 20.0 // completed jobs at which weight saturates to 1.0

	// negative-asymmetry tuning
	asymmetryPivot     = 0.4 // below this, normalized sentiment is amplified downward
	asymmetryAmplifier = 1.3 // amplification factor for diff below pivot

	// sentiment eligibility thresholds
	sentimentMinSamples    = 3
	sentimentMaxAbsMedian  = 1.0
	sentimentMaxAbsExtreme = 2.0
	sentimentMaxStddev     = 1.0

	// tag-authority normalization: how many unique reviewers saturate the log term
	authorityLogSaturateAt = 50.0
	authorityTimeSaturateDays = 30.0

	// sigmoid sharpness: x / (x + sigmoidK)
	sigmoidK = 3.0

	// fusion alpha when sentiment tags exist
	fusionAlphaWithSentiment = 0.4
)

// ─── Helpers ────────────────────────────────────────────────────────────────

// clamp01 clamps x to [0, 1].
func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// sigmoid is a smooth saturating map x → x/(x+K) with K=sigmoidK, returning [0,1].
// Negative inputs clamp to 0 (we never want a negative effective-N).
func sigmoid(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return x / (x + sigmoidK)
}

// ─── timeDecay ──────────────────────────────────────────────────────────────

// timeDecay returns an exponential decay coefficient in [0, 1] based on how
// far in the past blockTs is relative to nowTs. halfLifeDays controls the
// half-life; non-positive half-life falls back to 1.0 (no decay).
func timeDecay(blockTs, nowTs uint64, halfLifeDays float64) float64 {
	if halfLifeDays <= 0 {
		return 1.0
	}
	if nowTs <= blockTs {
		return 1.0
	}
	const secondsPerDay = 86400.0
	elapsedDays := float64(nowTs-blockTs) / secondsPerDay
	return math.Pow(0.5, elapsedDays/halfLifeDays)
}

// ─── dedupeFactor ───────────────────────────────────────────────────────────

// dedupeFactor returns the weight for the k-th repeat feedback from the same
// reviewer on the same (agent, tag). Invalid k (≤ 0) returns 0. k=1 → 1,
// k=2 → 0.5, k=3 → 0.25, etc.
func dedupeFactor(k int) float64 {
	if k <= 0 {
		return 0
	}
	return math.Pow(0.5, float64(k-1))
}

// ─── reviewerWeight ─────────────────────────────────────────────────────────

// reviewerWeight maps a reviewer's commerce activity to a credibility weight
// in [reviewerWeightFloor, 1.0]. nil global → 0.1 (no commerce history).
// We only use CompletedCount as the activity proxy in V1; can be expanded.
func reviewerWeight(global *model.CommerceScoreGlobal) float64 {
	if global == nil {
		return reviewerWeightFloor
	}
	saturation := float64(global.CompletedCount) / reviewerWeightSaturatedAt
	if saturation > 1 {
		saturation = 1
	}
	w := reviewerWeightBaselineFloor + reviewerWeightSpan*saturation
	return clamp01(w)
}

// ─── normalizeSentiment ─────────────────────────────────────────────────────

// normalizeSentiment maps a raw value in [-1, 1] to [0, 1]: n = (v+1)/2.
// Values outside [-1, 1] are clamped.
func normalizeSentiment(v float64) float64 {
	if v <= -1 {
		return 0
	}
	if v >= 1 {
		return 1
	}
	return (v + 1) / 2
}

// ─── applyNegativeAsymmetry ─────────────────────────────────────────────────

// applyNegativeAsymmetry amplifies normalized values below the pivot, leaving
// values at or above the pivot unchanged. Output is clamped to [0, 1].
//
//   n < 0.4: n_eff = 0.4 - (0.4 - n) * 1.3
//   n ≥ 0.4: n_eff = n
//
// Continuous at the pivot; never goes negative (clamped at 0).
func applyNegativeAsymmetry(n float64) float64 {
	if n >= asymmetryPivot {
		return clamp01(n)
	}
	out := asymmetryPivot - (asymmetryPivot-n)*asymmetryAmplifier
	return clamp01(out)
}

// ─── isSentimentEligible ────────────────────────────────────────────────────

// isSentimentEligible returns true iff the value distribution looks like a
// sentiment rating in [-1, 1]: ≥3 samples, median(|v|) ≤ 1, max(|v|) ≤ 2,
// stddev ≤ 1.
func isSentimentEligible(values []float64) bool {
	if len(values) < sentimentMinSamples {
		return false
	}
	abs := make([]float64, len(values))
	var sum, maxAbs float64
	for i, v := range values {
		a := math.Abs(v)
		abs[i] = a
		if a > maxAbs {
			maxAbs = a
		}
		sum += v
	}
	if maxAbs > sentimentMaxAbsExtreme {
		return false
	}
	// median of |v|
	sort.Float64s(abs)
	var med float64
	mid := len(abs) / 2
	if len(abs)%2 == 1 {
		med = abs[mid]
	} else {
		med = (abs[mid-1] + abs[mid]) / 2
	}
	if med > sentimentMaxAbsMedian {
		return false
	}
	// population stddev
	mean := sum / float64(len(values))
	var sq float64
	for _, v := range values {
		d := v - mean
		sq += d * d
	}
	stddev := math.Sqrt(sq / float64(len(values)))
	return stddev <= sentimentMaxStddev
}

// ─── tagAuthority ───────────────────────────────────────────────────────────

// tagAuthority computes a self-adaptive weight for a tag in [0, 1], based on
// breadth (unique reviewers), maturity (active days), and newcomer share.
//
//   authority = logTerm × timeTerm × (1 − newcomerDominance)
//   logTerm   = min(1, log(1+N) / log(1+50))
//   timeTerm  = min(1, activeDays / 30)
func tagAuthority(uniqueReviewers, activeDays int, newcomerDominance float64) float64 {
	if uniqueReviewers <= 0 {
		return 0
	}
	logTerm := math.Log1p(float64(uniqueReviewers)) / math.Log1p(authorityLogSaturateAt)
	if logTerm > 1 {
		logTerm = 1
	}
	timeTerm := float64(activeDays) / authorityTimeSaturateDays
	if timeTerm > 1 {
		timeTerm = 1
	}
	if timeTerm < 0 {
		timeTerm = 0
	}
	nd := clamp01(newcomerDominance)
	return clamp01(logTerm * timeTerm * (1 - nd))
}

// ─── hhi (Herfindahl–Hirschman Index) ──────────────────────────────────────

// hhi returns the concentration index of the given weights, normalized to
// [0, 1]. 1.0 means a single dominant entity; 1/n means perfectly equal.
// Empty or all-zero input returns 0.
func hhi(weights []float64) float64 {
	var total float64
	for _, w := range weights {
		if w > 0 {
			total += w
		}
	}
	if total <= 0 {
		return 0
	}
	var sumSq float64
	for _, w := range weights {
		if w <= 0 {
			continue
		}
		s := w / total
		sumSq += s * s
	}
	return clamp01(sumSq)
}

// ─── behavioralScore ────────────────────────────────────────────────────────

// behavioralScore computes Layer A — value-agnostic credit — from:
//   R              effective weighted unique reviewer count
//   revocationRate revoked / total, in [0, 1]
//   diversity      1 - hhi(reviewer weights), in [0, 1]
//
// Formula: σ(log(1+R)) · (1 − revocationRate) · diversity.
func behavioralScore(R, revocationRate, diversity float64) float64 {
	if R <= 0 {
		return 0
	}
	rr := clamp01(revocationRate)
	dv := clamp01(diversity)
	return clamp01(sigmoid(math.Log1p(R)) * (1 - rr) * dv)
}

// ─── fuseCreditScore ────────────────────────────────────────────────────────

// fuseCreditScore blends Layer A (behavioral) and Layer B (sentiment) into a
// final credit score, also returning the alpha actually used.
//
//   no sentiment tags → alpha=1.0, score = A
//   has sentiment tags → alpha=0.4, score = 0.4·A + 0.6·B
func fuseCreditScore(A, B float64, hasSentimentTags bool) (score, alpha float64) {
	if !hasSentimentTags {
		return clamp01(A), 1.0
	}
	alpha = fusionAlphaWithSentiment
	score = clamp01(alpha*A + (1-alpha)*B)
	return score, alpha
}
