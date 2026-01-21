package processor

//NewFeedback
//FeedbackRevoked
//ResponseAppended

import (
	"agent_identity/abi"
	"agent_identity/model"
	"errors"
	"fmt"
	"math/big"
	"time"

	"agent_identity/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var NewFeedbackTopic = common.HexToHash("0x801d7d4264128f6f43835850f1fabb91902c3543c3738f1f62fbf7e9fd80531d")
var ResponseAppendedTopic = common.HexToHash("0xb1c6be0b5b8aef6539e2fac0fd131a2faa7b49edf8e505b5eb0ad487d56051d4")
var FeedbackRevokedTopic = common.HexToHash("0x25156fd3288212246d8b008d5921fde376c71ed14ac2e072a506eb06fde6d09d")

type ReputationProcessor struct {
	execBlock             uint64
	execIndex             uint64
	identityExecBlockChan <-chan uint64
	identityExecBlock     uint64

	reputationRegistry *abi.ReputationRegistry
	reputationAddr     common.Address
	identityAddr       string

	fetchBlockInterval int64
	chainID            string
	logger             *logger.Logger
	ethClient          *ethclient.Client
}

func NewReputationProcessor(reputationAddr, identityAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger, identityExecBlockChan <-chan uint64) *ReputationProcessor {
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
		execBlock:             execBlock,
		execIndex:             execIndex,
		identityExecBlockChan: identityExecBlockChan,
		reputationRegistry:    reputationRegistry,
		reputationAddr:        common.HexToAddress(reputationAddr),
		fetchBlockInterval:    fetchBlockInterval,
		ethClient:             ethClient,
		logger:                _logger,
		chainID:               chainId.String(),
		identityAddr:          identityAddr,
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
		case identityBlock := <-p.identityExecBlockChan:
			// 接收identity processor发送的execBlock更新
			p.identityExecBlock = identityBlock
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

	agentUID, err := model.GetAgentUID(p.chainID, p.identityAddr, newFeedbackEvent.AgentId.String())
	if err != nil {
		if p.identityExecBlock > uint64(e.BlockNumber) && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			p.logger.WithFields(logrus.Fields{
				"error": err,
				"block": e.BlockNumber,
				"index": e.Index,
			}).Error("failed to get agent uid")
			return fmt.Errorf("failed to get agent uid: %w", err)
		}
	}

	blockTimestamp := uint64(e.BlockTimestamp)
	scoreStr := newFeedbackEvent.Value.String()
	scoreDecimals := newFeedbackEvent.ValueDecimals
	score := calculateScore(newFeedbackEvent.Value, newFeedbackEvent.ValueDecimals)

	feedback := &model.Feedback{
		ChainID:            p.chainID,
		AgentUID:           agentUID,
		AgentID:            newFeedbackEvent.AgentId.String(),
		ReputationRegistry: p.reputationAddr.Hex(),
		ClientAddress:      newFeedbackEvent.ClientAddress.String(),
		FeedbackIndex:      newFeedbackEvent.FeedbackIndex,
		Score:              score,
		Value:              scoreStr,
		ValueDecimals:      scoreDecimals,
		Tag1:               newFeedbackEvent.Tag1,
		Tag2:               newFeedbackEvent.Tag2,
		FeedbackURI:        newFeedbackEvent.FeedbackURI,
		FeedbackHash:       common.BytesToHash(newFeedbackEvent.FeedbackHash[:]).String(),
		Endpoint:           newFeedbackEvent.Endpoint,
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

	if err := model.UpdateFeedbackRevoked(p.chainID, feedbackRevokedEvent.AgentId.String(), p.reputationAddr.String(), feedbackRevokedEvent.ClientAddress.String(), feedbackRevokedEvent.FeedbackIndex); err != nil {
		return fmt.Errorf("failed to update feedback revoked: %w", err)
	}
	return nil
}

func (p *ReputationProcessor) dealWithResponseAppendedEvent(e types.Log) error {
	responseAppendedEvent, err := p.reputationRegistry.ParseResponseAppended(e)
	if err != nil {
		return fmt.Errorf("failed to parse response appended event: %w", err)
	}

	feedbackUID, agentUID, err := model.GetFeedbackUIDAndAgentUID(p.chainID, responseAppendedEvent.AgentId.String(), p.reputationAddr.String(), responseAppendedEvent.ClientAddress.String(), responseAppendedEvent.FeedbackIndex)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"error": err,
			"block": e.BlockNumber,
			"index": e.Index,
		}).Error("fail to get feedback uid and agent uid")
		return nil
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
		ResponseURI:   responseAppendedEvent.ResponseURI,
		ResponseHash:  common.BytesToHash(responseAppendedEvent.ResponseHash[:]).String(),
		BlockNumber:   uint64(e.BlockNumber),
		Index:         uint64(e.Index),
		TxHash:        e.TxHash.String(),
		Timestamps:    blockTimestamp,
	}

	if err := model.CreateResponse(p.chainID, response); err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	return nil
}

// todo: to be implemented
func (p *ReputationProcessor) fetchFeedbackAndResponse() {
}
