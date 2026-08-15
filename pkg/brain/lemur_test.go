package brain_test

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/brain"
)

func TestLeMURPostCallAnalytics(t *testing.T) {
	engine := brain.NewLeMUREngine("", nil)

	utterances := []brain.SpeakerUtterance{
		{Speaker: "Commercial", StartMS: 0, EndMS: 5000, Text: "Bonjour, bienvenue chez Lynxflow!"},
		{Speaker: "Prospect", StartMS: 5100, EndMS: 8000, Text: "Bonjour, je cherche une solution RAG."},
		{Speaker: "Commercial", StartMS: 8100, EndMS: 15000, Text: "Parfait, notre moteur Patter Go prend en charge Vertex Search et PGVector."},
	}

	resp, err := engine.ProcessPostCallAnalytics(context.Background(), "tenant_acme", "trx_998877", utterances)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.TranscriptID != "trx_998877" {
		t.Errorf("Expected TranscriptID trx_998877, got %s", resp.TranscriptID)
	}
	if resp.BANT == nil {
		t.Fatalf("Expected BANT qualification present")
	}
	if resp.BANT.DealScore <= 0 {
		t.Errorf("Expected DealScore > 0, got %f", resp.BANT.DealScore)
	}
	if resp.Analytics == nil {
		t.Fatalf("Expected AudioAnalytics present")
	}
	if len(resp.Analytics.TalkRatio) == 0 {
		t.Errorf("Expected non-empty TalkRatio analytics")
	}
}
