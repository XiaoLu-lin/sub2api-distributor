package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/distributor"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/middleware"
)

// OpsHandler exposes operator-facing management endpoints.
type OpsHandler struct {
	distributorService *distributor.Service
}

// NewOpsHandler creates the HTTP adapter for operator routes.
func NewOpsHandler(distributorService *distributor.Service) *OpsHandler {
	return &OpsHandler{distributorService: distributorService}
}

// ListDistributors returns all distributor profiles visible to operators.
func (h *OpsHandler) ListDistributors(c *gin.Context) {
	items, err := h.distributorService.ListDistributorProfiles(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// LookupUsers searches main-system users so operators can enable distributor access.
func (h *OpsHandler) LookupUsers(c *gin.Context) {
	items, err := h.distributorService.LookupUsers(c.Request.Context(), c.Query("q"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetDistributor returns both profile details and balance summary for one distributor.
func (h *OpsHandler) GetDistributor(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("用户 ID 不正确"))
		return
	}
	profile, err := h.distributorService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, distributor.ErrProfileNotFound) {
			status = http.StatusNotFound
		}
		respondError(c, status, err)
		return
	}
	summary, err := h.distributorService.GetSummary(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile, "summary": summary})
}

// UpdateDistributor creates or updates a distributor profile for the target user.
func (h *OpsHandler) UpdateDistributor(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("用户 ID 不正确"))
		return
	}
	var profile distributor.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("请求参数不正确"))
		return
	}
	profile.UserID = userID
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

// ListWithdrawals returns the full operator payout queue.
func (h *OpsHandler) ListWithdrawals(c *gin.Context) {
	items, err := h.distributorService.ListAllWithdrawals(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetWithdrawal returns one payout request together with its event timeline.
func (h *OpsHandler) GetWithdrawal(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("提现申请 ID 不正确"))
		return
	}
	item, err := h.distributorService.GetWithdrawalDetail(c.Request.Context(), requestID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, distributor.ErrWithdrawalNotFound) {
			status = http.StatusNotFound
		}
		respondError(c, status, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

type markPaidRequest struct {
	PaidChannel     string `json:"paid_channel"`
	PaidReferenceNo string `json:"paid_reference_no"`
	PaidRemark      string `json:"paid_remark"`
}

// MarkPaid records that the operator completed the offline payout.
func (h *OpsHandler) MarkPaid(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("提现申请 ID 不正确"))
		return
	}
	claims, _ := middleware.GetClaims(c)
	var req markPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("请求参数不正确"))
		return
	}
	if err := h.distributorService.MarkWithdrawalPaid(c.Request.Context(), claims.UserID, requestID, req.PaidChannel, req.PaidReferenceNo, req.PaidRemark); err != nil {
		respondError(c, operatorMutationStatus(err), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Cancel voids a paying withdrawal request from the operator console.
func (h *OpsHandler) Cancel(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("提现申请 ID 不正确"))
		return
	}
	claims, _ := middleware.GetClaims(c)
	if err := h.distributorService.CancelWithdrawalAsOperator(c.Request.Context(), claims.UserID, requestID); err != nil {
		respondError(c, operatorMutationStatus(err), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
