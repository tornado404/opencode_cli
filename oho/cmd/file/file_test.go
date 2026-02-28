package file

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

func TestFileListCmd(t *testing.T) {

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestFileListCmd(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "list root",
			path:    "/",
			wantErr: false,
		},
		{
			name:    "list specific path",
			path:    "/src",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
					return testutil.MockFileListResponse(), nil
				},
			}

			resp, err := mock.GetWithQuery(context.Background(), "/file", map[string]string{"path": tt.path})
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
		})
	}
}

func TestFileContentCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockFileContentResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/file/main.go/content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var content types.FileContent
	if err := json.Unmarshal(resp, &content); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if content.Path != "main.go" {
		t.Errorf("Expected path 'main.go', got %s", content.Path)
	}
}

func TestFileStatusCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockFileStatusResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/file/status")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var files []types.File
	if err := json.Unmarshal(resp, &files); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}
