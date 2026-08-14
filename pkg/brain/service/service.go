package service

import (
	"context"
	"fmt"
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
}

type brainService struct{}

func NewBrainService() BrainService {
	return &brainService{}
}

func (s *brainService) QueryBrain(ctx context.Context, tenantID string, req QueryBrainRequest) (*QueryBrainResponse, error) {
	answer := fmt.Sprintf("Reponse IA White-Label pour Tenant [%s]: Traitement de '%s'", tenantID, req.Query)
	return &QueryBrainResponse{
		Answer:     answer,
		Confidence: 0.98,
		Sources:    []string{"RAG Knowledgebase", "LeMUR v3"},
	}, nil
}
