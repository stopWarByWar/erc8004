package mock

import (
	serverLogic "agent_identity/server/api/logic"
	"agent_identity/server/api/types"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Mock chain meta should mirror config/config.yaml (data replacement, not runtime lookup).
const (
	mockChainLogoEth  = "https://bnbattest.s3.ap-southeast-1.amazonaws.com/agentCard/chain/eth.png"
	mockChainLogoBsc  = "https://bnbattest.s3.ap-southeast-1.amazonaws.com/agentCard/chain/bsc.png"
	mockChainLogoBase = "https://bnbattest.s3.ap-southeast-1.amazonaws.com/agentCard/chain/base.png"
)

func mockChainMeta(chainID string) (name string, logo string) {
	switch strings.TrimSpace(chainID) {
	case "1":
		return "Ethereum", mockChainLogoEth
	case "56":
		return "BSC", mockChainLogoBsc
	case "8453":
		return "Base", mockChainLogoBase
	default:
		return "", ""
	}
}

func seededRand(parts ...string) *rand.Rand {
	h := fnv.New64a()
	for _, p := range parts {
		_, _ = h.Write([]byte(p))
		_, _ = h.Write([]byte{0})
	}
	seed := int64(binary.LittleEndian.Uint64(h.Sum(nil)))
	return rand.New(rand.NewSource(seed))
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func paginate[T any](items []T, page, pageSize int) ([]T, int64) {
	page = clampInt(page, 1, math.MaxInt32)
	pageSize = clampInt(pageSize, 1, 100)
	total := int64(len(items))
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []T{}, total
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	out := make([]T, 0, end-start)
	out = append(out, items[start:end]...)
	return out, total
}

// -------------------- agent/identity/commerce/scores --------------------

func CommerceScores(uid uint64, chainID, contract string) []serverLogic.CommerceScoreResp {
	r := seededRand(fmt.Sprintf("%d", uid), chainID, contract, "scores")
	roles := []string{"client", "provider", "evaluator"}
	out := make([]serverLogic.CommerceScoreResp, 0, len(roles))
	for _, role := range roles {
		totalJobs := 20 + r.Intn(180)
		completed := r.Intn(totalJobs + 1)
		rejected := r.Intn(totalJobs - completed + 1)
		expired := r.Intn(10)
		successRate := 0.0
		if completed+rejected > 0 {
			successRate = float64(completed) / float64(completed+rejected)
		}
		totalVolumeUSD := 1000 + r.Float64()*50000
		weightedScoreUSD := 0.2 + r.Float64()*0.8
		out = append(out, serverLogic.CommerceScoreResp{
			Role:                    role,
			CompletedCount:          completed,
			RejectedCount:           rejected,
			ExpiredResponsibleCount: expired,
			SuccessRate:             round3(successRate),
			WeightedScore:           round3(0.2 + r.Float64()*0.8),
			TotalVolumeUSD:          round2(totalVolumeUSD),
			WeightedScoreUSD:        round3(weightedScoreUSD),
			CreatedCount:            5 + r.Intn(80),
			FundedCount:             5 + r.Intn(80),
			FundedRate:              round3(r.Float64()),
			CompletionRate:          round3(r.Float64()),
			EvaluatedCount:          r.Intn(30),
			ExpiredFromSubmittedCount: r.Intn(10),
			Responsiveness:            round3(r.Float64()),
			TotalJobs:                 totalJobs,
			TotalVolume:               round2(10 + r.Float64()*300),
			UniqueCounterparties:      1 + r.Intn(50),
			Confidence:                round3(0.4 + r.Float64()*0.6),
		})
	}
	return out
}

// -------------------- agent/identity/commerce/actions --------------------

type CommerceActionsQuery struct {
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

func CommerceActions(q CommerceActionsQuery) ([]types.CommerceActionDTO, int64) {
	r := seededRand(fmt.Sprintf("%d", q.UID), q.Role, q.ChainID, q.Contract, "actions")
	all := make([]types.CommerceActionDTO, 0, 120)
	actionTypes := []string{"job_created", "job_funded", "job_submitted", "job_completed", "job_rejected", "job_expired", "provider_set", "budget_set"}
	polarities := []string{"positive", "neutral", "negative"}
	certainties := []string{"definitive", "indicative"}
	now := uint64(time.Now().Unix())

	for i := 0; i < 120; i++ {
		act := actionTypes[r.Intn(len(actionTypes))]
		role := []string{"client", "provider", "evaluator"}[r.Intn(3)]
		pol := polarities[r.Intn(len(polarities))]
		cert := certainties[r.Intn(len(certainties))]
		ts := now - uint64(r.Intn(3600*24*30))
		budget := 10 + r.Float64()*500
		budgetUSD := budget * (0.8 + r.Float64()*1.2)
		jobID := uint64(1000 + r.Intn(200))
		agentAddr := fmt.Sprintf("0x%040x", r.Uint64())
		counterparty := fmt.Sprintf("0x%040x", r.Uint64())
		contract := q.Contract
		if contract == "" {
			contract = fmt.Sprintf("0x%040x", r.Uint64())
		}
		chainID := q.ChainID
		if chainID == "" {
			chainID = "1"
		}

		all = append(all, types.CommerceActionDTO{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            jobID,
			AgentUID:         q.UID,
			AgentAddress:     agentAddr,
			Role:             role,
			Action:           act,
			SignalPolarity:   pol,
			SignalWeight:     round3(0.1 + r.Float64()*0.9),
			SignalCertainty:  cert,
			JobBudget:        formatAmount8(budget),
			BudgetUSD:        round2(budgetUSD),
			PaymentToken:     firstNonEmpty(q.PaymentToken, fmt.Sprintf("0x%040x", r.Uint64())),
			PaymentDecimals:  18,
			TokenSymbol:      firstNonEmpty(q.TokenSymbol, []string{"USDC", "DAI", "ETH"}[r.Intn(3)]),
			Counterparty:     counterparty,
			Reason:           []string{"", "low_quality", "timeout", "scope_change"}[r.Intn(4)],
			Deliverable:      []string{"", "ipfs://Qm...", "https://example.com/result"}[r.Intn(3)],
			PreviousStatus:   []string{"", "open", "funded", "submitted"}[r.Intn(4)],
			HookAddress:      hookAddrOrEmpty(r, q.HookAddress),
			BlockNumber:      uint64(18000000 + r.Intn(100000)),
			TxHash:           fmt.Sprintf("0x%064x", r.Uint64()),
			LogIndex:         uint(r.Intn(10)),
			BlockTimestamp:   ts,
		})
	}

	filtered := filterCommerceActions(all, q)

	// sort
	if strings.ToLower(strings.TrimSpace(q.SortBy)) == "budget" {
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].BudgetUSD > filtered[j].BudgetUSD })
	} else if strings.ToLower(strings.TrimSpace(q.SortBy)) == "signal_weight" {
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].SignalWeight > filtered[j].SignalWeight })
	} else {
		// timestamp desc (default)
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].BlockTimestamp > filtered[j].BlockTimestamp })
	}

	pageItems, total := paginate(filtered, q.Page, q.PageSize)
	return pageItems, total
}

