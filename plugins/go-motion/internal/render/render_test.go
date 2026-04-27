package render

import "testing"

func TestBuildRenderPlan_HasCaptureAndEncodeCommands(t *testing.T) {
	plan := BuildRenderPlan("C:/plugin", "C:/job")
	if len(plan.Commands) < 2 {
		t.Fatalf("expected render and encode commands, got %d", len(plan.Commands))
	}
	if plan.Commands[0].Name != "capture-frames" {
		t.Fatalf("first command = %q, want %q", plan.Commands[0].Name, "capture-frames")
	}
	if plan.Commands[1].Name != "encode-mp4" {
		t.Fatalf("second command = %q, want %q", plan.Commands[1].Name, "encode-mp4")
	}
}
