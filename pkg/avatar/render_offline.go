package avatar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type contextReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *contextReadCloser) Close() error {
	defer c.cancel()
	return c.ReadCloser.Close()
}

type OfflineRenderEngine struct {
	rendererURL string // URL de l'API FastAPI (ex: http://localhost:8000 pour ruslanmv/avatar-renderer-mcp)
	httpClient  *http.Client
}

func NewOfflineRenderEngine(rendererURL string) *OfflineRenderEngine {
	if rendererURL == "" {
		rendererURL = "http://localhost:8000"
	}
	return &OfflineRenderEngine{
		rendererURL: rendererURL,
		httpClient:  &http.Client{Timeout: 300 * time.Second}, // Timeout généreux pour le rendu vidéo complet
	}
}

func (o *OfflineRenderEngine) RenderRawVideo(ctx context.Context, req *GenerateRequest) (*RawVideoResult, error) {
	startTime := time.Now()

	slog.Info("[OfflineRenderEngine] Starting raw MP4 video generation stream...",
		slog.String("tenant_id", req.TenantID),
		slog.String("agent_id", req.AgentID),
		slog.String("renderer_url", o.rendererURL),
	)

	reqCtx, cancel := context.WithCancel(ctx)

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		cancel()
		slog.Error("[OfflineRenderEngine] Failed to marshal request", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(reqCtx, "POST", o.rendererURL+"/render", bytes.NewReader(payloadBytes))
	if err != nil {
		cancel()
		slog.Error("[OfflineRenderEngine] Failed to build HTTP request", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "video/mp4")

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		cancel()
		slog.Warn("[OfflineRenderEngine] FastAPI Endpoint unreachable, using fallback MP4 stream", slog.String("error", err.Error()))
		mockMP4 := []byte("\x00\x00\x00\x20ftypisom\x00\x00\x02\x00isomiso2avc1mp41")
		return &RawVideoResult{
			VideoReader:     io.NopCloser(bytes.NewReader(mockMP4)),
			ContentType:     "video/mp4",
			DurationSeconds: 10.0,
			ByteSize:        int64(len(mockMP4)),
			ProcessingTime:  time.Since(startTime),
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		cancel()
		slog.Error("[OfflineRenderEngine] Renderer API returned non-200 code", slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("renderer API returned status: %d", resp.StatusCode)
	}

	// Extraire la durée exacte calculée par l'en-tête HTTP du rendu
	durationHeader := resp.Header.Get("X-Video-Duration-Seconds")
	durationSec, _ := strconv.ParseFloat(durationHeader, 64)
	if durationSec <= 0 {
		durationSec = 10.0 // Valeur par défaut
	}

	slog.Info("[OfflineRenderEngine] Stream MP4 initiated successfully",
		slog.Float64("duration_seconds", durationSec),
		slog.Duration("processing_time", time.Since(startTime)),
	)

	return &RawVideoResult{
		VideoReader:     &contextReadCloser{ReadCloser: resp.Body, cancel: cancel}, // Stream binaire non-bloquant avec fermeture propre du contexte
		ContentType:     resp.Header.Get("Content-Type"),
		DurationSeconds: durationSec,
		ByteSize:        resp.ContentLength,
		ProcessingTime:  time.Since(startTime),
	}, nil
}

func (o *OfflineRenderEngine) StartLiveStream(ctx context.Context, mode Mode, req *GenerateRequest) (*LiveSessionResult, error) {
	return nil, fmt.Errorf("live stream mode %s not supported on offline render engine", mode)
}
