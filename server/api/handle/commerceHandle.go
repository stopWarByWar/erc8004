package handle

import (
	serverLogic "agent_identity/server/api/logic"
	"agent_identity/model"
	apiMock "agent_identity/server/api/mock"
	"agent_identity/server/api/types"
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
	allowedJobStatus = map[string]struct{}{
		"": {}, "open": {}, "funded": {}, "submitted": {}, "completed": {}, "rejected": {}, "expired": {},
	}
	allowedJobSortBy = map[string]struct{}{
		"": {}, "updated_at": {}, "budget": {}, "budget_usd": {}, "paid_amount_usd": {},
	}
	allowedJobActionTypes = map[string]struct{}{
		// Job lifecycle & config events (from commerce_actions.action)
		"job_created": {}, "job_funded": {}, "job_submitted": {},
		"job_completed": {}, "job_rejected": {}, "job_expired": {},
		"provider_set": {}, "budget_set": {},
		// Settlement breakdown events
		"payment_released": {}, "platform_fee_paid": {}, "evaluator_fee_paid": {},
	}
	allowedSortOrder = map[string]struct{}{
		"": {}, "asc": {}, "desc": {},
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

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		serverUtils.SuccessResp(gin.H{"scores": apiMock.CommerceScores(uid, chainID, contract)}, c)
		return
	}

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

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		dto, total := apiMock.CommerceActions(apiMock.CommerceActionsQuery{
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
		serverUtils.SuccessResp(gin.H{"actions": dto, "total": total}, c)
		return
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
	dto := make([]types.CommerceActionDTO, 0, len(actionsRes))
	for _, a := range actionsRes {
		dto = append(dto, serverLogic.ActionToDTO(a))
	}
	serverUtils.SuccessResp(gin.H{"actions": dto, "total": total}, c)
}

func GetCommerceStatsHandler(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get uid", "Invalid Request", c)
		return
	}

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		now := uint64(time.Now().Unix())
		serverUtils.SuccessResp(gin.H{"stats": apiMock.CommerceStats(uid, now)}, c)
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

func GetCommerceJobsHandler(c *gin.Context) {
	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")
	status := c.Query("status")
	role := c.Query("role")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	agentAddress := c.Query("agent_address")
	counterparty := c.Query("counterparty")
	paymentToken := c.Query("payment_token")
	tokenSymbol := c.Query("token_symbol")

	if _, ok := allowedJobStatus[status]; !ok {
		serverUtils.ErrResp(nil, "invalid status", "Invalid Request", c)
		return
	}
	if role != "" {
		if _, ok := allowedRoles[role]; !ok {
			serverUtils.ErrResp(nil, "invalid role", "Invalid Request", c)
			return
		}
	}
	if _, ok := allowedJobSortBy[sortBy]; !ok {
		serverUtils.ErrResp(nil, "invalid sort_by", "Invalid Request", c)
		return
	}
	if _, ok := allowedSortOrder[sortOrder]; !ok {
		serverUtils.ErrResp(nil, "invalid sort_order", "Invalid Request", c)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

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

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		jobs, total := apiMock.CommerceJobs(apiMock.CommerceJobsQuery{
			ChainID:          chainID,
			CommerceContract: contract,
			Status:           status,
			Role:             role,
			AgentAddress:     agentAddress,
			Counterparty:     counterparty,
			PaymentToken:     paymentToken,
			TokenSymbol:      tokenSymbol,
			MinBudget:        minBudget,
			MaxBudget:        maxBudget,
			MinBudgetUSD:     minBudgetUSD,
			MaxBudgetUSD:     maxBudgetUSD,
			StartTime:        startTime,
			EndTime:          endTime,
			SortBy:           sortBy,
			SortOrder:        sortOrder,
			Page:             page,
			PageSize:         pageSize,
		})
		serverUtils.SuccessResp(gin.H{"jobs": jobs, "total": total}, c)
		return
	}

	jobs, total, err := serverLogic.GetCommerceJobs(serverLogic.CommerceJobsParams{
		ChainID:          chainID,
		CommerceContract: contract,
		Status:           status,
		Role:             role,
		AgentAddress:     agentAddress,
		Counterparty:     counterparty,
		PaymentToken:     paymentToken,
		TokenSymbol:      tokenSymbol,
		MinBudget:        minBudget,
		MaxBudget:        maxBudget,
		MinBudgetUSD:     minBudgetUSD,
		MaxBudgetUSD:     maxBudgetUSD,
		StartTime:        startTime,
		EndTime:          endTime,
		SortBy:           sortBy,
		SortOrder:        sortOrder,
		Page:             page,
		PageSize:         pageSize,
	})
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce jobs", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{"jobs": jobs, "total": total}, c)
}

func GetCommerceJobDetailHandler(c *gin.Context) {
	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")
	jobID, err := strconv.ParseUint(c.Query("job_id"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "invalid job_id", "Invalid Request", c)
		return
	}
	if chainID == "" || contract == "" {
		serverUtils.ErrResp(nil, "missing chain_id or commerce_contract", "Invalid Request", c)
		return
	}

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		job, evidence := apiMock.CommerceJobDetail(chainID, contract, jobID)
		serverUtils.SuccessResp(gin.H{"job": job, "evidence": evidence}, c)
		return
	}

	job, err := serverLogic.GetCommerceJobDetail(chainID, contract, jobID)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce job", "Internal Error", c)
		return
	}
	if job == nil {
		// Keep project convention: still 200 with failed message.
		serverUtils.ErrResp(nil, "job not found", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"job": job,
		"evidence": gin.H{
			"timeline_source":      "commerce_actions",
			"events_table_source":  "commerce_actions",
			"settlement_events":    []string{"PaymentReleased", "PlatformFeePaid", "EvaluatorFeePaid"},
		},
	}, c)
}

