package api

import (
	"net/http"
	"time"

	"agent_identity/config"
	"agent_identity/logger"

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

func UpdateGeneralInfo() {
	for {
		config.UpdateGeneralInfo()
		time.Sleep(5 * time.Minute)
	}
}
