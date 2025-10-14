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
	limit              int
	attestor           string
	logger             *logger.Logger
}

func NewCommentProcessor(startBlock uint64, limit int, fetchBlockInterval int64, attestor string, _logger *logger.Logger) *CommentProcessor {
	execBlock, execIndex, err := model.GetLatestAgentComment()
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
		attestor:           attestor,
		logger:             _logger,
	}
}

const commentSchemaID = "0x4a16aebbff5aeef5b4a39e5e4f900067db897619a8d75c4e9cf4c96c169546dc"

func (p *CommentProcessor) Process() {
	p.logger.WithFields(logrus.Fields{
		"block": p.execBlock,
		"index": p.execIndex,
	}).Info("start run comment processor")

	for {
		commentAttestations, err := model.GetUnInsertedCommentAttestation(p.execBlock, p.execIndex, p.limit, commentSchemaID, p.attestor)
		if err != nil {
			p.logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("fail to get uninserted comment attestations")
			time.Sleep(10 * time.Second)
			continue
		}

		var agentComments []*model.AgentComment
		for _, commentAttestation := range commentAttestations {
			agentComment := p.dealWithCommentAttestation(commentAttestation)
			if agentComment != nil {
				agentComments = append(agentComments, agentComment)
			}
		}

		if len(agentComments) > 0 {
			if err := model.InsertAgentComments(agentComments); err != nil {
				p.logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("fail to insert agent comments")
				time.Sleep(time.Duration(p.fetchBlockInterval) * time.Second)
				continue
			}

			p.execBlock = uint64(commentAttestations[len(commentAttestations)-1].Block)
			p.execIndex = uint64(commentAttestations[len(commentAttestations)-1].Index)
		}

		if len(commentAttestations) < p.limit {
			time.Sleep(time.Duration(p.fetchBlockInterval) * time.Second)
			continue
		}
	}

}

func (p *CommentProcessor) dealWithCommentAttestation(att *model.Attestation) *model.AgentComment {
	comment, err := DecodeCommentEvent(att.RawData)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("fail to decode comment event")
		return nil
	}

	return &model.AgentComment{
		CommentAttestationID: att.UID,
		Commenter:            comment.Commenter.Hex(),
		AgentClientID:        comment.AgentClientId.String(),
		AgentServerID:        comment.AgentServerId.String(),
		CommentText:          comment.Comment,
		Score:                comment.Score,
		Timestamps:           uint64(att.Timestamps),
		NewComment:           true,
		BlockNumber:          uint64(att.Block),
		Index:                uint64(att.Index),
		TxHash:               att.TransactionId,
		IsAuthorized:         comment.IsAuthorized,
	}
}
