package model

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

type commerceTestConfig struct {
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
}

func initCommerceTest(t *testing.T) {
	t.Helper()
	conf := &commerceTestConfig{}
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		t.Skipf("skipping: config.yaml not found (%v)", err)
	}
	if err := yaml.Unmarshal(data, conf); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}
	InitDB(conf.Dns, conf.OpenaiAPIKey)
	if err := EnsureCommerceActionsSchema(); err != nil {
		t.Fatalf("failed to ensure schema: %v", err)
	}
}

func TestCreateCommerceAction_Idempotent(t *testing.T) {
	initCommerceTest(t)

	action := &CommerceAction{
		ChainID:          "84532",
		CommerceContract: "0xTEST_CONTRACT",
		JobID:            99999,
		AgentUID:         1,
		AgentAddress:     "0xTEST_ADDRESS",
		Role:             RoleProvider,
		Action:           ActionJobCompleted,
		SignalPolarity:   PolarityPositive,
		SignalWeight:     1.0,
		SignalCertainty:  CertaintyDefinitive,
		BlockNumber:      100,
		TxHash:           "0xTEST_TX",
		LogIndex:         0,
		BlockTimestamp:    1700000000,
	}

	if err := CreateCommerceAction(action); err != nil {
		t.Fatalf("first insert failed: %v", err)
	}

	// second insert should be idempotent (no error, no duplicate)
	if err := CreateCommerceAction(action); err != nil {
		t.Fatalf("idempotent insert failed: %v", err)
	}

	actions, err := GetCommerceActionsByJobID("84532", "0xTEST_CONTRACT", 99999)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	count := 0
	for _, a := range actions {
		if a.TxHash == "0xTEST_TX" && a.Role == RoleProvider {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 record, got %d", count)
	}

	// cleanup
	db.Where("chain_id = ? AND commerce_contract = ? AND job_id = ? AND tx_hash = ?",
		"84532", "0xTEST_CONTRACT", 99999, "0xTEST_TX").
		Delete(&CommerceAction{})
}

func TestGetCommerceActionsByAgentUID_Pagination(t *testing.T) {
	initCommerceTest(t)

	for i := 0; i < 3; i++ {
		a := &CommerceAction{
			ChainID:          "84532",
			CommerceContract: "0xPAGE_TEST",
			JobID:            uint64(80000 + i),
			AgentUID:         99998,
			AgentAddress:     "0xPAGE_ADDR",
			Role:             RoleProvider,
			Action:           ActionJobCompleted,
			SignalPolarity:   PolarityPositive,
			SignalWeight:     1.0,
			SignalCertainty:  CertaintyDefinitive,
			BlockNumber:      uint64(200 + i),
			TxHash:           "0xPAGE_TX_" + string(rune('A'+i)),
			LogIndex:         0,
			BlockTimestamp:    uint64(1700000000 + i),
		}
		if err := CreateCommerceAction(a); err != nil {
			t.Fatalf("insert %d failed: %v", i, err)
		}
	}

	actions, total, err := GetCommerceActionsByAgentUID(99998, "", "", "", 1, 2)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if total < 3 {
		t.Errorf("expected total >= 3, got %d", total)
	}
	if len(actions) > 2 {
		t.Errorf("expected page size <= 2, got %d", len(actions))
	}

	// cleanup
	db.Where("agent_uid = ? AND commerce_contract = ?", 99998, "0xPAGE_TEST").Delete(&CommerceAction{})
}

func TestGetCommerceActionsByQuery_FiltersAndSort(t *testing.T) {
	initCommerceTest(t)

	uid := uint64(99996)
	contract := "0xQUERY_TEST"
	chainID := "84532"

	cleanup := func() {
		db.Where("agent_uid = ? AND commerce_contract = ?", uid, contract).Delete(&CommerceAction{})
	}
	cleanup()
	t.Cleanup(cleanup)

	// Seed 4 actions with varying fields.
	seed := []*CommerceAction{
		{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            91001,
			AgentUID:         uid,
			AgentAddress:     "0xADDR",
			Role:             RoleProvider,
			Action:           ActionJobSubmitted,
			SignalPolarity:   PolarityNeutral,
			SignalWeight:     0.3,
			SignalCertainty:  CertaintyIndicative,
			JobBudget:        10,
			HookAddress:      "",
			BlockNumber:      500,
			TxHash:           "0xTX_1",
			LogIndex:         1,
			BlockTimestamp:   1700000100,
		},
		{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            91002,
			AgentUID:         uid,
			AgentAddress:     "0xADDR",
			Role:             RoleProvider,
			Action:           ActionJobCompleted,
			SignalPolarity:   PolarityPositive,
			SignalWeight:     1.0,
			SignalCertainty:  CertaintyDefinitive,
			JobBudget:        200,
			Counterparty:     "0xCP_A",
			PreviousStatus:   StatusSubmitted,
			HookAddress:      "0xHOOK",
			BlockNumber:      501,
			TxHash:           "0xTX_2",
			LogIndex:         2,
			BlockTimestamp:   1700000200,
		},
		{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            91003,
			AgentUID:         uid,
			AgentAddress:     "0xADDR",
			Role:             RoleClient,
			Action:           ActionJobFunded,
			SignalPolarity:   PolarityNeutral,
			SignalWeight:     0.1,
			SignalCertainty:  CertaintyIndicative,
			JobBudget:        50,
			HookAddress:      "",
			BlockNumber:      502,
			TxHash:           "0xTX_3",
			LogIndex:         3,
			BlockTimestamp:   1700000300,
		},
		{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            91004,
			AgentUID:         uid,
			AgentAddress:     "0xADDR",
			Role:             RoleProvider,
			Action:           ActionJobRejected,
			SignalPolarity:   PolarityNegative,
			SignalWeight:     -0.5,
			SignalCertainty:  CertaintyDefinitive,
			JobBudget:        300,
			Counterparty:     "0xCP_B",
			PreviousStatus:   StatusFunded,
			HookAddress:      "0xHOOK",
			BlockNumber:      503,
			TxHash:           "0xTX_4",
			LogIndex:         4,
			BlockTimestamp:   1700000400,
		},
	}
	for i, a := range seed {
		if err := CreateCommerceAction(a); err != nil {
			t.Fatalf("seed %d insert failed: %v", i, err)
		}
	}

	// Filter: provider + definitive + hasHook
	hasHook := true
	actions, total, err := GetCommerceActionsByQuery(CommerceActionsQuery{
		UID:       uid,
		Role:      RoleProvider,
		Certainty: CertaintyDefinitive,
		HasHook:   &hasHook,
		Page:      1,
		PageSize:  10,
	})
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total=2, got %d", total)
	}
	if len(actions) != 2 {
		t.Fatalf("expected len=2, got %d", len(actions))
	}
	for _, a := range actions {
		if a.Role != RoleProvider || a.SignalCertainty != CertaintyDefinitive || a.HookAddress == "" {
			t.Fatalf("unexpected record: %+v", a)
		}
	}

	// Sort by budget desc: first should be budget=300
	actions, _, err = GetCommerceActionsByQuery(CommerceActionsQuery{
		UID:      uid,
		Role:     RoleProvider,
		SortBy:   "budget",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("sort query failed: %v", err)
	}
	if len(actions) == 0 || actions[0].JobBudget != 300 {
		t.Fatalf("expected first budget=300, got %+v", actions)
	}

	// Time range: keep only timestamps >= 1700000300
	start := uint64(1700000300)
	actions, total, err = GetCommerceActionsByQuery(CommerceActionsQuery{
		UID:       uid,
		StartTime:  &start,
		Page:      1,
		PageSize:  10,
		SortBy:    "timestamp",
	})
	if err != nil {
		t.Fatalf("time range query failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total=2 in range, got %d", total)
	}
	for _, a := range actions {
		if a.BlockTimestamp < start {
			t.Fatalf("unexpected timestamp: %+v", a)
		}
	}
}

func TestQueryPerf_ExplainUsesIndex_SoftCheck(t *testing.T) {
	initCommerceTest(t)

	uid := uint64(99995)
	contract := "0xEXPLAIN_TEST"
	chainID := "84532"

	cleanup := func() {
		db.Where("agent_uid = ? AND commerce_contract = ?", uid, contract).Delete(&CommerceAction{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for i := 0; i < 5; i++ {
		a := &CommerceAction{
			ChainID:          chainID,
			CommerceContract: contract,
			JobID:            uint64(92000 + i),
			AgentUID:         uid,
			AgentAddress:     "0xADDR",
			Role:             RoleProvider,
			Action:           ActionJobCompleted,
			SignalPolarity:   PolarityPositive,
			SignalWeight:     1.0,
			SignalCertainty:  CertaintyDefinitive,
			BlockNumber:      uint64(600 + i),
			TxHash:           "0xEXPLAIN_TX",
			LogIndex:         uint(i),
			BlockTimestamp:   uint64(1700001000 + i),
		}
		if err := CreateCommerceAction(a); err != nil {
			t.Fatalf("seed insert failed: %v", err)
		}
	}

	type row struct{ Plan string }
	var rows []row
	err := db.Raw(`EXPLAIN SELECT * FROM commerce_actions WHERE agent_uid = $1 ORDER BY block_timestamp DESC LIMIT 20 OFFSET 0`, uid).
		Scan(&rows).Error
	if err != nil {
		t.Skipf("EXPLAIN not supported in this env: %v", err)
	}
	// Soft check: we just ensure it doesn't error and returns a non-empty plan.
	if len(rows) == 0 {
		t.Fatalf("expected explain plan rows")
	}
}

func TestGetLatestCommerceAction(t *testing.T) {
	initCommerceTest(t)

	for i := 0; i < 2; i++ {
		a := &CommerceAction{
			ChainID:          "84532",
			CommerceContract: "0xLATEST_TEST",
			JobID:            uint64(70000 + i),
			AgentUID:         99997,
			AgentAddress:     "0xLATEST_ADDR",
			Role:             RoleClient,
			Action:           ActionJobCreated,
			SignalPolarity:   PolarityNeutral,
			SignalWeight:     0.1,
			SignalCertainty:  CertaintyIndicative,
			BlockNumber:      uint64(300 + i),
			TxHash:           "0xLATEST_TX_" + string(rune('A'+i)),
			LogIndex:         uint(i),
			BlockTimestamp:    uint64(1700000000 + i),
		}
		if err := CreateCommerceAction(a); err != nil {
			t.Fatalf("insert %d failed: %v", i, err)
		}
	}

	block, idx, err := GetLatestCommerceAction("84532", "0xLATEST_TEST")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if block != 301 || idx != 1 {
		t.Errorf("expected block=301/idx=1, got block=%d/idx=%d", block, idx)
	}

	// cleanup
	db.Where("agent_uid = ? AND commerce_contract = ?", 99997, "0xLATEST_TEST").Delete(&CommerceAction{})
}

// ─────────────── M5: Stub Agent ───────────────

func TestFindOrCreateAgentByWallet(t *testing.T) {
	initCommerceTest(t)

	wallet := "0xSTUB_TEST_WALLET_001"
	chainID := "84532"

	// cleanup from previous runs
	db.Where("agent_wallet = ? AND chain_id = ?", wallet, chainID).Delete(&Agent{})

	// first call should create
	uid1, err := FindOrCreateAgentByWallet(chainID, wallet)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if uid1 == 0 {
		t.Fatal("expected non-zero uid")
	}

	// second call should be idempotent
	uid2, err := FindOrCreateAgentByWallet(chainID, wallet)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if uid1 != uid2 {
		t.Errorf("expected same uid, got %d vs %d", uid1, uid2)
	}

	// verify stub properties
	var agent Agent
	if err := db.Where("uid = ?", uid1).First(&agent).Error; err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if agent.Active {
		t.Error("stub agent should have active=false")
	}
	if agent.Inserted {
		t.Error("stub agent should have inserted=false")
	}

	// cleanup
	db.Where("uid = ?", uid1).Delete(&Agent{})
}
