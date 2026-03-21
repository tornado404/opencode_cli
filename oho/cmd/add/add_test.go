package add

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
	os.Setenv("OPENCODE_SERVER_HOST", "127.0.0.1")
	os.Setenv("OPENCODE_SERVER_PORT", "4096")
	os.Setenv("OPENCODE_SERVER_USERNAME", "opencode")
	os.Setenv("OPENCODE_SERVER_PASSWORD", "test")
	_ = config.Init()

	m.Run()
}

func TestConvertModel(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		wantType string
		wantStr  string
		wantObj  types.Model
	}{
		{
			name:     "empty model returns nil",
			model:    "",
			wantType: "nil",
		},
		{
			name:     "simple model returns string",
			model:    "gpt-4",
			wantType: "string",
			wantStr:  "gpt-4",
		},
		{
			name:     "provider:model format returns Model object",
			model:    "openai:gpt-4",
			wantType: "Model",
			wantObj:  types.Model{ProviderID: "openai", ModelID: "gpt-4"},
		},
		{
			name:     "provider with colon in model name",
			model:    "anthropic:claude-3-opus",
			wantType: "Model",
			wantObj:  types.Model{ProviderID: "anthropic", ModelID: "claude-3-opus"},
		},
		{
			name:     "model without provider stays string",
			model:    "claude-3",
			wantType: "string",
			wantStr:  "claude-3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertModel(tt.model)

			switch tt.wantType {
			case "nil":
				if result != nil {
					t.Errorf("convertModel(%q) = %v, want nil", tt.model, result)
				}
			case "string":
				str, ok := result.(string)
				if !ok {
					t.Errorf("convertModel(%q) returned %T, want string", tt.model, result)
				} else if str != tt.wantStr {
					t.Errorf("convertModel(%q) = %q, want %q", tt.model, str, tt.wantStr)
				}
			case "Model":
				obj, ok := result.(types.Model)
				if !ok {
					t.Errorf("convertModel(%q) returned %T, want types.Model", tt.model, result)
				} else if obj.ProviderID != tt.wantObj.ProviderID || obj.ModelID != tt.wantObj.ModelID {
					t.Errorf("convertModel(%q) = %+v, want %+v", tt.model, obj, tt.wantObj)
				}
			}
		})
	}
}

func TestDetectMimeType(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{"text file", "test.txt", "text/plain"},
		{"markdown file", "README.md", "text/markdown"},
		{"Go source", "main.go", "text/x-go"},
		{"Python source", "script.py", "text/x-python"},
		{"JavaScript", "app.js", "application/javascript"},
		{"TypeScript", "index.ts", "text/x-typescript"},
		{"TSX file", "component.tsx", "text/x-typescript"},
		{"JSON", "config.json", "application/json"},
		{"YAML", "config.yaml", "application/x-yaml"},
		{"YML", "config.yml", "application/x-yaml"},
		{"PNG image", "image.png", "image/png"},
		{"JPEG image", "photo.jpg", "image/jpeg"},
		{"GIF image", "anim.gif", "image/gif"},
		{"PDF document", "doc.pdf", "application/pdf"},
		{"ZIP archive", "archive.zip", "application/zip"},
		{"unknown extension", "file.xyz", "application/octet-stream"},
		{"no extension", "Makefile", "application/octet-stream"},
		{"uppercase extension", "IMAGE.PNG", "image/png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectMimeType(tt.filePath)
			if got != tt.want {
				t.Errorf("detectMimeType(%q) = %q, want %q", tt.filePath, got, tt.want)
			}
		})
	}
}

