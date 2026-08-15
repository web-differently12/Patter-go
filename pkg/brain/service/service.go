package service

import (
	"context"
	"fmt"

	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/rag"
	"github.com/lynxflow/patter-go/pkg/skills"
)

type QueryBrainRequest struct {
	Query     string            `json:"query" binding:"required"`
	Context   map[string]string `json:"context,omitempty"`
	Channel   string            `json:"channel,omitempty"`
}

type QueryBrainResponse struct {
	Answer     string   `json:"answer"`
	Confidence float64  `json:"confidence"`
	Sources    []string `json:"sources,omitempty"`
}

type BrainService interface {
	QueryBrain(ctx context.Context, tenantID string, req QueryBrainRequest) (*QueryBrainResponse, error)
	SearchRAG(ctx context.Context, tenantID, callSID, query string, cfg rag.KnowledgeBaseConfig) (*rag.SearchResult, error)
	ExecuteHumanTransfer(ctx context.Context, tenantID, callSID, rawArgs, destination string) (string, error)
}

type brainService struct {
	ragRouter   *rag.UnifiedRAGRouter
	transferSkill *skills.HumanTransferSkill
}

func NewBrainService() BrainService {
	logger := core.GetLogger()
	router := rag.NewUnifiedRAGRouter(nil, nil, logger)
	skill := skills.NewHumanTransferSkill(nil, logger)

	return &brainService{
		ragRouter:   router,
		transferSkill: skill,
	}
}

func (s *brainService) QueryBrain(ctx context.Context, tenantID string, req QueryBrainRequest) (*QueryBrainResponse, error) {
	answer := fmt.Sprintf("Reponse IA White-Label pour Tenant [%s]: Traitement de '%s'", tenantID, req.Query)
	return &QueryBrainResponse{
		Answer:     answer,
		Confidence: 0.98,
		Sources:    []string{"RAG Knowledgebase", "LeMUR v3"},
	}, nil
}

func (s *brainService) SearchRAG(ctx context.Context, tenantID, callSID, query string, cfg rag.KnowledgeBaseConfig) (*rag.SearchResult, error) {
	return s.ragRouter.SearchAndDecide(ctx, tenantID, callSID, query, []float32{0.1, 0.2}, cfg)
}

func (s *brainService) ExecuteHumanTransfer(ctx context.Context, tenantID, callSID, rawArgs, destination string) (string, error) {
	return s.transferSkill.Execute(ctx, tenantID, callSID, rawArgs, destination)
}
