package avatar

import (
	"context"
	"fmt"
	"time"
)

type LiveKitTransport struct {
	host     string
	apiKey   string
	apiSecret string
}

func NewLiveKitTransport(host, apiKey, apiSecret string) *LiveKitTransport {
	return &LiveKitTransport{
		host:      host,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (l *LiveKitTransport) CreateVideoTrackRoom(ctx context.Context, roomName string) (*LiveSessionResult, error) {
	token := fmt.Sprintf("lk_token_%d", time.Now().Unix())
	return &LiveSessionResult{
		SessionID:  fmt.Sprintf("lk_session_%s", roomName),
		RoomName:   roomName,
		JoinToken:  token,
		LiveKitURL: l.host,
	}, nil
}
