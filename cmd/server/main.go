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
	brainCtrl "github.com/lynxflow/patter-go/pkg/brain/controller"
	brainSvc "github.com/lynxflow/patter-go/pkg/brain/service"
	campaignCtrl "github.com/lynxflow/patter-go/pkg/campaign/controller"
	campaignSvc "github.com/lynxflow/patter-go/pkg/campaign/service"
	"github.com/lynxflow/patter-go/pkg/core"
	instanceCtrl "github.com/lynxflow/patter-go/pkg/instance/controller"
	instanceSvc "github.com/lynxflow/patter-go/pkg/instance/service"
	messagingCtrl "github.com/lynxflow/patter-go/pkg/messaging/controller"
	messagingSvc "github.com/lynxflow/patter-go/pkg/messaging/service"
	voiceCtrl "github.com/lynxflow/patter-go/pkg/voice/controller"
	voiceSvc "github.com/lynxflow/patter-go/pkg/voice/service"
	waCtrl "github.com/lynxflow/patter-go/pkg/whatsapp/controller"
	waSvc "github.com/lynxflow/patter-go/pkg/whatsapp/service"
)

// @title           Patter Engine Gateway & Campaign Engine API
// @version         1.0
// @description     Omnichannel gateway and campaign engine (WhatsApp Evolution Go, Voice, SMS, AI Brain)
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

	// Initialize Services
	instSvc := instanceSvc.NewInstanceService()
	whatsappSvc := waSvc.NewWhatsAppService()
	campSvc := campaignSvc.NewCampaignService()
	vSvc := voiceSvc.NewVoiceService()
	msgSvc := messagingSvc.NewMessagingService()
	bSvc := brainSvc.NewBrainService()

	// Initialize Controllers
	instController := instanceCtrl.NewInstanceController(instSvc)
	waController := waCtrl.NewWhatsAppController(whatsappSvc)
	campaignController := campaignCtrl.NewCampaignController(campSvc)
	vController := voiceCtrl.NewVoiceController(vSvc)
	msgController := messagingCtrl.NewMessagingController(msgSvc)
	bController := brainCtrl.NewBrainController(bSvc)

	// API Gateway V1 Routes
	api := router.Group("/api/v1/gateway")
	{
		// Instances
		api.POST("/instances", instController.CreateInstance)
		api.GET("/instances", instController.ListInstances)
		api.GET("/instances/:id", instController.GetInstance)
		api.DELETE("/instances/:id", instController.DeleteInstance)

		// WhatsApp Evolution Proxy
		api.POST("/whatsapp/connect", waController.ConnectSession)
		api.GET("/whatsapp/qrcode", waController.GetQRCode)
		api.GET("/whatsapp/qrcode/:session", waController.GetQRCode)
		api.GET("/whatsapp/sessions", waController.ListSessions)
		api.POST("/whatsapp/message/send", waController.SendMessage)

		// Omnichannel Campaigns
		api.POST("/campaigns", campaignController.CreateCampaign)
		api.GET("/campaigns", campaignController.ListCampaigns)
		api.GET("/campaigns/:id", campaignController.GetCampaign)
		api.POST("/campaigns/:id/pause", campaignController.PauseCampaign)
		api.POST("/campaigns/:id/resume", campaignController.ResumeCampaign)

		// Voice Campaigns & Calls
		api.POST("/voice/campaign", campaignController.CreateVoiceCampaign)
		api.POST("/voice/call", vController.InitiateCall)

		// Messaging SMS/MMS
		api.POST("/messaging/sms", msgController.SendSMS)

		// AI Brain
		api.POST("/brain/query", bController.QueryBrain)
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
