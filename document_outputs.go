package gopherpdf

// SVG returns SVG for a page.
func (d *Document) SVG(pageNumber int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.SVG(pageNumber)
}
