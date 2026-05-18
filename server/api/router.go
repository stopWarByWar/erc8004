package api

import (
	"agent_identity/logger"
	"agent_identity/server/api/handle"
	serverLogic "agent_identity/server/api/logic"
	apiUtils "agent_identity/server/api/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(nlogger *logger.Logger, rc apiUtils.RuntimeConfig) {
	apiUtils.Init(nlogger, rc)
}

func Run(_cors []string, port string) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     _cors,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Content-Disposition", "multipart/form-data"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("agent/identity/list", handle.GetAgentCardListHandler)
	r.GET("agent/identity/detail", handle.GetAgentCardDetailHandler)
	// r.GET("agent/identity/trustModel", handle.GetTrustModelListHandler)
	// r.GET("agent/identity/filter/info", handle.GetFilterInfoHandler)
	// r.GET("agent/identity/search/skill", handle.GetAgentCardsSearchBySkillHandler)
	// r.GET("agent/identity/search/name", handle.GetAgentCardsSearchByNameHandler)
	r.POST("agent/identity/search/semantic", handle.GetAgentCardsSearchBySemanticHandler)

	r.GET("agent/identity/general/info", handle.GetAgentAmountForEachChainHandler)

	r.POST("agent/identity/set/profile", handle.UploadAgentProfileHandler)

	r.GET("agent/identity/detail/validation/list", handle.GetAgentValidationListHandler)
	r.GET("agent/identity/detail/validation/eval/latest", handle.GetAgentValidationEvalLatestHandler)
	r.GET("agent/validator/list", handle.GetValidatorListHandler)
	r.GET("agent/validator/detail/validation/list", handle.GetValidatorValidationListHandler)
	r.GET("agent/validator/detail/address", handle.GetValidatorByAddressHandler)

	r.GET("agent/network/list", handle.GetNetworkListHandler)

	r.GET("agent/identity/detail/feedback/scores", handle.GetFeedbackScoresHandler)
	r.GET("agent/identity/detail/feedbacks", handle.GetAgentFeedbacksHandler)
	r.POST("agent/identity/set/feedback", handle.UploadFeedbackHandler)

	// ERC-8004 Feedback Credit Score (design: docs/designs/202605120000_feedback-credit-score.html)
	r.GET("agent/feedback/credit", handle.GetFeedbackCreditHandler)

	r.GET("agent/identity/commerce/scores", handle.GetCommerceScoresHandler)
	r.GET("agent/identity/commerce/actions", handle.GetCommerceActionsHandler)
	r.GET("agent/identity/commerce/stats", handle.GetCommerceStatsHandler)

	// Job Browser / Job Detail (ERC-8183)
	r.GET("agent/commerce/jobs", handle.GetCommerceJobsHandler)
	r.GET("agent/commerce/jobs/detail", handle.GetCommerceJobDetailHandler)
	r.GET("agent/commerce/jobs/general", handle.GetCommerceJobsGeneralHandler)
	r.GET("agent/commerce/jobs/charts", handle.GetCommerceJobsChartsHandler)
	r.GET("agent/commerce/jobs/actions", handle.GetCommerceJobActionsHandler)
	r.GET("agent/commerce/jobs/filters", handle.GetCommerceJobsFiltersHandler)

	r.GET("agent/identity/filter/info", handle.GetFilterInfoHandler)
	r.GET("agent/identity/filter/search/skill", handle.GetSkillsForFilterHandler)
	r.GET("agent/leaderboard", handle.GetLeaderboardInfoHandler)

	// Passport / Agent Summary (ERC-8183)
	r.GET("agent/identity/commerce/passport", handle.GetPassportHandler)

	// Unified agent credits view: identity + Commerce (8183) + Feedback (8004).
	// Powers the redesigned agent detail page. See dev-doc/agent-credits.html.
	r.GET("agent/identity/credits", handle.GetAgentCreditsHandler)

	go apiUtils.UpdateGeneralInfo()
	apiUtils.StartJobsFiltersCache(context.Background(), 5*time.Minute)
	// ERC-8004 Feedback Credit Score cache preheater. See
	// docs/designs/202605120000_feedback-credit-score.html
	serverLogic.StartFeedbackCreditCron(context.Background(), 5*time.Minute)
	r.Run(fmt.Sprintf(":%s", port))
}
