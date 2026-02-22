package gopherpdf

import (
	"image"
	"strings"

	fitz "bank-statments/backend/pkg/gopher-pdf/internal/fitz"
)

// NeedsPassword reports whether the opened document still requires password auth.
func (d *Document) NeedsPassword() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.NeedsPassword(), nil
}

// HasPermission reports whether a document permission is granted.
func (d *Document) HasPermission(permission Permission) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.HasPermission(int(permission)), nil
}

// IsDirty reports whether a PDF document has unsaved changes.
func (d *Document) IsDirty() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.IsDirty(), nil
}

// CanSaveIncrementally reports whether incremental save is possible for a PDF.
func (d *Document) CanSaveIncrementally() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.CanSaveIncrementally(), nil
}

// IsFastWebAccess reports whether a PDF is linearized.
func (d *Document) IsFastWebAccess() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.IsFastWebAccess(), nil
}

// IsRepaired reports whether a PDF required repair on open.
func (d *Document) IsRepaired() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return false, err
	}
	return d.doc.IsRepaired(), nil
}

// PageLabel returns display label for a page.
func (d *Document) PageLabel(pageNumber int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.PageLabel(pageNumber)
}

// PageNumbersByLabel returns all page numbers matching a page label.
func (d *Document) PageNumbersByLabel(label string) ([]int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return nil, err
	}
	want := strings.TrimSpace(label)
	if want == "" {
		return []int{}, nil
	}
	total := d.doc.NumPage()
	out := make([]int, 0, 1)
	for i := 0; i < total; i++ {
		got, err := d.doc.PageLabel(i)
		if err != nil {
			return nil, err
		}
		if got == want {
			out = append(out, i)
		}
	}
	return out, nil
}

// PageLabels returns per-page display labels for the whole document.
func (d *Document) PageLabels() ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return nil, err
	}
	total := d.doc.NumPage()
	out := make([]string, 0, total)
	for i := 0; i < total; i++ {
		label, err := d.doc.PageLabel(i)
		if err != nil {
			return nil, err
		}
		out = append(out, label)
	}
	return out, nil
}

// GetPageLabels is a convenience alias for PageLabels.
func (d *Document) GetPageLabels() ([]string, error) {
	return d.PageLabels()
}

// GetPageNumbers is a convenience alias for PageNumbersByLabel.
func (d *Document) GetPageNumbers(label string) ([]int, error) {
	return d.PageNumbersByLabel(label)
}

// PageLabelRules returns explicit PDF page-label definition rules.
func (d *Document) PageLabelRules() ([]PageLabelRule, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return nil, err
	}

	raw, err := d.doc.PageLabelRules()
	if err != nil {
		return nil, err
	}
	out := make([]PageLabelRule, 0, len(raw))
	for _, rule := range raw {
		out = append(out, PageLabelRule{
			StartPage:    rule.StartPage,
			Prefix:       rule.Prefix,
			Style:        rule.Style,
			FirstPageNum: rule.FirstPageNum,
		})
	}
	return out, nil
}

// GetPageLabelRules is a convenience alias for PageLabelRules.
func (d *Document) GetPageLabelRules() ([]PageLabelRule, error) {
	return d.PageLabelRules()
}

// SetPageLabels replaces PDF page-label definition rules.
func (d *Document) SetPageLabels(rules []PageLabelRule) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return err
	}

	raw := make([]fitz.PageLabelRule, 0, len(rules))
	for _, rule := range rules {
		raw = append(raw, fitz.PageLabelRule{
			StartPage:    rule.StartPage,
			Prefix:       rule.Prefix,
			Style:        rule.Style,
			FirstPageNum: rule.FirstPageNum,
		})
	}
	return d.doc.SetPageLabels(raw)
}

// ResolveNames returns PDF name-destination entries.
func (d *Document) ResolveNames() (map[string]NamedDestination, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return nil, err
	}

	raw := d.doc.ResolveNames()
	out := make(map[string]NamedDestination, len(raw))
	for k, v := range raw {
		out[k] = NamedDestination{
			Page: v.Page,
			To:   [2]float64{v.X, v.Y},
			Zoom: v.Zoom,
		}
	}
	return out, nil
}

