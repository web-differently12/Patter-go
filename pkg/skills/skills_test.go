package skills_test

import (
	"context"
	"strings"
	"testing"

	"github.com/lynxflow/patter-go/pkg/skills"
)

type mockDispatcher struct {
	transferred bool
	destination string
	reason      string
}

func (m *mockDispatcher) TransferCall(ctx context.Context, tenantID, callSID, destinationNumber, reason string) error {
	m.transferred = true
	m.destination = destinationNumber
	m.reason = reason
	return nil
}

func TestHumanTransferSkill_ConsentRequired(t *testing.T) {
	disp := &mockDispatcher{}
	skill := skills.NewHumanTransferSkill(disp, nil)

	rawArgs := `{"reason": "EDIFACT v4 non supporte", "user_consented": false}`
	res, err := skill.Execute(context.Background(), "tenant_1", "call_123", rawArgs, "+33699887766")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(res, "CONSENTEMENT_REQUIS") {
		t.Errorf("Expected result to contain CONSENTEMENT_REQUIS, got: %s", res)
	}
	if disp.transferred {
		t.Errorf("Expected transfer NOT to be triggered when user_consented is false")
	}
}

func TestHumanTransferSkill_ConsentedExecute(t *testing.T) {
	disp := &mockDispatcher{}
	skill := skills.NewHumanTransferSkill(disp, nil)

	rawArgs := `{"reason": "EDIFACT v4 non supporte", "user_consented": true, "proposed_department": "TECHNICAL"}`
	res, err := skill.Execute(context.Background(), "tenant_1", "call_123", rawArgs, "+33699887766")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(res, "TRANSFERT_INITIÉ") {
		t.Errorf("Expected result to contain TRANSFERT_INITIÉ, got: %s", res)
	}
	if !disp.transferred {
		t.Errorf("Expected call transfer to be executed")
	}
	if disp.destination != "+33699887766" {
		t.Errorf("Expected destination +33699887766, got %s", disp.destination)
	}
	if disp.reason != "EDIFACT v4 non supporte" {
		t.Errorf("Expected reason 'EDIFACT v4 non supporte', got %s", disp.reason)
	}
}
