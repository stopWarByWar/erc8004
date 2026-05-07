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

var NewFeedbackTopic = common.HexToHash("0x6a4a61743519c9d648a14e6493f47dbe3ff1aa29e7785c96c8326a205e58febc")
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

	processInterval := 20 * time.Second
	ticker := time.NewTicker(processInterval)
	defer ticker.Stop()

	fetchFeedbackAndResponseInterval := 20 * time.Second
	fetchFeedbackAndResponseTicker := time.NewTicker(fetchFeedbackAndResponseInterval)
	defer fetchFeedbackAndResponseTicker.Stop()

	logExecBlockInterval := 60 * time.Second
	logExecBlockTicker := time.NewTicker(logExecBlockInterval)
	defer logExecBlockTicker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			currentBlock, err := p.ethClient.BlockNumber(ctx)
			if err != nil {
				p.logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("fail to get current block num")
				ticker.Reset(processInterval)
				continue
			}
			if p.execBlock < uint64(currentBlock) {
				p.process(int64(currentBlock))
			}
			ticker.Reset(processInterval)
		case <-fetchFeedbackAndResponseTicker.C:
			fetchFeedbackAndResponseTicker.Stop()
			p.fetchFeedbackAndResponse()
			fetchFeedbackAndResponseTicker.Reset(fetchFeedbackAndResponseInterval)
		case identityBlock := <-p.identityExecBlockChan:
			// 接收identity processor发送的execBlock更新
			p.identityExecBlock = identityBlock
		case <-logExecBlockTicker.C:
			logExecBlockTicker.Stop()
			p.logger.WithFields(logrus.Fields{
				"identityBlock": p.identityExecBlock,
				"block":         p.execBlock,
				"index":         p.execIndex,
			}).Info("reputation registry processor exec block")
			logExecBlockTicker.Reset(logExecBlockInterval)
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
					"error":  err,
					"block":  e.BlockNumber,
					"index":  e.Index,
					"txHash": e.TxHash.String(),
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
				"agentID":           newFeedbackEvent.AgentId.String(),
				"identityRegistry":  p.identityAddr,
				"chainID":           p.chainID,
				"error":             err,
				"block":             e.BlockNumber,
				"index":             e.Index,
				"identityExecBlock": p.identityExecBlock,
			}).Error("failed to get agent uid")
			return fmt.Errorf("failed to get agent uid: %w", err)
		}
	}

	blockTimestamp := uint64(e.BlockTimestamp)
	scoreStr := newFeedbackEvent.Value.String()
	scoreDecimals := newFeedbackEvent.ValueDecimals
	score, err := calculateScore(newFeedbackEvent.Value, newFeedbackEvent.ValueDecimals)
	if err != nil {
		return fmt.Errorf("failed to calculate score: %w", err)
	}

	feedback := &model.Feedback{
		ChainID:            p.chainID,
		AgentUID:           agentUID,
		AgentID:            newFeedbackEvent.AgentId.String(),
		ReputationRegistry: p.reputationAddr.String(),
		ClientAddress:      newFeedbackEvent.ClientAddress.String(),
		FeedbackIndex:      newFeedbackEvent.FeedbackIndex,
		FormatValue:        score,
		Value:              scoreStr,
		ValueDecimals:      uint(scoreDecimals),
		Tag1:               newFeedbackEvent.Tag1,
		Tag2:               newFeedbackEvent.Tag2,
		FeedbackURI:        newFeedbackEvent.FeedbackURI,
		FeedbackHash:       common.BytesToHash(newFeedbackEvent.FeedbackHash[:]).String(),
		BlockNumber:        uint64(e.BlockNumber),
		Index:              uint64(e.Index),
		TxHash:             e.TxHash.String(),
		Timestamps:         blockTimestamp,
		Endpoint:           newFeedbackEvent.Endpoint,
		IdentityRegistry:   p.identityAddr,
	}

	if err := model.CreateFeedback(feedback); err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event":            "new feedback",
		"agentID":          newFeedbackEvent.AgentId.String(),
		"identityRegistry": p.identityAddr,
		"chainID":          p.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with new feedback event")
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

	p.logger.WithFields(logrus.Fields{
		"event":            "feedback revoked",
		"agentID":          feedbackRevokedEvent.AgentId.String(),
		"identityRegistry": p.identityAddr,
		"chainID":          p.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with feedback revoked event")
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
	p.logger.WithFields(logrus.Fields{
		"event":            "response appended",
		"agentID":          responseAppendedEvent.AgentId.String(),
		"identityRegistry": p.identityAddr,
		"chainID":          p.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with response appended event")

	return nil
}

// todo: to be implemented
func (p *ReputationProcessor) fetchFeedbackAndResponse() {
}
