package handle

import (
	serverLogic "agent_identity/server/api/logic"
	serverUtils "agent_identity/server/api/utils"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetAgentValidationListHandler(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	agentUID := c.Query("agent_uid")
	agentUIDInt, err := strconv.ParseUint(agentUID, 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}
	filter := c.Query("filter")

	validationResponsesList, total, err := serverLogic.GetAgentValidationList(agentUIDInt, pageInt, pageSizeInt, filter)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get agent validation list", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"validation_list": validationResponsesList,
		"total":           total,
		"current_page":    pageInt,
	}, c)
}

func GetValidatorListHandler(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	validatorList, total, err := serverLogic.GetValidatorList(pageInt, pageSizeInt)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get validator list", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"validator_list": validatorList,
		"total":          total,
		"current_page":   pageInt,
	}, c)
}

func GetValidatorValidationListHandler(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get page_size", "Invalid Request", c)
		return
	}

	validatorAddress := c.Query("validator_address")
	filter := c.Query("filter")

	validationsList, total, err := serverLogic.GetValidatorValidationList(validatorAddress, pageInt, pageSizeInt, filter)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get validator validation list", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"validation_list": validationsList,
		"total":           total,
		"current_page":    pageInt,
	}, c)
}

func GetValidatorByAddressHandler(c *gin.Context) {
	validatorAddress := c.Query("validator_address")
	validatorAddress = common.HexToAddress(validatorAddress).String()
	validator, err := serverLogic.GetValidatorByAddress(validatorAddress)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get validator by address", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"validator": validator,
	}, c)
}

func GetAgentValidationEvalLatestHandler(c *gin.Context) {
	agentUIDStr := c.Query("agent_uid")
	agentUID, err := strconv.ParseUint(agentUIDStr, 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}
	resp, err := serverLogic.GetLatestAgentValidationEval(agentUID)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get agent validation eval latest", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"report":     resp.Report,
		"dimensions": resp.Dimensions,
	}, c)
}