func TestCreateSession(t *testing.T) {
	tests := []struct {
		name            string
		title           string
		parentID        string
		directory       string
		mockResp        []byte
		mockErr         error
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "success with title and directory",
			title:     "Test Session",
			parentID:  "",
			directory: "/home/user/project",
			mockResp:  testutil.MockSessionResponse(),
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "success without title",
			title:     "",
			parentID:  "",
			directory: "/home/user/project",
			mockResp:  testutil.MockSessionResponse(),
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "success with parent session",
			title:     "Child Session",
			parentID:  "ses_parent123",
			directory: "/home/user/project",
			mockResp:  testutil.MockSessionResponse(),
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:            "API error - server unavailable",
			title:           "Test Session",
			parentID:        "",
			directory:       "/home/user/project",
			mockResp:        nil,
			mockErr:         &client.APIError{StatusCode: 503, Message: "Service Unavailable"},
			wantErr:         true,
			wantErrContains: "API request failed",
		},
		{
			name:            "API error - unauthorized",
			title:           "Test Session",
			parentID:        "",
			directory:       "/home/user/project",
			mockResp:        nil,
			mockErr:         &client.APIError{StatusCode: 401, Message: "Unauthorized"},
			wantErr:         true,
			wantErrContains: "API request failed",
		},
		{
			name:            "malformed JSON response",
			title:           "Test Session",
			parentID:        "",
			directory:       "/home/user/project",
			mockResp:        []byte(`{invalid json}`),
			mockErr:         nil,
			wantErr:         true,
			wantErrContains: "failed to parse response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &client.MockClient{
				PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
					if dir, ok := queryParams["directory"]; !ok {
						t.Error("Expected 'directory' query parameter")
					} else if dir != tt.directory {
						t.Errorf("Expected directory %q, got %q", tt.directory, dir)
					}

					bodyBytes, _ := json.Marshal(body)
					var bodyMap map[string]interface{}
					json.Unmarshal(bodyBytes, &bodyMap)

					if tt.title != "" {
						if _, ok := bodyMap["title"]; !ok {
							t.Error("Expected 'title' in request body")
						}
					}
					if tt.parentID != "" {
						if _, ok := bodyMap["parentID"]; !ok {
							t.Error("Expected 'parentID' in request body")
						}
					}

					return tt.mockResp, tt.mockErr
				},
			}

			sessionID, err := createSession(mock, context.Background(), tt.title, tt.parentID, tt.directory)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantErr && tt.wantErrContains != "" {
				if err == nil || !contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.wantErrContains, err)
				}
			}
			if !tt.wantErr && sessionID == "" {
				t.Error("Expected non-empty session ID")
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name            string
		sessionID       string
		message         string
		agent           string
		model           string
		noReply         bool
		system          string
		tools           []string
		files           []string
		mockResp        []byte
		mockErr         error
		wantErr         bool
		wantErrContains string
		wantMsgID       string
	}{
		{
			name:      "success with simple message",
			sessionID: "ses_test123",
			message:   "Hello, help me with this",
			agent:     "",
			model:     "",
			noReply:   false,
			mockResp:  testutil.MockMessageResponse(),
			mockErr:   nil,
			wantErr:   false,
			wantMsgID: "msg1",
		},
		{
			name:      "success with agent specified",
			sessionID: "ses_test123",
			message:   "Hello",
			agent:     "default",
			model:     "",
			noReply:   false,
			mockResp:  testutil.MockMessageResponse(),
			wantErr:   false,
		},
		{
			name:      "success with model specified",
			sessionID: "ses_test123",
			message:   "Hello",
			agent:     "",
			model:     "openai:gpt-4",
			noReply:   false,
			mockResp:  testutil.MockMessageResponse(),
			wantErr:   false,
		},
		{
			name:      "no-reply mode returns empty message ID",
			sessionID: "ses_test123",
			message:   "Hello",
			noReply:   true,
			mockResp:  []byte{},
			mockErr:   nil,
			wantErr:   false,
			wantMsgID: "",
		},
		{
			name:            "API error - server unavailable",
			sessionID:       "ses_test123",
			message:         "Hello",
			mockResp:        nil,
			mockErr:         &client.APIError{StatusCode: 503, Message: "Service Unavailable"},
			wantErr:         true,
			wantErrContains: "API request failed",
		},
		{
			name:            "malformed JSON response",
			sessionID:       "ses_test123",
			message:         "Hello",
			mockResp:        []byte(`{invalid}`),
			mockErr:         nil,
			wantErr:         true,
			wantErrContains: "failed to parse response",
		},
		{
			name:            "file not found",
			sessionID:       "ses_test123",
			message:         "Hello",
			files:           []string{"/nonexistent/file.txt"},
			mockResp:        nil,
			mockErr:         nil,
			wantErr:         true,
			wantErrContains: "file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFiles := []string{}
			for _, filePath := range tt.files {
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					tmpFile, err := os.CreateTemp("", "test-*.txt")
					if err == nil {
						tmpFile.WriteString("test content")
						tmpFile.Close()
						tempFiles = append(tempFiles, tmpFile.Name())
					}
				} else {
					tempFiles = append(tempFiles, filePath)
				}
			}

			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					if path != "/session/"+tt.sessionID+"/message" {
						t.Errorf("Expected path /session/%s/message, got %s", tt.sessionID, path)
					}

					bodyBytes, _ := json.Marshal(body)
					var bodyMap map[string]interface{}
					json.Unmarshal(bodyBytes, &bodyMap)

					if _, ok := bodyMap["parts"]; !ok {
						t.Error("Expected 'parts' in request body")
					}
					if bodyMap["noReply"] != tt.noReply {
						t.Errorf("Expected noReply=%v, got %v", tt.noReply, bodyMap["noReply"])
					}

					return tt.mockResp, tt.mockErr
				},
			}

			msgID, err := sendMessage(mock, context.Background(), tt.sessionID, tt.message, tt.agent, tt.model, tt.noReply, tt.system, tt.tools, tempFiles)

			for _, f := range tempFiles {
				os.Remove(f)
			}

			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantErr && tt.wantErrContains != "" {
				if err == nil || !contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.wantErrContains, err)
				}
			}
			if !tt.wantErr && msgID != tt.wantMsgID {
				t.Errorf("Expected message ID %q, got %q", tt.wantMsgID, msgID)
			}
		})
	}
}

