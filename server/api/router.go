package api

import (
	"agent_identity/config"
	"agent_identity/logger"
	"agent_identity/model"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var generalInfo = make(map[string]any)

func InitRouter(nlogger *logger.Logger, _mock bool, _feedbackMock bool) {
	_logger = nlogger
	mock = _mock
	feedbackMock = _feedbackMock
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

	r.GET("agent/identity/list", GetAgentCardListHandler)
	r.GET("agent/identity/detail", GetAgentCardDetailHandler)
	r.GET("agent/identity/trustModel", GetTrustModelListHandler)
	r.GET("agent/identity/chains", GetChainListHandler)
	r.GET("agent/identity/search/skill", GetAgentCardsSearchBySkillHandler)
	r.GET("agent/identity/search/name", GetAgentCardsSearchByNameHandler)
	r.POST("agent/identity/search/semantic", GetAgentCardsSearchBySemanticHandler)
	r.GET("agent/identity/detail/comments", GetAgentCommentsHandler)
	r.GET("agent/identity/detail/feedbacks", GetAgentFeedbacksHandler)
	r.GET("agent/identity/general/info", GetGeneralInfoHandler)

	r.POST("agent/identity/set/feedback", UploadFeedbackHandler)
	r.POST("agent/identity/set/profile", UploadAgentProfileHandler)

	r.GET("agent/identity/detail/validation/list", GetAgentValidationListHandler)
	r.GET("agent/validator/list", GetValidatorListHandler)
	r.GET("agent/validator/detail/validation/list", GetValidatorValidationListHandler)

	r.GET("agent/network/list", GetNetworkListHandler)

	go updateGeneralInfo()
	r.Run(fmt.Sprintf(":%s", port))
}

func updateGeneralInfo() {
	for {
		agentsAmount, err := model.GetAgentAmountForEachChain()
		if err != nil {
			fmt.Println("failed to get agent amount for each chain", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		total := int64(0)
		for chainId, amount := range agentsAmount {
			chainInfo := config.GetChainInfo(chainId)
			generalInfo[strings.Replace(chainInfo.ChainName, " ", "_", -1)] = amount
			total += amount
		}
		generalInfo["total"] = total
		time.Sleep(5 * time.Minute)
	}
}
