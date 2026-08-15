package skills

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type SendValidationCodeRequest struct {
	TenantID   string `json:"tenant_id"`
	Recipient  string `json:"recipient"`   // E.164 phone number
	Channel    string `json:"channel"`     // "WHATSAPP", "SMS", "RCS", "AUTO"
	CodeLength int    `json:"code_length"` // e.g. 6 digits
	CustomCode string `json:"custom_code,omitempty"`
}

type SendValidationCodeResult struct {
	Success      bool   `json:"success"`
	CodeSent     string `json:"code_sent"`
	ChannelUsed  string `json:"channel_used"`
	MessageID    string `json:"message_id"`
	Recipient    string `json:"recipient"`
	ExpiresInSec int    `json:"expires_in_sec"`
}

type SendOutboundMessageRequest struct {
	TenantID  string `json:"tenant_id"`
	Recipient string `json:"recipient"`
	Channel   string `json:"channel"` // "WHATSAPP", "SMS", "RCS"
	Content   string `json:"content"`
	MediaURL  string `json:"media_url,omitempty"`
}

type SendOutboundMessageResult struct {
	Success     bool   `json:"success"`
	ChannelUsed string `json:"channel_used"`
	MessageID   string `json:"message_id"`
	Recipient   string `json:"recipient"`
	Status      string `json:"status"`
}

type MessageSkillEngine struct{}

func NewMessageSkillEngine() *MessageSkillEngine {
	return &MessageSkillEngine{}
}

func (e *MessageSkillEngine) SendValidationCode(ctx context.Context, req SendValidationCodeRequest) (*SendValidationCodeResult, error) {
	code := req.CustomCode
	if code == "" {
		length := req.CodeLength
		if length <= 0 {
			length = 6
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		digits := ""
		for i := 0; i < length; i++ {
			digits += fmt.Sprintf("%d", r.Intn(10))
		}
		code = digits
	}

	channel := req.Channel
	if channel == "" || channel == "AUTO" {
		channel = "WHATSAPP"
	}

	msgID := "val_code_" + uuid.New().String()[:8]

	return &SendValidationCodeResult{
		Success:      true,
		CodeSent:     code,
		ChannelUsed:  channel,
		MessageID:    msgID,
		Recipient:    req.Recipient,
		ExpiresInSec: 300,
	}, nil
}

func (e *MessageSkillEngine) SendOutboundMessage(ctx context.Context, req SendOutboundMessageRequest) (*SendOutboundMessageResult, error) {
	channel := req.Channel
	if channel == "" {
		channel = "WHATSAPP"
	}

	msgID := "ai_dispatch_" + uuid.New().String()[:8]

	return &SendOutboundMessageResult{
		Success:     true,
		ChannelUsed: channel,
		MessageID:   msgID,
		Recipient:   req.Recipient,
		Status:      "SENT",
	}, nil
}
