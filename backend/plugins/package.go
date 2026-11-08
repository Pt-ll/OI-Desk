package plugins

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"
)

const (
	maxPackageFiles        = 2048
	maxPackageBytes  int64 = 64 << 20
	maxManifestBytes       = 1 << 20
)

type PackageInfo struct {
	Manifest Manifest
	Files    int
	Bytes    int64
}

func InspectZip(filename string, coreVersion string) (PackageInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return PackageInfo{}, fmt.Errorf("open plugin package: %w", err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return PackageInfo{}, fmt.Errorf("stat plugin package: %w", err)
	}
	if stat.Size() <= 0 {
		return PackageInfo{}, fmt.Errorf("plugin package is empty")
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return PackageInfo{}, fmt.Errorf("read plugin zip: %w", err)
	}
	if len(reader.File) == 0 || len(reader.File) > maxPackageFiles {
		return PackageInfo{}, fmt.Errorf("plugin package contains an invalid number of files")
	}
	manifestData, files, total, err := collectZip(reader)
	if err != nil {
		return PackageInfo{}, err
	}
	manifest, err := ParseAndValidateManifest(manifestData, files, coreVersion)
	if err != nil {
		return PackageInfo{}, err
	}
	return PackageInfo{Manifest: manifest, Files: len(reader.File), Bytes: total}, nil
}

func collectZip(reader *zip.Reader) ([]byte, fs.FS, int64, error) {
	entries := make(map[string][]byte, len(reader.File))
	var manifest []byte
	var total int64
	for _, file := range reader.File {
		name := path.Clean(strings.ReplaceAll(file.Name, "\\", "/"))
		if name == "." || path.IsAbs(file.Name) || strings.ContainsRune(file.Name, '\\') || name != strings.ReplaceAll(file.Name, "\\", "/") || strings.HasPrefix(name, "../") || name == ".." {
			return nil, nil, 0, fmt.Errorf("unsafe plugin archive path %q", file.Name)
		}
		if file.Mode()&os.ModeSymlink != 0 {
			return nil, nil, 0, fmt.Errorf("symlink is not allowed in plugin archive: %q", name)
		}
		if file.FileInfo().IsDir() {
			continue
		}
		if _, exists := entries[name]; exists {
			return nil, nil, 0, fmt.Errorf("duplicate plugin archive path %q", name)
		}
		if file.UncompressedSize64 > uint64(maxPackageBytes) {
			return nil, nil, 0, fmt.Errorf("plugin file %q is too large", name)
		}
		if total+int64(file.UncompressedSize64) > maxPackageBytes {
			return nil, nil, 0, fmt.Errorf("plugin package exceeds size limit")
		}
		reader, err := file.Open()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("open plugin file %q: %w", name, err)
		}
		data, err := io.ReadAll(io.LimitReader(reader, maxPackageBytes+1))
		_ = reader.Close()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("read plugin file %q: %w", name, err)
		}
		if int64(len(data)) > maxPackageBytes {
			return nil, nil, 0, fmt.Errorf("plugin package exceeds size limit")
		}
		entries[name] = data
		total += int64(len(data))
		if name == "manifest.json" {
			if len(data) > maxManifestBytes {
				return nil, nil, 0, fmt.Errorf("manifest is too large")
			}
			manifest = data
		}
	}
	if manifest == nil {
		return nil, nil, 0, fmt.Errorf("plugin package has no manifest.json")
	}
	return manifest, mapFS(entries), total, nil
}

type mapFS map[string][]byte

func (m mapFS) Open(name string) (fs.File, error) {
	data, ok := m[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &memoryFile{name: name, reader: bytes.NewReader(data)}, nil
}

type memoryFile struct {
	name   string
	reader *bytes.Reader
}

func (f *memoryFile) Stat() (fs.FileInfo, error) {
	return memoryInfo{name: path.Base(f.name), size: int64(f.reader.Len())}, nil
}
func (f *memoryFile) Read(p []byte) (int, error) { return f.reader.Read(p) }
func (f *memoryFile) Close() error               { return nil }

type memoryInfo struct {
	name string
	size int64
}

func (i memoryInfo) Name() string           { return i.name }
func (i memoryInfo) Size() int64            { return i.size }
func (i memoryInfo) Mode() fs.FileMode      { return 0 }
func (i memoryInfo) ModTime() (t time.Time) { return time.Time{} }
func (i memoryInfo) IsDir() bool            { return false }
func (i memoryInfo) Sys() any               { return nil }
