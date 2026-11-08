package database

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

func TestOpenMigrateAndSettingsAreIdempotent(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "core.db")
	db, err := Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	var migrationCount int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&migrationCount); err != nil {
		t.Fatal(err)
	}
	if migrationCount != 1 {
		t.Fatalf("migration count = %d", migrationCount)
	}
	settings := NewSettingsStore(db)
	if err := settings.Set(ctx, "theme", `"dark"`, "string"); err != nil {
		t.Fatal(err)
	}
	if err := settings.Set(ctx, "theme", `"light"`, "string"); err != nil {
		t.Fatal(err)
	}
	got, err := settings.Get(ctx, "theme")
	if err != nil {
		t.Fatal(err)
	}
	if got.ValueJSON != `"light"` || got.ValueType != "string" {
		t.Fatalf("setting = %+v", got)
	}
	if err := settings.Delete(ctx, "theme"); err != nil {
		t.Fatal(err)
	}
	if _, err := settings.Get(ctx, "theme"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected missing setting, got %v", err)
	}
}

func TestSettingsRejectInvalidValues(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, filepath.Join(t.TempDir(), "core.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	settings := NewSettingsStore(db)
	if err := settings.Set(ctx, "bad", "not-json", "string"); err == nil {
		t.Fatal("invalid JSON accepted")
	}
	if err := settings.Set(ctx, "bad", `"ok"`, "duration"); err == nil {
		t.Fatal("invalid value type accepted")
	}
}
