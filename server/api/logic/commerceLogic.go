package logic

import (
	"context"
	"agent_identity/config"
	"agent_identity/model"
	"agent_identity/server/api/types"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
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

// ─────────────── Commerce Jobs (Job Browser) ───────────────

type CommerceJobsParams struct {
	ChainID          string
	CommerceContract string
	Status           string
	Role             string
	AgentAddress     string
	Counterparty     string
	PaymentToken     string
	TokenSymbol      string

	MinBudget    *float64
	MaxBudget    *float64
	MinBudgetUSD *float64
	MaxBudgetUSD *float64

	StartTime *uint64
	EndTime   *uint64

	SortBy    string
	SortOrder string

	Page     int
	PageSize int
}

type CommerceJobsChartsParams struct {
	CommerceJobsParams
	Window        string
	BucketSeconds uint64
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	if pageSize > 100 {
		return 100
	}
	return pageSize
}

func normalizeSortOrder(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return "desc"
	}
	if v != "asc" && v != "desc" {
		return "desc"
	}
	return v
}

func normalizeJobsSortBy(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	switch v {
	case "", "updated_at", "budget", "budget_usd", "paid_amount_usd":
		return v
	default:
		return ""
	}
}

func normalizeChartsWindow(v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "7d", nil
	}
	switch v {
	case "24h", "3d", "7d", "1month", "all", "30d", "custom":
		return v, nil
	default:
		return "", errors.New("invalid window")
	}
}

func defaultBucketSeconds(window string) uint64 {
	switch window {
	case "24h":
		return 3600
	case "3d":
		return 3 * 3600
	case "7d":
		return 6 * 3600
	case "1month":
		return 24 * 3600
	case "all":
		return 24 * 3600
	case "30d":
		return 24 * 3600
	default:
		return 0
	}
}

func jobToDTO(j model.CommerceJob) types.CommerceJobDTO {
	dto := types.CommerceJobDTO{
		ChainID:          j.ChainID,
		CommerceContract: j.CommerceContract,
		JobID:            j.JobID,
		Status:           j.Status,
		UpdatedAt:        j.UpdatedAt,
		Client:           j.Client,
		Provider:         j.Provider,
		Evaluator:        j.Evaluator,
		Description:      j.Description,
		HookAddress:      j.HookAddress,
		ExpiredAt:        j.ExpiredAt,
		SubmittedAt:      j.SubmittedAt,
		CompletedAt:      j.CompletedAt,
		PaymentToken:     j.PaymentToken,
		PaymentDecimals:  normalizePaymentDecimals(j.PaymentDecimals),
		TokenSymbol:      j.TokenSymbol,
		Budget:           formatAmount(j.Budget),
		BudgetUSD:        j.BudgetUSD,
		PaidAmount:       formatAmount(j.PaidAmount),
		PaidAmountUSD:    j.PaidAmountUSD,
		PlatformFeeAmount:  formatAmount(j.PlatformFeeAmount),
		PlatformFeeUSD:     j.PlatformFeeUSD,
		EvaluatorFeeAmount: formatAmount(j.EvaluatorFeeAmount),
		EvaluatorFeeUSD:    j.EvaluatorFeeUSD,
		LatestBlockNumber:  j.LatestBlockNumber,
		LatestTxHash:       j.LatestTxHash,
		LatestActionUID:    j.LatestActionUID,
	}
	if chain, ok := config.GetChainInfo(j.ChainID); ok {
		dto.ChainName = chain.ChainName
		dto.ChainLogo = chain.ChainLogo
	}
	return dto
}

func normalizePaymentDecimals(v uint) uint {
	if v == 0 {
		return 18
	}
	return v
}

