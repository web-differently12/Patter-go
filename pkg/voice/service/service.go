package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type NoiseSuppressionMode string

const (
	NoiseSuppressionKrisp   NoiseSuppressionMode = "krisp"
	NoiseSuppressionRNNoise NoiseSuppressionMode = "rnnoise"
	NoiseSuppressionOff     NoiseSuppressionMode = "off"
)

type AmbientSoundType string

const (
	AmbientSoundOffice    AmbientSoundType = "office_hum"
	AmbientSoundCafe      AmbientSoundType = "cafe_chatter"
	AmbientSoundQuietRoom AmbientSoundType = "quiet_room"
	AmbientSoundOff       AmbientSoundType = "off"
)

type VoiceProfile struct {
	ProfileID           string               `json:"profile_id"`
	TenantID            string               `json:"tenant_id"`
	Name                string               `json:"name" binding:"required"`
	TTSEngine           string               `json:"tts_engine"`           // "elevenlabs", "deepgram", "cartesia", "openai", "azure"
	VoiceID             string               `json:"voice_id"`             // e.g. "sarah_hd", "thomas_pro"
	STTEngine           string               `json:"stt_engine"`           // "whisper", "deepgram_nova2", "assemblyai"
	NoiseSuppression    NoiseSuppressionMode `json:"noise_suppression"`    // "krisp", "rnnoise", "off"
	AmbientBackground   AmbientSoundType     `json:"ambient_background"`   // "office_hum", "cafe_chatter", "off"
	PitchTuning         float64              `json:"pitch_tuning"`         // -5.0 to +5.0
	SpeechSpeed         float64              `json:"speech_speed"`         // 0.5x to 2.0x
	TurnTakingStyle     string               `json:"turn_taking_style"`    // "EAGER_INTERRUPT", "BALANCED", "PATIENT_LISTENER"
	EnableBackchanneling bool                 `json:"enable_backchanneling"`// "hmm", "d'accord", "je vois"
	BargeInSensitivityMS int                 `json:"barge_in_sensitivity_ms"`
	FallbackTTSEngine   string               `json:"fallback_tts_engine"`
}

type InitiateCallRequest struct {
	FromNumber          string  `json:"from_number" binding:"required"`
	ToNumber            string  `json:"to_number" binding:"required"`
	ProfileID           string  `json:"profile_id,omitempty"`
	SystemPrompt        string  `json:"system_prompt,omitempty"`
	EnableRAG           bool    `json:"enable_rag"`
	AsyncToolExecution  bool    `json:"async_tool_execution"`
	VoiceStyle          string  `json:"voice_style,omitempty"`
	SpeechSpeed         float64 `json:"speech_speed,omitempty"`
	BargeInSensitivity  int     `json:"barge_in_sensitivity_ms,omitempty"`
	EnergyVADThreshold  float64 `json:"energy_vad_threshold,omitempty"`
	RecordingEnabled    bool    `json:"recording_enabled"`
	TransferDestination string  `json:"human_transfer_destination,omitempty"`
}

type InitiateCallResponse struct {
	CallID              string `json:"call_id"`
	ProfileID           string `json:"profile_id,omitempty"`
	Status              string `json:"status"` // "QUEUED", "IN_PROGRESS", "COMPLETED"
	InitialFillerSpeech string `json:"initial_filler_speech,omitempty"`
}

type VoiceService interface {
	CreateVoiceProfile(ctx context.Context, tenantID string, profile VoiceProfile) (*VoiceProfile, error)
	GetVoiceProfile(ctx context.Context, tenantID, profileID string) (*VoiceProfile, error)
	ListVoiceProfiles(ctx context.Context, tenantID string) ([]*VoiceProfile, error)
	InitiateCall(ctx context.Context, tenantID string, req InitiateCallRequest) (*InitiateCallResponse, error)
}

type voiceService struct {
	mu       sync.RWMutex
	profiles map[string]*VoiceProfile
}

func NewVoiceService() VoiceService {
	svc := &voiceService{
		profiles: make(map[string]*VoiceProfile),
	}

	// Pre-populate default Voice Profile
	defaultID := "vp_default_hd"
	svc.profiles["default_tenant:"+defaultID] = &VoiceProfile{
		ProfileID:            defaultID,
		TenantID:             "default_tenant",
		Name:                 "Sarah HD — Commercial & Support",
		TTSEngine:            "elevenlabs",
		VoiceID:              "sarah_natural_fr",
		STTEngine:            "deepgram_nova2",
		NoiseSuppression:     NoiseSuppressionKrisp,
		AmbientBackground:    AmbientSoundQuietRoom,
		PitchTuning:          0.0,
		SpeechSpeed:          1.0,
		TurnTakingStyle:      "BALANCED",
		EnableBackchanneling:  true,
		BargeInSensitivityMS: 1200,
		FallbackTTSEngine:    "openai",
	}

	return svc
}

func (s *voiceService) CreateVoiceProfile(ctx context.Context, tenantID string, profile VoiceProfile) (*VoiceProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	profileID := "vp_" + uuid.New().String()[:8]
	profile.ProfileID = profileID
	profile.TenantID = tenantID
	if profile.TTSEngine == "" {
		profile.TTSEngine = "elevenlabs"
	}
	if profile.STTEngine == "" {
		profile.STTEngine = "deepgram_nova2"
	}
	if profile.TurnTakingStyle == "" {
		profile.TurnTakingStyle = "BALANCED"
	}

	s.profiles[tenantID+":"+profileID] = &profile
	return &profile, nil
}

func (s *voiceService) GetVoiceProfile(ctx context.Context, tenantID, profileID string) (*VoiceProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prof, ok := s.profiles[tenantID+":"+profileID]
	if !ok {
		return nil, fmt.Errorf("voice profile %s not found", profileID)
	}
	return prof, nil
}

func (s *voiceService) ListVoiceProfiles(ctx context.Context, tenantID string) ([]*VoiceProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*VoiceProfile
	for _, prof := range s.profiles {
		if prof.TenantID == tenantID {
			result = append(result, prof)
		}
	}
	return result, nil
}

func (s *voiceService) InitiateCall(ctx context.Context, tenantID string, req InitiateCallRequest) (*InitiateCallResponse, error) {
	callID := "call_" + uuid.New().String()[:8]
	return &InitiateCallResponse{
		CallID:              callID,
		ProfileID:           req.ProfileID,
		Status:              "QUEUED",
		InitialFillerSpeech: "Bonjour ! Je suis votre assistant virtuel, un instant je charge notre dossier...",
	}, nil
}
