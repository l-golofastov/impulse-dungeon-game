package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"Floors": 3,
		"Monsters": 4,
		"OpenAt": "09:30:15",
		"Duration": 2
	}`

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("ReadConfig returned error: %v", err)
	}

	if cfg.Floors != 3 {
		t.Fatalf("Floors = %d, want 3", cfg.Floors)
	}
	if cfg.Monsters != 4 {
		t.Fatalf("Monsters = %d, want 4", cfg.Monsters)
	}
	wantOpenAt := 9*time.Hour + 30*time.Minute + 15*time.Second
	if cfg.OpenAt != wantOpenAt {
		t.Fatalf("OpenAt = %v, want %v", cfg.OpenAt, wantOpenAt)
	}
	if cfg.CloseAt != wantOpenAt+2*time.Hour {
		t.Fatalf("CloseAt = %v, want %v", cfg.CloseAt, wantOpenAt+2*time.Hour)
	}
}

func TestReadConfigErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "invalid json",
			content: `{`,
		},
		{
			name: "invalid open time",
			content: `{
				"Floors": 2,
				"Monsters": 2,
				"OpenAt": "bad",
				"Duration": 1
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("write config: %v", err)
			}

			if _, err := ReadConfig(path); err == nil {
				t.Fatal("ReadConfig returned nil error")
			}
		})
	}
}

func TestReadConfigMissingFile(t *testing.T) {
	if _, err := ReadConfig(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("ReadConfig returned nil error for missing file")
	}
}
