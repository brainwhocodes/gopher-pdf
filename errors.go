package gopherpdf

import "errors"

var (
	// ErrNoPDF is returned for PDF-only APIs called on non-PDF documents.
	ErrNoPDF = errors.New("is no PDF")

	// ErrDocumentClosed is returned when an operation is attempted after Close.
	ErrDocumentClosed = errors.New("document is closed")
	// ErrPageOutOfRange is returned when a page index is outside [0, NumPages).
	ErrPageOutOfRange = errors.New("page index out of range")
	// ErrChapterOutOfRange is returned when a chapter index is outside [0, NumChapters).
	ErrChapterOutOfRange = errors.New("chapter index out of range")
)
