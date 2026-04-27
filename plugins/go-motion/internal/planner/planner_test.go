package planner

import "testing"

func TestPlanPrompt_BuildsScenesFromPrompt(t *testing.T) {
	plan := PlanPrompt("Create a clean product promo for a note-taking app")
	if len(plan.Scenes) == 0 {
		t.Fatal("expected at least one scene")
	}
	if plan.Title == "" {
		t.Fatal("expected title to be populated")
	}
	if plan.FPS != 24 {
		t.Fatalf("FPS = %d, want 24", plan.FPS)
	}
}
