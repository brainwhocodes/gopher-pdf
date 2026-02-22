package gopherpdf

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/png"
	"reflect"
	"strings"
	"testing"
)

func TestPyMuPDFParity_ReflowAndChapterHelpers(t *testing.T) {
	pdf, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open pdf: %v", err)
	}
	defer pdf.Close()

	isReflowable, err := pdf.IsReflowable()
	if err != nil {
		t.Fatalf("pdf is reflowable: %v", err)
	}
	if isReflowable {
		t.Fatalf("expected 2.pdf to be non-reflowable")
	}

	pdfPages, err := pdf.NumPages()
	if err != nil {
		t.Fatalf("pdf num pages: %v", err)
	}
	ch0Pages, err := pdf.ChapterPageCount(0)
	if err != nil {
		t.Fatalf("pdf chapter 0 pages: %v", err)
	}
	if ch0Pages != pdfPages {
		t.Fatalf("expected chapter 0 page count %d to equal num pages %d", ch0Pages, pdfPages)
	}
	if _, err := pdf.ChapterPageCount(1); err == nil {
		t.Fatalf("expected out-of-range chapter on pdf to fail")
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()

	isReflowable, err = epub.IsReflowable()
	if err != nil {
		t.Fatalf("epub is reflowable: %v", err)
	}
	if !isReflowable {
		t.Fatalf("expected Bezier.epub to be reflowable")
	}

	chapters, err := epub.NumChapters()
	if err != nil {
		t.Fatalf("epub num chapters: %v", err)
	}
	if chapters <= 0 {
		t.Fatalf("expected epub to have chapters")
	}

	sum := 0
	for i := 0; i < chapters; i++ {
		n, err := epub.ChapterPageCount(i)
		if err != nil {
			t.Fatalf("chapter %d pages: %v", i, err)
		}
		if n <= 0 {
			t.Fatalf("expected chapter %d to have pages, got %d", i, n)
		}
		sum += n
	}

	totalPages, err := epub.NumPages()
	if err != nil {
		t.Fatalf("epub num pages: %v", err)
	}
	if sum != totalPages {
		t.Fatalf("sum of chapter pages %d does not match total pages %d", sum, totalPages)
	}
	if _, err := epub.ChapterPageCount(chapters); err == nil {
		t.Fatalf("expected out-of-range chapter on epub to fail")
	}

	clampedMax, err := epub.ClampLocation(Location{Chapter: 999, Page: 999})
	if err != nil {
		t.Fatalf("clamp max location: %v", err)
	}
	last, err := epub.LastLocation()
	if err != nil {
		t.Fatalf("last location: %v", err)
	}
	if clampedMax != last {
		t.Fatalf("expected clamped max %+v to equal last location %+v", clampedMax, last)
	}

	clampedMin, err := epub.ClampLocation(Location{Chapter: -1, Page: -1})
	if err != nil {
		t.Fatalf("clamp min location: %v", err)
	}
	if clampedMin != (Location{Chapter: 0, Page: 0}) {
		t.Fatalf("expected clamped min location {0 0}, got %+v", clampedMin)
	}
}

