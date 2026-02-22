package gopherpdf

import (
	"testing"
)

// Mirrors PyMuPDF tests/test_nonpdf.py basic open/is_pdf coverage for EPUB.
func TestPyMuPDFParity_NonPDF_EpubOpens(t *testing.T) {
	doc, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer doc.Close()

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages epub: %v", err)
	}
	if pages < 1 {
		t.Fatalf("expected epub to expose pages, got %d", pages)
	}
}

// Mirrors tests/test_nonpdf.py::test_isnopdf.
func TestPyMuPDFParity_IsPDF(t *testing.T) {
	pdfDoc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open pdf: %v", err)
	}
	defer pdfDoc.Close()

	isPDF, err := pdfDoc.IsPDF()
	if err != nil {
		t.Fatalf("is pdf on pdf: %v", err)
	}
	if !isPDF {
		t.Fatalf("expected 2.pdf to be reported as PDF")
	}

	epubDoc, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epubDoc.Close()

	isPDF, err = epubDoc.IsPDF()
	if err != nil {
		t.Fatalf("is pdf on epub: %v", err)
	}
	if isPDF {
		t.Fatalf("expected Bezier.epub to be reported as non-PDF")
	}
}

// Mirrors PyMuPDF tests/test_nonpdf.py page-location and layout behavior.
func TestPyMuPDFParity_NonPDF_PageIDsAndLayout(t *testing.T) {
	doc, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer doc.Close()

	chapters, err := doc.NumChapters()
	if err != nil {
		t.Fatalf("num chapters: %v", err)
	}
	if chapters != 7 {
		t.Fatalf("expected 7 chapters, got %d", chapters)
	}

	last, err := doc.LastLocation()
	if err != nil {
		t.Fatalf("last location: %v", err)
	}
	if last != (Location{Chapter: 6, Page: 1}) {
		t.Fatalf("unexpected last location: %+v", last)
	}

	prev, err := doc.PrevLocation(Location{Chapter: 6, Page: 0})
	if err != nil {
		t.Fatalf("prev location: %v", err)
	}
	if prev != (Location{Chapter: 5, Page: 11}) {
		t.Fatalf("unexpected prev location: %+v", prev)
	}

	next, err := doc.NextLocation(Location{Chapter: 5, Page: 11})
	if err != nil {
		t.Fatalf("next location: %v", err)
	}
	if next != (Location{Chapter: 6, Page: 0}) {
		t.Fatalf("unexpected next location: %+v", next)
	}

	totalPages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	for i := 0; i < totalPages; i++ {
		loc, err := doc.LocationFromPageNumber(i)
		if err != nil {
			t.Fatalf("location from page number %d: %v", i, err)
		}
		pageNo, err := doc.PageNumberFromLocation(loc)
		if err != nil {
			t.Fatalf("page number from location %+v: %v", loc, err)
		}
		if pageNo != i {
			t.Fatalf("roundtrip mismatch: page=%d loc=%+v back=%d", i, loc, pageNo)
		}
	}

	mark, err := doc.MakeBookmark(Location{Chapter: 5, Page: 11})
	if err != nil {
		t.Fatalf("make bookmark: %v", err)
	}
	if err := doc.Layout(595, 842, 12); err != nil {
		t.Fatalf("layout: %v", err)
	}
	after, err := doc.LookupBookmark(mark)
	if err != nil {
		t.Fatalf("lookup bookmark: %v", err)
	}
	if after != (Location{Chapter: 5, Page: 6}) {
		t.Fatalf("unexpected location after layout: %+v", after)
	}
}

