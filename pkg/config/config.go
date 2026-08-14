package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port               string
	TwilioAccountSID   string
	TwilioAuthToken    string
	TwilioPhoneNumber  string
	RabbitMQURL        string
	RabbitMQQueueName  string
	RealtimeAPIURL     string
	RealtimeAPIKey     string
}

// LoadConfig loads application configurations from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:               getEnv("SERVER_PORT", "8080"),
		TwilioAccountSID:   getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:    getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber:  getEnv("TWILIO_PHONE_NUMBER", ""),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQQueueName:  getEnv("RABBITMQ_QUEUE", "agent_tool_calls"),
		RealtimeAPIURL:     getEnv("REALTIME_API_URL", "wss://api.openai.com/v1/realtime?model=gpt-4o-realtime-preview-2024-10-01"),
		RealtimeAPIKey:     getEnv("REALTIME_API_KEY", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
