package api

import (
	"agent_identity/logger"
	"agent_identity/server/api/handle"
	apiUtils "agent_identity/server/api/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(nlogger *logger.Logger) {
	apiUtils.Init(nlogger)
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

	r.GET("agent/commerce/scores", handle.GetCommerceScoresHandler)
	r.GET("agent/commerce/actions", handle.GetCommerceActionsHandler)
	r.GET("agent/commerce/stats", handle.GetCommerceStatsHandler)

	r.GET("agent/identity/filter/info", handle.GetFilterInfoHandler)
	r.GET("agent/identity/filter/search/skill", handle.GetSkillsForFilterHandler)
	r.GET("agent/leaderboard", handle.GetLeaderboardInfoHandler)
	go apiUtils.UpdateGeneralInfo()
	r.Run(fmt.Sprintf(":%s", port))
}
