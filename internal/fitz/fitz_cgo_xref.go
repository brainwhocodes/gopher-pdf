//go:build cgo && !nocgo

package fitz

/*
#include <mupdf/fitz.h>
#include <stdlib.h>

typedef struct pdf_document pdf_document;
typedef struct pdf_obj pdf_obj;

pdf_document *pdf_document_from_fz_document(fz_context *ctx, fz_document *doc);
pdf_obj *pdf_trailer(fz_context *ctx, pdf_document *doc);
int pdf_xref_len(fz_context *ctx, pdf_document *doc);
pdf_obj *safe_pdf_load_object(fz_context *ctx, pdf_document *pdf, int num);
void pdf_drop_obj(fz_context *ctx, pdf_obj *obj);
int safe_pdf_print_obj(fz_context *ctx, fz_output *out, pdf_obj *obj);
int pdf_is_indirect(fz_context *ctx, pdf_obj *obj);
pdf_obj *pdf_resolve_indirect(fz_context *ctx, pdf_obj *ref);
int pdf_is_dict(fz_context *ctx, pdf_obj *obj);
pdf_obj *pdf_dict_getp(fz_context *ctx, pdf_obj *dict, const char *path);
int pdf_is_null(fz_context *ctx, pdf_obj *obj);
int pdf_is_array(fz_context *ctx, pdf_obj *obj);
int pdf_is_int(fz_context *ctx, pdf_obj *obj);
int pdf_is_real(fz_context *ctx, pdf_obj *obj);
int pdf_is_bool(fz_context *ctx, pdf_obj *obj);
int pdf_is_name(fz_context *ctx, pdf_obj *obj);
int pdf_is_string(fz_context *ctx, pdf_obj *obj);
int pdf_to_int(fz_context *ctx, pdf_obj *obj);
int pdf_to_num(fz_context *ctx, pdf_obj *obj);
int pdf_to_bool(fz_context *ctx, pdf_obj *obj);
const char *pdf_to_name(fz_context *ctx, pdf_obj *obj);
const char *pdf_to_text_string(fz_context *ctx, pdf_obj *obj);
int safe_pdf_dict_putp_parsed(fz_context *ctx, pdf_document *pdf, pdf_obj *dict, const char *path, const char *value);
int safe_pdf_update_object(fz_context *ctx, pdf_document *pdf, int xref, pdf_obj *obj);
int safe_pdf_dict_del_path(fz_context *ctx, pdf_obj *dict, const char *path);
int pdf_obj_num_is_stream(fz_context *ctx, pdf_document *doc, int num);
fz_buffer *safe_pdf_load_stream_number(fz_context *ctx, pdf_document *pdf, int xref);
fz_buffer *safe_pdf_load_raw_stream_number(fz_context *ctx, pdf_document *pdf, int xref);
int pdf_dict_len(fz_context *ctx, pdf_obj *obj);
pdf_obj *pdf_dict_get_key(fz_context *ctx, pdf_obj *dict, int i);
*/
import "C"

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func (f *Document) xrefObjLocked(xref int) (*C.pdf_obj, bool, error) {
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return nil, false, ErrOpenDocument
	}
	if xref == -1 {
		return C.pdf_trailer(f.ctx, pdf), false, nil
	}
	xrefLen := int(C.pdf_xref_len(f.ctx, pdf))
	if xref < 1 || xref >= xrefLen {
		return nil, false, fmt.Errorf("xref out of range: %d", xref)
	}
	obj := C.safe_pdf_load_object(f.ctx, pdf, C.int(xref))
	if obj == nil {
		return nil, false, ErrOpenDocument
	}
	return obj, true, nil
}

func (f *Document) pdfObjStringLocked(obj *C.pdf_obj) (string, error) {
	if obj == nil {
		return "null", nil
	}

	buf := C.fz_new_buffer(f.ctx, 1024)
	defer C.fz_drop_buffer(f.ctx, buf)

	out := C.fz_new_output_with_buffer(f.ctx, buf)
	defer C.fz_drop_output(f.ctx, out)

	if C.safe_pdf_print_obj(f.ctx, out, obj) == 0 {
		return "", ErrOpenDocument
	}
	C.fz_close_output(f.ctx, out)

	n := C.fz_buffer_storage(f.ctx, buf, nil)
	return C.GoStringN(C.fz_string_from_buffer(f.ctx, buf), C.int(n)), nil
}

func hasInvalidPDFNameChars(s string) bool {
	for _, r := range s {
		if r == 0 || unicode.IsSpace(r) {
			return true
		}
		if strings.ContainsRune("()<>[]{}%/", r) {
			return true
		}
	}
	return false
}

