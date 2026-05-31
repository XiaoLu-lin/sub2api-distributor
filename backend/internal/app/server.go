package app

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/auth"
	"github.com/lhl/sub2api-distributor/backend/internal/config"
	"github.com/lhl/sub2api-distributor/backend/internal/distributor"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/handlers"
	"github.com/lhl/sub2api-distributor/backend/internal/httpapi/routes"
)

// NewServer wires all application services, HTTP handlers, routes, and local
// development middleware into a runnable HTTP server.
func NewServer(db *sql.DB, cfg config.Config) *http.Server {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isAllowedOrigin(origin, cfg.CORSAllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			c.Header("Access-Control-Max-Age", "43200")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authService := auth.NewService(db, cfg.JWTSecret)
	distributorService := distributor.NewService(db)

	authHandler := handlers.NewAuthHandler(authService)
	portalHandler := handlers.NewPortalHandler(distributorService)
	opsHandler := handlers.NewOpsHandler(distributorService)

	routes.Register(router, authService, authHandler, portalHandler, opsHandler)
	registerFrontend(router, cfg.StaticDir)

	return &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}
}

func registerFrontend(router *gin.Engine, staticDir string) {
	staticDir = strings.TrimSpace(staticDir)
	if staticDir == "" {
		return
	}

	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return
		}
		return
	}

	router.Static("/assets", filepath.Join(staticDir, "assets"))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"message": "接口不存在"})
			return
		}
		if path == "/health" {
			c.JSON(http.StatusNotFound, gin.H{"message": "请求失败"})
			return
		}
		c.File(indexPath)
	})
}

// isAllowedDevOrigin constrains browser CORS access to the known local
// development origins used by the distributor frontend.
func isAllowedDevOrigin(origin string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}

	switch origin {
	case "http://127.0.0.1:5173",
		"http://localhost:5173",
		"http://127.0.0.1:5176",
		"http://localhost:5176",
		"http://127.0.0.1:5177",
		"http://localhost:5177":
		return true
	default:
		return false
	}
}

func isAllowedOrigin(origin string, configured []string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}

	for _, allowed := range configured {
		if origin == strings.TrimSpace(allowed) {
			return true
		}
	}

	return isAllowedDevOrigin(origin)
}
