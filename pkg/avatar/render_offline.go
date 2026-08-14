package avatar

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"time"
)

type OfflineRenderEngine struct {
	rendererURL string // e.g. http://localhost:8000 (ruslanmv/avatar-renderer-mcp)
	httpClient  *http.Client
}

func NewOfflineRenderEngine(rendererURL string) *OfflineRenderEngine {
	if rendererURL == "" {
		rendererURL = "http://localhost:8000"
	}
	return &OfflineRenderEngine{
		rendererURL: rendererURL,
		httpClient:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (o *OfflineRenderEngine) RenderRawVideo(ctx context.Context, req *GenerateRequest) (*RawVideoResult, error) {
	startTime := time.Now()

	// Simulate rendering on-the-fly MP4 byte stream via self-hosted FastAPI /render endpoint (FOMM / Hallo3 / GFPGAN)
	// Metered strictly per second of generated video with zero mandatory storage.
	mockMP4Data := []byte("\x00\x00\x00\x20ftypisom\x00\x00\x02\x00isomiso2avc1mp41") // Mock MP4 header

	durationSec := 10.0 // 10 seconds of metered video
	reader := io.NopCloser(bytes.NewReader(mockMP4Data))

	return &RawVideoResult{
		VideoReader:     reader,
		ContentType:     "video/mp4",
		DurationSeconds: durationSec,
		ByteSize:        int64(len(mockMP4Data)),
		ProcessingTime:  time.Since(startTime),
	}, nil
}

func (o *OfflineRenderEngine) StartLiveStream(ctx context.Context, mode Mode, req *GenerateRequest) (*LiveSessionResult, error) {
	return nil, fmt.Errorf("live stream mode %s not supported on offline render engine", mode)
}
