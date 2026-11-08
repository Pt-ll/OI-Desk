package main

import (
	"context"
	"path/filepath"
	"testing"

	foundation "github.com/Pt-ll/OI-Desk/backend/app"
)

func TestSettingBindingsRoundTrip(t *testing.T) {
	runtime, err := foundation.Bootstrap(context.Background(), "OI-Desk", filepath.Join(t.TempDir(), "data"))
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	app := &App{ctx: context.Background(), foundation: runtime}
	if err := app.SetSetting("theme", `"light"`, "string"); err != nil {
		t.Fatal(err)
	}
	setting, err := app.GetSetting("theme")
	if err != nil {
		t.Fatal(err)
	}
	if setting.ValueJSON != `"light"` || setting.ValueType != "string" {
		t.Fatalf("setting = %+v", setting)
	}
}
