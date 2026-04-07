package processor

import (
	"agent_identity/abi"
	"agent_identity/model"
	"context"
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
)

var (
	JobCreatedTopic   = common.HexToHash("0xb0f0239bfdd96453e24733e18bfc24b70d8fadf123dd977473518dd577ee79b9")
	JobFundedTopic    = common.HexToHash("0x") // filled below
	JobSubmittedTopic = common.HexToHash("0x")
	JobCompletedTopic = common.HexToHash("0x")
	JobRejectedTopic  = common.HexToHash("0x")
	JobExpiredTopic   = common.HexToHash("0x")
	ProviderSetTopic  = common.HexToHash("0x")
	BudgetSetTopic    = common.HexToHash("0x")
)

func init() {
	acABI, err := abi.AgenticCommerceMetaData.GetAbi()
	if err != nil {
		panic(fmt.Sprintf("failed to parse AgenticCommerce ABI: %v", err))
	}
	JobCreatedTopic = acABI.Events["JobCreated"].ID
	JobFundedTopic = acABI.Events["JobFunded"].ID
	JobSubmittedTopic = acABI.Events["JobSubmitted"].ID
	JobCompletedTopic = acABI.Events["JobCompleted"].ID
	JobRejectedTopic = acABI.Events["JobRejected"].ID
	JobExpiredTopic = acABI.Events["JobExpired"].ID
	ProviderSetTopic = acABI.Events["ProviderSet"].ID
	BudgetSetTopic = acABI.Events["BudgetSet"].ID
}

type CommerceProcessor struct {
	execBlock          uint64
	execIndex          uint64
	commerceContract   *abi.AgenticCommerce
	commerceAddr       common.Address
	fetchBlockInterval int64
	chainID            string
	logger             *logger.Logger
	ethClient          *ethclient.Client
}

func NewCommerceProcessor(commerceAddr string, ethClient *ethclient.Client, fetchBlockInterval int64, startBlock uint64, _logger *logger.Logger) *CommerceProcessor {
	chainId, err := ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	contract, err := abi.NewAgenticCommerce(common.HexToAddress(commerceAddr), ethClient)
	if err != nil {
		panic(err)
	}

	execBlock, execIndex, err := model.GetLatestCommerceAction(chainId.String(), common.HexToAddress(commerceAddr).String())
	if err != nil {
		panic(err)
	}

	if startBlock > 0 {
		execBlock = startBlock
		execIndex = 0
	}

	return &CommerceProcessor{
		execBlock:          execBlock,
		execIndex:          execIndex,
		commerceContract:   contract,
		commerceAddr:       common.HexToAddress(commerceAddr),
		fetchBlockInterval: fetchBlockInterval,
		ethClient:          ethClient,
		logger:             _logger,
		chainID:            chainId.String(),
	}
}

func (p *CommerceProcessor) Process() {
	p.logger.WithFields(logrus.Fields{
		"block": p.execBlock,
		"index": p.execIndex,
	}).Info("start commerce processor")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	logTicker := time.NewTicker(60 * time.Second)
	defer logTicker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			currentBlock, err := p.ethClient.BlockNumber(context.Background())
			if err != nil {
				p.logger.WithField("error", err).Error("fail to get current block num")
				ticker.Reset(20 * time.Second)
				continue
			}
			if p.execBlock < currentBlock {
				p.process(int64(currentBlock))
			}
			ticker.Reset(20 * time.Second)
		case <-logTicker.C:
			logTicker.Stop()
			p.logger.WithFields(logrus.Fields{
				"block": p.execBlock,
				"index": p.execIndex,
			}).Info("commerce processor exec block")
			logTicker.Reset(60 * time.Second)
		}
	}
}

