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

var ValidationRequestTopic = common.HexToHash("0x530436c3634a98e1e626b0898be2f1e9980cc1bd2a78c07a0aba52d0a48a5059")
var ValidationResponseTopic = common.HexToHash("0xafddf629e874ccc3963b6a888c477bd464a6c8525024fc88759ea3b2326349ae")

type ValidationRegistryProcessor struct {
	execBlock uint64
	execIndex uint64

	identityExecBlockChan <-chan uint64 // 用于接收identity processor发送的execBlock的channel
	identityExecBlock     uint64        // 记录identity processor执行到的区块号

	validationRegistry *abi.ValidationRegistry
	validationAddr     common.Address
	identityAddr       string
	fetchBlockInterval int64
	chainID            string
	logger             *logger.Logger
	ethClient          *ethclient.Client
}

func NewValidationRegistryProcessor(validationAddr, identityAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger, identityExecBlockChan <-chan uint64) *ValidationRegistryProcessor {
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	execBlock, execIndex, err := model.GetLatestValidationRequestAndResponse(chainId.String(), validationAddr)
	if err != nil {
		panic(err)
	}
	validationRegistry, err := abi.NewValidationRegistry(common.HexToAddress(validationAddr), ethClient)
	if err != nil {
		panic(err)
	}

	if startBlock > 0 {
		execBlock = startBlock
		execIndex = 0
	}

	return &ValidationRegistryProcessor{
		execBlock:             execBlock,
		execIndex:             execIndex,
		identityExecBlockChan: identityExecBlockChan,
		validationRegistry:    validationRegistry,
		validationAddr:        common.HexToAddress(validationAddr),
		identityAddr:          identityAddr,
		fetchBlockInterval:    fetchBlockInterval,
		ethClient:             ethClient,
		logger:                _logger,
		chainID:               chainId.String(),
	}
}

func (p *ValidationRegistryProcessor) Process() {
	p.logger.WithFields(logrus.Fields{
		"block": p.execBlock,
		"index": p.execIndex,
	}).Info("start run validation registry processor")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	logExecBlockTicker := time.NewTicker(60 * time.Second)
	defer logExecBlockTicker.Stop()

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
		case identityBlock := <-p.identityExecBlockChan:
			// 接收identity processor发送的execBlock更新
			p.identityExecBlock = identityBlock
		case <-logExecBlockTicker.C:
			p.logger.WithFields(logrus.Fields{
				"identityBlock": p.identityExecBlock,
				"block":         p.execBlock,
				"index":         p.execIndex,
			}).Info("validation registry processor exec block")
		}
	}
}

func (p *ValidationRegistryProcessor) process(currentBlockNum int64) {
	fromBlock := int64(p.execBlock) + 1
loop:
	for {
		toBlock := fromBlock + p.fetchBlockInterval
		events, err := p.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{p.validationAddr},
			Topics:    [][]common.Hash{{ValidationRequestTopic, ValidationResponseTopic}},
		})
		if err != nil {
			p.logger.WithFields(logrus.Fields{
				"from":  fromBlock,
				"to":    toBlock,
				"error": err,
			}).Error("fail to get validation registry event from chain")
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

func (p *ValidationRegistryProcessor) dealWithEvent(e types.Log) error {
	switch e.Topics[0] {
	case ValidationRequestTopic:
		return p.dealWithValidationRequestEvent(e)
	case ValidationResponseTopic:
		return p.dealWithValidationResponseEvent(e)
	default:
		return fmt.Errorf("unknown event topic: %s", e.Topics[0])
	}
}

func (p *ValidationRegistryProcessor) dealWithValidationRequestEvent(e types.Log) error {
	event, err := p.validationRegistry.ParseValidationRequest(e)
	if err != nil {
		return err
	}

	agentUID, err := model.GetAgentUID(p.chainID, p.identityAddr, event.AgentId.String())
	if err != nil {
		if p.identityExecBlock > uint64(e.BlockNumber) && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			p.logger.WithFields(logrus.Fields{
				"agentID":           event.AgentId.String(),
				"identityRegistry":  p.identityAddr,
				"chainID":           p.chainID,
				"error":             err,
				"block":             e.BlockNumber,
				"index":             e.Index,
				"identityExecBlock": p.identityExecBlock,
			}).Error("failed to get agent uid")
			return err
		}
	}

	err = model.InsertValidation(&model.Validation{
		AgentUID:           agentUID,
		ChainID:            p.chainID,
		AgentID:            event.AgentId.String(),
		ValidationRegistry: p.validationAddr.String(),
		ValidatorAddress:   event.ValidatorAddress.String(),
		RequestHash:        common.BytesToHash(event.RequestHash[:]).String(),
		RequestURI:         event.RequestURI,
		BlockNumber:        uint64(e.BlockNumber),
		Index:              uint64(e.Index),
		RequestTxHash:      e.TxHash.String(),
		Timestamps:         uint64(e.BlockTimestamp),
	})

	p.logger.WithFields(logrus.Fields{
		"event":            "validation request",
		"agentID":          event.AgentId.String(),
		"identityRegistry": p.identityAddr,
		"chainID":          p.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with validation request event")

	return err
}

func (p *ValidationRegistryProcessor) dealWithValidationResponseEvent(e types.Log) error {
	event, err := p.validationRegistry.ParseValidationResponse(e)
	if err != nil {
		return err
	}

	err = model.UpdateValidation(&model.Validation{
		RequestHash:    common.BytesToHash(event.RequestHash[:]).String(),
		Response:       int(event.Response),
		ResponseURI:    event.ResponseURI,
		ResponseHash:   common.BytesToHash(event.ResponseHash[:]).String(),
		Tag1:           event.Tag,
		ResponseTxHash: e.TxHash.String(),
		Timestamps:     uint64(e.BlockTimestamp),
		BlockNumber:    uint64(e.BlockNumber),
		Index:          uint64(e.Index),
	})

	p.logger.WithFields(logrus.Fields{
		"event":            "validation response",
		"agentID":          event.AgentId.String(),
		"identityRegistry": p.identityAddr,
		"chainID":          p.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with validation response event")
	return err
}
