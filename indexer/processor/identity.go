package processor

import (
	"agent_identity/abi"
	agentcard "agent_identity/agentCard"
	"agent_identity/model"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"agent_identity/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var RegisteredTopic = common.HexToHash("0xca52e62c367d81bb2e328eb795f7c7ba24afb478408a26c0e201d155c449bc4a")
var UriUpdatedTopic = common.HexToHash("0x3a2c7fffc2cba7582c690e3b82c453ea02a308326a98a3ad7576c606336409fb")
var TransferOwnerShipTopic = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
var SetMetaDataTopic = common.HexToHash("0x2c149ed548c6d2993cd73efe187df6eccabe4538091b33adbd25fafdb8a1468b")

type IdentityProcessor struct {
	execBlock uint64
	execIndex uint64
	mu        sync.RWMutex // 保护execBlock和execIndex的并发访问

	execBlockChan chan<- uint64 // 用于向reputation processor发送execBlock的channel

	ethClient        *ethclient.Client
	identityRegistry *abi.IdentityRegistry
	identityAddr     common.Address

	fetchBlockInterval int64

	chainID string
	logger  *logger.Logger
}

func NewCreateAgentProcessor(identityAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger, execBlockChan chan<- uint64) *IdentityProcessor {
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		panic(err)
	}

	execBlock, execIndex, err := model.GetLatestAgent(chainId.String(), identityAddr)
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
		execBlockChan:      execBlockChan,
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
	idx.mu.RLock()
	execBlock := idx.execBlock
	execIndex := idx.execIndex
	idx.mu.RUnlock()
	idx.logger.WithFields(logrus.Fields{
		"block": execBlock,
		"index": execIndex,
	}).Info("start run identity registry processor")

	processInterval := 20 * time.Second
	ticker := time.NewTicker(1)
	defer ticker.Stop()

	fetchAgentProfileInterval := 20 * time.Second
	fetchAgentCardTicker := time.NewTicker(fetchAgentProfileInterval)
	defer fetchAgentCardTicker.Stop()

	logExecBlockInterval := 60 * time.Second
	logExecBlockTicker := time.NewTicker(2)
	defer logExecBlockTicker.Stop()
	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			currentBlock, err := idx.ethClient.BlockNumber(ctx)
			if err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("fail to get current block num")
				continue
			}
			idx.mu.RLock()
			execBlock := idx.execBlock
			idx.mu.RUnlock()

			if idx.execBlockChan != nil {
				select {
				case idx.execBlockChan <- execBlock:
				default:
				}
			}

			if execBlock < uint64(currentBlock) {
				idx.process(int64(currentBlock))
			}
			ticker.Reset(processInterval)
		case <-fetchAgentCardTicker.C:
			fetchAgentCardTicker.Stop()
			idx.setAgentCardInserted()
			fetchAgentCardTicker.Reset(fetchAgentProfileInterval)
		case <-logExecBlockTicker.C:
			logExecBlockTicker.Stop()
			idx.logger.WithFields(logrus.Fields{
				"block": idx.execBlock,
				"index": idx.execIndex,
			}).Info("identity registry processor exec block")
			logExecBlockTicker.Reset(logExecBlockInterval)
		}
	}
}

