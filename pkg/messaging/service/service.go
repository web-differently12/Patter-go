package service

import (
	"context"

	"github.com/google/uuid"
)

type SendSMSRequest struct {
	FromNumber string `json:"from_number" binding:"required"`
	ToNumber   string `json:"to_number" binding:"required"`
	Message    string `json:"message" binding:"required"`
	MediaURL   string `json:"media_url,omitempty"`
}

type SendSMSResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type MessagingService interface {
	SendSMS(ctx context.Context, tenantID string, req SendSMSRequest) (*SendSMSResponse, error)
}

type messagingService struct{}

func NewMessagingService() MessagingService {
	return &messagingService{}
}

func (s *messagingService) SendSMS(ctx context.Context, tenantID string, req SendSMSRequest) (*SendSMSResponse, error) {
	msgID := "sms_" + uuid.New().String()[:8]
	return &SendSMSResponse{
		MessageID: msgID,
		Status:    "SENT",
	}, nil
}
