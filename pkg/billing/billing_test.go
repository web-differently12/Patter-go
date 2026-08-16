package billing_test

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/billing"
)

func TestPricingGridAndMetering(t *testing.T) {
	pg := billing.NewPricingGrid()
	wallet := billing.NewWalletLedgerService(nil)
	meter := billing.NewUsageMeterService(pg, wallet, nil)
	bridge := billing.NewHyperswitchLagoService(pg, wallet, nil)
	ctx := context.Background()

	tenantID := "tenant_test_billing"

	// 1. Get initial wallet
	w, err := wallet.GetWallet(ctx, tenantID)
	if err != nil {
		t.Fatalf("unexpected error getting wallet: %v", err)
	}
	if w.BalanceFiat <= 0 {
		t.Errorf("expected positive initial welcome credit, got %f", w.BalanceFiat)
	}

	// 2. Set custom tenant markup (+30%)
	pg.SetTenantConfig(billing.TenantPricingConfig{
		TenantID:         tenantID,
		MarkupPercentage: 30.0,
		FixedFeePerCall:  0.01,
		GraceLimitFiat:   -0.50,
	})

	// 3. Ingest consumption chunk event (STT + LLM + TTS + SIP)
	evt := billing.MeteringEvent{
		TenantID:         tenantID,
		SessionID:        "sess_test_01",
		ChannelType:      billing.ChannelVoiceCall,
		Provider:         "openai",
		ModelName:        "gpt-4o",
		PromptTokens:     1000,
		CompletionTokens: 200,
		AudioDurationMs:  30000,
		TTSCharacters:    250,
		SIPDurationSec:   30,
	}

	result, err := meter.IngestEvent(ctx, evt)
	if err != nil {
		t.Fatalf("unexpected error ingesting event: %v", err)
	}

	if result.RawCostUSD <= 0 {
		t.Errorf("expected positive raw cost, got %f", result.RawCostUSD)
	}
	if result.BilledCostUSD <= result.RawCostUSD {
		t.Errorf("expected billed cost (%f) to be strictly greater than raw cost (%f) due to +30%% markup", result.BilledCostUSD, result.RawCostUSD)
	}
	if result.ProfitUSD <= 0 {
		t.Errorf("expected positive profit margin, got %f", result.ProfitUSD)
	}

	// 4. Test Hyperswitch Topup Session
	topupResp, err := bridge.InitHyperswitchTopup(ctx, tenantID, billing.CreateTopupRequest{
		AmountFiat:    100.0,
		PaymentMethod: billing.PaymentBitcoinLightning,
	})
	if err != nil {
		t.Fatalf("unexpected error initializing topup: %v", err)
	}
	if topupResp.PaymentID == "" {
		t.Errorf("expected valid payment ID")
	}
	if topupResp.SatoshisEquivalent <= 0 {
		t.Errorf("expected positive satoshis equivalent for bitcoin lightning payment")
	}

	// 5. Test Admin Provider Catalog Report
	reports, err := bridge.GetAdminProviderCostCatalog(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting catalog: %v", err)
	}
	if len(reports) == 0 {
		t.Errorf("expected non-empty provider catalog reports")
	}
}