func filterCommerceActions(items []types.CommerceActionDTO, q CommerceActionsQuery) []types.CommerceActionDTO {
	actionSet := make(map[string]struct{}, len(q.Actions))
	for _, a := range q.Actions {
		actionSet[strings.ToLower(strings.TrimSpace(a))] = struct{}{}
	}
	polSet := make(map[string]struct{}, len(q.Polarities))
	for _, p := range q.Polarities {
		polSet[strings.ToLower(strings.TrimSpace(p))] = struct{}{}
	}
	role := strings.ToLower(strings.TrimSpace(q.Role))
	certainty := strings.ToLower(strings.TrimSpace(q.Certainty))

	out := make([]types.CommerceActionDTO, 0, len(items))
	for _, it := range items {
		if role != "" && strings.ToLower(it.Role) != role {
			continue
		}
		if q.ChainID != "" && it.ChainID != q.ChainID {
			continue
		}
		if q.Contract != "" && it.CommerceContract != q.Contract {
			continue
		}
		if q.Counterparty != "" && !strings.EqualFold(it.Counterparty, q.Counterparty) {
			continue
		}
		if certainty != "" && strings.ToLower(it.SignalCertainty) != certainty {
			continue
		}
		if len(actionSet) > 0 {
			if _, ok := actionSet[strings.ToLower(it.Action)]; !ok {
				continue
			}
		}
		if len(polSet) > 0 {
			if _, ok := polSet[strings.ToLower(it.SignalPolarity)]; !ok {
				continue
			}
		}
		if q.HasHook != nil {
			has := strings.TrimSpace(it.HookAddress) != ""
			if has != *q.HasHook {
				continue
			}
		}
		if q.PaymentToken != "" && !strings.EqualFold(it.PaymentToken, q.PaymentToken) {
			continue
		}
		if q.TokenSymbol != "" && !strings.EqualFold(it.TokenSymbol, q.TokenSymbol) {
			continue
		}
		if q.MinBudget != nil || q.MaxBudget != nil {
			amt := parseAmountOrZero(it.JobBudget)
			if q.MinBudget != nil && amt < *q.MinBudget {
				continue
			}
			if q.MaxBudget != nil && amt > *q.MaxBudget {
				continue
			}
		}
		if q.MinBudgetUSD != nil && it.BudgetUSD < *q.MinBudgetUSD {
			continue
		}
		if q.MaxBudgetUSD != nil && it.BudgetUSD > *q.MaxBudgetUSD {
			continue
		}
		if q.StartTime != nil && it.BlockTimestamp < *q.StartTime {
			continue
		}
		if q.EndTime != nil && it.BlockTimestamp > *q.EndTime {
			continue
		}
		if q.HookAddress != "" && !strings.EqualFold(it.HookAddress, q.HookAddress) {
			continue
		}

		out = append(out, it)
	}
	return out
}

