package lsp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/types"
)

func TestLSPStatusCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockLSPStatusResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/lsp/status")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var status []types.LSPStatus
	if err := json.Unmarshal(resp, &status); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(status) == 0 {
		t.Error("Expected LSP status but got none")
	}
}
