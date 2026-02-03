package handle

import (
	"agent_identity/model"
	"agent_identity/server/api/logic"
	apiUtils "agent_identity/server/api/utils"

	"agent_identity/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// func GetChainListHandler(c *gin.Context) {
// 	apiUtils.SuccessResp(gin.H{
// 		"chain_list": config.ChainList,
// 	}, c)
// }

// func GetTrustModelListHandler(c *gin.Context) {
// 	apiUtils.SuccessResp(gin.H{
// 		"trust_model_list": []string{agentcard.TrustModelFeedback, agentcard.TrustModelInferenceValidation, agentcard.TrustModelTeeAttestation},
// 	}, c)
// }

// func GetGeneralInfoHandler(c *gin.Context) {
// 	apiUtils.SuccessResp(gin.H(apiUtils.GetGeneralInfo()), c)
// }

func GetAgentAmountForEachChainHandler(c *gin.Context) {
	chainInfos, total := logic.GetAgentAmountForEachChain()
	apiUtils.SuccessResp(gin.H{
		"chains": chainInfos,
		"total":  total,
	}, c)
}

func GetFilterInfoHandler(c *gin.Context) {
	apiUtils.SuccessResp(gin.H{
		"filter_info": config.GetGeneralInfo(),
	}, c)
}

func GetSkillsForFilterHandler(c *gin.Context) {
	skilkName := c.Query("skill_name")
	skills, total, err := model.SearchSkills(skilkName, 0, 10)
	if err != nil {
		apiUtils.ErrResp(logrus.Fields{
			"error": err.Error(),
		}, "failed to get skills for filter", "Internal Error", c)
		return
	}
	apiUtils.SuccessResp(gin.H{
		"skills": skills,
		"total":  total,
	}, c)
}
