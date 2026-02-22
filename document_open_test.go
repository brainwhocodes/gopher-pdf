package gopherpdf

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestOpenBytesAndCoreAPIs(t *testing.T) {
	pdfBytes := mustReadResource(t, "test_2548.pdf")
	doc, err := OpenBytes(pdfBytes)
	if err != nil {
		t.Fatalf("open bytes: %v", err)
	}
	defer doc.Close()

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if pages < 1 {
		t.Fatalf("expected pages > 0, got %d", pages)
	}

	meta, err := doc.Metadata()
	if err != nil {
		t.Fatalf("metadata: %v", err)
	}
	if meta == nil {
		t.Fatalf("expected non-nil metadata")
	}

	bound, err := doc.Bound(0)
	if err != nil {
		t.Fatalf("bound page 0: %v", err)
	}
	if bound.Dx() <= 0 || bound.Dy() <= 0 {
		t.Fatalf("invalid bounds: %v", bound)
	}

	text, err := doc.Text(0)
	if err != nil {
		t.Fatalf("text page 0: %v", err)
	}
	if strings.TrimSpace(text) == "" {
		t.Fatalf("expected non-empty text on page 0")
	}

	html, err := doc.HTML(0, true)
	if err != nil {
		t.Fatalf("html page 0: %v", err)
	}
	if !strings.Contains(strings.ToLower(html), "<html") {
		t.Fatalf("expected html output, got: %.120s", html)
	}

	svg, err := doc.SVG(0)
	if err != nil {
		t.Fatalf("svg page 0: %v", err)
	}
	if !strings.Contains(strings.ToLower(svg), "<svg") {
		t.Fatalf("expected svg output, got: %.120s", svg)
	}

	_, err = doc.Links(0)
	if err != nil {
		t.Fatalf("links page 0: %v", err)
	}

	page, err := doc.RenderPage(0, 120)
	if err != nil {
		t.Fatalf("render page 0: %v", err)
	}
	if page.Width <= 0 || page.Height <= 0 || len(page.PNG) == 0 {
		t.Fatalf("invalid rendered page: %+v", page)
	}
	if !bytes.HasPrefix(page.PNG, []byte("\x89PNG")) {
		t.Fatalf("expected png bytes")
	}
}

func TestPyMuPDFParity_DocumentStateAndAliases(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	if doc.IsClosed() {
		t.Fatalf("expected opened doc to report IsClosed=false")
	}
	if !strings.HasSuffix(doc.Name(), "2.pdf") {
		t.Fatalf("unexpected document name: %q", doc.Name())
	}

	numPages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	pageCount, err := doc.PageCount()
	if err != nil {
		t.Fatalf("page count alias: %v", err)
	}
	if pageCount != numPages {
		t.Fatalf("page count alias mismatch: got %d want %d", pageCount, numPages)
	}

	hasLinks, err := doc.HasLinks()
	if err != nil {
		t.Fatalf("has links: %v", err)
	}
	if !hasLinks {
		t.Fatalf("expected 2.pdf to have links")
	}

	if err := doc.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if !doc.IsClosed() {
		t.Fatalf("expected closed doc to report IsClosed=true")
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()

	chapters, err := epub.NumChapters()
	if err != nil {
		t.Fatalf("num chapters: %v", err)
	}
	chapterCount, err := epub.ChapterCount()
	if err != nil {
		t.Fatalf("chapter count alias: %v", err)
	}
	if chapterCount != chapters {
		t.Fatalf("chapter count alias mismatch: got %d want %d", chapterCount, chapters)
	}
}

func TestOpenWithPassword_UnprotectedPDF(t *testing.T) {
	doc, err := OpenWithPassword(resourcePath(t, "2.pdf"), "secret")
	if err != nil {
		t.Fatalf("open with password should succeed for unprotected file: %v", err)
	}
	defer doc.Close()

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if pages != 47 {
		t.Fatalf("expected 47 pages, got %d", pages)
	}
}

// Mirrors part of PyMuPDF tests/test_general.py::test_open2.
func TestPyMuPDFParity_Open2_MultiFormatFileOpen(t *testing.T) {
	files := []string{
		"test_open2.pdf",
		"test_open2.epub",
		"test_open2.xps",
		"test_open2.docx",
		"test_open2.cbz",
		"test_open2.mobi",
		"test_open2.fb2",
		"test_open2.html",
		"test_open2.xhtml",
		"test_open2.xml",
		"test_open2.jpg",
		"test_open2.svg",
	}
	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			doc, err := Open(resourcePath(t, name))
			if err != nil {
				t.Fatalf("open file: %v", err)
			}
			defer doc.Close()

			pages, err := doc.NumPages()
			if err != nil {
				t.Fatalf("num pages: %v", err)
			}
			if pages < 1 {
				t.Fatalf("expected pages > 0")
			}
		})
	}
}

