package app

import (
	"context"
	"path/filepath"
	"testing"
)

func TestBootstrapCreatesFoundationAndCloses(t *testing.T) {
	dataDir := filepath.Join(t.TempDir(), "data")
	runtime, err := Bootstrap(context.Background(), "OI-Desk", dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.Paths.Database != filepath.Join(dataDir, "core", "oidesk.db") {
		t.Fatalf("database path = %q", runtime.Paths.Database)
	}
	if err := runtime.Settings.Set(context.Background(), "theme", `"dark"`, "string"); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
}