func (p *CommerceProcessor) process(currentBlockNum int64) {
	fromBlock := int64(p.execBlock) + 1
	topics := []common.Hash{
		JobCreatedTopic, JobFundedTopic, JobSubmittedTopic,
		JobCompletedTopic, JobRejectedTopic, JobExpiredTopic,
		ProviderSetTopic, BudgetSetTopic,
	}

loop:
	for {
		toBlock := fromBlock + p.fetchBlockInterval
		events, err := p.ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{p.commerceAddr},
			Topics:    [][]common.Hash{topics},
		})
		if err != nil {
			p.logger.WithFields(logrus.Fields{"from": fromBlock, "to": toBlock, "error": err}).Error("fail to get commerce events")
			time.Sleep(10 * time.Second)
			continue loop
		}

		for _, e := range events {
			if p.execBlock > e.BlockNumber {
				continue
			} else if p.execBlock == e.BlockNumber && p.execIndex >= uint64(e.Index) {
				continue
			}
			if err := p.dealWithEvent(e); err != nil {
				fields := logrus.Fields{"error": err, "block": e.BlockNumber, "index": e.Index, "txHash": e.TxHash.String()}
				if isRetryable(err) {
					p.logger.WithFields(fields).Warn("retryable commerce event error; will retry")
					time.Sleep(2 * time.Second)
					continue
				}
				// Permanent error: skip this log to avoid getting stuck.
				p.logger.WithFields(fields).Error("permanent commerce event error; skipping log")
				p.execBlock = e.BlockNumber
				p.execIndex = uint64(e.Index)
				continue
			}
			p.execBlock = e.BlockNumber
			p.execIndex = uint64(e.Index)
		}

		if toBlock < currentBlockNum {
			fromBlock = toBlock
			continue loop
		}
		p.execBlock = uint64(currentBlockNum)
		p.execIndex = 0
		return
	}
}

func (p *CommerceProcessor) dealWithEvent(e types.Log) error {
	switch e.Topics[0] {
	case JobCreatedTopic:
		return p.handleJobCreated(e)
	case JobFundedTopic:
		return p.handleJobFunded(e)
	case JobSubmittedTopic:
		return p.handleJobSubmitted(e)
	case JobCompletedTopic:
		return p.handleJobCompleted(e)
	case JobRejectedTopic:
		return p.handleJobRejected(e)
	case JobExpiredTopic:
		return p.handleJobExpired(e)
	case ProviderSetTopic:
		return p.handleProviderSet(e)
	case BudgetSetTopic:
		return p.handleBudgetSet(e)
	default:
		return fmt.Errorf("unknown commerce event topic: %s", e.Topics[0])
	}
}

func (p *CommerceProcessor) handleJobCreated(e types.Log) error {
	ev, err := p.commerceContract.ParseJobCreated(e)
	if err != nil {
		return fmt.Errorf("parse JobCreated: %w", err)
	}
	clientUID, err := model.FindOrCreateAgentByWallet(p.chainID, ev.Client.String())
	if err != nil {
		return fmt.Errorf("find/create client agent: %w", err)
	}

	sig := model.DetermineSignal(model.ActionJobCreated, "", model.RoleClient)
	if sig.ShouldWrite {
		if err := model.CreateCommerceAction(&model.CommerceAction{
			ChainID: p.chainID, CommerceContract: p.commerceAddr.String(),
			JobID: ev.JobId.Uint64(), AgentUID: clientUID, AgentAddress: ev.Client.String(),
			Role: model.RoleClient, Action: model.ActionJobCreated,
			SignalPolarity: sig.Polarity, SignalWeight: sig.Weight, SignalCertainty: sig.Certainty,
			HookAddress: ev.Hook.String(),
			BlockNumber: e.BlockNumber, TxHash: e.TxHash.String(), LogIndex: e.Index, BlockTimestamp: e.BlockTimestamp,
		}); err != nil {
			return fmt.Errorf("create action: %w", err)
		}
	}

	p.logger.WithFields(logrus.Fields{"event": "JobCreated", "jobId": ev.JobId, "block": e.BlockNumber}).Info("commerce event")
	return nil
}

func (p *CommerceProcessor) handleJobFunded(e types.Log) error {
	ev, err := p.commerceContract.ParseJobFunded(e)
	if err != nil {
		return fmt.Errorf("parse JobFunded: %w", err)
	}
	clientUID, err := model.FindOrCreateAgentByWallet(p.chainID, ev.Client.String())
	if err != nil {
		return fmt.Errorf("find/create client agent: %w", err)
	}

	budget := bigIntToFloat(ev.Amount)
	hookAddr, hookErr := p.tryGetHookAddress(ev.JobId.Uint64())
	if hookErr != nil && isRetryable(hookErr) {
		return hookErr
	}
	sig := model.DetermineSignal(model.ActionJobFunded, "", model.RoleClient)
	if sig.ShouldWrite {
		if err := model.CreateCommerceAction(&model.CommerceAction{
			ChainID: p.chainID, CommerceContract: p.commerceAddr.String(),
			JobID: ev.JobId.Uint64(), AgentUID: clientUID, AgentAddress: ev.Client.String(),
			Role: model.RoleClient, Action: model.ActionJobFunded,
			SignalPolarity: sig.Polarity, SignalWeight: sig.Weight, SignalCertainty: sig.Certainty,
			JobBudget:   budget,
			HookAddress: hookAddr,
			BlockNumber: e.BlockNumber, TxHash: e.TxHash.String(), LogIndex: e.Index, BlockTimestamp: e.BlockTimestamp,
		}); err != nil {
			return fmt.Errorf("create action: %w", err)
		}
	}

	p.logger.WithFields(logrus.Fields{"event": "JobFunded", "jobId": ev.JobId, "block": e.BlockNumber}).Info("commerce event")
	return nil
}

