package model

import "testing"

func TestDetermineSignal_TerminalCompleted(t *testing.T) {
	tests := []struct {
		role            string
		wantPolarity    string
		wantWeight      float64
		wantCertainty   string
		wantShouldWrite bool
	}{
		{RoleProvider, PolarityPositive, 1.0, CertaintyDefinitive, true},
		{RoleClient, PolarityPositive, 1.0, CertaintyDefinitive, true},
		{RoleEvaluator, PolarityNeutral, 0.5, CertaintyDefinitive, true},
	}
	for _, tt := range tests {
		sig := DetermineSignal(ActionJobCompleted, StatusSubmitted, tt.role)
		if sig.ShouldWrite != tt.wantShouldWrite {
			t.Fatalf("Completed/%s: ShouldWrite = %v, want %v", tt.role, sig.ShouldWrite, tt.wantShouldWrite)
		}
		if sig.Polarity != tt.wantPolarity {
			t.Errorf("Completed/%s: Polarity = %s, want %s", tt.role, sig.Polarity, tt.wantPolarity)
		}
		if sig.Weight != tt.wantWeight {
			t.Errorf("Completed/%s: Weight = %f, want %f", tt.role, sig.Weight, tt.wantWeight)
		}
		if sig.Certainty != tt.wantCertainty {
			t.Errorf("Completed/%s: Certainty = %s, want %s", tt.role, sig.Certainty, tt.wantCertainty)
		}
	}
}

