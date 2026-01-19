package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"agent_identity/config"
	"agent_identity/logger"
	"agent_identity/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var _logger *logger.Logger

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

var generalInfo = make(map[string]any)

func GetGeneralInfo() map[string]any {
	return generalInfo
}

func UpdateGeneralInfo() {
	for {
		agentsAmount, err := model.GetAgentAmountForEachChain()
		if err != nil {
			fmt.Println("failed to get agent amount for each chain", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		total := int64(0)
		for chainId, amount := range agentsAmount {
			chainInfo, ok := config.GetChainInfo(chainId)
			if !ok {
				continue
			}
			chainName := strings.Replace(chainInfo.ChainName, " ", "_", -1)
			if len(chainName) > 0 {
				generalInfo[chainName] = amount
				total += amount
			}
			config.SetChainAgentAmount(chainId, amount)

		}
		generalInfo["total"] = total
		time.Sleep(5 * time.Minute)
	}
}