func (p *CommerceProcessor) handleJobSubmitted(e types.Log) error {
	ev, err := p.commerceContract.ParseJobSubmitted(e)
	if err != nil {
		return fmt.Errorf("parse JobSubmitted: %w", err)
	}
	providerUID, err := model.FindOrCreateAgentByWallet(p.chainID, ev.Provider.String())
	if err != nil {
		return fmt.Errorf("find/create provider agent: %w", err)
	}

	hookAddr, hookErr := p.tryGetHookAddress(ev.JobId.Uint64())
	if hookErr != nil && isRetryable(hookErr) {
		return hookErr
	}
	sig := model.DetermineSignal(model.ActionJobSubmitted, "", model.RoleProvider)
	if sig.ShouldWrite {
		if err := model.CreateCommerceAction(&model.CommerceAction{
			ChainID: p.chainID, CommerceContract: p.commerceAddr.String(),
			JobID: ev.JobId.Uint64(), AgentUID: providerUID, AgentAddress: ev.Provider.String(),
			Role: model.RoleProvider, Action: model.ActionJobSubmitted,
			SignalPolarity: sig.Polarity, SignalWeight: sig.Weight, SignalCertainty: sig.Certainty,
			Deliverable: common.BytesToHash(ev.Deliverable[:]).String(),
			HookAddress: hookAddr,
			BlockNumber: e.BlockNumber, TxHash: e.TxHash.String(), LogIndex: e.Index, BlockTimestamp: e.BlockTimestamp,
		}); err != nil {
			return fmt.Errorf("create action: %w", err)
		}
	}

	p.logger.WithFields(logrus.Fields{"event": "JobSubmitted", "jobId": ev.JobId, "block": e.BlockNumber}).Info("commerce event")
	return nil
}

// handleJobCompleted writes definitive signals for provider, client, evaluator
func (p *CommerceProcessor) handleJobCompleted(e types.Log) error {
	ev, err := p.commerceContract.ParseJobCompleted(e)
	if err != nil {
		return fmt.Errorf("parse JobCompleted: %w", err)
	}

	job, err := p.getJobWithRetry(ev.JobId.Uint64(), 3)
	if err != nil {
		return err
	}

	budget := bigIntToFloat(job.Budget)
	reason := common.BytesToHash(ev.Reason[:]).String()
	prevStatus := model.StatusSubmitted

	return p.writeTerminalSignals(e, ev.JobId.Uint64(), job, budget, reason, model.ActionJobCompleted, prevStatus)
}

func (p *CommerceProcessor) handleJobRejected(e types.Log) error {
	ev, err := p.commerceContract.ParseJobRejected(e)
	if err != nil {
		return fmt.Errorf("parse JobRejected: %w", err)
	}

	job, err := p.getJobWithRetry(ev.JobId.Uint64(), 3)
	if err != nil {
		return err
	}

	budget := bigIntToFloat(job.Budget)
	reason := common.BytesToHash(ev.Reason[:]).String()

	// The job is now in Rejected state (status=4). We need to infer previous status
	// from on-chain data: look at whether a submit event exists for this job.
	prevStatus := p.inferPreviousStatusForRejected(ev.JobId.Uint64(), budget)

	return p.writeTerminalSignals(e, ev.JobId.Uint64(), job, budget, reason, model.ActionJobRejected, prevStatus)
}

