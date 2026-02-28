package client

import (
	"context"
	"fmt"
)

// APIError API 错误
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error [%d]: %s", e.StatusCode, e.Message)
}

// MockClient implements ClientInterface for testing
type MockClient struct {
	GetFunc          func(ctx context.Context, path string) ([]byte, error)
	GetWithQueryFunc func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error)
	PostFunc         func(ctx context.Context, path string, body interface{}) ([]byte, error)
	PutFunc          func(ctx context.Context, path string, body interface{}) ([]byte, error)
	PatchFunc        func(ctx context.Context, path string, body interface{}) ([]byte, error)
	DeleteFunc       func(ctx context.Context, path string) ([]byte, error)
	SSEStreamFunc    func(ctx context.Context, path string) (<-chan []byte, <-chan error, error)
}

func (m *MockClient) Get(ctx context.Context, path string) ([]byte, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, path)
	}
	return nil, nil
}

func (m *MockClient) GetWithQuery(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
	if m.GetWithQueryFunc != nil {
		return m.GetWithQueryFunc(ctx, path, queryParams)
	}
	return nil, nil
}

func (m *MockClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	if m.PostFunc != nil {
		return m.PostFunc(ctx, path, body)
	}
	return nil, nil
}

func (m *MockClient) Put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, path, body)
	}
	return nil, nil
}

func (m *MockClient) Patch(ctx context.Context, path string, body interface{}) ([]byte, error) {
	if m.PatchFunc != nil {
		return m.PatchFunc(ctx, path, body)
	}
	return nil, nil
}

func (m *MockClient) Delete(ctx context.Context, path string) ([]byte, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, path)
	}
	return nil, nil
}

func (m *MockClient) SSEStream(ctx context.Context, path string) (<-chan []byte, <-chan error, error) {
	if m.SSEStreamFunc != nil {
		return m.SSEStreamFunc(ctx, path)
	}
	return nil, nil, nil
}

// Ensure MockClient implements ClientInterface
var _ ClientInterface = (*MockClient)(nil)
