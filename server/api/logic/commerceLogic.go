package logic

import (
	"agent_identity/model"
	"agent_identity/server/api/types"
	"fmt"
)

type CommerceScoreResp struct {
	Role                      string  `json:"role"`
	CompletedCount            int     `json:"completed_count"`
	RejectedCount             int     `json:"rejected_count"`
	ExpiredResponsibleCount   int     `json:"expired_responsible_count"`
	SuccessRate               float64 `json:"success_rate"`
	WeightedScore             float64 `json:"weighted_score"`
	TotalVolumeUSD            float64 `json:"total_volume_usd"`
	WeightedScoreUSD          float64 `json:"weighted_score_usd"`
	CreatedCount              int     `json:"created_count"`
	FundedCount               int     `json:"funded_count"`
	FundedRate                float64 `json:"funded_rate"`
	CompletionRate            float64 `json:"completion_rate"`
	EvaluatedCount            int     `json:"evaluated_count"`
	ExpiredFromSubmittedCount int     `json:"expired_from_submitted_count"`
	Responsiveness            float64 `json:"responsiveness"`
	TotalJobs                 int     `json:"total_jobs"`
	TotalVolume               float64 `json:"total_volume"`
	UniqueCounterparties      int     `json:"unique_counterparties"`
	Confidence                float64 `json:"confidence"`
}

func GetCommerceScores(uid uint64, chainID, contract string) ([]CommerceScoreResp, error) {
	var result []CommerceScoreResp

	if chainID != "" && contract != "" {
		scores, err := model.GetCommerceScoresByAgentUIDAndContract(uid, chainID, contract)
		if err != nil {
			return nil, fmt.Errorf("fail to get commerce scores: %v", err)
		}
		for _, s := range scores {
			result = append(result, segmentedToResp(s))
		}
	} else {
		globals, err := model.GetCommerceScoreGlobal(uid)
		if err != nil {
			return nil, fmt.Errorf("fail to get global commerce scores: %v", err)
		}
		for _, g := range globals {
			result = append(result, globalToResp(g))
		}
	}

	if result == nil {
		result = []CommerceScoreResp{}
	}
	return result, nil
}

func GetCommerceActions(uid uint64, role, action, certainty string, page, pageSize int) ([]model.CommerceAction, int64, error) {
	actions, total, err := model.GetCommerceActionsByAgentUID(uid, role, action, certainty, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get commerce actions: %v", err)
	}
	return actions, total, nil
}

type CommerceActionsParams struct {
	UID          uint64
	Role         string
	Actions      []string
	Certainty    string
	Polarities   []string
	ChainID      string
	Contract     string
	Counterparty string
	HookAddress  string
	HasHook      *bool
	MinBudget    *float64
	MaxBudget    *float64
	PaymentToken string
	TokenSymbol  string
	MinBudgetUSD *float64
	MaxBudgetUSD *float64
	StartTime    *uint64
	EndTime      *uint64
	SortBy       string
	Page         int
	PageSize     int
}