func (idx *IdentityProcessor) process(currentBlockNum int64) {
	idx.mu.RLock()
	fromBlock := int64(idx.execBlock) + 1
	currentExecBlock := idx.execBlock
	currentExecIndex := idx.execIndex
	idx.mu.RUnlock()

loop:
	for {
		toBlock := fromBlock + idx.fetchBlockInterval
		events, err := idx.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{idx.identityAddr},
			Topics:    [][]common.Hash{{RegisteredTopic, UriUpdatedTopic, TransferOwnerShipTopic, SetMetaDataTopic}},
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
			if currentExecBlock > e.BlockNumber {
				continue
			} else if currentExecBlock == e.BlockNumber && currentExecIndex >= uint64(e.Index) {
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

			idx.mu.Lock()
			idx.execBlock = uint64(e.BlockNumber)
			idx.execIndex = uint64(e.Index)
			currentExecBlock = idx.execBlock
			currentExecIndex = idx.execIndex
			idx.mu.Unlock()
		}

		if toBlock < currentBlockNum {
			fromBlock = toBlock
			continue loop
		} else {
			idx.mu.Lock()
			idx.execBlock = uint64(currentBlockNum)
			idx.execIndex = 0
			idx.mu.Unlock()
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
	case SetMetaDataTopic:
		return idx.dealWithSetMetaDataEvent(e)
	default:
		return fmt.Errorf("unknown event topic: %s", e.Topics[0])
	}
}

func (idx *IdentityProcessor) dealWithSetMetaDataEvent(e types.Log) error {
	event, err := idx.identityRegistry.ParseMetadataSet(e)
	if err != nil {
		return fmt.Errorf("failed to parse set meta data event: %w", err)
	}

	if event.MetadataKey == "agentWallet" {
		err := model.UpdateAgentWallet(idx.chainID, idx.identityAddr.String(), event.AgentId.String(), common.BytesToAddress(event.MetadataValue).String())
		if err != nil {
			return fmt.Errorf("failed to update agent wallet: %w", err)
		}
		return nil
	}

	return model.CreateMetadata(&model.Metadata{
		ChainID:          idx.chainID,
		IdentityRegistry: idx.identityAddr.String(),
		AgentID:          event.AgentId.String(),
		Key:              event.MetadataKey,
		Value:            hex.EncodeToString(event.MetadataValue),
	})
}

func (idx *IdentityProcessor) dealWithAgentRegisteredEvent(e types.Log) error {
	agentRegisteredEvent, err := idx.identityRegistry.ParseRegistered(e)
	if err != nil {
		return fmt.Errorf("failed to parse agent registered event: %w", err)
	}

	registry := &model.Agent{
		AgentID:          agentRegisteredEvent.AgentId.String(),
		IdentityRegistry: idx.identityAddr.String(),
		Owner:            agentRegisteredEvent.Owner.String(),
		AgentURI:         agentRegisteredEvent.AgentURI,
		ChainID:          idx.chainID,
		BlockNumber:      uint64(e.BlockNumber),
		Index:            uint64(e.Index),
		TxHash:           e.TxHash.String(),
		Timestamps:       uint64(e.BlockTimestamp),
		Inserted:         false,
	}

	if err := model.CreateAgent(registry); err != nil {
		return fmt.Errorf("failed to create agent registry: %w", err)
	}

	idx.logger.WithFields(logrus.Fields{
		"event":            "register agent",
		"agentID":          agentRegisteredEvent.AgentId.String(),
		"identityRegistry": idx.identityAddr.String(),
		"chainID":          idx.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with agent registered event")
	return nil
}

func (idx *IdentityProcessor) dealWithUriUpdatedEvent(e types.Log) error {
	event, err := idx.identityRegistry.ParseURIUpdated(e)
	if err != nil {
		return fmt.Errorf("failed to parse auth feedback event: %w", err)
	}

	if err := model.UpdateAgentTokenURL(idx.chainID, idx.identityAddr.String(), event.AgentId.String(), event.NewURI, uint64(e.BlockNumber), uint64(e.Index)); err != nil {
		return fmt.Errorf("failed to update agent token url: %w", err)
	}
	idx.logger.WithFields(logrus.Fields{
		"event":            "update agent uri",
		"agentID":          event.AgentId.String(),
		"identityRegistry": idx.identityAddr.String(),
		"chainID":          idx.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with agent uri updated event")

	return nil
}

func (idx *IdentityProcessor) dealWithTransferOwnerShipEvent(e types.Log) error {
	event, err := idx.identityRegistry.ParseTransfer(e)
	if err != nil {
		return fmt.Errorf("failed to parse transfer owner ship event: %w", err)
	}

	if err := model.TransferOwnerShip(idx.chainID, idx.identityAddr.String(), event.TokenId.String(), event.To.String(), uint64(e.BlockNumber), uint64(e.Index)); err != nil {
		return fmt.Errorf("failed to transfer owner ship: %w", err)
	}
	idx.logger.WithFields(logrus.Fields{
		"event":            "transfer owner ship",
		"agentID":          event.TokenId.String(),
		"identityRegistry": idx.identityAddr.String(),
		"chainID":          idx.chainID,
		"blockNumber":      uint64(e.BlockNumber),
		"index":            uint64(e.Index),
		"txHash":           e.TxHash.String(),
		"timestamps":       uint64(e.BlockTimestamp),
	}).Info("deal with transfer owner ship event")
	return nil
}

func (idx *IdentityProcessor) setAgentCardInserted() {
	var limit = 100
	for {
		agents, err := model.GetUnInsertedAgents(idx.chainID, idx.identityAddr.String(), limit)
		if err != nil {
			idx.logger.WithFields(logrus.Fields{
				"error":            err,
				"chainID":          idx.chainID,
				"identityRegistry": idx.identityAddr.String(),
			}).Error("failed to get un-inserted agents")
			return
		}

		if len(agents) == 0 {
			break
		}

		vectorUpserts := make([]model.AgentVectorUpsert, 0, len(agents))
		for _, agent := range agents {
			agentProfile, err := agentcard.GetAgentProfile(agent.AgentURI)
			if err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error":            err,
					"chainID":          idx.chainID,
					"identityRegistry": idx.identityAddr.String(),
					"agentID":          agent.AgentID,
					"agentURI":         agent.AgentURI,
				}).Error("failed to get agent profile from agent url")
			}

			// upload agent to gemini file api
			if agentProfile != nil {
				if _, err := model.UpdateAgent(agent.ChainID, agent.IdentityRegistry, agent.AgentID, agentProfile); err != nil {
					idx.logger.WithFields(logrus.Fields{
						"error":            err,
						"chainID":          agent.ChainID,
						"identityRegistry": agent.IdentityRegistry,
						"agentID":          agent.AgentID,
						"agentURI":         agent.AgentURI,
					}).Error("failed to update agent")
					continue
				}
				vectorUpserts = append(vectorUpserts, model.AgentVectorUpsert{
					AgentUID:         agent.UID,
					IdentityRegistry: agent.IdentityRegistry,
					ChainID:          agent.ChainID,
					CreateTimestamp:  agent.Timestamps,
					Content:          agent.Description,
					Metadata:         nil,
				})
			}

			if err := model.UpdateAgentInserted([]uint64{agent.UID}); err != nil {
				idx.logger.WithFields(logrus.Fields{
					"error":            err,
					"chainID":          idx.chainID,
					"identityRegistry": idx.identityAddr.String(),
					"agentUID":         agent.UID,
					"agentURI":         agent.AgentURI,
				}).Error("failed to update agent registry inserted")
				continue
			}
		}

		// Batch upsert vectors for this page (best-effort, doesn't block inserted status updates above).
		if len(vectorUpserts) > 0 {
			if err := model.InsertAgentVectors(vectorUpserts, 100); err != nil {
				idx.logger.WithFields(logrus.Fields{"error": err}).Error("failed to batch upsert agent vectors")
			}
		}
	}
}
