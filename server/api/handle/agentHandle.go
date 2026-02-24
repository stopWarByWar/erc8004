package handle

import (
	serverLogic "agent_identity/server/api/logic"
	serverTypes "agent_identity/server/api/types"
	serverUtils "agent_identity/server/api/utils"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type AgentListRequest struct {
	Name         *string   `json:"name"`
	TrustModel   *[]string `json:"trust_model"`
	Chains       *[]string `json:"chains"`
	Skills       *[]string `json:"skills"`
	X402Support  *bool     `json:"x402_support"`
	Active       *bool     `json:"active"`
	HaveFeedback *bool     `json:"have_feedback"`
	Page         int       `json:"page" binding:"required"`
	PageSize     int       `json:"page_size" binding:"required"`
}

func GetAgentCardListHandler(c *gin.Context) {
	var req AgentListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serverUtils.ErrResp(logrus.Fields{
			"error":   err.Error(),
			"request": req,
		}, "fail to bind request", "Invalid Request", c)
		return
	}

	var agents []*serverTypes.AgentResponse
	var total int64

	agents, total, err := serverLogic.GetAgentListByFilter(req.Page, req.PageSize, req.Name, req.TrustModel, req.Chains, req.Skills, req.X402Support, req.Active, req.HaveFeedback)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{
			"error":   err.Error(),
			"request": req,
		}, "fail to get agent card list", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": req.Page,
	}, c)
}

func GetAgentCardDetailHandler(c *gin.Context) {
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get agent uid", "Invalid Request", c)
		return
	}
	agentCard, err := serverLogic.GetCardResponse(agentUID)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get agent card detail", "Internal Error", c)
		return
	}
	serverUtils.SuccessResp(gin.H{
		"agent": agentCard,
	}, c)
}

const maxMultipartMemory = 20 << 20 // 20MB

func UploadAgentProfileHandler(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(maxMultipartMemory); err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to parse multipart form", "Invalid Request", c)
		return
	}

	var request serverTypes.UploadAgentProfileRequest
	if err := c.ShouldBindWith(&request, binding.FormMultipart); err != nil {
		serverUtils.ErrResp(logrus.Fields{
			"error":   err.Error(),
			"request": request,
		}, "fail to bind request", "Invalid Request", c)
		return
	}

	file, err := c.FormFile("logo")
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to get logo file", "Invalid Request", c)
		return
	}

	src, err := file.Open()
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to open logo file", "Invalid Request", c)
		return
	}
	defer src.Close()

	logoData, err := io.ReadAll(src)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to read logo file", "Invalid Request", c)
		return
	}

	tokenURI, err := serverLogic.UploadAgentProfile(request, logoData)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{"error": err.Error()}, "fail to upload agent profile", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"tokenURI": tokenURI,
	}, c)
}

// SemanticSearchRequest 语义搜索请求体
type SemanticSearchRequest struct {
	Desc         string    `json:"desc" binding:"required"`
	Limit        int       `json:"limit" binding:"required"`
	Threshold    float64   `json:"threshold" binding:"required"`
	TrustModel   *[]string `json:"trust_models"`
	Chains       *[]string `json:"chains"`
	Skills       *[]string `json:"skills"`
	X402Support  *bool     `json:"x402_support"`
	Active       *bool     `json:"active"`
	HaveFeedback *bool     `json:"have_feedback"`
}

func GetAgentCardsSearchBySemanticHandler(c *gin.Context) {
	var req SemanticSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serverUtils.ErrResp(logrus.Fields{
			"error":   err.Error(),
			"request": req,
		}, "fail to bind request", "Invalid Request", c)
		return
	}

	// 基本参数校验
	if req.Limit <= 0 {
		serverUtils.ErrResp(logrus.Fields{
			"error": "limit must be greater than 0",
			"limit": req.Limit,
		}, "invalid request", "Invalid Request", c)
		return
	}

	agents, err := serverLogic.FilterSearchAgentListBySemantic(req.Desc, req.Limit, req.Threshold, req.TrustModel, req.Chains, req.Skills, req.X402Support, req.Active, req.HaveFeedback)
	if err != nil {
		serverUtils.ErrResp(logrus.Fields{
			"error":   err.Error(),
			"request": req,
		}, "fail to get agent card list by semantic", "Internal Error", c)
		return
	}

	serverUtils.SuccessResp(gin.H{
		"agent_list": agents,
	}, c)
}