func validateXrefSetKeyAndValue(key, value string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("bad 'key'")
	}
	keyParts := strings.Split(key, "/")
	for _, part := range keyParts {
		if part == "" || hasInvalidPDFNameChars(part) {
			return fmt.Errorf("bad 'key'")
		}
	}

	if value == "" {
		return fmt.Errorf("bad 'value'")
	}
	if strings.HasPrefix(value, "/") {
		name := value[1:]
		if name == "" || hasInvalidPDFNameChars(name) {
			return fmt.Errorf("bad 'value'")
		}
	}

	return nil
}

func (f *Document) xrefGetKeyTypedLocked(xref int, key string) (string, string, error) {
	obj, drop, err := f.xrefObjLocked(xref)
	if err != nil {
		return "", "", err
	}
	if drop {
		defer C.pdf_drop_obj(f.ctx, obj)
	}
	if C.pdf_is_indirect(f.ctx, obj) != 0 {
		resolved := C.pdf_resolve_indirect(f.ctx, obj)
		if resolved != nil {
			obj = resolved
		}
	}
	if C.pdf_is_dict(f.ctx, obj) == 0 {
		return "null", "null", nil
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	subObj := C.pdf_dict_getp(f.ctx, obj, ckey)
	if subObj == nil || C.pdf_is_null(f.ctx, subObj) != 0 {
		return "null", "null", nil
	}

	if C.pdf_is_indirect(f.ctx, subObj) != 0 {
		return "xref", fmt.Sprintf("%d 0 R", int(C.pdf_to_num(f.ctx, subObj))), nil
	}
	if C.pdf_is_array(f.ctx, subObj) != 0 {
		text, err := f.pdfObjStringLocked(subObj)
		return "array", text, err
	}
	if C.pdf_is_dict(f.ctx, subObj) != 0 {
		text, err := f.pdfObjStringLocked(subObj)
		return "dict", text, err
	}
	if C.pdf_is_int(f.ctx, subObj) != 0 {
		return "int", strconv.Itoa(int(C.pdf_to_int(f.ctx, subObj))), nil
	}
	if C.pdf_is_real(f.ctx, subObj) != 0 {
		text, err := f.pdfObjStringLocked(subObj)
		return "float", text, err
	}
	if C.pdf_is_bool(f.ctx, subObj) != 0 {
		if C.pdf_to_bool(f.ctx, subObj) != 0 {
			return "bool", "true", nil
		}
		return "bool", "false", nil
	}
	if C.pdf_is_name(f.ctx, subObj) != 0 {
		return "name", "/" + C.GoString(C.pdf_to_name(f.ctx, subObj)), nil
	}
	if C.pdf_is_string(f.ctx, subObj) != 0 {
		return "string", C.GoString(C.pdf_to_text_string(f.ctx, subObj)), nil
	}

	text, err := f.pdfObjStringLocked(subObj)
	return "unknown", text, err
}

// XrefLength returns the xref table length for PDFs.
func (f *Document) XrefLength() int {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return 0
	}
	return int(C.pdf_xref_len(f.ctx, pdf))
}

// XrefObject returns a printable PDF object representation for xref.
// Use xref=-1 for trailer.
func (f *Document) XrefObject(xref int) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	obj, drop, err := f.xrefObjLocked(xref)
	if err != nil {
		return "", err
	}
	if drop {
		defer C.pdf_drop_obj(f.ctx, obj)
	}
	if C.pdf_is_indirect(f.ctx, obj) != 0 {
		resolved := C.pdf_resolve_indirect(f.ctx, obj)
		if resolved != nil {
			obj = resolved
		}
	}
	return f.pdfObjStringLocked(obj)
}

// XrefGetKey returns the value string for a key in a trailer/object dictionary.
func (f *Document) XrefGetKey(xref int, key string) (string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	_, value, err := f.xrefGetKeyTypedLocked(xref, key)
	return value, err
}

// XrefGetKeyTyped returns key value as (type, value), mirroring PyMuPDF.
func (f *Document) XrefGetKeyTyped(xref int, key string) (string, string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.xrefGetKeyTypedLocked(xref, key)
}

// XrefSetKey sets a key to a PDF object value expression in trailer/object xref.
func (f *Document) XrefSetKey(xref int, key, value string) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if err := validateXrefSetKeyAndValue(key, value); err != nil {
		return err
	}

	pdf := f.pdfDocLocked()
	if pdf == nil {
		return ErrOpenDocument
	}

	obj, drop, err := f.xrefObjLocked(xref)
	if err != nil {
		return err
	}
	if drop {
		defer C.pdf_drop_obj(f.ctx, obj)
	}
	if C.pdf_is_indirect(f.ctx, obj) != 0 {
		resolved := C.pdf_resolve_indirect(f.ctx, obj)
		if resolved != nil {
			obj = resolved
		}
	}
	if C.pdf_is_dict(f.ctx, obj) == 0 {
		return ErrOpenDocument
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	if strings.EqualFold(strings.TrimSpace(value), "null") {
		if C.safe_pdf_dict_del_path(f.ctx, obj, ckey) == 0 {
			return fmt.Errorf("bad 'value'")
		}
		if xref > 0 && C.safe_pdf_update_object(f.ctx, pdf, C.int(xref), obj) == 0 {
			return ErrOpenDocument
		}
		return nil
	}

	if C.safe_pdf_dict_putp_parsed(f.ctx, pdf, obj, ckey, cvalue) == 0 {
		return fmt.Errorf("bad 'value'")
	}
	if xref > 0 && C.safe_pdf_update_object(f.ctx, pdf, C.int(xref), obj) == 0 {
		return ErrOpenDocument
	}

	return nil
}

