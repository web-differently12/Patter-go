package engine

import (
	"encoding/base64"
	"testing"
)

func TestBufferPoolRecycling(t *testing.T) {
	// Retrieve a buffer from sync.Pool
	buf := bufferPool.Get().([]byte)
	if len(buf) != 8192 {
		t.Errorf("expected buffer size 8192, got %d", len(buf))
	}

	// Test decoding raw PCM u-law audio Base64 chunks (Twilio format)
	// Let's encode a mock PCM chunk
	mockData := []byte("some-mock-twilio-pcmu-audio-payload-chunk")
	base64Payload := base64.StdEncoding.EncodeToString(mockData)

	decodedLen, err := base64.StdEncoding.Decode(buf, []byte(base64Payload))
	if err != nil {
		t.Fatalf("failed to decode base64 payload: %v", err)
	}

	decodedString := string(buf[:decodedLen])
	if decodedString != string(mockData) {
		t.Errorf("expected decoded string %q, got %q", string(mockData), decodedString)
	}

	// Put buffer back into sync.Pool
	bufferPool.Put(buf)
}
