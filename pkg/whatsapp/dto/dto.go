package dto

type ConnectWhatsAppRequest struct {
	SessionName string `json:"session_name" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type QRCodeResponse struct {
	SessionName string `json:"session_name"`
	QRCode      string `json:"qrcode_base64"` // Base64 encoded QR Code or string SVG
	PairingCode string `json:"pairing_code,omitempty"`
	Status      string `json:"status"` // "SCAN_QR_CODE", "WORKING", "CONNECTING"
	ExpiresIn   int    `json:"expires_in"` // in seconds
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
	Recipient   string `json:"recipient" binding:"required"` // Phone number in E.164
	Message     string `json:"message" binding:"required"`
	MediaURL    string `json:"media_url,omitempty"`
}

type SendMessageResponse struct {
	MessageID   string `json:"message_id"`
	SessionName string `json:"session_name"`
	Recipient   string `json:"recipient"`
	Status      string `json:"status"` // "SENT", "PENDING", "FAILED"
}
