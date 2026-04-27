package mcp

import "testing"

func TestListTools_ContainsGenerateVideo(t *testing.T) {
	srv := NewServer("C:/plugin", t.TempDir())
	tools := srv.ListTools()
	if len(tools) == 0 {
		t.Fatal("expected tools")
	}

	found := false
	for _, tool := range tools {
		if tool.Name == "generate_video" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected generate_video tool")
	}
}

func TestDispatchRequest_ToolsList(t *testing.T) {
	srv := NewServer("C:/plugin", t.TempDir())
	resp, shouldReply, err := srv.DispatchRequest("tools/list", nil)
	if err != nil {
		t.Fatalf("DispatchRequest() error = %v", err)
	}
	if !shouldReply {
		t.Fatal("expected reply for tools/list")
	}
	tools, ok := resp["tools"].([]Tool)
	if !ok {
		t.Fatalf("expected tools response, got %#v", resp)
	}
	if len(tools) == 0 {
		t.Fatal("expected at least one tool")
	}
}
