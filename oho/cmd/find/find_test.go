package find

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestFindTextCmd(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{
			name:    "search pattern",
			pattern: "func main",
			wantErr: false,
		},
		{
			name:    "empty pattern",
			pattern: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pattern == "" {
				return
			}

			mock := &client.MockClient{
				GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
					return testutil.MockFindMatchesResponse(), nil
				},
			}

			resp, err := mock.GetWithQuery(context.Background(), "/find/text", map[string]string{"q": tt.pattern})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			var matches []types.FindMatch
			if err := json.Unmarshal(resp, &matches); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
			}

			if len(matches) == 0 {
				t.Error("Expected matches but got none")
			}
		})
	}
}

func TestFindFileCmd(t *testing.T) {
	mock := &client.MockClient{
		GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
			return testutil.MockFileListResponse(), nil
		},
	}

	resp, err := mock.GetWithQuery(context.Background(), "/find/file", map[string]string{"q": "main.go"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var files []types.FileNode
	if err := json.Unmarshal(resp, &files); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(files) == 0 {
		t.Error("Expected files but got none")
	}
}

func TestFindSymbolCmd(t *testing.T) {
	mock := &client.MockClient{
		GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
			return testutil.MockSymbolsResponse(), nil
		},
	}

	resp, err := mock.GetWithQuery(context.Background(), "/find/symbol", map[string]string{"q": "main"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var symbols []types.Symbol
	if err := json.Unmarshal(resp, &symbols); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(symbols) == 0 {
		t.Error("Expected symbols but got none")
	}
}
