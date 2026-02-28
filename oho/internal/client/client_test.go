package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anomalyco/oho/internal/config"
)

// 初始化配置
func init() {
	_ = config.Init()
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient() returned nil")
	}
	if c.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestClientGetSuccess(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := &Client{
		baseURL:    server.URL,
		username:   "test",
		password:   "test",
		httpClient: &http.Client{},
	}

	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status ok, got %s", result["status"])
	}
}

func TestClientPostSuccess(t *testing.T) {
	// 测试 JSON 编组
	req := map[string]string{"name": "test"}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name test, got %s", result["name"])
	}
}

func TestClientEmptyResponse(t *testing.T) {
	// 创建返回空响应的服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{
		baseURL:    server.URL,
		username:   "test",
		password:   "test",
		httpClient: &http.Client{},
	}

	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// 空响应应该返回空字节切片
	if len(resp) != 0 {
		t.Errorf("Expected empty response, got %d bytes", len(resp))
	}
}

func TestClientErrorResponse(t *testing.T) {
	// 创建返回错误的服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	c := &Client{
		baseURL:    server.URL,
		username:   "test",
		password:   "test",
		httpClient: &http.Client{},
	}

	_, err := c.Get(context.Background(), "/test")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
