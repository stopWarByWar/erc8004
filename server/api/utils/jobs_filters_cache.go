package api

import (
	"context"
	"agent_identity/config"
	"agent_identity/model"
	"agent_identity/server/api/types"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// JobsFiltersSnapshot is an immutable snapshot stored in memory.
type JobsFiltersSnapshot struct {
	Filters     types.CommerceJobsFilters
	RefreshedAt time.Time
}

var (
	jobsFiltersStarted atomic.Bool
	jobsFiltersValue   atomic.Value // *JobsFiltersSnapshot
	jobsFiltersMu      sync.Mutex   // guards first build

	// Injectable builder for unit tests.
	jobsFiltersBuilder = buildJobsFilters
)

// StartJobsFiltersCache starts background refresh for jobs filters options.
// It is safe to call multiple times; only the first call takes effect.
func StartJobsFiltersCache(ctx context.Context, refreshInterval time.Duration) {
	if refreshInterval <= 0 {
		refreshInterval = 5 * time.Minute
	}
	if !jobsFiltersStarted.CompareAndSwap(false, true) {
		return
	}

	go func() {
		// Best-effort initial refresh.
		_, _ = refreshJobsFilters(ctx)

		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = refreshJobsFilters(ctx)
			}
		}
	}()
}

// GetJobsFiltersSnapshot returns the latest snapshot from memory.
// If cache is empty, it will synchronously build it once as a fallback.
func GetJobsFiltersSnapshot(ctx context.Context) (*JobsFiltersSnapshot, error) {
	if v := jobsFiltersValue.Load(); v != nil {
		return v.(*JobsFiltersSnapshot), nil
	}

	jobsFiltersMu.Lock()
	defer jobsFiltersMu.Unlock()

	// Double check after acquiring lock.
	if v := jobsFiltersValue.Load(); v != nil {
		return v.(*JobsFiltersSnapshot), nil
	}

	return refreshJobsFilters(ctx)
}

func refreshJobsFilters(ctx context.Context) (*JobsFiltersSnapshot, error) {
	filters, err := jobsFiltersBuilder(ctx)
	if err != nil {
		// Keep old snapshot if any.
		if v := jobsFiltersValue.Load(); v != nil {
			return v.(*JobsFiltersSnapshot), err
		}
		return nil, err
	}

	snap := &JobsFiltersSnapshot{
		Filters:     *filters,
		RefreshedAt: time.Now(),
	}
	jobsFiltersValue.Store(snap)
	return snap, nil
}

func buildJobsFilters(ctx context.Context) (*types.CommerceJobsFilters, error) {
	_ = ctx // reserved for future DB timeout/cancellation

	chainIDs, err := model.GetCommerceJobsDistinctChainIDs()
	if err != nil {
		return nil, fmt.Errorf("distinct chain_ids: %w", err)
	}
	contracts, err := model.GetCommerceJobsDistinctCommerceContracts()
	if err != nil {
		return nil, fmt.Errorf("distinct commerce_contracts: %w", err)
	}
	tokens, err := model.GetCommerceJobsDistinctPaymentTokens()
	if err != nil {
		return nil, fmt.Errorf("distinct payment_tokens: %w", err)
	}

	chains := make([]types.CommerceJobsFilterChain, 0, len(chainIDs))
	for _, id := range chainIDs {
		item := types.CommerceJobsFilterChain{ChainID: id}
		if chain, ok := config.GetChainInfo(id); ok {
			item.ChainName = chain.ChainName
			item.ChainLogo = chain.ChainLogo
		}
		chains = append(chains, item)
	}
	sort.Slice(chains, func(i, j int) bool { return chains[i].ChainID < chains[j].ChainID })

	return &types.CommerceJobsFilters{
		Chains:            chains,
		CommerceContracts: contracts,
		PaymentTokens:     tokens,
		LastUpdated:       uint64(time.Now().Unix()),
	}, nil
}

