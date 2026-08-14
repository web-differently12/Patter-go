package routes

import (
	"vocal-engine/pkg/engine"
	"vocal-engine/pkg/telephony"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter registers endpoints for telephony handlers and websocket streaming
func SetupRouter(tc *telephony.CallController, se *engine.StreamEngine) *gin.Engine {
	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Telephony group
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/call/outbound", tc.HandleOutboundCall)
		apiV1.GET("/call/twiml", tc.HandleTwiML)
	}

	// Bi-directional stream endpoint
	r.GET("/ws/twilio/stream", se.HandleTwilioStream)

	return r
}
