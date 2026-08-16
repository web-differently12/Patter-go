package service

import (
	"context"
	"testing"
)

func TestSendRCS(t *testing.T) {
	svc := NewMessagingService()
	ctx := context.Background()

	resp, err := svc.SendRCS(ctx, "tenant_test", SendRCSRequest{
		FromNumber: "+15551234567",
		ToNumber:   "+33612345678",
		Title:      "Offre Exclusive",
		Message:    "Profitez de -20% aujourd'hui!",
		Buttons: []RCSButton{
			{Title: "Voir l'offre", Type: "URL", Payload: "https://patter.ai/promo"},
		},
		FallbackSMS: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ChannelUsed != "RCS" {
		t.Errorf("expected RCS channel, got %s", resp.ChannelUsed)
	}
	if !resp.RCSCapable {
		t.Errorf("expected RCS capable true")
	}
}
