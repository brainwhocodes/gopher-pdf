package gopherpdf

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func mustReadResource(t *testing.T, name string) []byte {
	t.Helper()
	path := resourcePath(t, name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read resource %s: %v", path, err)
	}
	return b
}

func resourcePath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("could not determine caller path")
	}
	moduleRoot := filepath.Dir(file)
	return filepath.Join(moduleRoot, "resources", name)
}
