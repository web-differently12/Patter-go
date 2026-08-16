package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lynxflow/patter-go/docs"
	avatarCtrl "github.com/lynxflow/patter-go/pkg/avatar"
	brainCtrl "github.com/lynxflow/patter-go/pkg/brain/controller"
	brainSvc "github.com/lynxflow/patter-go/pkg/brain/service"
	campaignCtrl "github.com/lynxflow/patter-go/pkg/campaign/controller"
	campaignSvc "github.com/lynxflow/patter-go/pkg/campaign/service"
	"github.com/lynxflow/patter-go/pkg/core"
	instanceCtrl "github.com/lynxflow/patter-go/pkg/instance/controller"
	instanceSvc "github.com/lynxflow/patter-go/pkg/instance/service"
	mcpCtrl "github.com/lynxflow/patter-go/pkg/mcp"
	messagingCtrl "github.com/lynxflow/patter-go/pkg/messaging/controller"
	messagingSvc "github.com/lynxflow/patter-go/pkg/messaging/service"
	"github.com/lynxflow/patter-go/pkg/rtc"
	"github.com/lynxflow/patter-go/pkg/sip"
	voiceCtrl "github.com/lynxflow/patter-go/pkg/voice/controller"
	voiceSvc "github.com/lynxflow/patter-go/pkg/voice/service"
	waCtrl "github.com/lynxflow/patter-go/pkg/whatsapp/controller"
	waSvc "github.com/lynxflow/patter-go/pkg/whatsapp/service"
)

