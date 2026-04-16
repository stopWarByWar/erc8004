package logic

import (
	"agent_identity/model"
	"testing"
)

func TestJobToDTO_AmountStringsAndDecimals(t *testing.T) {
	j := model.CommerceJob{
		ChainID:          "1",
		CommerceContract: "0xC",
		JobID:            1,
		Status:           model.StatusCompleted,
		UpdatedAt:        1,
		PaymentToken:     "0xT",
		PaymentDecimals:  6,
		TokenSymbol:      "USDC",
		Budget:           12.34000000,
		PaidAmount:       1.23000000,
		PlatformFeeAmount: 0.10000000,
		EvaluatorFeeAmount: 0,
	}
	dto := jobToDTO(j)
	if dto.PaymentDecimals != 6 {
		t.Fatalf("expected payment_decimals=6, got %d", dto.PaymentDecimals)
	}
	if dto.Budget != "12.34" {
		t.Fatalf("expected budget=12.34, got %q", dto.Budget)
	}
	if dto.PaidAmount != "1.23" {
		t.Fatalf("expected paid_amount=1.23, got %q", dto.PaidAmount)
	}
	if dto.PlatformFeeAmount != "0.1" {
		t.Fatalf("expected platform_fee_amount=0.1, got %q", dto.PlatformFeeAmount)
	}
	if dto.EvaluatorFeeAmount != "0" {
		t.Fatalf("expected evaluator_fee_amount=0, got %q", dto.EvaluatorFeeAmount)
	}
}

func TestActionToDTO_DefaultDecimals(t *testing.T) {
	a := model.CommerceAction{
		ChainID:          "1",
		CommerceContract: "0xC",
		JobID:            1,
		AgentUID:         1,
		AgentAddress:     "0xA",
		Role:             model.RoleClient,
		Action:           model.ActionBudgetSet,
		JobBudget:        100.00000000,
		PaymentToken:     "0xT",
		PaymentDecimals:  0, // missing -> default 18
		TokenSymbol:      "ETH",
		BlockNumber:      1,
		TxHash:           "0xTX",
		LogIndex:         1,
		BlockTimestamp:   1,
	}
	dto := ActionToDTO(a)
	if dto.PaymentDecimals != 18 {
		t.Fatalf("expected default payment_decimals=18, got %d", dto.PaymentDecimals)
	}
	if dto.JobBudget != "100" {
		t.Fatalf("expected job_budget=100, got %q", dto.JobBudget)
	}
}

