package avatar

import (
	"context"
	"fmt"
	"time"
)

type MuseTalkEngine struct {
	endpoint string
}

func NewMuseTalkEngine(endpoint string) *MuseTalkEngine {
	return &MuseTalkEngine{endpoint: endpoint}
}

func (m *MuseTalkEngine) StartInferenceStream(ctx context.Context, req *GenerateRequest) (*LiveSessionResult, error) {
	sessionID := fmt.Sprintf("musetalk_sess_%d", time.Now().UnixNano())
	return &LiveSessionResult{
		SessionID:  sessionID,
		RoomName:   "musetalk_room_" + req.AgentID,
		JoinToken:  "musetalk_jwt_" + req.TenantID,
		LiveKitURL: "wss://livekit.patter.ai/musetalk",
	}, nil
}
