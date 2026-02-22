package gopherpdf

import (
	"bytes"
	"context"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestMuPDFRenderer_CountPages_FromFixture(t *testing.T) {
	renderer := NewMuPDFRenderer()
	pdf := mustReadFixture(t)

	count, err := renderer.CountPages(context.Background(), pdf, "")
	if err != nil {
		t.Fatalf("count pages: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected 5 pages, got %d", count)
	}
}

func TestMuPDFRenderer_Render_FromFixture(t *testing.T) {
	renderer := NewMuPDFRenderer()
	pdf := mustReadFixture(t)

	pages, err := renderer.Render(context.Background(), pdf, RenderOptions{DPI: 120, Timeout: 60 * time.Second})
	if err != nil {
		t.Fatalf("render pages: %v", err)
	}
	if len(pages) != 5 {
		t.Fatalf("expected 5 rendered pages, got %d", len(pages))
	}
	first := pages[0]
	if first.PageNumber != 1 {
		t.Fatalf("expected first page number 1, got %d", first.PageNumber)
	}
	if first.Width <= 0 || first.Height <= 0 {
		t.Fatalf("invalid first page dimensions %dx%d", first.Width, first.Height)
	}
	if len(first.PNG) == 0 {
		t.Fatalf("first page PNG is empty")
	}
	if _, err := png.Decode(bytes.NewReader(first.PNG)); err != nil {
		t.Fatalf("decode first page PNG: %v", err)
	}
}

func TestMuPDFRenderer_CountPages_AcceptsPasswordOnUnprotectedPDF(t *testing.T) {
	renderer := NewMuPDFRenderer()
	pdf := mustReadFixture(t)

	count, err := renderer.CountPages(context.Background(), pdf, "secret")
	if err != nil {
		t.Fatalf("count pages with password: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected 5 pages, got %d", count)
	}
}

func TestMuPDFRenderer_Render_AcceptsPasswordOnUnprotectedPDF(t *testing.T) {
	renderer := NewMuPDFRenderer()
	pdf := mustReadFixture(t)

	pages, err := renderer.Render(context.Background(), pdf, RenderOptions{
		DPI:      110,
		Password: "secret",
		Timeout:  60 * time.Second,
	})
	if err != nil {
		t.Fatalf("render with password: %v", err)
	}
	if len(pages) != 5 {
		t.Fatalf("expected 5 rendered pages, got %d", len(pages))
	}
}

func mustReadFixture(t *testing.T) []byte {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("could not determine test file location")
	}
	dir := filepath.Dir(file)
	path := filepath.Join(dir, "resources", "multimodal_bank_statement_scan.pdf")
	b, err := os.ReadFile(path)
	if err == nil {
		return b
	}
	t.Fatalf("could not read multimodal fixture at %s: %v", path, err)
	return nil
}
