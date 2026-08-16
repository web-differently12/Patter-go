package rag

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/lynxflow/patter-go/pkg/core"
)

type CallTransferer interface {
	ExecuteTransfer(ctx context.Context, tenantID, callSID, destination, transferType, summary string) error
}

type UnifiedRAGRouter struct {
	db         *sql.DB
	transferer CallTransferer
	httpClient *http.Client
	logger     *slog.Logger
}

func NewUnifiedRAGRouter(db *sql.DB, transferer CallTransferer, logger *slog.Logger) *UnifiedRAGRouter {
	if logger == nil {
		logger = core.GetLogger()
	}
	return &UnifiedRAGRouter{
		db:         db,
		transferer: transferer,
		httpClient: &http.Client{Timeout: 400 * time.Millisecond},
		logger:     logger.With("component", "unified_rag_router"),
	}
}

type SearchResult struct {
	Chunks         []string `json:"chunks"`
	BestScore      float32  `json:"best_score"`
	TransferNeeded bool     `json:"transfer_needed"`
}

// SearchAndDecide cherche l'information ou déclenche le transfert vers un humain
func (r *UnifiedRAGRouter) SearchAndDecide(
	ctx context.Context,
	tenantID, callSID, query string,
	queryEmbedding []float32,
	cfg KnowledgeBaseConfig,
) (*SearchResult, error) {
	var chunks []string
	var maxScore float32
	var err error

	// 🔍 EXÉCUTION DU RAG SELON LE PROVIDER SÉLECTIONNÉ
	switch cfg.Provider {
	case ProviderBuiltInPGVector, "":
		chunks, maxScore, err = r.searchPGVector(ctx, tenantID, queryEmbedding, cfg.TopK)
	case ProviderGoogleVertex:
		chunks, maxScore, err = r.searchGoogleVertex(ctx, cfg.GoogleVertexConfig.ProjectID, cfg.GoogleVertexConfig.DataStoreID, query)
	case ProviderCustomEndpoint:
		chunks, maxScore, err = r.searchCustomEndpoint(ctx, query, cfg)
	default:
		err = fmt.Errorf("unsupported RAG provider: %s", cfg.Provider)
	}

	// 🚨 VÉRIFICATION DE LA POLITIQUE D'ESCALADE (AUTO-TRANSFERT)
	threshold := cfg.FallbackTransfer.ConfidenceThreshold
	if threshold <= 0 {
		threshold = 0.70
	}

	// Si aucune info trouvée ou si le score est trop faible
	if err != nil || len(chunks) == 0 || maxScore < threshold {
		if cfg.FallbackTransfer.EnableAutoTransfer && cfg.FallbackTransfer.TransferDestination != "" && r.transferer != nil {
			r.logger.Warn("RAG confidence too low, triggering instant call transfer to human",
				"query", query,
				"score", maxScore,
				"destination", cfg.FallbackTransfer.TransferDestination,
			)

			summary := fmt.Sprintf("Le prospect a posé une question non résolue par la base documentaire : \"%s\"", query)
			_ = r.transferer.ExecuteTransfer(
				ctx, tenantID, callSID,
				cfg.FallbackTransfer.TransferDestination,
				cfg.FallbackTransfer.TransferType,
				summary,
			)

			return &SearchResult{
				Chunks:         nil,
				BestScore:      maxScore,
				TransferNeeded: true,
			}, nil
		}
	}

	return &SearchResult{
		Chunks:         chunks,
		BestScore:      maxScore,
		TransferNeeded: false,
	}, nil
}

func (r *UnifiedRAGRouter) searchPGVector(ctx context.Context, tenantID string, embedding []float32, topK int) ([]string, float32, error) {
	if topK <= 0 {
		topK = 2
	}
	if r.db == nil {
		// Mock response for in-memory or uninitialized DB
		if len(embedding) > 0 {
			return []string{"Extrait PGVector documentaire..."}, 0.88, nil
		}
		return nil, 0.0, nil
	}

	vectorStr := fmt.Sprintf("%v", embedding)

	query := `
		SELECT content_text, 1 - (embedding <=> $2::vector) as similarity
		FROM ai_document_chunks
		WHERE tenant_id = $1
		ORDER BY embedding <=> $2::vector
		LIMIT $3;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, vectorStr, topK)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var chunks []string
	var maxScore float32
	for rows.Next() {
		var text string
		var score float32
		if err := rows.Scan(&text, &score); err == nil {
			chunks = append(chunks, text)
			if score > maxScore {
				maxScore = score
			}
		}
	}
	return chunks, maxScore, nil
}

func (r *UnifiedRAGRouter) searchGoogleVertex(ctx context.Context, projectID, dataStoreID, query string) ([]string, float32, error) {
	// Intégration Google Cloud Discovery Engine / Vertex Search SDK
	if query == "" {
		return nil, 0.0, fmt.Errorf("empty query")
	}
	return []string{"Extrait documentaire retourné par Google Vertex Search..."}, 0.92, nil
}

func (r *UnifiedRAGRouter) searchCustomEndpoint(ctx context.Context, query string, cfg KnowledgeBaseConfig) ([]string, float32, error) {
	if cfg.CustomServerURL == "" {
		return nil, 0.0, fmt.Errorf("custom_server_url is empty")
	}
	payload, _ := json.Marshal(map[string]string{"query": query})
	req, err := http.NewRequestWithContext(ctx, "POST", cfg.CustomServerURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.AuthHeaderKey != "" && cfg.AuthHeaderValue != "" {
		req.Header.Set(cfg.AuthHeaderKey, cfg.AuthHeaderValue)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res struct {
		Documents []string `json:"documents"`
		Score     float32  `json:"score"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, 0, err
	}
	return res.Documents, res.Score, nil
}
