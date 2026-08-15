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

type RCSButton struct {
	Title   string `json:"title"`
	Type    string `json:"type"` // "URL", "CALL", "REPLY"
	Payload string `json:"payload"`
}

type SendRCSRequest struct {
	FromNumber  string      `json:"from_number" binding:"required"`
	ToNumber    string      `json:"to_number" binding:"required"`
	Title       string      `json:"title,omitempty"`
	Message     string      `json:"message" binding:"required"`
	MediaURL    string      `json:"media_url,omitempty"`
	Buttons     []RCSButton `json:"buttons,omitempty"`
	FallbackSMS bool        `json:"fallback_sms"`
}

type SendRCSResponse struct {
	MessageID   string `json:"message_id"`
	ChannelUsed string `json:"channel_used"` // "RCS" or "SMS_FALLBACK"
	Status      string `json:"status"`
	RCSCapable  bool   `json:"rcs_capable"`
}

type MessagingService interface {
	SendSMS(ctx context.Context, tenantID string, req SendSMSRequest) (*SendSMSResponse, error)
	SendRCS(ctx context.Context, tenantID string, req SendRCSRequest) (*SendRCSResponse, error)
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

func (s *messagingService) SendRCS(ctx context.Context, tenantID string, req SendRCSRequest) (*SendRCSResponse, error) {
	msgID := "rcs_" + uuid.New().String()[:8]
	return &SendRCSResponse{
		MessageID:   msgID,
		ChannelUsed: "RCS",
		Status:      "DELIVERED",
		RCSCapable:  true,
	}, nil
}
