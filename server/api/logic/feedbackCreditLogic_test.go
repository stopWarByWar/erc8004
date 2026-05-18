package logic

import (
	"math"
	"testing"

	"agent_identity/model"
)

// floatEq compares two floats within an epsilon tolerance.
func floatEq(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

// ─── timeDecay ───────────────────────────────────────────────────────────────

func TestTimeDecay(t *testing.T) {
	const day = uint64(86400)
	tests := []struct {
		name        string
		blockTs     uint64
		nowTs       uint64
		halfLifeDay float64
		want        float64
		eps         float64
	}{
		{"now == block → 1.0", 1000, 1000, 180, 1.0, 1e-9},
		{"future block (now < block) → 1.0", 2000, 1000, 180, 1.0, 1e-9},
		{"one half-life ago → 0.5", 1000, 1000 + 180*day, 180, 0.5, 1e-9},
		{"two half-lives ago → 0.25", 1000, 1000 + 360*day, 180, 0.25, 1e-9},
		{"three half-lives ago → 0.125", 1000, 1000 + 540*day, 180, 0.125, 1e-9},
		{"zero half-life is invalid → 1.0 fallback", 1000, 2000, 0, 1.0, 1e-9},
		{"negative half-life is invalid → 1.0 fallback", 1000, 2000, -1, 1.0, 1e-9},
		{"shorter half-life decays faster", 0, 90 * day, 90, 0.5, 1e-9},
	}
	for _, tt := range tests {
		got := timeDecay(tt.blockTs, tt.nowTs, tt.halfLifeDay)
		if !floatEq(got, tt.want, tt.eps) {
			t.Errorf("%s: timeDecay(%d,%d,%.0f)=%.6f want %.6f", tt.name, tt.blockTs, tt.nowTs, tt.halfLifeDay, got, tt.want)
		}
		if got < 0 || got > 1 {
			t.Errorf("%s: result %.6f out of [0,1]", tt.name, got)
		}
	}
}

// ─── dedupeFactor ───────────────────────────────────────────────────────────

func TestDedupeFactor(t *testing.T) {
	tests := []struct {
		k    int
		want float64
	}{
		{0, 0},
		{-1, 0},
		{1, 1.0},
		{2, 0.5},
		{3, 0.25},
		{4, 0.125},
		{10, 1.0 / 512},
	}
	for _, tt := range tests {
		got := dedupeFactor(tt.k)
		if !floatEq(got, tt.want, 1e-12) {
			t.Errorf("dedupeFactor(%d)=%.10f want %.10f", tt.k, got, tt.want)
		}
	}
}

// ─── reviewerWeight ─────────────────────────────────────────────────────────

func TestReviewerWeight(t *testing.T) {
	tests := []struct {
		name   string
		global *model.CommerceScoreGlobal
		want   float64
		eps    float64
	}{
		{"nil → 0.1 (no history)", nil, 0.1, 1e-9},
		{"zero counts → 0.3 (baseline floor)", &model.CommerceScoreGlobal{}, 0.3, 1e-9},
		{"10 completed → 0.3 + 0.7*0.5 = 0.65", &model.CommerceScoreGlobal{CompletedCount: 10}, 0.65, 1e-9},
		{"20 completed → 0.3 + 0.7*1.0 = 1.0", &model.CommerceScoreGlobal{CompletedCount: 20}, 1.0, 1e-9},
		{"100 completed → cap at 1.0", &model.CommerceScoreGlobal{CompletedCount: 100}, 1.0, 1e-9},
		{"5 completed → 0.3 + 0.7*0.25 = 0.475", &model.CommerceScoreGlobal{CompletedCount: 5}, 0.475, 1e-9},
	}
	for _, tt := range tests {
		got := reviewerWeight(tt.global)
		if !floatEq(got, tt.want, tt.eps) {
			t.Errorf("%s: reviewerWeight=%.6f want %.6f", tt.name, got, tt.want)
		}
		if got < 0 || got > 1 {
			t.Errorf("%s: result %.6f out of [0,1]", tt.name, got)
		}
	}
}

// ─── detectScale ────────────────────────────────────────────────────────────

func TestDetectScale(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   string // expected scale name, "" for nil
	}{
		{"empty → nil", []float64{}, ""},
		{"signed [-1,0,1] → signed_unit", []float64{-1, 0, 1}, "signed_unit"},
		{"unit cluster {0.4,0.6,0.9} → unit", []float64{0.4, 0.6, 0.9}, "unit"},
		{"five-star {4,5,5,3} → five_star", []float64{4, 5, 5, 3}, "five_star"},
		{"ten-point {6,7,8,5,9} → ten_point", []float64{6, 7, 8, 5, 9}, "ten_point"},
		{"percent {60,70,80,85,75} → percent", []float64{60, 70, 80, 85, 75}, "percent"},
		{"way out of range {150,200} → nil", []float64{150, 200}, ""},
		{"unit boundary {0,1} → unit", []float64{0, 1}, "unit"},
		{"signed boundary {-1,1} → signed_unit", []float64{-1, 1}, "signed_unit"},
		{"unit upper micro-overshoot 1.04 → unit (5% tol)", []float64{0.5, 1.04}, "unit"},
		{"any negative in unsigned-only data → signed_unit", []float64{-0.1, 0.5, 0.7}, "signed_unit"},
		{"percent {0,50,100} → percent (boundary)", []float64{0, 50, 100}, "percent"},
	}
	for _, tt := range tests {
		got := detectScale(tt.values)
		gotName := ""
		if got != nil {
			gotName = got.Name
		}
		if gotName != tt.want {
			t.Errorf("%s: detectScale(%v)=%q want %q", tt.name, tt.values, gotName, tt.want)
		}
	}
}

