package sip_test

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/sip"
)

func TestSIPTrunkRegistrationAndCall(t *testing.T) {
	svc := sip.NewSIPPBXService(nil)

	trunkReq := sip.SIPTrunkConfig{
		Provider:   sip.TrunkOVHTelecom,
		ServerHost: "sip.ovh.fr",
		Username:   "0033988776655",
		Password:   "secret_password",
	}

	tr, err := svc.RegisterTrunk(context.Background(), "tenant_ovh", trunkReq)
	if err != nil {
		t.Fatalf("Unexpected error registering trunk: %v", err)
	}

	if !tr.Registered {
		t.Errorf("Expected trunk to be registered")
	}

	callReq := sip.SIPCallRequest{
		TrunkID:   tr.TrunkID,
		FromUser:  "0033988776655",
		ToURI:     "+33612345678",
		EnableRTP: true,
	}

	call, err := svc.InitiateSIPCall(context.Background(), "tenant_ovh", callReq)
	if err != nil {
		t.Fatalf("Unexpected error initiating SIP call: %v", err)
	}

	if call.Status != "CONNECTED" {
		t.Errorf("Expected CONNECTED status, got %s", call.Status)
	}
}
