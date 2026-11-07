package phase0

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSQLiteDriverTransactionBackupAndReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "probe.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	for _, query := range []string{
		"PRAGMA foreign_keys=ON", "PRAGMA journal_mode=WAL", "PRAGMA busy_timeout=5000",
		"CREATE TABLE schema_migrations(version INTEGER PRIMARY KEY)",
		"INSERT INTO schema_migrations VALUES(1)",
		"CREATE TABLE parent(id INTEGER PRIMARY KEY)",
		"CREATE TABLE child(id INTEGER PRIMARY KEY, parent_id INTEGER REFERENCES parent(id))",
		"INSERT INTO parent VALUES(1)",
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
	}
	if _, err := db.Exec("INSERT INTO child VALUES(1,999)"); err == nil {
		t.Fatal("foreign key violation accepted")
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.Exec("INSERT INTO child VALUES(1,1)"); err != nil {
		t.Fatal(err)
	}
	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	var count int
	if err = db.QueryRow("SELECT count(*) FROM child").Scan(&count); err != nil || count != 0 {
		t.Fatalf("rollback: %d, %v", count, err)
	}
	if _, err = db.Exec("INSERT INTO child VALUES(1,1)"); err != nil {
		t.Fatal(err)
	}
	backup := filepath.Join(dir, "backup.db")
	if _, err = db.Exec("VACUUM INTO '" + strings.ReplaceAll(backup, "'", "''") + "'"); err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{path, backup} {
		copy, err := sql.Open("sqlite", file)
		if err != nil {
			t.Fatal(err)
		}
		err = copy.QueryRow("SELECT count(*) FROM child").Scan(&count)
		_ = copy.Close()
		if err != nil || count != 1 {
			t.Fatalf("restore %s: count=%d err=%v", file, count, err)
		}
	}
}
