package rtc_test

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/rtc"
)

func TestRecallAIBotLifecycle(t *testing.T) {
	svc := rtc.NewRecallAIService("", nil)

	req := rtc.CreateMeetingBotRequest{
		MeetingURL:               "https://meet.google.com/xyz-uvwx-rst",
		BotName:                  "Assistant Lynxflow",
		Platform:                 rtc.PlatformGoogleMeet,
		RecordingMode:            rtc.RecordingSpeakerView,
		EnableRealtimeTranscript: true,
		Language:                 "fr",
	}

	bot, err := svc.CreateMeetingBot(context.Background(), "tenant_test", req)
	if err != nil {
		t.Fatalf("Unexpected error creating bot: %v", err)
	}

	if bot.BotID == "" {
		t.Errorf("Expected valid BotID, got empty")
	}
	if bot.Language != "fr" {
		t.Errorf("Expected language 'fr', got %s", bot.Language)
	}

	bots, err := svc.ListBots(context.Background(), "tenant_test")
	if err != nil {
		t.Fatalf("Unexpected error listing bots: %v", err)
	}

	if len(bots) != 1 {
		t.Errorf("Expected 1 bot, got %d", len(bots))
	}

	err = svc.LeaveMeeting(context.Background(), "tenant_test", bot.BotID)
	if err != nil {
		t.Fatalf("Unexpected error leaving meeting: %v", err)
	}

	status, _ := svc.GetBotStatus(context.Background(), "tenant_test", bot.BotID)
	if status.Status != rtc.StatusLeft {
		t.Errorf("Expected StatusLeft, got %s", status.Status)
	}
}

func TestMeetingProfileConfiguration(t *testing.T) {
	svc := rtc.NewRecallAIService("", nil)
	ctx := context.Background()

	prof, err := svc.CreateMeetingProfile(ctx, "tenant_test", rtc.MeetingProfile{
		Name:                        "Profil Visio Commercial BANT",
		EnableSpeakerDiarization:    true,
		EnableActionItemsExtraction: true,
		EnableParticipantSentiment:  true,
		EnableLiveTranslation:       true,
		TargetTranslationLanguage:   "fr",
		EnableScreenShareRecording:  true,
		SummaryTemplate:             "BANT_QUALIFICATION",
		AutoLeaveOnSilenceMinutes:   10,
		AutoLeaveWhenEveryoneLeft:   true,
	})

	if err != nil {
		t.Fatalf("Unexpected error creating meeting profile: %v", err)
	}
	if prof.ProfileID == "" {
		t.Fatalf("Expected valid ProfileID")
	}

	fetched, err := svc.GetMeetingProfile(ctx, "tenant_test", prof.ProfileID)
	if err != nil {
		t.Fatalf("Unexpected error fetching meeting profile: %v", err)
	}
	if fetched.SummaryTemplate != "BANT_QUALIFICATION" {
		t.Errorf("Expected BANT_QUALIFICATION summary template")
	}
}
