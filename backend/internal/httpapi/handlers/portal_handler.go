package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/distributor"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/middleware"
)

// PortalHandler exposes distributor-facing portal endpoints.
type PortalHandler struct {
	distributorService *distributor.Service
}

// NewPortalHandler creates the HTTP adapter for portal routes.
func NewPortalHandler(distributorService *distributor.Service) *PortalHandler {
	return &PortalHandler{distributorService: distributorService}
}

// Dashboard returns the authenticated distributor's current rebate summary.
func (h *PortalHandler) Dashboard(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	summary, err := h.distributorService.GetSummary(c.Request.Context(), claims.UserID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}

// Invitees returns users attributed to the current distributor's invitation chain.
func (h *PortalHandler) Invitees(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	items, err := h.distributorService.ListInvitees(c.Request.Context(), claims.UserID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// InviteMeta returns the real invite code the distributor can share externally.
func (h *PortalHandler) InviteMeta(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	meta, err := h.distributorService.GetInviteMeta(c.Request.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, distributor.ErrInviteMetaNotFound) {
			status = http.StatusNotFound
		}
		respondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, meta)
}

// Rebates returns accrue ledger entries that contribute to distributor earnings.
func (h *PortalHandler) Rebates(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	items, err := h.distributorService.ListRebates(c.Request.Context(), claims.UserID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// Withdrawals returns the current distributor's payout request history.
func (h *PortalHandler) Withdrawals(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	items, err := h.distributorService.ListWithdrawals(c.Request.Context(), claims.UserID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createWithdrawalRequest struct {
	Amount float64 `json:"amount"`
	Remark string  `json:"remark"`
}

// CreateWithdrawal submits a new payout request and immediately places it into paying status.
func (h *PortalHandler) CreateWithdrawal(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	var req createWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("请求参数不正确"))
		return
	}
	item, err := h.distributorService.CreateWithdrawal(c.Request.Context(), claims.UserID, req.Amount, req.Remark)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, distributor.ErrInvalidWithdrawalAmount), errors.Is(err, distributor.ErrWithdrawalAmountTooLarge):
			status = http.StatusBadRequest
		}
		respondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// CancelWithdrawal lets the distributor void one of their own paying requests.
func (h *PortalHandler) CancelWithdrawal(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("提现申请 ID 不正确"))
		return
	}
	if err := h.distributorService.CancelWithdrawal(c.Request.Context(), claims.UserID, requestID); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, distributor.ErrWithdrawalNotFound):
			status = http.StatusNotFound
		case errors.Is(err, distributor.ErrInvalidWithdrawalTransition):
			status = http.StatusBadRequest
		}
		respondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetProfile returns the distributor's settlement profile for display or editing.
func (h *PortalHandler) GetProfile(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	profile, err := h.distributorService.GetProfile(c.Request.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, distributor.ErrProfileNotFound) {
			status = http.StatusNotFound
		}
		respondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile saves the distributor's settlement profile using the authenticated user id.
func (h *PortalHandler) UpdateProfile(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	var profile distributor.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("请求参数不正确"))
		return
	}
	profile.UserID = claims.UserID
	if err := h.distributorService.UpsertProfile(c.Request.Context(), profile); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, distributor.ErrInvalidProfile) {
			status = http.StatusBadRequest
		}
		respondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
