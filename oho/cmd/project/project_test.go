package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/testutil"
	"github.com/anomalyco/oho/internal/types"
)

func TestProjectListCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockProjectsResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/project")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var projects []types.Project
	if err := json.Unmarshal(resp, &projects); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}

func TestProjectCurrentCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			project := types.Project{
				ID:   "proj1",
				Name: "Current Project",
				Path: "/home/user/project",
				Vcs:  "git",
			}
			return json.Marshal(project)
		},
	}

	resp, err := mock.Get(context.Background(), "/project/current")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var project types.Project
	if err := json.Unmarshal(resp, &project); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if project.ID != "proj1" {
		t.Errorf("Expected project ID 'proj1', got %s", project.ID)
	}
}

func TestPathCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockPathResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/project/path")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var path types.Path
	if err := json.Unmarshal(resp, &path); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if path.Current == "" {
		t.Error("Expected current path but got empty")
	}
}

func TestVCSCmd(t *testing.T) {
	mock := &client.MockClient{
		GetFunc: func(ctx context.Context, path string) ([]byte, error) {
			return testutil.MockVCSResponse(), nil
		},
	}

	resp, err := mock.Get(context.Background(), "/project/vcs")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var vcs types.VcsInfo
	if err := json.Unmarshal(resp, &vcs); err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if vcs.Type != "git" {
		t.Errorf("Expected VCS type 'git', got %s", vcs.Type)
	}
}

func TestInstanceDisposeCmd(t *testing.T) {
	mock := &client.MockClient{
		PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
			return testutil.MockBoolResponse(true), nil
		},
	}

	resp, err := mock.Post(context.Background(), "/instance/dispose", nil)
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
