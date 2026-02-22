// Package fitz provides wrapper for the [MuPDF](http://mupdf.com/) fitz library
// that can extract pages from PDF, EPUB, MOBI, DOCX, XLSX and PPTX documents as IMG, TXT, HTML or SVG.
package fitz

import (
	"errors"
	"unsafe"
)

// Errors.
var (
	ErrNoSuchFile      = errors.New("fitz: no such file")
	ErrCreateContext   = errors.New("fitz: cannot create context")
	ErrOpenDocument    = errors.New("fitz: cannot open document")
	ErrEmptyBytes      = errors.New("fitz: cannot send empty bytes")
	ErrOpenMemory      = errors.New("fitz: cannot open memory")
	ErrLoadPage        = errors.New("fitz: cannot load page")
	ErrRunPageContents = errors.New("fitz: cannot run page contents")
	ErrRunPageAnnots   = errors.New("fitz: cannot run page annotations")
	ErrPageMissing     = errors.New("fitz: page missing")
	ErrCreatePixmap    = errors.New("fitz: cannot create pixmap")
	ErrPixmapSamples   = errors.New("fitz: cannot get pixmap samples")
	ErrNeedsPassword   = errors.New("fitz: document needs password")
	ErrLoadOutline     = errors.New("fitz: cannot load outline")
	ErrCopyRectangle   = errors.New("fitz: cannot copy rectangle text")
	ErrCGODisabled     = errors.New("fitz: cgo is required; build with CGO_ENABLED=1")
)

// MaxStore is maximum size in bytes of the resource store, before it will start evicting cached resources such as fonts and images.
var MaxStore = 256 << 20

// FzVersion is used for experimental purego implementation, it must be exactly the same as libmupdf shared library version.
// It is also possible to set `FZ_VERSION` environment variable.
var FzVersion = "1.24.9"

// Outline type.
type Outline struct {
	// Hierarchy level of the entry (starting from 1).
	Level int
	// Title of outline item.
	Title string
	// Destination in the document to be displayed when this outline item is activated.
	URI string
	// The page number of an internal link.
	Page int
	// Top.
	Top float64
}

// Link type.
type Link struct {
	URI string
}

// Rect is a floating-point rectangle in page coordinates.
type Rect struct {
	X0 float64
	Y0 float64
	X1 float64
	Y1 float64
}

// NamedDest describes a resolved named destination in a PDF.
type NamedDest struct {
	Page int
	X    float64
	Y    float64
	Zoom float64
}

// PageLabelRule describes one PDF page-label definition rule.
type PageLabelRule struct {
	StartPage    int
	Prefix       string
	Style        string
	FirstPageNum int
}

// Location identifies a chapter/page tuple for reflowable documents.
type Location struct {
	Chapter int
	Page    int
}

// Bookmark is an opaque bookmark token from MuPDF.
type Bookmark uintptr

// Structured text extraction flags. Values match MuPDF FZ_STEXT_*.
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

// Permission flags. Values match MuPDF fz_permission enum character codes.
const (
	PermissionPrint         = int('p')
	PermissionCopy          = int('c')
	PermissionEdit          = int('e')
	PermissionAnnotate      = int('n')
	PermissionForm          = int('f')
	PermissionAccessibility = int('y')
	PermissionAssemble      = int('a')
	PermissionPrintHQ       = int('h')
)

// Page box selectors. Values match MuPDF fz_box_type enum.
const (
	PageBoxMedia   = 0
	PageBoxCrop    = 1
	PageBoxBleed   = 2
	PageBoxTrim    = 3
	PageBoxArt     = 4
	PageBoxUnknown = 5
)

func bytePtrToString(p *byte) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*byte)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	}

	return string(unsafe.Slice(p, n))
}