// formatAmount formats a numeric(36,8) value stored in float64 into a stable string.
// It avoids scientific notation and trims trailing zeros.
func formatAmount(v float64) string {
	if v == 0 || math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	// Keep exactly 8 decimals (matches schema), then trim.
	s := fmt.Sprintf("%.8f", v)
	// Trim trailing zeros and dot.
	for len(s) > 1 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 1 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

func ActionToDTO(a model.CommerceAction) types.CommerceActionDTO {
	return types.CommerceActionDTO{
		ChainID:          a.ChainID,
		CommerceContract: a.CommerceContract,
		JobID:            a.JobID,
		AgentUID:         a.AgentUID,
		AgentAddress:     a.AgentAddress,
		Role:             a.Role,
		Action:           a.Action,
		SignalPolarity:   a.SignalPolarity,
		SignalWeight:     a.SignalWeight,
		SignalCertainty:  a.SignalCertainty,
		JobBudget:        formatAmount(a.JobBudget),
		BudgetUSD:        a.BudgetUSD,
		PaymentToken:     a.PaymentToken,
		PaymentDecimals:  normalizePaymentDecimals(a.PaymentDecimals),
		TokenSymbol:      a.TokenSymbol,
		Counterparty:     a.Counterparty,
		Reason:           a.Reason,
		Deliverable:      a.Deliverable,
		PreviousStatus:   a.PreviousStatus,
		HookAddress:      a.HookAddress,
		BlockNumber:      a.BlockNumber,
		TxHash:           a.TxHash,
		LogIndex:         a.LogIndex,
		BlockTimestamp:   a.BlockTimestamp,
	}
}

func GetCommerceJobs(p CommerceJobsParams) ([]types.CommerceJobDTO, int64, error) {
	q := model.CommerceJobsQuery{
		ChainID:       p.ChainID,
		Contract:      p.CommerceContract,
		Status:        p.Status,
		Role:          p.Role,
		AgentAddress:  p.AgentAddress,
		Counterparty:  p.Counterparty,
		PaymentToken:  p.PaymentToken,
		TokenSymbol:   p.TokenSymbol,
		MinBudget:     p.MinBudget,
		MaxBudget:     p.MaxBudget,
		MinBudgetUSD:  p.MinBudgetUSD,
		MaxBudgetUSD:  p.MaxBudgetUSD,
		StartTime:     p.StartTime,
		EndTime:       p.EndTime,
		SortBy:        normalizeJobsSortBy(p.SortBy),
		SortOrder:     normalizeSortOrder(p.SortOrder),
		Page:          normalizePage(p.Page),
		PageSize:      normalizePageSize(p.PageSize),
	}
	jobs, total, err := model.GetCommerceJobs(q)
	if err != nil {
		return nil, 0, fmt.Errorf("get commerce jobs: %w", err)
	}
	out := make([]types.CommerceJobDTO, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, jobToDTO(j))
	}
	return out, total, nil
}

func GetCommerceJobDetail(chainID, contract string, jobID uint64) (*types.CommerceJobDTO, error) {
	j, err := model.GetCommerceJobByID(chainID, contract, jobID)
	if err != nil {
		return nil, fmt.Errorf("get commerce job detail: %w", err)
	}
	if j == nil {
		return nil, nil
	}
	dto := jobToDTO(*j)
	return &dto, nil
}

type CommerceJobsGeneral struct {
	Summary types.CommerceJobsGeneralSummaryResp       `json:"summary"`
	Distributions types.CommerceJobsGeneralDistributionsResp `json:"distributions"`
}

func GetCommerceJobsGeneral(p CommerceJobsParams, includeDistributions bool) (*CommerceJobsGeneral, error) {
	q := model.CommerceJobsQuery{
		ChainID:       p.ChainID,
		Contract:      p.CommerceContract,
		Status:        p.Status,
		Role:          p.Role,
		AgentAddress:  p.AgentAddress,
		Counterparty:  p.Counterparty,
		PaymentToken:  p.PaymentToken,
		TokenSymbol:   p.TokenSymbol,
		MinBudget:     p.MinBudget,
		MaxBudget:     p.MaxBudget,
		MinBudgetUSD:  p.MinBudgetUSD,
		MaxBudgetUSD:  p.MaxBudgetUSD,
		StartTime:     p.StartTime,
		EndTime:       p.EndTime,
		Page:          1,
		PageSize:      1,
	}
	summary, tokenDist, ccDist, fees, err := model.GetCommerceJobsGeneral(q, includeDistributions)
	if err != nil {
		return nil, fmt.Errorf("get commerce jobs general: %w", err)
	}

	// Enrich chain meta for UI display (chain name/logo).
	for i := range ccDist {
		if ccDist[i].ChainID == "" {
			continue
		}
		if chain, ok := config.GetChainInfo(ccDist[i].ChainID); ok {
			ccDist[i].ChainName = chain.ChainName
			ccDist[i].ChainLogo = chain.ChainLogo
		}
	}

	grouped := groupChainContracts(ccDist)
	return &CommerceJobsGeneral{
		Summary: types.CommerceJobsGeneralSummaryResp{
			JobsCount:       summary.JobsCount,
			PaidVolumeUSD:   summary.PaidVolumeUSD,
			BudgetVolumeUSD: summary.BudgetVolumeUSD,
			OutcomeMix: map[string]int64{
				"completed": summary.OutcomeCompleted,
				"rejected":  summary.OutcomeRejected,
				"expired":   summary.OutcomeExpired,
			},
			ActiveMix: map[string]int64{
				"open":      summary.ActiveOpen,
				"funded":    summary.ActiveFunded,
				"submitted": summary.ActiveSubmitted,
			},
			LastUpdated: summary.LastUpdated,
		},
		Distributions: types.CommerceJobsGeneralDistributionsResp{
			Token:         tokenDist,
			ChainContracts: grouped,
			Fees:          fees,
		},
	}, nil
}

