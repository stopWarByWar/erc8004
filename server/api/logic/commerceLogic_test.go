package logic

import (
	"agent_identity/model"
	"fmt"
	"testing"
)

func TestGetCommerceScores_Global(t *testing.T) {
	initTest(t)
	uid := uint64(1)
	scores, err := GetCommerceScores(uid, "", "")
	if err != nil {
		t.Fatalf("GetCommerceScores error: %v", err)
	}
	if scores == nil {
		t.Fatalf("expected non-nil slice")
	}
	fmt.Printf("global scores for uid=%d: %+v\n", uid, scores)
}

func TestGetCommerceScores_Segmented(t *testing.T) {
	initTest(t)
	uid := uint64(1)
	scores, err := GetCommerceScores(uid, "84532", "0xTEST")
	if err != nil {
		t.Fatalf("GetCommerceScores segmented error: %v", err)
	}
	fmt.Printf("segmented scores for uid=%d: %+v\n", uid, scores)
}

func TestGetCommerceActions(t *testing.T) {
	initTest(t)

	// Seed one action via model to validate query->logic wiring.
	uid := uint64(424242)
	seed := &model.CommerceAction{
		ChainID:          "84532",
		CommerceContract: "0xLOGIC_TEST",
		JobID:            424242,
		AgentUID:         uid,
		AgentAddress:     "0xLOGIC_ADDR",
		Role:             model.RoleProvider,
		Action:           model.ActionJobCompleted,
		SignalPolarity:   model.PolarityPositive,
		SignalWeight:     1.0,
		SignalCertainty:  model.CertaintyDefinitive,
		JobBudget:        123,
		Counterparty:     "0xCP",
		PreviousStatus:   model.StatusSubmitted,
		HookAddress:      "0xHOOK",
		BlockNumber:      777,
		TxHash:           "0xLOGIC_TX",
		LogIndex:         7,
		BlockTimestamp:   1700000001,
	}
	if err := model.CreateCommerceAction(seed); err != nil {
		t.Fatalf("seed insert failed: %v", err)
	}

	actions, total, err := GetCommerceActions(uid, model.RoleProvider, model.ActionJobCompleted, model.CertaintyDefinitive, 1, 10)
	if err != nil {
		t.Fatalf("GetCommerceActions error: %v", err)
	}
	if total < 1 || len(actions) < 1 {
		t.Fatalf("expected >=1 action, got total=%d len=%d", total, len(actions))
	}
	fmt.Printf("actions total=%d, page=%+v\n", total, actions)
}

func TestGetCommerceScoreSummary(t *testing.T) {
	initTest(t)
	uid := uint64(1)
	summary, err := GetCommerceScoreSummary(uid)
	if err != nil {
		t.Fatalf("GetCommerceScoreSummary error: %v", err)
	}
	fmt.Printf("summary for uid=%d: %+v\n", uid, summary)
}
