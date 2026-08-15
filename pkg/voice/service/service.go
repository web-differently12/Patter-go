package service

import (
	"context"

	"github.com/google/uuid"
)

type InitiateCallRequest struct {
	FromNumber          string  `json:"from_number" binding:"required"`
	ToNumber            string  `json:"to_number" binding:"required"`
	SystemPrompt        string  `json:"system_prompt,omitempty"`
	EnableRAG           bool    `json:"enable_rag"`
	AsyncToolExecution  bool    `json:"async_tool_execution"`
	VoiceStyle          string  `json:"voice_style,omitempty"`
	SpeechSpeed         float64 `json:"speech_speed,omitempty"`          // e.g. 1.0, 1.2x
	BargeInSensitivity  int     `json:"barge_in_sensitivity_ms,omitempty"` // e.g. 1500ms
	EnergyVADThreshold  float64 `json:"energy_vad_threshold,omitempty"`  // e.g. 0.02 energy level
	RecordingEnabled    bool    `json:"recording_enabled"`
	TransferDestination string  `json:"human_transfer_destination,omitempty"`
}

type InitiateCallResponse struct {
	CallID              string `json:"call_id"`
	Status              string `json:"status"` // "QUEUED", "IN_PROGRESS", "COMPLETED"
	InitialFillerSpeech string `json:"initial_filler_speech,omitempty"`
}

type VoiceService interface {
	InitiateCall(ctx context.Context, tenantID string, req InitiateCallRequest) (*InitiateCallResponse, error)
}

type voiceService struct{}

func NewVoiceService() VoiceService {
	return &voiceService{}
}

func (s *voiceService) InitiateCall(ctx context.Context, tenantID string, req InitiateCallRequest) (*InitiateCallResponse, error) {
	callID := "call_" + uuid.New().String()[:8]
	return &InitiateCallResponse{
		CallID:              callID,
		Status:              "QUEUED",
		InitialFillerSpeech: "Bonjour ! Je suis votre assistant virtuel, un instant je charge notre dossier...",
	}, nil
}