func TestPyMuPDFParity_PageLabelCollectionHelpers(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open pdf: %v", err)
	}
	defer doc.Close()

	labels, err := doc.PageLabels()
	if err != nil {
		t.Fatalf("page labels: %v", err)
	}
	total, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if len(labels) != total {
		t.Fatalf("expected %d labels, got %d", total, len(labels))
	}
	if labels[0] == "" {
		t.Fatalf("expected first page label to be non-empty")
	}

	numbers, err := doc.PageNumbersByLabel(labels[0])
	if err != nil {
		t.Fatalf("page numbers by label: %v", err)
	}
	if len(numbers) == 0 {
		t.Fatalf("expected at least one page for label %q", labels[0])
	}
	for _, pageNumber := range numbers {
		got, err := doc.PageLabel(pageNumber)
		if err != nil {
			t.Fatalf("page label for %d: %v", pageNumber, err)
		}
		if got != labels[0] {
			t.Fatalf("expected page %d label %q, got %q", pageNumber, labels[0], got)
		}
	}

	none, err := doc.PageNumbersByLabel("__no-such-label__")
	if err != nil {
		t.Fatalf("page numbers by missing label: %v", err)
	}
	if len(none) != 0 {
		t.Fatalf("expected no page numbers for missing label, got %v", none)
	}

	blank, err := doc.PageNumbersByLabel("   ")
	if err != nil {
		t.Fatalf("page numbers by blank label: %v", err)
	}
	if len(blank) != 0 {
		t.Fatalf("expected no page numbers for blank label, got %v", blank)
	}

	labelsAlias, err := doc.GetPageLabels()
	if err != nil {
		t.Fatalf("get page labels alias: %v", err)
	}
	if !reflect.DeepEqual(labels, labelsAlias) {
		t.Fatalf("expected GetPageLabels alias to match PageLabels")
	}

	numbersAlias, err := doc.GetPageNumbers(labels[0])
	if err != nil {
		t.Fatalf("get page numbers alias: %v", err)
	}
	if !reflect.DeepEqual(numbers, numbersAlias) {
		t.Fatalf("expected GetPageNumbers alias to match PageNumbersByLabel")
	}
}

