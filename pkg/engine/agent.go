package engine

import (
	"errors"
)

// Mode defines the voice execution pattern
type Mode string

const (
	ModeRealtime Mode = "realtime" // Bidirectional WebSockets OpenAI/Groq
	ModePipeline Mode = "pipeline" // Orchestrated (STT -> LLM -> TTS)
)

// Modality defines the modes of communication
type Modality string

const (
	ModalityText  Modality = "text"
	ModalityAudio Modality = "audio"
)

// Agent represents the structural AI agent configuration mirroring Patter TS features
type Agent struct {
	ID             string          `json:"id" example:"agent-99"`
	Name           string          `json:"name" example:"Acme Receptionist"`
	TenantID       string          `json:"tenantId" example:"tenant-123"`
	Mode           Mode            `json:"mode" example:"realtime"`
	SystemPrompt   string          `json:"systemPrompt" example:"You are a helpful office receptionist."`
	FirstMessage   string          `json:"firstMessage" example:"Hello! Thank you for calling Acme. How can I help you?"`
	Voice          string          `json:"voice" example:"alloy"`
	STTProviderID  string          `json:"sttProviderId" example:"deepgram"`
	TTSProviderID  string          `json:"ttsProviderId" example:"elevenlabs"`
	LLMProviderID  string          `json:"llmProviderId" example:"openai"`
	FallbackChain  []string        `json:"fallbackChain" example:"[\"groq\", \"anthropic\"]"`
	Tools          []AgentTool     `json:"tools"`
	MaxDurationSec int             `json:"maxDurationSec" example:"600"`
	VADSensitivity float64         `json:"vadSensitivity" example:"0.5"`
}

// AgentTool represents a function declaration schema for Tool Calling
type AgentTool struct {
	Name        string                 `json:"name" example:"get_customer_info"`
	Description string                 `json:"description" example:"Get details for a customer"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// AgentRepository defines an interface for managing agents metadata persistence
type AgentRepository interface {
	GetAgentByID(ctx interface{}, id string) (*Agent, error)
	CreateAgent(ctx interface{}, agent *Agent) error
	ListAgents(ctx interface{}) ([]*Agent, error)
}

// MemoryAgentRepository implements a simple in-memory repository for testing and fallback use
type MemoryAgentRepository struct {
	agents map[string]*Agent
}

func NewMemoryAgentRepository() *MemoryAgentRepository {
	return &MemoryAgentRepository{
		agents: make(map[string]*Agent),
	}
}

func (m *MemoryAgentRepository) GetAgentByID(ctx interface{}, id string) (*Agent, error) {
	agent, ok := m.agents[id]
	if !ok {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

func (m *MemoryAgentRepository) CreateAgent(ctx interface{}, agent *Agent) error {
	m.agents[agent.ID] = agent
	return nil
}

func (m *MemoryAgentRepository) ListAgents(ctx interface{}) ([]*Agent, error) {
	var list []*Agent
	for _, a := range m.agents {
		list = append(list, a)
	}
	return list, nil
}
