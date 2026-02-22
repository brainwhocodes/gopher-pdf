//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

func (f *Document) pdfDocLocked() *C.pdf_document {
	return C.pdf_document_from_fz_document(f.ctx, f.doc)
}

// IsDirty reports whether a PDF has unsaved changes.
func (f *Document) IsDirty() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return false
	}
	return int(C.pdf_has_unsaved_changes(f.ctx, pdf)) != 0
}

// CanSaveIncrementally reports whether incremental save is possible for a PDF.
func (f *Document) CanSaveIncrementally() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return false
	}
	return int(C.pdf_can_be_saved_incrementally(f.ctx, pdf)) != 0
}

// IsFastWebAccess reports whether a PDF is linearized (fast web view).
func (f *Document) IsFastWebAccess() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return false
	}
	return int(C.pdf_doc_was_linearized(f.ctx, pdf)) != 0
}

// IsRepaired reports whether a PDF required repair on open.
func (f *Document) IsRepaired() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return false
	}
	return int(C.pdf_was_repaired(f.ctx, pdf)) != 0
}

// PageLabel returns the display label for a page.
func (f *Document) PageLabel(pageNumber int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber < 0 || pageNumber >= f.NumPage() {
		return "", ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return "", ErrLoadPage
	}
	defer C.fz_drop_page(f.ctx, page)

	buf := make([]byte, 256)
	cbuf := (*C.char)(unsafe.Pointer(&buf[0]))
	out := C.fz_page_label(f.ctx, page, cbuf, C.int(len(buf)))
	if out == nil {
		return "", nil
	}

	return C.GoString(out), nil
}

func parsePageLabelRulesDump(dump string) []PageLabelRule {
	out := make([]PageLabelRule, 0, 8)
	for _, line := range strings.Split(dump, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) != 4 {
			continue
		}
		startPage, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		firstNum, err := strconv.Atoi(parts[3])
		if err != nil {
			continue
		}
		if firstNum <= 0 {
			firstNum = 1
		}
		out = append(out, PageLabelRule{
			StartPage:    startPage,
			Style:        parts[1],
			Prefix:       parts[2],
			FirstPageNum: firstNum,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartPage < out[j].StartPage
	})
	return out
}

func pageLabelStyleCode(style string) (C.int, error) {
	style = strings.TrimSpace(style)
	style = strings.TrimPrefix(style, "/")
	switch style {
	case "":
		return C.int(0), nil
	case "D":
		return C.int('D'), nil
	case "R":
		return C.int('R'), nil
	case "r":
		return C.int('r'), nil
	case "A":
		return C.int('A'), nil
	case "a":
		return C.int('a'), nil
	default:
		return 0, fmt.Errorf("invalid page label style %q", style)
	}
}

// PageLabelRules returns the explicit PDF page-label definition rules.
func (f *Document) PageLabelRules() ([]PageLabelRule, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if f.pdfDocLocked() == nil {
		return nil, ErrOpenDocument
	}

	cDump := C.page_label_rules_dump(f.ctx, f.doc)
	if cDump == nil {
		return []PageLabelRule{}, nil
	}
	defer C.free(unsafe.Pointer(cDump))

	return parsePageLabelRulesDump(C.GoString(cDump)), nil
}

// SetPageLabels replaces PDF page-label definition rules.
func (f *Document) SetPageLabels(rules []PageLabelRule) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	pdf := f.pdfDocLocked()
	if pdf == nil {
		return ErrOpenDocument
	}

	// Remove all existing rules first to match PyMuPDF "replace" semantics.
	cDump := C.page_label_rules_dump(f.ctx, f.doc)
	if cDump != nil {
		existing := parsePageLabelRulesDump(C.GoString(cDump))
		C.free(unsafe.Pointer(cDump))
		for _, rule := range existing {
			if C.safe_pdf_delete_page_labels(f.ctx, pdf, C.int(rule.StartPage)) == 0 {
				return ErrOpenDocument
			}
		}
	}

	if len(rules) == 0 {
		return nil
	}

	sorted := append([]PageLabelRule(nil), rules...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].StartPage < sorted[j].StartPage
	})

	pageCount := f.NumPage()
	for _, rule := range sorted {
		if rule.StartPage < 0 || rule.StartPage >= pageCount {
			return fmt.Errorf("start page out of range: %d", rule.StartPage)
		}

		style, err := pageLabelStyleCode(rule.Style)
		if err != nil {
			return err
		}

		firstPageNum := rule.FirstPageNum
		if firstPageNum <= 1 {
			firstPageNum = 1
		}

		ok := 0
		withCString(rule.Prefix, func(cPrefix *C.char) {
			ok = int(C.safe_pdf_set_page_labels(
				f.ctx,
				pdf,
				C.int(rule.StartPage),
				style,
				cPrefix,
				C.int(firstPageNum),
			))
		})
		if ok == 0 {
			return ErrOpenDocument
		}
	}

	return nil
}

// ResolveNames returns named destination map for PDF documents.
func (f *Document) ResolveNames() map[string]NamedDest {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	cDump := C.resolve_names_dump(f.ctx, f.doc)
	if cDump == nil {
		return map[string]NamedDest{}
	}
	defer C.free(unsafe.Pointer(cDump))

	dump := C.GoString(cDump)
	out := make(map[string]NamedDest)
	for _, line := range strings.Split(dump, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 5 {
			continue
		}
		page, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		x, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			continue
		}
		zoom, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			continue
		}
		out[parts[0]] = NamedDest{
			Page: page,
			X:    x,
			Y:    y,
			Zoom: zoom,
		}
	}

	return out
}

// CatalogXref returns the catalog object xref number for PDFs.
func (f *Document) CatalogXref() int {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	pdf := f.pdfDocLocked()
	if pdf == nil {
		return 0
	}
	trailer := C.pdf_trailer(f.ctx, pdf)
	if trailer == nil {
		return 0
	}
	var root *C.pdf_obj
	withCString("Root", func(cRoot *C.char) {
		root = C.pdf_dict_gets(f.ctx, trailer, cRoot)
	})
	if root == nil {
		return 0
	}
	if C.pdf_is_indirect(f.ctx, root) == 0 {
		return 0
	}
	return int(C.pdf_to_num(f.ctx, root))
}

// PageXref returns page object xref number by 0-based page index.
func (f *Document) PageXref(pageNumber int) int {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	pdf := f.pdfDocLocked()
	if pdf == nil {
		return 0
	}
	if pageNumber < 0 {
		return 0
	}
	pageObj := C.pdf_lookup_page_obj(f.ctx, pdf, C.int(pageNumber))
	if pageObj == nil {
		return 0
	}
	if C.pdf_is_indirect(f.ctx, pageObj) == 0 {
		return 0
	}
	return int(C.pdf_to_num(f.ctx, pageObj))
}

// FormFieldCount returns AcroForm field count; returns -1 when unavailable/non-PDF.
func (f *Document) FormFieldCount() int {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	pdf := f.pdfDocLocked()
	if pdf == nil {
		return -1
	}
	return int(C.safe_pdf_form_field_count(f.ctx, pdf))
}
