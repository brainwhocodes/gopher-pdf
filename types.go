package gopherpdf

import "time"

// Text extraction flags. Values mirror MuPDF FZ_STEXT_* semantics.
const (
	TextPreserveLigatures       = 1
	TextPreserveWhitespace      = 2
	TextPreserveImages          = 4
	TextInhibitSpaces           = 8
	TextDehyphenate             = 16
	TextPreserveSpans           = 32
	TextMediaboxClip            = 64
	TextUseCIDForUnknownUnicode = 128
)

// Predefined text flag bundles mirroring PyMuPDF TEXTFLAGS_* presets.
const (
	TextFlagsWords   = TextPreserveLigatures | TextPreserveWhitespace | TextMediaboxClip | TextUseCIDForUnknownUnicode
	TextFlagsBlocks  = TextFlagsWords
	TextFlagsDict    = TextPreserveLigatures | TextPreserveWhitespace | TextMediaboxClip | TextPreserveImages | TextUseCIDForUnknownUnicode
	TextFlagsRawDict = TextFlagsDict
	TextFlagsSearch  = TextPreserveLigatures | TextPreserveWhitespace | TextMediaboxClip | TextDehyphenate | TextUseCIDForUnknownUnicode
	TextFlagsHTML    = TextFlagsDict
	TextFlagsXHTML   = TextFlagsDict
	TextFlagsXML     = TextFlagsWords
	TextFlagsText    = TextFlagsWords
)

// Permission identifies document operation permissions.
type Permission int

// Permission flags. Values mirror MuPDF fz_permission semantics.
const (
	PermissionPrint         Permission = Permission('p')
	PermissionCopy          Permission = Permission('c')
	PermissionEdit          Permission = Permission('e')
	PermissionAnnotate      Permission = Permission('n')
	PermissionForm          Permission = Permission('f')
	PermissionAccessibility Permission = Permission('y')
	PermissionAssemble      Permission = Permission('a')
	PermissionPrintHQ       Permission = Permission('h')
)

// PageBox identifies which PDF page box to query.
type PageBox int

// Page box selectors. Values mirror MuPDF fz_box_type semantics.
const (
	PageBoxMedia   PageBox = 0
	PageBoxCrop    PageBox = 1
	PageBoxBleed   PageBox = 2
	PageBoxTrim    PageBox = 3
	PageBoxArt     PageBox = 4
	PageBoxUnknown PageBox = 5
)

// Link mirrors a subset of link data exposed by MuPDF bindings.
type Link struct {
	URI string
}

// Location identifies chapter/page tuple for reflowable documents.
type Location struct {
	Chapter int
	Page    int
}

// Bookmark is an opaque bookmark token.
type Bookmark uint64

// ResolvedLink is the resolved target of an internal document link.
type ResolvedLink struct {
	Location Location
	X        float64
	Y        float64
}

// RectF is a floating-point rectangle in page coordinates.
type RectF struct {
	X0 float64
	Y0 float64
	X1 float64
	Y1 float64
}

// Contains reports whether r fully contains other.
func (r RectF) Contains(other RectF) bool {
	return r.X0 <= other.X0 && r.Y0 <= other.Y0 && r.X1 >= other.X1 && r.Y1 >= other.Y1
}

// Outline mirrors table-of-contents entries.
type Outline struct {
	Level int
	Title string
	URI   string
	Page  int
	Top   float64
}

// TOCEntry is a simplified table-of-contents entry with 1-based page numbering.
type TOCEntry struct {
	Level int
	Title string
	Page  int
}

// NamedDestination describes a resolved PDF name destination.
type NamedDestination struct {
	Page int
	To   [2]float64
	Zoom float64
}

// PageLabelRule describes one explicit PDF page-label definition rule.
type PageLabelRule struct {
	StartPage    int
	Prefix       string
	Style        string
	FirstPageNum int
}

// XrefKeyValue is a typed xref key lookup result.
type XrefKeyValue struct {
	Type  string
	Value string
}

// RenderedPage is one rasterized page from a PDF document.
type RenderedPage struct {
	PageNumber int
	Width      int
	Height     int
	DPI        float64
	RenderMS   int
	PNG        []byte
}

// RenderOptions controls page rasterization behavior.
type RenderOptions struct {
	DPI      int
	Password string
	Timeout  time.Duration
}

// Defaults applies sane render defaults for OCR ingestion.
func (o RenderOptions) Defaults() RenderOptions {
	if o.DPI <= 0 {
		o.DPI = 200
	}
	if o.Timeout <= 0 {
		o.Timeout = 90 * time.Second
	}
	return o
}
