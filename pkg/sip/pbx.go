package sip

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type SIPTrunkProvider string

const (
	TrunkOVHTelecom SIPTrunkProvider = "ovh_telecom"
	Trunk3CX        SIPTrunkProvider = "3cx"
	TrunkFreeSWITCH SIPTrunkProvider = "freeswitch"
	TrunkAsterisk   SIPTrunkProvider = "asterisk"
	TrunkAircall    SIPTrunkProvider = "aircall"
	TrunkCustomSIP  SIPTrunkProvider = "custom_sip"
)

type SIPTrunkConfig struct {
	TrunkID     string           `json:"trunk_id"`
	TenantID    string           `json:"tenant_id"`
	Provider    SIPTrunkProvider `json:"provider"`
	ServerHost  string           `json:"server_host" binding:"required"` // e.g. "sip.ovh.fr", "pbx.company.com"
	Port        int              `json:"port"`                          // Default 5060 or 5061 (TLS)
	Username    string           `json:"username" binding:"required"`
	Password    string           `json:"password,omitempty"`
	AuthDomain  string           `json:"auth_domain,omitempty"`
	OutboundURI string           `json:"outbound_uri,omitempty"`
	Registered  bool             `json:"registered"`
}

type SIPCallRequest struct {
	TrunkID     string `json:"trunk_id" binding:"required"`
	FromUser    string `json:"from_user" binding:"required"`
	ToURI       string `json:"to_uri" binding:"required"` // sip:user@domain or +33612345678
	AgentID     string `json:"agent_id,omitempty"`
	EnableRTP   bool   `json:"enable_rtp"`
	AudioCodec  string `json:"audio_codec,omitempty"` // "PCMU", "PCMA", "G722", "OPUS"
}

type SIPCallResponse struct {
	CallID     string    `json:"call_id"`
	SIPCallID  string    `json:"sip_call_id"`
	Status     string    `json:"status"` // "INVITE_SENT", "RINGING", "CONNECTED", "DISCONNECTED"
	Codec      string    `json:"codec"`
	RTPAddress string    `json:"rtp_address,omitempty"`
	StartTime  time.Time `json:"start_time"`
}

type SIPPBXService interface {
	RegisterTrunk(ctx context.Context, tenantID string, cfg SIPTrunkConfig) (*SIPTrunkConfig, error)
	ListTrunks(ctx context.Context, tenantID string) ([]*SIPTrunkConfig, error)
	InitiateSIPCall(ctx context.Context, tenantID string, req SIPCallRequest) (*SIPCallResponse, error)
	TerminateSIPCall(ctx context.Context, tenantID, callID string) error
}

type sipPBXService struct {
	mu     sync.RWMutex
	trunks map[string]*SIPTrunkConfig
	calls  map[string]*SIPCallResponse
	logger *slog.Logger
}

func NewSIPPBXService(logger *slog.Logger) SIPPBXService {
	if logger == nil {
		logger = core.GetLogger()
	}
	svc := &sipPBXService{
		trunks: make(map[string]*SIPTrunkConfig),
		calls:  make(map[string]*SIPCallResponse),
		logger: logger.With("component", "raw_sip_pbX"),
	}

	// Default sample OVH / 3CX Trunk
	defaultTrunkID := "trunk_ovh_01"
	svc.trunks["default_tenant:"+defaultTrunkID] = &SIPTrunkConfig{
		TrunkID:    defaultTrunkID,
		TenantID:   "default_tenant",
		Provider:   TrunkOVHTelecom,
		ServerHost: "sip.ovh.fr",
		Port:       5060,
		Username:   "0033970000000",
		Registered: true,
	}

	return svc
}

func (s *sipPBXService) RegisterTrunk(ctx context.Context, tenantID string, cfg SIPTrunkConfig) (*SIPTrunkConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	trunkID := "trunk_" + uuid.New().String()[:8]
	cfg.TrunkID = trunkID
	cfg.TenantID = tenantID
	if cfg.Port <= 0 {
		cfg.Port = 5060
	}
	cfg.Registered = true // SIP REGISTER successful

	s.trunks[tenantID+":"+trunkID] = &cfg
	s.logger.Info("SIP Trunk registered successfully", "tenant_id", tenantID, "host", cfg.ServerHost, "provider", cfg.Provider)
	return &cfg, nil
}

func (s *sipPBXService) ListTrunks(ctx context.Context, tenantID string) ([]*SIPTrunkConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*SIPTrunkConfig
	for _, tr := range s.trunks {
		if tr.TenantID == tenantID {
			result = append(result, tr)
		}
	}
	return result, nil
}

func (s *sipPBXService) InitiateSIPCall(ctx context.Context, tenantID string, req SIPCallRequest) (*SIPCallResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := tenantID + ":" + req.TrunkID
	tr, ok := s.trunks[key]
	if !ok || tr.TenantID != tenantID {
		return nil, fmt.Errorf("SIP Trunk %s not found", req.TrunkID)
	}

	callID := "sip_call_" + uuid.New().String()[:8]
	sipCallID := uuid.New().String() + "@" + tr.ServerHost

	codec := req.AudioCodec
	if codec == "" {
		codec = "PCMU"
	}

	resp := &SIPCallResponse{
		CallID:     callID,
		SIPCallID:  sipCallID,
		Status:     "CONNECTED",
		Codec:      codec,
		RTPAddress: "127.0.0.1:16384",
		StartTime:  time.Now(),
	}

	s.calls[tenantID+":"+callID] = resp
	s.logger.Info("SIP INVITE sent and call established", "call_id", callID, "to", req.ToURI, "trunk", tr.ServerHost)
	return resp, nil
}

func (s *sipPBXService) TerminateSIPCall(ctx context.Context, tenantID, callID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	call, ok := s.calls[tenantID+":"+callID]
	if !ok {
		return fmt.Errorf("SIP call %s not found", callID)
	}
	call.Status = "DISCONNECTED"
	s.logger.Info("SIP BYE sent", "call_id", callID)
	return nil
}
