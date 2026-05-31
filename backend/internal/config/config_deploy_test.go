package config

import "testing"

func TestLoadUsesProductionOrientedEnvOverrides(t *testing.T) {
	t.Setenv("SERVER_PORT", "8080")
	t.Setenv("DATABASE_DSN", "postgres://dist_user:secret@postgres:5432/sub2api?sslmode=disable")
	t.Setenv("JWT_SECRET", "deploy-secret")
	t.Setenv("APP_ENV", "production")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://dist.example.com,https://ops.example.com")
	t.Setenv("STATIC_DIR", "/app/web")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppEnv != "production" {
		t.Fatalf("cfg.AppEnv = %q, want %q", cfg.AppEnv, "production")
	}
	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Fatalf("len(cfg.CORSAllowedOrigins) = %d, want 2", len(cfg.CORSAllowedOrigins))
	}
	if cfg.CORSAllowedOrigins[0] != "https://dist.example.com" {
		t.Fatalf("cfg.CORSAllowedOrigins[0] = %q, want dist origin", cfg.CORSAllowedOrigins[0])
	}
	if cfg.StaticDir != "/app/web" {
		t.Fatalf("cfg.StaticDir = %q, want %q", cfg.StaticDir, "/app/web")
	}
}