// -------------------- agent/identity/commerce/stats --------------------

func CommerceStats(uid uint64, now uint64) *types.CommerceStats {
	r := seededRand(fmt.Sprintf("%d", uid), fmt.Sprintf("%d", now), "stats")
	stats := &types.CommerceStats{
		ActionBreakdown:    make(types.ActionBreakdown),
		TimeSeries:         types.TimeSeriesStats{},
		BudgetDistribution: make(types.BudgetDistribution),
	}

	roles := []string{"client", "provider", "evaluator"}
	for _, role := range roles {
		ac := types.ActionCounts{
			JobCreated:   r.Intn(200),
			JobFunded:    r.Intn(200),
			JobSubmitted: r.Intn(200),
			JobCompleted: r.Intn(200),
			JobRejected:  r.Intn(200),
			JobExpired:   r.Intn(200),
		}
		// Keep only non-empty roles like real logic.
		if ac.JobCreated+ac.JobFunded+ac.JobSubmitted+ac.JobCompleted+ac.JobRejected+ac.JobExpired == 0 {
			continue
		}
		stats.ActionBreakdown[role] = ac
	}

	stats.TimeSeries.Hours24 = timeBuckets(r, now, 24, 3600)
	stats.TimeSeries.Days7 = timeBuckets(r, now, 28, 21600)
	stats.TimeSeries.Days30 = timeBuckets(r, now, 30, 86400)

	contracts := []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
	}
	for _, role := range roles {
		stats.BudgetDistribution[role] = make(map[string]types.ContractBudgetStats)
		for _, c := range contracts {
			total := 10 + r.Intn(200)
			stats.BudgetDistribution[role][c] = types.ContractBudgetStats{
				TotalCount: total,
				Buckets: types.BudgetBuckets{
					Small:  &types.BudgetBucketStat{Count: total / 3, MaxAmount: formatAmount8(10 + r.Float64()*50)},
					Medium: &types.BudgetBucketStat{Count: total / 3, MaxAmount: formatAmount8(100 + r.Float64()*200)},
					Large:  &types.BudgetBucketStat{Count: total - 2*(total/3), MaxAmount: formatAmount8(500 + r.Float64()*1500)},
				},
			}
		}
	}

	// Remove empty roles.
	for role, contractsMap := range stats.BudgetDistribution {
		if len(contractsMap) == 0 {
			delete(stats.BudgetDistribution, role)
		}
	}

	return stats
}

func timeBuckets(r *rand.Rand, now uint64, n int, step uint64) []types.TimeBucket {
	out := make([]types.TimeBucket, 0, n)
	for i := 0; i < n; i++ {
		t := int64(now - uint64((n-1-i))*step)
		completed := r.Intn(20)
		rejected := r.Intn(10)
		if completed == 0 && rejected == 0 {
			continue
		}
		sr := 0.0
		if completed+rejected > 0 {
			sr = float64(completed) / float64(completed+rejected)
		}
		out = append(out, types.TimeBucket{
			Bucket:         t,
			CompletedCount: completed,
			RejectedCount:  rejected,
			SuccessRate:    round3(sr),
		})
	}
	return out
}

// -------------------- commerce/jobs --------------------

