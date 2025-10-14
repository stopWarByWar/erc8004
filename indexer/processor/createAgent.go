package processor

import (
	"agent_identity/abi"
	agentcard "agent_identity/agentCard"
	"agent_identity/model"
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"agent_identity/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type CreateAgentProcessor struct {
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

func NewCreateAgentProcessor(identityAddr, reputationAddr, validationAddr, commentAddr, rpcURL string, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger) *CreateAgentProcessor {
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

	return &CreateAgentProcessor{
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

func (idx *CreateAgentProcessor) Process() {
	idx.logger.WithFields(logrus.Fields{
		"block": idx.execBlock,
		"index": idx.execIndex,
	}).Info("start run identity registry processor")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	fetchAgentCardTicker := time.NewTicker(20 * time.Second)
	defer fetchAgentCardTicker.Stop()

	for {
		select {
		case <-ticker.C:
			currentBlock, err := idx.ethClient.BlockNumber(ctx)
			if err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("fail to get current block num")
				continue
			}
			if idx.execBlock < uint64(currentBlock) {
				idx.process(int64(currentBlock))
			}
		case <-fetchAgentCardTicker.C:
			idx.setAgentCardInserted()
		}
	}
}

var AgentRegisteredTopic = common.HexToHash("0xaef1fcdf962a03943b677ae3b92ef1b34c972aeef794ed6e9ba5f599b461883b")
var AuthFeedbackTopic = common.HexToHash("")

func (idx *CreateAgentProcessor) process(currentBlockNum int64) {
	fromBlock := int64(idx.execBlock) + 1

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

		authFeedbackEvents, err := idx.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{idx.reputationAddr},
			Topics:    [][]common.Hash{{AuthFeedbackTopic}},
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

		events := sortEvents(registryIdentityEvents, authFeedbackEvents)

		for _, e := range events {
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
				return
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
			return
		}
	}
}

func (idx *CreateAgentProcessor) dealWithEvent(e types.Log) error {
	switch e.Topics[0] {
	case AgentRegisteredTopic:
		return idx.dealWithAgentRegisteredEvent(e)
	case AuthFeedbackTopic:
		return idx.dealWithAuthFeedbackEvent(e)
	default:
		return fmt.Errorf("unknown event topic: %s", e.Topics[0])
	}
}

func (idx *CreateAgentProcessor) dealWithAgentRegisteredEvent(e types.Log) error {
	agentRegisteredEvent, err := idx.identityRegistry.ParseAgentRegistered(e)
	if err != nil {
		return fmt.Errorf("failed to parse agent registered event: %w", err)
	}

	// 创建 agent registry 记录
	registry := &model.AgentRegistry{
		AgentID:      agentRegisteredEvent.AgentId.String(),
		AgentAddress: agentRegisteredEvent.AgentAddress.String(),
		AgentDomain:  agentRegisteredEvent.AgentDomain,
		BlockNumber:  uint64(e.BlockNumber),
		Index:        uint64(e.Index),
		TxHash:       e.TxHash.String(),
		Timestamps:   uint64(time.Now().Unix()),
	}

	if err := model.CreateAgentRegistry(registry); err != nil {
		return fmt.Errorf("failed to create agent registry: %w", err)
	}

	return nil
}

func (idx *CreateAgentProcessor) dealWithAuthFeedbackEvent(e types.Log) error {
	authFeedbackEvent, err := idx.reputationRegistry.ParseAuthFeedback(e)
	if err != nil {
		return fmt.Errorf("failed to parse auth feedback event: %w", err)
	}

	agentClient, err := idx.identityRegistry.GetAgent(&bind.CallOpts{}, authFeedbackEvent.AgentClientId)
	if err != nil {
		return fmt.Errorf("failed to get agent client address: %w", err)
	}
	agentServer, err := idx.identityRegistry.GetAgent(&bind.CallOpts{}, authFeedbackEvent.AgentServerId)
	if err != nil {
		return fmt.Errorf("failed to get agent server address: %w", err)
	}

	if err := model.InsertAuthFeedback(
		model.AuthFeedback{
			AuthFeedbackID:     common.BytesToHash(authFeedbackEvent.FeedbackAuthId[:]).String(),
			AgentClientID:      authFeedbackEvent.AgentClientId.String(),
			AgentServerID:      authFeedbackEvent.AgentServerId.String(),
			AgentClientAddress: agentClient.AgentAddress.String(),
			AgentServerAddress: agentServer.AgentAddress.String(),

			BlockNumber: int64(e.BlockNumber),
			Index:       int64(e.Index),
			TxHash:      e.TxHash.String(),
		},
	); err != nil {
		return fmt.Errorf("failed to insert auth feedback: %w", err)
	}

	return nil
}

func (idx *CreateAgentProcessor) setAgentCardInserted() {
	var limit = 100
	for {
		agentRegistries, err := model.GetUnInsertedAgentRegistry(limit)
		if err != nil {
			idx.logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to get un inserted agent registry")
			return
		}

		if len(agentRegistries) == 0 {
			break
		}

		var insertedAgentIDs []string

		for _, agentRegistry := range agentRegistries {
			agentCards, err := agentcard.GetAgentCard(agentRegistry.AgentAddress, agentRegistry.AgentID, agentRegistry.AgentDomain)
			if err != nil {
				idx.logger.WithFields(logrus.Fields{
					"agent_address": agentRegistry.AgentAddress,
					"agent_id":      agentRegistry.AgentID,
					"domain":        agentRegistry.AgentDomain,
					"error":         err,
				}).Error("failed to get agent card from domain")

				if strings.Contains(err.Error(), "invalid domain") {
					insertedAgentIDs = append(insertedAgentIDs, agentRegistry.AgentID)
				}
			}

			var insertErrors []error
			for i, card := range agentCards {
				if err := model.InsertAgentCard(*card); err != nil {
					insertErrors = append(insertErrors, fmt.Errorf("failed to insert agent card %d: %w", i, err))
					idx.logger.WithFields(logrus.Fields{
						"agent_id": card.AgentID,
						"error":    err,
					}).Error("failed to insert agent card")
					continue
				}
				insertedAgentIDs = append(insertedAgentIDs, card.AgentID)
			}

			if len(insertErrors) > 0 {
				idx.logger.WithFields(logrus.Fields{
					"total_cards": len(agentCards),
					"errors":      len(insertErrors),
				}).Error("some agent cards failed to insert")
			}

			if len(insertedAgentIDs) > 0 {
				if err := model.UpdateAgentRegistryInserted(insertedAgentIDs); err != nil {
					idx.logger.WithFields(logrus.Fields{
						"agent_ids": insertedAgentIDs,
						"error":     err,
					}).Error("failed to update agent registry inserted")
					continue
				}
			}
		}

		if len(agentRegistries) < limit {
			return
		}

	}
}
