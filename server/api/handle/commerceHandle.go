package handle

import (
	serverLogic "agent_identity/server/api/logic"
	"agent_identity/model"
	serverUtils "agent_identity/server/api/utils"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	allowedRoles      = map[string]struct{}{"provider": {}, "client": {}, "evaluator": {}}
	allowedCertainty  = map[string]struct{}{"definitive": {}, "indicative": {}}
	allowedPolarities = map[string]struct{}{"positive": {}, "negative": {}, "neutral": {}}
	allowedSortBy     = map[string]struct{}{"": {}, "timestamp": {}, "budget": {}, "signal_weight": {}}
	allowedActions    = map[string]struct{}{
		"job_created": {}, "job_funded": {}, "job_submitted": {},
		"job_completed": {}, "job_rejected": {}, "job_expired": {},
		"provider_set": {}, "budget_set": {},
	}
)

func GetCommerceScoresHandler(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get uid", "Invalid Request", c)
		return
	}

	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")

	scores, err := serverLogic.GetCommerceScores(uid, chainID, contract)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce scores", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{"scores": scores}, c)
}

func GetCommerceActionsHandler(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get uid", "Invalid Request", c)
		return
	}

	role := c.Query("role")
	certainty := c.Query("certainty")
	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")
	counterparty := c.Query("counterparty")
	hookAddress := c.Query("hook_address")
	sortBy := c.Query("sort_by")

	actions := splitCSV(c.Query("action"))
	polarities := splitCSV(c.Query("polarity"))

	// Enum validations (return 400 instead of 500 for bad inputs)
	if role != "" {
		if _, ok := allowedRoles[role]; !ok {
			serverUtils.ErrResp(nil, "invalid role", "Invalid Request", c)
			return
		}
	}
	if certainty != "" {
		if _, ok := allowedCertainty[certainty]; !ok {
			serverUtils.ErrResp(nil, "invalid certainty", "Invalid Request", c)
			return
		}
	}
	if _, ok := allowedSortBy[sortBy]; !ok {
		serverUtils.ErrResp(nil, "invalid sort_by", "Invalid Request", c)
		return
	}
	for _, a := range actions {
		if _, ok := allowedActions[a]; !ok {
			serverUtils.ErrResp(nil, "invalid action", "Invalid Request", c)
			return
		}
	}
	for _, p := range polarities {
		if _, ok := allowedPolarities[p]; !ok {
			serverUtils.ErrResp(nil, "invalid polarity", "Invalid Request", c)
			return
		}
	}

	var hasHook *bool
	if raw := c.Query("has_hook"); raw != "" {
		v, perr := strconv.ParseBool(raw)
		if perr != nil {
			serverUtils.ErrResp(nil, "invalid has_hook", "Invalid Request", c)
			return
		}
		hasHook = &v
	}

	var minBudget *float64
	if raw := c.Query("min_budget"); raw != "" {
		v, perr := strconv.ParseFloat(raw, 64)
		if perr != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			serverUtils.ErrResp(nil, "invalid min_budget", "Invalid Request", c)
			return
		}
		minBudget = &v
	}
	var maxBudget *float64
	if raw := c.Query("max_budget"); raw != "" {
		v, perr := strconv.ParseFloat(raw, 64)
		if perr != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			serverUtils.ErrResp(nil, "invalid max_budget", "Invalid Request", c)
			return
		}
		maxBudget = &v
	}

	paymentToken := c.Query("payment_token")
	tokenSymbol := c.Query("token_symbol")

	var minBudgetUSD *float64
	if raw := c.Query("min_budget_usd"); raw != "" {
		v, perr := strconv.ParseFloat(raw, 64)
		if perr != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			serverUtils.ErrResp(nil, "invalid min_budget_usd", "Invalid Request", c)
			return
		}
		minBudgetUSD = &v
	}
	var maxBudgetUSD *float64
	if raw := c.Query("max_budget_usd"); raw != "" {
		v, perr := strconv.ParseFloat(raw, 64)
		if perr != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			serverUtils.ErrResp(nil, "invalid max_budget_usd", "Invalid Request", c)
			return
		}
		maxBudgetUSD = &v
	}

	var startTime *uint64
	if raw := c.Query("start_time"); raw != "" {
		v, perr := strconv.ParseUint(raw, 10, 64)
		if perr != nil {
			serverUtils.ErrResp(nil, "invalid start_time", "Invalid Request", c)
			return
		}
		startTime = &v
	}
	var endTime *uint64
	if raw := c.Query("end_time"); raw != "" {
		v, perr := strconv.ParseUint(raw, 10, 64)
		if perr != nil {
			serverUtils.ErrResp(nil, "invalid end_time", "Invalid Request", c)
			return
		}
		endTime = &v
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	actionsRes, total, err := serverLogic.GetCommerceActionsV2(serverLogic.CommerceActionsParams{
		UID:          uid,
		Role:         role,
		Actions:      actions,
		Certainty:    certainty,
		Polarities:   polarities,
		ChainID:      chainID,
		Contract:     contract,
		Counterparty: counterparty,
		HookAddress:  hookAddress,
		HasHook:      hasHook,
		MinBudget:    minBudget,
		MaxBudget:    maxBudget,
		PaymentToken: paymentToken,
		TokenSymbol:  tokenSymbol,
		MinBudgetUSD: minBudgetUSD,
		MaxBudgetUSD: maxBudgetUSD,
		StartTime:    startTime,
		EndTime:      endTime,
		SortBy:       sortBy,
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce actions", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{"actions": actionsRes, "total": total}, c)
}

func GetCommerceStatsHandler(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get uid", "Invalid Request", c)
		return
	}

	// Get latest block_timestamp for this agent to serve as "now"
	latestAction, err := model.GetLatestCommerceActionTimestampByUID(uid)
	now := uint64(0)
	if err == nil && latestAction > 0 {
		now = latestAction
	} else {
		// Fallback to system time
		now = uint64(time.Now().Unix())
	}

	stats, err := serverLogic.GetCommerceStats(serverLogic.CommerceStatsParams{
		UID: uid,
		Now: now,
	})
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce stats", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{"stats": stats}, c)
}

func splitCSV(v string) []string {
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
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
