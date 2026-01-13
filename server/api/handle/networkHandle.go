package handle

import (
	serverLogic "agent_identity/server/api/logic"
	apiUtils "agent_identity/server/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetNetworkListHandler(c *gin.Context) {
	networkList, err := serverLogic.GetNetworkList()
	if err != nil {
		apiUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get network list", "Internal Error", c)
		return
	}
	apiUtils.SuccessResp(gin.H{
		"network_list": networkList,
	}, c)
}