// ─── normalizeSentiment ─────────────────────────────────────────────────────

func TestNormalizeSentiment(t *testing.T) {
	unit := &ratingScale{Name: "unit", Min: 0, Max: 1}
	signed := &ratingScale{Name: "signed_unit", Min: -1, Max: 1}
	fiveStar := &ratingScale{Name: "five_star", Min: 0, Max: 5}
	percent := &ratingScale{Name: "percent", Min: 0, Max: 100}

	tests := []struct {
		name  string
		v     float64
		scale *ratingScale
		want  float64
	}{
		// signed_unit (original behavior preserved)
		{"signed -1 → 0", -1.0, signed, 0.0},
		{"signed 0 → 0.5", 0.0, signed, 0.5},
		{"signed 1 → 1", 1.0, signed, 1.0},
		{"signed 0.5 → 0.75", 0.5, signed, 0.75},
		{"signed -0.5 → 0.25", -0.5, signed, 0.25},
		{"signed -2 clamps to 0", -2.0, signed, 0.0},
		{"signed 2 clamps to 1", 2.0, signed, 1.0},
		// unit
		{"unit 0 → 0", 0.0, unit, 0.0},
		{"unit 0.5 → 0.5", 0.5, unit, 0.5},
		{"unit 1 → 1", 1.0, unit, 1.0},
		// five_star
		{"five_star 4 → 0.8", 4.0, fiveStar, 0.8},
		{"five_star 5 → 1", 5.0, fiveStar, 1.0},
		// percent
		{"percent 80 → 0.8", 80, percent, 0.8},
		{"percent 0 → 0", 0, percent, 0.0},
		{"percent 100 → 1", 100, percent, 1.0},
		{"percent 150 clamps to 1", 150, percent, 1.0},
		// nil scale: defensive fallback
		{"nil scale → 0.5 neutral", 42, nil, 0.5},
	}
	for _, tt := range tests {
		got := normalizeSentiment(tt.v, tt.scale)
		if !floatEq(got, tt.want, 1e-9) {
			t.Errorf("%s: normalizeSentiment(%.4f, %v)=%.6f want %.6f",
				tt.name, tt.v, tt.scale, got, tt.want)
		}
		if got < 0 || got > 1 {
			t.Errorf("%s: result %.6f out of [0,1]", tt.name, got)
		}
	}
}

