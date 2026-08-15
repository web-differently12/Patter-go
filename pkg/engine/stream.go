package engine

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"vocal-engine/pkg/config"
	"vocal-engine/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Twilio Stream WebSocket protocol messages
type TwilioMedia struct {
	Track   string `json:"track"`
	Chunk   string `json:"chunk"`
	Payload string `json:"payload"`
}

type TwilioEvent struct {
	Event     string       `json:"event"`
	Sequence  string       `json:"sequenceNumber"`
	Start     *TwilioStart `json:"start,omitempty"`
	Media     *TwilioMedia `json:"media,omitempty"`
	StreamSID string       `json:"streamSid,omitempty"`
}

type TwilioStart struct {
	StreamSID   string `json:"streamSid"`
	AccountSID  string `json:"accountSid"`
	CallSID     string `json:"callSid"`
	Tracks      []string `json:"tracks"`
	CustomParameters map[string]string `json:"customParameters"`
}

// Twilio outbound stream clear message format (used for interrupting playbacks)
type TwilioClear struct {
	Event     string `json:"event"` // Must be "clear"
	StreamSID string `json:"streamSid"`
}

// Twilio outbound stream media message format (to send audio to call)
type TwilioOutboundMedia struct {
	Event     string             `json:"event"` // Must be "media"
	StreamSID string             `json:"streamSid"`
	Media     TwilioOutboundData `json:"media"`
}

type TwilioOutboundData struct {
	Payload string `json:"payload"` // base64 encoded audio
}

