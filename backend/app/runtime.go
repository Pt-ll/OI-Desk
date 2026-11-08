package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Pt-ll/OI-Desk/backend/database"
)

type Runtime struct {
	Paths    Paths
	DB       *sql.DB
	Settings *database.SettingsStore
	Logger   *slog.Logger
	logFile  *os.File
}

func Bootstrap(ctx context.Context, appName, dataDirOverride string) (*Runtime, error) {
	paths, err := ResolvePaths(appName)
	if err != nil {
		return nil, err
	}
	if dataDirOverride != "" {
		paths.DataDir = dataDirOverride
		paths.Database = filepath.Join(dataDirOverride, "core", "oidesk.db")
		paths.PluginDir = filepath.Join(dataDirOverride, "plugins")
		paths.CacheDir = filepath.Join(dataDirOverride, "cache")
		paths.LogDir = filepath.Join(dataDirOverride, "logs")
		paths.BackupDir = filepath.Join(dataDirOverride, "backups")
	}
	if err := paths.Ensure(); err != nil {
		return nil, err
	}
	db, err := database.Open(ctx, paths.Database)
	if err != nil {
		return nil, err
	}
	closeWithError := func(err error) (*Runtime, error) {
		_ = db.Close()
		return nil, err
	}
	if err := database.ApplyMigrations(ctx, db); err != nil {
		return closeWithError(err)
	}
	logFile, err := os.OpenFile(filepath.Join(paths.LogDir, "oidesk.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return closeWithError(fmt.Errorf("open application log: %w", err))
	}
	return &Runtime{Paths: paths, DB: db, Settings: database.NewSettingsStore(db), Logger: NewLogger(logFile), logFile: logFile}, nil
}

func (r *Runtime) Close() error {
	if r == nil {
		return nil
	}
	var firstErr error
	if r.logFile != nil {
		if err := r.logFile.Close(); err != nil {
			firstErr = err
		}
	}
	if r.DB != nil {
		if err := r.DB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
