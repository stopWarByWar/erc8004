package handle

import (
	"agent_identity/server/api/logic"
	apiMock "agent_identity/server/api/mock"
	serverUtils "agent_identity/server/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetPassportHandler(c *gin.Context) {
	uidStr := c.Query("uid")
	if uidStr == "" {
		serverUtils.ErrResp(logrus.Fields{"error": "uid is required"}, "fail to get uid", "Invalid Request", c)
		return
	}

	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to parse uid", "Invalid Request", c)
		return
	}

	if serverUtils.IsMockEnabled() || c.Query("mock") == "true" {
		serverUtils.SetMockHeader(c)
		passport := apiMock.PassportDetail(uid)
		serverUtils.SuccessResp(gin.H{"passport": passport}, c)
		return
	}

	passport, err := logic.GetPassportData(uid)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get passport", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"passport": passport,
	}, c)
}

func GetPassportSharedHandler(c *gin.Context) {
	uidStr := c.Query("uid")
	if uidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"status":  "failed",
			"message": "uid is required",
		})
		return
	}

	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"status":  "failed",
			"message": "invalid uid format",
		})
		return
	}

	if serverUtils.IsMockEnabled() || c.Query("mock") == "true" {
		serverUtils.SetMockHeader(c)
		passport := apiMock.PassportDetail(uid)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"status":  "success",
			"message": "",
			"data": gin.H{
				"passport": passport,
			},
		})
		return
	}

	passport, err := logic.GetPassportData(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"status":  "success",
		"message": "",
		"data": gin.H{
			"passport": passport,
		},
	})
}