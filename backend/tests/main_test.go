// Package tests contains a meta-test that runs every backend test suite
// (auth, users, followers, posts, groups, chat, notifications, integration)
// in one single `go test` invocation.
//
// Usage (from the backend/ directory):
//
//	RUN_ALL_SUITES=1 go test ./tests/ -run TestRunAllSuites -count=1 -timeout 1800s -v
//
// The `RUN_ALL_SUITES=1` environment variable guards the test: without it the
// test is skipped, so a plain `go test ./tests/...` (which auto-discovers this
// test) does NOT spawn every suite a second time.
//
// Each suite is spawned as its own `go test` subprocess so every suite keeps
// its isolation (fresh temp DB per test) and a failure in one suite does not
// stop the others from running.
package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// suites lists every test suite (directory) under tests/.
var suites = []string{
	"auth",
	"users",
	"followers",
	"posts",
	"groups",
	"chat",
	"notifications",
	"integration",
}

// backendDir locates the backend/ directory (the module root containing go.mod).
// It tries a few common relative locations depending on where the test is run from.
func backendDir() string {
	candidates := []string{
		".",        // running from backend/
		"..",       // running from backend/tests/
		"backend",  // running from the repository root
		"../..",    // running from backend/tests/setup/
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if info, err := os.Stat(filepath.Join(abs, "go.mod")); err == nil && !info.IsDir() {
			return abs
		}
	}
	return ""
}

// TestRunAllSuites runs every test suite under tests/ sequentially and
// fails if any of them fail. The output of each suite is printed as it
// completes, followed by a final summary.
func TestRunAllSuites(t *testing.T) {
	if os.Getenv("RUN_ALL_SUITES") == "" {
		t.Skip("skipped: set RUN_ALL_SUITES=1 to run all test suites at once")
	}

	dir := backendDir()
	if dir == "" {
		t.Fatal("could not locate backend/ directory (no go.mod found)")
	}

	failed := 0
	for _, suite := range suites {
		t.Run(suite, func(t *testing.T) {
			cmd := exec.Command("go", "test", "./tests/"+suite+"/...", "-count=1", "-timeout", "180s")
			cmd.Dir = dir

			start := time.Now()
			out, err := cmd.CombinedOutput()
			elapsed := time.Since(start).Round(time.Millisecond)

			fmt.Printf("\n=== %s suite (%s) ===\n%s\n", suite, elapsed, out)

			if err != nil {
				failed++
				t.Errorf("suite %q failed: %v", suite, err)
			}
		})
	}

	if failed > 0 {
		t.Fatalf("%d of %d suites failed", failed, len(suites))
	}
	fmt.Printf("\nAll %d test suites passed.\n", len(suites))
}
