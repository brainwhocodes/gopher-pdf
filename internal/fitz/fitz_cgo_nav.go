//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

func locationToC(loc Location) C.fz_location {
	return C.fz_location{
		chapter: C.int(loc.Chapter),
		page:    C.int(loc.Page),
	}
}

func locationFromC(loc C.fz_location) Location {
	return Location{
		Chapter: int(loc.chapter),
		Page:    int(loc.page),
	}
}

// IsReflowable returns whether document supports reflow layout.
func (f *Document) IsReflowable() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return int(C.fz_is_document_reflowable(f.ctx, f.doc)) != 0
}

// Layout applies layout to reflowable documents.
func (f *Document) Layout(width, height, em float64) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	C.fz_layout_document(f.ctx, f.doc, C.float(width), C.float(height), C.float(em))
}

// CountChapters returns chapter count.
func (f *Document) CountChapters() int {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return int(C.fz_count_chapters(f.ctx, f.doc))
}

// CountChapterPages returns page count within a chapter.
func (f *Document) CountChapterPages(chapter int) int {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return int(C.fz_count_chapter_pages(f.ctx, f.doc, C.int(chapter)))
}

// LastLocation returns location of the last page.
func (f *Document) LastLocation() Location {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return locationFromC(C.fz_last_page(f.ctx, f.doc))
}

// NextLocation returns the next location from a given location.
func (f *Document) NextLocation(loc Location) Location {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return locationFromC(C.fz_next_page(f.ctx, f.doc, locationToC(loc)))
}

// PreviousLocation returns the previous location from a given location.
func (f *Document) PreviousLocation(loc Location) Location {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return locationFromC(C.fz_previous_page(f.ctx, f.doc, locationToC(loc)))
}

// ClampLocation clamps a location to valid chapter/page range.
func (f *Document) ClampLocation(loc Location) Location {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return locationFromC(C.fz_clamp_location(f.ctx, f.doc, locationToC(loc)))
}

// LocationFromPageNumber resolves chapter/page location from global page index.
func (f *Document) LocationFromPageNumber(number int) Location {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return locationFromC(C.fz_location_from_page_number(f.ctx, f.doc, C.int(number)))
}

// PageNumberFromLocation resolves global page index from chapter/page location.
func (f *Document) PageNumberFromLocation(loc Location) int {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return int(C.fz_page_number_from_location(f.ctx, f.doc, locationToC(loc)))
}

// MakeBookmark creates a bookmark for a location.
func (f *Document) MakeBookmark(loc Location) Bookmark {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return Bookmark(uintptr(C.fz_make_bookmark(f.ctx, f.doc, locationToC(loc))))
}

// LookupBookmark resolves a bookmark to location.
func (f *Document) LookupBookmark(mark Bookmark) Location {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return locationFromC(C.fz_lookup_bookmark(f.ctx, f.doc, C.fz_bookmark(mark)))
}

// ResolveLink resolves a URI to a location and optional coordinates.
func (f *Document) ResolveLink(uri string) (Location, float64, float64, bool) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	var (
		x   C.float
		y   C.float
		loc C.fz_location
	)
	withCString(uri, func(curi *C.char) {
		loc = C.fz_resolve_link(f.ctx, f.doc, curi, &x, &y)
	})
	out := locationFromC(loc)
	ok := out.Chapter >= 0 && out.Page >= 0
	return out, float64(x), float64(y), ok
}

// NumPage returns total number of pages in document.
func (f *Document) NumPage() int {
	return int(C.fz_count_pages(f.ctx, f.doc))
}

// Links returns slice of links for given page number.
func (f *Document) Links(pageNumber int) ([]Link, error) {
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

	links := C.fz_load_links(f.ctx, page)
	defer C.fz_drop_link(f.ctx, links)

	linkCount := 0
	for currLink := links; currLink != nil; currLink = currLink.next {
		linkCount++
	}

	if linkCount == 0 {
		return nil, nil
	}

	gLinks := make([]Link, linkCount)

	currLink := links
	for i := 0; i < linkCount; i++ {
		gLinks[i] = Link{
			URI: C.GoString(currLink.uri),
		}
		currLink = currLink.next
	}

	return gLinks, nil
}