func groupChainContracts(flat []model.CommerceJobsChainContractDistributionItem) []types.CommerceJobsChainContractsItem {
	if len(flat) == 0 {
		return []types.CommerceJobsChainContractsItem{}
	}

	byChain := make(map[string]*types.CommerceJobsChainContractsItem)
	chainOrder := make([]string, 0, 8)

	for _, it := range flat {
		cid := strings.TrimSpace(it.ChainID)
		if cid == "" {
			continue
		}
		g, ok := byChain[cid]
		if !ok {
			g = &types.CommerceJobsChainContractsItem{
				ChainID:          cid,
				ChainName:        it.ChainName,
				ChainLogo:        it.ChainLogo,
				ERC8183Contracts: []types.CommerceJobsChainContractItem{},
			}
			byChain[cid] = g
			chainOrder = append(chainOrder, cid)
		} else {
			// Ensure chain meta is filled when later rows contain it.
			if g.ChainName == "" && it.ChainName != "" {
				g.ChainName = it.ChainName
			}
			if g.ChainLogo == "" && it.ChainLogo != "" {
				g.ChainLogo = it.ChainLogo
			}
		}

		contract := strings.TrimSpace(it.CommerceContract)
		if contract == "" {
			continue
		}
		g.ERC8183Contracts = append(g.ERC8183Contracts, types.CommerceJobsChainContractItem{
			CommerceContract: contract,
			JobsCount:        it.JobsCount,
			BudgetVolumeUSD:  it.BudgetVolumeUSD,
			PaidVolumeUSD:    it.PaidVolumeUSD,
			Drilldown: map[string]any{
				"chain_id":          cid,
				"commerce_contract": contract,
			},
		})
	}

	out := make([]types.CommerceJobsChainContractsItem, 0, len(chainOrder))
	for _, cid := range chainOrder {
		g := byChain[cid]
		if g == nil || len(g.ERC8183Contracts) == 0 {
			continue
		}
		// Deterministic order: paid desc then jobs desc.
		sort.Slice(g.ERC8183Contracts, func(i, j int) bool {
			if g.ERC8183Contracts[i].PaidVolumeUSD != g.ERC8183Contracts[j].PaidVolumeUSD {
				return g.ERC8183Contracts[i].PaidVolumeUSD > g.ERC8183Contracts[j].PaidVolumeUSD
			}
			if g.ERC8183Contracts[i].JobsCount != g.ERC8183Contracts[j].JobsCount {
				return g.ERC8183Contracts[i].JobsCount > g.ERC8183Contracts[j].JobsCount
			}
			return g.ERC8183Contracts[i].CommerceContract < g.ERC8183Contracts[j].CommerceContract
		})
		out = append(out, *g)
	}
	return out
}

func GetCommerceJobsCharts(p CommerceJobsChartsParams) (any, error) {
	window, err := normalizeChartsWindow(p.Window)
	if err != nil {
		return nil, err
	}
	bucket := p.BucketSeconds
	if bucket == 0 {
		bucket = defaultBucketSeconds(window)
	}
	if window == "custom" && bucket == 0 {
		return nil, errors.New("missing bucket_seconds for custom window")
	}
	if bucket == 0 {
		return nil, errors.New("invalid bucket_seconds")
	}
	ws, we := uint64(0), uint64(0)
	if p.StartTime != nil {
		ws = *p.StartTime
	}
	if p.EndTime != nil {
		we = *p.EndTime
	}
	q := model.CommerceJobsQuery{
		ChainID:       p.ChainID,
		Contract:      p.CommerceContract,
		Status:        p.Status,
		Role:          p.Role,
		AgentAddress:  p.AgentAddress,
		Counterparty:  p.Counterparty,
		PaymentToken:  p.PaymentToken,
		TokenSymbol:   p.TokenSymbol,
		MinBudget:     p.MinBudget,
		MaxBudget:     p.MaxBudget,
		MinBudgetUSD:  p.MinBudgetUSD,
		MaxBudgetUSD:  p.MaxBudgetUSD,
		StartTime:     p.StartTime,
		EndTime:       p.EndTime,
		Page:          1,
		PageSize:      1,
	}
	charts, err := model.GetCommerceJobsCharts(q, ws, we, bucket)
	if err != nil {
		return nil, err
	}
	// Enrich chain meta for UI display (chain name/logo).
	// NOTE: chart chain/contract distribution is grouped in model (see model.CommerceJobsCharts).
	for i := range charts.Distribution.ChainContracts {
		if charts.Distribution.ChainContracts[i].ChainID == "" {
			continue
		}
		if chain, ok := config.GetChainInfo(charts.Distribution.ChainContracts[i].ChainID); ok {
			charts.Distribution.ChainContracts[i].ChainName = chain.ChainName
			charts.Distribution.ChainContracts[i].ChainLogo = chain.ChainLogo
		}
	}
	return charts, nil
}