// ─── applyNegativeAsymmetry ─────────────────────────────────────────────────

func TestApplyNegativeAsymmetry(t *testing.T) {
	tests := []struct {
		name string
		n    float64
		want float64
		eps  float64
	}{
		{"continuity at 0.4", 0.4, 0.4, 1e-9},
		{"above 0.4 unchanged: 0.5", 0.5, 0.5, 1e-9},
		{"above 0.4 unchanged: 1.0", 1.0, 1.0, 1e-9},
		{"below 0.4 amplified: 0.2 → 0.14", 0.2, 0.14, 1e-9},
		{"below 0.4 amplified: 0.3 → 0.27", 0.3, 0.27, 1e-9},
		{"clamp lower bound at 0", 0.0, 0.0, 1e-9},
		{"near zero clamped to 0", 0.05, 0.0, 1e-9},
	}
	for _, tt := range tests {
		got := applyNegativeAsymmetry(tt.n)
		if !floatEq(got, tt.want, tt.eps) {
			t.Errorf("%s: applyNegativeAsymmetry(%.2f)=%.6f want %.6f", tt.name, tt.n, got, tt.want)
		}
		if got < 0 || got > 1 {
			t.Errorf("%s: result %.6f out of [0,1]", tt.name, got)
		}
	}

	// Monotonicity: f should be non-decreasing in n on [0,1]
	prev := -1.0
	for n := 0.0; n <= 1.0; n += 0.05 {
		got := applyNegativeAsymmetry(n)
		if got < prev-1e-9 {
			t.Errorf("non-monotonic: f(%.2f)=%.4f < prev %.4f", n, got, prev)
		}
		prev = got
	}
}

// ─── isSentimentEligible ────────────────────────────────────────────────────

func TestIsSentimentEligible(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		wantOK    bool
		wantScale string // "" when wantOK=false
	}{
		// sample-count gate
		{"empty → false", []float64{}, false, ""},
		{"single sample → false", []float64{0.5}, false, ""},
		{"two samples → false", []float64{0.5, 0.7}, false, ""},

		// signed_unit (backward compat with original [-1,1] design)
		{"signed cluster {0.5,0.7,-0.3} → signed_unit", []float64{0.5, 0.7, -0.3}, true, "signed_unit"},
		{"signed extremes {-1,0,1} → signed_unit", []float64{-1.0, 0.0, 1.0}, true, "signed_unit"},

		// new scales
		{"percent {60,70,80,90,85} → percent", []float64{60, 70, 80, 90, 85}, true, "percent"},
		{"five_star {4,5,5,3,4} → five_star", []float64{4, 5, 5, 3, 4}, true, "five_star"},
		{"ten_point {6,7,8,5,9} → ten_point", []float64{6, 7, 8, 5, 9}, true, "ten_point"},
		{"unit {0.4,0.6,0.9} → unit", []float64{0.4, 0.6, 0.9}, true, "unit"},

		// out-of-range / no preset fits
		{"{0.5,0.6,850} → false (no preset)", []float64{0.5, 0.6, 850}, false, ""},
		{"{150,200,180} → false (above all presets)", []float64{150, 200, 180}, false, ""},

		// stddev gate (scale-relative: > 50% × span). For bounded data
		// within a preset, max possible population stddev is 0.5×span (perfect
		// bimodal), so the gate only ever rejects values that *also* breach
		// the upper bound — covered by the no-preset cases above.
		{"high-variance signed {-2,0,2} → false (no preset fits)",
			[]float64{-2.0, 0.0, 2.0}, false, ""},

		// metric-style distributions still correctly rejected
		{"latency cluster {200..1200} → false", []float64{200, 350, 500, 800, 1200}, false, ""},
	}
	for _, tt := range tests {
		gotOK, gotScale := isSentimentEligible(tt.values)
		if gotOK != tt.wantOK {
			t.Errorf("%s: ok=%v want %v", tt.name, gotOK, tt.wantOK)
		}
		gotName := ""
		if gotScale != nil {
			gotName = gotScale.Name
		}
		if gotName != tt.wantScale {
			t.Errorf("%s: scale=%q want %q", tt.name, gotName, tt.wantScale)
		}
	}
}

