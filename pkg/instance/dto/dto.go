package dto

import "time"

type CreateInstanceRequest struct {
	InstanceName string `json:"instance_name" binding:"required"`
	Description  string `json:"description,omitempty"`
	WebhookURL   string `json:"webhook_url,omitempty"`
}

type InstanceResponse struct {
	InstanceID   string    `json:"instance_id"`
	TenantID     string    `json:"tenant_id"`
	InstanceName string    `json:"instance_name"`
	Status       string    `json:"status"` // "CREATED", "CONNECTED", "DISCONNECTED"
	WebhookURL   string    `json:"webhook_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type InstanceStatusResponse struct {
	InstanceID string `json:"instance_id"`
	TenantID   string `json:"tenant_id"`
	Status     string `json:"status"`
	State      string `json:"state"` // "WORKING", "DISCONNECTED", "CONNECTING"
}
