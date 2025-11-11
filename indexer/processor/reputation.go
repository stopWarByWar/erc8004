package processor

//NewFeedback
//FeedbackRevoked
//ResponseAppended

import (
	"agent_identity/abi"
	"agent_identity/model"
	"fmt"
	"math/big"
	"time"

	"agent_identity/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var NewFeedbackTopic = common.HexToHash("")
var ResponseAppendedTopic = common.HexToHash("")
var FeedbackRevokedTopic = common.HexToHash("")

type ReputationProcessor struct {
	execBlock uint64
	execIndex uint64

	reputationRegistry *abi.ReputationRegistry
	reputationAddr     common.Address

	fetchBlockInterval int64
	chainID            string
	logger             *logger.Logger
	ethClient          *ethclient.Client
}

func NewReputationProcessor(reputationAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger) *ReputationProcessor {
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	execBlock, execIndex, err := model.GetLatestFeedbackAndResponse(chainId.String(), reputationAddr)
	if err != nil {
		panic(err)
	}
	reputationRegistry, err := abi.NewReputationRegistry(common.HexToAddress(reputationAddr), ethClient)
	if err != nil {
		panic(err)
	}

	if startBlock > 0 {
		execBlock = startBlock
		execIndex = 0
	}

	return &ReputationProcessor{
		execBlock:          execBlock,
		execIndex:          execIndex,
		reputationRegistry: reputationRegistry,
		reputationAddr:     common.HexToAddress(reputationAddr),
		fetchBlockInterval: fetchBlockInterval,
		ethClient:          ethClient,
		logger:             _logger,
		chainID:            chainId.String(),
	}
}

func (p *ReputationProcessor) Process() {
	p.logger.WithFields(logrus.Fields{
		"block": p.execBlock,
		"index": p.execIndex,
	}).Info("start run reputation registry processor")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	fetchFeedbackAndResponseTicker := time.NewTicker(20 * time.Second)
	defer fetchFeedbackAndResponseTicker.Stop()

	for {
		select {
		case <-ticker.C:
			currentBlock, err := p.ethClient.BlockNumber(ctx)
			if err != nil {
				p.logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("fail to get current block num")
				continue
			}
			if p.execBlock < uint64(currentBlock) {
				p.process(int64(currentBlock))
			}
		case <-fetchFeedbackAndResponseTicker.C:
			p.fetchFeedbackAndResponse()
		}
	}
}

func (p *ReputationProcessor) process(currentBlockNum int64) {
	fromBlock := int64(p.execBlock) + 1
loop:
	for {
		toBlock := fromBlock + p.fetchBlockInterval
		events, err := p.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{p.reputationAddr},
			Topics:    [][]common.Hash{{NewFeedbackTopic, ResponseAppendedTopic, FeedbackRevokedTopic}},
		})
		if err != nil {
			p.logger.WithFields(logrus.Fields{
				"from":  fromBlock,
				"to":    toBlock,
				"error": err,
			}).Error("fail to get reputation registry event from chain")
			time.Sleep(10 * time.Second)
			continue loop
		}

		for _, e := range events {
			if uint64(p.execBlock) > e.BlockNumber {
				continue
			} else if uint64(p.execBlock) == e.BlockNumber && p.execIndex >= uint64(e.Index) {
				continue
			}

			if err := p.dealWithEvent(e); err != nil {
				p.logger.WithFields(logrus.Fields{
					"error": err,
					"block": e.BlockNumber,
					"index": e.Index,
				}).Error("fail to deal with event")
				return
			}

			p.execBlock = uint64(e.BlockNumber)
			p.execIndex = uint64(e.Index)
		}

		if toBlock < currentBlockNum {
			fromBlock = toBlock
			continue loop
		} else {
			p.execBlock = uint64(currentBlockNum)
			p.execIndex = 0
			return
		}
	}
}

func (p *ReputationProcessor) dealWithEvent(e types.Log) error {
	switch e.Topics[0] {
	case NewFeedbackTopic:
		return p.dealWithNewFeedbackEvent(e)
	case ResponseAppendedTopic:
		return p.dealWithResponseAppendedEvent(e)
	case FeedbackRevokedTopic:
		return p.dealWithFeedbackRevokedEvent(e)
	default:
		return fmt.Errorf("unknown event topic: %s", e.Topics[0])
	}
}

