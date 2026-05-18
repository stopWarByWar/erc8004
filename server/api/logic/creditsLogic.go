// Package logic — unified agent credits endpoint.
//
// GetAgentCredits assembles a single response that powers the redesigned
// agent detail page: identity + Commerce (ERC-8183) topline + Feedback Credit
// (ERC-8004) topline + a merged tag pool. It reuses existing query layers
// (model.GetAgentByUID, GetCommerceScoreGlobal, GetFeedbackCredit) so this is
// a pure aggregation — no new SQL.
//
// Design: see dev-doc/agent-credits.html for the visual mapping and
// docs/designs/202605120000_feedback-credit-score.html for the score math.
package logic

import (
	"fmt"

	"agent_identity/config"
	"agent_identity/model"
)

// ─── Response types ────────────────────────────────────────────────────────

// AgentCreditsResp is the single response that backs the redesigned agent
// detail page. Each top-level field maps to one card on the page.
type AgentCreditsResp struct {
	Identity *AgentCreditsIdentity   `json:"identity"`
	Commerce *CommerceCreditsTopline `json:"commerce"`
	Feedback *FeedbackCreditResp     `json:"feedback"`
}

// AgentCreditsIdentity backs the header card (avatar / name / chain / wallet
// / verification). Shape stays close to PassportResponse for FE reuse.
type AgentCreditsIdentity struct {
	UID              string   `json:"uid"`             // "0x123" — same hex form as Passport
	AgentID          string   `json:"agent_id"`        // on-chain agent_id (decimal string)
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Avatar           string   `json:"avatar"`          // agents.image
	AgentWallet      string   `json:"agent_wallet"`
	Owner            string   `json:"owner"`
	ChainID          string   `json:"chain_id"`
	ChainName        string   `json:"chain_name"`
	ChainLogo        string   `json:"chain_logo"`
	IdentityRegistry string   `json:"identity_registry"`
	RegisteredAt     int64    `json:"registered_at"`   // agents.timestamps (unix sec)
	Verification     string   `json:"verification_level"` // "basic" | "advanced"
	Skills           []string `json:"skills"`
}

// CommerceCreditsTopline backs the left "Commerce Score" card.
// Score / Confidence are headline numbers; Roles carries the per-role detail
// so the card can show a 2-3 line breakdown beneath the headline.
type CommerceCreditsTopline struct {
	// Score is the headline Commerce score in [0, 1]. nil = no commerce data
	// (cold-start). FE should render as `score × 100`.
	Score *float64 `json:"score"`
	// Confidence in [0, 1] — peak confidence across roles.
	Confidence float64 `json:"confidence"`
	// TotalJobs is the sum of TotalJobs across all roles.
	TotalJobs int `json:"total_jobs"`
	// TotalVolumeUSD is the sum of TotalVolumeUSD across all roles.
	TotalVolumeUSD float64             `json:"total_volume_usd"`
	Roles          []CommerceScoreResp `json:"roles"` // reuse existing per-role struct
}

// ─── Main entry point ─────────────────────────────────────────────────────

// GetAgentCredits returns the merged identity + commerce + feedback view for
// the given agent UID. Cold-start cases are represented by nil scores or
// empty slices, never errors.
func GetAgentCredits(uid uint64) (*AgentCreditsResp, error) {
	identity, err := buildAgentCreditsIdentity(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to build identity: %w", err)
	}

	commerce, err := buildCommerceCreditsTopline(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to build commerce topline: %w", err)
	}

	feedback, err := GetFeedbackCredit(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to build feedback credit: %w", err)
	}

	return &AgentCreditsResp{
		Identity: identity,
		Commerce: commerce,
		Feedback: feedback,
	}, nil
}

// ─── Internal builders ────────────────────────────────────────────────────

// buildAgentCreditsIdentity reads the agents row + skills + chain config +
// verification level and assembles the header card payload.
func buildAgentCreditsIdentity(uid uint64) (*AgentCreditsIdentity, error) {
	agent, err := model.GetAgentByUID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	if agent == nil || agent.AgentID == "" {
		return nil, fmt.Errorf("agent not found for uid %d", uid)
	}

	chainName := agent.ChainID
	chainLogo := ""
	if info, ok := config.GetChainInfo(agent.ChainID); ok {
		chainName = info.ChainName
		chainLogo = info.ChainLogo
	}

	skills, err := model.GetOASFSkillsByAgentUID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	skillNames := make([]string, 0, len(skills))
	for _, s := range skills {
		if s.SkillName != "" {
			skillNames = append(skillNames, s.SkillName)
		}
	}

	return &AgentCreditsIdentity{
		UID:              fmt.Sprintf("0x%x", uid),
		AgentID:          agent.AgentID,
		Name:             agent.Name,
		Description:      agent.Description,
		Avatar:           agent.Image,
		AgentWallet:      agent.AgentWallet,
		Owner:            agent.Owner,
		ChainID:          agent.ChainID,
		ChainName:        chainName,
		ChainLogo:        chainLogo,
		IdentityRegistry: agent.IdentityRegistry,
		RegisteredAt:     int64(agent.Timestamps),
		Verification:     determineVerificationLevel(uid, agent),
		Skills:           skillNames,
	}, nil
}

// buildCommerceCreditsTopline aggregates per-role global commerce scores into
// the headline card payload. Score is the peak WeightedScore across roles;
// totals are summed. Returns a non-nil topline even with zero roles so the FE
// can render a clean cold-start ("—") state.
func buildCommerceCreditsTopline(uid uint64) (*CommerceCreditsTopline, error) {
	globals, err := model.GetCommerceScoreGlobal(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get commerce score global: %w", err)
	}

	topline := &CommerceCreditsTopline{
		Roles: make([]CommerceScoreResp, 0, len(globals)),
	}
	var peakScore float64
	var peakConf float64
	hasScore := false
	for _, g := range globals {
		topline.Roles = append(topline.Roles, globalToResp(g))
		topline.TotalJobs += g.TotalJobs
		topline.TotalVolumeUSD += g.TotalVolumeUSD
		if g.WeightedScore > peakScore {
			peakScore = g.WeightedScore
			hasScore = true
		}
		if g.Confidence > peakConf {
			peakConf = g.Confidence
		}
	}
	if hasScore {
		topline.Score = &peakScore
	}
	topline.Confidence = peakConf
	return topline, nil
}
