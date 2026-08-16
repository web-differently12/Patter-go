package rag

type KnowledgeBaseProvider string

const (
	ProviderBuiltInPGVector KnowledgeBaseProvider = "builtin_pgvector"
	ProviderGoogleVertex    KnowledgeBaseProvider = "google_vertex_search"
	ProviderPinecone        KnowledgeBaseProvider = "pinecone"
	ProviderQdrant          KnowledgeBaseProvider = "qdrant"
	ProviderCustomEndpoint  KnowledgeBaseProvider = "custom_endpoint"
)

// FallbackTransferPolicy définit la règle de transfert automatique vers un humain
type FallbackTransferPolicy struct {
	EnableAutoTransfer  bool    `json:"enable_auto_transfer"`  // Activer l'escalade vers un humain
	ConfidenceThreshold float32 `json:"confidence_threshold"`  // Ex: 0.75 (en dessous -> transfert)
	TransferDestination string  `json:"transfer_destination"`  // Numéro E.164 ou SIP URI (+33612345678)
	TransferType        string  `json:"transfer_type"`         // "WARM" (accompagné) ou "COLD" (direct)
	TransferMessage     string  `json:"transfer_message"`      // Phrase prononcée avant de transférer
	HoldMusicURL        string  `json:"hold_music_url"`        // Musique d'attente
}

// KnowledgeBaseConfig regroupe tous les providers RAG et la politique d'escalade
type KnowledgeBaseConfig struct {
	Provider KnowledgeBaseProvider `json:"provider"`

	// 1. Google ADK / Vertex Search Config
	GoogleVertexConfig struct {
		ProjectID   string `json:"project_id,omitempty"`
		Location    string `json:"location,omitempty"`    // "global" ou "europe-west1"
		DataStoreID string `json:"datastore_id,omitempty"`
	} `json:"google_vertex_config,omitempty"`

	// 2. Custom Endpoint Config
	CustomServerURL string `json:"custom_server_url,omitempty"`
	AuthHeaderKey   string `json:"auth_header_key,omitempty"`
	AuthHeaderValue string `json:"auth_header_value,omitempty"`

	// 3. Règle d'Escalade & Auto-Transfert d'Appel
	FallbackTransfer FallbackTransferPolicy `json:"fallback_transfer"`

	TopK int `json:"top_k"`
}