func (p *CommerceProcessor) handleJobExpired(e types.Log) error {
	ev, err := p.commerceContract.ParseJobExpired(e)
	if err != nil {
		return fmt.Errorf("parse JobExpired: %w", err)
	}

	job, err := p.getJobWithRetry(ev.JobId.Uint64(), 3)
	if err != nil {
		return err
	}

	budget := bigIntToFloat(job.Budget)

	// Infer previous status: check if there's a submit action for this job
	prevStatus := p.inferPreviousStatusForExpired(ev.JobId.Uint64())

	return p.writeTerminalSignals(e, ev.JobId.Uint64(), job, budget, "", model.ActionJobExpired, prevStatus)
}

// writeTerminalSignals fans out definitive signals for all three roles
func (p *CommerceProcessor) writeTerminalSignals(e types.Log, jobID uint64, job struct {
	Id          *big.Int
	Client      common.Address
	Provider    common.Address
	Evaluator   common.Address
	Description string
	Budget      *big.Int
	ExpiredAt   *big.Int
	Status      uint8
	Hook        common.Address
}, budget float64, reason, action, prevStatus string) error {

	roles := []struct {
		role    string
		address common.Address
	}{
		{model.RoleProvider, job.Provider},
		{model.RoleClient, job.Client},
		{model.RoleEvaluator, job.Evaluator},
	}

	for _, r := range roles {
		sig := model.DetermineSignal(action, prevStatus, r.role)
		if !sig.ShouldWrite {
			continue
		}

		uid, err := model.FindOrCreateAgentByWallet(p.chainID, r.address.String())
		if err != nil {
			return fmt.Errorf("find/create %s agent: %w", r.role, err)
		}

		counterparty := ""
		switch r.role {
		case model.RoleProvider:
			counterparty = job.Client.String()
		case model.RoleClient:
			counterparty = job.Provider.String()
		case model.RoleEvaluator:
			counterparty = job.Provider.String()
		}

		if err := model.CreateCommerceAction(&model.CommerceAction{
			ChainID: p.chainID, CommerceContract: p.commerceAddr.String(),
			JobID: jobID, AgentUID: uid, AgentAddress: r.address.String(),
			Role: r.role, Action: action,
			SignalPolarity: sig.Polarity, SignalWeight: sig.Weight, SignalCertainty: sig.Certainty,
			JobBudget: budget, Counterparty: counterparty, Reason: reason,
			PreviousStatus: prevStatus, HookAddress: job.Hook.String(),
			BlockNumber: e.BlockNumber, TxHash: e.TxHash.String(), LogIndex: e.Index, BlockTimestamp: e.BlockTimestamp,
		}); err != nil {
			return fmt.Errorf("create %s action: %w", r.role, err)
		}
	}

	p.logger.WithFields(logrus.Fields{"event": action, "jobId": jobID, "block": e.BlockNumber}).Info("commerce terminal event")
	return nil
}

// inferPreviousStatusForRejected uses existing action records to determine
// if the job was in Open, Funded, or Submitted state before rejection.
func (p *CommerceProcessor) inferPreviousStatusForRejected(jobID uint64, budget float64) string {
	actions, _ := model.GetCommerceActionsByJobID(p.chainID, p.commerceAddr.String(), jobID)
	hasSubmit := false
	hasFund := false
	for _, a := range actions {
		if a.Action == model.ActionJobSubmitted {
			hasSubmit = true
		}
		if a.Action == model.ActionJobFunded {
			hasFund = true
		}
	}
	if hasSubmit {
		return model.StatusSubmitted
	}
	if hasFund || budget > 0 {
		return model.StatusFunded
	}
	return model.StatusOpen
}

func (p *CommerceProcessor) inferPreviousStatusForExpired(jobID uint64) string {
	actions, _ := model.GetCommerceActionsByJobID(p.chainID, p.commerceAddr.String(), jobID)
	for _, a := range actions {
		if a.Action == model.ActionJobSubmitted {
			return model.StatusSubmitted
		}
	}
	return model.StatusFunded
}

func bigIntToFloat(v *big.Int) float64 {
	if v == nil {
		return 0
	}
	f, _ := new(big.Float).SetInt(v).Float64()
	return f
}

type retryableError struct{ err error }

func (e retryableError) Error() string { return e.err.Error() }
func (e retryableError) Unwrap() error { return e.err }

func retryable(err error) error {
	if err == nil {
		return nil
	}
	return retryableError{err: err}
}

func isRetryable(err error) bool {
	var re retryableError
	return errors.As(err, &re)
}

func (p *CommerceProcessor) tryGetHookAddress(jobID uint64) (string, error) {
	job, err := p.getJobWithRetry(jobID, 3)
	if err != nil {
		return "", err
	}
	return job.Hook.String(), nil
}