// Mirrors stream + filetype behavior from test_open2 for tricky text formats.
func TestPyMuPDFParity_Open2_StreamWithFileTypeHint(t *testing.T) {
	for _, tc := range []struct {
		name     string
		fileType string
	}{
		{name: "test_open2.html", fileType: ".html"},
		{name: "test_open2.xhtml", fileType: ".xhtml"},
		{name: "test_open2.xml", fileType: ".xml"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := mustReadResource(t, tc.name)
			if _, err := OpenBytes(b); err == nil {
				t.Fatalf("expected plain OpenBytes to fail without explicit file type")
			}
			doc, err := OpenBytesWithFileType(b, tc.fileType)
			if err != nil {
				t.Fatalf("open bytes with file type: %v", err)
			}
			defer doc.Close()
			pages, err := doc.NumPages()
			if err != nil {
				t.Fatalf("num pages: %v", err)
			}
			if pages < 1 {
				t.Fatalf("expected pages > 0")
			}
		})
	}
}

func TestPyMuPDFParity_OpenReaderWithFileTypeHint(t *testing.T) {
	b := mustReadResource(t, "test_open2.html")
	doc, err := OpenReaderWithFileType(bytes.NewReader(b), "html")
	if err != nil {
		t.Fatalf("open reader with file type: %v", err)
	}
	defer doc.Close()
	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if pages < 1 {
		t.Fatalf("expected pages > 0")
	}
}

func TestPyMuPDFParity_OpenReaderWithFileTypeAndPasswordHint(t *testing.T) {
	b := mustReadResource(t, "test_open2.xhtml")
	doc, err := OpenReaderWithFileTypeAndPassword(bytes.NewReader(b), ".xhtml", "secret")
	if err != nil {
		t.Fatalf("open reader with file type and password: %v", err)
	}
	defer doc.Close()
	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if pages < 1 {
		t.Fatalf("expected pages > 0")
	}
}

// Mirrors core open exception behavior from PyMuPDF tests/test_general.py::test_open.
func TestPyMuPDFParity_OpenExceptions(t *testing.T) {
	if _, err := Open(resourcePath(t, "this-file-does-not-exist.pdf")); err == nil {
		t.Fatalf("expected open missing file to fail")
	}
	if _, err := OpenBytes(nil); err == nil {
		t.Fatalf("expected open empty bytes to fail")
	}
	if _, err := OpenReader(bytes.NewReader(nil)); err == nil {
		t.Fatalf("expected open empty reader to fail")
	}
	if _, err := OpenBytesWithFileType(nil, "pdf"); err == nil {
		t.Fatalf("expected open empty bytes with filetype to fail")
	}
}

func TestInvalidPageReturnsError(t *testing.T) {
	pdfBytes := mustReadResource(t, "test_2548.pdf")
	doc, err := OpenBytes(pdfBytes)
	if err != nil {
		t.Fatalf("open bytes: %v", err)
	}
	defer doc.Close()

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if _, err := doc.Text(pages); err == nil {
		t.Fatalf("expected out-of-range page to error")
	} else if !errors.Is(err, ErrPageOutOfRange) {
		t.Fatalf("expected ErrPageOutOfRange, got %v", err)
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	pdfBytes := mustReadResource(t, "test_2548.pdf")
	doc, err := OpenBytes(pdfBytes)
	if err != nil {
		t.Fatalf("open bytes: %v", err)
	}
	if err := doc.Close(); err != nil {
		t.Fatalf("first close: %v", err)
	}
	if err := doc.Close(); err != nil {
		t.Fatalf("second close should be nil: %v", err)
	}
	if _, err := doc.NumPages(); err == nil {
		t.Fatalf("expected operations on closed doc to error")
	} else if !errors.Is(err, ErrDocumentClosed) {
		t.Fatalf("expected ErrDocumentClosed, got %v", err)
	}
}
