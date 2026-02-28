package tui

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
)

func TestTUIOpenHelpCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/tui/open-help", nil)
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

func TestTUIShowToastCmd(t *testing.T) {
	tests := []struct {
		name    string
		message string
		title   string
		variant string
	}{
		{
			name:    "show toast with message",
			message: "Test message",
			title:   "Test",
			variant: "info",
		},
		{
			name:    "show toast error",
			message: "Error occurred",
			title:   "Error",
			variant: "error",
		},
		{
			name:    "show toast success",
			message: "Operation successful",
			title:   "Success",
			variant: "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockBoolResponse(true), nil
				},
			}

			req := map[string]interface{}{
				"message": tt.message,
				"title":   tt.title,
				"variant": tt.variant,
			}

			resp, err := mock.Post(context.Background(), "/tui/show-toast", req)
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