func GetCommerceJobsGeneralHandler(c *gin.Context) {
	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")
	status := c.Query("status")
	role := c.Query("role")
	agentAddress := c.Query("agent_address")
	counterparty := c.Query("counterparty")
	paymentToken := c.Query("payment_token")
	tokenSymbol := c.Query("token_symbol")

	if _, ok := allowedJobStatus[status]; !ok {
		serverUtils.ErrResp(nil, "invalid status", "Invalid Request", c)
		return
	}
	if role != "" {
		if _, ok := allowedRoles[role]; !ok {
			serverUtils.ErrResp(nil, "invalid role", "Invalid Request", c)
			return
		}
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

	includeDistributions := true
	if raw := c.Query("include_distributions"); raw != "" {
		v, perr := strconv.ParseBool(raw)
		if perr != nil {
			serverUtils.ErrResp(nil, "invalid include_distributions", "Invalid Request", c)
			return
		}
		includeDistributions = v
	}

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		summary, distributions := apiMock.CommerceJobsGeneral()
		serverUtils.SuccessResp(gin.H{"summary": summary, "distributions": distributions}, c)
		return
	}

	resp, err := serverLogic.GetCommerceJobsGeneral(serverLogic.CommerceJobsParams{
		ChainID:          chainID,
		CommerceContract: contract,
		Status:           status,
		Role:             role,
		AgentAddress:     agentAddress,
		Counterparty:     counterparty,
		PaymentToken:     paymentToken,
		TokenSymbol:      tokenSymbol,
		MinBudget:        minBudget,
		MaxBudget:        maxBudget,
		MinBudgetUSD:     minBudgetUSD,
		MaxBudgetUSD:     maxBudgetUSD,
		StartTime:        startTime,
		EndTime:          endTime,
		Page:             1,
		PageSize:         1,
	}, includeDistributions)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce jobs general", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{"summary": resp.Summary, "distributions": resp.Distributions}, c)
}

