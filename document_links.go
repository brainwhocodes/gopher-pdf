package gopherpdf

// ResolveLink resolves an internal document URI.
func (d *Document) ResolveLink(uri string) (ResolvedLink, bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensureOpenLocked(); err != nil {
		return ResolvedLink{}, false, err
	}
	loc, x, y, ok := d.doc.ResolveLink(uri)
	out := ResolvedLink{
		Location: Location{Chapter: loc.Chapter, Page: loc.Page},
		X:        x,
		Y:        y,
	}
	return out, ok, nil
}

// Links returns link targets for a page.
func (d *Document) Links(pageNumber int) ([]Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return nil, err
	}
	links, err := d.doc.Links(pageNumber)
	if err != nil {
		return nil, err
	}
	out := make([]Link, 0, len(links))
	for _, link := range links {
		out = append(out, Link{URI: link.URI})
	}
	return out, nil
}

// HasLinks reports whether any page contains at least one link.
func (d *Document) HasLinks() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return false, err
	}
	pages := d.doc.NumPage()
	for i := 0; i < pages; i++ {
		links, err := d.doc.Links(i)
		if err != nil {
			return false, err
		}
		if len(links) > 0 {
			return true, nil
		}
	}
	return false, nil
}

// HasAnnotations reports whether any page has non-link, non-widget annotations.
func (d *Document) HasAnnotations() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePDFLocked(); err != nil {
		return false, err
	}
	pages := d.doc.NumPage()
	for i := 0; i < pages; i++ {
		has, err := d.doc.PageHasAnnots(i)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}
