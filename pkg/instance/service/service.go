package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/instance/dto"
)

type InstanceService interface {
	CreateInstance(ctx context.Context, tenantID string, req dto.CreateInstanceRequest) (*dto.InstanceResponse, error)
	GetInstance(ctx context.Context, tenantID, instanceID string) (*dto.InstanceResponse, error)
	ListInstances(ctx context.Context, tenantID string) ([]*dto.InstanceResponse, error)
	DeleteInstance(ctx context.Context, tenantID, instanceID string) error
}

type instanceService struct {
	mu        sync.RWMutex
	instances map[string]*dto.InstanceResponse
}

func NewInstanceService() InstanceService {
	return &instanceService{
		instances: make(map[string]*dto.InstanceResponse),
	}
}

func (s *instanceService) CreateInstance(ctx context.Context, tenantID string, req dto.CreateInstanceRequest) (*dto.InstanceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	instID := fmt.Sprintf("inst_%s_%s", tenantID, uuid.New().String()[:8])
	inst := &dto.InstanceResponse{
		InstanceID:   instID,
		TenantID:     tenantID,
		InstanceName: req.InstanceName,
		Status:       "CREATED",
		WebhookURL:   req.WebhookURL,
		CreatedAt:    time.Now(),
	}

	s.instances[instID] = inst
	return inst, nil
}

func (s *instanceService) GetInstance(ctx context.Context, tenantID, instanceID string) (*dto.InstanceResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	inst, ok := s.instances[instanceID]
	if !ok || inst.TenantID != tenantID {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}
	return inst, nil
}

func (s *instanceService) ListInstances(ctx context.Context, tenantID string) ([]*dto.InstanceResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*dto.InstanceResponse
	for _, inst := range s.instances {
		if inst.TenantID == tenantID {
			result = append(result, inst)
		}
	}
	return result, nil
}

func (s *instanceService) DeleteInstance(ctx context.Context, tenantID, instanceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	inst, ok := s.instances[instanceID]
	if !ok || inst.TenantID != tenantID {
		return fmt.Errorf("instance %s not found", instanceID)
	}

	delete(s.instances, instanceID)
	return nil
}
