package phase0

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"testing"

	"github.com/zalando/go-keyring"
)

// This opt-in probe touches only its own random, disposable keychain item.
func TestSystemKeyringRoundTrip(t *testing.T) {
	if os.Getenv("OIDESK_KEYRING_PROBE") != "1" {
		t.Skip("set OIDESK_KEYRING_PROBE=1 on an interactive desktop session")
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		t.Fatal(err)
	}
	account := "phase0-" + hex.EncodeToString(random[:])
	const service = "org.oidesk.phase0.probe"
	const secret = "disposable-probe-value"
	if err := keyring.Set(service, account, secret); err != nil {
		t.Fatalf("keychain unavailable (no plaintext fallback): %v", err)
	}
	t.Cleanup(func() {
		if err := keyring.Delete(service, account); err != nil && !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("probe cleanup: %v", err)
		}
	})
	got, err := keyring.Get(service, account)
	if err != nil {
		t.Fatal(err)
	}
	if got != secret {
		t.Fatal("keychain round trip mismatch")
	}
	if err := keyring.Delete(service, account); err != nil {
		t.Fatal(err)
	}
	if _, err := keyring.Get(service, account); !errors.Is(err, keyring.ErrNotFound) {
		t.Fatalf("deleted credential still accessible: %v", err)
	}
}
