package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAgentCardListHandler(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	limit := c.Query("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ErrResp(nil, "fail to get limit", "Invalid Request", c)
	}
	agentCardList, err := GetCardList(pageInt, pageSizeInt, limitInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list", "Internal Error", c)
	}

	SuccessResp(gin.H{
		"campaigns": agentCardList,
	}, c)
}

func GetAgentCardDetailHandler(c *gin.Context) {
	agentID := c.Query("agent_id")
	agentCard, err := GetCardResponse(agentID)
	if err != nil {
		ErrResp(nil, "fail to get agent card detail", "Internal Error", c)
	}
	SuccessResp(gin.H{
		"agent_card": agentCard,
	}, c)
}

func GetAgentCardListByTrustModelHandler(c *gin.Context) {
	trustModel := c.QueryArray("trust_model")
	agentCardList, err := GetCardResponseByTrustModel(trustModel)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by trust model", "Internal Error", c)
	}
	SuccessResp(gin.H{
		"agent_card_list": agentCardList,
	}, c)
}
