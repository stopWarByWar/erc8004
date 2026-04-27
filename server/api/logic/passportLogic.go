package logic

import (
	"agent_identity/config"
	"agent_identity/model"
	"agent_identity/server/api/types"
	"fmt"
	"strings"
)

// GetPassportData returns a comprehensive passport/agent summary for the given uid.
func GetPassportData(uid uint64) (*types.PassportResponse, error) {
	return getRealPassport(uid)
}

// getRealPassport queries real DB data to build a passport response.
func getRealPassport(uid uint64) (*types.PassportResponse, error) {
	// 1. Get agent basic info
	agent, err := model.GetAgentByUID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	if agent == nil || agent.AgentID == "" {
		return nil, fmt.Errorf("agent not found for uid %d", uid)
	}

	uidStr := fmt.Sprintf("0x%x", uid)

	// 2. Resolve chain name and logo
	chainName := agent.ChainID
	chainLogo := ""
	if chainInfo, ok := config.GetChainInfo(agent.ChainID); ok {
		chainName = chainInfo.ChainName
		chainLogo = chainInfo.ChainLogo
	}

	// 3. Get skills from OASFSkill table
	skills, err := model.GetOASFSkillsByAgentUID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	var skillNames []string
	for _, s := range skills {
		if s.SkillName != "" {
			skillNames = append(skillNames, s.SkillName)
		}
	}

	// 4. Get basic stats from commerce_jobs table
	basicStats, err := getPassportBasicStats(uid, agent.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	// 5. Get commerce scores from commerce_scores_global table
	commerceScores, err := getPassportCommerceScores(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get commerce scores: %w", err)
	}

	// 6. Determine verification level
	verification := determineVerificationLevel(uid, agent)

	return &types.PassportResponse{
		UID:           uidStr,
		Name:          agent.Name,
		Avatar:        agent.Image,
		ChainID:       agent.ChainID,
		ChainName:     chainName,
		ChainLogo:     chainLogo,
		CreatedAt:     int64(agent.Timestamps),
		BasicStats:    *basicStats,
		CommerceScore: commerceScores,
		Skills:        skillNames,
		Verification:  verification,
		ShareURL:      fmt.Sprintf("/passport/%s/shared", uidStr),
	}, nil
}

// getPassportBasicStats computes basic job statistics for an agent.
func getPassportBasicStats(uid uint64, chainID string) (*types.PassportBasicStats, error) {
	providerAddress := fmt.Sprintf("0x%x", uid)

	jobs, err := model.GetCommerceJobsByProviderAddress(providerAddress, chainID)
	if err != nil {
		return nil, err
	}

	totalJobs := len(jobs)
	completedJobs := 0
	activeJobs := 0
	for _, job := range jobs {
		switch job.Status {
		case "Completed":
			completedJobs++
		case "Active", "Funded", "Submitted", "Paid", "InProgress":
			activeJobs++
		}
	}

	reputationScore, feedbackCount, err := model.GetAgentReputationStats(uid)
	if err != nil {
		reputationScore = 0
	}

	return &types.PassportBasicStats{
		TotalJobs:       totalJobs,
		CompletedJobs:   completedJobs,
		ActiveJobs:      activeJobs,
		ReputationScore: reputationScore,
		FeedbackCount:   feedbackCount,
	}, nil
}

// getPassportCommerceScores fetches global commerce scores for all roles.
func getPassportCommerceScores(uid uint64) (map[string]types.CommerceScore, error) {
	scores, err := model.GetCommerceScoreGlobal(uid)
	if err != nil {
		return nil, err
	}

	result := make(map[string]types.CommerceScore)
	for _, s := range scores {
		result[strings.ToLower(s.Role)] = types.CommerceScore{
			Role:                       strings.ToLower(s.Role),
			CompletedCount:             s.CompletedCount,
			RejectedCount:              s.RejectedCount,
			ExpiredResponsibleCount:     s.ExpiredResponsibleCount,
			SuccessRate:                s.SuccessRate,
			WeightedScore:              s.WeightedScore,
			TotalVolumeUSD:             s.TotalVolumeUSD,
			WeightedScoreUSD:           s.WeightedVolumeUSD,
			TotalJobs:                  s.TotalJobs,
			TotalVolume:                s.TotalVolume,
			UniqueCounterparties:       s.UniqueCounterparties,
			Confidence:                 s.Confidence,
			CreatedCount:               s.CreatedCount,
			FundedCount:               s.FundedCount,
			FundedRate:                s.FundedRate,
			CompletionRate:             s.CompletionRate,
			EvaluatedCount:             s.EvaluatedCount,
			ExpiredFromSubmittedCount:  s.ExpiredFromSubmittedCount,
			Responsiveness:             s.Responsiveness,
		}
	}
	return result, nil
}

// determineVerificationLevel checks if an agent has valid attestations.
// "advanced" if attestations exist, "basic" otherwise.
func determineVerificationLevel(uid uint64, _ *model.Agent) string {
	// Check via attestation count using a model query
	count, err := model.GetAttestationCountByRecipient(uid)
	if err == nil && count > 0 {
		return "advanced"
	}
	return "basic"
}
