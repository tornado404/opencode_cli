package command

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestCommandListCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockCommandsResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/command")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var commands []types.Command
	if err := json.Unmarshal(resp, &commands); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(commands) == 0 {
		t.Error("Expected commands but got none")
	}
}
