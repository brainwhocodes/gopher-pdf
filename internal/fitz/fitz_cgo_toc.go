//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

// ToC returns the table of contents (also known as outline).
func (f *Document) ToC() ([]Outline, error) {
	data := make([]Outline, 0)

	outline := C.load_outline(f.ctx, f.doc)
	if outline == nil {
		return []Outline{}, nil
	}
	defer C.fz_drop_outline(f.ctx, outline)

	var walk func(outline *C.fz_outline, level int)

	walk = func(outline *C.fz_outline, level int) {
		for outline != nil {
			res := Outline{}
			res.Level = level
			res.Title = C.GoString(outline.title)
			res.URI = C.GoString(outline.uri)
			res.Page = int(outline.page.page)
			res.Top = float64(outline.y)
			data = append(data, res)

			if outline.down != nil {
				walk(outline.down, level+1)
			}
			outline = outline.next
		}
	}

	walk(outline, 1)
	return data, nil
}

