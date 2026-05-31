package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("JWT_SECRET", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ServerPort != "8091" {
		t.Fatalf("cfg.ServerPort = %q, want %q", cfg.ServerPort, "8091")
	}
	if cfg.DatabaseDSN != "postgres://sub2api@localhost:5432/sub2api?sslmode=disable" {
		t.Fatalf("cfg.DatabaseDSN = %q, want default DSN", cfg.DatabaseDSN)
	}
	if cfg.JWTSecret != "sub2api-distributor-dev-secret" {
		t.Fatalf("cfg.JWTSecret = %q, want default secret", cfg.JWTSecret)
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("SERVER_PORT", "9100")
	t.Setenv("DATABASE_DSN", "postgres://demo@localhost:5432/demo?sslmode=disable")
	t.Setenv("JWT_SECRET", "override-secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ServerPort != "9100" {
		t.Fatalf("cfg.ServerPort = %q, want %q", cfg.ServerPort, "9100")
	}
	if cfg.DatabaseDSN != "postgres://demo@localhost:5432/demo?sslmode=disable" {
		t.Fatalf("cfg.DatabaseDSN = %q, want override value", cfg.DatabaseDSN)
	}
	if cfg.JWTSecret != "override-secret" {
		t.Fatalf("cfg.JWTSecret = %q, want %q", cfg.JWTSecret, "override-secret")
	}
}
