// Package handle — unified agent credits endpoint.
//
// GET /agent/identity/credits?uid=<uint> returns the merged identity +
// Commerce (ERC-8183) + Feedback Credit (ERC-8004) view used by the
// redesigned agent detail page. See dev-doc/agent-credits.html for the
// visual mapping between API fields and on-page elements.
package handle

import (
	"strconv"

	serverLogic "agent_identity/server/api/logic"
	serverUtils "agent_identity/server/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetAgentCreditsHandler serves /agent/identity/credits?uid=N.
//
// Errors:
//   - missing/invalid uid → 400 "Invalid Request"
//   - agent not found / db error → 500 "Internal Error"
func GetAgentCreditsHandler(c *gin.Context) {
	uidRaw := c.Query("uid")
	uid, err := strconv.ParseUint(uidRaw, 10, 64)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"uid": uidRaw}, "fail to parse uid", "Invalid Request", c)
		return
	}

	resp, err := serverLogic.GetAgentCredits(uid)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"uid": uid, "error": err}, "fail to get agent credits", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"identity": resp.Identity,
		"commerce": resp.Commerce,
		"feedback": resp.Feedback,
	}, c)
}