func (p *CommerceProcessor) getJobWithRetry(jobID uint64, attempts int) (struct {
	Id          *big.Int
	Client      common.Address
	Provider    common.Address
	Evaluator   common.Address
	Description string
	Budget      *big.Int
	ExpiredAt   *big.Int
	Status      uint8
	Hook        common.Address
}, error) {
	var last error
	for i := 0; i < attempts; i++ {
		job, err := p.commerceContract.GetJob(nil, new(big.Int).SetUint64(jobID))
		if err == nil {
			return job, nil
		}
		last = err
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}
	return struct {
		Id          *big.Int
		Client      common.Address
		Provider    common.Address
		Evaluator   common.Address
		Description string
		Budget      *big.Int
		ExpiredAt   *big.Int
		Status      uint8
		Hook        common.Address
	}{}, retryable(fmt.Errorf("getJob(%d): %w", jobID, last))
}

func (p *CommerceProcessor) handleProviderSet(e types.Log) error {
	ev, err := p.commerceContract.ParseProviderSet(e)
	if err != nil {
		return fmt.Errorf("parse ProviderSet: %w", err)
	}
	providerUID, err := model.FindOrCreateAgentByWallet(p.chainID, ev.Provider.String())
	if err != nil {
		return fmt.Errorf("find/create provider agent: %w", err)
	}

	hookAddr, hookErr := p.tryGetHookAddress(ev.JobId.Uint64())
	if hookErr != nil && isRetryable(hookErr) {
		return hookErr
	}

	sig := model.DetermineSignal(model.ActionProviderSet, "", model.RoleProvider)
	if sig.ShouldWrite {
		if err := model.CreateCommerceAction(&model.CommerceAction{
			ChainID: p.chainID, CommerceContract: p.commerceAddr.String(),
			JobID: ev.JobId.Uint64(), AgentUID: providerUID, AgentAddress: ev.Provider.String(),
			Role: model.RoleProvider, Action: model.ActionProviderSet,
			SignalPolarity: sig.Polarity, SignalWeight: sig.Weight, SignalCertainty: sig.Certainty,
			HookAddress: hookAddr,
			BlockNumber: e.BlockNumber, TxHash: e.TxHash.String(), LogIndex: e.Index, BlockTimestamp: e.BlockTimestamp,
		}); err != nil {
			return fmt.Errorf("create action: %w", err)
		}
	}
	p.logger.WithFields(logrus.Fields{"event": "ProviderSet", "jobId": ev.JobId, "block": e.BlockNumber}).Info("commerce event")
	return nil
}

func (p *CommerceProcessor) handleBudgetSet(e types.Log) error {
	ev, err := p.commerceContract.ParseBudgetSet(e)
	if err != nil {
		return fmt.Errorf("parse BudgetSet: %w", err)
	}
	// Budget is a client-side context change; use client role for timeline.
	job, jobErr := p.getJobWithRetry(ev.JobId.Uint64(), 3)
	if jobErr != nil {
		return jobErr
	}
	clientUID, err := model.FindOrCreateAgentByWallet(p.chainID, job.Client.String())
	if err != nil {
		return fmt.Errorf("find/create client agent: %w", err)
	}
	budget := bigIntToFloat(ev.Amount)
	sig := model.DetermineSignal(model.ActionBudgetSet, "", model.RoleClient)
	if sig.ShouldWrite {
		if err := model.CreateCommerceAction(&model.CommerceAction{
			ChainID: p.chainID, CommerceContract: p.commerceAddr.String(),
			JobID: ev.JobId.Uint64(), AgentUID: clientUID, AgentAddress: job.Client.String(),
			Role: model.RoleClient, Action: model.ActionBudgetSet,
			SignalPolarity: sig.Polarity, SignalWeight: sig.Weight, SignalCertainty: sig.Certainty,
			JobBudget: budget, HookAddress: job.Hook.String(),
			BlockNumber: e.BlockNumber, TxHash: e.TxHash.String(), LogIndex: e.Index, BlockTimestamp: e.BlockTimestamp,
		}); err != nil {
			return fmt.Errorf("create action: %w", err)
		}
	}
	p.logger.WithFields(logrus.Fields{"event": "BudgetSet", "jobId": ev.JobId, "block": e.BlockNumber}).Info("commerce event")
	return nil
}
