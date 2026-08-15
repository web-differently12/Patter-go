package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("TWILIO_ACCOUNT_SID", "ACmocked")
	os.Setenv("TWILIO_AUTH_TOKEN", "mockedtoken")
	os.Setenv("TWILIO_PHONE_NUMBER", "+12345")
	os.Setenv("RABBITMQ_URL", "amqp://test")
	os.Setenv("RABBITMQ_QUEUE", "test_queue")
	os.Setenv("REALTIME_API_URL", "wss://test-openai")
	os.Setenv("REALTIME_API_KEY", "testkey")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("TWILIO_ACCOUNT_SID")
		os.Unsetenv("TWILIO_AUTH_TOKEN")
		os.Unsetenv("TWILIO_PHONE_NUMBER")
		os.Unsetenv("RABBITMQ_URL")
		os.Unsetenv("RABBITMQ_QUEUE")
		os.Unsetenv("REALTIME_API_URL")
		os.Unsetenv("REALTIME_API_KEY")
	}()

	cfg := LoadConfig()

	if cfg.Port != "9090" {
		t.Errorf("expected Port 9090, got %s", cfg.Port)
	}
	if cfg.TwilioAccountSID != "ACmocked" {
		t.Errorf("expected TwilioAccountSID ACmocked, got %s", cfg.TwilioAccountSID)
	}
	if cfg.TwilioAuthToken != "mockedtoken" {
		t.Errorf("expected TwilioAuthToken mockedtoken, got %s", cfg.TwilioAuthToken)
	}
	if cfg.TwilioPhoneNumber != "+12345" {
		t.Errorf("expected TwilioPhoneNumber +12345, got %s", cfg.TwilioPhoneNumber)
	}
	if cfg.RabbitMQURL != "amqp://test" {
		t.Errorf("expected RabbitMQURL amqp://test, got %s", cfg.RabbitMQURL)
	}
	if cfg.RabbitMQQueueName != "test_queue" {
		t.Errorf("expected RabbitMQQueueName test_queue, got %s", cfg.RabbitMQQueueName)
	}
	if cfg.RealtimeAPIURL != "wss://test-openai" {
		t.Errorf("expected RealtimeAPIURL wss://test-openai, got %s", cfg.RealtimeAPIURL)
	}
	if cfg.RealtimeAPIKey != "testkey" {
		t.Errorf("expected RealtimeAPIKey testkey, got %s", cfg.RealtimeAPIKey)
	}
}
