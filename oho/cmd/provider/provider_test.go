package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestProviderListCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			type ProviderListResult struct {
				All       []types.Provider  `json:"all"`
				Default   map[string]string `json:"default"`
				Connected []string          `json:"connected"`
			}
			result := ProviderListResult{
				All: []types.Provider{
					{ID: "openai", Name: "OpenAI", BaseURL: "https://api.openai.com"},
					{ID: "anthropic", Name: "Anthropic", BaseURL: "https://api.anthropic.com"},
				},
				Default:   map[string]string{"default": "gpt-4"},
				Connected: []string{"openai"},
			}
			return json.Marshal(result)
		},
	}

	resp, err := mock.Get(context.Background(), "/provider")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var result struct {
		All       []types.Provider  `json:"all"`
		Default   map[string]string `json:"default"`
		Connected []string          `json:"connected"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(result.All) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(result.All))
	}
}

func TestProviderAuthCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			methods := map[string][]types.ProviderAuthMethod{
				"openai": {
					{Type: "api_key", Required: true, Description: "API Key authentication"},
				},
			}
			return json.Marshal(methods)
		},
	}

	resp, err := mock.Get(context.Background(), "/provider/auth")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var methods map[string][]types.ProviderAuthMethod
	if err := json.Unmarshal(resp, &methods); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(methods) == 0 {
		t.Error("Expected auth methods but got none")
	}
}

func TestProviderOAuthAuthorizeCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			auth := types.ProviderAuthAuthorization{
				URL:           "https://auth.example.com/authorize",
				State:         "state123",
				CodeChallenge: "challenge123",
			}
			return json.Marshal(auth)
		},
	}

	resp, err := mock.Post(context.Background(), "/provider/openai/oauth/authorize", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var auth types.ProviderAuthAuthorization
	if err := json.Unmarshal(resp, &auth); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if auth.URL == "" {
		t.Error("Expected auth URL but got empty string")
	}
}

func TestProviderOAuthCallbackCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/provider/openai/oauth/callback", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var success bool
	if err := json.Unmarshal(resp, &success); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if !success {
		t.Error("Expected success=true")
	}
}