type CommerceJobsQuery struct {
	ChainID          string
	CommerceContract string
	Status           string
	Role             string
	AgentAddress     string
	Counterparty     string
	PaymentToken     string
	TokenSymbol      string
	MinBudget        *float64
	MaxBudget        *float64
	MinBudgetUSD     *float64
	MaxBudgetUSD     *float64
	StartTime        *uint64
	EndTime          *uint64
	SortBy           string
	SortOrder        string
	Page             int
	PageSize         int
}

func CommerceJobs(q CommerceJobsQuery) ([]types.CommerceJobDTO, int64) {
	r := seededRand(q.ChainID, q.CommerceContract, q.Status, q.Role, q.AgentAddress, q.Counterparty, "jobs")
	now := uint64(time.Now().Unix())
	statuses := []string{"open", "funded", "submitted", "completed", "rejected", "expired"}
	all := make([]types.CommerceJobDTO, 0, 200)
	for i := 0; i < 200; i++ {
		st := statuses[r.Intn(len(statuses))]
		chainID := firstNonEmpty(q.ChainID, "1")
		chainName, chainLogo := mockChainMeta(chainID)
		contract := firstNonEmpty(q.CommerceContract, fmt.Sprintf("0x%040x", r.Uint64()))
		jobID := uint64(10000 + i)
		updatedAt := now - uint64(r.Intn(3600*24*90))
		paymentToken := firstNonEmpty(q.PaymentToken, fmt.Sprintf("0x%040x", r.Uint64()))
		tokenSymbol := firstNonEmpty(q.TokenSymbol, []string{"USDC", "DAI", "ETH"}[r.Intn(3)])
		budget := 20 + r.Float64()*2000
		budgetUSD := budget * (0.8 + r.Float64()*1.3)
		paid := budget * (0.1 + r.Float64()*0.9)
		paidUSD := paid * (0.8 + r.Float64()*1.3)
		platformFee := paid * 0.02
		evaluatorFee := paid * 0.01

		all = append(all, types.CommerceJobDTO{
			JobUID:           900000 + uint64(i),
			ChainID:          chainID,
			ChainName:        chainName,
			ChainLogo:        chainLogo,
			CommerceContract: contract,
			JobID:            jobID,
			Status:           st,
			UpdatedAt:        updatedAt,
			Client:           firstNonEmpty(q.AgentAddress, fmt.Sprintf("0x%040x", r.Uint64())),
			Provider:         fmt.Sprintf("0x%040x", r.Uint64()),
			Evaluator:        fmt.Sprintf("0x%040x", r.Uint64()),
			Description:      fmt.Sprintf("Mock job #%d: build feature %d", jobID, r.Intn(1000)),
			HookAddress:      hookAddrOrEmpty(r, ""),
			ExpiredAt:        updatedAt + uint64(r.Intn(3600*24*7)),
			SubmittedAt:      updatedAt - uint64(r.Intn(3600*24*2)),
			CompletedAt:      updatedAt + uint64(r.Intn(3600*24*2)),
			PaymentToken:     paymentToken,
			PaymentDecimals:  18,
			TokenSymbol:      tokenSymbol,
			Budget:           formatAmount8(budget),
			BudgetUSD:        round2(budgetUSD),
			PaidAmount:       formatAmount8(paid),
			PaidAmountUSD:    round2(paidUSD),
			PlatformFeeAmount:  formatAmount8(platformFee),
			PlatformFeeUSD:     round2(platformFee * (0.8 + r.Float64()*1.3)),
			EvaluatorFeeAmount: formatAmount8(evaluatorFee),
			EvaluatorFeeUSD:    round2(evaluatorFee * (0.8 + r.Float64()*1.3)),
			LatestBlockNumber:  uint64(18000000 + r.Intn(100000)),
			LatestTxHash:       fmt.Sprintf("0x%064x", r.Uint64()),
			LatestActionUID:    uint64(1 + r.Intn(100000)),
		})
	}

	filtered := filterJobs(all, q)

	// sorting
	sortBy := strings.ToLower(strings.TrimSpace(q.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(q.SortOrder))
	desc := sortOrder != "asc"
	lessU64 := func(a, b uint64) bool { return a < b }
	lessF64 := func(a, b float64) bool { return a < b }
	switch sortBy {
	case "budget":
		sort.Slice(filtered, func(i, j int) bool {
			if desc {
				return !lessF64(parseAmountOrZero(filtered[i].Budget), parseAmountOrZero(filtered[j].Budget))
			}
			return lessF64(parseAmountOrZero(filtered[i].Budget), parseAmountOrZero(filtered[j].Budget))
		})
	case "budget_usd":
		sort.Slice(filtered, func(i, j int) bool {
			if desc {
				return filtered[i].BudgetUSD > filtered[j].BudgetUSD
			}
			return filtered[i].BudgetUSD < filtered[j].BudgetUSD
		})
	case "paid_amount_usd":
		sort.Slice(filtered, func(i, j int) bool {
			if desc {
				return filtered[i].PaidAmountUSD > filtered[j].PaidAmountUSD
			}
			return filtered[i].PaidAmountUSD < filtered[j].PaidAmountUSD
		})
	default: // updated_at
		sort.Slice(filtered, func(i, j int) bool {
			if desc {
				return !lessU64(filtered[i].UpdatedAt, filtered[j].UpdatedAt)
			}
			return lessU64(filtered[i].UpdatedAt, filtered[j].UpdatedAt)
		})
	}

	pageItems, total := paginate(filtered, q.Page, q.PageSize)
	return pageItems, total
}

func filterJobs(items []types.CommerceJobDTO, q CommerceJobsQuery) []types.CommerceJobDTO {
	status := strings.ToLower(strings.TrimSpace(q.Status))
	role := strings.ToLower(strings.TrimSpace(q.Role))
	out := make([]types.CommerceJobDTO, 0, len(items))
	for _, it := range items {
		if q.ChainID != "" && it.ChainID != q.ChainID {
			continue
		}
		if q.CommerceContract != "" && it.CommerceContract != q.CommerceContract {
			continue
		}
		if status != "" && strings.ToLower(it.Status) != status {
			continue
		}
		// role: we treat as "my role" by comparing address.
		if role != "" && q.AgentAddress != "" {
			addr := strings.ToLower(q.AgentAddress)
			switch role {
			case "client":
				if strings.ToLower(it.Client) != addr {
					continue
				}
			case "provider":
				if strings.ToLower(it.Provider) != addr {
					continue
				}
			case "evaluator":
				if strings.ToLower(it.Evaluator) != addr {
					continue
				}
			}
		}
		if q.AgentAddress != "" {
			// If caller passes agent_address but no role, keep jobs where address appears in any party.
			addr := strings.ToLower(q.AgentAddress)
			if strings.ToLower(it.Client) != addr && strings.ToLower(it.Provider) != addr && strings.ToLower(it.Evaluator) != addr {
				continue
			}
		}
		if q.Counterparty != "" {
			// Minimal: match any party.
			cp := strings.ToLower(q.Counterparty)
			if strings.ToLower(it.Client) != cp && strings.ToLower(it.Provider) != cp && strings.ToLower(it.Evaluator) != cp {
				continue
			}
		}
		if q.PaymentToken != "" && !strings.EqualFold(it.PaymentToken, q.PaymentToken) {
			continue
		}
		if q.TokenSymbol != "" && !strings.EqualFold(it.TokenSymbol, q.TokenSymbol) {
			continue
		}
		if q.MinBudget != nil || q.MaxBudget != nil {
			amt := parseAmountOrZero(it.Budget)
			if q.MinBudget != nil && amt < *q.MinBudget {
				continue
			}
			if q.MaxBudget != nil && amt > *q.MaxBudget {
				continue
			}
		}
		if q.MinBudgetUSD != nil && it.BudgetUSD < *q.MinBudgetUSD {
			continue
		}
		if q.MaxBudgetUSD != nil && it.BudgetUSD > *q.MaxBudgetUSD {
			continue
		}
		if q.StartTime != nil && it.UpdatedAt < *q.StartTime {
			continue
		}
		if q.EndTime != nil && it.UpdatedAt > *q.EndTime {
			continue
		}
		out = append(out, it)
	}
	return out
}

func CommerceJobDetail(chainID, contract string, jobID uint64) (*types.CommerceJobDTO, *types.CommerceJobDetailEvidence) {
	r := seededRand(chainID, contract, fmt.Sprintf("%d", jobID), "job_detail")
	now := uint64(time.Now().Unix())
	chainName, chainLogo := mockChainMeta(firstNonEmpty(chainID, "1"))
	job := &types.CommerceJobDTO{
		JobUID:           900000 + jobID,
		ChainID:          firstNonEmpty(chainID, "1"),
		ChainName:        chainName,
		ChainLogo:        chainLogo,
		CommerceContract: firstNonEmpty(contract, fmt.Sprintf("0x%040x", r.Uint64())),
		JobID:            jobID,
		Status:           []string{"open", "funded", "submitted", "completed"}[r.Intn(4)],
		UpdatedAt:        now - uint64(r.Intn(3600*24*60)),
		Client:           fmt.Sprintf("0x%040x", r.Uint64()),
		Provider:         fmt.Sprintf("0x%040x", r.Uint64()),
		Evaluator:        fmt.Sprintf("0x%040x", r.Uint64()),
		Description:      fmt.Sprintf("Mock job #%d detail", jobID),
		HookAddress:      hookAddrOrEmpty(r, ""),
		PaymentToken:     fmt.Sprintf("0x%040x", r.Uint64()),
		PaymentDecimals:  18,
		TokenSymbol:      []string{"USDC", "DAI", "ETH"}[r.Intn(3)],
		Budget:           formatAmount8(100 + r.Float64()*2000),
		BudgetUSD:        round2(100 + r.Float64()*5000),
		PaidAmount:       formatAmount8(20 + r.Float64()*1800),
		PaidAmountUSD:    round2(20 + r.Float64()*4000),
		PlatformFeeAmount:  formatAmount8(1 + r.Float64()*50),
		PlatformFeeUSD:     round2(1 + r.Float64()*80),
		EvaluatorFeeAmount: formatAmount8(1 + r.Float64()*30),
		EvaluatorFeeUSD:    round2(1 + r.Float64()*50),
		LatestBlockNumber:  uint64(18000000 + r.Intn(100000)),
		LatestTxHash:       fmt.Sprintf("0x%064x", r.Uint64()),
		LatestActionUID:    uint64(1 + r.Intn(100000)),
	}
	ev := &types.CommerceJobDetailEvidence{
		TimelineSource:    "commerce_actions",
		EventsTableSource: "commerce_actions",
		SettlementEvents:  []string{"PaymentReleased", "PlatformFeePaid", "EvaluatorFeePaid"},
	}
	return job, ev
}

func CommerceJobsGeneral() (types.CommerceJobsGeneralSummaryResp, types.CommerceJobsGeneralDistributionsResp) {
	// Keep schema-compatible but minimal.
	chainName, chainLogo := mockChainMeta("8453")
	return types.CommerceJobsGeneralSummaryResp{
			JobsCount:       1234,
			PaidVolumeUSD:   567890.12,
			BudgetVolumeUSD: 901234.56,
			OutcomeMix:      map[string]int64{"completed": 700, "rejected": 100, "expired": 50},
			ActiveMix:       map[string]int64{"open": 200, "funded": 120, "submitted": 64},
			LastUpdated:     uint64(time.Now().Unix()),
		}, types.CommerceJobsGeneralDistributionsResp{
			Token: map[string]any{
				"USDC": map[string]any{"jobs_count": 800, "paid_volume_usd": 400000.0},
				"DAI":  map[string]any{"jobs_count": 300, "paid_volume_usd": 120000.0},
				"ETH":  map[string]any{"jobs_count": 134, "paid_volume_usd": 47890.12},
			},
			ChainContracts: []types.CommerceJobsChainContractsItem{
				{
					ChainID:   "8453",
					ChainName: chainName,
					ChainLogo: chainLogo,
					ERC8183Contracts: []types.CommerceJobsChainContractItem{
						{
							CommerceContract: "0x1111111111111111111111111111111111111111",
							JobsCount:        600,
							BudgetVolumeUSD:  500000.0,
							PaidVolumeUSD:    300000.0,
							Drilldown:        map[string]any{"chain_id": "8453", "commerce_contract": "0x1111111111111111111111111111111111111111"},
						},
					},
				},
			},
			Fees: map[string]any{
				"platform_fee_usd": 12345.67,
				"evaluator_fee_usd": 4567.89,
			},
		}
}

func CommerceJobsCharts() any {
	now := uint64(time.Now().Unix())
	r := seededRand("charts", fmt.Sprintf("%d", now))

	// Build buckets (7d, 6h) as default-like.
	buckets := make([]map[string]any, 0, 10)
	step := uint64(6 * 3600)
	start := now - step*uint64(9)
	for i := 0; i < 10; i++ {
		b := start + uint64(i)*step
		count := int64(5 + r.Intn(30))
		paid := round2(500 + r.Float64()*5000)
		buckets = append(buckets, map[string]any{
			"bucket":    b,
			"count":     count,
			"drilldown": map[string]any{"start_time": b, "end_time": b + step},
		})
		_ = paid
	}

	paidBuckets := make([]map[string]any, 0, len(buckets))
	feeBuckets := make([]map[string]any, 0, len(buckets))
	for _, it := range buckets {
		b := it["bucket"].(uint64)
		paid := round2(500 + r.Float64()*8000)
		fee := round2(paid * 0.02)
		paidBuckets = append(paidBuckets, map[string]any{
			"bucket":         b,
			"value_usd":      paid,
			"paid_volume_usd": paid, // keep compatibility
			"drilldown":      map[string]any{"status": "completed", "start_time": b, "end_time": b + step},
		})
		feeBuckets = append(feeBuckets, map[string]any{
			"bucket":           b,
			"platform_fee_usd": fee,
			"value_usd":        fee,
			"drilldown":        map[string]any{"status": "completed", "start_time": b, "end_time": b + step},
		})
	}

	tokenDist := []map[string]any{
		{"token_symbol": "USDC", "jobs_count": 80, "budget_volume_usd": 60000, "paid_volume_usd": 40000, "drilldown": map[string]any{"token_symbol": "USDC"}},
		{"token_symbol": "ETH", "jobs_count": 40, "budget_volume_usd": 20123.45, "paid_volume_usd": 12340.12, "drilldown": map[string]any{"token_symbol": "ETH"}},
	}
	chainName, chainLogo := mockChainMeta("8453")
	chainContracts := []map[string]any{
		{
			"chain_id":   "8453",
			"chain_name": chainName,
			"chain_logo": chainLogo,
			"erc8183_contracts": []map[string]any{
				{
					"commerce_contract": "0x1111111111111111111111111111111111111111",
					"jobs_count":        120,
					"budget_volume_usd": 80123.45,
					"paid_volume_usd":   52340.12,
					"drilldown":         map[string]any{"chain_id": "8453", "commerce_contract": "0x1111111111111111111111111111111111111111"},
				},
			},
		},
	}

	return map[string]any{
		"activity": map[string]any{
			"active_jobs_over_time": buckets,
			"created_jobs_over_time": []map[string]any{
				{"bucket": start, "count": 3, "drilldown": map[string]any{"start_time": start, "end_time": start + step}},
			},
			"unique_clients_over_time": []map[string]any{
				{"bucket": start, "count": 4, "drilldown": map[string]any{"start_time": start, "end_time": start + step}},
			},
			"unique_providers_over_time": []map[string]any{
				{"bucket": start, "count": 5, "drilldown": map[string]any{"start_time": start, "end_time": start + step}},
			},
		},
		"health": map[string]any{
			"outcome_mix": []map[string]any{
				{"status": "completed", "count": 50, "drilldown": map[string]any{"status": "completed"}},
				{"status": "rejected", "count": 10, "drilldown": map[string]any{"status": "rejected"}},
				{"status": "expired", "count": 5, "drilldown": map[string]any{"status": "expired"}},
			},
			"funnel": []map[string]any{
				{"stage": "open", "count": 30, "drilldown": map[string]any{"status": "open"}},
				{"stage": "funded", "count": 15, "drilldown": map[string]any{"status": "funded"}},
				{"stage": "submitted", "count": 10, "drilldown": map[string]any{"status": "submitted"}},
				{"stage": "completed", "count": 50, "drilldown": map[string]any{"status": "completed"}},
			},
		},
		"volume": map[string]any{
			"paid_volume_usd_over_time":     paidBuckets,
			"platform_fee_usd_over_time":    feeBuckets,
			"fees_breakdown": []map[string]any{
				{"type": "platform_fee_usd", "value": 1200.5, "drilldown": map[string]any{"status": "completed"}},
				{"type": "evaluator_fee_usd", "value": 800.25, "drilldown": map[string]any{"status": "completed"}},
			},
		},
		"distribution": map[string]any{
			"token_distribution":          tokenDist,
			"chain_contracts":             chainContracts,
			"budget_histogram_usd": []map[string]any{
				{"range": map[string]any{"min": 0, "max": 100}, "count": 20, "drilldown": map[string]any{"min_budget_usd": 0, "max_budget_usd": 100}},
			},
		},
		"evidence": map[string]any{
			"top_events": []map[string]any{
				{
					"kind":             "PaymentReleased",
					"chain_id":          "8453",
					"commerce_contract": "0x1111111111111111111111111111111111111111",
					"job_id":            12,
					"value_usd":         1350,
					"drilldown":         map[string]any{"chain_id": "8453", "commerce_contract": "0x1111111111111111111111111111111111111111", "job_id": 12},
				},
			},
		},
	}
}

func CommerceJobsFilters() types.CommerceJobsFilters {
	now := uint64(time.Now().Unix())
	ethName, ethLogo := mockChainMeta("1")
	baseName, baseLogo := mockChainMeta("8453")
	return types.CommerceJobsFilters{
		Chains: []types.CommerceJobsFilterChain{
			{ChainID: "8453", ChainName: baseName, ChainLogo: baseLogo},
			{ChainID: "1", ChainName: ethName, ChainLogo: ethLogo},
		},
		CommerceContracts: []string{
			"0x1111111111111111111111111111111111111111",
			"0x2222222222222222222222222222222222222222",
		},
		PaymentTokens: []string{
			"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", // USDC (example)
			"0x0000000000000000000000000000000000000000", // native placeholder
		},
		LastUpdated: now,
	}
}

type CommerceJobActionsQuery struct {
	ChainID          string
	CommerceContract string
	JobID            uint64
	ActionTypes      []string
	StartTime        *uint64
	EndTime          *uint64
	Page             int
	PageSize         int
}

func CommerceJobActions(q CommerceJobActionsQuery) ([]types.CommerceJobActionDTO, int64) {
	r := seededRand(q.ChainID, q.CommerceContract, fmt.Sprintf("%d", q.JobID), "job_actions")
	now := uint64(time.Now().Unix())
	allTypes := []string{"JobCreated", "JobFunded", "JobSubmitted", "JobCompleted", "JobRejected", "JobExpired", "ProviderSet", "BudgetSet", "PaymentReleased", "PlatformFeePaid", "EvaluatorFeePaid"}

	typeSet := make(map[string]struct{}, len(q.ActionTypes))
	for _, t := range q.ActionTypes {
		typeSet[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}

	all := make([]types.CommerceJobActionDTO, 0, 80)
	for i := 0; i < 80; i++ {
		at := allTypes[r.Intn(len(allTypes))]
		ts := now - uint64(r.Intn(3600*24*30))
		row := types.CommerceJobActionDTO{
			ChainID:          firstNonEmpty(q.ChainID, "1"),
			CommerceContract: firstNonEmpty(q.CommerceContract, "0x1111111111111111111111111111111111111111"),
			JobID:            q.JobID,
			ActionType:       at,
			BlockTimestamp:   ts,
			BlockNumber:      uint64(18000000 + r.Intn(100000)),
			TxHash:           fmt.Sprintf("0x%064x", r.Uint64()),
			LogIndex:         uint(r.Intn(20)),
			Actor:            fmt.Sprintf("0x%040x", r.Uint64()),
			Role:             []string{"client", "provider", "evaluator", "platform"}[r.Intn(4)],
			PaymentToken:     fmt.Sprintf("0x%040x", r.Uint64()),
			PaymentDecimals:  18,
			TokenSymbol:      []string{"USDC", "DAI", "ETH"}[r.Intn(3)],
		}
		if strings.Contains(strings.ToLower(at), "fee") || strings.EqualFold(at, "PaymentReleased") {
			v := round2(10 + r.Float64()*1000)
			row.Amount = &v
			v2 := round2(v * (0.8 + r.Float64()*1.2))
			row.AmountUSD = &v2
		}
		all = append(all, row)
	}

	// filter
	filtered := make([]types.CommerceJobActionDTO, 0, len(all))
	for _, it := range all {
		if len(typeSet) > 0 {
			if _, ok := typeSet[strings.ToLower(it.ActionType)]; !ok {
				continue
			}
		}
		if q.StartTime != nil && it.BlockTimestamp < *q.StartTime {
			continue
		}
		if q.EndTime != nil && it.BlockTimestamp > *q.EndTime {
			continue
		}
		filtered = append(filtered, it)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].BlockTimestamp > filtered[j].BlockTimestamp })
	pageItems, total := paginate(filtered, q.Page, q.PageSize)
	return pageItems, total
}

// -------------------- helpers --------------------

func firstNonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func hookAddrOrEmpty(r *rand.Rand, preferred string) string {
	if strings.TrimSpace(preferred) != "" {
		return preferred
	}
	if r.Intn(100) < 30 {
		return fmt.Sprintf("0x%040x", r.Uint64())
	}
	return ""
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}

func formatAmount8(v float64) string {
	if v == 0 || math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	s := fmt.Sprintf("%.8f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

func parseAmountOrZero(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconvParseFloat(s)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return v
}

// isolated helper to avoid pulling strconv in many places if we later swap parser
func strconvParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