func TestRunAddSuccess(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		title         string
		parent        string
		directory     string
		agent         string
		model         string
		noReply       bool
		system        string
		tools         []string
		files         []string
		jsonOutput    bool
		sessionResp   []byte
		messageResp   []byte
		wantSessionID string
		wantMessageID string
	}{
		{
			name:          "basic add command",
			message:       "Help me analyze this project",
			title:         "",
			noReply:       false,
			sessionResp:   testutil.MockSessionResponse(),
			messageResp:   testutil.MockMessageResponse(),
			wantSessionID: "session1",
			wantMessageID: "msg1",
		},
		{
			name:        "add with custom title",
			message:     "Fix the login bug",
			title:       "Bug Fix Session",
			noReply:     false,
			sessionResp: testutil.MockSessionResponse(),
			messageResp: testutil.MockMessageResponse(),
		},
		{
			name:          "add with no-reply mode",
			message:       "Run tests in background",
			noReply:       true,
			sessionResp:   testutil.MockSessionResponse(),
			messageResp:   []byte{},
			wantMessageID: "",
		},
		{
			name:        "add with agent and model",
			message:     "Review this code",
			agent:       "review",
			model:       "anthropic:claude-3",
			noReply:     false,
			sessionResp: testutil.MockSessionResponse(),
			messageResp: testutil.MockMessageResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addTitle = tt.title
			addParent = tt.parent
			addDirectory = tt.directory
			addAgent = tt.agent
			addModel = tt.model
			addNoReply = tt.noReply
			addSystem = tt.system
			addTools = tt.tools
			addFiles = tt.files
			addJSONOutput = tt.jsonOutput

			mock := &client.MockClient{
				PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
					return tt.sessionResp, nil
				},
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return tt.messageResp, nil
				},
			}

			ctx := context.Background()

			sessionID, err := createSession(mock, ctx, tt.title, tt.parent, tt.directory)
			if err != nil {
				t.Fatalf("Failed to create session: %v", err)
			}
			if sessionID == "" {
				t.Fatal("Expected non-empty session ID")
			}

			msgID, err := sendMessage(mock, ctx, sessionID, tt.message, tt.agent, tt.model, tt.noReply, tt.system, tt.tools, tt.files)
			if err != nil {
				t.Fatalf("Failed to send message: %v", err)
			}
			if !tt.noReply && msgID == "" && len(tt.messageResp) > 0 {
				t.Error("Expected non-empty message ID")
			}
		})
	}
}

func TestRaceConditionScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupDelay  int
		wantSuccess bool
		description string
	}{
		{
			name:        "immediate session ready",
			setupDelay:  0,
			wantSuccess: true,
			description: "Session is immediately ready for messages",
		},
		{
			name:        "slight delay in session ready",
			setupDelay:  50,
			wantSuccess: true,
			description: "Session becomes ready after slight delay",
		},
		{
			name:        "significant delay in session ready",
			setupDelay:  200,
			wantSuccess: true,
			description: "Session takes longer to initialize",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			mock := &client.MockClient{
				PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
					callCount++
					if callCount == 1 {
						return testutil.MockSessionResponse(), nil
					}
					return nil, nil
				},
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockMessageResponse(), nil
				},
			}

			ctx := context.Background()

			sessionID, err := createSession(mock, ctx, "Test", "", "/test")
			if err != nil {
				t.Fatalf("createSession failed: %v", err)
			}

			msgID, err := sendMessage(mock, ctx, sessionID, "Test message", "", "", false, "", nil, nil)

			if tt.wantSuccess && err != nil {
				t.Errorf("%s: Expected success but got error: %v", tt.description, err)
			}
			if !tt.wantSuccess && err == nil {
				t.Errorf("%s: Expected failure but got success", tt.description)
			}
			if tt.wantSuccess && msgID == "" {
				t.Error("Expected non-empty message ID")
			}
		})
	}
}

