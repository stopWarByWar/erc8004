package logic

import (
	"agent_identity/helper"
	"agent_identity/model"
	"agent_identity/server/api/types"
	serverTypes "agent_identity/server/api/types"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func GetAgentFeedbacksList(agentUID uint64, page, pageSize int) ([]*model.FeedbackResp, int64, error) {
	feedbacks, total, err := model.GetFeedbacksByAgentUID(agentUID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get agent feedbacks list: %v", err)
	}
	return feedbacks, total, nil
}

func SetFeedback(request serverTypes.UploadFeedbackRequest) (string, string, error) {
	agent, err := model.GetAgentByUID(request.UID)
	if err != nil {
		return "", "", fmt.Errorf("fail to get agent: %w", err)
	}

	agentID, err := strconv.ParseUint(agent.AgentID, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("fail to parse agent id:%s, error: %w", agent.AgentID, err)
	}

	agentRegistry := common.HexToAddress(agent.IdentityRegistry).String()
	clientAddress := common.HexToAddress(request.ClientAddress).String()

	feedback := &types.Feedback{
		AgentRegistry: fmt.Sprintf("eip155:%s:%s", agent.ChainID, agentRegistry),
		AgentId:       int64(agentID),
		ClientAddress: fmt.Sprintf("eip155:%s:%s", agent.ChainID, clientAddress),
		CreatedAt:     strconv.FormatInt(time.Now().Unix(), 10),
		Score:         request.Score,
		Tag1:          request.Tag1,
		Tag2:          request.Tag2,
		Context:       request.Context,
		Task:          request.Task,
		Capability:    request.Capability,
		Endpoint:      request.Endpoint,
		Domain:        request.Domain,
		Name:          request.Name,
	}

	if request.ProofOfPayment != nil {
		feedback.ProofOfPayment = &types.ProofOfPayment{
			FromAddress: common.HexToAddress(request.ProofOfPayment.FromAddress).String(),
			ToAddress:   common.HexToAddress(request.ProofOfPayment.ToAddress).String(),
			ChainId:     agent.ChainID,
			TxHash:      request.ProofOfPayment.TxHash,
		}
	}

	feedbackData, err := json.Marshal(feedback)
	if err != nil {
		return "", "", fmt.Errorf("fail to marshal feedback: %w", err)
	}

	feedbackURI, err := helper.GetHelper().UploadFeedbackToS3(agent.ChainID, agentRegistry, agent.AgentID, clientAddress, feedbackData)
	if err != nil {
		return "", "", fmt.Errorf("fail to upload feedback to s3: %w", err)
	}

	feedbackHash := common.BytesToHash(sha256.New().Sum(feedbackData)).String()

	return feedbackURI, feedbackHash, nil
}