func TestDetermineSignal_RejectedFromOpen(t *testing.T) {
	// Provider: no record; Client: neutral; Evaluator: no record
	sig := DetermineSignal(ActionJobRejected, StatusOpen, RoleProvider)
	if sig.ShouldWrite {
		t.Error("Rejected(Open)/Provider should not write")
	}
	sig = DetermineSignal(ActionJobRejected, StatusOpen, RoleClient)
	if !sig.ShouldWrite {
		t.Fatal("Rejected(Open)/Client should write")
	}
	if sig.Polarity != PolarityNeutral || sig.Weight != 0.2 {
		t.Errorf("Rejected(Open)/Client: got %s/%.1f, want neutral/0.2", sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobRejected, StatusOpen, RoleEvaluator)
	if sig.ShouldWrite {
		t.Error("Rejected(Open)/Evaluator should not write")
	}
}

func TestDetermineSignal_RejectedFromFunded(t *testing.T) {
	sig := DetermineSignal(ActionJobRejected, StatusFunded, RoleProvider)
	if !sig.ShouldWrite || sig.Polarity != PolarityNegative || sig.Weight != 0.5 {
		t.Errorf("Rejected(Funded)/Provider: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobRejected, StatusFunded, RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral {
		t.Errorf("Rejected(Funded)/Client: got write=%v %s", sig.ShouldWrite, sig.Polarity)
	}
	sig = DetermineSignal(ActionJobRejected, StatusFunded, RoleEvaluator)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Weight != 0.5 {
		t.Errorf("Rejected(Funded)/Evaluator: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
}

func TestDetermineSignal_RejectedFromSubmitted(t *testing.T) {
	sig := DetermineSignal(ActionJobRejected, StatusSubmitted, RoleProvider)
	if !sig.ShouldWrite || sig.Polarity != PolarityNegative || sig.Weight != 1.0 {
		t.Errorf("Rejected(Submitted)/Provider: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobRejected, StatusSubmitted, RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Weight != 0.3 {
		t.Errorf("Rejected(Submitted)/Client: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobRejected, StatusSubmitted, RoleEvaluator)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Weight != 0.5 {
		t.Errorf("Rejected(Submitted)/Evaluator: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
}

func TestDetermineSignal_ExpiredFromFunded(t *testing.T) {
	sig := DetermineSignal(ActionJobExpired, StatusFunded, RoleProvider)
	if !sig.ShouldWrite || sig.Polarity != PolarityNegative || sig.Weight != 1.0 {
		t.Errorf("Expired(Funded)/Provider: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobExpired, StatusFunded, RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Weight != 0.3 {
		t.Errorf("Expired(Funded)/Client: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobExpired, StatusFunded, RoleEvaluator)
	if sig.ShouldWrite {
		t.Error("Expired(Funded)/Evaluator should not write")
	}
}

func TestDetermineSignal_ExpiredFromSubmitted(t *testing.T) {
	sig := DetermineSignal(ActionJobExpired, StatusSubmitted, RoleProvider)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Weight != 0.3 {
		t.Errorf("Expired(Submitted)/Provider: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobExpired, StatusSubmitted, RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityNegative || sig.Weight != 0.5 {
		t.Errorf("Expired(Submitted)/Client: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
	sig = DetermineSignal(ActionJobExpired, StatusSubmitted, RoleEvaluator)
	if !sig.ShouldWrite || sig.Polarity != PolarityNegative || sig.Weight != 0.7 {
		t.Errorf("Expired(Submitted)/Evaluator: got write=%v %s/%.1f", sig.ShouldWrite, sig.Polarity, sig.Weight)
	}
}

func TestDetermineSignal_ProcessEvents(t *testing.T) {
	// job_created -> Client neutral indicative
	sig := DetermineSignal(ActionJobCreated, "", RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Weight != 0.1 || sig.Certainty != CertaintyIndicative {
		t.Errorf("Created/Client: got write=%v %s/%.1f/%s", sig.ShouldWrite, sig.Polarity, sig.Weight, sig.Certainty)
	}
	sig = DetermineSignal(ActionJobCreated, "", RoleProvider)
	if sig.ShouldWrite {
		t.Error("Created/Provider should not write")
	}

	// job_funded -> Client positive indicative
	sig = DetermineSignal(ActionJobFunded, "", RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityPositive || sig.Weight != 0.3 || sig.Certainty != CertaintyIndicative {
		t.Errorf("Funded/Client: got write=%v %s/%.1f/%s", sig.ShouldWrite, sig.Polarity, sig.Weight, sig.Certainty)
	}

	// job_submitted -> Provider positive indicative
	sig = DetermineSignal(ActionJobSubmitted, "", RoleProvider)
	if !sig.ShouldWrite || sig.Polarity != PolarityPositive || sig.Weight != 0.3 || sig.Certainty != CertaintyIndicative {
		t.Errorf("Submitted/Provider: got write=%v %s/%.1f/%s", sig.ShouldWrite, sig.Polarity, sig.Weight, sig.Certainty)
	}

	// provider_set -> Provider neutral indicative (audit)
	sig = DetermineSignal(ActionProviderSet, "", RoleProvider)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Certainty != CertaintyIndicative {
		t.Errorf("ProviderSet/Provider: got write=%v %s/%s", sig.ShouldWrite, sig.Polarity, sig.Certainty)
	}

	// budget_set -> Client neutral indicative (audit)
	sig = DetermineSignal(ActionBudgetSet, "", RoleClient)
	if !sig.ShouldWrite || sig.Polarity != PolarityNeutral || sig.Certainty != CertaintyIndicative {
		t.Errorf("BudgetSet/Client: got write=%v %s/%s", sig.ShouldWrite, sig.Polarity, sig.Certainty)
	}
}

func TestDetermineSignal_AllTerminalAreDefinitive(t *testing.T) {
	terminalActions := []string{ActionJobCompleted, ActionJobRejected, ActionJobExpired}
	statuses := []string{StatusOpen, StatusFunded, StatusSubmitted}
	roles := []string{RoleProvider, RoleClient, RoleEvaluator}

	for _, action := range terminalActions {
		for _, status := range statuses {
			for _, role := range roles {
				sig := DetermineSignal(action, status, role)
				if sig.ShouldWrite && sig.Certainty != CertaintyDefinitive {
					t.Errorf("%s/%s/%s: certainty = %s, want definitive", action, status, role, sig.Certainty)
				}
			}
		}
	}
}
