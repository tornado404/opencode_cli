package session

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

func TestSessionListCmd(t *testing.T) {

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestSessionListCmd(t *testing.T) {
	tests := []struct {
		name       string
		mockResp   []byte
		mockErr    error
		statusCode int
		wantErr    bool
	}{
		{
			name:       "success with sessions",
			mockResp:   testutil.MockSessionsResponse(),
			mockErr:    nil,
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "empty sessions",
			mockResp:   testutil.MockResponse([]types.Session{}),
			mockErr:    nil,
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "server error",
			mockResp:   nil,
			mockErr:    &client.APIError{StatusCode: 500, Message: "Internal Error"},
			statusCode: 500,
			wantErr:    true,
		},
		{
			name:       "invalid JSON",
			mockResp:   []byte("invalid json"),
			mockErr:    nil,
			statusCode: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				GetFunc: func(ctx context.Context, path string) ([]byte, error) {
					return tt.mockResp, tt.mockErr
				},
			}

			resp, err := mock.Get(context.Background(), "/session")
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && resp != nil {
				var sessions []types.Session
				if err := json.Unmarshal(resp, &sessions); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
			}
		})
	}
}

func TestSessionCreateCmd(t *testing.T) {
	tests := []struct {
		name       string
		parentID   string
		title      string
		mockResp   []byte
		mockErr    error
		statusCode int
		wantErr    bool
	}{
		{
			name:       "create simple session",
			title:      "Test Session",
			mockResp:   testutil.MockSessionResponse(),
			mockErr:    nil,
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "create with parent",
			parentID:   "parent-123",
			title:      "Child Session",
			mockResp:   testutil.MockSessionResponse(),
			mockErr:    nil,
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "server error",
			title:      "Test Session",
			mockResp:   nil,
			mockErr:    &client.APIError{StatusCode: 500, Message: "Internal Error"},
			statusCode: 500,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return tt.mockResp, tt.mockErr
				},
			}

			resp, err := mock.Post(context.Background(), "/session", map[string]interface{}{"title": tt.title})
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && resp != nil {
				var session types.Session
				if err := json.Unmarshal(resp, &session); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
			}
		})
	}
}

func TestSessionStatusCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockSessionStatusResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/session/status")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var status map[string]types.SessionStatus
	if err := json.Unmarshal(resp, &status); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(status) != 2 {
		t.Errorf("Expected 2 sessions, got %d", len(status))
	}
}

func TestSessionGetCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockSessionResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/session/session1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var session types.Session
	if err := json.Unmarshal(resp, &session); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if session.ID != "session1" {
		t.Errorf("Expected session ID 'session1', got %s", session.ID)
	}
}

func TestSessionDeleteCmd(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		wantErr bool
	}{
		{
			name:    "delete success",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "delete error",
			mockErr: &client.APIError{StatusCode: 404, Message: "Session not found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				DeleteFunc: func(ctx context.Context, path string) ([]byte, error) {
					return testutil.MockBoolResponse(true), tt.mockErr
				},
			}

			resp, err := mock.Delete(context.Background(), "/session/session1")
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && resp != nil {
				var deleted bool
				if err := json.Unmarshal(resp, &deleted); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
				if !deleted {
					t.Error("Expected deleted=true")
				}
			}
		})
	}
}

func TestSessionUpdateCmd(t *testing.T) {
	mock := &client.MockClient{
		PatchFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockSessionResponse(), nil
		},
	}

	resp, err := mock.Patch(context.Background(), "/session/session1", map[string]interface{}{"title": "New Title"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var session types.Session
	if err := json.Unmarshal(resp, &session); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionChildrenCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockSessionsResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/session/session1/children")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var sessions []types.Session
	if err := json.Unmarshal(resp, &sessions); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionTodoCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockTodoResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/session/session1/todo")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var todos []types.Todo
	if err := json.Unmarshal(resp, &todos); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestSessionForkCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockSessionResponse(), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/fork", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var session types.Session
	if err := json.Unmarshal(resp, &session); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionAbortCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/abort", nil)
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

func TestSessionShareCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockSessionResponse(), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/share", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var session types.Session
	if err := json.Unmarshal(resp, &session); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionUnshareCmd(t *testing.T) {
	mock := &client.MockClient{
		DeleteFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockSessionResponse(), nil
		},
	}

	resp, err := mock.Delete(context.Background(), "/session/session1/share")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var session types.Session
	if err := json.Unmarshal(resp, &session); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionDiffCmd(t *testing.T) {
	mock := &client.MockClient{
		GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
			return testutil.MockDiffResponse(), nil
		},
	}

	resp, err := mock.GetWithQuery(context.Background(), "/session/session1/diff", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var diffs []types.FileDiff
	if err := json.Unmarshal(resp, &diffs); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionSummarizeCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/summarize", map[string]interface{}{"providerID": "openai", "modelID": "gpt-4"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var success bool
	if err := json.Unmarshal(resp, &success); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionRevertCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/revert", map[string]interface{}{"messageID": "msg1"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var success bool
	if err := json.Unmarshal(resp, &success); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionUnrevertCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/unrevert", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var success bool
	if err := json.Unmarshal(resp, &success); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestSessionPermissionsCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/permissions/perm1", map[string]interface{}{"response": "allow"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var success bool
	if err := json.Unmarshal(resp, &success); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

// Helper function to create test server
func createTestServer(handlers map[string]http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	for path, handler := range handlers {
		mux.Handle(path, handler)
	}
	return httptest.NewServer(mux)
}
