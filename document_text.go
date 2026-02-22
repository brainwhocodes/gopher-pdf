package gopherpdf

import (
	"encoding/json"
	"fmt"
	"strings"

	fitz "github.com/brainwhocodes/gopher-pdf/internal/fitz"
)

// Text returns plain text for a page.
func (d *Document) Text(pageNumber int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.Text(pageNumber)
}

// TextWithFlags returns plain text for a page with MuPDF text extraction flags.
func (d *Document) TextWithFlags(pageNumber int, flags int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.TextWithFlags(pageNumber, flags)
}

// HTML returns HTML for a page.
func (d *Document) HTML(pageNumber int, header bool) (string, error) {
	return d.HTMLWithFlags(pageNumber, header, TextFlagsHTML)
}

// HTMLWithFlags returns HTML for a page with MuPDF text extraction flags.
func (d *Document) HTMLWithFlags(pageNumber int, header bool, flags int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.HTMLWithFlags(pageNumber, header, flags)
}

// XHTML returns XHTML for a page.
func (d *Document) XHTML(pageNumber int, header bool) (string, error) {
	return d.XHTMLWithFlags(pageNumber, header, TextFlagsXHTML)
}

// XHTMLWithFlags returns XHTML for a page with MuPDF text extraction flags.
func (d *Document) XHTMLWithFlags(pageNumber int, header bool, flags int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.XHTMLWithFlags(pageNumber, header, flags)
}

// XML returns XML for a page.
func (d *Document) XML(pageNumber int) (string, error) {
	return d.XMLWithFlags(pageNumber, TextFlagsXML)
}

// XMLWithFlags returns XML for a page with MuPDF text extraction flags.
func (d *Document) XMLWithFlags(pageNumber int, flags int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.XMLWithFlags(pageNumber, flags)
}

// JSON returns structured JSON text extraction for a page.
func (d *Document) JSON(pageNumber int) (string, error) {
	return d.JSONWithFlags(pageNumber, TextFlagsDict)
}

// JSONWithFlags returns structured JSON text extraction for a page with MuPDF text flags.
func (d *Document) JSONWithFlags(pageNumber int, flags int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.JSONWithFlags(pageNumber, flags, 1.0)
}

// SearchCount returns number of text hits for a needle on a page.
func (d *Document) SearchCount(pageNumber int, needle string, flags int, maxHits int) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return 0, err
	}
	return d.doc.SearchCount(pageNumber, needle, flags, maxHits)
}

// Search returns rectangles of text hits on a page.
// If clip is non-nil, only hits fully inside clip are returned.
func (d *Document) Search(pageNumber int, needle string, flags int, maxHits int, clip *RectF) ([]RectF, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return nil, err
	}

	var fitzClip *fitz.Rect
	if clip != nil {
		fitzClip = &fitz.Rect{
			X0: clip.X0,
			Y0: clip.Y0,
			X1: clip.X1,
			Y1: clip.Y1,
		}
	}

	hits, err := d.doc.Search(pageNumber, needle, flags, maxHits, fitzClip)
	if err != nil {
		return nil, err
	}
	out := make([]RectF, 0, len(hits))
	for _, hit := range hits {
		out = append(out, RectF{
			X0: hit.X0,
			Y0: hit.Y0,
			X1: hit.X1,
			Y1: hit.Y1,
		})
	}
	return out, nil
}

// TextBox returns text inside a rectangular clip area on a page.
func (d *Document) TextBox(pageNumber int, area RectF, flags int) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.ensurePageLocked(pageNumber); err != nil {
		return "", err
	}
	return d.doc.TextBox(pageNumber, fitz.Rect{
		X0: area.X0,
		Y0: area.Y0,
		X1: area.X1,
		Y1: area.Y1,
	}, flags)
}

// GetText provides a PyMuPDF-like mode-based text extraction API.
//
// Supported modes:
// - `text`
// - `words` (word tokens)
// - `blocks` (raw blocks from JSON extraction)
// - `html`
// - `xhtml`
// - `xml`
// - `json`
// - `dict` (parsed json map)
// - `rawjson`
// - `rawdict` (parsed json map)
func (d *Document) GetText(pageNumber int, mode string, flags int) (any, error) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if flags == 0 {
		flags = defaultTextFlagsForMode(mode)
	}

	switch mode {
	case "", "text":
		return d.TextWithFlags(pageNumber, flags)
	case "words":
		dictAny, err := d.GetText(pageNumber, "dict", flags)
		if err != nil {
			return nil, err
		}
		dict, _ := dictAny.(map[string]any)
		blocksAny, _ := dict["blocks"].([]any)
		words := make([]string, 0, 128)
		for _, blockAny := range blocksAny {
			block, ok := blockAny.(map[string]any)
			if !ok {
				continue
			}
			linesAny, _ := block["lines"].([]any)
			for _, lineAny := range linesAny {
				line, ok := lineAny.(map[string]any)
				if !ok {
					continue
				}
				text, _ := line["text"].(string)
				words = append(words, strings.Fields(text)...)
			}
		}
		return words, nil
	case "blocks":
		dictAny, err := d.GetText(pageNumber, "dict", flags)
		if err != nil {
			return nil, err
		}
		dict, _ := dictAny.(map[string]any)
		blocksAny, _ := dict["blocks"].([]any)
		return blocksAny, nil
	case "html":
		return d.HTMLWithFlags(pageNumber, true, flags)
	case "xhtml":
		return d.XHTMLWithFlags(pageNumber, true, flags)
	case "xml":
		return d.XMLWithFlags(pageNumber, flags)
	case "json", "rawjson":
		return d.JSONWithFlags(pageNumber, flags)
	case "dict", "rawdict":
		text, err := d.JSONWithFlags(pageNumber, flags)
		if err != nil {
			return nil, err
		}
		var out map[string]any
		if err := json.Unmarshal([]byte(text), &out); err != nil {
			return nil, fmt.Errorf("decode json text mode %q: %w", mode, err)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported text mode %q", mode)
	}
}

func defaultTextFlagsForMode(mode string) int {
	switch mode {
	case "words":
		return TextFlagsWords
	case "blocks":
		return TextFlagsBlocks
	case "dict":
		return TextFlagsDict
	case "rawdict":
		return TextFlagsRawDict
	case "json":
		return TextFlagsDict
	case "rawjson":
		return TextFlagsRawDict
	case "html":
		return TextFlagsHTML
	case "xhtml":
		return TextFlagsXHTML
	case "xml":
		return TextFlagsXML
	case "", "text":
		return TextFlagsText
	default:
		return 0
	}
}

// GetPageText is a PyMuPDF-style convenience alias for GetText.
func (d *Document) GetPageText(pageNumber int, mode string, flags int) (any, error) {
	return d.GetText(pageNumber, mode, flags)
}

// SearchPageFor is a PyMuPDF-style convenience alias for Search.
func (d *Document) SearchPageFor(pageNumber int, needle string, flags int, maxHits int, clip *RectF) ([]RectF, error) {
	return d.Search(pageNumber, needle, flags, maxHits, clip)
}
