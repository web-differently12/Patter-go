package avatar

import (
	"context"
	"fmt"
	"time"
)

type SimliLiveEngine struct {
	apiKey string
}

func NewSimliLiveEngine(apiKey string) *SimliLiveEngine {
	return &SimliLiveEngine{apiKey: apiKey}
}

func (s *SimliLiveEngine) StartSession(ctx context.Context, req *GenerateRequest) (*LiveSessionResult, error) {
	sessionID := fmt.Sprintf("simli_sess_%d", time.Now().UnixNano())
	return &LiveSessionResult{
		SessionID: sessionID,
		RoomName:  "simli_room_" + req.AgentID,
		JoinToken: "simli_token_" + req.TenantID,
	}, nil
}
