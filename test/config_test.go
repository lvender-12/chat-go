package test

import (
	"os"
	"path/filepath"
	"testing"

	"chat-go/internal/config"
)

func TestLoadConfig_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{
		"app": {
			"name": "chat-go",
			"host": "localhost",
			"port": 3000,
			"log_level": "debug"
		},
		"database": {
			"driver": "mysql",
			"host": "localhost",
			"port": 3306,
			"user": "root",
			"password": "secret",
			"name": "chat",
			"charset": "utf8mb4",
			"parse_time": true
		},
		"jwt": {
			"secret": "jwt-secret",
			"exp": 24
		},
		"migration": {
			"enabled": true,
			"path": "migrations"
		}
	}`

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.App.Name != "chat-go" {
		t.Errorf("App.Name = %q, want %q", cfg.App.Name, "chat-go")
	}
	if cfg.Database.Driver != "mysql" {
		t.Errorf("Database.Driver = %q, want %q", cfg.Database.Driver, "mysql")
	}
	if cfg.JWT.Secret != "jwt-secret" {
		t.Errorf("JWT.Secret = %q, want %q", cfg.JWT.Secret, "jwt-secret")
	}
	if cfg.JWT.Exp != 24 {
		t.Errorf("JWT.Exp = %d, want %d", cfg.JWT.Exp, 24)
	}
	if !cfg.Migration.Enabled {
		t.Fatal("Migration.Enabled = false, want true")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := config.LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want error")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want error")
	}
}
