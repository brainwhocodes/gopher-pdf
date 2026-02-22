package gopherpdf

import (
	"strings"
)

// Metadata returns document metadata.
func (d *Document) Metadata() (map[string]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return nil, err
	}
	return d.doc.Metadata(), nil
}

// SetMetadata updates document metadata keys supported by MuPDF.
// Standard logical keys: title, author, subject, keywords, creator, producer,
// creationDate, modDate, trapped.
func (d *Document) SetMetadata(values map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return err
	}
	d.doc.SetMetadata(values)
	return nil
}

// IsPDF reports whether the opened document format is PDF.
func (d *Document) IsPDF() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	format := strings.ToUpper(strings.TrimSpace(d.doc.Metadata()["format"]))
	return strings.HasPrefix(format, "PDF"), nil
}

// IsEncrypted reports whether the document is currently encrypted/locked.
// This mirrors PyMuPDF's "is_encrypted" behavior after authentication state.
func (d *Document) IsEncrypted() (bool, error) {
	return d.NeedsPassword()
}
