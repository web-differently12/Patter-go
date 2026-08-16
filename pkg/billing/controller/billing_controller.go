package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/billing"
	"github.com/lynxflow/patter-go/pkg/core"
)

type BillingController struct {
	pricingGrid  *billing.PricingGrid
	meterService billing.UsageMeterService
	walletLedger billing.WalletLedgerService
	bridgeSvc    billing.HyperswitchLagoService
}

func NewBillingController(
	pricingGrid *billing.PricingGrid,
	meterService billing.UsageMeterService,
	walletLedger billing.WalletLedgerService,
	bridgeSvc billing.HyperswitchLagoService,
) *BillingController {
	return &BillingController{
		pricingGrid:  pricingGrid,
		meterService: meterService,
		walletLedger: walletLedger,
		bridgeSvc:    bridgeSvc,
	}
}

// GetWallet godoc
// @Summary      Visualiser le solde du Wallet & l'autonomie estimée
// @Tags         Gateway - Billing & Wallet
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} billing.TenantWallet
// @Router       /api/v1/gateway/billing/wallet [get]
func (ctrl *BillingController) GetWallet(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	wallet, err := ctrl.walletLedger.GetWallet(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Burn rate estimation: ~ $0.025 / voice minute => minutes left
	minutesLeft := int(wallet.BalanceFiat / 0.025)
	if minutesLeft < 0 {
		minutesLeft = 0
	}

	core.Success(c, gin.H{
		"wallet":               wallet,
		"estimated_autonomy":   gin.H{"voice_call_minutes": minutesLeft, "video_meetings": int(minutesLeft / 5)},
		"account_status_label": getAccountStatusLabel(wallet.BalanceFiat),
	})
}

// InitTopup godoc
// @Summary      Recharger le wallet via Hyperswitch Dropin ou Bitcoin Lightning
// @Tags         Gateway - Billing & Wallet
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body billing.CreateTopupRequest true "Top-up Config"
// @Success      200 {object} billing.HyperswitchSessionResponse
// @Router       /api/v1/gateway/billing/topup [post]
func (ctrl *BillingController) InitTopup(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	var req billing.CreateTopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := ctrl.bridgeSvc.InitHyperswitchTopup(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	core.Success(c, res)
}

// ConfigureAutoReload godoc
// @Summary      Paramétrer la recharge automatique au franchissement de seuil
// @Tags         Gateway - Billing & Wallet
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} billing.TenantWallet
// @Router       /api/v1/gateway/billing/auto-reload [post]
func (ctrl *BillingController) ConfigureAutoReload(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	var body struct {
		Enabled   bool    `json:"enabled"`
		Threshold float64 `json:"threshold"`
		Amount    float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	w, err := ctrl.walletLedger.SetAutoReloadConfig(c.Request.Context(), tenantID, body.Enabled, body.Threshold, body.Amount)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	core.Success(c, w)
}

// GetMarginsConfig godoc
// @Summary      Consulter ou simuler la règle de marge agence
// @Tags         Gateway - Billing & Wallet
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} billing.TenantPricingConfig
// @Router       /api/v1/gateway/billing/margins [get]
func (ctrl *BillingController) GetMarginsConfig(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	cfg := ctrl.pricingGrid.GetTenantConfig(tenantID)
	core.Success(c, cfg)
}

// UpdateMarginsConfig godoc
// @Summary      Mettre à jour la marge agence (% markup)
// @Tags         Gateway - Billing & Wallet
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} billing.TenantPricingConfig
// @Router       /api/v1/gateway/billing/margins [post]
func (ctrl *BillingController) UpdateMarginsConfig(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	var body struct {
		MarkupPercentage float64 `json:"markup_percentage"`
		FixedFeePerCall  float64 `json:"fixed_fee_per_call"`
		GraceLimitFiat   float64 `json:"grace_limit_fiat"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	cfg := billing.TenantPricingConfig{
		TenantID:         tenantID,
		MarkupPercentage: body.MarkupPercentage,
		FixedFeePerCall:  body.FixedFeePerCall,
		GraceLimitFiat:   body.GraceLimitFiat,
	}
	ctrl.pricingGrid.SetTenantConfig(cfg)
	core.Success(c, cfg)
}

// IngestMeteringEvent godoc
// @Summary      Ingérer un chunk de consommation (STT ms, LLM tokens, TTS chars, SIP sec)
// @Tags         Gateway - Billing & Wallet
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} billing.MeteringResult
// @Router       /api/v1/gateway/billing/metering/event [post]
func (ctrl *BillingController) IngestMeteringEvent(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	var evt billing.MeteringEvent
	if err := c.ShouldBindJSON(&evt); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	evt.TenantID = tenantID

	res, err := ctrl.meterService.IngestEvent(c.Request.Context(), evt)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	_ = ctrl.bridgeSvc.PushLagoUsageEvent(c.Request.Context(), evt, res)
	core.Success(c, res)
}

// GetLiveSessionMetering godoc
// @Summary      Visualiser le décompte en direct d'une session
// @Tags         Gateway - Billing & Wallet
// @Produce      json
// @Param        session_id query string true "ID de session"
// @Success      200 {object} billing.MeteringResult
// @Router       /api/v1/gateway/billing/metering/live [get]
func (ctrl *BillingController) GetLiveSessionMetering(c *gin.Context) {
	sessionID := c.Query("session_id")
	res, ok := ctrl.meterService.GetLiveSessionMetrics(sessionID)
	if !ok {
		core.Error(c, http.StatusNotFound, "live session metering data not found")
		return
	}
	core.Success(c, res)
}

// GetTransactionHistory godoc
// @Summary      Consulter le grand livre comptable des débits & rechargements
// @Tags         Gateway - Billing & Wallet
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {array} billing.WalletTransaction
// @Router       /api/v1/gateway/billing/history [get]
func (ctrl *BillingController) GetTransactionHistory(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	history, err := ctrl.walletLedger.GetTransactionHistory(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	core.Success(c, history)
}

// GetAdminProviderCostCatalog godoc
// @Summary      Inspecter le catalogue des coûts réels fournisseurs et options en temps réel
// @Tags         Gateway - Billing & Wallet Admin
// @Produce      json
// @Success      200 {array} billing.AdminProviderCostReport
// @Router       /api/v1/gateway/billing/admin/provider-costs [get]
func (ctrl *BillingController) GetAdminProviderCostCatalog(c *gin.Context) {
	reports, err := ctrl.bridgeSvc.GetAdminProviderCostCatalog(c.Request.Context())
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	core.Success(c, reports)
}

func getAccountStatusLabel(balance float64) string {
	if balance >= 10.0 {
		return "ACTIVE_HEALTHY"
	} else if balance > 0.0 {
		return "LOW_BALANCE_WARNING"
	}
	return "SUSPENDED_ZERO_BALANCE"
}
