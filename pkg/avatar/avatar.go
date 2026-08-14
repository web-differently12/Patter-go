package avatar

import (
	"context"
	"io"
	"time"
)

type Mode string

const (
	ModeLiveStreamSimli   Mode = "live_stream_simli"
	ModeLiveStreamLiveKit Mode = "live_stream_livekit_musetalk" // MuseTalk streamé via LiveKit WebRTC
	ModeOfflineRenderRaw  Mode = "offline_render_raw_mp4"       // Rendu brut à la volée (ruslanmv)
)

type GenerateRequest struct {
	TenantID     string            `json:"tenant_id"`
	AgentID      string            `json:"agent_id"`
	FaceImageURL string            `json:"face_image_url"` // Photo de l'Employé IA
	ScriptText   string            `json:"script_text"`    // Texte prononcé
	VoiceID      string            `json:"voice_id"`       // Voix TTS (ElevenLabs, CosyVoice...)
	QualityMode  string            `json:"quality_mode"`   // "auto", "high_quality", "cinematic"
	Enhancements []string          `json:"enhancements"`   // ["emotion_expressions", "eye_gaze_blink", ...]
	Variables    map[string]string `json:"variables"`      // Variables résolues (first_name, company)
}

// Données brutes renvoyées à l'appelant (Zéro stockage obligatoire côté Go)
type RawVideoResult struct {
	VideoReader     io.ReadCloser `json:"-"`                // Stream binaire brut MP4
	ContentType     string        `json:"content_type"`     // "video/mp4"
	DurationSeconds float64       `json:"duration_seconds"` // Durée exacte pour facturation / métriques
	ByteSize        int64         `json:"byte_size"`
	ProcessingTime  time.Duration `json:"processing_time"`
}

// Session WebRTC Live Stream
type LiveSessionResult struct {
	SessionID  string `json:"session_id"`
	RoomName   string `json:"room_name"`
	JoinToken  string `json:"join_token"`
	LiveKitURL string `json:"livekit_url,omitempty"`
}

// Interface centralisée du moteur Avatar
type Engine interface {
	// 1. Génération Vidéo brute à la volée (Auto-hébergé ruslanmv/avatar-renderer, MP4 streamé sans stockage obligatoire)
	RenderRawVideo(ctx context.Context, req *GenerateRequest) (*RawVideoResult, error)
	// 2. Lancement d'un Stream WebRTC interactif (Simli ou LiveKit+MuseTalk)
	StartLiveStream(ctx context.Context, mode Mode, req *GenerateRequest) (*LiveSessionResult, error)
}
