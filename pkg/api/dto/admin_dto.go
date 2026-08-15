package dto

type TenantConfigDTO struct {
	TenantID             string            `json:"tenant_id" example:"tenant-123"`
	TwilioAccountSID     string            `json:"twilio_account_sid,omitempty" example:"AC12345"`
	TwilioAuthToken      string            `json:"twilio_auth_token,omitempty" example:"secret"`
	TelnyxAPIKey         string            `json:"telnyx_api_key,omitempty" example:"KEY123"`
	SimliAPIKey          string            `json:"simli_api_key,omitempty" example:"SIMLI123"`
	OpenRouterAPIKey     string            `json:"openrouter_api_key,omitempty" example:"sk-or-123"`
	ElevenLabsAPIKey     string            `json:"elevenlabs_api_key,omitempty" example:"eleven123"`
	AssemblyAIAPIKey     string            `json:"assemblyai_api_key,omitempty" example:"assembly123"`
	LiveKitHost          string            `json:"livekit_host,omitempty" example:"wss://livekit.example.com"`
	LiveKitAPIKey        string            `json:"livekit_api_key,omitempty" example:"lk_key"`
	LiveKitAPISecret     string            `json:"livekit_api_secret,omitempty" example:"lk_secret"`
	EvolutionGoServerURL string            `json:"evolution_go_server_url,omitempty" example:"http://evolution-go:8080"`
	Webhooks             map[string]string `json:"webhooks,omitempty"`
}

type TenantBrandingDTO struct {
	TenantID      string `json:"tenant_id" example:"tenant-123"`
	CustomDomain  string `json:"custom_domain" example:"ai.acme.com"`
	DisplayName   string `json:"display_name" example:"Acme AI Assistant"`
	LogoURL       string `json:"logo_url" example:"https://acme.com/logo.png"`
	DefaultVoice  string `json:"default_voice" example:"alloy"`
	PrimaryColor  string `json:"primary_color" example:"#3B82F6"`
}