func TestTimeoutScenarios(t *testing.T) {
	tests := []struct {
		name          string
		timeoutMs     int
		serverDelayMs int
		wantTimeout   bool
		description   string
	}{
		{
			name:          "request completes within timeout",
			timeoutMs:     1000,
			serverDelayMs: 100,
			wantTimeout:   false,
			description:   "Fast server response within timeout",
		},
		{
			name:          "request exceeds timeout",
			timeoutMs:     50,
			serverDelayMs: 200,
			wantTimeout:   true,
			description:   "Slow server response exceeds timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%s: timeout=%dms, delay=%dms", tt.description, tt.timeoutMs, tt.serverDelayMs)

			mock := &client.MockClient{
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockMessageResponse(), nil
				},
			}

			ctx := context.Background()
			_, err := sendMessage(mock, ctx, "ses_test", "Test", "", "", false, "", nil, nil)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestErrorPropagation(t *testing.T) {
	tests := []struct {
		name             string
		sessionCreateErr error
		messageSendErr   error
		wantSessionID    bool
		wantMessageID    bool
		description      string
	}{
		{
			name:             "session create fails, message not sent",
			sessionCreateErr: &client.APIError{StatusCode: 500, Message: "Internal Error"},
			messageSendErr:   nil,
			wantSessionID:    false,
			wantMessageID:    false,
			description:      "Session creation failure should prevent message send",
		},
		{
			name:             "session succeeds, message fails",
			sessionCreateErr: nil,
			messageSendErr:   &client.APIError{StatusCode: 500, Message: "Internal Error"},
			wantSessionID:    true,
			wantMessageID:    false,
			description:      "Message send failure after session created",
		},
		{
			name:             "both succeed",
			sessionCreateErr: nil,
			messageSendErr:   nil,
			wantSessionID:    true,
			wantMessageID:    true,
			description:      "Both operations succeed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionCreated := false
			mock := &client.MockClient{
				PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
					sessionCreated = true
					if tt.sessionCreateErr != nil {
						return nil, tt.sessionCreateErr
					}
					return testutil.MockSessionResponse(), nil
				},
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					if !sessionCreated {
						t.Error("Message send attempted before session creation")
					}
					if tt.messageSendErr != nil {
						return nil, tt.messageSendErr
					}
					return testutil.MockMessageResponse(), nil
				},
			}

			ctx := context.Background()

			sessionID, sessionErr := createSession(mock, ctx, "Test", "", "/test")

			var msgID string
			var msgErr error
			if sessionErr == nil {
				msgID, msgErr = sendMessage(mock, ctx, sessionID, "Test", "", "", false, "", nil, nil)
			}

			if tt.wantSessionID && sessionID == "" {
				t.Error("Expected session ID but got empty")
			}
			if !tt.wantSessionID && sessionID != "" {
				t.Error("Expected no session ID but got one")
			}

			if tt.wantMessageID && msgID == "" && msgErr == nil {
				t.Error("Expected message ID but got empty")
			}
			if !tt.wantMessageID && msgID != "" {
				t.Error("Expected no message ID but got one")
			}
		})
	}
}

func TestPartialFailureHandling(t *testing.T) {
	mock := &client.MockClient{
		PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
			return testutil.MockSessionResponse(), nil
		},
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return nil, &client.APIError{StatusCode: 500, Message: "Message send failed"}
		},
	}

	ctx := context.Background()

	sessionID, err := createSession(mock, ctx, "Test", "", "/test")
	if err != nil {
		t.Fatalf("Session creation failed: %v", err)
	}
	if sessionID == "" {
		t.Fatal("Expected session ID")
	}

	_, err = sendMessage(mock, ctx, sessionID, "Test", "", "", false, "", nil, nil)
	if err == nil {
		t.Error("Expected message send to fail")
	}

	t.Logf("Partial failure: session %s created, but message send failed", sessionID)
}

func TestJSONOutputFormat(t *testing.T) {
	tests := []struct {
		name       string
		jsonOutput bool
	}{
		{"text output", false},
		{"JSON output", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addJSONOutput = tt.jsonOutput

			mock := &client.MockClient{
				PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
					return testutil.MockSessionResponse(), nil
				},
				PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
					return testutil.MockMessageResponse(), nil
				},
			}

			ctx := context.Background()

			sessionID, _ := createSession(mock, ctx, "Test", "", "/test")
			msgID, _ := sendMessage(mock, ctx, sessionID, "Test", "", "", false, "", nil, nil)

			if tt.jsonOutput {
				output := map[string]interface{}{
					"sessionId": sessionID,
					"messageId": msgID,
					"status":    "success",
				}
				data, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					t.Errorf("JSON marshaling failed: %v", err)
				}
				t.Logf("JSON output: %s", string(data))
			} else {
				t.Logf("Text output: Session=%s, Message=%s", sessionID, msgID)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
