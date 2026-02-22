package gopherpdf

import (
	"fmt"
	"io"
	"strings"
	"sync"

	fitz "github.com/brainwhocodes/gopher-pdf/internal/fitz"
)

// Document provides a Go-facing API over MuPDF document operations.
type Document struct {
	mu     sync.Mutex
	doc    *fitz.Document
	name   string
	closed bool
}

// Open opens a document from a file path.
func Open(path string) (*Document, error) {
	return OpenWithPassword(path, "")
}

// OpenBytes opens a document from in-memory bytes.
func OpenBytes(b []byte) (*Document, error) {
	return OpenBytesWithPassword(b, "")
}

// OpenBytesWithFileType opens in-memory bytes using an explicit file type hint
// (for example "html", ".html", "xml", ".xml").
func OpenBytesWithFileType(b []byte, fileType string) (*Document, error) {
	return OpenBytesWithFileTypeAndPassword(b, fileType, "")
}

// OpenWithPassword opens a document from path with optional password.
func OpenWithPassword(path, password string) (*Document, error) {
	doc, err := fitz.New(path)
	return wrapOpenResult("open document", path, doc, err, password)
}

// OpenBytesWithPassword opens in-memory bytes with optional password.
func OpenBytesWithPassword(b []byte, password string) (*Document, error) {
	doc, err := fitz.NewFromMemory(b)
	return wrapOpenResult("open document from bytes", "", doc, err, password)
}

// OpenBytesWithFileTypeAndPassword opens in-memory bytes with optional
// file-type hint and password.
func OpenBytesWithFileTypeAndPassword(b []byte, fileType, password string) (*Document, error) {
	magic := strings.TrimSpace(fileType)
	if strings.HasPrefix(magic, ".") {
		magic = magic[1:]
	}
	doc, err := fitz.NewFromMemoryWithMagic(b, magic)
	return wrapOpenResult(fmt.Sprintf("open document from bytes with file type %q", fileType), "", doc, err, password)
}

// OpenReader opens a document from an io.Reader.
func OpenReader(r io.Reader) (*Document, error) {
	return OpenReaderWithPassword(r, "")
}

// OpenReaderWithFileType opens from reader bytes with an explicit file type hint.
func OpenReaderWithFileType(r io.Reader, fileType string) (*Document, error) {
	return OpenReaderWithFileTypeAndPassword(r, fileType, "")
}

// OpenReaderWithPassword opens from reader bytes with optional password.
func OpenReaderWithPassword(r io.Reader, password string) (*Document, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read document from reader: %w", err)
	}
	return OpenBytesWithPassword(b, password)
}

// OpenReaderWithFileTypeAndPassword opens from reader bytes with explicit file
// type hint and optional password.
func OpenReaderWithFileTypeAndPassword(r io.Reader, fileType, password string) (*Document, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read document from reader: %w", err)
	}
	return OpenBytesWithFileTypeAndPassword(b, fileType, password)
}

func wrapOpenResult(source, name string, doc *fitz.Document, err error, password string) (*Document, error) {
	doc, err = openFitzMaybeAuth(source, doc, err, password)
	if err != nil {
		return nil, err
	}
	return &Document{doc: doc, name: name}, nil
}

// Close releases document resources. It is safe to call multiple times.
func (d *Document) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return nil
	}
	d.closed = true
	return d.doc.Close()
}

// Name returns the document name/path when opened from file.
func (d *Document) Name() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.name
}

// IsClosed reports whether Close has been called.
func (d *Document) IsClosed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.closed || d.doc == nil
}