// ─── tagAuthority ───────────────────────────────────────────────────────────

func TestTagAuthority(t *testing.T) {
	t.Run("zero reviewers → 0", func(t *testing.T) {
		got := tagAuthority(0, 30, 0)
		if !floatEq(got, 0, 1e-9) {
			t.Errorf("got %.4f want 0", got)
		}
	})

	t.Run("monotonic in uniqueReviewers", func(t *testing.T) {
		prev := -1.0
		for _, n := range []int{1, 2, 5, 10, 20, 50, 100} {
			got := tagAuthority(n, 30, 0)
			if got < prev-1e-9 {
				t.Errorf("non-monotonic at n=%d: %.4f < prev %.4f", n, got, prev)
			}
			prev = got
		}
	})

	t.Run("clamped to [0,1]", func(t *testing.T) {
		for _, n := range []int{1, 10, 100, 1000, 10000} {
			for _, d := range []int{0, 10, 30, 365} {
				for _, nd := range []float64{0, 0.5, 1.0} {
					got := tagAuthority(n, d, nd)
					if got < 0 || got > 1 {
						t.Errorf("authority(%d,%d,%.1f)=%.4f out of [0,1]", n, d, nd, got)
					}
				}
			}
		}
	})

	t.Run("activeDays<30 penalty", func(t *testing.T) {
		small := tagAuthority(50, 10, 0)
		mature := tagAuthority(50, 30, 0)
		if !(small < mature) {
			t.Errorf("expected small (%.4f) < mature (%.4f)", small, mature)
		}
		// activeDays >= 30 should saturate
		later := tagAuthority(50, 365, 0)
		if !floatEq(later, mature, 1e-9) {
			t.Errorf("expected saturation: 365d (%.4f) == 30d (%.4f)", later, mature)
		}
	})

	t.Run("newcomerDominance penalty", func(t *testing.T) {
		base := tagAuthority(50, 30, 0)
		half := tagAuthority(50, 30, 0.5)
		full := tagAuthority(50, 30, 1.0)
		if !(half < base) {
			t.Errorf("expected half (%.4f) < base (%.4f)", half, base)
		}
		if !floatEq(full, 0, 1e-9) {
			t.Errorf("newcomerDominance=1 should fully zero out, got %.4f", full)
		}
	})
}

// ─── sigmoid ────────────────────────────────────────────────────────────────

func TestSigmoid(t *testing.T) {
	tests := []struct {
		x, want float64
	}{
		{0, 0},
		{-1, 0}, // negatives clamp to 0
		{3, 0.5},
		{1, 0.25},
		{9, 0.75},
		{1e6, 1.0 - 1e-6}, // approaches 1 for large x
	}
	for _, tt := range tests {
		got := sigmoid(tt.x)
		if math.Abs(got-tt.want) > 1e-3 {
			t.Errorf("sigmoid(%.2f)=%.6f want ≈%.6f", tt.x, got, tt.want)
		}
		if got < 0 || got > 1 {
			t.Errorf("sigmoid(%.2f)=%.6f out of [0,1]", tt.x, got)
		}
	}
}

// ─── hhi (Herfindahl–Hirschman Index) ───────────────────────────────────────

