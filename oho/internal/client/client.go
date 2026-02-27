package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/anomalyco/oho/internal/config"
)

// Client OpenCode API 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
}

// NewClient 创建新的 API 客户端
func NewClient() *Client {
	cfg := config.Get()
	
	return &Client{
		baseURL:  config.GetBaseURL(),
		username: cfg.Username,
		password: cfg.Password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request 发送 HTTP 请求
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败：%w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}

	// 设置认证
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	// 设置请求头
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败：%w", err)
	}

	// 检查状态码
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API 错误 [%d]: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// RequestWithQuery 发送带查询参数的 HTTP 请求
func (c *Client) RequestWithQuery(ctx context.Context, method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	// 构建查询参数
	if len(queryParams) > 0 {
		u, err := url.Parse(c.baseURL + path)
		if err != nil {
			return nil, err
		}
		
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		path = u.String()
	}

	return c.Request(ctx, method, path, body)
}

// Get 发送 GET 请求
func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	return c.Request(ctx, http.MethodGet, path, nil)
}

// GetWithQuery 发送带查询参数的 GET 请求
func (c *Client) GetWithQuery(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
	return c.RequestWithQuery(ctx, http.MethodGet, path, queryParams, nil)
}

// Post 发送 POST 请求
func (c *Client) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.Request(ctx, http.MethodPost, path, body)
}

// Put 发送 PUT 请求
func (c *Client) Put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.Request(ctx, http.MethodPut, path, body)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.Request(ctx, http.MethodPatch, path, body)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(ctx context.Context, path string) ([]byte, error) {
	return c.Request(ctx, http.MethodDelete, path, nil)
}

// SSEStream 服务器发送事件流
func (c *Client) SSEStream(ctx context.Context, path string) (<-chan []byte, <-chan error, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, nil, err
	}

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("SSE 错误 [%d]: %s", resp.StatusCode, string(body))
	}

	eventChan := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		defer close(eventChan)
		defer close(errChan)
		defer resp.Body.Close()

		buf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := resp.Body.Read(buf)
				if n > 0 {
					data := make([]byte, n)
					copy(data, buf[:n])
					eventChan <- data
				}
				if err != nil {
					if err != io.EOF {
						errChan <- err
					}
					return
				}
			}
		}
	}()

	return eventChan, errChan, nil
}
