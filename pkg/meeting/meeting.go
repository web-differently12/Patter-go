package meeting

import (
	"context"
	"fmt"
	"time"
)

type MeetingProviderType string

const (
	MeetingProviderRecallAI MeetingProviderType = "recall_ai"
	MeetingProviderLiveKit  MeetingProviderType = "livekit_webrtc"
	MeetingProviderDyte     MeetingProviderType = "dyte_interactive"
)

type MeetingSession struct {
	SessionID   string              `json:"session_id"`
	Provider    MeetingProviderType `json:"provider"`
	RoomURL     string              `json:"room_url"`
	AccessToken string              `json:"access_token"`
	Status      string              `json:"status"` // "ACTIVE", "RECORDING", "ENDED"
	CreatedAt   time.Time           `json:"created_at"`
}

type Service interface {
	CreateRoomSession(ctx context.Context, provider MeetingProviderType, roomName string) (*MeetingSession, error)
	DispatchMeetingBot(ctx context.Context, meetingURL string, botName string) (string, error)
}

type MeetingService struct{}

func NewMeetingService() *MeetingService {
	return &MeetingService{}
}

func (m *MeetingService) CreateRoomSession(ctx context.Context, provider MeetingProviderType, roomName string) (*MeetingSession, error) {
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
	token := fmt.Sprintf("jwt_token_%s", roomName)

	switch provider {
	case MeetingProviderLiveKit:
		return &MeetingSession{
			SessionID:   sessionID,
			Provider:    provider,
			RoomURL:     "wss://livekit.patter.ai/" + roomName,
			AccessToken: token,
			Status:      "ACTIVE",
			CreatedAt:   time.Now(),
		}, nil
	case MeetingProviderDyte:
		return &MeetingSession{
			SessionID:   sessionID,
			Provider:    provider,
			RoomURL:     "https://app.dyte.io/v2/meeting?id=" + roomName,
			AccessToken: token,
			Status:      "ACTIVE",
			CreatedAt:   time.Now(),
		}, nil
	default:
		return &MeetingSession{
			SessionID:   sessionID,
			Provider:    MeetingProviderRecallAI,
			RoomURL:     "https://recall.ai/bot_room/" + roomName,
			AccessToken: token,
			Status:      "ACTIVE",
			CreatedAt:   time.Now(),
		}, nil
	}
}

func (m *MeetingService) DispatchMeetingBot(ctx context.Context, meetingURL string, botName string) (string, error) {
	botID := fmt.Sprintf("bot_recall_%d", time.Now().Unix())
	return botID, nil
}
