package service

import (
	"context"

	"github.com/google/uuid"
)

type VoiceCallRequest struct {
	FromNumber string `json:"from_number" binding:"required"`
	ToNumber   string `json:"to_number" binding:"required"`
	AgentID    string `json:"agent_id,omitempty"`
	Prompt     string `json:"prompt,omitempty"`
}

type VoiceCallResponse struct {
	CallID     string `json:"call_id"`
	Status     string `json:"status"` // "QUEUED", "RINGING", "IN_PROGRESS"
	FromNumber string `json:"from_number"`
	ToNumber   string `json:"to_number"`
}

type VoiceService interface {
	InitiateCall(ctx context.Context, tenantID string, req VoiceCallRequest) (*VoiceCallResponse, error)
}

type voiceService struct{}

func NewVoiceService() VoiceService {
	return &voiceService{}
}

func (s *voiceService) InitiateCall(ctx context.Context, tenantID string, req VoiceCallRequest) (*VoiceCallResponse, error) {
	callID := "call_" + uuid.New().String()[:8]
	return &VoiceCallResponse{
		CallID:     callID,
		Status:     "QUEUED",
		FromNumber: req.FromNumber,
		ToNumber:   req.ToNumber,
	}, nil
}
