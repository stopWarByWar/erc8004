package handle

import (
	agentcard "agent_identity/agentCard"
	"agent_identity/config"
	apiUtils "agent_identity/server/api/utils"

	"github.com/gin-gonic/gin"
)

func GetChainListHandler(c *gin.Context) {
	apiUtils.SuccessResp(gin.H{
		"chain_list": config.ChainList,
	}, c)
}

func GetTrustModelListHandler(c *gin.Context) {
	apiUtils.SuccessResp(gin.H{
		"trust_model_list": []string{agentcard.TrustModelFeedback, agentcard.TrustModelInferenceValidation, agentcard.TrustModelTeeAttestation},
	}, c)
}

func GetGeneralInfoHandler(c *gin.Context) {
	apiUtils.SuccessResp(gin.H(apiUtils.GetGeneralInfo()), c)
}