// BuildCommerceJobsFilters builds default filter options for Job Browser.
// It is used by the in-memory cache refresher.
func BuildCommerceJobsFilters(ctx context.Context) (*types.CommerceJobsFilters, error) {
	_ = ctx // reserved for future DB timeouts/cancellation

	chainIDs, err := model.GetCommerceJobsDistinctChainIDs()
	if err != nil {
		return nil, fmt.Errorf("distinct chain_ids: %w", err)
	}
	contracts, err := model.GetCommerceJobsDistinctCommerceContracts()
	if err != nil {
		return nil, fmt.Errorf("distinct commerce_contracts: %w", err)
	}
	tokens, err := model.GetCommerceJobsDistinctPaymentTokens()
	if err != nil {
		return nil, fmt.Errorf("distinct payment_tokens: %w", err)
	}

	chains := make([]types.CommerceJobsFilterChain, 0, len(chainIDs))
	for _, id := range chainIDs {
		item := types.CommerceJobsFilterChain{ChainID: id}
		if chain, ok := config.GetChainInfo(id); ok {
			item.ChainName = chain.ChainName
			item.ChainLogo = chain.ChainLogo
		}
		chains = append(chains, item)
	}

	// Stable order for UI.
	sort.Slice(chains, func(i, j int) bool { return chains[i].ChainID < chains[j].ChainID })

	return &types.CommerceJobsFilters{
		Chains:            chains,
		CommerceContracts: contracts,
		PaymentTokens:     tokens,
		LastUpdated:       uint64(time.Now().Unix()),
	}, nil
}

// ─────────────── Commerce Job Actions (Job Detail) ───────────────

type CommerceJobActionsParams struct {
	ChainID          string
	CommerceContract string
	JobID            uint64

	ActionTypes []string
	StartTime   *uint64
	EndTime     *uint64

	Page     int
	PageSize int
}

func normalizeActionTypes(v []string) []string {
	if len(v) == 0 {
		return nil
	}
	out := make([]string, 0, len(v))
	for _, s := range v {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func actionToUIType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "job_created":
		return "JobCreated"
	case "job_funded":
		return "JobFunded"
	case "job_submitted":
		return "JobSubmitted"
	case "job_completed":
		return "JobCompleted"
	case "job_rejected":
		return "JobRejected"
	case "job_expired":
		return "JobExpired"
	case "provider_set":
		return "ProviderSet"
	case "budget_set":
		return "BudgetSet"
	case "payment_released":
		return "PaymentReleased"
	case "platform_fee_paid":
		return "PlatformFeePaid"
	case "evaluator_fee_paid":
		return "EvaluatorFeePaid"
	default:
		return raw
	}
}

func roleNormalize(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	switch v {
	case "client", "provider", "evaluator", "platform", "unknown":
		return v
	default:
		return "unknown"
	}
}

func GetCommerceJobActions(p CommerceJobActionsParams) ([]types.CommerceJobActionDTO, int64, error) {
	if p.ChainID == "" || p.CommerceContract == "" || p.JobID == 0 {
		return nil, 0, errors.New("missing chain_id/commerce_contract/job_id")
	}

	page := normalizePage(p.Page)
	pageSize := normalizePageSize(p.PageSize)
	actionTypes := normalizeActionTypes(p.ActionTypes)

	actions, total, err := model.GetCommerceJobActions(model.CommerceJobActionsQuery{
		ChainID:     p.ChainID,
		Contract:    p.CommerceContract,
		JobID:       p.JobID,
		ActionTypes: actionTypes,
		StartTime:   p.StartTime,
		EndTime:     p.EndTime,
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("get commerce job actions: %w", err)
	}

	out := make([]types.CommerceJobActionDTO, 0, len(actions))
	for _, a := range actions {
		var amount *float64
		var amountUSD *float64
		if a.JobBudget > 0 {
			v := a.JobBudget
			amount = &v
		}
		if a.BudgetUSD > 0 {
			v := a.BudgetUSD
			amountUSD = &v
		}
		out = append(out, types.CommerceJobActionDTO{
			ChainID:          a.ChainID,
			CommerceContract: a.CommerceContract,
			JobID:            a.JobID,
			ActionType:       actionToUIType(a.Action),
			BlockTimestamp:   a.BlockTimestamp,
			BlockNumber:      a.BlockNumber,
			TxHash:           a.TxHash,
			LogIndex:         a.LogIndex,
			Actor:            a.AgentAddress,
			Role:             roleNormalize(a.Role),
			PaymentToken:     a.PaymentToken,
			PaymentDecimals:  a.PaymentDecimals,
			TokenSymbol:      a.TokenSymbol,
			Amount:           amount,
			AmountUSD:        amountUSD,
		})
	}
	return out, total, nil
}