func GetCommerceActionsV2(p CommerceActionsParams) ([]model.CommerceAction, int64, error) {
	actions, total, err := model.GetCommerceActionsByQuery(model.CommerceActionsQuery{
		UID:          p.UID,
		Role:         p.Role,
		Actions:      p.Actions,
		Certainty:    p.Certainty,
		Polarities:   p.Polarities,
		ChainID:      p.ChainID,
		Contract:     p.Contract,
		Counterparty: p.Counterparty,
		HookAddress:  p.HookAddress,
		HasHook:      p.HasHook,
		MinBudget:    p.MinBudget,
		MaxBudget:    p.MaxBudget,
		PaymentToken: p.PaymentToken,
		TokenSymbol:  p.TokenSymbol,
		MinBudgetUSD: p.MinBudgetUSD,
		MaxBudgetUSD: p.MaxBudgetUSD,
		StartTime:    p.StartTime,
		EndTime:      p.EndTime,
		SortBy:       p.SortBy,
		Page:         p.Page,
		PageSize:     p.PageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get commerce actions: %v", err)
	}
	return actions, total, nil
}

func GetCommerceScoreSummary(uid uint64) ([]CommerceScoreResp, error) {
	globals, err := model.GetCommerceScoreGlobal(uid)
	if err != nil {
		return nil, fmt.Errorf("fail to get global commerce scores: %v", err)
	}
	result := make([]CommerceScoreResp, 0, len(globals))
	for _, g := range globals {
		result = append(result, globalToResp(g))
	}
	return result, nil
}

func segmentedToResp(s model.CommerceScore) CommerceScoreResp {
	resp := CommerceScoreResp{
		Role: s.Role, CompletedCount: s.CompletedCount, RejectedCount: s.RejectedCount,
		ExpiredResponsibleCount: s.ExpiredResponsibleCount, SuccessRate: s.SuccessRate,
		WeightedScore: s.WeightedScore, CreatedCount: s.CreatedCount, FundedCount: s.FundedCount,
		FundedRate: s.FundedRate, CompletionRate: s.CompletionRate,
		EvaluatedCount: s.EvaluatedCount, ExpiredFromSubmittedCount: s.ExpiredFromSubmittedCount,
		Responsiveness: s.Responsiveness, TotalJobs: s.TotalJobs, TotalVolume: s.TotalVolume,
		UniqueCounterparties: s.UniqueCounterparties, Confidence: s.Confidence,
	}
	if s.TotalVolumeUSD > 0 {
		resp.TotalVolumeUSD = s.TotalVolumeUSD
		resp.WeightedScoreUSD = s.WeightedVolumeUSD / s.TotalVolumeUSD
	}
	return resp
}

func globalToResp(g model.CommerceScoreGlobal) CommerceScoreResp {
	resp := CommerceScoreResp{
		Role: g.Role, CompletedCount: g.CompletedCount, RejectedCount: g.RejectedCount,
		ExpiredResponsibleCount: g.ExpiredResponsibleCount, SuccessRate: g.SuccessRate,
		WeightedScore: g.WeightedScore, CreatedCount: g.CreatedCount, FundedCount: g.FundedCount,
		FundedRate: g.FundedRate, CompletionRate: g.CompletionRate,
		EvaluatedCount: g.EvaluatedCount, ExpiredFromSubmittedCount: g.ExpiredFromSubmittedCount,
		Responsiveness: g.Responsiveness, TotalJobs: g.TotalJobs, TotalVolume: g.TotalVolume,
		UniqueCounterparties: g.UniqueCounterparties, Confidence: g.Confidence,
	}
	if g.TotalVolumeUSD > 0 {
		resp.TotalVolumeUSD = g.TotalVolumeUSD
		resp.WeightedScoreUSD = g.WeightedVolumeUSD / g.TotalVolumeUSD
	}
	return resp
}

// ─────────────── Commerce Stats ───────────────

type CommerceStatsParams struct {
	UID uint64
	Now uint64 // Unix timestamp of the latest known action, or system time
}

// GetCommerceStats returns full stats for an agent.
func GetCommerceStats(p CommerceStatsParams) (*types.CommerceStats, error) {
	// 1. Action breakdown
	actionCounts, err := model.GetCommerceActionCountsByRole(p.UID)
	if err != nil {
		return nil, fmt.Errorf("fail to get action counts: %v", err)
	}
	stats := &types.CommerceStats{
		ActionBreakdown:    make(types.ActionBreakdown),
		TimeSeries:         types.TimeSeriesStats{},
		BudgetDistribution: make(types.BudgetDistribution),
	}
	// Build action breakdown, filtering empty roles
	for role, actions := range actionCounts {
		hasData := false
		for _, c := range actions {
			if c > 0 {
				hasData = true
				break
			}
		}
		if !hasData {
			continue
		}
		ac := types.ActionCounts{
			JobCreated:   actions["job_created"],
			JobFunded:   actions["job_funded"],
			JobSubmitted: actions["job_submitted"],
			JobCompleted: actions["job_completed"],
			JobRejected:  actions["job_rejected"],
			JobExpired:   actions["job_expired"],
		}
		stats.ActionBreakdown[role] = ac
	}

	// 2. Time series — three windows
	windows := []struct {
		name       string
		startOffset uint64
		bucketDur   uint64
		out        *[]types.TimeBucket
	}{
		{"24h", 86400, 3600, &stats.TimeSeries.Hours24},
		{"7d", 604800, 21600, &stats.TimeSeries.Days7},
		{"30d", 2592000, 86400, &stats.TimeSeries.Days30},
	}
	for _, w := range windows {
		windowStart := p.Now - w.startOffset
		buckets, err := model.GetCommerceTimeSeries(p.UID, windowStart, p.Now, w.bucketDur)
		if err != nil {
			return nil, fmt.Errorf("fail to get time series (%s): %v", w.name, err)
		}
		// Filter out empty buckets (success_rate and counts both 0)
		filtered := make([]types.TimeBucket, 0, len(buckets))
		for _, b := range buckets {
			if b.CompletedCount > 0 || b.RejectedCount > 0 {
				filtered = append(filtered, types.TimeBucket{
					Bucket:         int64(b.BucketTime),
					CompletedCount: b.CompletedCount,
					RejectedCount:  b.RejectedCount,
					SuccessRate:    b.SuccessRate,
				})
			}
		}
		*w.out = filtered
	}

	// 3. Budget distribution — compute bucketing in Go
	rawItems, err := model.GetCommerceBudgetItems(p.UID)
	if err != nil {
		return nil, fmt.Errorf("fail to get budget items: %v", err)
	}
	stats.BudgetDistribution = computeBudgetDistribution(rawItems)

	return stats, nil
}

// computeBudgetDistribution groups raw budget items by role and contract,
// then percentile-buckets them.
func computeBudgetDistribution(items []model.RawBudgetItem) types.BudgetDistribution {
	if len(items) == 0 {
		return make(types.BudgetDistribution)
	}

	// Group by role+contract
	type key struct {
		role     string
		contract string
	}
	groups := make(map[key][]float64)
	for _, item := range items {
		k := key{item.Role, item.Contract}
		groups[k] = append(groups[k], item.Budget)
	}

	result := make(types.BudgetDistribution)

	for k, budgets := range groups {
		totalCount := len(budgets)
		if totalCount == 0 {
			continue
		}

		// Compute 33rd and 66th percentiles
		p33 := percentile(budgets, 0.33)
		p66 := percentile(budgets, 0.66)

		var smallCount, mediumCount, largeCount int
		var smallMax, mediumMax, largeMax float64

		if totalCount < 3 {
			// Flat bucket: no segmentation
			maxAmt := budgets[0]
			for _, b := range budgets[1:] {
				if b > maxAmt {
					maxAmt = b
				}
			}
			result[k.role] = map[string]types.ContractBudgetStats{
				k.contract: {
					TotalCount: totalCount,
					Buckets: types.BudgetBuckets{
						Flat: map[string]any{
							"count":       totalCount,
							"max_amount": fmt.Sprintf("%f", maxAmt),
						},
					},
				},
			}
			continue
		}

		// Bucket assignment
		var small, medium, large []float64
		for _, b := range budgets {
			if b < p33 {
				small = append(small, b)
			} else if b < p66 {
				medium = append(medium, b)
			} else {
				large = append(large, b)
			}
		}
		smallCount = len(small)
		mediumCount = len(medium)
		largeCount = len(large)
		var smallStat, mediumStat, largeStat *types.BudgetBucketStat
		if len(small) > 0 {
			smallMax = small[len(small)-1]
			smallStat = &types.BudgetBucketStat{Count: smallCount, MaxAmount: fmt.Sprintf("%f", smallMax)}
		}
		if len(medium) > 0 {
			mediumMax = medium[len(medium)-1]
			mediumStat = &types.BudgetBucketStat{Count: mediumCount, MaxAmount: fmt.Sprintf("%f", mediumMax)}
		}
		if len(large) > 0 {
			largeMax = large[len(large)-1]
			largeStat = &types.BudgetBucketStat{Count: largeCount, MaxAmount: fmt.Sprintf("%f", largeMax)}
		}

		if result[k.role] == nil {
			result[k.role] = make(map[string]types.ContractBudgetStats)
		}
		result[k.role][k.contract] = types.ContractBudgetStats{
			TotalCount: totalCount,
			Buckets: types.BudgetBuckets{
				Small:  smallStat,
				Medium: mediumStat,
				Large:  largeStat,
			},
		}
	}

	// Remove roles with no data
	for role, contracts := range result {
		if len(contracts) == 0 {
			delete(result, role)
		}
	}
	return result
}

// percentile returns the p-th percentile of a sorted slice (0 < p < 1).
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	n := float64(len(sorted))
	idx := p * (n - 1)
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	frac := idx - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}
