package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"vocal-engine/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ToolCallPayload is the structured message published to RabbitMQ
type ToolCallPayload struct {
	TenantID  string          `json:"tenantId"`
	CallSID   string          `json:"callSid"`
	ToolName  string          `json:"toolName"`
	Arguments json.RawMessage `json:"arguments"`
	Timestamp string          `json:"timestamp"`
}

// Publisher defines an interface to publish tool calls asynchronously
type Publisher interface {
	PublishToolCall(ctx context.Context, tenantID, callSID, toolName string, arguments json.RawMessage) error
	Close() error
}

// RabbitMQPublisher implements the Publisher interface using RabbitMQ
type RabbitMQPublisher struct {
	cfg        *config.Config
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQPublisher connects to RabbitMQ and initializes the publisher
func NewRabbitMQPublisher(cfg *config.Config) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		cfg.RabbitMQQueueName, // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQPublisher{
		cfg:        cfg,
		connection: conn,
		channel:    ch,
	}, nil
}

// PublishToolCall publishes a tool call payload asynchronously without blocking the main engine loop
func (p *RabbitMQPublisher) PublishToolCall(ctx context.Context, tenantID, callSID, toolName string, arguments json.RawMessage) error {
	payload := ToolCallPayload{
		TenantID:  tenantID,
		CallSID:   callSID,
		ToolName:  toolName,
		Arguments: arguments,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal tool call payload: %w", err)
	}

	// We use a context with timeout for publishing to ensure it doesn't block forever
	pubCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(pubCtx,
		"",                 // exchange
		p.cfg.RabbitMQQueueName, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}

	log.Printf("[RabbitMQ] Published tool call for tool: %s, tenant: %s, callSid: %s", toolName, tenantID, callSID)
	return nil
}

// Close gracefully closes RabbitMQ channel and connection
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.connection != nil {
		return p.connection.Close()
	}
	return nil
}
