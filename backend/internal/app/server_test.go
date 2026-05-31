package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/lhl/sub2api-distributor/backend/internal/config"
)

func TestIsAllowedDevOrigin(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		origin string
		want   bool
	}{
		{name: "allow localhost 5173", origin: "http://localhost:5173", want: true},
		{name: "allow 127.0.0.1 5176", origin: "http://127.0.0.1:5176", want: true},
		{name: "allow localhost 5177", origin: "http://localhost:5177", want: true},
		{name: "allow 127.0.0.1 5177", origin: "http://127.0.0.1:5177", want: true},
		{name: "reject empty", origin: "", want: false},
		{name: "reject unknown port", origin: "http://127.0.0.1:5180", want: false},
		{name: "reject foreign host", origin: "https://example.com", want: false},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isAllowedDevOrigin(tt.origin)
			if got != tt.want {
				t.Fatalf("isAllowedDevOrigin(%q) = %v, want %v", tt.origin, got, tt.want)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	server := NewServer(nil, config.Config{
		ServerPort:         "8091",
		AppEnv:             "development",
		CORSAllowedOrigins: nil,
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body == "" {
		t.Fatal("health body is empty")
	}
}

func TestServesFrontendIndexWhenStaticDirConfigured(t *testing.T) {
	t.Parallel()

	staticDir := t.TempDir()
	indexPath := filepath.Join(staticDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>portal</body></html>"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	server := NewServer(nil, config.Config{
		ServerPort:         "8091",
		AppEnv:             "production",
		CORSAllowedOrigins: []string{"https://dist.example.com"},
		StaticDir:          staticDir,
	})

	req := httptest.NewRequest(http.MethodGet, "/portal/dashboard", nil)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body != "<html><body>portal</body></html>" {
		t.Fatalf("body = %q, want frontend index", body)
	}
}
