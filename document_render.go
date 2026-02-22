package gopherpdf

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

// Image renders a page at default DPI as RGBA pixels.
func (d *Document) Image(pageNumber int) (*image.RGBA, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return nil, err
	}
	return d.doc.Image(pageNumber)
}

// ImageDPI renders a page as RGBA pixels at the given DPI.
func (d *Document) ImageDPI(pageNumber int, dpi float64) (*image.RGBA, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return nil, err
	}
	if dpi <= 0 {
		dpi = 200
	}
	return d.doc.ImageDPI(pageNumber, dpi)
}

// GetPagePixmap is a PyMuPDF-style convenience alias for page rendering.
// It returns an RGBA pixel buffer at the requested DPI.
func (d *Document) GetPagePixmap(pageNumber int, dpi float64) (*image.RGBA, error) {
	return d.ImageDPI(pageNumber, dpi)
}

// RenderPagePNG renders a page to PNG bytes at the given DPI.
func (d *Document) RenderPagePNG(pageNumber int, dpi float64) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return nil, err
	}
	if dpi <= 0 {
		dpi = 200
	}
	return d.doc.ImagePNG(pageNumber, dpi)
}

// RenderPage renders a page and includes dimensions.
func (d *Document) RenderPage(pageNumber int, dpi float64) (RenderedPage, error) {
	pngBytes, err := d.RenderPagePNG(pageNumber, dpi)
	if err != nil {
		return RenderedPage{}, err
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(pngBytes))
	if err != nil {
		return RenderedPage{}, fmt.Errorf("decode rendered png: %w", err)
	}
	if dpi <= 0 {
		dpi = 200
	}
	return RenderedPage{
		PageNumber: pageNumber,
		Width:      cfg.Width,
		Height:     cfg.Height,
		DPI:        dpi,
		RenderMS:   0,
		PNG:        pngBytes,
	}, nil
}