func TestHHI(t *testing.T) {
	tests := []struct {
		name    string
		weights []float64
		want    float64
		eps     float64
	}{
		{"empty → 0", []float64{}, 0, 1e-9},
		{"all zeros → 0", []float64{0, 0, 0}, 0, 1e-9},
		{"single reviewer → 1 (max concentration)", []float64{1.0}, 1.0, 1e-9},
		{"two equal → 0.5", []float64{1.0, 1.0}, 0.5, 1e-9},
		{"four equal → 0.25", []float64{2, 2, 2, 2}, 0.25, 1e-9},
		{"dominant: 9 vs 1 → 0.82", []float64{9, 1}, 0.82, 1e-9},
		{"three equal → 1/3", []float64{1, 1, 1}, 1.0 / 3, 1e-9},
	}
	for _, tt := range tests {
		got := hhi(tt.weights)
		if !floatEq(got, tt.want, tt.eps) {
			t.Errorf("%s: hhi(%v)=%.6f want %.6f", tt.name, tt.weights, got, tt.want)
		}
		if got < 0 || got > 1 {
			t.Errorf("%s: hhi out of [0,1]: %.6f", tt.name, got)
		}
	}
}

// ─── behavioralScore ────────────────────────────────────────────────────────

func TestBehavioralScore(t *testing.T) {
	tests := []struct {
		name           string
		R              float64
		revocationRate float64
		diversity      float64
		// We only check broad properties + a known point
		mustPositive bool
		mustBeZero   bool
	}{
		{"R=0 → 0", 0, 0, 1, false, true},
		{"diversity=0 → 0", 5, 0, 0, false, true},
		{"all revoked → 0", 5, 1, 1, false, true},
		{"healthy: R=10, no revoke, full diversity", 10, 0, 1, true, false},
		{"healthy bigger R should be ≥ smaller R", 100, 0, 1, true, false},
	}
	for _, tt := range tests {
		got := behavioralScore(tt.R, tt.revocationRate, tt.diversity)
		if got < 0 || got > 1 {
			t.Errorf("%s: behavioralScore=%.4f out of [0,1]", tt.name, got)
		}
		if tt.mustBeZero && !floatEq(got, 0, 1e-9) {
			t.Errorf("%s: expected 0, got %.4f", tt.name, got)
		}
		if tt.mustPositive && !(got > 0) {
			t.Errorf("%s: expected > 0, got %.4f", tt.name, got)
		}
	}

	// Monotonicity in R
	prev := -1.0
	for _, R := range []float64{1, 2, 5, 10, 20, 50, 100} {
		got := behavioralScore(R, 0, 1)
		if got < prev-1e-9 {
			t.Errorf("non-monotonic in R: R=%.1f → %.4f < prev %.4f", R, got, prev)
		}
		prev = got
	}

	// Revocation penalty
	high := behavioralScore(20, 0, 1)
	low := behavioralScore(20, 0.5, 1)
	if !(low < high) {
		t.Errorf("revocation should lower score: low=%.4f high=%.4f", low, high)
	}
}

// ─── fuseCreditScore ────────────────────────────────────────────────────────

func TestFuseCreditScore(t *testing.T) {
	tests := []struct {
		name             string
		A, B             float64
		hasSentimentTags bool
		wantScore        float64
		wantAlpha        float64
		eps              float64
	}{
		{"no sentiment tags → α=1, pure Layer A", 0.7, 0, false, 0.7, 1.0, 1e-9},
		{"no sentiment tags ignores B", 0.7, 0.99, false, 0.7, 1.0, 1e-9},
		{"with sentiment tags → α=0.4", 0.5, 0.8, true, 0.4*0.5 + 0.6*0.8, 0.4, 1e-9},
		{"both layers equal → score equals A=B", 0.6, 0.6, true, 0.6, 0.4, 1e-9},
	}
	for _, tt := range tests {
		score, alpha := fuseCreditScore(tt.A, tt.B, tt.hasSentimentTags)
		if !floatEq(score, tt.wantScore, tt.eps) {
			t.Errorf("%s: score=%.6f want %.6f", tt.name, score, tt.wantScore)
		}
		if !floatEq(alpha, tt.wantAlpha, tt.eps) {
			t.Errorf("%s: alpha=%.6f want %.6f", tt.name, alpha, tt.wantAlpha)
		}
		if score < 0 || score > 1 {
			t.Errorf("%s: score out of [0,1]: %.6f", tt.name, score)
		}
	}
}
