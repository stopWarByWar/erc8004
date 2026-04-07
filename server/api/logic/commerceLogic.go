package logic

import (
	"agent_identity/model"
	"fmt"
)

type CommerceScoreResp struct {
	Role                      string  `json:"role"`
	CompletedCount            int     `json:"completed_count"`
	RejectedCount             int     `json:"rejected_count"`
	ExpiredResponsibleCount   int     `json:"expired_responsible_count"`
	SuccessRate               float64 `json:"success_rate"`
	WeightedScore             float64 `json:"weighted_score"`
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
	return CommerceScoreResp{
		Role: s.Role, CompletedCount: s.CompletedCount, RejectedCount: s.RejectedCount,
		ExpiredResponsibleCount: s.ExpiredResponsibleCount, SuccessRate: s.SuccessRate,
		WeightedScore: s.WeightedScore, CreatedCount: s.CreatedCount, FundedCount: s.FundedCount,
		FundedRate: s.FundedRate, CompletionRate: s.CompletionRate,
		EvaluatedCount: s.EvaluatedCount, ExpiredFromSubmittedCount: s.ExpiredFromSubmittedCount,
		Responsiveness: s.Responsiveness, TotalJobs: s.TotalJobs, TotalVolume: s.TotalVolume,
		UniqueCounterparties: s.UniqueCounterparties, Confidence: s.Confidence,
	}
}

func globalToResp(g model.CommerceScoreGlobal) CommerceScoreResp {
	return CommerceScoreResp{
		Role: g.Role, CompletedCount: g.CompletedCount, RejectedCount: g.RejectedCount,
		ExpiredResponsibleCount: g.ExpiredResponsibleCount, SuccessRate: g.SuccessRate,
		WeightedScore: g.WeightedScore, CreatedCount: g.CreatedCount, FundedCount: g.FundedCount,
		FundedRate: g.FundedRate, CompletionRate: g.CompletionRate,
		EvaluatedCount: g.EvaluatedCount, ExpiredFromSubmittedCount: g.ExpiredFromSubmittedCount,
		Responsiveness: g.Responsiveness, TotalJobs: g.TotalJobs, TotalVolume: g.TotalVolume,
		UniqueCounterparties: g.UniqueCounterparties, Confidence: g.Confidence,
	}
}
