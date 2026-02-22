//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"unsafe"
)

// Document represents fitz document.
type Document struct {
	ctx    *C.struct_fz_context
	data   []byte // binds data to the Document lifecycle avoiding premature GC
	doc    *C.struct_fz_document
	mtx    sync.Mutex
	stream *C.fz_stream
}

// New returns new fitz document.
func New(filename string) (f *Document, err error) {
	f = &Document{}

	filename, err = filepath.Abs(filename)
	if err != nil {
		return
	}

	if _, e := os.Stat(filename); e != nil {
		err = ErrNoSuchFile
		return
	}

	f.ctx = (*C.struct_fz_context)(unsafe.Pointer(C.fz_new_context_imp(nil, nil, C.store(MaxStore), C.fz_version)))
	if f.ctx == nil {
		err = ErrCreateContext
		return
	}

	C.fz_register_document_handlers(f.ctx)

	withCString(filename, func(cfilename *C.char) {
		f.doc = C.open_document(f.ctx, cfilename)
	})
	if f.doc == nil {
		err = ErrOpenDocument
		return
	}

	ret := C.fz_needs_password(f.ctx, f.doc)
	v := int(ret) != 0
	if v {
		err = ErrNeedsPassword
	}

	return
}

// NewFromMemory returns new fitz document from byte slice.
func NewFromMemory(b []byte) (f *Document, err error) {
	return NewFromMemoryWithMagic(b, "")
}

// NewFromMemoryWithMagic returns new fitz document from byte slice and optional magic file type.
func NewFromMemoryWithMagic(b []byte, magic string) (f *Document, err error) {
	if len(b) == 0 {
		return nil, ErrEmptyBytes
	}
	f = &Document{}

	f.ctx = (*C.struct_fz_context)(unsafe.Pointer(C.fz_new_context_imp(nil, nil, C.store(MaxStore), C.fz_version)))
	if f.ctx == nil {
		err = ErrCreateContext
		return
	}

	C.fz_register_document_handlers(f.ctx)

	f.stream = C.fz_open_memory(f.ctx, (*C.uchar)(&b[0]), C.size_t(len(b)))
	if f.stream == nil {
		err = ErrOpenMemory
		return
	}

	if magic == "" {
		magic = contentType(b)
	}
	if magic == "" {
		err = ErrOpenMemory
		return
	}

	f.data = b


	withCString(magic, func(cmagic *C.char) {
		f.doc = C.open_document_with_stream(f.ctx, cmagic, f.stream)
	})
	if f.doc == nil {
		err = ErrOpenDocument
	}

	ret := C.fz_needs_password(f.ctx, f.doc)
	v := int(ret) != 0
	if v {
		err = ErrNeedsPassword
	}

	return
}

// NewFromReader returns new fitz document from io.Reader.
func NewFromReader(r io.Reader) (f *Document, err error) {
	b, e := io.ReadAll(r)
	if e != nil {
		err = e
		return
	}

	f, err = NewFromMemory(b)
	return
}

// NeedsPassword reports whether the document requires a password.
func (f *Document) NeedsPassword() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return int(C.fz_needs_password(f.ctx, f.doc)) != 0
}

// AuthenticatePassword attempts to authenticate document password.
func (f *Document) AuthenticatePassword(password string) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	ok := false
	withCString(password, func(cpassword *C.char) {
		ok = int(C.fz_authenticate_password(f.ctx, f.doc, cpassword)) != 0
	})
	return ok
}

// HasPermission reports whether the given permission flag is granted.
func (f *Document) HasPermission(permission int) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return int(C.fz_has_permission(f.ctx, f.doc, C.fz_permission(permission))) != 0
}

// Close closes the underlying fitz document.
func (f *Document) Close() error {
	if f.stream != nil {
		C.fz_drop_stream(f.ctx, f.stream)
	}

	C.fz_drop_document(f.ctx, f.doc)
	C.fz_drop_context(f.ctx)

	f.data = nil

	return nil
}
