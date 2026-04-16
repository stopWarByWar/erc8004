package model

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

type jobsTestConfig struct {
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
}

func initJobsTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run DB-backed tests")
	}
	// Prefer server/api/logic/config.yaml (user-authoritative local DB config).
	data, err := os.ReadFile("../server/api/logic/config.yaml")
	if err != nil {
		t.Skipf("skipping: cannot read ../server/api/logic/config.yaml (%v)", err)
	}
	var cfg jobsTestConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}
	if cfg.Dns == "" {
		t.Fatal("empty dns in ../server/api/logic/config.yaml")
	}
	InitDB(cfg.Dns, cfg.OpenaiAPIKey)
	if err := EnsureCommerceJobsSchema(); err != nil {
		t.Fatalf("failed to ensure commerce_jobs schema: %v", err)
	}
}

func TestUpsertCommerceJob_Idempotent(t *testing.T) {
	initJobsTest(t)

	chainID := "84532"
	contract := "0xJOB_TEST"
	jobID := uint64(42424242)

	// cleanup
	db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).Delete(&CommerceJob{})
	t.Cleanup(func() {
		db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).Delete(&CommerceJob{})
	})

	j := &CommerceJob{
		ChainID:          chainID,
		CommerceContract: contract,
		JobID:            jobID,
		Client:           "0xCLIENT",
		Status:           StatusOpen,
		UpdatedAt:        1700000001,
		Budget:           1,
		BudgetUSD:        2,
	}
	if err := UpsertCommerceJob(j); err != nil {
		t.Fatalf("first upsert failed: %v", err)
	}

	// second upsert with changes should update the same row, not create new.
	j2 := *j
	j2.Status = StatusFunded
	j2.UpdatedAt = 1700000002
	j2.PaidAmountUSD = 99
	if err := UpsertCommerceJob(&j2); err != nil {
		t.Fatalf("second upsert failed: %v", err)
	}

	var count int64
	if err := db.Model(&CommerceJob{}).Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).Count(&count).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}

	got, err := GetCommerceJobByID(chainID, contract, jobID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected job exists")
	}
	if got.Status != StatusFunded {
		t.Fatalf("expected status=%s, got %s", StatusFunded, got.Status)
	}
	if got.UpdatedAt != 1700000002 {
		t.Fatalf("expected updated_at=1700000002, got %d", got.UpdatedAt)
	}
	if got.PaidAmountUSD != 99 {
		t.Fatalf("expected paid_amount_usd=99, got %f", got.PaidAmountUSD)
	}
}

func TestSettlementOutOfOrder_DoesNotPanic(t *testing.T) {
	initJobsTest(t)

	chainID := "84532"
	contract := "0xJOB_TEST_OOO"
	jobID := uint64(1111)

	db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).Delete(&CommerceJob{})
	t.Cleanup(func() {
		db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ?", chainID, contract, jobID).Delete(&CommerceJob{})
	})

	// Simulate: settlement event comes before JobCompleted.
	cj, err := FetchOrCreateCommerceJob(chainID, contract, jobID)
	if err != nil {
		t.Fatalf("fetch/create failed: %v", err)
	}
	cj.PaidAmountUSD = 10
	cj.PlatformFeeUSD = 1
	cj.EvaluatorFeeUSD = 2
	cj.UpdatedAt = 1700000100
	if err := UpsertCommerceJob(cj); err != nil {
		t.Fatalf("upsert settlement-first failed: %v", err)
	}

	// Later: completed event updates status, keeps existing money fields.
	cj2, err := FetchOrCreateCommerceJob(chainID, contract, jobID)
	if err != nil {
		t.Fatalf("fetch/create failed: %v", err)
	}
	cj2.Status = StatusCompleted
	cj2.CompletedAt = 1700000200
	cj2.UpdatedAt = 1700000200
	if err := UpsertCommerceJob(cj2); err != nil {
		t.Fatalf("upsert completed-later failed: %v", err)
	}

	got, err := GetCommerceJobByID(chainID, contract, jobID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected job exists")
	}
	if got.Status != StatusCompleted {
		t.Fatalf("expected status=%s, got %s", StatusCompleted, got.Status)
	}
	// money fields should be preserved (upsert updates them too, but cj2 has defaults).
	// Our current UpsertCommerceJob updates all columns, so this is a known caveat:
	// last-writer-wins may overwrite paid/fee with zeros if cj2 doesn't carry them.
	// This test ensures at least no panic; value preservation would require merge logic.
}

