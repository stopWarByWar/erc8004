package server

import (
	"net/http"

	"agent_identity/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var _logger *logger.Logger

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
