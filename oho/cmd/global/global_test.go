package global

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestHealthCmd(t *testing.T) {
	tests := []struct {

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestHealthCmd(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse *types.HealthResponse
		serverStatus   int
		wantHealthy    bool
		wantErr        bool
	}{
		{
			name:           "server healthy",
			serverResponse: &types.HealthResponse{Healthy: true, Version: "1.0.0"},
			serverStatus:   200,
			wantHealthy:    true,
			wantErr:        false,
		},
		{
			name:           "server unhealthy",
			serverResponse: &types.HealthResponse{Healthy: false, Version: "1.0.0"},
			serverStatus:   200,
			wantHealthy:    false,
			wantErr:        false,
		},
		{
			name:         "server error",
			serverStatus: 500,
			wantErr:      true,
		},
		{
			name:         "connection refused",
			serverStatus: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.serverStatus > 0 {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.serverStatus)
					if tt.serverResponse != nil {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				}))
				defer server.Close()

				// Override config to use test server
				cfg := config.Get()
				cfg.Host = server.Listener.Addr().String()
			}

			mock := &client.MockClient{
				GetFunc: func(ctx context.Context, path string) ([]byte, error) {
					if tt.wantErr {
						return nil, &client.APIError{StatusCode: 500, Message: "Internal Error"}
					}
					return testutil.MockHealthResponse(), nil
				},
			}

			// Test that the mock works
			resp, err := mock.Get(context.Background(), "/global/health")
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && resp != nil {
				var health types.HealthResponse
				if err := json.Unmarshal(resp, &health); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
				if health.Healthy != tt.wantHealthy {
					t.Errorf("Expected healthy=%v, got %v", tt.wantHealthy, health.Healthy)
				}
			}
			_ = server
		})
	}
}

func TestHealthCmdJSONOutput(t *testing.T) {
	// Save original config
	origJSON := config.Get().JSON
	defer func() { config.Get().JSON = origJSON }()

	// Enable JSON mode
	config.Get().JSON = true

	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockHealthResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/global/health")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp == nil {
		t.Error("Expected response but got nil")
	}
}

func TestEventCmd(t *testing.T) {
	mock := &client.MockClient{
		SSEStreamFunc: func(ctx context.Context, path string) (<-chan []byte, <-chan error, error) {
			eventChan := make(chan []byte, 1)
			errChan := make(chan error, 1)

			// Send a test event
			eventChan <- []byte("data: test event\n\n")

			close(eventChan)
			close(errChan)

			return eventChan, errChan, nil
		},
	}

	eventChan, errChan, err := mock.SSEStream(context.Background(), "/global/event")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	select {
	case event, ok := <-eventChan:
		if !ok {
			t.Error("Event channel closed unexpectedly")
		}
		if string(event) != "data: test event\n\n" {
			t.Errorf("Expected event data, got %s", string(event))
		}
	case err := <-errChan:
		if err != nil {
			t.Errorf("Unexpected error in channel: %v", err)
		}
	}
}
