//go:build darwin || linux

package phase0

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestCPPProcessGroupTimeout(t *testing.T) {
	if os.Getenv("OIDESK_CPP_PROBE") != "1" {
		t.Skip("set OIDESK_CPP_PROBE=1 to compile and run C++")
	}
	compiler := os.Getenv("CXX")
	if compiler == "" {
		compiler = "c++"
	}
	dir := t.TempDir()
	source := filepath.Join(dir, "probe.cpp")
	binary := filepath.Join(dir, "probe")
	marker := filepath.Join(dir, "escaped-child")
	code := `#include <unistd.h>
#include <fstream>
#include <iostream>
int main(int argc, char** argv) {
  if (argc < 2) { std::cout << "phase0-ok\n"; return 0; }
  if (fork() == 0) { usleep(700000); std::ofstream(argv[1]) << "escaped"; return 0; }
  std::cout << "started\n" << std::flush;
  sleep(30);
}`
	if err := os.WriteFile(source, []byte(code), 0600); err != nil {
		t.Fatal(err)
	}
	compileCtx, stopCompile := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopCompile()
	if out, err := exec.CommandContext(compileCtx, compiler, "-std=c++17", source, "-o", binary).CombinedOutput(); err != nil {
		t.Fatalf("compile: %v\n%s", err, out)
	}
	if out, err := exec.Command(binary).CombinedOutput(); err != nil || string(out) != "phase0-ok\n" {
		t.Fatalf("run: %v %s", err, out)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, marker)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = time.Second
	out, err := cmd.CombinedOutput()
	if err == nil || ctx.Err() == nil || !bytes.Contains(out, []byte("started")) {
		t.Fatalf("timeout: %v ctx=%v output=%s", err, ctx.Err(), out)
	}
	time.Sleep(800 * time.Millisecond)
	if _, err = os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("descendant survived process-group cancellation: %v", err)
	}
}
