package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var pluginIDPattern = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*(?:\.[a-z0-9]+(?:[._-][a-z0-9]+)*)+$`)
var semverPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

const CoreVersion = "1.0.0"

var AllowedPermissions = map[string]struct{}{
	"storage": {}, "clipboard": {}, "notification": {}, "theme": {},
	"contest.read": {}, "problem.read": {}, "submission.read": {}, "algorithm.read": {},
	"ai.chat": {}, "ai.summarize": {}, "ai.embedding": {}, "network": {},
	"file.export": {}, "file.import": {},
}

var AllowedAICapabilities = map[string]struct{}{"chat": {}, "summarize": {}, "embedding": {}}

type Manifest struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Version        string     `json:"version"`
	Description    string     `json:"description"`
	Author         string     `json:"author"`
	Entry          string     `json:"entry"`
	Icon           string     `json:"icon,omitempty"`
	MinCoreVersion string     `json:"minCoreVersion"`
	Permissions    []string   `json:"permissions"`
	AI             AIManifest `json:"ai"`
}

type AIManifest struct {
	Enabled      bool     `json:"enabled"`
	Capabilities []string `json:"capabilities"`
	DefaultModel *string  `json:"defaultModel"`
}

func ParseAndValidateManifest(data []byte, pluginFS fs.FS, coreVersion string) (Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("parse manifest: %w", err)
	}
	if err := validateManifest(manifest, pluginFS, coreVersion); err != nil {
		return Manifest{}, err
	}
	manifest.Permissions = uniqueSorted(manifest.Permissions)
	manifest.AI.Capabilities = uniqueSorted(manifest.AI.Capabilities)
	return manifest, nil
}

func validateManifest(manifest Manifest, pluginFS fs.FS, coreVersion string) error {
	if !pluginIDPattern.MatchString(manifest.ID) {
		return fmt.Errorf("invalid plugin id %q", manifest.ID)
	}
	if len(manifest.ID) > 128 || len(manifest.Name) == 0 || len(manifest.Name) > 128 {
		return errors.New("plugin id/name is empty or too long")
	}
	if !semverPattern.MatchString(manifest.Version) {
		return fmt.Errorf("invalid plugin version %q", manifest.Version)
	}
	if !semverPattern.MatchString(manifest.MinCoreVersion) {
		return fmt.Errorf("invalid minCoreVersion %q", manifest.MinCoreVersion)
	}
	if compareVersion(coreVersion, manifest.MinCoreVersion) < 0 {
		return fmt.Errorf("plugin requires core %s, current core is %s", manifest.MinCoreVersion, coreVersion)
	}
	if err := validateRelativeAsset(manifest.Entry, "entry"); err != nil {
		return err
	}
	if pluginFS == nil {
		return errors.New("plugin filesystem is nil")
	}
	if info, err := fs.Stat(pluginFS, manifest.Entry); err != nil || info.IsDir() {
		return fmt.Errorf("entry file %q does not exist or is a directory", manifest.Entry)
	}
	if manifest.Icon != "" {
		if err := validateRelativeAsset(manifest.Icon, "icon"); err != nil {
			return err
		}
		if info, err := fs.Stat(pluginFS, manifest.Icon); err != nil || info.IsDir() {
			return fmt.Errorf("icon file %q does not exist or is a directory", manifest.Icon)
		}
	}
	for _, permission := range manifest.Permissions {
		if _, ok := AllowedPermissions[permission]; !ok {
			return fmt.Errorf("unknown permission %q", permission)
		}
	}
	for _, capability := range manifest.AI.Capabilities {
		if _, ok := AllowedAICapabilities[capability]; !ok {
			return fmt.Errorf("unknown AI capability %q", capability)
		}
	}
	if len(manifest.AI.Capabilities) > 0 && !manifest.AI.Enabled {
		return errors.New("AI capabilities require ai.enabled=true")
	}
	return nil
}

func validateRelativeAsset(value, field string) error {
	if value == "" || path.IsAbs(value) || strings.ContainsRune(value, '\\') {
		return fmt.Errorf("%s must be a relative POSIX path", field)
	}
	clean := path.Clean(value)
	if clean == "." || clean != value || clean == ".." || strings.HasPrefix(clean, "../") {
		return fmt.Errorf("%s contains an unsafe path", field)
	}
	return nil
}

func uniqueSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func compareVersion(left, right string) int {
	parse := func(value string) [3]int {
		base := strings.SplitN(value, "+", 2)[0]
		base = strings.SplitN(base, "-", 2)[0]
		parts := strings.Split(base, ".")
		var result [3]int
		for i := range result {
			result[i], _ = strconv.Atoi(parts[i])
		}
		return result
	}
	a, b := parse(left), parse(right)
	for i := range a {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
