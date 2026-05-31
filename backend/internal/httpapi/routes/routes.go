package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/auth"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/handlers"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/middleware"
)

// Register mounts every API route used by the distributor and operator consoles.
func Register(
	router *gin.Engine,
	authService *auth.Service,
	authHandler *handlers.AuthHandler,
	portalHandler *handlers.PortalHandler,
	opsHandler *handlers.OpsHandler,
) {
	api := router.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)

		authenticated := api.Group("")
		authenticated.Use(middleware.RequireAuth(authService))
		authenticated.POST("/auth/logout", authHandler.Logout)
		authenticated.GET("/me", authHandler.Me)

		portal := authenticated.Group("/portal")
		// Portal routes are only available to enabled distributor users.
		portal.Use(middleware.RequirePortalRole(auth.PortalRoleDistributor))
		portal.GET("/dashboard", portalHandler.Dashboard)
		portal.GET("/invite-meta", portalHandler.InviteMeta)
		portal.GET("/invitees", portalHandler.Invitees)
		portal.GET("/rebates", portalHandler.Rebates)
		portal.GET("/withdrawals", portalHandler.Withdrawals)
		portal.POST("/withdrawals", portalHandler.CreateWithdrawal)
		portal.POST("/withdrawals/:id/cancel", portalHandler.CancelWithdrawal)
		portal.GET("/settlement-profile", portalHandler.GetProfile)
		portal.PUT("/settlement-profile", portalHandler.UpdateProfile)

		ops := authenticated.Group("/ops")
		// Operator routes are reserved for admin users from the main system.
		ops.Use(middleware.RequirePortalRole(auth.PortalRoleOperator))
		ops.GET("/distributors", opsHandler.ListDistributors)
		ops.GET("/users/lookup", opsHandler.LookupUsers)
		ops.GET("/distributors/:userId", opsHandler.GetDistributor)
		ops.PUT("/distributors/:userId/profile", opsHandler.UpdateDistributor)
		ops.GET("/withdrawals", opsHandler.ListWithdrawals)
		ops.GET("/withdrawals/:id", opsHandler.GetWithdrawal)
		ops.POST("/withdrawals/:id/mark-paid", opsHandler.MarkPaid)
		ops.POST("/withdrawals/:id/cancel", opsHandler.Cancel)
	}
}
