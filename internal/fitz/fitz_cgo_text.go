//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

func (f *Document) stextPageLocked(pageNumber int, flags C.int) (*C.fz_stext_page, *C.fz_page, error) {
	if pageNumber < 0 || pageNumber >= f.NumPage() {
		return nil, nil, ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return nil, nil, ErrLoadPage
	}

	var bounds C.fz_rect
	bounds = C.fz_bound_page(f.ctx, page)

	var ctm C.fz_matrix
	ctm = C.fz_scale(C.float(72.0/72), C.float(72.0/72))

	text := C.fz_new_stext_page(f.ctx, bounds)
	if text == nil {
		C.fz_drop_page(f.ctx, page)
		return nil, nil, ErrRunPageContents
	}

	var opts C.fz_stext_options
	opts.flags = flags

	device := C.fz_new_stext_device(f.ctx, text, &opts)
	C.fz_enable_device_hints(f.ctx, device, C.FZ_NO_CACHE)

	var cookie C.fz_cookie
	ret := C.run_page_contents(f.ctx, page, device, ctm, &cookie)
	C.fz_close_device(f.ctx, device)
	C.fz_drop_device(f.ctx, device)
	if ret == 0 {
		C.fz_drop_stext_page(f.ctx, text)
		C.fz_drop_page(f.ctx, page)
		return nil, nil, ErrRunPageContents
	}

	return text, page, nil
}

// Text returns text for given page number.
func (f *Document) Text(pageNumber int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.textLocked(pageNumber, 0)
}

// TextWithFlags returns text for a page with MuPDF stext flags applied.
func (f *Document) TextWithFlags(pageNumber int, flags int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.textLocked(pageNumber, C.int(flags))
}

func (f *Document) textLocked(pageNumber int, flags C.int) (string, error) {
	text, page, err := f.stextPageLocked(pageNumber, flags)
	if err != nil {
		return "", err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	buf := C.fz_new_buffer_from_stext_page(f.ctx, text)
	defer C.fz_drop_buffer(f.ctx, buf)

	return C.GoString(C.fz_string_from_buffer(f.ctx, buf)), nil
}

// HTML returns html for given page number.
func (f *Document) HTML(pageNumber int, header bool) (string, error) {
	return f.HTMLWithFlags(pageNumber, header, TextFlagsHTML)
}

// HTMLWithFlags returns html for given page number with explicit stext flags.
func (f *Document) HTMLWithFlags(pageNumber int, header bool, flags int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.htmlLocked(pageNumber, header, C.int(flags))
}

// XHTML returns xhtml for given page number.
func (f *Document) XHTML(pageNumber int, header bool) (string, error) {
	return f.XHTMLWithFlags(pageNumber, header, TextFlagsXHTML)
}

// XHTMLWithFlags returns xhtml for given page number with explicit stext flags.
func (f *Document) XHTMLWithFlags(pageNumber int, header bool, flags int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.xhtmlLocked(pageNumber, header, C.int(flags))
}

// XML returns xml for given page number.
func (f *Document) XML(pageNumber int) (string, error) {
	return f.XMLWithFlags(pageNumber, TextFlagsXML)
}

// XMLWithFlags returns xml for given page number with explicit stext flags.
func (f *Document) XMLWithFlags(pageNumber int, flags int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.xmlLocked(pageNumber, C.int(flags))
}

// JSON returns structured JSON text extraction for given page number.
func (f *Document) JSON(pageNumber int, scale float64) (string, error) {
	return f.JSONWithFlags(pageNumber, TextFlagsDict, scale)
}

// JSONWithFlags returns structured JSON text extraction with explicit stext flags.
func (f *Document) JSONWithFlags(pageNumber int, flags int, scale float64) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	if scale <= 0 {
		scale = 1.0
	}
	return f.jsonLocked(pageNumber, C.int(flags), C.float(scale))
}

