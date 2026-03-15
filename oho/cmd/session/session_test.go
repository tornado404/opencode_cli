package session

import (
	"context"
	"encoding/json"
	"os"
	"sort"
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

func TestApplyFieldFilter(t *testing.T) {
	testSessions := []types.Session{
		{
			ID:        "ses_abc123",
			Title:     "多帧融合提取回合状态信息",
			ProjectID: "1f01524d641bdfc5f4e43134de956b66c0b1332b",
			Directory: "/mnt/d/code/rl_mockgame",
			Time:      types.SessionTime{Created: 1773537883643, Updated: 1773538142930},
			Model:     "gpt-4",
		},
		{
			ID:        "ses_def456",
			Title:     "初始化 Git 仓库并推送 GitHub",
			ProjectID: "1f01524d641bdfc5f4e43134de956b66c0b1332b",
			Directory: "/mnt/d/code/rl_mockgame",
			Time:      types.SessionTime{Created: 1773411633123, Updated: 1773412006307},
			Model:     "gpt-3.5-turbo",
		},
		{
			ID:        "ses_ghi789",
			Title:     "Another Project",
			ProjectID: "proj2",
			Directory: "/home/user/project2",
			Time:      types.SessionTime{Created: 1773400000000, Updated: 1773400000000},
			Model:     "claude-3",
		},
	}

	tests := []struct {
		name          string
		setupFilters  func()
		expectedCount int
		expectedIDs   []string
		description   string
	}{
		{
			name: "no filters",
			setupFilters: func() {
				filterID = ""
				filterTitle = ""
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = ""
				filterDirectory = ""
			},
			expectedCount: 3,
			description:   "无过滤条件时返回所有会话",
		},
		{
			name: "filter by ID exact",
			setupFilters: func() {
				filterID = "ses_abc123"
				filterTitle = ""
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = ""
				filterDirectory = ""
			},
			expectedCount: 1,
			expectedIDs:   []string{"ses_abc123"},
			description:   "按 ID 精确匹配",
		},
		{
			name: "filter by ID fuzzy",
			setupFilters: func() {
				filterID = "abc"
				filterTitle = ""
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = ""
				filterDirectory = ""
			},
			expectedCount: 1,
			expectedIDs:   []string{"ses_abc123"},
			description:   "按 ID 模糊查询",
		},
		{
			name: "filter by title fuzzy",
			setupFilters: func() {
				filterID = ""
				filterTitle = "Git"
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = ""
				filterDirectory = ""
			},
			expectedCount: 1,
			expectedIDs:   []string{"ses_def456"},
			description:   "按标题模糊查询",
		},
		{
			name: "filter by projectID",
			setupFilters: func() {
				filterID = ""
				filterTitle = ""
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = "1f01524d"
				filterDirectory = ""
			},
			expectedCount: 2,
			expectedIDs:   []string{"ses_abc123", "ses_def456"},
			description:   "按项目 ID 模糊查询",
		},
		{
			name: "filter by directory",
			setupFilters: func() {
				filterID = ""
				filterTitle = ""
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = ""
				filterDirectory = "rl_mockgame"
			},
			expectedCount: 2,
			expectedIDs:   []string{"ses_abc123", "ses_def456"},
			description:   "按目录模糊查询",
		},
		{
			name: "filter by created timestamp",
			setupFilters: func() {
				filterID = ""
				filterTitle = ""
				filterCreated = 1773537883643
				filterUpdated = 0
				filterProjectID = ""
				filterDirectory = ""
			},
			expectedCount: 1,
			expectedIDs:   []string{"ses_abc123"},
			description:   "按创建时间精确匹配",
		},
		{
			name: "filter by updated timestamp",
			setupFilters: func() {
				filterID = ""
				filterTitle = ""
				filterCreated = 0
				filterUpdated = 1773412006307
				filterProjectID = ""
				filterDirectory = ""
			},
			expectedCount: 1,
			expectedIDs:   []string{"ses_def456"},
			description:   "按更新时间精确匹配",
		},
		{
			name: "combined filters",
			setupFilters: func() {
				filterID = ""
				filterTitle = "Git"
				filterCreated = 0
				filterUpdated = 0
				filterProjectID = "1f01524d"
				filterDirectory = "rl_mockgame"
			},
			expectedCount: 1,
			expectedIDs:   []string{"ses_def456"},
			description:   "组合多个过滤条件",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFilters()
			result := applyFieldFilter(testSessions)
			if len(result) != tt.expectedCount {
				t.Errorf("%s: 期望数量 %d, 实际 %d", tt.description, tt.expectedCount, len(result))
			}
			if tt.expectedIDs != nil {
				if len(result) != len(tt.expectedIDs) {
					t.Errorf("%s: 期望 ID 数量 %d, 实际 %d", tt.description, len(tt.expectedIDs), len(result))
				}
				for i, expectedID := range tt.expectedIDs {
					if i < len(result) && result[i].ID != expectedID {
						t.Errorf("%s: 期望 ID %s, 实际 %s", tt.description, expectedID, result[i].ID)
					}
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{"empty substr", "hello", "", true},
		{"empty string", "", "test", true},
		{"both empty", "", "", true},
		{"exact match", "hello", "hello", true},
		{"substring", "hello world", "world", true},
		{"case insensitive", "Hello World", "hello", true},
		{"case insensitive 2", "HELLO", "hello", true},
		{"no match", "hello", "world", false},
		{"chinese", "多帧融合", "融合", true},
		{"chinese no match", "多帧融合", "Git", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%q, %q) = %v, 期望 %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestSessionListCmdWithFilters(t *testing.T) {
	tests := []struct {
		name            string
		filterID        string
		filterTitle     string
		filterCreated   int64
		filterUpdated   int64
		filterProjectID string
		filterDirectory string
		mockResp        []byte
		expectedCount   int
		expectedIDs     []string
		description     string
	}{
		{
			name:            "filter by ID only",
			filterID:        "ses_abc",
			filterTitle:     "",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "",
			filterDirectory: "",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   0,
			description:     "按 ID 过滤，无匹配结果",
		},
		{
			name:            "filter by title only",
			filterID:        "",
			filterTitle:     "Test",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "",
			filterDirectory: "",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   2,
			expectedIDs:     []string{"session1", "session2"},
			description:     "按标题模糊查询",
		},
		{
			name:            "filter by projectID only",
			filterID:        "",
			filterTitle:     "",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "proj1",
			filterDirectory: "",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   1,
			expectedIDs:     []string{"session1"},
			description:     "按项目 ID 过滤",
		},
		{
			name:            "filter by directory only",
			filterID:        "",
			filterTitle:     "",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "",
			filterDirectory: "project1",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   1,
			expectedIDs:     []string{"session1"},
			description:     "按目录过滤",
		},
		{
			name:            "combined title and projectID",
			filterID:        "",
			filterTitle:     "Session 1",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "proj1",
			filterDirectory: "",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   1,
			expectedIDs:     []string{"session1"},
			description:     "组合标题和项目 ID 过滤",
		},
		{
			name:            "combined all filters",
			filterID:        "session1",
			filterTitle:     "Session 1",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "proj1",
			filterDirectory: "project1",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   1,
			expectedIDs:     []string{"session1"},
			description:     "组合所有过滤条件",
		},
		{
			name:            "no match combined filters",
			filterID:        "session1",
			filterTitle:     "Session 2",
			filterCreated:   0,
			filterUpdated:   0,
			filterProjectID: "proj1",
			filterDirectory: "project1",
			mockResp:        testutil.MockSessionsResponse(),
			expectedCount:   0,
			description:     "组合过滤无匹配结果",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置过滤参数
			filterID = tt.filterID
			filterTitle = tt.filterTitle
			filterCreated = tt.filterCreated
			filterUpdated = tt.filterUpdated
			filterProjectID = tt.filterProjectID
			filterDirectory = tt.filterDirectory

			// 解析模拟响应
			var sessions []types.Session
			if err := json.Unmarshal(tt.mockResp, &sessions); err != nil {
				t.Fatalf("Failed to unmarshal mock response: %v", err)
			}

			// 应用过滤
			result := applyFieldFilter(sessions)

			// 验证数量
			if len(result) != tt.expectedCount {
				t.Errorf("%s: 期望数量 %d, 实际 %d", tt.description, tt.expectedCount, len(result))
			}

			// 验证 ID
			if tt.expectedIDs != nil {
				if len(result) != len(tt.expectedIDs) {
					t.Errorf("%s: 期望 ID 数量 %d, 实际 %d", tt.description, len(tt.expectedIDs), len(result))
				}
				for i, expectedID := range tt.expectedIDs {
					if i < len(result) && result[i].ID != expectedID {
						t.Errorf("%s: 期望 ID %s, 实际 %s", tt.description, expectedID, result[i].ID)
					}
				}
			}
		})
	}
}

func TestSessionListCmdIntegration(t *testing.T) {
	// 测试 listCmd 与状态过滤和字段过滤的组合
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			if path == "/session" {
				return testutil.MockSessionsResponse(), nil
			}
			if path == "/session/status" {
				return testutil.MockSessionStatusResponse(), nil
			}
			return nil, nil
		},
	}

	// 测试场景 1: 无过滤
	func() {
		filterID = ""
		filterTitle = ""
		filterCreated = 0
		filterUpdated = 0
		filterProjectID = ""
		filterDirectory = ""
		statusFilter = ""
		runningOnly = false
		limit = 0
		offset = 0

		resp, err := mock.Get(context.Background(), "/session")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		var sessions []types.Session
		if err := json.Unmarshal(resp, &sessions); err != nil {
			t.Errorf("Failed to unmarshal: %v", err)
			return
		}

		result := applyFieldFilter(sessions)
		if len(result) != 2 {
			t.Errorf("无过滤时应返回 2 个会话，实际 %d", len(result))
		}
	}()

	// 测试场景 2: 状态过滤 + 字段过滤
	func() {
		filterID = ""
		filterTitle = "Session 1"
		filterCreated = 0
		filterUpdated = 0
		filterProjectID = ""
		filterDirectory = ""
		statusFilter = ""
		runningOnly = false

		resp, err := mock.Get(context.Background(), "/session")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		var sessions []types.Session
		if err := json.Unmarshal(resp, &sessions); err != nil {
			t.Errorf("Failed to unmarshal: %v", err)
			return
		}

		result := applyFieldFilter(sessions)
		if len(result) != 1 {
			t.Errorf("标题过滤后应返回 1 个会话，实际 %d", len(result))
		}
		if len(result) > 0 && result[0].ID != "session1" {
			t.Errorf("期望 session1，实际 %s", result[0].ID)
		}
	}()
}

func TestSessionListCmdSortAndPagination(t *testing.T) {
	testSessions := []types.Session{
		{ID: "s1", Title: "First", Time: types.SessionTime{Created: 1000, Updated: 3000}},
		{ID: "s2", Title: "Second", Time: types.SessionTime{Created: 2000, Updated: 2000}},
		{ID: "s3", Title: "Third", Time: types.SessionTime{Created: 3000, Updated: 1000}},
	}

	tests := []struct {
		name        string
		sortBy      string
		sortOrder   string
		limit       int
		offset      int
		expectedIDs []string
		description string
	}{
		{
			name:        "sort by updated desc",
			sortBy:      "updated",
			sortOrder:   "desc",
			limit:       0,
			offset:      0,
			expectedIDs: []string{"s1", "s2", "s3"},
			description: "按更新时间降序",
		},
		{
			name:        "sort by updated asc",
			sortBy:      "updated",
			sortOrder:   "asc",
			limit:       0,
			offset:      0,
			expectedIDs: []string{"s3", "s2", "s1"},
			description: "按更新时间升序",
		},
		{
			name:        "sort by created desc",
			sortBy:      "created",
			sortOrder:   "desc",
			limit:       0,
			offset:      0,
			expectedIDs: []string{"s3", "s2", "s1"},
			description: "按创建时间降序",
		},
		{
			name:        "sort by created asc",
			sortBy:      "created",
			sortOrder:   "asc",
			limit:       0,
			offset:      0,
			expectedIDs: []string{"s1", "s2", "s3"},
			description: "按创建时间升序",
		},
		{
			name:        "pagination limit 2",
			sortBy:      "updated",
			sortOrder:   "desc",
			limit:       2,
			offset:      0,
			expectedIDs: []string{"s1", "s2"},
			description: "限制返回 2 条",
		},
		{
			name:        "pagination offset 1 limit 2",
			sortBy:      "updated",
			sortOrder:   "desc",
			limit:       2,
			offset:      1,
			expectedIDs: []string{"s2", "s3"},
			description: "偏移 1 条，限制 2 条",
		},
		{
			name:        "pagination offset beyond length",
			sortBy:      "updated",
			sortOrder:   "desc",
			limit:       2,
			offset:      10,
			expectedIDs: []string{},
			description: "偏移超出长度",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制测试数据
			sessions := make([]types.Session, len(testSessions))
			copy(sessions, testSessions)

			// 设置排序参数
			sortBy = tt.sortBy
			sortOrder = tt.sortOrder
			limit = tt.limit
			offset = tt.offset

			// 应用排序
			sort.Slice(sessions, func(i, j int) bool {
				var less bool
				switch sortBy {
				case "created":
					less = sessions[i].Time.Created < sessions[j].Time.Created
				case "updated", "":
					less = sessions[i].Time.Updated < sessions[j].Time.Updated
				default:
					less = sessions[i].Time.Updated < sessions[j].Time.Updated
				}
				if sortOrder == "desc" {
					return !less
				}
				return less
			})

			// 应用分页
			if limit > 0 {
				start := offset
				if start > len(sessions) {
					start = len(sessions)
				}
				end := start + limit
				if end > len(sessions) {
					end = len(sessions)
				}
				sessions = sessions[start:end]
			}

			// 验证结果
			if len(sessions) != len(tt.expectedIDs) {
				t.Errorf("%s: 期望数量 %d, 实际 %d", tt.description, len(tt.expectedIDs), len(sessions))
			}
			for i, expectedID := range tt.expectedIDs {
				if i < len(sessions) && sessions[i].ID != expectedID {
					t.Errorf("%s: 期望 ID %s, 实际 %s", tt.description, expectedID, sessions[i].ID)
				}
			}
		})
	}
}

func TestSessionListCmdStatusFilter(t *testing.T) {
	testSessions := []types.Session{
		{ID: "s1", Title: "Session 1"},
		{ID: "s2", Title: "Session 2"},
		{ID: "s3", Title: "Session 3"},
	}

	statusMap := map[string]types.SessionStatus{
		"s1": {Status: "idle", IsReady: true, IsWorking: false},
		"s2": {Status: "working", IsReady: true, IsWorking: true, MessageID: "msg1"},
		"s3": {Status: "error", IsReady: true, IsWorking: false},
	}

	tests := []struct {
		name          string
		runningOnly   bool
		statusFilter  string
		expectedCount int
		expectedIDs   []string
		description   string
	}{
		{
			name:          "running only",
			runningOnly:   true,
			statusFilter:  "",
			expectedCount: 1,
			expectedIDs:   []string{"s2"},
			description:   "只显示运行中的会话",
		},
		{
			name:          "status running",
			runningOnly:   false,
			statusFilter:  "running",
			expectedCount: 1,
			expectedIDs:   []string{"s2"},
			description:   "过滤 running 状态",
		},
		{
			name:          "status completed",
			runningOnly:   false,
			statusFilter:  "completed",
			expectedCount: 2,
			expectedIDs:   []string{"s1", "s3"},
			description:   "过滤 completed 状态",
		},
		{
			name:          "status error",
			runningOnly:   false,
			statusFilter:  "error",
			expectedCount: 1,
			expectedIDs:   []string{"s3"},
			description:   "过滤 error 状态",
		},
		{
			name:          "no status filter",
			runningOnly:   false,
			statusFilter:  "",
			expectedCount: 3,
			expectedIDs:   []string{"s1", "s2", "s3"},
			description:   "无状态过滤",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filteredSessions []types.Session

			if tt.runningOnly {
				for _, session := range testSessions {
					if status, exists := statusMap[session.ID]; exists && status.IsWorking {
						filteredSessions = append(filteredSessions, session)
					}
				}
			} else if tt.statusFilter != "" {
				for _, session := range testSessions {
					status, exists := statusMap[session.ID]
					match := false
					switch tt.statusFilter {
					case "running":
						match = exists && status.IsWorking
					case "completed", "idle":
						if !exists {
							match = true
						} else {
							match = !status.IsWorking
						}
					case "error":
						match = exists && status.Status == "error"
					case "aborted":
						match = exists && status.Status == "aborted"
					default:
						match = true
					}
					if match {
						filteredSessions = append(filteredSessions, session)
					}
				}
			} else {
				filteredSessions = testSessions
			}

			if len(filteredSessions) != tt.expectedCount {
				t.Errorf("%s: 期望数量 %d, 实际 %d", tt.description, tt.expectedCount, len(filteredSessions))
			}
			if tt.expectedIDs != nil {
				for i, expectedID := range tt.expectedIDs {
					if i < len(filteredSessions) && filteredSessions[i].ID != expectedID {
						t.Errorf("%s: 期望 ID %s, 实际 %s", tt.description, expectedID, filteredSessions[i].ID)
					}
				}
			}
		})
	}
}
