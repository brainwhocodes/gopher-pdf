package gopherpdf

import (
	"fmt"
	"strings"
)

func (d *Document) ensureOpenLocked() error {
	if d.closed || d.doc == nil {
		return ErrDocumentClosed
	}
	return nil
}

func (d *Document) ensurePageLocked(pageNumber int) error {
	if err := d.ensureOpenLocked(); err != nil {
		return err
	}
	if pageNumber < 0 || pageNumber >= d.doc.NumPage() {
		return fmt.Errorf("%w: %d", ErrPageOutOfRange, pageNumber)
	}
	return nil
}

func (d *Document) ensurePDFLocked() error {
	if err := d.ensureOpenLocked(); err != nil {
		return err
	}
	format := strings.ToUpper(strings.TrimSpace(d.doc.Metadata()["format"]))
	if !strings.HasPrefix(format, "PDF") {
		return ErrNoPDF
	}
	return nil
}
