package gopherpdf

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// Mirrors PyMuPDF tests/test_toc.py simple TOC expectations.
func TestPyMuPDFParity_SimpleToC(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	toc, err := doc.ToC()
	if err != nil {
		t.Fatalf("toc: %v", err)
	}

	simpleTOCRaw := string(mustReadResource(t, "simple_toc.txt"))
	wantCount := strings.Count(simpleTOCRaw, "[")
	if len(toc) != wantCount {
		t.Fatalf("toc length mismatch: got %d want %d", len(toc), wantCount)
	}
	if toc[0].Title != "HAUPTÜBERSICHT" {
		t.Fatalf("unexpected first toc title: %q", toc[0].Title)
	}
	if toc[len(toc)-1].Title != "Vorschau" {
		t.Fatalf("unexpected last toc title: %q", toc[len(toc)-1].Title)
	}
}

// Mirrors get_toc(simple=True) behavior from tests/test_toc.py.
func TestPyMuPDFParity_ToCSimpleAndGetToC(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	simple, err := doc.ToCSimple()
	if err != nil {
		t.Fatalf("toc simple: %v", err)
	}
	simpleTOCRaw := string(mustReadResource(t, "simple_toc.txt"))
	wantCount := strings.Count(simpleTOCRaw, "[")
	if len(simple) != wantCount {
		t.Fatalf("simple toc length mismatch: got %d want %d", len(simple), wantCount)
	}
	if simple[0].Page != -1 {
		t.Fatalf("expected first simple toc page to be -1 for external link, got %d", simple[0].Page)
	}

	getSimpleAny, err := doc.GetToC(true)
	if err != nil {
		t.Fatalf("get toc simple: %v", err)
	}
	getSimple, ok := getSimpleAny.([]TOCEntry)
	if !ok {
		t.Fatalf("expected GetToC(true) to return []TOCEntry")
	}
	if len(getSimple) != len(simple) {
		t.Fatalf("get toc simple length mismatch: got %d want %d", len(getSimple), len(simple))
	}

	getFullAny, err := doc.GetToC(false)
	if err != nil {
		t.Fatalf("get toc full: %v", err)
	}
	getFull, ok := getFullAny.([]Outline)
	if !ok {
		t.Fatalf("expected GetToC(false) to return []Outline")
	}
	if len(getFull) != len(simple) {
		t.Fatalf("get toc full length mismatch: got %d want %d", len(getFull), len(simple))
	}

	doc2788, err := Open(resourcePath(t, "test_2788.pdf"))
	if err != nil {
		t.Fatalf("open path test_2788: %v", err)
	}
	defer doc2788.Close()
	simple2788, err := doc2788.ToCSimple()
	if err != nil {
		t.Fatalf("toc simple test_2788: %v", err)
	}
	if len(simple2788) != 1 {
		t.Fatalf("expected 1 simple toc entry for test_2788, got %d", len(simple2788))
	}
	if simple2788[0].Page != 2 {
		t.Fatalf("expected simple toc page to be 2, got %d", simple2788[0].Page)
	}
}

// Mirrors tests/test_toc.py::test_full_toc title/level ordering.
func TestPyMuPDFParity_FullToC(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	toc, err := doc.ToC()
	if err != nil {
		t.Fatalf("toc: %v", err)
	}

	raw := string(mustReadResource(t, "full_toc.txt"))
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(toc) != len(lines) {
		t.Fatalf("full toc length mismatch: got %d want %d", len(toc), len(lines))
	}

	lineRe := regexp.MustCompile(`^\[(\d+), '(.+)', -?\d+,`)
	for i, line := range lines {
		m := lineRe.FindStringSubmatch(line)
		if len(m) != 3 {
			t.Fatalf("unexpected full_toc line format at %d: %q", i, line)
		}
		wantLevel := m[1]
		wantTitle := m[2]
		if fmt.Sprintf("%d", toc[i].Level) != wantLevel {
			t.Fatalf("level mismatch at %d: got %d want %s", i, toc[i].Level, wantLevel)
		}
		if toc[i].Title != wantTitle {
			t.Fatalf("title mismatch at %d: got %q want %q", i, toc[i].Title, wantTitle)
		}
	}
}

// Mirrors tests/test_toc.py::test_circular robustness (must not hang/crash).
func TestPyMuPDFParity_CircularToC(t *testing.T) {
	doc, err := Open(resourcePath(t, "circular-toc.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	if _, err := doc.ToC(); err != nil {
		t.Fatalf("toc circular should not error: %v", err)
	}
}

// Mirrors tests/test_toc.py::test_2788 named destination handling.
func TestPyMuPDFParity_ToCNamedDest2788(t *testing.T) {
	doc, err := Open(resourcePath(t, "test_2788.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	toc, err := doc.ToC()
	if err != nil {
		t.Fatalf("toc: %v", err)
	}
	if len(toc) != 1 {
		t.Fatalf("expected 1 toc entry, got %d", len(toc))
	}
	entry := toc[0]
	if entry.Level != 1 || entry.Title != "page2" {
		t.Fatalf("unexpected toc entry: %+v", entry)
	}
	if entry.Page != 1 {
		t.Fatalf("unexpected toc page value: %d", entry.Page)
	}
	if !strings.Contains(entry.URI, "#nameddest=") {
		t.Fatalf("expected named destination URI, got %q", entry.URI)
	}

	rl, ok, err := doc.ResolveLink(entry.URI)
	if err != nil {
		t.Fatalf("resolve toc uri: %v", err)
	}
	if !ok {
		t.Fatalf("expected toc URI to resolve")
	}
	if rl.Location != (Location{Chapter: 0, Page: entry.Page}) {
		t.Fatalf("resolved location mismatch: got %+v want page %d", rl.Location, entry.Page)
	}

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	for i := 0; i < pages; i++ {
		if _, err := doc.Links(i); err != nil {
			t.Fatalf("links page %d: %v", i, err)
		}
	}
}

// Mirrors tests/test_toc.py::test_3820 consistency check.
func TestPyMuPDFParity_ToCPageConsistency3820(t *testing.T) {
	doc, err := Open(resourcePath(t, "test-3820.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	toc, err := doc.ToC()
	if err != nil {
		t.Fatalf("toc: %v", err)
	}
	if len(toc) == 0 {
		t.Fatalf("expected toc entries")
	}
	for i, item := range toc {
		if item.Page < 0 {
			continue
		}
		rl, ok, err := doc.ResolveLink(item.URI)
		if err != nil {
			t.Fatalf("resolve toc[%d] uri: %v", i, err)
		}
		if !ok {
			t.Fatalf("expected toc[%d] uri to resolve: %q", i, item.URI)
		}
		if rl.Location.Page != item.Page {
			t.Fatalf("toc[%d] page mismatch: toc page=%d resolved=%d uri=%q", i, item.Page, rl.Location.Page, item.URI)
		}
	}
}

func TestToC(t *testing.T) {
	pdfBytes := mustReadResource(t, "test_toc_count.pdf")
	doc, err := OpenBytes(pdfBytes)
	if err != nil {
		t.Fatalf("open bytes: %v", err)
	}
	defer doc.Close()

	toc, err := doc.ToC()
	if err != nil {
		t.Fatalf("toc: %v", err)
	}
	if len(toc) < 1 {
		t.Fatalf("expected at least one toc item")
	}
}
