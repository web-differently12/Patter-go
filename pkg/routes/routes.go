package routes

import (
	"vocal-engine/pkg/api/admin"
	"vocal-engine/pkg/api/gateway"
	"vocal-engine/pkg/api/middleware"
	"vocal-engine/pkg/engine"
	"vocal-engine/pkg/telephony"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter registers endpoints for telephony, agents, admin BYOK configuration, and white-label gateway APIs
func SetupRouter(tc *telephony.CallController, ac *telephony.AgentController, se *engine.StreamEngine) *gin.Engine {
	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve the static frontend web dashboard files
	r.StaticFile("/", "./web/index.html")
	r.Static("/web", "./web")

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Instantiate Controllers
	adminCtrl := admin.NewAdminController()
	brainCtrl := gateway.NewBrainController()
	voiceCtrl := gateway.NewVoiceGatewayController()
	waCtrl := gateway.NewWhatsAppGatewayController()
	msgCtrl := gateway.NewMessagingGatewayController()
	avatarCtrl := gateway.NewAvatarGatewayController()
	meetingCtrl := gateway.NewMeetingGatewayController()

	// 1. ADMIN SPACE (`/api/v1/admin/...`)
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.AdminAuthMiddleware())
	{
		adminGroup.GET("/tenants/:tenant_id/config", adminCtrl.GetTenantConfig)
		adminGroup.PUT("/tenants/:tenant_id/config", adminCtrl.UpdateTenantConfig)
		adminGroup.GET("/tenants/:tenant_id/branding", adminCtrl.GetTenantBranding)
	}

	// 2. WHITE-LABEL GATEWAY SPACE (`/api/v1/gateway/...`)
	gatewayGroup := r.Group("/api/v1/gateway")
	gatewayGroup.Use(middleware.TenantContextMiddleware())
	{
		// Brain
		gatewayGroup.POST("/brain/chat", brainCtrl.Chat)
		gatewayGroup.POST("/brain/report", brainCtrl.Report)

		// Voice
		gatewayGroup.POST("/voice/outbound", voiceCtrl.Outbound)
		gatewayGroup.POST("/voice/campaign", voiceCtrl.Campaign)
		gatewayGroup.POST("/voice/inbound/webhook", voiceCtrl.InboundWebhook)

		// WhatsApp
		gatewayGroup.POST("/whatsapp/send", waCtrl.Send)
		gatewayGroup.POST("/whatsapp/campaign", waCtrl.Campaign)
		gatewayGroup.GET("/whatsapp/templates", waCtrl.GetTemplates)
		gatewayGroup.POST("/whatsapp/webhook", waCtrl.Webhook)

		// Messaging SMS/MMS
		gatewayGroup.POST("/messaging/send", msgCtrl.Send)
		gatewayGroup.POST("/messaging/campaign", msgCtrl.Campaign)
		gatewayGroup.POST("/messaging/webhook", msgCtrl.Webhook)

		// Avatar & Video
		gatewayGroup.POST("/avatar/live", avatarCtrl.Live)
		gatewayGroup.POST("/avatar/offline", avatarCtrl.Offline)

		// Meeting
		gatewayGroup.POST("/meeting/schedule", meetingCtrl.Schedule)
		gatewayGroup.POST("/meeting/bot", meetingCtrl.Bot)
	}

	// 3. CORE TELEPHONY & AGENT APIS (`/api/v1/...`)
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/call/outbound", tc.HandleOutboundCall)
		apiV1.POST("/call/inbound", ac.HandleInboundCall)
		apiV1.GET("/call/twiml", tc.HandleTwiML)

		apiV1.POST("/agents", ac.HandleCreateAgent)
		apiV1.GET("/agents", ac.HandleListAgents)
	}

	// Bi-directional stream endpoint
	r.GET("/ws/twilio/stream", se.HandleTwilioStream)

	return r
}