// OpenAI/Groq Realtime API messages structure
type RealtimeEvent struct {
	Type           string          `json:"type"`
	EventID        string          `json:"event_id,omitempty"`
	Item           *RealtimeItem   `json:"item,omitempty"`
	Delta          string          `json:"delta,omitempty"`
	Audio          string          `json:"audio,omitempty"` // For input_audio_buffer.append
	OutputItem     *RealtimeItem   `json:"output_item,omitempty"`
	FunctionCallID string          `json:"function_call_id,omitempty"`
	Response       *RealtimeResp   `json:"response,omitempty"`
	CallID         string          `json:"call_id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Arguments      string          `json:"arguments,omitempty"`
}

type RealtimeItem struct {
	Type    string          `json:"type"`
	ID      string          `json:"id"`
	Role    string          `json:"role,omitempty"`
	Name    string          `json:"name,omitempty"`
	CallID  string          `json:"call_id,omitempty"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"content,omitempty"`
}

type RealtimeResp struct {
	Status string `json:"status"`
}

// SessionUpdateEvent is sent on connection initialization
type SessionUpdateEvent struct {
	Type    string         `json:"type"`
	Session SessionContent `json:"session"`
}

type SessionContent struct {
	Modalities    []string `json:"modalities"`
	Instructions  string   `json:"instructions"`
	Voice         string   `json:"voice"`
	InputAudioFormat  string `json:"input_audio_format"`
	OutputAudioFormat string `json:"output_audio_format"`
	Tools         []Tool   `json:"tools,omitempty"`
}

type Tool struct {
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// Pool structure for high performance buffer recycling to optimize Garbage Collector pressure
var (
	// bufferPool returns pre-allocated byte slices of fixed size (8192 bytes for base64 operations)
	bufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 8192)
		},
	}
)

type StreamEngine struct {
	cfg       *config.Config
	upgrader  websocket.Upgrader
	publisher rabbitmq.Publisher
}

func NewStreamEngine(cfg *config.Config, pub rabbitmq.Publisher) *StreamEngine {
	return &StreamEngine{
		cfg: cfg,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		publisher: pub,
	}
}

// HandleTwilioStream upgrades the HTTP connection and routes bi-directional audio
func (se *StreamEngine) HandleTwilioStream(c *gin.Context) {
	prompt := c.Query("prompt")
	tenantID := c.Query("tenantId")

	// Upgrade the Gin HTTP request to Gorilla WebSocket connection
	twilioConn, err := se.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[StreamEngine] Error upgrading Twilio websocket: %v", err)
		return
	}
	defer twilioConn.Close()

	log.Printf("[StreamEngine] Twilio connected. Tenant: %s", tenantID)

	// Context to coordinate goroutine lifecycles
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Audio channel connecting Goroutine 1 (Read Twilio) -> Goroutine 2 (LLM Realtime WebSocket)
	audioIn := make(chan []byte, 100)

	// Control channel for barge-in: signals back to the outbound writer to flush queues and discard current buffers
	interruptionChan := make(chan struct{}, 5)

	// Store StreamSID and CallSID once starting event received from Twilio
	var streamSID string
	var callSID string
	var mu sync.RWMutex

	setSIDs := func(sSID, cSID string) {
		mu.Lock()
		defer mu.Unlock()
		streamSID = sSID
		callSID = cSID
	}

	getSIDs := func() (string, string) {
		mu.RLock()
		defer mu.RUnlock()
		return streamSID, callSID
	}

	// Dynamic, thread-safe audio out channel to Twilio.
	// Since write operations on standard WebSocket connections are not thread-safe, all outbound
	// writes to Twilio are handled in a dedicated goroutine.
	type outboundMessage struct {
		isClear bool
		payload string
	}
	audioOut := make(chan outboundMessage, 200)

	// GOROUTINE 1: Read voice chunks from Twilio stream, Base64 decode them and stream to audioIn
	go func() {
		defer func() {
			close(audioIn)
			cancel()
			log.Printf("[Goroutine 1] Stopped reading from Twilio")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := twilioConn.ReadMessage()
				if err != nil {
					log.Printf("[Goroutine 1] Error reading Twilio WS: %v", err)
					return
				}

				var event TwilioEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					log.Printf("[Goroutine 1] Unmarshal error: %v", err)
					continue
				}

				switch event.Event {
				case "start":
					if event.Start != nil {
						setSIDs(event.Start.StreamSID, event.Start.CallSID)
						log.Printf("[Goroutine 1] Stream started. StreamSID: %s, CallSID: %s", event.Start.StreamSID, event.Start.CallSID)
					}
				case "media":
					if event.Media != nil && event.Media.Payload != "" {
						decodedMaxLen := base64.StdEncoding.DecodedLen(len(event.Media.Payload))

						// Retrieve buffer from sync.Pool to reduce garbage collection overhead
						buf := bufferPool.Get().([]byte)

						if decodedMaxLen > len(buf) {
							// If for any reason the payload is larger than our pooled buffer, allocate dynamically to avoid panic
							buf = make([]byte, decodedMaxLen)
						}

						decodedLen, err := base64.StdEncoding.Decode(buf, []byte(event.Media.Payload))
						if err != nil {
							log.Printf("[Goroutine 1] Base64 decoding error: %v", err)
							bufferPool.Put(buf)
							continue
						}

						// Create a precise copy for sending through the channel
						chunk := make([]byte, decodedLen)
						copy(chunk, buf[:decodedLen])

						// Put buffer back into sync.Pool
						bufferPool.Put(buf)

						select {
						case audioIn <- chunk:
						case <-ctx.Done():
							return
						}
					}
				case "stop":
					log.Printf("[Goroutine 1] Call stop received")
					return
				}
			}
		}
	}()

	// GOROUTINE 2: Bi-directional pipeline with OpenAI / Groq Realtime WebSocket API
	go func() {
		defer func() {
			cancel()
			log.Printf("[Goroutine 2] Stopped Realtime API WS link")
		}()

		// Realtime API connection url
		rtURL := se.cfg.RealtimeAPIURL
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+se.cfg.RealtimeAPIKey)
		headers.Add("OpenAI-Beta", "realtime=v1")

		dialer := websocket.DefaultDialer
		rtConn, _, err := dialer.DialContext(ctx, rtURL, headers)
		if err != nil {
			log.Printf("[Goroutine 2] Failed connecting to Realtime API: %v", err)
			return
		}
		defer rtConn.Close()

		log.Println("[Goroutine 2] Connected to Realtime AI Agent.")

		// Thread-safe helper to write message to OpenAI Realtime websocket connection
		var rtWriteMu sync.Mutex
		safeWriteRtMsg := func(messageType int, payload []byte) error {
			rtWriteMu.Lock()
			defer rtWriteMu.Unlock()
			return rtConn.WriteMessage(messageType, payload)
		}

		// Configure the Realtime Session parameters (modalities, voice, system instructions)
		initSessionEvent := SessionUpdateEvent{
			Type: "session.update",
			Session: SessionContent{
				Modalities:        []string{"text", "audio"},
				Instructions:      prompt,
				Voice:             "alloy",
				InputAudioFormat:  "g711_ulaw", // Twilio streams default format
				OutputAudioFormat: "g711_ulaw", // Request G711 u-law back from LLM for direct output
				Tools: []Tool{
					{
						Type:        "function",
						Name:        "get_customer_info",
						Description: "Retrieves customer information based on their phone number.",
						Parameters: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"customerId": map[string]interface{}{
									"type":        "string",
									"description": "Unique customer identifier",
								},
							},
							"required": []string{"customerId"},
						},
					},
				},
			},
		}

		sessionPayload, err := json.Marshal(initSessionEvent)
		if err == nil {
			safeWriteRtMsg(websocket.TextMessage, sessionPayload)
		}

		// Sub-goroutine inside Goroutine 2 to read user audio chunks from channel and stream them to Realtime WebSocket
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case chunk, ok := <-audioIn:
					if !ok {
						return
					}

					base64Audio := base64.StdEncoding.EncodeToString(chunk)
					appendEvent := map[string]interface{}{
						"type":  "input_audio_buffer.append",
						"audio": base64Audio,
					}

					payload, err := json.Marshal(appendEvent)
					if err != nil {
						continue
					}

					err = safeWriteRtMsg(websocket.TextMessage, payload)
					if err != nil {
						log.Printf("[Goroutine 2 Sub] Error sending audio to Realtime API: %v", err)
						return
					}
				}
			}
		}()

		// Realtime WebSocket reader loop
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := rtConn.ReadMessage()
				if err != nil {
					log.Printf("[Goroutine 2 Reader] Error reading from Realtime API: %v", err)
					return
				}

				var rtEvent RealtimeEvent
				if err := json.Unmarshal(msg, &rtEvent); err != nil {
					continue
				}

				switch rtEvent.Type {
				case "input_audio_buffer.speech_started":
					// Barge-in Interruption detected: User started speaking while AI is playing voice.
					log.Println("[Barge-In] Speech started! Interrupting AI speech...")
					select {
					case interruptionChan <- struct{}{}:
					default:
					}

					// CRITICAL: Send response.cancel client event to the Realtime API to stop speech generation instantly
					cancelEvent := map[string]interface{}{
						"type": "response.cancel",
					}
					cancelPayload, err := json.Marshal(cancelEvent)
					if err == nil {
						_ = safeWriteRtMsg(websocket.TextMessage, cancelPayload)
					}

				case "response.audio.delta":
					if rtEvent.Delta != "" {
						// Send AI response voice chunk directly to Twilio outbound voice channel
						select {
						case audioOut <- outboundMessage{isClear: false, payload: rtEvent.Delta}:
						case <-ctx.Done():
							return
						}
					}

				case "response.function_call_arguments.done":
					// IA Requested a tool execution
					if rtEvent.Name != "" && rtEvent.Arguments != "" {
						_, cSID := getSIDs()

						// Launch an asynchronous goroutine to process the tool call publishing
						// to avoid blocking the real-time audio pipeline.
						go func(tID, callID, toolName string, args string) {
							err := se.publisher.PublishToolCall(ctx, tID, callID, toolName, json.RawMessage(args))
							if err != nil {
								log.Printf("[Tool Call] Failed to publish tool call: %v", err)
							}
						}(tenantID, cSID, rtEvent.Name, rtEvent.Arguments)
					}
				}
			}
		}
	}()

	// GOROUTINE 3: Twilio Outbound WebSocket Writer with Barge-In State management
	// Single writer ensures thread safety on Twilio WS connection.
	go func() {
		defer log.Printf("[Goroutine 3] Stopped Twilio writer")

		for {
			select {
			case <-ctx.Done():
				return
			case <-interruptionChan:
				// Flush any queued outbound audio from the channel to prevent old response chunks from playing
				cleared := false
				for !cleared {
					select {
					case <-audioOut:
					default:
						cleared = true
					}
				}

				sSID, _ := getSIDs()
				if sSID != "" {
					// Build and send 'clear' event to Twilio to stop playing audio instantly
					clearMsg := TwilioClear{
						Event:     "clear",
						StreamSID: sSID,
					}
					payload, err := json.Marshal(clearMsg)
					if err == nil {
						twilioConn.WriteMessage(websocket.TextMessage, payload)
					}
				}

			case msg, ok := <-audioOut:
				if !ok {
					return
				}

				sSID, _ := getSIDs()
				if sSID != "" {
					outboundMsg := TwilioOutboundMedia{
						Event:     "media",
						StreamSID: sSID,
						Media: TwilioOutboundData{
							Payload: msg.payload,
						},
					}

					payload, err := json.Marshal(outboundMsg)
					if err == nil {
						err = twilioConn.WriteMessage(websocket.TextMessage, payload)
						if err != nil {
							log.Printf("[Goroutine 3] Write Twilio message error: %v", err)
							return
						}
					}
				}
			}
		}
	}()

	// Wait until connection closes or is cancelled
	<-ctx.Done()
}
