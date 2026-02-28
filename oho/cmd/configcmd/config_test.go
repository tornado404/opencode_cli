package configcmd

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestMain(m *testing.M) {
	// Initialize config for tests
	os.Setenv("OPENCODE_SERVER_HOST", "127.0.0.1")
	os.Setenv("OPENCODE_SERVER_PORT", "4096")
	os.Setenv("OPENCODE_SERVER_USERNAME", "opencode")
	os.Setenv("OPENCODE_SERVER_PASSWORD", "test")
	config.Init()

	m.Run()
}

func TestConfigGetCmd(t *testing.T) {

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestConfigGetCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockConfigResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/config")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var cfg types.Config
	if err := json.Unmarshal(resp, &cfg); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if cfg.DefaultModel != "gpt-4" {
		t.Errorf("Expected default model 'gpt-4', got %s", cfg.DefaultModel)
	}
	if cfg.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got %s", cfg.Theme)
	}
}

func TestConfigSetCmd(t *testing.T) {
	tests := []struct {
		name    string
		updates map[string]interface{}
		wantErr bool
	}{
		{
			name:    "update theme",
			updates: map[string]interface{}{"theme": "light"},
			wantErr: false,
		},
		{
			name:    "update language",
			updates: map[string]interface{}{"language": "zh"},
			wantErr: false,
		},
		{
			name:    "update model",
			updates: map[string]interface{}{"defaultModel": "gpt-3.5-turbo"},
			wantErr: false,
		},
		{
			name:    "update max tokens",
			updates: map[string]interface{}{"maxTokens": 8192},
			wantErr: false,
		},
		{
			name:    "update temperature",
			updates: map[string]interface{}{"temperature": 0.5},
			wantErr: false,
		},
		{
			name:    "empty updates",
			updates: map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				PatchFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					if len(tt.updates) == 0 {
						return nil, &client.APIError{StatusCode: 400, Message: "No updates provided"}
					}
					return testutil.MockConfigResponse(), nil
				},
			}

			resp, err := mock.Patch(context.Background(), "/config", tt.updates)
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && resp != nil {
				var cfg types.Config
				if err := json.Unmarshal(resp, &cfg); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
			}
		})
	}
}

func TestConfigProvidersCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			// Return as array directly (matching actual API)
			return testutil.MockProvidersResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/config/providers")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var providers []types.Provider
	if err := json.Unmarshal(resp, &providers); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(providers) == 0 {
		t.Error("Expected providers but got none")
	}
}
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockProvidersResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/config/providers")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var result struct {
		Providers []types.Provider  `json:"providers"`
		Default   map[string]string `json:"default"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(result.Providers) == 0 {
		t.Error("Expected providers but got none")
	}
}
