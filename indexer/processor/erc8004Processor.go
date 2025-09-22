package processor

import (
	"agent_identity/abi"
	"agent_identity/agentcard"
	"agent_identity/model"
	"context"
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

type Processor struct {
	execBlock          uint64
	execIndex          uint64
	identityRegistry   *abi.IdentityRegistry
	reputationRegistry *abi.ReputationRegistry
	validationRegistry *abi.ValidationRegistry
	comment            *abi.Comment
	identityAddr       common.Address
	reputationAddr     common.Address
	validationAddr     common.Address
	commentAddr        common.Address

	fetchBlockInterval int64
	ethClient          *ethclient.Client
	logger             *logger.Logger
}

func NewProcessor(identityAddr, reputationAddr, validationAddr, commentAddr, rpcURL string, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger) *Processor {
	execBlock, execIndex, err := model.GetLatestAgentRegistry()
	if err != nil {
		panic(err)
	}

	if startBlock > 0 {
		execBlock = startBlock
		execIndex = 0
	}

	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}

	comment, err := abi.NewComment(common.HexToAddress(commentAddr), ethClient)
	if err != nil {
		panic(err)
	}

	reputation, err := abi.NewReputationRegistry(common.HexToAddress(reputationAddr), ethClient)
	if err != nil {
		panic(err)
	}

	validation, err := abi.NewValidationRegistry(common.HexToAddress(validationAddr), ethClient)
	if err != nil {
		panic(err)
	}

	identity, err := abi.NewIdentityRegistry(common.HexToAddress(identityAddr), ethClient)
	if err != nil {
		panic(err)
	}

	return &Processor{
		execBlock:          execBlock,
		execIndex:          execIndex,
		identityRegistry:   identity,
		reputationRegistry: reputation,
		validationRegistry: validation,
		comment:            comment,
		identityAddr:       common.HexToAddress(identityAddr),
		reputationAddr:     common.HexToAddress(reputationAddr),
		validationAddr:     common.HexToAddress(validationAddr),
		commentAddr:        common.HexToAddress(commentAddr),
		ethClient:          ethClient,
		fetchBlockInterval: fetchBlockInterval,
		logger:             _logger,
	}
}

var ctx = context.Background()

func (idx *Processor) Process() {
	idx.logger.WithFields(logrus.Fields{
		"block": idx.execBlock,
		"index": idx.execIndex,
	}).Info("start run identity registry processor")

	for {
		currentBlock, err := idx.ethClient.BlockNumber(ctx)
		if err != nil {
			idx.logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("fail to get current block num")
			time.Sleep(5 * time.Second)
		}
		if idx.execBlock < uint64(currentBlock) {
			idx.process(int64(currentBlock))
		}
		time.Sleep(20 * time.Second)
	}
}

var AgentRegisteredTopic = common.HexToHash("0xaef1fcdf962a03943b677ae3b92ef1b34c972aeef794ed6e9ba5f599b461883b")

func (idx *Processor) process(currentBlockNum int64) {
	fromBlock := int64(idx.execBlock)

loop:
	for {
		toBlock := fromBlock + idx.fetchBlockInterval
		registryIdentityEvents, err := idx.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{idx.identityAddr},
			Topics:    [][]common.Hash{{AgentRegisteredTopic}},
		})
		if err != nil {
			idx.logger.WithFields(logrus.Fields{
				"from":  fromBlock,
				"to":    toBlock,
				"error": err,
			}).Error("fail to get schemaRegistry event from chain")
			time.Sleep(10 * time.Second)
			continue loop
		}

		for _, e := range registryIdentityEvents {
			if uint64(idx.execBlock) > e.BlockNumber {
				continue
			} else if uint64(idx.execBlock) == e.BlockNumber && idx.execIndex >= uint64(e.Index) {
				continue
			}

			if err := idx.dealWithEvent(e); err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error": err,
					"block": e.BlockNumber,
					"index": e.Index,
				}).Error("fail to deal with event")
			}

			idx.execBlock = uint64(e.BlockNumber)
			idx.execIndex = uint64(e.Index)
		}

		if toBlock < currentBlockNum {
			fromBlock = toBlock
			continue loop
		} else {
			idx.execBlock = uint64(currentBlockNum)
			idx.execIndex = 0
			break
		}
	}
}

func (idx *Processor) dealWithEvent(e types.Log) error {
	switch e.Topics[0] {
	case AgentRegisteredTopic:
		return idx.dealWithAgentRegisteredEvent(e)
	default:
		return fmt.Errorf("unknown event topic: %s", e.Topics[0])
	}
}

func (idx *Processor) dealWithAgentRegisteredEvent(e types.Log) error {
	agentRegisteredEvent, err := idx.identityRegistry.ParseAgentRegistered(e)
	if err != nil {
		return err
	}

	agentCards, err := agentcard.GetAgentCard(agentRegisteredEvent.AgentAddress.String(), agentRegisteredEvent.AgentId.String(), agentRegisteredEvent.AgentDomain)
	if err != nil {
		return err
	}

	for _, agentCard := range agentCards {
		if err := model.InsertAgentCard(*agentCard); err != nil {
			return err
		}
	}

	return model.CreateAgentRegistry(&model.AgentRegistry{
		AgentID:      agentRegisteredEvent.AgentId.String(),
		AgentAddress: agentRegisteredEvent.AgentAddress.String(),
		AgentDomain:  agentRegisteredEvent.AgentDomain,
		BlockNumber:  uint64(e.BlockNumber),
		Index:        uint64(e.Index),
		TxHash:       e.TxHash.String(),
		Timestamps:   uint64(time.Now().Unix()),
	})
}
