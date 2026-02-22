package gopherpdf

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"time"

	fitz "bank-statments/backend/pkg/gopher-pdf/internal/fitz"
)

// MuPDFRenderer rasterizes PDF pages using a Go MuPDF binding.
// This implementation has no Python runtime dependency.
type MuPDFRenderer struct{}

func NewMuPDFRenderer() *MuPDFRenderer {
	return &MuPDFRenderer{}
}

func openMemoryWithPassword(pdf []byte, password string) (*fitz.Document, error) {
	doc, err := fitz.NewFromMemory(pdf)
	return openFitzMaybeAuth("open pdf from memory", doc, err, password)
}

// CountPages returns an exact page count from PDF bytes.
func (r *MuPDFRenderer) CountPages(_ context.Context, pdf []byte, password string) (int, error) {
	if len(pdf) == 0 {
		return 0, fmt.Errorf("pdf bytes are empty")
	}
	doc, err := openMemoryWithPassword(pdf, password)
	if err != nil {
		return 0, err
	}
	defer doc.Close()

	count := doc.NumPage()
	if count < 1 {
		return 0, fmt.Errorf("renderer returned invalid page count %d", count)
	}
	return count, nil
}

// Render rasterizes all pages as PNG bytes at the requested DPI.
func (r *MuPDFRenderer) Render(ctx context.Context, pdf []byte, opts RenderOptions) ([]RenderedPage, error) {
	if len(pdf) == 0 {
		return nil, fmt.Errorf("pdf bytes are empty")
	}
	opts = opts.Defaults()

	execCtx := ctx
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	doc, err := openMemoryWithPassword(pdf, opts.Password)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	pageCount := doc.NumPage()
	if pageCount < 1 {
		return nil, fmt.Errorf("renderer returned invalid page count %d", pageCount)
	}

	out := make([]RenderedPage, 0, pageCount)
	for i := 0; i < pageCount; i++ {
		select {
		case <-execCtx.Done():
			return nil, execCtx.Err()
		default:
		}

		started := time.Now()
		pngBytes, err := doc.ImagePNG(i, float64(opts.DPI))
		if err != nil {
			return nil, fmt.Errorf("render page %d: %w", i+1, err)
		}

		cfg, err := png.DecodeConfig(bytes.NewReader(pngBytes))
		if err != nil {
			return nil, fmt.Errorf("decode rendered page %d: %w", i+1, err)
		}

		out = append(out, RenderedPage{
			PageNumber: i + 1,
			Width:      cfg.Width,
			Height:     cfg.Height,
			DPI:        float64(opts.DPI),
			RenderMS:   int(time.Since(started).Milliseconds()),
			PNG:        pngBytes,
		})
	}

	return out, nil
}
