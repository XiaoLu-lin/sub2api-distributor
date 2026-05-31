package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/auth"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/middleware"
)

// AuthHandler exposes login-session endpoints for the distributor frontend.
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates the HTTP adapter for authentication requests.
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login verifies credentials and returns a role-aware portal session token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("请求参数不正确"))
		return
	}

	user, portalRole, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			status = http.StatusUnauthorized
		case errors.Is(err, auth.ErrInactiveUser), errors.Is(err, auth.ErrDistributorNotEnabled):
			status = http.StatusForbidden
		}
		respondError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":       token,
		"user":        user,
		"portal_role": portalRole,
	})
}

// Logout is a best-effort stateless logout endpoint kept for frontend symmetry.
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Me returns the currently authenticated principal as parsed from the bearer token.
func (h *AuthHandler) Me(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, errors.New("未登录或登录已失效"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":     claims.UserID,
		"email":       claims.Email,
		"portal_role": claims.PortalRole,
	})
}
