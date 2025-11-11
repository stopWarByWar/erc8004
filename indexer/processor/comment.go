package processor

import (
	"agent_identity/logger"
	"agent_identity/model"
	"time"

	"github.com/sirupsen/logrus"
)

type CommentProcessor struct {
	execBlock          uint64
	execIndex          uint64
	fetchBlockInterval int64

	chainID string

	limit  int
	logger *logger.Logger
}

func NewCommentProcessor(chainID string, startBlock uint64, limit int, fetchBlockInterval int64, _logger *logger.Logger) *CommentProcessor {

	execBlock, execIndex, err := model.GetLatestAgentComment(chainID)
	if err != nil {
		panic(err)
	}

	if startBlock > 0 {
		execBlock = startBlock
		execIndex = 0
	}

	return &CommentProcessor{
		execBlock:          execBlock,
		execIndex:          execIndex,
		limit:              limit,
		fetchBlockInterval: fetchBlockInterval,
		logger:             _logger,
		chainID:            chainID,
	}
}

// todo: new comment schema ID
const commentSchemaID = ""

func (p *CommentProcessor) Process() {
	p.logger.WithFields(logrus.Fields{
		"block": p.execBlock,
		"index": p.execIndex,
	}).Info("start run comment processor")

	for {
		commentAttestations, err := model.GetUnInsertedCommentAttestation(p.execBlock, p.execIndex, p.limit, commentSchemaID)
		if err != nil {
			p.logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("fail to get uninserted comment attestations")
			time.Sleep(10 * time.Second)
			continue
		}

		p.logger.WithFields(logrus.Fields{
			"total": len(commentAttestations),
		}).Debug("get uninserted comment attestations")

		var agentComments []*model.AgentComment
		for _, commentAttestation := range commentAttestations {
			agentComment := p.dealWithCommentAttestation(commentAttestation)
			if agentComment != nil {
				agentComments = append(agentComments, agentComment)
			}
		}

		if len(agentComments) > 0 {
			if err := model.CreateAgentComments(agentComments); err != nil {
				p.logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("fail to insert agent comments")
				time.Sleep(10 * time.Second)
				continue
			}

			p.execBlock = uint64(commentAttestations[len(commentAttestations)-1].Block)
			p.execIndex = uint64(commentAttestations[len(commentAttestations)-1].Index)
		}

		if len(commentAttestations) < p.limit {
			time.Sleep(10 * time.Second)
			continue
		}
	}

}

func (p *CommentProcessor) dealWithCommentAttestation(att *model.Attestation) *model.AgentComment {
	comment, err := DecodeCommentEvent(att.Recipient, att.RawData)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"attestation_id": att.UID,
			"error":          err,
		}).Warn("fail to decode comment event")
		return nil
	}

	agentUID, err := model.GetAgentUID(p.chainID, comment.IdentityRegistry, comment.AgentID.String())
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"error":             err,
			"chain_id":          p.chainID,
			"agent_id":          comment.AgentID.String(),
			"identity_registry": comment.IdentityRegistry,
		}).Error("fail to get agent uid")
		return nil
	}

	return &model.AgentComment{
		CommentAttestationID: att.UID,
		AgentUID:             agentUID,
		ChainID:              p.chainID,
		IdentityRegistry:     comment.IdentityRegistry,
		AgentID:              comment.AgentID.String(),
		Commenter:            comment.Commenter.String(),
		CommentText:          comment.CommentText,
		Score:                comment.Score,
		Timestamps:           uint64(att.Timestamps),
		Block:                uint64(att.Block),
		Index:                uint64(att.Index),
		TxHash:               att.TransactionId,
		Revoked:              att.Revoked,
	}
}