// @title           Patter Engine Gateway & Campaign Engine API
// @version         1.0
// @description     Omnichannel gateway, campaign engine, Direct Universal SIP Trunking PBX, MCP (Model Context Protocol), Unified RAG Router, Recall.ai Meeting Bots & Evolution Go Proxy
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/lynxflow/patter-go

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
func main() {
	logger := core.InitLogger()
	logger.Info("Starting Patter Core Engine Gateway...")

	router := gin.Default()
	router.Use(core.CORSMiddleware())
	router.Use(core.AuthMiddleware())
	router.Use(core.TenantRLSMiddleware())

	// Initialize Services
	instSvc := instanceSvc.NewInstanceService()
	whatsappSvc := waSvc.NewWhatsAppService()
	campSvc := campaignSvc.NewCampaignService()
	vSvc := voiceSvc.NewVoiceService()
	msgSvc := messagingSvc.NewMessagingService()
	bSvc := brainSvc.NewBrainService()
	mcpSvc := mcpCtrl.NewMCPService(logger)
	meetingEngineSvc := rtc.NewMeetingEngineService(os.Getenv("PATTER_MEETING_ENGINE_KEY"), logger)
	sipSvc := sip.NewSIPPBXService(logger)
	avatarSvc := avatarCtrl.NewAvatarEngineService(os.Getenv("WAVESPEED_API_KEY"), logger)

	// Initialize Controllers
	instController := instanceCtrl.NewInstanceController(instSvc)
	waController := waCtrl.NewWhatsAppController(whatsappSvc)
	campaignController := campaignCtrl.NewCampaignController(campSvc)
	vController := voiceCtrl.NewVoiceController(vSvc)
	msgController := messagingCtrl.NewMessagingController(msgSvc)
	bController := brainCtrl.NewBrainController(bSvc)
	recallController := rtc.NewRecallController(meetingEngineSvc)
	sipController := sip.NewSIPController(sipSvc)
	mcpController := mcpCtrl.NewMCPController(mcpSvc)
	avatarController := avatarCtrl.NewAvatarController(avatarSvc)

	// API Gateway V1 Routes
	api := router.Group("/api/v1/gateway")
	{
		// Contract-First TypeScript Schema Endpoint
		api.GET("/schema/typescript", core.ServeTypeScriptSchema)

		// n8n / Make.com Community Node 1-Click Schema Endpoint
		api.GET("/integrations/n8n/node-schema", mcpController.GetN8NCommunityNodeSchema)

		// Instances
		api.POST("/instances", instController.CreateInstance)
		api.GET("/instances", instController.ListInstances)
		api.GET("/instances/:id", instController.GetInstance)
		api.DELETE("/instances/:id", instController.DeleteInstance)

		// Direct Universal SIP Trunking PBX (OVH, 3CX, FreeSWITCH, Asterisk, Aircall)
		api.POST("/sip/trunks", sipController.RegisterTrunk)
		api.GET("/sip/trunks", sipController.ListTrunks)
		api.POST("/sip/call", sipController.InitiateSIPCall)

		// Model Context Protocol (MCP) Integration & Patter Bridge
		api.POST("/mcp/servers", mcpController.RegisterServer)
		api.GET("/mcp/servers", mcpController.ListServers)
		api.GET("/mcp/servers/:id/tools", mcpController.ListTools)
		api.POST("/mcp/tools/execute", mcpController.ExecuteTool)

		// Avatar Video Engine (WaveSpeed API & Local GPU Renderer)
		api.POST("/avatar/stream", avatarController.StreamRealtimeAvatar)
		api.POST("/avatar/render", avatarController.RenderMassAvatarVideo)

		// WhatsApp Engine Proxy & Number Checker
		api.POST("/whatsapp/check-number", waController.CheckNumberExists)
		api.POST("/whatsapp/connect", waController.ConnectSession)
		api.GET("/whatsapp/qrcode", waController.GetQRCode)
		api.GET("/whatsapp/qrcode/:session", waController.GetQRCode)
		api.GET("/whatsapp/sessions", waController.ListSessions)
		api.POST("/whatsapp/message/send", waController.SendMessage)
		api.POST("/whatsapp/media/send", waController.SendMedia)
		api.POST("/whatsapp/location/send", waController.SendLocation)
		api.POST("/whatsapp/contact/send", waController.SendContact)

		// Omnichannel Campaigns
		api.POST("/campaigns", campaignController.CreateCampaign)
		api.GET("/campaigns", campaignController.ListCampaigns)
		api.GET("/campaigns/:id", campaignController.GetCampaign)
		api.POST("/campaigns/:id/pause", campaignController.PauseCampaign)
		api.POST("/campaigns/:id/resume", campaignController.ResumeCampaign)

		// Voice Campaigns, Profiles & Calls
		api.POST("/voice/campaign", campaignController.CreateVoiceCampaign)
		api.POST("/voice/profiles", vController.CreateVoiceProfile)
		api.GET("/voice/profiles", vController.ListVoiceProfiles)
		api.POST("/voice/call", vController.InitiateCall)

		// Messaging SMS/MMS & RCS
		api.POST("/messaging/sms", msgController.SendSMS)
		api.POST("/messaging/rcs", msgController.SendRCS)

		// AI Brain, Prompt Generator Copilot, Unified RAG, Calendar & AssemblyAI LeMUR v3
		api.POST("/brain/query", bController.QueryBrain)
		api.POST("/brain/prompt/enhance", bController.EnhancePrompt)
		api.POST("/brain/rag/search", bController.SearchRAG)
		api.POST("/brain/transfer", bController.HumanTransfer)
		api.POST("/brain/calendar/availability", bController.CalendarAvailability)
		api.POST("/brain/calendar/book", bController.CalendarBook)
		api.POST("/brain/lemur/process", bController.ProcessLeMUR)

		// Meeting Bots & Profiles (Patter White-Label Meeting Engine)
		api.POST("/rtc/profiles", recallController.CreateMeetingProfile)
		api.GET("/rtc/profiles", recallController.ListMeetingProfiles)
		api.POST("/rtc/bot", recallController.CreateBot)
		api.GET("/rtc/bots", recallController.ListBots)
		api.POST("/rtc/bots/:id/leave", recallController.LeaveMeeting)
	}

	// Serve Static Dashboard Web UI
	router.Static("/web", "./web")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	// Swagger API Docs Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server running on port %s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Listen error: %s\n", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	logger.Info("Server exiting")
}
