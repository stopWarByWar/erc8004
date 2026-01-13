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
var ValidationResponseTopic = common.HexToHash("0xf224d3d5ad74301be48e4d51ca5f1b24c7946875887327585becc59165297dcf")

type ValidationRegistryProcessor struct {
	execBlock uint64
	execIndex uint64

	validationRegistry *abi.ValidationRegistry
	validationAddr     common.Address

	fetchBlockInterval int64
	chainID            string
	logger             *logger.Logger
	ethClient          *ethclient.Client
}

func NewValidationRegistryProcessor(validationAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger) *ValidationRegistryProcessor {
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
		execBlock:          execBlock,
		execIndex:          execIndex,
		validationRegistry: validationRegistry,
		validationAddr:     common.HexToAddress(validationAddr),
		fetchBlockInterval: fetchBlockInterval,
		ethClient:          ethClient,
		logger:             _logger,
		chainID:            chainId.String(),
	}
}

func (p *ValidationRegistryProcessor) Process() {
	p.logger.WithFields(logrus.Fields{
		"block": p.execBlock,
		"index": p.execIndex,
	}).Info("start run validation registry processor")

	for {
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
		time.Sleep(20 * time.Second)
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

	agentUID, err := model.GetAgentUID(p.chainID, p.validationAddr.Hex(), event.AgentId.String())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return model.InsertValidation(&model.Validation{
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
}

func (p *ValidationRegistryProcessor) dealWithValidationResponseEvent(e types.Log) error {
	event, err := p.validationRegistry.ParseValidationResponse(e)
	if err != nil {
		return err
	}

	//todo: update
	return model.UpdateValidation(&model.Validation{
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
}
