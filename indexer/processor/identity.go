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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var RegisteredTopic = common.HexToHash("0xca52e62c367d81bb2e328eb795f7c7ba24afb478408a26c0e201d155c449bc4a")
var UriUpdatedTopic = common.HexToHash("0xb41beef75d9f8d55b985319b459e96f82453580af381391f1ad531eb8f8b5a3a")
var TransferOwnerShipTopic = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

type IdentityProcessor struct {
	execBlock uint64
	execIndex uint64

	ethClient        *ethclient.Client
	identityRegistry *abi.IdentityRegistry
	identityAddr     common.Address

	fetchBlockInterval int64

	chainID string
	logger  *logger.Logger
}

func NewCreateAgentProcessor(identityAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger) *IdentityProcessor {
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		panic(err)
	}

	execBlock, execIndex, err := model.GetLatestAgentRegistry(chainId.String(), identityAddr)
	if err != nil {
		panic(err)
	}

	if startBlock > 0 {
		execBlock = startBlock
		execIndex = 0
	}
	identity, err := abi.NewIdentityRegistry(common.HexToAddress(identityAddr), ethClient)
	if err != nil {
		panic(err)
	}

	return &IdentityProcessor{
		execBlock:          execBlock,
		execIndex:          execIndex,
		identityRegistry:   identity,
		identityAddr:       common.HexToAddress(identityAddr),
		fetchBlockInterval: fetchBlockInterval,
		logger:             _logger,
		chainID:            chainId.String(),
		ethClient:          ethClient,
	}
}

var ctx = context.Background()

func (idx *IdentityProcessor) Process() {
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

func (idx *IdentityProcessor) process(currentBlockNum int64) {
	fromBlock := int64(idx.execBlock) + 1

loop:
	for {
		toBlock := fromBlock + idx.fetchBlockInterval
		events, err := idx.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{idx.identityAddr},
			Topics:    [][]common.Hash{{RegisteredTopic, UriUpdatedTopic, TransferOwnerShipTopic}},
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

func (idx *IdentityProcessor) dealWithEvent(e types.Log) error {
	switch e.Topics[0] {
	case RegisteredTopic:
		return idx.dealWithAgentRegisteredEvent(e)
	case UriUpdatedTopic:
		return idx.dealWithUriUpdatedEvent(e)
	case TransferOwnerShipTopic:
		return idx.dealWithTransferOwnerShipEvent(e)
	default:
		return fmt.Errorf("unknown event topic: %s", e.Topics[0])
	}
}

func (idx *IdentityProcessor) dealWithAgentRegisteredEvent(e types.Log) error {
	agentRegisteredEvent, err := idx.identityRegistry.ParseRegistered(e)
	if err != nil {
		return fmt.Errorf("failed to parse agent registered event: %w", err)
	}

	// 创建 agent registry 记录
	registry := &model.AgentRegistry{
		AgentID:          agentRegisteredEvent.AgentId.String(),
		IdentityRegistry: idx.identityAddr.Hex(),
		Owner:            agentRegisteredEvent.Owner.String(),
		TokenURL:         agentRegisteredEvent.TokenURI,
		ChainID:          idx.chainID,
		BlockNumber:      uint64(e.BlockNumber),
		Index:            uint64(e.Index),
		TxHash:           e.TxHash.String(),
		Timestamps:       uint64(e.BlockTimestamp),
	}

	if err := model.CreateAgentRegistry(registry); err != nil {
		return fmt.Errorf("failed to create agent registry: %w", err)
	}
	return nil
}

func (idx *IdentityProcessor) dealWithUriUpdatedEvent(e types.Log) error {
	event, err := idx.identityRegistry.ParseUriUpdated(e)
	if err != nil {
		return fmt.Errorf("failed to parse auth feedback event: %w", err)
	}

	if err := model.UpdateAgentTokenURL(idx.chainID, idx.identityAddr.Hex(), event.AgentId.String(), event.NewUri, uint64(e.BlockNumber), uint64(e.Index)); err != nil {
		return fmt.Errorf("failed to update agent token url: %w", err)
	}

	return nil
}

func (idx *IdentityProcessor) dealWithTransferOwnerShipEvent(e types.Log) error {
	event, err := idx.identityRegistry.ParseTransfer(e)
	if err != nil {
		return fmt.Errorf("failed to parse transfer owner ship event: %w", err)
	}

	if err := model.TransferOwnerShip(idx.chainID, idx.identityAddr.Hex(), event.TokenId.String(), event.To.String(), uint64(e.BlockNumber), uint64(e.Index)); err != nil {
		return fmt.Errorf("failed to transfer owner ship: %w", err)
	}
	return nil
}

func (idx *IdentityProcessor) setAgentCardInserted() {
	var limit = 100
	for {
		agentRegistries, err := model.GetUnInsertedAgentRegistry(idx.chainID, idx.identityAddr.Hex(), limit)
		if err != nil {
			idx.logger.WithFields(logrus.Fields{
				"error":            err,
				"chainID":          idx.chainID,
				"identityRegistry": idx.identityAddr.Hex(),
			}).Error("failed to get un inserted agent registry")
			return
		}

		if len(agentRegistries) == 0 {
			break
		}

		for _, agentRegistry := range agentRegistries {
			agent, err := agentcard.GetAgentCardFromTokenURL(agentRegistry.Owner, agentRegistry.AgentID, agentRegistry.TokenURL, idx.chainID, idx.identityAddr.Hex(), agentRegistry.Timestamps)
			if err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error":            err,
					"chainID":          idx.chainID,
					"identityRegistry": idx.identityAddr.Hex(),
					"agentID":          agentRegistry.AgentID,
					"tokenURL":         agentRegistry.TokenURL,
				}).Error("failed to get agent card from token url")
				if strings.HasPrefix(err.Error(), "network error:") {
					return
				} else {
					continue
				}
			}

			if err := model.InsertAgentCard(agent); err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error":            err,
					"chainID":          idx.chainID,
					"identityRegistry": idx.identityAddr.Hex(),
					"agentID":          agentRegistry.AgentID,
					"tokenURL":         agentRegistry.TokenURL,
				}).Error("failed to insert agent card")
				return
			}

			if err := model.UpdateAgentRegistryInserted([]string{agentRegistry.AgentID}); err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error":            err,
					"chainID":          idx.chainID,
					"identityRegistry": idx.identityAddr.Hex(),
					"agentID":          agentRegistry.AgentID,
					"tokenURL":         agentRegistry.TokenURL,
				}).Error("failed to update agent registry inserted")
				return
			}
		}

		if len(agentRegistries) < limit {
			return
		}
	}
}