func GetCommerceJobsChartsHandler(c *gin.Context) {
	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")
	status := c.Query("status")
	paymentToken := c.Query("payment_token")
	tokenSymbol := c.Query("token_symbol")
	role := c.Query("role")
	agentAddress := c.Query("agent_address")
	counterparty := c.Query("counterparty")

	if _, ok := allowedJobStatus[status]; !ok {
		serverUtils.ErrResp(nil, "invalid status", "Invalid Request", c)
		return
	}
	if role != "" {
		if _, ok := allowedRoles[role]; !ok {
			serverUtils.ErrResp(nil, "invalid role", "Invalid Request", c)
			return
		}
	}

	window := c.DefaultQuery("window", "7d")
	switch window {
	case "24h", "3d", "7d", "1month", "all", "30d", "custom":
		// ok
	default:
		serverUtils.ErrResp(nil, "invalid window", "Invalid Request", c)
		return
	}
	var bucketSeconds uint64
	if raw := c.Query("bucket_seconds"); raw != "" {
		v, perr := strconv.ParseUint(raw, 10, 64)
		if perr != nil || v == 0 {
			serverUtils.ErrResp(nil, "invalid bucket_seconds", "Invalid Request", c)
			return
		}
		bucketSeconds = v
	}
	// Default bucket seconds by window
	if bucketSeconds == 0 {
		switch window {
		case "24h":
			bucketSeconds = 3600
		case "3d":
			bucketSeconds = 3 * 3600
		case "7d":
			bucketSeconds = 6 * 3600
		case "1month":
			bucketSeconds = 24 * 3600
		case "all":
			bucketSeconds = 24 * 3600
		case "30d":
			// Backward compatible alias; prefer 1month in new clients.
			bucketSeconds = 24 * 3600
		case "custom":
			serverUtils.ErrResp(nil, "missing bucket_seconds for custom window", "Invalid Request", c)
			return
		default:
			serverUtils.ErrResp(nil, "invalid window", "Invalid Request", c)
			return
		}
	}

	// For now, we treat window as optional and rely on start/end if provided.
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

	ws, we := uint64(0), uint64(0)
	if startTime != nil {
		ws = *startTime
	}
	if endTime != nil {
		we = *endTime
	}

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		serverUtils.SuccessResp(gin.H{"charts": apiMock.CommerceJobsCharts()}, c)
		return
	}

	charts, err := serverLogic.GetCommerceJobsCharts(serverLogic.CommerceJobsChartsParams{
		CommerceJobsParams: serverLogic.CommerceJobsParams{
			ChainID:          chainID,
			CommerceContract: contract,
			Status:           status,
			Role:             role,
			AgentAddress:     agentAddress,
			Counterparty:     counterparty,
			PaymentToken:     paymentToken,
			TokenSymbol:      tokenSymbol,
			MinBudget:        minBudget,
			MaxBudget:        maxBudget,
			MinBudgetUSD:     minBudgetUSD,
			MaxBudgetUSD:     maxBudgetUSD,
			StartTime:        func() *uint64 { if ws == 0 { return nil }; return &ws }(),
			EndTime:          func() *uint64 { if we == 0 { return nil }; return &we }(),
			Page:             1,
			PageSize:         1,
		},
		Window:        window,
		BucketSeconds: bucketSeconds,
	})
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce jobs charts", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{"charts": charts}, c)
}

func GetCommerceJobActionsHandler(c *gin.Context) {
	chainID := c.Query("chain_id")
	contract := c.Query("commerce_contract")
	jobID, err := strconv.ParseUint(c.Query("job_id"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "invalid job_id", "Invalid Request", c)
		return
	}
	if chainID == "" || contract == "" {
		serverUtils.ErrResp(nil, "missing chain_id or commerce_contract", "Invalid Request", c)
		return
	}

	actionTypes := splitCSV(c.Query("action_type"))
	for _, a := range actionTypes {
		if _, ok := allowedJobActionTypes[a]; !ok {
			serverUtils.ErrResp(nil, "invalid action_type", "Invalid Request", c)
			return
		}
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

	if serverUtils.IsMockEnabled() {
		serverUtils.SetMockHeader(c)
		actions, total := apiMock.CommerceJobActions(apiMock.CommerceJobActionsQuery{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            jobID,
			ActionTypes:      actionTypes,
			StartTime:        startTime,
			EndTime:          endTime,
			Page:             page,
			PageSize:         pageSize,
		})
		serverUtils.SuccessResp(gin.H{"actions": actions, "total": total}, c)
		return
	}

	actions, total, err := serverLogic.GetCommerceJobActions(serverLogic.CommerceJobActionsParams{
		ChainID:          chainID,
		CommerceContract: contract,
		JobID:            jobID,
		ActionTypes:      actionTypes,
		StartTime:        startTime,
		EndTime:          endTime,
		Page:             page,
		PageSize:         pageSize,
	})
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get commerce job actions", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{"actions": actions, "total": total}, c)
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
