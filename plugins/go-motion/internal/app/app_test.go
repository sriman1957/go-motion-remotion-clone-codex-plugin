package app

import "testing"

func TestNewDefaultConfig_HasPluginName(t *testing.T) {
	cfg := NewDefaultConfig()
	if cfg.PluginName != "go-motion" {
		t.Fatalf("PluginName = %q, want %q", cfg.PluginName, "go-motion")
	}
}
