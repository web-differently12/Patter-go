package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// GenerateTypeScriptInterfaces returns Contract-First TypeScript types generated directly from Go DTOs
func GenerateTypeScriptInterfaces() string {
	return `// Contract-First Auto-Generated TypeScript Interfaces for Patter Engine Gateway
export interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

export interface ConnectWhatsAppRequest {
  session_name: string;
  phone_number?: string;
}

export interface CheckNumberRequest {
  phone_number: string;
}

export interface CheckNumberResponse {
  phone_number: string;
  jid: string;
  exists: boolean;
  is_in_whatsapp: boolean;
}

export interface SendRCSRequest {
  from_number: string;
  to_number: string;
  title?: string;
  message: string;
  media_url?: string;
  buttons?: Array<{ title: string; type: 'URL' | 'CALL' | 'REPLY'; payload: string }>;
  fallback_sms: boolean;
}

export interface MCPServerConfig {
  server_id: string;
  tenant_id: string;
  name: string;
  url: string;
  transport: 'streamable_http' | 'sse' | 'websocket' | 'nango_unified_bridge';
  status: 'CONNECTED' | 'DISCONNECTED';
}

export interface CreateMeetingBotRequest {
  meeting_url: string;
  bot_name?: string;
  avatar_url?: string;
  platform?: 'zoom' | 'google_meet' | 'teams' | 'webex';
  recording_mode?: 'speaker_view' | 'gallery_view' | 'audio_only';
  enable_realtime_transcript: boolean;
  enable_realtime_audio_stream: boolean;
  language?: string;
}
`
}

func ServeTypeScriptSchema(c *gin.Context) {
	c.Header("Content-Type", "application/typescript")
	c.String(http.StatusOK, GenerateTypeScriptInterfaces())
}
