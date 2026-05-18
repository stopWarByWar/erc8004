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
	sentimentMinSamples = 3
	// Scale-relative stddev gate: values are considered "rating-like" only if
	// their population stddev is ≤ this fraction of the matched scale's span.
	// Replaces the old absolute stddev=1.0 threshold so 0-100 / 0-5 / 0-10
	// distributions are treated on the same footing as the original [-1, 1].
	sentimentMaxStddevFraction = 0.5
	// 5% tolerance on the upper bound for floating-point / micro-overshoot
	// (e.g. percent value 100.5). For unsigned scales the lower bound stays
	// strict at 0 so a single negative value forces fallback to signed_unit.
	scaleBoundTolerance = 0.05

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

// ─── ratingScale ────────────────────────────────────────────────────────────

// ratingScale is one of the known rating ranges we auto-detect. Sentiment
// detection is multi-scale: we accept [-1,1] (signed_unit), [0,1] (unit),
// [0,5] (five_star), [0,10] (ten_point) and [0,100] (percent). After
// detection every value is normalized to [0,1] before fusion, so the final
// credit_score range stays [0,1] (commerce-aligned).
type ratingScale struct {
	Name string  // "unit" | "five_star" | "ten_point" | "percent" | "signed_unit"
	Min  float64
	Max  float64
}

// ratingScalePresets is intentionally ordered: positive-only scales come
// first in ascending span order, so positive-clustered data (e.g. {0.5, 0.7})
// matches unit rather than the wider signed_unit. signed_unit comes last as
// the fallback for any distribution containing negative values.
var ratingScalePresets = []ratingScale{
	{Name: "unit", Min: 0, Max: 1},
	{Name: "five_star", Min: 0, Max: 5},
	{Name: "ten_point", Min: 0, Max: 10},
	{Name: "percent", Min: 0, Max: 100},
	{Name: "signed_unit", Min: -1, Max: 1},
}

// detectScale returns the smallest preset that contains every observed value
// (with a small upper-bound tolerance for floating-point overshoot), or nil
// when no preset fits. Empty input → nil.
//
// Lower-bound semantics: for unsigned scales (Min=0) we require min ≥ 0
// strictly — so any negative value forces signed_unit (or nil). For
// signed_unit we allow a symmetric tolerance.
func detectScale(values []float64) *ratingScale {
	if len(values) == 0 {
		return nil
	}
	minV, maxV := values[0], values[0]
	for _, v := range values {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	for i := range ratingScalePresets {
		s := &ratingScalePresets[i]
		span := s.Max - s.Min
		hiBound := s.Max + scaleBoundTolerance*span
		// Unsigned scales: strict lower bound to keep negatives out.
		loBound := s.Min
		if s.Min < 0 {
			loBound = s.Min - scaleBoundTolerance*span
		}
		if minV >= loBound && maxV <= hiBound {
			return s
		}
	}
	return nil
}

// ─── normalizeSentiment ─────────────────────────────────────────────────────

// normalizeSentiment linearly maps v from scale.[Min, Max] to [0, 1] and
// clamps. A nil scale falls back to a neutral 0.5 (defensive — callers
// should only normalize after isSentimentEligible returned a non-nil scale).
func normalizeSentiment(v float64, scale *ratingScale) float64 {
	if scale == nil {
		return 0.5
	}
	span := scale.Max - scale.Min
	if span <= 0 {
		return 0.5
	}
	return clamp01((v - scale.Min) / span)
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

// isSentimentEligible decides whether a value distribution looks like ratings
// on one of the known scales (unit / five_star / ten_point / percent /
// signed_unit). On success it returns the matched scale so callers can
// normalize without re-running detection.
//
// Rules:
//   1. ≥ sentimentMinSamples (3) samples.
//   2. All values fit one preset scale (see detectScale).
//   3. Population stddev ≤ sentimentMaxStddevFraction × span — keeps
//      bimodal/noisy distributions out even when they technically fit the
//      range (e.g. {0, 100, 0, 100} on percent).
func isSentimentEligible(values []float64) (bool, *ratingScale) {
	if len(values) < sentimentMinSamples {
		return false, nil
	}
	scale := detectScale(values)
	if scale == nil {
		return false, nil
	}
	mean := meanOf(values)
	stddev := populationStddev(values, mean)
	span := scale.Max - scale.Min
	if stddev > sentimentMaxStddevFraction*span {
		return false, nil
	}
	return true, scale
}

// meanOf returns the arithmetic mean. len(values) must be > 0.
func meanOf(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// populationStddev returns sqrt(Σ(v-mean)²/n). len(values) must be > 0.
func populationStddev(values []float64, mean float64) float64 {
	var sq float64
	for _, v := range values {
		d := v - mean
		sq += d * d
	}
	return math.Sqrt(sq / float64(len(values)))
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
