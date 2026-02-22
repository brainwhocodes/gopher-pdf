package gopherpdf

import (
	"fmt"

	fitz "github.com/brainwhocodes/gopher-pdf/internal/fitz"
)

// NumPages returns total number of pages.
func (d *Document) NumPages() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return 0, err
	}
	return d.doc.NumPage(), nil
}

// PageCount is a convenience alias for NumPages.
func (d *Document) PageCount() (int, error) {
	return d.NumPages()
}

// IsReflowable reports whether the document supports reflow layout operations.
func (d *Document) IsReflowable() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.IsReflowable(), nil
}

// Layout applies reflow layout parameters for reflowable documents.
func (d *Document) Layout(width, height, em float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return err
	}
	if width <= 0 {
		width = 595 // A4 width at 72 DPI.
	}
	if height <= 0 {
		height = 842 // A4 height at 72 DPI.
	}
	if em <= 0 {
		em = 12
	}
	d.doc.Layout(width, height, em)
	return nil
}

// NumChapters returns chapter count for the document.
func (d *Document) NumChapters() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return 0, err
	}
	return d.doc.CountChapters(), nil
}

// ChapterCount is a convenience alias for NumChapters.
func (d *Document) ChapterCount() (int, error) {
	return d.NumChapters()
}

// ChapterPageCount returns page count for a chapter.
func (d *Document) ChapterPageCount(chapter int) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return 0, err
	}
	chapters := d.doc.CountChapters()
	if chapter < 0 || chapter >= chapters {
		return 0, fmt.Errorf("%w: %d", ErrChapterOutOfRange, chapter)
	}
	return d.doc.CountChapterPages(chapter), nil
}

// LastLocation returns the last valid chapter/page location.
func (d *Document) LastLocation() (Location, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return Location{}, err
	}
	loc := d.doc.LastLocation()
	return Location{Chapter: loc.Chapter, Page: loc.Page}, nil
}

// NextLocation returns next location from a given location.
func (d *Document) NextLocation(loc Location) (Location, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return Location{}, err
	}
	out := d.doc.NextLocation(fitz.Location{Chapter: loc.Chapter, Page: loc.Page})
	return Location{Chapter: out.Chapter, Page: out.Page}, nil
}

// PrevLocation returns previous location from a given location.
func (d *Document) PrevLocation(loc Location) (Location, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return Location{}, err
	}
	out := d.doc.PreviousLocation(fitz.Location{Chapter: loc.Chapter, Page: loc.Page})
	return Location{Chapter: out.Chapter, Page: out.Page}, nil
}

// ClampLocation clamps a location to valid chapter/page range.
func (d *Document) ClampLocation(loc Location) (Location, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return Location{}, err
	}
	out := d.doc.ClampLocation(fitz.Location{Chapter: loc.Chapter, Page: loc.Page})
	return Location{Chapter: out.Chapter, Page: out.Page}, nil
}

// LocationFromPageNumber resolves chapter/page location from global page index.
func (d *Document) LocationFromPageNumber(pageNumber int) (Location, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return Location{}, err
	}
	pages := d.doc.NumPage()
	if pageNumber < 0 || pageNumber >= pages {
		return Location{}, fmt.Errorf("%w: %d", ErrPageOutOfRange, pageNumber)
	}
	out := d.doc.LocationFromPageNumber(pageNumber)
	return Location{Chapter: out.Chapter, Page: out.Page}, nil
}

// PageNumberFromLocation resolves global page index from chapter/page tuple.
func (d *Document) PageNumberFromLocation(loc Location) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return 0, err
	}
	return d.doc.PageNumberFromLocation(fitz.Location{Chapter: loc.Chapter, Page: loc.Page}), nil
}

// MakeBookmark creates bookmark token for a location.
func (d *Document) MakeBookmark(loc Location) (Bookmark, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return 0, err
	}
	mark := d.doc.MakeBookmark(fitz.Location{Chapter: loc.Chapter, Page: loc.Page})
	return Bookmark(uint64(mark)), nil
}

// LookupBookmark resolves bookmark token into a location.
func (d *Document) LookupBookmark(mark Bookmark) (Location, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return Location{}, err
	}
	out := d.doc.LookupBookmark(fitz.Bookmark(mark))
	return Location{Chapter: out.Chapter, Page: out.Page}, nil
}

// FindBookmark is a convenience alias for LookupBookmark.
func (d *Document) FindBookmark(mark Bookmark) (Location, error) {
	return d.LookupBookmark(mark)
}
