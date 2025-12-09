package api

import (
	agentcard "agent_identity/agentCard"
	"agent_identity/config"
	"agent_identity/helper"
	"agent_identity/model"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"agent_identity/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

var mock = false
var feedbackMock = false

const maxMultipartMemory = 20 << 20 // 20MB

func GetAgentCardListHandler(c *gin.Context) {
	if mock {
		page := c.Query("page")
		pageSize := c.Query("page_size")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			ErrResp(nil, "fail to get page", "Invalid Request", c)
			return
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			ErrResp(nil, "fail to get page_size", "Invalid Request", c)
			return
		}
		SuccessResp(gin.H{
			"agent_list":   mockAgents(pageSizeInt),
			"total":        int64(len(mockAgents(pageSizeInt))),
			"current_page": pageInt,
		}, c)
		return
	}

	page := c.Query("page")
	pageSize := c.Query("page_size")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	resp, total, err := GetAgentList(pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list", "Internal Error", c)
		return
	}

	SuccessResp(gin.H{
		"agent_list":   resp,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

func GetAgentCardDetailHandler(c *gin.Context) {
	if mock {
		SuccessResp(gin.H{
			"agent": mockAgent(),
		}, c)
		return
	}

	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}
	agentCard, err := GetCardResponse(agentUID)
	if err != nil {
		ErrResp(nil, "fail to get agent card detail", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent": agentCard,
	}, c)
}

func GetAgentCardListByFilterHandler(c *gin.Context) {
	if mock {
		page := c.Query("page")
		pageSize := c.Query("page_size")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			ErrResp(nil, "fail to get page", "Invalid Request", c)
			return
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			ErrResp(nil, "fail to get page_size", "Invalid Request", c)
			return
		}
		SuccessResp(gin.H{
			"agent_list":   mockAgents(pageSizeInt),
			"total":        int64(len(mockAgents(pageSizeInt))),
			"current_page": pageInt,
		}, c)
		return
	}

	trustModels := c.QueryArray("trust_model")
	chains := c.QueryArray("chains")

	trustModelIDs := make([]string, 0)
	for _, v := range trustModels {
		for _, id := range strings.Split(v, ",") {
			if id != "" {
				trustModelIDs = append(trustModelIDs, id)
			}
		}
		fmt.Println(v)
	}

	chainIDs := make([]string, 0)
	for _, v := range chains {
		for _, id := range strings.Split(v, ",") {
			if id != "" {
				chainIDs = append(chainIDs, id)
			}
		}
		fmt.Println(v)
	}

	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	_logger.WithFields(logrus.Fields{
		"page":        pageInt,
		"page_size":   pageSizeInt,
		"trust_model": trustModelIDs,
		"chain":       chainIDs,
	}).Info("GetCardResponseByFilter")

	agents, total, err := GetAgentListByFilter(pageInt, pageSizeInt, trustModelIDs, chainIDs)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by filter", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

func GetChainListHandler(c *gin.Context) {
	SuccessResp(gin.H{
		"chain_list": config.ChainList,
	}, c)
}

func GetTrustModelListHandler(c *gin.Context) {
	SuccessResp(gin.H{
		"trust_model_list": []string{agentcard.TrustModelFeedback, agentcard.TrustModelInferenceValidation, agentcard.TrustModelTeeAttestation},
	}, c)
}

func GetAgentCardsSearchBySkillHandler(c *gin.Context) {
	if mock {
		page := c.Query("page")
		pageSize := c.Query("page_size")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			ErrResp(nil, "fail to get page", "Invalid Request", c)
			return
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			ErrResp(nil, "fail to get page_size", "Invalid Request", c)
			return
		}
		SuccessResp(gin.H{
			"agent_list":   mockAgents(pageSizeInt),
			"total":        int64(len(mockAgents(pageSizeInt))),
			"current_page": pageInt,
		}, c)
		return
	}

	skill := c.Query("skill")
	page := c.Query("page")
	pageSize := c.Query("page_size")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	agents, total, err := SearchAgentListBySkill(skill, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by skill", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

func GetAgentCardsSearchByNameHandler(c *gin.Context) {
	if mock {
		page := c.Query("page")
		pageSize := c.Query("page_size")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			ErrResp(nil, "fail to get page", "Invalid Request", c)
			return
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			ErrResp(nil, "fail to get page_size", "Invalid Request", c)
			return
		}
		SuccessResp(gin.H{
			"agent_list":   mockAgents(pageSizeInt),
			"total":        int64(len(mockAgents(pageSizeInt))),
			"current_page": pageInt,
		}, c)
		return
	}
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	agents, total, err := SearchAgentListByName(name, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by name", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

func GetAgentCardsFilterSearchByNameHandler(c *gin.Context) {
	if mock {
		page := c.Query("page")
		pageSize := c.Query("page_size")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			ErrResp(nil, "fail to get page", "Invalid Request", c)
			return
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			ErrResp(nil, "fail to get page_size", "Invalid Request", c)
			return
		}
		SuccessResp(gin.H{
			"agent_list":   mockAgents(pageSizeInt),
			"total":        int64(len(mockAgents(pageSizeInt))),
			"current_page": pageInt,
		}, c)
		return
	}
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}

	trustModels := c.QueryArray("trust_model")
	chains := c.QueryArray("chains")

	trustModelIDs := make([]string, 0)
	for _, v := range trustModels {
		for _, id := range strings.Split(v, ",") {
			if id != "" {
				trustModelIDs = append(trustModelIDs, id)
			}
		}
	}

	chainIDs := make([]string, 0)
	for _, v := range chains {
		for _, id := range strings.Split(v, ",") {
			if id != "" {
				chainIDs = append(chainIDs, id)
			}
		}
	}

	agents, total, err := FilterSearchAgentListByFilter(name, pageInt, pageSizeInt, trustModelIDs, chainIDs)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by name", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
	}, c)
}

