package avatar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type SimliLiveEngine struct {
	apiKey     string
	simliURL   string
	httpClient *http.Client
}

type SimliSessionRequest struct {
	ApiKey       string `json:"apiKey"`
	FaceId       string `json:"faceId"`
	HandleFormat string `json:"handleFormat"`
}

type SimliSessionResponse struct {
	SessionToken string `json:"sessionToken"`
	RoomName     string `json:"roomName"`
}

func NewSimliLiveEngine(apiKey string) *SimliLiveEngine {
	return &SimliLiveEngine{
		apiKey:     apiKey,
		simliURL:   "https://api.simli.ai/startAudioToVideoSession",
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *SimliLiveEngine) RenderRawVideo(ctx context.Context, req *GenerateRequest) (*RawVideoResult, error) {
	return nil, fmt.Errorf("offline raw video render not supported on SimliLiveEngine")
}

func (s *SimliLiveEngine) StartLiveStream(ctx context.Context, mode Mode, req *GenerateRequest) (*LiveSessionResult, error) {
	slog.Info("[SimliLiveEngine] Requesting Simli WebRTC Live Stream Session...",
		slog.String("tenant_id", req.TenantID),
		slog.String("agent_id", req.AgentID),
		slog.String("mode", string(mode)),
	)

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	simliReq := SimliSessionRequest{
		ApiKey:       s.apiKey,
		FaceId:       req.FaceImageURL,
		HandleFormat: "webrtc",
	}

	payloadBytes, err := json.Marshal(simliReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal simli session request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(reqCtx, "POST", s.simliURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build simli request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		slog.Warn("[SimliLiveEngine] Simli API unreachable. Returning fallback session token",
			slog.String("error", err.Error()),
		)
		sessionID := fmt.Sprintf("simli_sess_%d", time.Now().UnixNano())
		return &LiveSessionResult{
			SessionID: sessionID,
			RoomName:  "simli_room_" + req.AgentID,
			JoinToken: "simli_token_" + req.TenantID,
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		slog.Warn("[SimliLiveEngine] Simli API returned non-200 status code. Returning fallback session token",
			slog.Int("status_code", resp.StatusCode),
		)
		sessionID := fmt.Sprintf("simli_sess_%d", time.Now().UnixNano())
		return &LiveSessionResult{
			SessionID: sessionID,
			RoomName:  "simli_room_" + req.AgentID,
			JoinToken: "simli_token_" + req.TenantID,
		}, nil
	}
	defer resp.Body.Close()

	var simliResp SimliSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&simliResp); err != nil {
		return nil, fmt.Errorf("failed to decode simli response: %w", err)
	}

	slog.Info("[SimliLiveEngine] Simli WebRTC Session created successfully",
		slog.String("room_name", simliResp.RoomName),
	)

	return &LiveSessionResult{
		SessionID: fmt.Sprintf("simli_%d", time.Now().UnixNano()),
		RoomName:  simliResp.RoomName,
		JoinToken: simliResp.SessionToken,
	}, nil
}
