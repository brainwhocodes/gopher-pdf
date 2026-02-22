package gopherpdf

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func moduleRootDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller failed")
	}
	dir := filepath.Dir(file)
	for i := 0; i < 20; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		next := filepath.Dir(dir)
		if next == dir {
			break
		}
		dir = next
	}
	t.Fatalf("could not locate module root (go.mod)")
	return ""
}

func writeTestResultFile(t *testing.T, name string, data []byte) {
	t.Helper()

	dir := os.Getenv("GOPHERPDF_TEST_RESULTS_DIR")
	if dir == "" {
		return
	}

	// `go test` runs with a package working directory; treat relative paths as relative
	// to the module root (backend/).
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(moduleRootDir(t), dir)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir test results dir: %v", err)
	}

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write test result %s: %v", path, err)
	}
}
