package agent

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestAgentListCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockAgentsResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/agent")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var agents []types.Agent
	if err := json.Unmarshal(resp, &agents); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(agents) == 0 {
		t.Error("Expected agents but got none")
	}

	if agents[0].ID != "default" {
		t.Errorf("Expected first agent ID 'default', got %s", agents[0].ID)
	}
}
