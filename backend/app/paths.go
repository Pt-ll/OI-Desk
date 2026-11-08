package app

import (
	"fmt"
	"os"
	"path/filepath"
)

type Paths struct {
	DataDir   string
	Database  string
	PluginDir string
	CacheDir  string
	LogDir    string
	BackupDir string
}

func ResolvePaths(appName string) (Paths, error) {
	if appName == "" {
		return Paths{}, fmt.Errorf("application name is empty")
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return Paths{}, fmt.Errorf("resolve user config directory: %w", err)
	}
	data := filepath.Join(base, appName)
	cacheBase, err := os.UserCacheDir()
	if err != nil {
		return Paths{}, fmt.Errorf("resolve user cache directory: %w", err)
	}
	return Paths{
		DataDir:   data,
		Database:  filepath.Join(data, "core", "oidesk.db"),
		PluginDir: filepath.Join(data, "plugins"),
		CacheDir:  filepath.Join(cacheBase, appName),
		LogDir:    filepath.Join(data, "logs"),
		BackupDir: filepath.Join(data, "backups"),
	}, nil
}

func (p Paths) Ensure() error {
	for _, dir := range []string{p.DataDir, filepath.Dir(p.Database), p.PluginDir, p.CacheDir, p.LogDir, p.BackupDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("create application directory %q: %w", dir, err)
		}
	}
	return nil
}
