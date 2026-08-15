package service

import (
	"context"
	"fmt"
	"strings"
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
	SendMedia(ctx context.Context, tenantID string, req dto.SendMediaRequest) (*dto.SendMessageResponse, error)
	SendLocation(ctx context.Context, tenantID string, req dto.SendLocationRequest) (*dto.SendMessageResponse, error)
	SendContact(ctx context.Context, tenantID string, req dto.SendContactRequest) (*dto.SendMessageResponse, error)
	DisconnectSession(ctx context.Context, tenantID, sessionName string) error
	CheckNumberExists(ctx context.Context, tenantID string, req dto.CheckNumberRequest) (*dto.CheckNumberResponse, error)
}

type whatsappService struct {
	mu       sync.RWMutex
	sessions map[string]*dto.SessionStatusResponse
}

func NewWhatsAppService() WhatsAppService {
	svc := &whatsappService{
		sessions: make(map[string]*dto.SessionStatusResponse),
	}
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

func (s *whatsappService) CheckNumberExists(ctx context.Context, tenantID string, req dto.CheckNumberRequest) (*dto.CheckNumberResponse, error) {
	cleanDigits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, req.PhoneNumber)

	jid := cleanDigits + "@s.whatsapp.net"
	exists := len(cleanDigits) >= 10 && len(cleanDigits) <= 15

	return &dto.CheckNumberResponse{
		PhoneNumber:  req.PhoneNumber,
		JID:          jid,
		Exists:       exists,
		IsInWhatsApp: exists,
	}, nil
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
			State:       "WORKING",
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

	msgID := "msg_wa_" + uuid.New().String()[:8]
	return &dto.SendMessageResponse{
		MessageID:   msgID,
		SessionName: req.SessionName,
		Recipient:   req.Recipient,
		Status:      "SENT",
	}, nil
}

func (s *whatsappService) SendMedia(ctx context.Context, tenantID string, req dto.SendMediaRequest) (*dto.SendMessageResponse, error) {
	msgID := "msg_media_" + uuid.New().String()[:8]
	return &dto.SendMessageResponse{
		MessageID:   msgID,
		SessionName: req.SessionName,
		Recipient:   req.Recipient,
		Status:      "SENT",
	}, nil
}

func (s *whatsappService) SendLocation(ctx context.Context, tenantID string, req dto.SendLocationRequest) (*dto.SendMessageResponse, error) {
	msgID := "msg_loc_" + uuid.New().String()[:8]
	return &dto.SendMessageResponse{
		MessageID:   msgID,
		SessionName: req.SessionName,
		Recipient:   req.Recipient,
		Status:      "SENT",
	}, nil
}

func (s *whatsappService) SendContact(ctx context.Context, tenantID string, req dto.SendContactRequest) (*dto.SendMessageResponse, error) {
	msgID := "msg_cnt_" + uuid.New().String()[:8]
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
