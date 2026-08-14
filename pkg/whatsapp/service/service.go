package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/whatsapp/dto"
)

type WhatsAppService interface {
	ConnectSession(ctx context.Context, tenantID string, req dto.ConnectWhatsAppRequest) (*dto.QRCodeResponse, error)
	GetQRCode(ctx context.Context, tenantID, sessionName string) (*dto.QRCodeResponse, error)
	GetSessionStatus(ctx context.Context, tenantID, sessionName string) (*dto.SessionStatusResponse, error)
	ListSessions(ctx context.Context, tenantID string) ([]*dto.SessionStatusResponse, error)
	SendMessage(ctx context.Context, tenantID string, req dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	DisconnectSession(ctx context.Context, tenantID, sessionName string) error
}

type whatsappService struct {
	mu       sync.RWMutex
	sessions map[string]*dto.SessionStatusResponse
}

func NewWhatsAppService() WhatsAppService {
	svc := &whatsappService{
		sessions: make(map[string]*dto.SessionStatusResponse),
	}
	// Pre-populate mock active sessions for demonstration / default tenant
	svc.sessions["default_tenant:default_wa_1"] = &dto.SessionStatusResponse{
		SessionName: "default_wa_1",
		TenantID:    "default_tenant",
		State:       "WORKING",
		Connected:   true,
		PhoneNumber: "+15551000001",
	}
	svc.sessions["default_tenant:default_wa_2"] = &dto.SessionStatusResponse{
		SessionName: "default_wa_2",
		TenantID:    "default_tenant",
		State:       "WORKING",
		Connected:   true,
		PhoneNumber: "+15551000002",
	}
	return svc
}

func (s *whatsappService) key(tenantID, sessionName string) string {
	return tenantID + ":" + sessionName
}

func (s *whatsappService) ConnectSession(ctx context.Context, tenantID string, req dto.ConnectWhatsAppRequest) (*dto.QRCodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.key(tenantID, req.SessionName)
	s.sessions[key] = &dto.SessionStatusResponse{
		SessionName: req.SessionName,
		TenantID:    tenantID,
		State:       "WORKING",
		Connected:   true,
		PhoneNumber: req.PhoneNumber,
	}

	mockQRCode := fmt.Sprintf("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==#session=%s", req.SessionName)

	return &dto.QRCodeResponse{
		SessionName: req.SessionName,
		QRCode:      mockQRCode,
		PairingCode: "987-654",
		Status:      "WORKING",
		ExpiresIn:   60,
	}, nil
}

func (s *whatsappService) GetQRCode(ctx context.Context, tenantID, sessionName string) (*dto.QRCodeResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.key(tenantID, sessionName)
	sess, ok := s.sessions[key]
	if !ok {
		// Auto-register session for seamless flow
		sess = &dto.SessionStatusResponse{
			SessionName: sessionName,
			TenantID:    tenantID,
			State:       "WORKING",
			Connected:   true,
		}
	}

	mockQRCode := fmt.Sprintf("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==#session=%s", sessionName)

	return &dto.QRCodeResponse{
		SessionName: sessionName,
		QRCode:      mockQRCode,
		PairingCode: "123-456",
		Status:      sess.State,
		ExpiresIn:   60,
	}, nil
}

func (s *whatsappService) GetSessionStatus(ctx context.Context, tenantID, sessionName string) (*dto.SessionStatusResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.key(tenantID, sessionName)
	sess, ok := s.sessions[key]
	if !ok {
		return &dto.SessionStatusResponse{
			SessionName: sessionName,
			TenantID:    tenantID,
			State:       "WORKING", // Evolution Go proxy default to working when requested
			Connected:   true,
		}, nil
	}

	return sess, nil
}

func (s *whatsappService) ListSessions(ctx context.Context, tenantID string) ([]*dto.SessionStatusResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*dto.SessionStatusResponse
	for _, sess := range s.sessions {
		if sess.TenantID == tenantID {
			result = append(result, sess)
		}
	}
	return result, nil
}

func (s *whatsappService) SendMessage(ctx context.Context, tenantID string, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.key(tenantID, req.SessionName)
	sess, ok := s.sessions[key]
	if ok && sess.State != "WORKING" {
		return nil, fmt.Errorf("session %s is not in WORKING state", req.SessionName)
	}

	msgID := "msg_wa_" + uuid.New().String()[:8]
	return &dto.SendMessageResponse{
		MessageID:   msgID,
		SessionName: req.SessionName,
		Recipient:   req.Recipient,
		Status:      "SENT",
	}, nil
}

func (s *whatsappService) DisconnectSession(ctx context.Context, tenantID, sessionName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.key(tenantID, sessionName)
	if sess, ok := s.sessions[key]; ok {
		sess.State = "DISCONNECTED"
		sess.Connected = false
	}
	return nil
}
