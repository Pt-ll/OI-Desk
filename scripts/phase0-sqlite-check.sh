#!/bin/sh
set -eu
probe_dir="$(mktemp -d "${TMPDIR:-/tmp}/oidesk-sqlite.XXXXXX")"
trap 'rm -f "$probe_dir/core.db" "$probe_dir/core.db-wal" "$probe_dir/core.db-shm" "$probe_dir/backup.db" "$probe_dir/backup.db-wal" "$probe_dir/backup.db-shm"; rmdir "$probe_dir"' EXIT HUP INT TERM
sqlite3 "$probe_dir/core.db" <<'SQL'
PRAGMA foreign_keys=ON;
PRAGMA journal_mode=WAL;
CREATE TABLE parent(id INTEGER PRIMARY KEY);
CREATE TABLE child(id INTEGER PRIMARY KEY, parent_id INTEGER NOT NULL REFERENCES parent(id));
INSERT INTO parent VALUES(1);
BEGIN;
INSERT INTO child VALUES(1,1);
COMMIT;
SELECT sqlite_version();
SELECT foreign_keys FROM pragma_foreign_keys;
SELECT journal_mode FROM pragma_journal_mode;
SELECT count(*) FROM child;
SQL
sqlite3 "$probe_dir/core.db" ".backup '$probe_dir/backup.db'"
test "$(sqlite3 "$probe_dir/backup.db" 'SELECT count(*) FROM child;')" = 1
echo 'backup restore: PASS'
