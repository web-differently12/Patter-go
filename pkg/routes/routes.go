package routes

import (
	"vocal-engine/pkg/engine"
	"vocal-engine/pkg/telephony"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter registers endpoints for telephony handlers, agent management, and websocket streaming
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

	// Telephony & Agents API Group
	apiV1 := r.Group("/api/v1")
	{
		// Outbound/Inbound Calling
		apiV1.POST("/call/outbound", tc.HandleOutboundCall)
		apiV1.POST("/call/inbound", ac.HandleInboundCall)
		apiV1.GET("/call/twiml", tc.HandleTwiML)

		// Agent Management
		apiV1.POST("/agents", ac.HandleCreateAgent)
		apiV1.GET("/agents", ac.HandleListAgents)
	}

	// Bi-directional stream endpoint
	r.GET("/ws/twilio/stream", se.HandleTwilioStream)

	return r
}
