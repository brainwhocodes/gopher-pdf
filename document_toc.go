package gopherpdf

import (
	"strings"
)

// ToC returns table of contents entries.
func (d *Document) ToC() ([]Outline, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return nil, err
	}
	items, err := d.doc.ToC()
	if err != nil {
		return nil, err
	}
	out := make([]Outline, 0, len(items))
	for _, item := range items {
		out = append(out, Outline{
			Level: item.Level,
			Title: item.Title,
			URI:   item.URI,
			Page:  item.Page,
			Top:   item.Top,
		})
	}
	return out, nil
}

// ToCSimple returns simplified TOC entries with 1-based page numbers.
func (d *Document) ToCSimple() ([]TOCEntry, error) {
	toc, err := d.ToC()
	if err != nil {
		return nil, err
	}
	out := make([]TOCEntry, 0, len(toc))
	for _, item := range toc {
		page := -1
		if item.Page >= 0 && strings.HasPrefix(item.URI, "#") {
			page = item.Page + 1
		}
		out = append(out, TOCEntry{
			Level: item.Level,
			Title: item.Title,
			Page:  page,
		})
	}
	return out, nil
}

// GetToC is a PyMuPDF-style convenience method.
// If simple is true it returns []TOCEntry, otherwise []Outline.
func (d *Document) GetToC(simple bool) (any, error) {
	if simple {
		return d.ToCSimple()
	}
	return d.ToC()
}
