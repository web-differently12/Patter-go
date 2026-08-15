package dto

type ConnectWhatsAppRequest struct {
	SessionName string `json:"session_name" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type QRCodeResponse struct {
	SessionName string `json:"session_name"`
	QRCode      string `json:"qrcode_base64"`
	PairingCode string `json:"pairing_code,omitempty"`
	Status      string `json:"status"` // "SCAN_QR_CODE", "WORKING", "CONNECTING"
	ExpiresIn   int    `json:"expires_in"`
}

type SessionStatusResponse struct {
	SessionName string `json:"session_name"`
	TenantID    string `json:"tenant_id"`
	State       string `json:"state"` // "WORKING", "DISCONNECTED", "CONNECTING", "PAIRING"
	Connected   bool   `json:"connected"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type SendMessageRequest struct {
	SessionName string `json:"session_name" binding:"required"`
	Recipient   string `json:"recipient" binding:"required"` // E.164 format
	Message     string `json:"message" binding:"required"`
}

type SendMediaRequest struct {
	SessionName string `json:"session_name" binding:"required"`
	Recipient   string `json:"recipient" binding:"required"`
	MediaType   string `json:"media_type" binding:"required"` // "image", "video", "document", "audio", "sticker"
	MediaURL    string `json:"media_url" binding:"required"`
	Caption     string `json:"caption,omitempty"`
}

type SendLocationRequest struct {
	SessionName string  `json:"session_name" binding:"required"`
	Recipient   string  `json:"recipient" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Name        string  `json:"name,omitempty"`
	Address     string  `json:"address,omitempty"`
}

type SendContactRequest struct {
	SessionName string `json:"session_name" binding:"required"`
	Recipient   string `json:"recipient" binding:"required"`
	ContactName string `json:"contact_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
}

type SendMessageResponse struct {
	MessageID   string `json:"message_id"`
	SessionName string `json:"session_name"`
	Recipient   string `json:"recipient"`
	Status      string `json:"status"` // "SENT", "PENDING", "FAILED"
}
