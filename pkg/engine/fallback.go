package engine

import (
	"context"
	"errors"
	"log"
)

// FallbackManager orchestrates LLM call routing with built-in mid-call failover capabilities
type FallbackManager struct {
	providers map[string]LLMProvider
}

// NewFallbackManager creates a new FallbackManager
func NewFallbackManager(providers ...LLMProvider) *FallbackManager {
	provMap := make(map[string]LLMProvider)
	for _, p := range providers {
		provMap[p.GetProviderID()] = p
	}
	return &FallbackManager{
		providers: provMap,
	}
}

// ExecuteWithFallback executes an LLM generation with automated sequential failover mid-call
func (fm *FallbackManager) ExecuteWithFallback(ctx context.Context, primaryID string, fallbackChain []string, systemPrompt string, userMessage string, tools []AgentTool) (*LLMResponse, string, error) {
	// 1. Try Primary LLM Provider
	if primaryProvider, ok := fm.providers[primaryID]; ok {
		log.Printf("[FallbackManager] Executing completion using primary provider: %s", primaryID)
		resp, err := primaryProvider.GenerateCompletion(ctx, systemPrompt, userMessage, tools)
		if err == nil {
			return resp, primaryID, nil
		}
		log.Printf("[FallbackManager] Primary provider %s failed: %v. Initiating fallback failover chain...", primaryID, err)
	} else {
		log.Printf("[FallbackManager] Primary provider %s not found in registry. Running fallbacks...", primaryID)
	}

	// 2. Iterate through Fallback Chain sequentially
	for _, providerID := range fallbackChain {
		if provider, ok := fm.providers[providerID]; ok {
			log.Printf("[FallbackManager] Attempting fallback provider: %s", providerID)
			resp, err := provider.GenerateCompletion(ctx, systemPrompt, userMessage, tools)
			if err == nil {
				log.Printf("[FallbackManager] Failover successful! Active provider: %s", providerID)
				return resp, providerID, nil
			}
			log.Printf("[FallbackManager] Fallback provider %s failed: %v", providerID, err)
		}
	}

	return nil, "", errors.New("all providers in fallback chain failed")
}
