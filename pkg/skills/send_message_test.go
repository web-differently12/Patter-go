package skills

import (
	"context"
	"testing"
)

func TestMessageSkillEngine(t *testing.T) {
	engine := NewMessageSkillEngine()
	ctx := context.Background()

	valRes, err := engine.SendValidationCode(ctx, SendValidationCodeRequest{
		TenantID:   "tenant_test",
		Recipient:  "+33612345678",
		Channel:    "WHATSAPP",
		CodeLength: 6,
	})
	if err != nil {
		t.Fatalf("unexpected validation code skill error: %v", err)
	}
	if !valRes.Success || len(valRes.CodeSent) != 6 {
		t.Errorf("expected 6-digit code, got %s", valRes.CodeSent)
	}

	outRes, err := engine.SendOutboundMessage(ctx, SendOutboundMessageRequest{
		TenantID:  "tenant_test",
		Recipient: "+33612345678",
		Channel:   "RCS",
		Content:   "Voici votre confirmation de rendez-vous",
	})
	if err != nil {
		t.Fatalf("unexpected outbound message skill error: %v", err)
	}
	if !outRes.Success || outRes.ChannelUsed != "RCS" {
		t.Errorf("expected successful RCS outbound dispatch")
	}
}
