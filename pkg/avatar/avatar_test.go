package avatar_test

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/avatar"
)

func TestAvatarEngineWaveSpeed(t *testing.T) {
	svc := avatar.NewAvatarEngineService("", nil)
	ctx := context.Background()

	streamResp, err := svc.StreamRealtimeAvatar(ctx, "tenant_test", avatar.StreamAvatarRequest{
		AvatarID:          "avatar_sarah_hd",
		Mode:              avatar.ModeWaveSpeedRealtimeStreaming,
		PromptText:        "Bonjour, bienvenue sur notre plateforme!",
		InteractiveStream: true,
	})

	if err != nil {
		t.Fatalf("unexpected error creating realtime avatar stream: %v", err)
	}
	if streamResp.JobID == "" {
		t.Errorf("expected generated JobID")
	}
	if streamResp.ProviderUsed != "wavespeed_api_v2" {
		t.Errorf("expected wavespeed_api_v2 provider, got %s", streamResp.ProviderUsed)
	}

	renderResp, err := svc.RenderMassAvatarVideo(ctx, "tenant_test", avatar.RenderAvatarVideoRequest{
		AvatarID: "avatar_mass_prod",
		Prompts:  []string{"Prompt 1", "Prompt 2"},
	})

	if err != nil {
		t.Fatalf("unexpected error queuing mass avatar video render: %v", err)
	}
	if renderResp.ProviderUsed != "local_gpu_renderer_mcp" {
		t.Errorf("expected local_gpu_renderer_mcp provider, got %s", renderResp.ProviderUsed)
	}
}
