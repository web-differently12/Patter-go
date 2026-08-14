package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "vocal-engine/docs"
	"vocal-engine/pkg/config"
	"vocal-engine/pkg/engine"
	"vocal-engine/pkg/rabbitmq"
	"vocal-engine/pkg/routes"
	"vocal-engine/pkg/telephony"
)

// @title Vocal Engine AI Agent Microservice API
// @version 1.0
// @description Microservice engine handling low latency bidirectionnal Twilio-to-LLM audio routing and asynchronous tool calling.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	log.Println("[Main] Starting Vocal Engine Microservice...")

	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize RabbitMQ Publisher
	var publisher rabbitmq.Publisher
	var err error
	if cfg.RabbitMQURL != "" {
		log.Printf("[Main] Attempting to connect to RabbitMQ at: %s", cfg.RabbitMQURL)
		publisher, err = rabbitmq.NewRabbitMQPublisher(cfg)
		if err != nil {
			log.Printf("[Main] [WARNING] Failed to connect to RabbitMQ: %v. Running in mockup fallback mode.", err)
			publisher = &mockPublisher{}
		} else {
			log.Println("[Main] Connected to RabbitMQ successfully.")
			defer publisher.Close()
		}
	} else {
		log.Println("[Main] RabbitMQ URL not set. Initializing mockup fallback publisher.")
		publisher = &mockPublisher{}
	}

	// 3. Initialize Shared State Repositories
	agentRepo := engine.NewMemoryAgentRepository()

	// Seed dummy agent for testing
	dummyAgent := &engine.Agent{
		ID:             "agent-receptionist",
		Name:           "Acme reception voice agent",
		TenantID:       "tenant-1",
		Mode:           engine.ModeRealtime,
		SystemPrompt:   "You are a welcoming office secretary receptionist at Acme Corporation.",
		FirstMessage:   "Hello! Thank you for calling Acme. How can I guide you today?",
		Voice:          "alloy",
		STTProviderID:  "deepgram",
		TTSProviderID:  "elevenlabs",
		LLMProviderID:  "openai",
		FallbackChain:  []string{"groq", "anthropic"},
		VADSensitivity: 0.25,
		Tools:          []engine.AgentTool{},
	}
	_ = agentRepo.CreateAgent(context.Background(), dummyAgent)

	// 4. Initialize Domain Handlers and Controllers
	callController := telephony.NewCallController(cfg)
	agentController := telephony.NewAgentController(agentRepo)
	streamEngine := engine.NewStreamEngine(cfg, publisher)

	// 5. Initialize Router with newly added controllers
	router := routes.SetupRouter(callController, agentController, streamEngine)

	// 6. Configure Graceful HTTP Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Run server in a goroutine so that it doesn't block
	go func() {
		log.Printf("[Main] Listening and serving HTTP on port: %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Main] Listen error: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[Main] Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[Main] Server forced to shutdown: ", err)
	}

	log.Println("[Main] Microservice exited successfully.")
}

// mockPublisher is used if RabbitMQ is not available or disabled (fallback development mode)
type mockPublisher struct{}

func (m *mockPublisher) PublishToolCall(ctx context.Context, tenantID, callSID, toolName string, arguments json.RawMessage) error {
	log.Printf("[MockPublisher] Intercepted Tool Call -> Tool: %s, tenant: %s, callSid: %s. Args: %s", toolName, tenantID, callSID, string(arguments))
	return nil
}

func (m *mockPublisher) Close() error {
	return nil
}
