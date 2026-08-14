package engine

import (
	"context"
	"log"
	"math"
	"sync"
)

// PipelineOrchestrator coordinates Pipeline Mode: STT -> Fallback LLM -> TTS orchestration
type PipelineOrchestrator struct {
	stt      STTProvider
	tts      TTSProvider
	fallback *FallbackManager
	mu       sync.Mutex
}

// NewPipelineOrchestrator creates a new orchestrator
func NewPipelineOrchestrator(stt STTProvider, tts TTSProvider, fallback *FallbackManager) *PipelineOrchestrator {
	return &PipelineOrchestrator{
		stt:      stt,
		tts:      tts,
		fallback: fallback,
	}
}

// IsSpeechDetected performs energy-based Voice Activity Detection (VAD) on u-law/PCM audio chunks
func (po *PipelineOrchestrator) IsSpeechDetected(chunk []byte, threshold float64) bool {
	if len(chunk) == 0 {
		return false
	}

	// Convert G711 u-law sample values to PCM linear energy level approximation
	var sumSquares float64
	for _, b := range chunk {
		// u-law decode approximation
		u := ^b
		sign := (u & 0x80) != 0
		exponent := (u >> 4) & 0x07
		mantissa := u & 0x0F
		sample := int16(mantissa << 3)
		if exponent > 0 {
			sample += 0x84
			sample <<= exponent - 1
		}
		if sign {
			sample = -sample
		}

		val := float64(sample) / 32768.0
		sumSquares += val * val
	}

	rms := math.Sqrt(sumSquares / float64(len(chunk)))
	return rms > threshold
}

// RunStreamSession orchestrates the audio streaming loop asynchronously
func (po *PipelineOrchestrator) RunStreamSession(ctx context.Context, agent *Agent, inboundAudio <-chan []byte, outboundAudio chan<- []byte) {
	log.Printf("[PipelineOrchestrator] Starting stream session for agent: %s (%s mode)", agent.Name, agent.Mode)

	sttInput, sttOutput, err := po.stt.StartStream(ctx, 8000, 1)
	if err != nil {
		log.Printf("[PipelineOrchestrator] Failed to start STT stream: %v", err)
		return
	}
	defer close(sttInput)

	// Sub-goroutine to feed inbound audio into STT stream and perform live VAD
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case chunk, ok := <-inboundAudio:
				if !ok {
					return
				}

				// Apply Voice Activity Detection (VAD)
				isSpeaking := po.IsSpeechDetected(chunk, agent.VADSensitivity)
				if isSpeaking {
					// Audio chunk has live user voice, forward to STT engine
					select {
					case sttInput <- chunk:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Loop to handle transcription results, generate completions using the LLM with failover fallback chain,
	// and feed text segments into the TTS engine to stream outbound audio back to Twilio.
	for {
		select {
		case <-ctx.Done():
			return
		case transcript, ok := <-sttOutput:
			if !ok {
				return
			}

			if transcript.IsFinal && transcript.Text != "" {
				log.Printf("[PipelineOrchestrator] Final transcription received: %q. Routing to LLM...", transcript.Text)

				// 1. Get LLM response with automated sequential failover capabilities
				resp, activeProvider, err := po.fallback.ExecuteWithFallback(ctx, agent.LLMProviderID, agent.FallbackChain, agent.SystemPrompt, transcript.Text, agent.Tools)
				if err != nil {
					log.Printf("[PipelineOrchestrator] Failover chain failed: %v", err)
					continue
				}

				log.Printf("[PipelineOrchestrator] Selected LLM: %s. Response text: %q", activeProvider, resp.Text)

				// 2. Synthesize text response using TTS engine
				ttsAudioCh, err := po.tts.SynthesizeText(ctx, resp.Text, agent.Voice)
				if err != nil {
					log.Printf("[PipelineOrchestrator] TTS synthesis failed: %v", err)
					continue
				}

				// 3. Stream output audio blocks back to Twilio outbound voice channel
				for voiceChunk := range ttsAudioCh {
					select {
					case outboundAudio <- voiceChunk:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}
}
