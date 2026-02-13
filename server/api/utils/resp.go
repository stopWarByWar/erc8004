package api

import (
	"net/http"
	"time"

	"agent_identity/config"
	"agent_identity/logger"
	"agent_identity/model"
	"agent_identity/server/api/types"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var _logger *logger.Logger
var leaderboardInfo types.LeaderboardInfo

func Init(nlogger *logger.Logger) {
	_logger = nlogger
}

func ErrResp(errorInfos logrus.Fields, msg, reply string, c *gin.Context) {
	if errorInfos != nil {
		_logger.WithFields(errorInfos).Error(msg)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   1,
		"status": "failed",
		"msg":    reply,
	})
}

func SuccessResp(data gin.H, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"status": "success",
		"msg":    "",
		"data":   data,
	})
}

func SuccessRespWithMsg(code int, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"status": "success",
		"msg":    msg,
	})
}

func UpdateGeneralInfo() {
	for {
		config.UpdateFilterInfo()
		UpdateLeaderboardInfo()
		time.Sleep(5 * time.Minute)
	}
}

func UpdateLeaderboardInfo() {
	filterInfo := config.GetFilterInfo()
	for _, chainInfo := range filterInfo.Networks {
		leaderboardInfo.AgentAmount += int64(chainInfo.AgentAmount)
	}

	leaderboardInfo.NetworkAmount = int64(len(filterInfo.Networks))

	feedbackAmount, err := model.GetAgentAmountWithFeedback()
	if err != nil {
		_logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get agent amount with feedback")
		return
	}
	leaderboardInfo.FeedbackAmount = feedbackAmount

	agentAmountWithIn7Days, err := model.GetAgentAmountWithIn7Days()
	if err != nil {
		_logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get agent amount with in 7 days")
		return
	}
	leaderboardInfo.AgentAmountWithIn7Days = agentAmountWithIn7Days

	newCreatedAgents, err := model.GetAgentListBy(0, 10, []int8{-1, 0, 0})
	if err != nil {
		_logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get new created agents")
		return
	}
	leaderboardInfo.Leaderboard = append(leaderboardInfo.Leaderboard, types.LeaderboardAgentInfo{
		Name: "Newest Agents",
		Data: formatSimpleAgentInfo(newCreatedAgents),
	})

	trendingAgents, err := model.GetAgentListBy(0, 10, []int8{0, -1, 0})
	if err != nil {
		_logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get trending agents")
		return
	}

	leaderboardInfo.Leaderboard = append(leaderboardInfo.Leaderboard, types.LeaderboardAgentInfo{
		Name: "Trending Agents",
		Data: formatSimpleAgentInfo(trendingAgents),
	})

	agentsWithNewestFeedback, err := model.GetAgentListBy(0, 10, []int8{0, 0, -1})
	if err != nil {
		_logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get agents with newest feedback")
		return
	}
	leaderboardInfo.Leaderboard = append(leaderboardInfo.Leaderboard, types.LeaderboardAgentInfo{
		Name: "Active Agents",
		Data: formatSimpleAgentInfo(agentsWithNewestFeedback),
	})
}

func GetLeaderboardInfo() types.LeaderboardInfo {
	return leaderboardInfo
}

func formatSimpleAgentInfo(agents []*model.Agent) []types.SimpleAgentInfo {
	var simpleAgentInfos []types.SimpleAgentInfo
	for _, agent := range agents {
		chainInfo, ok := config.GetChainInfo(agent.ChainID)
		if !ok {
			continue
		}
		simpleAgentInfos = append(simpleAgentInfos, types.SimpleAgentInfo{
			UID:              agent.UID,
			AgentID:          agent.AgentID,
			AgentName:        agent.Name,
			AgentDescription: agent.Description,
			ChainID:          agent.ChainID,
			ChainName:        chainInfo.ChainName,
			ChainLogo:        chainInfo.ChainLogo,
		})
	}
	return simpleAgentInfos
}
