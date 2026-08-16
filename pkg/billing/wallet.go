package billing

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type TransactionType string

const (
	TxUsageDebit       TransactionType = "USAGE_DEBIT"
	TxHyperswitchTopup TransactionType = "HYPERSWITCH_TOPUP"
	TxBitcoinTopup     TransactionType = "BITCOIN_TOPUP"
	TxRefund           TransactionType = "REFUND"
)

type TenantWallet struct {
	TenantID            string    `json:"tenant_id"`
	BalanceFiat         float64   `json:"balance_fiat"`     // USD / EUR balance
	BalanceSatoshis     int64     `json:"balance_satoshis"` // Bitcoin Lightning Satoshis
	Currency            string    `json:"currency"`         // "EUR" or "USD"
	AutoReloadEnabled   bool      `json:"auto_reload_enabled"`
	AutoReloadThreshold float64   `json:"auto_reload_threshold"` // e.g. 10.00
	AutoReloadAmount    float64   `json:"auto_reload_amount"`    // e.g. 50.00
	UpdatedAt           time.Time `json:"updated_at"`
}

type WalletTransaction struct {
	ID           string                 `json:"id"`
	TenantID     string                 `json:"tenant_id"`
	SessionID    string                 `json:"session_id,omitempty"`
	Type         TransactionType        `json:"type"`
	Amount       float64                `json:"amount"`
	BalanceAfter float64                `json:"balance_after"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

type WalletLedgerService interface {
	GetWallet(ctx context.Context, tenantID string) (*TenantWallet, error)
	DeductUsage(ctx context.Context, tenantID, sessionID string, amount float64, evt MeteringEvent) (float64, error)
	TopUpWallet(ctx context.Context, tenantID string, amountFiat float64, satoshis int64, txType TransactionType, reference string) (*TenantWallet, error)
	GetTransactionHistory(ctx context.Context, tenantID string) ([]*WalletTransaction, error)
	SetAutoReloadConfig(ctx context.Context, tenantID string, enabled bool, threshold, amount float64) (*TenantWallet, error)
}

type walletLedgerService struct {
	mu           sync.RWMutex
	wallets      map[string]*TenantWallet
	transactions map[string][]*WalletTransaction
	logger       *slog.Logger
}

func NewWalletLedgerService(logger *slog.Logger) WalletLedgerService {
	if logger == nil {
		logger = core.GetLogger()
	}
	svc := &walletLedgerService{
		wallets:      make(map[string]*TenantWallet),
		transactions: make(map[string][]*WalletTransaction),
		logger:       logger.With("component", "wallet_ledger"),
	}

	// Seed default tenant wallet for instant demo/testing
	defaultTenant := "default_tenant"
	svc.wallets[defaultTenant] = &TenantWallet{
		TenantID:            defaultTenant,
		BalanceFiat:         50.00,
		BalanceSatoshis:     125000, // ~125k sats
		Currency:            "EUR",
		AutoReloadEnabled:   true,
		AutoReloadThreshold: 10.00,
		AutoReloadAmount:    50.00,
		UpdatedAt:           time.Now(),
	}

	return svc
}

func (s *walletLedgerService) GetWallet(ctx context.Context, tenantID string) (*TenantWallet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, ok := s.wallets[tenantID]
	if !ok {
		// Auto-initialize new tenant wallet with $10 welcome credit
		w = &TenantWallet{
			TenantID:            tenantID,
			BalanceFiat:         10.00,
			BalanceSatoshis:     25000,
			Currency:            "EUR",
			AutoReloadEnabled:   false,
			AutoReloadThreshold: 5.00,
			AutoReloadAmount:    25.00,
			UpdatedAt:           time.Now(),
		}
		s.wallets[tenantID] = w
	}
	return w, nil
}

func (s *walletLedgerService) DeductUsage(ctx context.Context, tenantID, sessionID string, amount float64, evt MeteringEvent) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.wallets[tenantID]
	if !ok {
		w = &TenantWallet{
			TenantID:    tenantID,
			BalanceFiat: 10.00,
			Currency:    "EUR",
			UpdatedAt:   time.Now(),
		}
		s.wallets[tenantID] = w
	}

	w.BalanceFiat -= amount
	w.UpdatedAt = time.Now()

	tx := &WalletTransaction{
		ID:           "tx_" + uuid.New().String()[:8],
		TenantID:     tenantID,
		SessionID:    sessionID,
		Type:         TxUsageDebit,
		Amount:       -amount,
		BalanceAfter: w.BalanceFiat,
		Metadata: map[string]interface{}{
			"provider":   evt.Provider,
			"model_name": evt.ModelName,
			"channel":    evt.ChannelType,
		},
		CreatedAt: time.Now(),
	}

	s.transactions[tenantID] = append(s.transactions[tenantID], tx)
	return w.BalanceFiat, nil
}

func (s *walletLedgerService) TopUpWallet(ctx context.Context, tenantID string, amountFiat float64, satoshis int64, txType TransactionType, reference string) (*TenantWallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.wallets[tenantID]
	if !ok {
		w = &TenantWallet{
			TenantID:  tenantID,
			Currency:  "EUR",
			UpdatedAt: time.Now(),
		}
		s.wallets[tenantID] = w
	}

	w.BalanceFiat += amountFiat
	w.BalanceSatoshis += satoshis
	w.UpdatedAt = time.Now()

	tx := &WalletTransaction{
		ID:           "tx_" + uuid.New().String()[:8],
		TenantID:     tenantID,
		Type:         txType,
		Amount:       amountFiat,
		BalanceAfter: w.BalanceFiat,
		Metadata: map[string]interface{}{
			"reference": reference,
			"satoshis":  satoshis,
		},
		CreatedAt: time.Now(),
	}

	s.transactions[tenantID] = append(s.transactions[tenantID], tx)
	s.logger.Info("Wallet successfully topped up", "tenant_id", tenantID, "amount_fiat", amountFiat, "satoshis", satoshis, "new_balance", w.BalanceFiat)
	return w, nil
}

func (s *walletLedgerService) GetTransactionHistory(ctx context.Context, tenantID string) ([]*WalletTransaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	txs, ok := s.transactions[tenantID]
	if !ok {
		return []*WalletTransaction{}, nil
	}
	return txs, nil
}

func (s *walletLedgerService) SetAutoReloadConfig(ctx context.Context, tenantID string, enabled bool, threshold, amount float64) (*TenantWallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.wallets[tenantID]
	if !ok {
		return nil, fmt.Errorf("wallet not found for tenant %s", tenantID)
	}

	w.AutoReloadEnabled = enabled
	w.AutoReloadThreshold = threshold
	w.AutoReloadAmount = amount
	w.UpdatedAt = time.Now()

	s.logger.Info("Auto-reload settings updated", "tenant_id", tenantID, "enabled", enabled, "threshold", threshold, "amount", amount)
	return w, nil
}
