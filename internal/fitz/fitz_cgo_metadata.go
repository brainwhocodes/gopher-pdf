//go:build cgo && !nocgo

package fitz

/*
#include "fitz_cgo.h"
*/
import "C"

import (
	"bytes"
	"unsafe"
)

// Metadata returns the map with standard metadata.
func (f *Document) Metadata() map[string]string {
	data := make(map[string]string)

	lookup := func(key string) string {
		buf := make([]byte, 256)
		withCString(key, func(ckey *C.char) {
			C.fz_lookup_metadata(f.ctx, f.doc, ckey, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)))
		})
		if n := bytes.IndexByte(buf, 0); n >= 0 {
			buf = buf[:n]
		}
		return string(buf)
	}

	data["format"] = lookup("format")
	data["encryption"] = lookup("encryption")
	data["title"] = lookup("info:Title")
	data["author"] = lookup("info:Author")
	data["subject"] = lookup("info:Subject")
	data["keywords"] = lookup("info:Keywords")
	data["creator"] = lookup("info:Creator")
	data["producer"] = lookup("info:Producer")
	data["creationDate"] = lookup("info:CreationDate")
	data["modDate"] = lookup("info:ModDate")
	data["trapped"] = lookup("info:Trapped")

	return data
}

// SetMetadata updates standard document metadata fields.
func (f *Document) SetMetadata(values map[string]string) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if len(values) == 0 {
		pdf := C.pdf_document_from_fz_document(f.ctx, f.doc)
		if pdf != nil {
			trailer := C.pdf_trailer(f.ctx, pdf)
			if trailer != nil {
				withCString("Info", func(cinfo *C.char) {
					_ = C.safe_pdf_dict_del_path(f.ctx, trailer, cinfo)
				})
			}
		}
		return
	}

	getValue := func(logicalKey, rawKey string) string {
		if values == nil {
			return ""
		}
		if v, ok := values[logicalKey]; ok {
			return v
		}
		if v, ok := values[rawKey]; ok {
			return v
		}
		return ""
	}

	set := func(key, value string) {
		withCString(key, func(ckey *C.char) {
			withCString(value, func(cvalue *C.char) {
				C.fz_set_metadata(f.ctx, f.doc, ckey, cvalue)
			})
		})
	}

	set("info:Title", getValue("title", "info:Title"))
	set("info:Author", getValue("author", "info:Author"))
	set("info:Subject", getValue("subject", "info:Subject"))
	set("info:Keywords", getValue("keywords", "info:Keywords"))
	set("info:Creator", getValue("creator", "info:Creator"))
	set("info:Producer", getValue("producer", "info:Producer"))
	set("info:CreationDate", getValue("creationDate", "info:CreationDate"))
	set("info:ModDate", getValue("modDate", "info:ModDate"))
	set("info:Trapped", getValue("trapped", "info:Trapped"))
}