func TestPyMuPDFParity_TextFormatWrapperMethods(t *testing.T) {
	doc, err := Open(resourcePath(t, "symbol-list.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	html, err := doc.HTML(0, true)
	if err != nil {
		t.Fatalf("html: %v", err)
	}
	htmlWithFlags, err := doc.HTMLWithFlags(0, true, TextFlagsHTML)
	if err != nil {
		t.Fatalf("html with flags: %v", err)
	}
	if html != htmlWithFlags {
		t.Fatalf("expected HTML and HTMLWithFlags(TextFlagsHTML) to match")
	}
	if !strings.Contains(strings.ToLower(html), "<html") {
		t.Fatalf("expected html output")
	}

	xhtml, err := doc.XHTML(0, true)
	if err != nil {
		t.Fatalf("xhtml: %v", err)
	}
	xhtmlWithFlags, err := doc.XHTMLWithFlags(0, true, TextFlagsXHTML)
	if err != nil {
		t.Fatalf("xhtml with flags: %v", err)
	}
	if xhtml != xhtmlWithFlags {
		t.Fatalf("expected XHTML and XHTMLWithFlags(TextFlagsXHTML) to match")
	}

	xmlText, err := doc.XML(0)
	if err != nil {
		t.Fatalf("xml: %v", err)
	}
	xmlWithFlags, err := doc.XMLWithFlags(0, TextFlagsXML)
	if err != nil {
		t.Fatalf("xml with flags: %v", err)
	}
	if xmlText != xmlWithFlags {
		t.Fatalf("expected XML and XMLWithFlags(TextFlagsXML) to match")
	}
	if !strings.Contains(strings.ToLower(xmlText), "<page") {
		t.Fatalf("expected xml page payload")
	}

	jsonText, err := doc.JSON(0)
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	jsonWithFlags, err := doc.JSONWithFlags(0, TextFlagsDict)
	if err != nil {
		t.Fatalf("json with flags: %v", err)
	}
	if jsonText != jsonWithFlags {
		t.Fatalf("expected JSON and JSONWithFlags(TextFlagsDict) to match")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(jsonText), &parsed); err != nil {
		t.Fatalf("json parse: %v", err)
	}
	blocks, _ := parsed["blocks"].([]any)
	if len(blocks) == 0 {
		t.Fatalf("expected json output to include blocks")
	}
}

func TestPyMuPDFParity_RenderPagePNG(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	pngBytes, err := doc.RenderPagePNG(0, 72)
	if err != nil {
		t.Fatalf("render page png: %v", err)
	}
	writeTestResultFile(t, "render_2.pdf_p01_dpi072.png", pngBytes)
	if len(pngBytes) == 0 {
		t.Fatalf("expected non-empty png bytes")
	}
	if !bytes.HasPrefix(pngBytes, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}) {
		t.Fatalf("expected png signature")
	}

	cfg, err := png.DecodeConfig(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("decode png config: %v", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		t.Fatalf("expected positive rendered dimensions, got %dx%d", cfg.Width, cfg.Height)
	}

	rendered, err := doc.RenderPage(0, 72)
	if err != nil {
		t.Fatalf("render page: %v", err)
	}
	if rendered.Width != cfg.Width || rendered.Height != cfg.Height {
		t.Fatalf("expected RenderPage dimensions %dx%d to match PNG %dx%d", rendered.Width, rendered.Height, cfg.Width, cfg.Height)
	}

	defaultDPI, err := doc.RenderPage(0, 0)
	if err != nil {
		t.Fatalf("render page with default dpi: %v", err)
	}
	if defaultDPI.DPI != 200 {
		t.Fatalf("expected default dpi to be 200, got %v", defaultDPI.DPI)
	}
}

func TestPyMuPDFParity_XrefGetKeyValue(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	kv, err := doc.XrefGetKeyValue(-1, "Root")
	if err != nil {
		t.Fatalf("xref get key value: %v", err)
	}
	if kv.Type == "" || kv.Value == "" {
		t.Fatalf("expected typed xref key/value, got %+v", kv)
	}

	typ, value, err := doc.XrefGetKeyTyped(-1, "Root")
	if err != nil {
		t.Fatalf("xref get key typed: %v", err)
	}
	if kv.Type != typ || kv.Value != value {
		t.Fatalf("expected XrefGetKeyValue to mirror XrefGetKeyTyped, got %+v vs (%q,%q)", kv, typ, value)
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()
	if _, err := epub.XrefGetKeyValue(-1, "Root"); !errors.Is(err, ErrNoPDF) {
		t.Fatalf("expected ErrNoPDF on non-pdf xref get key value, got %v", err)
	}
}

func TestPyMuPDFParity_DocumentConvenienceAliases(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	aliasText, err := doc.GetPageText(0, "text", 0)
	if err != nil {
		t.Fatalf("get page text alias: %v", err)
	}
	baseText, err := doc.GetText(0, "text", 0)
	if err != nil {
		t.Fatalf("get text base: %v", err)
	}
	if aliasText != baseText {
		t.Fatalf("expected GetPageText alias to match GetText")
	}

	aliasHits, err := doc.SearchPageFor(0, "Portfolio", TextFlagsSearch, 0, nil)
	if err != nil {
		t.Fatalf("search page for alias: %v", err)
	}
	baseHits, err := doc.Search(0, "Portfolio", TextFlagsSearch, 0, nil)
	if err != nil {
		t.Fatalf("search base: %v", err)
	}
	if !reflect.DeepEqual(aliasHits, baseHits) {
		t.Fatalf("expected SearchPageFor alias to match Search")
	}

	aliasPixmap, err := doc.GetPagePixmap(0, 72)
	if err != nil {
		t.Fatalf("get page pixmap alias: %v", err)
	}
	basePixmap, err := doc.ImageDPI(0, 72)
	if err != nil {
		t.Fatalf("image dpi base: %v", err)
	}
	if aliasPixmap.Rect.Dx() != basePixmap.Rect.Dx() || aliasPixmap.Rect.Dy() != basePixmap.Rect.Dy() {
		t.Fatalf("expected GetPagePixmap alias dimensions %dx%d to match ImageDPI %dx%d", aliasPixmap.Rect.Dx(), aliasPixmap.Rect.Dy(), basePixmap.Rect.Dx(), basePixmap.Rect.Dy())
	}
}
