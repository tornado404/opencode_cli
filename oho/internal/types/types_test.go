package types

import (
	"encoding/json"
	"testing"
)

func TestSessionJSON(t *testing.T) {
	jsonData := `{
		"id": "test-session-123",
		"title": "Test Session",
		"parentId": "parent-456",
		"createdAt": 1699999999,
		"updatedAt": 1700000000,
		"model": "gpt-4",
		"agent": "default"
	}`

	var s Session
	err := json.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if s.ID != "test-session-123" {
		t.Errorf("ID = %v, want test-session-123", s.ID)
	}
	if s.Title != "Test Session" {
		t.Errorf("Title = %v, want Test Session", s.Title)
	}
	if s.ParentID != "parent-456" {
		t.Errorf("ParentID = %v, want parent-456", s.ParentID)
	}
	if s.Model != "gpt-4" {
		t.Errorf("Model = %v, want gpt-4", s.Model)
	}
}

func TestSessionStatusJSON(t *testing.T) {
	jsonData := `{
		"status": "working",
		"isReady": false,
		"isWorking": true,
		"messageId": "msg-789"
	}`

	var ss SessionStatus
	err := json.Unmarshal([]byte(jsonData), &ss)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if ss.Status != "working" {
		t.Errorf("Status = %v, want working", ss.Status)
	}
	if ss.IsReady != false {
		t.Errorf("IsReady = %v, want false", ss.IsReady)
	}
	if ss.IsWorking != true {
		t.Errorf("IsWorking = %v, want true", ss.IsWorking)
	}
	if ss.MessageID != "msg-789" {
		t.Errorf("MessageID = %v, want msg-789", ss.MessageID)
	}
}

func TestMessageJSON(t *testing.T) {
	jsonData := `{
		"id": "msg-123",
		"sessionId": "session-456",
		"role": "user",
		"createdAt": 1699999999,
		"content": "Hello, world!"
	}`

	var m Message
	err := json.Unmarshal([]byte(jsonData), &m)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if m.ID != "msg-123" {
		t.Errorf("ID = %v, want msg-123", m.ID)
	}
	if m.SessionID != "session-456" {
		t.Errorf("SessionID = %v, want session-456", m.SessionID)
	}
	if m.Role != "user" {
		t.Errorf("Role = %v, want user", m.Role)
	}
	if m.Content != "Hello, world!" {
		t.Errorf("Content = %v, want Hello, world!", m.Content)
	}
}

func TestPartJSON(t *testing.T) {
	jsonData := `{
		"type": "text",
		"data": "Some text content"
	}`

	var p Part
	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if p.Type != "text" {
		t.Errorf("Type = %v, want text", p.Type)
	}
}

func TestMessageWithPartsJSON(t *testing.T) {
	jsonData := `{
		"info": {
			"id": "msg-123",
			"sessionId": "session-456",
			"role": "assistant",
			"createdAt": 1699999999
		},
		"parts": [
			{"type": "text", "data": "Response text"}
		]
	}`

	var mwp MessageWithParts
	err := json.Unmarshal([]byte(jsonData), &mwp)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if mwp.Info.ID != "msg-123" {
		t.Errorf("Info.ID = %v, want msg-123", mwp.Info.ID)
	}
	if len(mwp.Parts) != 1 {
		t.Errorf("Parts length = %v, want 1", len(mwp.Parts))
	}
}

func TestConfigJSON(t *testing.T) {
	jsonData := `{
		"providers": {"openai": {"key": "test"}},
		"defaultModel": "gpt-4",
		"theme": "dark",
		"language": "en",
		"autoApprove": ["tool.use"],
		"maxTokens": 4000,
		"temperature": 0.7
	}`

	var c Config
	err := json.Unmarshal([]byte(jsonData), &c)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if c.DefaultModel != "gpt-4" {
		t.Errorf("DefaultModel = %v, want gpt-4", c.DefaultModel)
	}
	if c.Theme != "dark" {
		t.Errorf("Theme = %v, want dark", c.Theme)
	}
	if c.MaxTokens != 4000 {
		t.Errorf("MaxTokens = %v, want 4000", c.MaxTokens)
	}
	if c.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", c.Temperature)
	}
}

func TestProviderJSON(t *testing.T) {
	jsonData := `{
		"id": "openai",
		"name": "OpenAI",
		"baseURL": "https://api.openai.com",
		"models": ["gpt-4", "gpt-3.5-turbo"],
		"authType": "apiKey"
	}`

	var p Provider
	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if p.ID != "openai" {
		t.Errorf("ID = %v, want openai", p.ID)
	}
	if p.Name != "OpenAI" {
		t.Errorf("Name = %v, want OpenAI", p.Name)
	}
	if len(p.Models) != 2 {
		t.Errorf("Models length = %v, want 2", len(p.Models))
	}
}

func TestHealthResponseJSON(t *testing.T) {
	jsonData := `{
		"healthy": true,
		"version": "1.0.0"
	}`

	var hr HealthResponse
	err := json.Unmarshal([]byte(jsonData), &hr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if !hr.Healthy {
		t.Error("Healthy = false, want true")
	}
	if hr.Version != "1.0.0" {
		t.Errorf("Version = %v, want 1.0.0", hr.Version)
	}
}

func TestFileNodeJSON(t *testing.T) {
	jsonData := `{
		"name": "main.go",
		"path": "/path/to/main.go",
		"type": "file",
		"children": []
	}`

	var fn FileNode
	err := json.Unmarshal([]byte(jsonData), &fn)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if fn.Name != "main.go" {
		t.Errorf("Name = %v, want main.go", fn.Name)
	}
	if fn.Type != "file" {
		t.Errorf("Type = %v, want file", fn.Type)
	}
}

func TestFileJSON(t *testing.T) {
	jsonData := `{
		"path": "/path/to/file.go",
		"status": "modified"
	}`

	var f File
	err := json.Unmarshal([]byte(jsonData), &f)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if f.Path != "/path/to/file.go" {
		t.Errorf("Path = %v, want /path/to/file.go", f.Path)
	}
	if f.Status != "modified" {
		t.Errorf("Status = %v, want modified", f.Status)
	}
}

func TestSymbolJSON(t *testing.T) {
	jsonData := `{
		"name": "main",
		"kind": "function",
		"path": "/path/to/main.go",
		"line": 10,
		"column": 5,
		"container": "main"
	}`

	var s Symbol
	err := json.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if s.Name != "main" {
		t.Errorf("Name = %v, want main", s.Name)
	}
	if s.Kind != "function" {
		t.Errorf("Kind = %v, want function", s.Kind)
	}
	if s.Line != 10 {
		t.Errorf("Line = %v, want 10", s.Line)
	}
}
