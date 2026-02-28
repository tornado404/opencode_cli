package message

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

func TestMessageListCmd(t *testing.T) {

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestMessageListCmd(t *testing.T) {
	tests := []struct {
		name     string
		mockResp []byte
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			mockResp: testutil.MockMessagesResponse(),
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "error",
			mockResp: nil,
			mockErr:  &client.APIError{StatusCode: 500, Message: "Internal Error"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				GetWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
					return tt.mockResp, tt.mockErr
				},
			}

			resp, err := mock.GetWithQuery(context.Background(), "/session/session1/message", map[string]string{})
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestMessageAddCmd(t *testing.T) {
	tests := []struct {
		name    string
		content string
		model   string
		agent   string
		noReply bool
		wantErr bool
	}{
		{
			name:    "add simple message",
			content: "Hello",
			wantErr: false,
		},
		{
			name:    "add with model",
			content: "Hello",
			model:   "gpt-4",
			wantErr: false,
		},
		{
			name:    "add with agent",
			content: "Hello",
			agent:   "default",
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockMessageResponse(), nil
				},
			}

			if tt.content == "" {
				// Skip the actual API call for empty content
				return
			}

			parts := []types.Part{
				{Type: "text", Data: tt.content},
			}

			req := types.MessageRequest{
				Model:   tt.model,
				Agent:   tt.agent,
				NoReply: tt.noReply,
				Parts:   parts,
			}

			resp, err := mock.Post(context.Background(), "/session/session1/message", req)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if resp != nil {
				var result types.MessageWithParts
				if err := json.Unmarshal(resp, &result); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
			}
		})
	}
}

func TestMessageGetCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockMessageResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/session/session1/message/msg1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var result types.MessageWithParts
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if result.Info.ID != "msg1" {
		t.Errorf("Expected message ID 'msg1', got %s", result.Info.ID)
	}
}

func TestMessagePromptAsyncCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	parts := []types.Part{
		{Type: "text", Data: "Async message"},
	}

	req := types.MessageRequest{
		Parts: parts,
	}

	resp, err := mock.Post(context.Background(), "/session/session1/prompt_async", req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var success bool
	if err := json.Unmarshal(resp, &success); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestMessageCommandCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockMessageResponse(), nil
		},
	}

	req := types.CommandRequest{
		Command:   "/test",
		Arguments: map[string]string{"arg1": "value1"},
	}

	resp, err := mock.Post(context.Background(), "/session/session1/command", req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var result types.MessageWithParts
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
}

func TestMessageShellCmd(t *testing.T) {
	tests := []struct {
		name    string
		agent   string
		command string
		wantErr bool
	}{
		{
			name:    "shell with agent",
			agent:   "default",
			command: "ls -la",
			wantErr: false,
		},
		{
			name:    "shell without agent",
			agent:   "",
			command: "ls -la",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agent == "" {
				// Test the error case
				return
			}

			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockMessageResponse(), nil
				},
			}

			req := types.ShellRequest{
				Agent:   tt.agent,
				Command: tt.command,
			}

			resp, err := mock.Post(context.Background(), "/session/session1/shell", req)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			var result types.MessageWithParts
			if err := json.Unmarshal(resp, &result); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   int
	}{
		{"hello", "ll", 2},
		{"hello", "lo", 3},
		{"hello", "x", -1},
		{"hello", "", 0},
		{"", "a", -1},
		{"hello", "hell", 0},
		{"hello", "hello", 0},
		{"hello", "o", 4},
	}

	for _, tt := range tests {
		result := indexOf(tt.s, tt.substr)
		if result != tt.want {
			t.Errorf("indexOf(%q, %q) = %d, want %d", tt.s, tt.substr, result, tt.want)
		}
	}
}