// XrefLength returns the PDF xref table length.
func (d *Document) XrefLength() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return 0, err
	}
	return d.doc.XrefLength(), nil
}

// XrefObject returns a printable PDF object for xref number.
// Use xref=-1 for the trailer object.
func (d *Document) XrefObject(xref int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return "", err
	}
	return d.doc.XrefObject(xref)
}

// XrefGetKey returns a key value from trailer/object dictionary.
func (d *Document) XrefGetKey(xref int, key string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return "", err
	}
	return d.doc.XrefGetKey(xref, key)
}

// XrefGetKeyTyped returns a key value as (type, value), mirroring PyMuPDF xref_get_key.
func (d *Document) XrefGetKeyTyped(xref int, key string) (string, string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return "", "", err
	}
	return d.doc.XrefGetKeyTyped(xref, key)
}

// XrefGetKeyValue returns a typed key-value pair for xref dictionary lookup.
func (d *Document) XrefGetKeyValue(xref int, key string) (XrefKeyValue, error) {
	typ, value, err := d.XrefGetKeyTyped(xref, key)
	if err != nil {
		return XrefKeyValue{}, err
	}
	return XrefKeyValue{Type: typ, Value: value}, nil
}

// XrefSetKey sets a key to a PDF object expression on trailer/object xref.
func (d *Document) XrefSetKey(xref int, key, value string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return err
	}
	return d.doc.XrefSetKey(xref, key, value)
}

// XrefIsFont reports whether xref is a font dictionary object.
func (d *Document) XrefIsFont(xref int) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return false, err
	}
	return d.doc.XrefIsFont(xref)
}

// XrefIsImage reports whether xref is an image object.
func (d *Document) XrefIsImage(xref int) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return false, err
	}
	return d.doc.XrefIsImage(xref)
}

// XrefIsForm reports whether xref is a Form XObject.
func (d *Document) XrefIsForm(xref int) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return false, err
	}
	return d.doc.XrefIsForm(xref)
}

// XrefIsStream reports whether xref contains a stream.
func (d *Document) XrefIsStream(xref int) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return false, err
	}
	return d.doc.XrefIsStream(xref)
}

// XrefStream returns decompressed xref stream bytes.
func (d *Document) XrefStream(xref int) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return nil, err
	}
	return d.doc.XrefStream(xref)
}

// XrefStreamRaw returns raw (compressed) xref stream bytes.
func (d *Document) XrefStreamRaw(xref int) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return nil, err
	}
	return d.doc.XrefStreamRaw(xref)
}

// XrefGetKeys returns dictionary keys for trailer/object xref.
func (d *Document) XrefGetKeys(xref int) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return nil, err
	}
	return d.doc.XrefGetKeys(xref)
}

// CatalogXref returns the PDF catalog object xref number.
func (d *Document) CatalogXref() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return 0, err
	}
	return d.doc.CatalogXref(), nil
}

// PageXref returns page object xref number by 0-based page index.
func (d *Document) PageXref(pageNumber int) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return 0, err
	}
	if err := d.ensurePDFLocked(); err != nil {
		return 0, err
	}
	return d.doc.PageXref(pageNumber), nil
}

// FormFieldCount returns AcroForm field count. Returns 0 for non-form PDFs.
func (d *Document) FormFieldCount() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return 0, err
	}
	n := d.doc.FormFieldCount()
	if n < 0 {
		return 0, nil
	}
	return n, nil
}

// IsFormPDF reports whether this PDF has at least one AcroForm field.
func (d *Document) IsFormPDF() (bool, error) {
	n, err := d.FormFieldCount()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Bound returns page bounds in pixels at native document units.
func (d *Document) Bound(pageNumber int) (image.Rectangle, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return image.Rectangle{}, err
	}
	return d.doc.Bound(pageNumber)
}

// BoundBox returns bounds for the selected page box.
func (d *Document) BoundBox(pageNumber int, box PageBox) (image.Rectangle, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return image.Rectangle{}, err
	}
	return d.doc.BoundBox(pageNumber, int(box))
}