// SearchCount returns the number of text hits for a needle on a page.
func (f *Document) SearchCount(pageNumber int, needle string, flags int, maxHits int) (int, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if strings.TrimSpace(needle) == "" {
		return 0, nil
	}
	if maxHits <= 0 {
		maxHits = 512
	}

	text, page, err := f.stextPageLocked(pageNumber, C.int(flags))
	if err != nil {
		return 0, err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	hitMarks := make([]C.int, maxHits)
	hitBBoxes := make([]C.fz_quad, maxHits)
	var n C.int
	withCString(needle, func(cNeedle *C.char) {
		n = C.fz_search_stext_page(
			f.ctx,
			text,
			cNeedle,
			(*C.int)(unsafe.Pointer(&hitMarks[0])),
			(*C.fz_quad)(unsafe.Pointer(&hitBBoxes[0])),
			C.int(maxHits),
		)
	})

	return int(n), nil
}

// Search returns bounding rectangles of text hits for a needle on a page.
// If clip is non-nil, only hits fully inside the clip are returned.
func (f *Document) Search(pageNumber int, needle string, flags int, maxHits int, clip *Rect) ([]Rect, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if strings.TrimSpace(needle) == "" {
		return nil, nil
	}
	if maxHits <= 0 {
		maxHits = 512
	}

	text, page, err := f.stextPageLocked(pageNumber, C.int(flags))
	if err != nil {
		return nil, err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	hitMarks := make([]C.int, maxHits)
	hitBBoxes := make([]C.fz_quad, maxHits)
	n := 0
	withCString(needle, func(cNeedle *C.char) {
		n = int(C.fz_search_stext_page(
			f.ctx,
			text,
			cNeedle,
			(*C.int)(unsafe.Pointer(&hitMarks[0])),
			(*C.fz_quad)(unsafe.Pointer(&hitBBoxes[0])),
			C.int(maxHits),
		))
	})
	if n <= 0 {
		return nil, nil
	}
	if n > maxHits {
		n = maxHits
	}

	contains := func(haystack, needle Rect) bool {
		return haystack.X0 <= needle.X0 &&
			haystack.Y0 <= needle.Y0 &&
			haystack.X1 >= needle.X1 &&
			haystack.Y1 >= needle.Y1
	}

	out := make([]Rect, 0, n)
	for i := 0; i < n; i++ {
		r := C.fz_rect_from_quad(hitBBoxes[i])
		hit := Rect{
			X0: float64(r.x0),
			Y0: float64(r.y0),
			X1: float64(r.x1),
			Y1: float64(r.y1),
		}
		if clip != nil && !contains(*clip, hit) {
			continue
		}
		out = append(out, hit)
	}

	return out, nil
}

// TextBox returns text inside a rectangular clip area on a page.
func (f *Document) TextBox(pageNumber int, area Rect, flags int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	text, page, err := f.stextPageLocked(pageNumber, C.int(flags))
	if err != nil {
		return "", err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	cArea := C.fz_rect{
		x0: C.float(area.X0),
		y0: C.float(area.Y0),
		x1: C.float(area.X1),
		y1: C.float(area.Y1),
	}
	cText := C.copy_rectangle(f.ctx, text, cArea, C.int(0))
	if cText == nil {
		return "", ErrCopyRectangle
	}
	defer C.fz_free(f.ctx, unsafe.Pointer(cText))

	return C.GoString(cText), nil
}

func (f *Document) htmlLocked(pageNumber int, header bool, flags C.int) (string, error) {
	text, page, err := f.stextPageLocked(pageNumber, flags)
	if err != nil {
		return "", err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	buf := C.fz_new_buffer(f.ctx, 1024)
	defer C.fz_drop_buffer(f.ctx, buf)

	out := C.fz_new_output_with_buffer(f.ctx, buf)
	defer C.fz_drop_output(f.ctx, out)

	if header {
		C.fz_print_stext_header_as_html(f.ctx, out)
	}
	C.fz_print_stext_page_as_html(f.ctx, out, text, C.int(pageNumber))
	if header {
		C.fz_print_stext_trailer_as_html(f.ctx, out)
	}

	C.fz_close_output(f.ctx, out)

	return C.GoString(C.fz_string_from_buffer(f.ctx, buf)), nil
}

func (f *Document) xhtmlLocked(pageNumber int, header bool, flags C.int) (string, error) {
	text, page, err := f.stextPageLocked(pageNumber, flags)
	if err != nil {
		return "", err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	buf := C.fz_new_buffer(f.ctx, 1024)
	defer C.fz_drop_buffer(f.ctx, buf)

	out := C.fz_new_output_with_buffer(f.ctx, buf)
	defer C.fz_drop_output(f.ctx, out)

	if header {
		C.fz_print_stext_header_as_xhtml(f.ctx, out)
	}
	C.fz_print_stext_page_as_xhtml(f.ctx, out, text, C.int(pageNumber))
	if header {
		C.fz_print_stext_trailer_as_xhtml(f.ctx, out)
	}

	C.fz_close_output(f.ctx, out)

	return C.GoString(C.fz_string_from_buffer(f.ctx, buf)), nil
}

func (f *Document) xmlLocked(pageNumber int, flags C.int) (string, error) {
	text, page, err := f.stextPageLocked(pageNumber, flags)
	if err != nil {
		return "", err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	buf := C.fz_new_buffer(f.ctx, 1024)
	defer C.fz_drop_buffer(f.ctx, buf)

	out := C.fz_new_output_with_buffer(f.ctx, buf)
	defer C.fz_drop_output(f.ctx, out)

	C.fz_print_stext_page_as_xml(f.ctx, out, text, C.int(pageNumber))
	C.fz_close_output(f.ctx, out)

	return C.GoString(C.fz_string_from_buffer(f.ctx, buf)), nil
}

func (f *Document) jsonLocked(pageNumber int, flags C.int, scale C.float) (string, error) {
	text, page, err := f.stextPageLocked(pageNumber, flags)
	if err != nil {
		return "", err
	}
	defer C.fz_drop_stext_page(f.ctx, text)
	defer C.fz_drop_page(f.ctx, page)

	buf := C.fz_new_buffer(f.ctx, 1024)
	defer C.fz_drop_buffer(f.ctx, buf)

	out := C.fz_new_output_with_buffer(f.ctx, buf)
	defer C.fz_drop_output(f.ctx, out)

	C.fz_print_stext_page_as_json(f.ctx, out, text, scale)
	C.fz_close_output(f.ctx, out)

	return C.GoString(C.fz_string_from_buffer(f.ctx, buf)), nil
}
