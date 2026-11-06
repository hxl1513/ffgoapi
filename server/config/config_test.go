package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "8080")
	os.Setenv("GIN_MODE", "test")
	os.Setenv("DB_TYPE", "sqlite")
	os.Setenv("DB_NAME", "test.db")
	os.Setenv("SESSION_SECRET", "test-secret")

	cfg, err := LoadFromEnv()
	if err != nil { t.Fatalf("LoadFromEnv() error = %v", err) }

	if cfg.Server.Port != 8080 { t.Errorf("Port = %d, want 8080", cfg.Server.Port) }
	if cfg.Server.Mode != "test" { t.Errorf("Mode = %s, want test", cfg.Server.Mode) }
	if cfg.Database.Type != "sqlite" { t.Errorf("Database.Type = %s, want sqlite", cfg.Database.Type) }
	if cfg.Database.DBName != "test.db" { t.Errorf("Database.DBName = %s, want test.db", cfg.Database.DBName) }
	if cfg.Auth.SessionSecret != "test-secret" { t.Errorf("Auth.SessionSecret = %s, want test-secret", cfg.Auth.SessionSecret) }
}

func TestLoadFromEnvDefaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("SESSION_SECRET", "test-secret")

	cfg, err := LoadFromEnv()
	if err != nil { t.Fatalf("LoadFromEnv() error = %v", err) }

	if cfg.Server.Port != 3000 { t.Errorf("Port = %d, want 3000 (default)", cfg.Server.Port) }
	if cfg.Server.Mode != "release" { t.Errorf("Mode = %s, want release (default)", cfg.Server.Mode) }
	if cfg.Server.ReadTimeout != 10*time.Second { t.Errorf("ReadTimeout = %v, want 10s", cfg.Server.ReadTimeout) }
	if cfg.Database.Type != "sqlite" { t.Errorf("Database.Type = %s, want sqlite", cfg.Database.Type) }
	if cfg.RateLimit.DefaultMax != 100 { t.Errorf("RateLimit.DefaultMax = %d, want 100", cfg.RateLimit.DefaultMax) }
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Server:   ServerConfig{Port: 3000},
				Database: DatabaseConfig{Type: "sqlite", DBName: "test.db"},
				Auth:     AuthConfig{SessionSecret: "secret"},
			},
			wantErr: false,
		},
		{
			name: "missing session secret",
			config: &Config{
				Server:   ServerConfig{Port: 3000},
				Database: DatabaseConfig{Type: "sqlite", DBName: "test.db"},
				Auth:     AuthConfig{SessionSecret: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr { t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr) }
		})
	}
}
