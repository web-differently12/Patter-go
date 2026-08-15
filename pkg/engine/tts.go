package engine

import (
	"context"
)

// TTSProvider defines the contract for real-time text-to-speech generation
type TTSProvider interface {
	GetProviderID() string
	// SynthesizeText converts incoming text characters to raw G711/PCM audio streams asynchronously
	SynthesizeText(ctx context.Context, text string, voice string) (<-chan []byte, error)
}

// ElevenLabsTTS represents the ElevenLabs implementation
type ElevenLabsTTS struct {
	apiKey string
}

func NewElevenLabsTTS(apiKey string) *ElevenLabsTTS {
	return &ElevenLabsTTS{apiKey: apiKey}
}

func (e *ElevenLabsTTS) GetProviderID() string {
	return "elevenlabs"
}

func (e *ElevenLabsTTS) SynthesizeText(ctx context.Context, text string, voice string) (<-chan []byte, error) {
	audioCh := make(chan []byte, 10)

	go func() {
		defer close(audioCh)
		// Send mock PCM u-law frame payload
		mockPCM := []byte{0xff, 0x00, 0xff, 0x00, 0xaa, 0x55}
		select {
		case <-ctx.Done():
		case audioCh <- mockPCM:
		}
	}()

	return audioCh, nil
}

// CartesiaTTS represents the Cartesia implementation
type CartesiaTTS struct {
	apiKey string
}

func NewCartesiaTTS(apiKey string) *CartesiaTTS {
	return &CartesiaTTS{apiKey: apiKey}
}

func (c *CartesiaTTS) GetProviderID() string {
	return "cartesia"
}

func (c *CartesiaTTS) SynthesizeText(ctx context.Context, text string, voice string) (<-chan []byte, error) {
	audioCh := make(chan []byte, 10)

	go func() {
		defer close(audioCh)
		mockPCM := []byte{0x7f, 0x80, 0x7f, 0x80}
		select {
		case <-ctx.Done():
		case audioCh <- mockPCM:
		}
	}()

	return audioCh, nil
}

// OpenAITTS represents the OpenAI TTS implementation
type OpenAITTS struct {
	apiKey string
}

func NewOpenAITTS(apiKey string) *OpenAITTS {
	return &OpenAITTS{apiKey: apiKey}
}

func (o *OpenAITTS) GetProviderID() string {
	return "openai_tts"
}

func (o *OpenAITTS) SynthesizeText(ctx context.Context, text string, voice string) (<-chan []byte, error) {
	audioCh := make(chan []byte, 10)

	go func() {
		defer close(audioCh)
		mockPCM := []byte{0x00, 0xff, 0x00, 0xff}
		select {
		case <-ctx.Done():
		case audioCh <- mockPCM:
		}
	}()

	return audioCh, nil
}
