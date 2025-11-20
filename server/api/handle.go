package api

import (
	agentcard "agent_identity/agentCard"
	"agent_identity/helper"
	"agent_identity/model"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

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

func GetAgentCardListByTrustModelHandler(c *gin.Context) {
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

	trustModel := c.QueryArray("trust_model")
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
		"trust_model": trustModel,
	}).Info("GetCardResponseByTrustModel")

	agents, total, err := GetAgentListByTrustModel(pageInt, pageSizeInt, trustModel)
	if err != nil {
		ErrResp(nil, "fail to get agent card list by trust model", "Internal Error", c)
		return
	}
	SuccessResp(gin.H{
		"agent_list":   agents,
		"total":        total,
		"current_page": pageInt,
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

	agent, err := model.GetAgentByUID(request.UID)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to get agent", "Internal Error", c)
		return
	}

	feedbackAuth := request.FeedbackAuth

	agentID, err := strconv.ParseUint(agent.AgentID, 10, 64)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to parse agent id", "Internal Error", c)
		return
	}
	chainID, err := strconv.ParseInt(feedbackAuth.ChainId, 10, 64)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to parse chain id", "Internal Error", c)
		return
	}
	encodedFeedbackAuth, err := EncodeFeedbackAuth(
		big.NewInt(int64(agentID)),
		common.HexToAddress(feedbackAuth.ClientAddress),
		feedbackAuth.IndexLimit,
		big.NewInt(int64(feedbackAuth.Expiry)),
		big.NewInt(chainID),
		common.HexToAddress(feedbackAuth.IdentityRegistry),
		common.HexToAddress(feedbackAuth.SignerAddress),
	)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to encode feedback auth", "Internal Error", c)
		return
	}

	signerAddr, err := ExtractAddressFromSignature(request.FeedbackAuthSignature, encodedFeedbackAuth)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to extract signer address", "Internal Error", c)
		return
	}

	if signerAddr != common.HexToAddress(feedbackAuth.SignerAddress) {
		ErrResp(logrus.Fields{"error": err}, "invalid signer address", "Invalid Signature", c)
		return
	}

	feedbackAuthData := hex.EncodeToString(append(encodedFeedbackAuth, common.HexToHash(signerAddr.String()).Bytes()...))

	feedback := &types.Feedback{
		AgentRegistry: agent.IdentityRegistry,
		AgentId:       int64(agentID),
		ClientAddress: feedbackAuth.ClientAddress,
		CreatedAt:     strconv.FormatInt(time.Now().Unix(), 10),
		FeedbackAuth:  feedbackAuthData,
		Score:         request.Score,
		Tag1:          &request.Tag1,
		Tag2:          &request.Tag2,
		Context:       &request.Context,
		Task:          &request.Task,
		Capability:    &request.Capability,
	}

	feedbackData, err := json.Marshal(feedback)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to marshal feedback", "Internal Error", c)
		return
	}

	feedbackURI, err := helper.GetHelper().UploadFeedbackToS3(request.FeedbackAuth.ChainId, agent.IdentityRegistry, agent.AgentID, feedbackAuth.ClientAddress, feedbackAuth.IndexLimit, feedbackData)
	if err != nil {
		ErrResp(logrus.Fields{"error": err}, "fail to upload feedback to s3", "Internal Error", c)
		return
	}

	feedbackHash := common.BytesToHash(sha256.New().Sum(feedbackData)).String()

	SuccessResp(gin.H{
		"feedbackURI":  feedbackURI,
		"authData":     feedbackAuthData,
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

	logoURI, err := helper.GetHelper().UploadLogoToS3(request.ChainID, request.IdentityRegistry, request.AgentID, logoData)
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
		Type:        request.Type,
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
				Endpoint: fmt.Sprintf("eip155:%s:%s", request.ChainID, request.AgentWalletAddress),
			},
		},
		Registrations: []types.Registration{
			{
				AgentId:       int64(agentID),
				AgentRegistry: fmt.Sprintf("eip155:%s:%s", request.ChainID, request.IdentityRegistry),
			},
		},
		SupportedTrust: request.SupportedTrust,
	}

	agentProfileData, err := json.Marshal(agentProfile)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to marshal agent profile", "Internal Error", c)
		return
	}

	tokenURI, err := helper.GetHelper().UploadAgentProfileToS3(request.ChainID, request.IdentityRegistry, request.AgentID, agentProfileData)
	if err != nil {
		ErrResp(logrus.Fields{"error": err.Error()}, "fail to upload agent profile to s3", "Internal Error", c)
		return
	}

	SuccessResp(gin.H{
		"tokenURI": tokenURI,
	}, c)
}
