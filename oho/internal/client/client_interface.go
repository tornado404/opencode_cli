package client

import (
	"context"
)

// ClientInterface 定义客户端接口，便于测试
type ClientInterface interface {
	Get(ctx context.Context, path string) ([]byte, error)
	GetWithQuery(ctx context.Context, path string, queryParams map[string]string) ([]byte, error)
	Post(ctx context.Context, path string, body interface{}) ([]byte, error)
	Put(ctx context.Context, path string, body interface{}) ([]byte, error)
	Patch(ctx context.Context, path string, body interface{}) ([]byte, error)
	Delete(ctx context.Context, path string) ([]byte, error)
	SSEStream(ctx context.Context, path string) (<-chan []byte, <-chan error, error)
}

// 确保 Client 实现 ClientInterface
var _ ClientInterface = (*Client)(nil)