// XrefIsFont reports whether xref is a font dictionary object.
func (f *Document) XrefIsFont(xref int) (bool, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	typ, value, err := f.xrefGetKeyTypedLocked(xref, "Type")
	if err != nil {
		return false, err
	}
	return typ == "name" && value == "/Font", nil
}

// XrefIsImage reports whether xref is an image object.
func (f *Document) XrefIsImage(xref int) (bool, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	typ, value, err := f.xrefGetKeyTypedLocked(xref, "Subtype")
	if err != nil {
		return false, err
	}
	return typ == "name" && value == "/Image", nil
}

// XrefIsForm reports whether xref is a Form XObject.
func (f *Document) XrefIsForm(xref int) (bool, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	typ, value, err := f.xrefGetKeyTypedLocked(xref, "Subtype")
	if err != nil {
		return false, err
	}
	return typ == "name" && value == "/Form", nil
}

// XrefIsStream reports whether xref contains a stream.
func (f *Document) XrefIsStream(xref int) (bool, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	pdf := f.pdfDocLocked()
	if pdf == nil {
		return false, ErrOpenDocument
	}
	xrefLen := int(C.pdf_xref_len(f.ctx, pdf))
	if xref < 1 || xref >= xrefLen {
		return false, fmt.Errorf("xref out of range: %d", xref)
	}
	return int(C.pdf_obj_num_is_stream(f.ctx, pdf, C.int(xref))) != 0, nil
}

func (f *Document) xrefStreamLocked(xref int, raw bool) ([]byte, error) {
	pdf := f.pdfDocLocked()
	if pdf == nil {
		return nil, ErrOpenDocument
	}
	xrefLen := int(C.pdf_xref_len(f.ctx, pdf))
	if xref < 1 || xref >= xrefLen {
		return nil, fmt.Errorf("xref out of range: %d", xref)
	}
	if int(C.pdf_obj_num_is_stream(f.ctx, pdf, C.int(xref))) == 0 {
		return nil, nil
	}

	var buf *C.fz_buffer
	if raw {
		buf = C.safe_pdf_load_raw_stream_number(f.ctx, pdf, C.int(xref))
	} else {
		buf = C.safe_pdf_load_stream_number(f.ctx, pdf, C.int(xref))
	}
	if buf == nil {
		return nil, ErrOpenDocument
	}
	defer C.fz_drop_buffer(f.ctx, buf)

	n := C.fz_buffer_storage(f.ctx, buf, nil)
	if n == 0 {
		return []byte{}, nil
	}
	data := C.GoBytes(unsafe.Pointer(C.fz_string_from_buffer(f.ctx, buf)), C.int(n))
	return data, nil
}

// XrefStream returns decompressed xref stream bytes (nil if xref is not a stream).
func (f *Document) XrefStream(xref int) ([]byte, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.xrefStreamLocked(xref, false)
}

// XrefStreamRaw returns raw xref stream bytes (nil if xref is not a stream).
func (f *Document) XrefStreamRaw(xref int) ([]byte, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.xrefStreamLocked(xref, true)
}

// XrefGetKeys returns dictionary keys for trailer/object xref.
func (f *Document) XrefGetKeys(xref int) ([]string, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	obj, drop, err := f.xrefObjLocked(xref)
	if err != nil {
		return nil, err
	}
	if drop {
		defer C.pdf_drop_obj(f.ctx, obj)
	}
	if C.pdf_is_indirect(f.ctx, obj) != 0 {
		resolved := C.pdf_resolve_indirect(f.ctx, obj)
		if resolved != nil {
			obj = resolved
		}
	}
	if C.pdf_is_dict(f.ctx, obj) == 0 {
		return []string{}, nil
	}

	n := int(C.pdf_dict_len(f.ctx, obj))
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		keyObj := C.pdf_dict_get_key(f.ctx, obj, C.int(i))
		if keyObj == nil {
			continue
		}
		if C.pdf_is_name(f.ctx, keyObj) != 0 {
			out = append(out, C.GoString(C.pdf_to_name(f.ctx, keyObj)))
		}
	}
	return out, nil
}