func GetAgentCommentsHandler(c *gin.Context) {
	if mock {
		SuccessResp(gin.H{
			"comments": mockComments(),
			"total":    int64(len(mockComments())),
		}, c)
		return
	}
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}

	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	if pageInt <= 0 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	if pageSizeInt <= 0 {
		pageSizeInt = 10
	}
	comments, total, err := model.GetCommentsByAgentUID(agentUID, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent comments", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"comments": comments,
		"total":    total,
	}, c)
}

func GetAgentFeedbacksHandler(c *gin.Context) {
	if feedbackMock {
		SuccessResp(gin.H{
			"feedbacks": mockFeedbacks(),
			"total":     int64(len(mockFeedbacks())),
		}, c)
		return
	}
	agentUID, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		ErrResp(nil, "fail to get agent uid", "Invalid Request", c)
		return
	}
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ErrResp(nil, "fail to get page", "Invalid Request", c)
		return
	}
	if pageInt <= 0 {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ErrResp(nil, "fail to get page_size", "Invalid Request", c)
		return
	}
	if pageSizeInt <= 0 {
		pageSizeInt = 10
	}

	feedbacks, total, err := model.GetFeedbacksByAgentUID(agentUID, pageInt, pageSizeInt)
	if err != nil {
		ErrResp(nil, "fail to get agent feedbacks", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"feedbacks": feedbacks,
		"total":     total,
	}, c)

}

func UploadFeedbackHandler(c *gin.Context) {
	var request UploadFeedbackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to bind request", "Invalid Request", c)
		return
	}

	feedbackURI, feedbackHash, err := SetFeedback(request)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to set feedback", "Internal Error", c)
		return
	}

	SuccessResp(gin.H{
		"feedbackURI":  feedbackURI,
		"feedbackHash": feedbackHash,
	}, c)
}

func UploadAgentProfileHandler(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(maxMultipartMemory); err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to parse multipart form", "Invalid Request", c)
		return
	}

	var request UploadAgentProfileRequest
	if err := c.ShouldBindWith(&request, binding.FormMultipart); err != nil {
		ErrResp(logrus.Fields{
			"error":   err.Error(),
			"request": request,
		}, "fail to bind request", "Invalid Request", c)
		return
	}

	file, err := c.FormFile("logo")
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to get logo file", "Invalid Request", c)
		return
	}
	src, err := file.Open()
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to open logo file", "Internal Error", c)
		return
	}
	defer src.Close()

	logoData, err := io.ReadAll(src)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to read logo file", "Internal Error", c)
		return
	}

	logoURI, err := helper.GetHelper().UploadLogoToS3(request.ChainID, common.HexToAddress(request.IdentityRegistry).String(), request.AgentID, logoData)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to upload logo to s3", "Internal Error", c)
		return
	}

	agentID, err := strconv.ParseUint(request.AgentID, 10, 64)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to parse agent id", "Invalid Request", c)
		return
	}

	agentProfile := &types.AgentProfile{
		Type:        "https://eips.ethereum.org/EIPS/eip-8004#registration-v1",
		Name:        request.Name,
		Description: request.Description,
		Image:       logoURI,
		Endpoints: []types.Endpoint{
			{
				Name:     "A2A",
				Endpoint: request.A2AEndpoint,
			},
			{
				Name:     "agentWallet",
				Endpoint: fmt.Sprintf("eip155:%s:%s", request.ChainID, common.HexToAddress(request.AgentWalletAddress).String()),
			},
		},
		Registrations: []types.Registration{
			{
				AgentId:       int64(agentID),
				AgentRegistry: fmt.Sprintf("eip155:%s:%s", request.ChainID, common.HexToAddress(request.IdentityRegistry).String()),
			},
		},
		SupportedTrust: request.SupportedTrust,
	}

	agentProfileData, err := json.Marshal(agentProfile)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to marshal agent profile", "Internal Error", c)
		return
	}

	tokenURI, err := helper.GetHelper().UploadAgentProfileToS3(request.ChainID, common.HexToAddress(request.IdentityRegistry).String(), request.AgentID, agentProfileData)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to upload agent profile to s3", "Internal Error", c)
		return
	}

	SuccessResp(gin.H{
		"tokenURI": tokenURI,
	}, c)
}
