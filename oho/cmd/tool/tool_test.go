package tool

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestToolIDsCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockToolIDsResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/tool/ids")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var ids types.ToolIDs
	if err := json.Unmarshal(resp, &ids); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(ids.IDs) == 0 {
		t.Error("Expected tool IDs but got none")
	}
}

func TestToolListCmd(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		model    string
		wantErr  bool
	}{
		{
			name:    "list all tools",
			wantErr: false,
		},
		{
			name:     "list tools with provider",
			provider: "openai",
			wantErr:  false,
		},
		{
			name:    "list tools with model",
			model:   "gpt-4",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
					return testutil.MockToolListResponse(), nil
				},
			}

			queryParams := map[string]string{}
			if tt.provider != "" {
				queryParams["provider"] = tt.provider
			}
			if tt.model != "" {
				queryParams["model"] = tt.model
			}

			resp, err := mock.GetWithQuery(context.Background(), "/tool", queryParams)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			var tools types.ToolList
			if err := json.Unmarshal(resp, &tools); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
			}

			if len(tools.Tools) == 0 {
				t.Error("Expected tools but got none")
			}
		})
	}
}
