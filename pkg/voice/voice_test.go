package voice

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/voice/service"
)

func TestVoiceProfileConfiguration(t *testing.T) {
	svc := service.NewVoiceService()
	ctx := context.Background()

	prof, err := svc.CreateVoiceProfile(ctx, "tenant_test", service.VoiceProfile{
		Name:                 "Thomas Pro — Support Technique",
		TTSEngine:            "elevenlabs",
		VoiceID:              "thomas_pro_hd",
		STTEngine:            "deepgram_nova2",
		NoiseSuppression:     service.NoiseSuppressionKrisp,
		AmbientBackground:    service.AmbientSoundOffice,
		PitchTuning:          1.5,
		SpeechSpeed:          1.1,
		TurnTakingStyle:      "EAGER_INTERRUPT",
		EnableBackchanneling:  true,
		BargeInSensitivityMS: 1000,
	})

	if err != nil {
		t.Fatalf("unexpected error creating voice profile: %v", err)
	}
	if prof.ProfileID == "" {
		t.Fatalf("expected generated profile ID")
	}

	fetched, err := svc.GetVoiceProfile(ctx, "tenant_test", prof.ProfileID)
	if err != nil {
		t.Fatalf("unexpected error fetching voice profile: %v", err)
	}
	if fetched.TTSEngine != "elevenlabs" {
		t.Errorf("expected elevenlabs TTS engine, got %s", fetched.TTSEngine)
	}
	if fetched.NoiseSuppression != service.NoiseSuppressionKrisp {
		t.Errorf("expected krisp noise suppression")
	}
}