func (p *ReputationProcessor) dealWithNewFeedbackEvent(e types.Log) error {
	newFeedbackEvent, err := p.reputationRegistry.ParseNewFeedback(e)
	if err != nil {
		return fmt.Errorf("failed to parse new feedback event: %w", err)
	}

	agentUID, err := model.GetAgentUID(p.chainID, p.reputationAddr.Hex(), newFeedbackEvent.AgentId.String())
	if err != nil {
		return fmt.Errorf("failed to get agent uid: %w", err)
	}

	blockTimestamp := uint64(e.BlockTimestamp)

	feedback := &model.Feedback{
		ChainID:            p.chainID,
		AgentUID:           agentUID,
		AgentID:            newFeedbackEvent.AgentId.String(),
		ReputationRegistry: p.reputationAddr.Hex(),
		ClientAddress:      newFeedbackEvent.ClientAddress.String(),
		FeedbackIndex:      newFeedbackEvent.FeedbackIndex,
		Score:              newFeedbackEvent.Score,
		Tag1:               common.BytesToHash(newFeedbackEvent.Tag1[:]).String(),
		Tag2:               common.BytesToHash(newFeedbackEvent.Tag2[:]).String(),
		FeedbackURI:        newFeedbackEvent.FeedbackUri,
		FeedbackHash:       common.BytesToHash(newFeedbackEvent.FeedbackHash[:]).String(),
		Revoked:            false,
		BlockNumber:        uint64(e.BlockNumber),
		Index:              uint64(e.Index),
		TxHash:             e.TxHash.String(),
		Timestamps:         blockTimestamp,
	}

	if err := model.CreateFeedback(feedback); err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	return nil
}

func (p *ReputationProcessor) dealWithFeedbackRevokedEvent(e types.Log) error {
	feedbackRevokedEvent, err := p.reputationRegistry.ParseFeedbackRevoked(e)
	if err != nil {
		return fmt.Errorf("failed to parse feedback revoked event: %w", err)
	}

	if err := model.UpdateFeedbackRevoked(p.chainID, feedbackRevokedEvent.AgentId.String(), feedbackRevokedEvent.ClientAddress.String(), feedbackRevokedEvent.FeedbackIndex); err != nil {
		return fmt.Errorf("failed to update feedback revoked: %w", err)
	}
	return nil
}

func (p *ReputationProcessor) dealWithResponseAppendedEvent(e types.Log) error {
	responseAppendedEvent, err := p.reputationRegistry.ParseResponseAppended(e)
	if err != nil {
		return fmt.Errorf("failed to parse response appended event: %w", err)
	}

	feedbackUID, agentUID, err := model.GetFeedbackUIDAndAgentUID(p.chainID, responseAppendedEvent.AgentId.String(), responseAppendedEvent.ClientAddress.String(), responseAppendedEvent.FeedbackIndex)
	if err != nil {
		return fmt.Errorf("failed to get feedback uid and agent uid: %w", err)
	}

	blockTimestamp := uint64(e.BlockTimestamp)

	response := &model.Response{
		ChainID:       p.chainID,
		AgentUID:      agentUID,
		FeedbackUID:   feedbackUID,
		AgentID:       responseAppendedEvent.AgentId.String(),
		ClientAddress: responseAppendedEvent.ClientAddress.String(),
		FeedbackIndex: responseAppendedEvent.FeedbackIndex,
		Responder:     responseAppendedEvent.Responder.String(),
		ResponseURI:   responseAppendedEvent.ResponseUri,
		ResponseHash:  common.BytesToHash(responseAppendedEvent.ResponseHash[:]).String(),
		BlockNumber:   uint64(e.BlockNumber),
		Index:         uint64(e.Index),
		TxHash:        e.TxHash.String(),
		Timestamps:    blockTimestamp,
	}

	if err := model.CreateResponse(response); err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	return nil
}

// todo: to be implemented
func (p *ReputationProcessor) fetchFeedbackAndResponse() {
}
