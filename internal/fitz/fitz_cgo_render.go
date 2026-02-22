//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

import (
	"image"
	"unsafe"
)

// Image returns image for given page number.
func (f *Document) Image(pageNumber int) (*image.RGBA, error) {
	return f.ImageDPI(pageNumber, 300.0)
}

// ImageDPI returns image for given page number and DPI.
func (f *Document) ImageDPI(pageNumber int, dpi float64) (*image.RGBA, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber >= f.NumPage() {
		return nil, ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return nil, ErrLoadPage
	}

	defer C.fz_drop_page(f.ctx, page)

	var bounds C.fz_rect
	bounds = C.fz_bound_page(f.ctx, page)

	var ctm C.fz_matrix
	ctm = C.fz_scale(C.float(dpi/72), C.float(dpi/72))

	var bbox C.fz_irect
	bounds = C.fz_transform_rect(bounds, ctm)
	bbox = C.fz_round_rect(bounds)

	pixmap := C.fz_new_pixmap_with_bbox(f.ctx, C.fz_device_rgb(f.ctx), bbox, nil, 1)
	if pixmap == nil {
		return nil, ErrCreatePixmap
	}

	C.fz_clear_pixmap_with_value(f.ctx, pixmap, C.int(0xff))
	defer C.fz_drop_pixmap(f.ctx, pixmap)

	device := C.fz_new_draw_device(f.ctx, ctm, pixmap)
	C.fz_enable_device_hints(f.ctx, device, C.FZ_NO_CACHE)
	defer C.fz_drop_device(f.ctx, device)

	drawMatrix := C.fz_identity
	ret := C.run_page_contents(f.ctx, page, device, drawMatrix, nil)
	if ret == 0 {
		return nil, ErrRunPageContents
	}

	C.fz_close_device(f.ctx, device)

	pixels := C.fz_pixmap_samples(f.ctx, pixmap)
	if pixels == nil {
		return nil, ErrPixmapSamples
	}

	img := image.NewRGBA(image.Rect(int(bbox.x0), int(bbox.y0), int(bbox.x1), int(bbox.y1)))
	copy(img.Pix, C.GoBytes(unsafe.Pointer(pixels), C.int(4*bbox.x1*bbox.y1)))

	return img, nil
}

// ImagePNG returns image for given page number as PNG bytes.
func (f *Document) ImagePNG(pageNumber int, dpi float64) ([]byte, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber >= f.NumPage() {
		return nil, ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return nil, ErrLoadPage
	}

	defer C.fz_drop_page(f.ctx, page)

	var bounds C.fz_rect
	bounds = C.fz_bound_page(f.ctx, page)

	var ctm C.fz_matrix
	ctm = C.fz_scale(C.float(dpi/72), C.float(dpi/72))

	var bbox C.fz_irect
	bounds = C.fz_transform_rect(bounds, ctm)
	bbox = C.fz_round_rect(bounds)

	pixmap := C.fz_new_pixmap_with_bbox(f.ctx, C.fz_device_rgb(f.ctx), bbox, nil, 1)
	if pixmap == nil {
		return nil, ErrCreatePixmap
	}

	C.fz_clear_pixmap_with_value(f.ctx, pixmap, C.int(0xff))
	defer C.fz_drop_pixmap(f.ctx, pixmap)

	device := C.fz_new_draw_device(f.ctx, ctm, pixmap)
	C.fz_enable_device_hints(f.ctx, device, C.FZ_NO_CACHE)
	defer C.fz_drop_device(f.ctx, device)

	drawMatrix := C.fz_identity
	ret := C.run_page_contents(f.ctx, page, device, drawMatrix, nil)
	if ret == 0 {
		return nil, ErrRunPageContents
	}

	C.fz_close_device(f.ctx, device)

	buf := C.fz_new_buffer_from_pixmap_as_png(f.ctx, pixmap, C.fz_default_color_params)
	defer C.fz_drop_buffer(f.ctx, buf)

	size := C.fz_buffer_storage(f.ctx, buf, nil)
	str := C.GoStringN(C.fz_string_from_buffer(f.ctx, buf), C.int(size))

	return []byte(str), nil
}

// SVG returns svg document for given page number.
func (f *Document) SVG(pageNumber int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber >= f.NumPage() {
		return "", ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return "", ErrLoadPage
	}

	defer C.fz_drop_page(f.ctx, page)

	var bounds C.fz_rect
	bounds = C.fz_bound_page(f.ctx, page)

	var ctm C.fz_matrix
	ctm = C.fz_scale(C.float(72.0/72), C.float(72.0/72))
	bounds = C.fz_transform_rect(bounds, ctm)

	buf := C.fz_new_buffer(f.ctx, 1024)
	defer C.fz_drop_buffer(f.ctx, buf)

	out := C.fz_new_output_with_buffer(f.ctx, buf)
	defer C.fz_drop_output(f.ctx, out)

	device := C.fz_new_svg_device(f.ctx, out, bounds.x1-bounds.x0, bounds.y1-bounds.y0, C.FZ_SVG_TEXT_AS_PATH, 1)
	C.fz_enable_device_hints(f.ctx, device, C.FZ_NO_CACHE)
	defer C.fz_drop_device(f.ctx, device)

	var cookie C.fz_cookie
	ret := C.run_page_contents(f.ctx, page, device, ctm, &cookie)
	if ret == 0 {
		return "", ErrRunPageContents
	}

	C.fz_close_device(f.ctx, device)
	C.fz_close_output(f.ctx, out)

	str := C.GoString(C.fz_string_from_buffer(f.ctx, buf))

	return str, nil
}

// Bound gives the Bounds of a given Page in the document.
func (f *Document) Bound(pageNumber int) (image.Rectangle, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber < 0 || pageNumber >= f.NumPage() {
		return image.Rectangle{}, ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return image.Rectangle{}, ErrLoadPage
	}

	defer C.fz_drop_page(f.ctx, page)

	var bounds C.fz_rect
	bounds = C.fz_bound_page(f.ctx, page)

	return image.Rect(int(bounds.x0), int(bounds.y0), int(bounds.x1), int(bounds.y1)), nil
}

// BoundBox gives the bounds of a selected page box.
func (f *Document) BoundBox(pageNumber int, box int) (image.Rectangle, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber < 0 || pageNumber >= f.NumPage() {
		return image.Rectangle{}, ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return image.Rectangle{}, ErrLoadPage
	}
	defer C.fz_drop_page(f.ctx, page)

	bounds := C.fz_bound_page_box(f.ctx, page, C.fz_box_type(box))
	return image.Rect(int(bounds.x0), int(bounds.y0), int(bounds.x1), int(bounds.y1)), nil
}

// PageHasAnnots reports whether non-widget annotations produce drawable output.
func (f *Document) PageHasAnnots(pageNumber int) (bool, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if pageNumber < 0 || pageNumber >= f.NumPage() {
		return false, ErrPageMissing
	}

	page := C.load_page(f.ctx, f.doc, C.int(pageNumber))
	if page == nil {
		return false, ErrLoadPage
	}
	defer C.fz_drop_page(f.ctx, page)

	rect := C.fz_empty_rect
	device := C.fz_new_bbox_device(f.ctx, &rect)
	if device == nil {
		return false, ErrRunPageAnnots
	}
	defer C.fz_drop_device(f.ctx, device)

	ret := C.run_page_annots(f.ctx, page, device, C.fz_identity, nil)
	C.fz_close_device(f.ctx, device)
	if ret == 0 {
		return false, ErrRunPageAnnots
	}

	return int(C.fz_is_empty_rect(rect)) == 0, nil
}

