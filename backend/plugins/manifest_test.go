package plugins

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func validManifest() []byte {
	value := map[string]any{
		"id": "org.oidesk.hello", "name": "Hello", "version": "1.0.0", "description": "probe",
		"author": "OI-Desk", "entry": "index.html", "minCoreVersion": CoreVersion,
		"permissions": []string{"theme", "storage", "theme"},
		"ai":          map[string]any{"enabled": false, "capabilities": []string{}},
	}
	data, _ := json.Marshal(value)
	return data
}

func TestManifestValidationNormalizesPermissions(t *testing.T) {
	manifest, err := ParseAndValidateManifest(validManifest(), fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte("ok")}}, CoreVersion)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Permissions) != 2 || manifest.Permissions[0] != "storage" || manifest.Permissions[1] != "theme" {
		t.Fatalf("permissions = %#v", manifest.Permissions)
	}
}

func TestManifestValidationRejectsUnsafeOrUnauthorizedInput(t *testing.T) {
	for name, raw := range map[string][]byte{
		"bad id":         []byte(`{"id":"hello","name":"Hello","version":"1.0.0","entry":"index.html","minCoreVersion":"1.0.0","permissions":[]}`),
		"bad entry":      []byte(`{"id":"org.oidesk.hello","name":"Hello","version":"1.0.0","entry":"../index.html","minCoreVersion":"1.0.0","permissions":[]}`),
		"bad permission": []byte(`{"id":"org.oidesk.hello","name":"Hello","version":"1.0.0","entry":"index.html","minCoreVersion":"1.0.0","permissions":["filesystem.raw"]}`),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseAndValidateManifest(raw, fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte("ok")}}, CoreVersion); err == nil {
				t.Fatal("invalid manifest accepted")
			}
		})
	}
}

func TestInspectZipRejectsTraversalAndAcceptsValidPackage(t *testing.T) {
	valid := filepath.Join(t.TempDir(), "valid.zip")
	writeZip(t, valid, map[string]string{"manifest.json": string(validManifest()), "index.html": "hello"})
	info, err := InspectZip(valid, CoreVersion)
	if err != nil {
		t.Fatal(err)
	}
	if info.Manifest.ID != "org.oidesk.hello" || info.Files != 2 {
		t.Fatalf("package = %+v", info)
	}
	unsafe := filepath.Join(t.TempDir(), "unsafe.zip")
	writeZip(t, unsafe, map[string]string{"manifest.json": string(validManifest()), "../escape.js": "bad"})
	if _, err := InspectZip(unsafe, CoreVersion); err == nil {
		t.Fatal("path traversal accepted")
	}
}

func writeZip(t *testing.T, filename string, files map[string]string) {
	t.Helper()
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	for name, contents := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(contents)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
