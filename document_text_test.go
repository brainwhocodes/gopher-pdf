package gopherpdf

import (
	"encoding/json"
	"strings"
	"testing"
)

// Mirrors coverage intent of PyMuPDF tests/test_textextract.py::test_extract1.
func TestPyMuPDFParity_TextExtractFormats(t *testing.T) {
	doc, err := Open(resourcePath(t, "symbol-list.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	for _, mode := range []string{"text", "blocks", "words", "dict", "json", "rawdict", "rawjson", "html", "xhtml", "xml"} {
		got, err := doc.GetText(0, mode, 0)
		if err != nil {
			t.Fatalf("mode %q: %v", mode, err)
		}
		switch mode {
		case "text":
			s, _ := got.(string)
			if strings.TrimSpace(s) == "" {
				t.Fatalf("mode %q returned empty text", mode)
			}
		case "html":
			s, _ := got.(string)
			if !strings.Contains(strings.ToLower(s), "<html") {
				t.Fatalf("mode %q did not return html", mode)
			}
		case "xhtml":
			s, _ := got.(string)
			if !strings.Contains(strings.ToLower(s), "<html") {
				t.Fatalf("mode %q did not return xhtml", mode)
			}
		case "xml":
			s, _ := got.(string)
			if !strings.Contains(strings.ToLower(s), "<page") {
				t.Fatalf("mode %q did not return xml page payload", mode)
			}
		case "json", "rawjson":
			s, _ := got.(string)
			var out map[string]any
			if err := json.Unmarshal([]byte(s), &out); err != nil {
				t.Fatalf("mode %q invalid json: %v", mode, err)
			}
			blocks, _ := out["blocks"].([]any)
			if len(blocks) == 0 {
				t.Fatalf("mode %q returned no blocks", mode)
			}
		case "dict", "rawdict":
			out, _ := got.(map[string]any)
			if out == nil {
				t.Fatalf("mode %q did not return a map", mode)
			}
			blocks, _ := out["blocks"].([]any)
			if len(blocks) == 0 {
				t.Fatalf("mode %q returned no blocks", mode)
			}
		case "blocks":
			blocks, _ := got.([]any)
			if len(blocks) == 0 {
				t.Fatalf("mode %q returned no blocks", mode)
			}
		case "words":
			words, _ := got.([]string)
			if len(words) == 0 {
				t.Fatalf("mode %q returned no words", mode)
			}
		}
	}
}

// Mirrors behavior under PyMuPDF tests/test_textextract.py::test_2954.
func TestPyMuPDFParity_TextFlagsCIDForUnknownUnicode(t *testing.T) {
	doc, err := Open(resourcePath(t, "test_2954.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	flags0 := TextPreserveWhitespace | TextPreserveLigatures | TextMediaboxClip
	text0, err := doc.TextWithFlags(0, flags0)
	if err != nil {
		t.Fatalf("text flags0: %v", err)
	}
	text1, err := doc.TextWithFlags(0, flags0|TextUseCIDForUnknownUnicode)
	if err != nil {
		t.Fatalf("text flags1: %v", err)
	}

	countFFFD := func(s string) int { return strings.Count(s, string(rune(0xfffd))) }
	n0 := countFFFD(text0)
	n1 := countFFFD(text1)
	if n0 == 0 {
		t.Fatalf("expected flags0 to include replacement chars, got 0")
	}
	if n1 != 0 {
		t.Fatalf("expected flags1 to remove replacement chars, got %d", n1)
	}
}

// Mirrors PyMuPDF tests/test_textsearch.py core search count behavior.
func TestPyMuPDFParity_TextSearch(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	hits, err := doc.SearchCount(0, "mupdf", 0, 128)
	if err != nil {
		t.Fatalf("search count: %v", err)
	}
	if hits < 1 {
		t.Fatalf("expected at least one search hit for mupdf, got %d", hits)
	}

	doc2, err := Open(resourcePath(t, "text-find-ligatures.pdf"))
	if err != nil {
		t.Fatalf("open path ligatures: %v", err)
	}
	defer doc2.Close()

	allHits, err := doc2.SearchCount(0, "flag", 0, 128)
	if err != nil {
		t.Fatalf("search ligature all hits: %v", err)
	}
	preservedHits, err := doc2.SearchCount(0, "flag", TextPreserveLigatures, 128)
	if err != nil {
		t.Fatalf("search ligature preserved hits: %v", err)
	}
	if allHits != 2 {
		t.Fatalf("expected all-hits to equal 2, got %d", allHits)
	}
	if preservedHits != 1 {
		t.Fatalf("expected preserved-hits to equal 1, got %d", preservedHits)
	}
}

// Mirrors tests/test_textsearch.py::test_search1.
func TestPyMuPDFParity_TextSearchAndTextBox(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	const needle = "mupdf"
	hits, err := doc.Search(0, needle, TextFlagsSearch, 128, nil)
	if err != nil {
		t.Fatalf("search hits: %v", err)
	}
	if len(hits) == 0 {
		t.Fatalf("expected at least one search hit")
	}
	for i, hit := range hits {
		boxText, err := doc.TextBox(0, hit, TextFlagsSearch)
		if err != nil {
			t.Fatalf("textbox for hit %d: %v", i, err)
		}
		if !strings.Contains(strings.ToLower(boxText), needle) {
			t.Fatalf("textbox for hit %d did not contain needle: %q", i, boxText)
		}
	}
}

// Mirrors tests/test_textsearch.py::test_search2 clip behavior.
func TestPyMuPDFParity_TextSearchClip(t *testing.T) {
	doc, err := Open(resourcePath(t, "github_sample.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	clip := RectF{
		X0: 40.5,
		Y0: 228.31436157226562,
		X1: 346.5226135253906,
		Y1: 239.5338592529297,
	}
	hits, err := doc.Search(0, "the", TextFlagsSearch, 128, &clip)
	if err != nil {
		t.Fatalf("search with clip: %v", err)
	}
	if len(hits) != 2 {
		t.Fatalf("expected 2 clipped hits, got %d", len(hits))
	}
	for i, hit := range hits {
		if !clip.Contains(hit) {
			t.Fatalf("hit %d not inside clip: hit=%+v clip=%+v", i, hit, clip)
		}
	}
}

func TestPyMuPDFParity_TextFlagBundles(t *testing.T) {
	if TextFlagsSearch != (TextPreserveLigatures | TextPreserveWhitespace | TextMediaboxClip | TextDehyphenate | TextUseCIDForUnknownUnicode) {
		t.Fatalf("unexpected TextFlagsSearch value: %d", TextFlagsSearch)
	}
	if TextFlagsDict != (TextPreserveLigatures | TextPreserveWhitespace | TextMediaboxClip | TextPreserveImages | TextUseCIDForUnknownUnicode) {
		t.Fatalf("unexpected TextFlagsDict value: %d", TextFlagsDict)
	}
	if TextFlagsText != (TextPreserveLigatures | TextPreserveWhitespace | TextMediaboxClip | TextUseCIDForUnknownUnicode) {
		t.Fatalf("unexpected TextFlagsText value: %d", TextFlagsText)
	}
}

func TestGetTextUnsupportedModeReturnsError(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()
	if _, err := doc.GetText(0, "does-not-exist", 0); err == nil {
		t.Fatalf("expected unsupported mode to return error")
	}
}
