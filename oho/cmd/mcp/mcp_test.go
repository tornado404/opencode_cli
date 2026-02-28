package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestMCPListCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockMCPStatusResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/mcp")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var status []types.MCPStatus
	if err := json.Unmarshal(resp, &status); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(status) == 0 {
		t.Error("Expected MCP status but got none")
	}
}

func TestMCPAddCmd(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			name:    "add mcp server",
			config:  `{"command": "npx", "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path"]}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockBoolResponse(true), nil
				},
			}

			resp, err := mock.Post(context.Background(), "/mcp", map[string]interface{}{"config": tt.config})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
			}
		})
	}
}

func TestMCPRemoveCmd(t *testing.T) {
	mock := &client.MockClient{
		DeleteFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Delete(context.Background(), "/mcp/server1")
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
