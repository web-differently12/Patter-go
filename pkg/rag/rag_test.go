package rag_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lynxflow/patter-go/pkg/rag"
)

type mockTransferer struct {
	transferred bool
	destination string
	summary     string
}

func (m *mockTransferer) ExecuteTransfer(ctx context.Context, tenantID, callSID, destination, transferType, summary string) error {
	m.transferred = true
	m.destination = destination
	m.summary = summary
	return nil
}

func TestSearchAndDecide_HighConfidence(t *testing.T) {
	router := rag.NewUnifiedRAGRouter(nil, nil, nil)

	cfg := rag.KnowledgeBaseConfig{
		Provider: rag.ProviderGoogleVertex,
		FallbackTransfer: rag.FallbackTransferPolicy{
			EnableAutoTransfer:  true,
			ConfidenceThreshold: 0.75,
			TransferDestination: "+33612345678",
		},
		TopK: 2,
	}

	res, err := router.SearchAndDecide(context.Background(), "tenant_1", "call_123", "Comment configurer RAG ?", []float32{0.1, 0.2}, cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.TransferNeeded {
		t.Errorf("Expected TransferNeeded false for high confidence, got true")
	}
	if len(res.Chunks) == 0 {
		t.Errorf("Expected chunks returned, got empty")
	}
}

func TestSearchAndDecide_LowConfidenceAutoTransfer(t *testing.T) {
	mockTrans := &mockTransferer{}
	router := rag.NewUnifiedRAGRouter(nil, mockTrans, nil)

	// Custom endpoint returning score below threshold
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"documents": []string{"Info incertaine"},
			"score":     0.45,
		})
	}))
	defer ts.Close()

	cfg := rag.KnowledgeBaseConfig{
		Provider:        rag.ProviderCustomEndpoint,
		CustomServerURL: ts.URL,
		FallbackTransfer: rag.FallbackTransferPolicy{
			EnableAutoTransfer:  true,
			ConfidenceThreshold: 0.75,
			TransferDestination: "+33612345678",
			TransferType:        "WARM",
		},
	}

	res, err := router.SearchAndDecide(context.Background(), "tenant_1", "call_123", "Question complexe introuvable", []float32{}, cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !res.TransferNeeded {
		t.Errorf("Expected TransferNeeded true when confidence score (0.45) < threshold (0.75)")
	}
	if !mockTrans.transferred {
		t.Errorf("Expected mockTransferer to be invoked for call transfer")
	}
	if mockTrans.destination != "+33612345678" {
		t.Errorf("Expected destination +33612345678, got %s", mockTrans.destination)
	}
}
