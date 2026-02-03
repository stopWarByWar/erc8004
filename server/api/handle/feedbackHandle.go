package handle

import (
	serverLogic "agent_identity/server/api/logic"
	serverTypes "agent_identity/server/api/types"
	serverUtils "agent_identity/server/api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetAgentFeedbacksHandler(c *gin.Context) {
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}

	tag1 := c.Query("tag1")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	if pageInt <= 0 {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	if pageSizeInt <= 0 {
		pageSizeInt = 10
	}

	feedbacks, total, err := serverLogic.GetAgentFeedbacksList(agentUID, tag1, pageInt, pageSizeInt)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get agent feedbacks", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"feedbacks": feedbacks,
		"total":     total,
	}, c)

}

func GetFeedbackScoresHandler(c *gin.Context) {
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}

	scores, err := serverLogic.GetAgentScoreForEachTag1(agentUID, 0, 50)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to get feedback scores", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"scores": scores,
	}, c)
}

func UploadFeedbackHandler(c *gin.Context) {
	var request serverTypes.UploadFeedbackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err}, "fail to bind request", "Invalid Request", c)
		return
	}

	feedbackURI, feedbackHash, err := serverLogic.SetFeedback(request)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to set feedback", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"feedbackURI":  feedbackURI,
		"feedbackHash": feedbackHash,
	}, c)
}
