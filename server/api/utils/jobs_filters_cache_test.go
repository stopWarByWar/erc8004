package api

import (
	"context"
	"agent_identity/server/api/types"
	"testing"
	"sync/atomic"
)

func TestJobsFiltersCache_OnDemandBuildOnce(t *testing.T) {
	// Reset globals for test.
	jobsFiltersStarted.Store(false)
	jobsFiltersValue = atomic.Value{}

	origBuilder := jobsFiltersBuilder
	calls := 0
	jobsFiltersBuilder = func(ctx context.Context) (*types.CommerceJobsFilters, error) {
		calls++
		return &types.CommerceJobsFilters{
			Chains: []types.CommerceJobsFilterChain{
				{ChainID: "1", ChainName: "Ethereum"},
			},
			CommerceContracts: []string{"0xabc"},
			PaymentTokens:     []string{"0xdef"},
			LastUpdated:       123,
		}, nil
	}
	t.Cleanup(func() {
		jobsFiltersBuilder = origBuilder
	})

	// First call should build.
	s1, err := GetJobsFiltersSnapshot(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s1 == nil || len(s1.Filters.Chains) != 1 || s1.Filters.Chains[0].ChainID != "1" {
		t.Fatalf("unexpected snapshot: %+v", s1)
	}
	if calls != 1 {
		t.Fatalf("expected builder called once, got %d", calls)
	}

	// Second call should be served from memory.
	s2, err := GetJobsFiltersSnapshot(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s2 == nil {
		t.Fatalf("nil snapshot on second call")
	}
	if calls != 1 {
		t.Fatalf("expected builder called once total, got %d", calls)
	}
}

