package gopherpdf

import (
	"strconv"
	"strings"
	"testing"
)

func TestPyMuPDFParity_XrefAccess(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	xrefLen, err := doc.XrefLength()
	if err != nil {
		t.Fatalf("xref length: %v", err)
	}
	if xrefLen < 2 {
		t.Fatalf("unexpected xref length: %d", xrefLen)
	}

	trailer, err := doc.XrefObject(-1)
	if err != nil {
		t.Fatalf("xref trailer object: %v", err)
	}
	if !strings.Contains(trailer, "/Root") {
		t.Fatalf("expected trailer object to contain /Root")
	}

	keys, err := doc.XrefGetKeys(-1)
	if err != nil {
		t.Fatalf("xref get keys trailer: %v", err)
	}
	foundRoot := false
	for _, k := range keys {
		if k == "Root" {
			foundRoot = true
			break
		}
	}
	if !foundRoot {
		t.Fatalf("expected trailer keys to include Root")
	}

	root, err := doc.XrefGetKey(-1, "Root")
	if err != nil {
		t.Fatalf("xref get key trailer Root: %v", err)
	}
	if !strings.Contains(root, "R") {
		t.Fatalf("expected trailer Root to be indirect reference, got %q", root)
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()
	if _, err := epub.XrefLength(); err == nil {
		t.Fatalf("expected xref length on non-pdf to return error")
	}
	if _, err := epub.XrefObject(-1); err == nil {
		t.Fatalf("expected xref object on non-pdf to return error")
	}
	if _, err := epub.XrefGetKey(-1, "Info"); err == nil {
		t.Fatalf("expected xref get key on non-pdf to return error")
	}
	if _, err := epub.XrefGetKeys(-1); err == nil {
		t.Fatalf("expected xref get keys on non-pdf to return error")
	}
}

// Mirrors typed get/set key behavior in tests/test_object_manipulation.py.
func TestPyMuPDFParity_XrefSetAndTypedGet(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	xrefLen, err := doc.XrefLength()
	if err != nil {
		t.Fatalf("xref length: %v", err)
	}
	keyType, value, err := doc.XrefGetKeyTyped(-1, "Size")
	if err != nil {
		t.Fatalf("xref get key typed trailer Size: %v", err)
	}
	if keyType != "int" {
		t.Fatalf("expected key type int, got %q", keyType)
	}
	if value != strconv.Itoa(xrefLen) {
		t.Fatalf("expected trailer size %d, got %q", xrefLen, value)
	}

	pageXref, err := doc.PageXref(0)
	if err != nil {
		t.Fatalf("page xref: %v", err)
	}

	if err := doc.XrefSetKey(pageXref, "Rotate", "90"); err != nil {
		t.Fatalf("xref set key Rotate: %v", err)
	}
	keyType, value, err = doc.XrefGetKeyTyped(pageXref, "Rotate")
	if err != nil {
		t.Fatalf("xref get key typed Rotate: %v", err)
	}
	if keyType != "int" || value != "90" {
		t.Fatalf("unexpected Rotate key: type=%q value=%q", keyType, value)
	}

	if err := doc.XrefSetKey(pageXref, "my_rotate/something", "90"); err != nil {
		t.Fatalf("xref set key nested int: %v", err)
	}
	keyType, value, err = doc.XrefGetKeyTyped(pageXref, "my_rotate/something")
	if err != nil {
		t.Fatalf("xref get key typed nested int: %v", err)
	}
	if keyType != "int" || value != "90" {
		t.Fatalf("unexpected nested int key: type=%q value=%q", keyType, value)
	}

	if err := doc.XrefSetKey(pageXref, "my_rotate", "/90"); err != nil {
		t.Fatalf("xref set key name: %v", err)
	}
	keyType, value, err = doc.XrefGetKeyTyped(pageXref, "my_rotate")
	if err != nil {
		t.Fatalf("xref get key typed name: %v", err)
	}
	if keyType != "name" || value != "/90" {
		t.Fatalf("unexpected name key: type=%q value=%q", keyType, value)
	}

	if err := doc.XrefSetKey(pageXref, "MediaBox", "[-30 -20 595 842]"); err != nil {
		t.Fatalf("xref set key MediaBox: %v", err)
	}
	keyType, value, err = doc.XrefGetKeyTyped(pageXref, "MediaBox")
	if err != nil {
		t.Fatalf("xref get key typed MediaBox: %v", err)
	}
	if keyType != "array" {
		t.Fatalf("expected MediaBox type array, got %q", keyType)
	}
	if value != "[-30 -20 595 842]" {
		t.Fatalf("unexpected MediaBox value: %q", value)
	}

	if err := doc.XrefSetKey(pageXref, "my rotate", "90"); err == nil || err.Error() != "bad 'key'" {
		t.Fatalf("expected bad 'key' validation error, got %v", err)
	}
	if err := doc.XrefSetKey(pageXref, "my_rotate", "/9/0"); err == nil || err.Error() != "bad 'value'" {
		t.Fatalf("expected bad 'value' validation error, got %v", err)
	}
	if err := doc.XrefSetKey(-1, "Info", "null"); err != nil {
		t.Fatalf("xref set key trailer Info null: %v", err)
	}
	keyType, value, err = doc.XrefGetKeyTyped(-1, "Info")
	if err != nil {
		t.Fatalf("xref get key typed trailer Info after null: %v", err)
	}
	if keyType != "null" || value != "null" {
		t.Fatalf("unexpected trailer Info after null set: type=%q value=%q", keyType, value)
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()
	if _, _, err := epub.XrefGetKeyTyped(-1, "Info"); err == nil {
		t.Fatalf("expected xref get key typed on non-pdf to return error")
	}
	if err := epub.XrefSetKey(-1, "Info", "null"); err == nil {
		t.Fatalf("expected xref set key on non-pdf to return error")
	}
}

// Mirrors low-level box manipulation intent from tests/test_general.py::test_2736.
func TestPyMuPDFParity_XrefSetKeyAffectsPageBoxes(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	pageXref, err := doc.PageXref(0)
	if err != nil {
		t.Fatalf("page xref: %v", err)
	}

	if err := doc.XrefSetKey(pageXref, "MediaBox", "[-30 -20 595 842]"); err != nil {
		t.Fatalf("xref set key MediaBox: %v", err)
	}
	if err := doc.XrefSetKey(pageXref, "CropBox", "[-20 -10 595 842]"); err != nil {
		t.Fatalf("xref set key CropBox: %v", err)
	}

	keyType, value, err := doc.XrefGetKeyTyped(pageXref, "MediaBox")
	if err != nil {
		t.Fatalf("xref get key typed MediaBox: %v", err)
	}
	if keyType != "array" || value != "[-30 -20 595 842]" {
		t.Fatalf("unexpected MediaBox value: type=%q value=%q", keyType, value)
	}

	keyType, value, err = doc.XrefGetKeyTyped(pageXref, "CropBox")
	if err != nil {
		t.Fatalf("xref get key typed CropBox: %v", err)
	}
	if keyType != "array" || value != "[-20 -10 595 842]" {
		t.Fatalf("unexpected CropBox value: type=%q value=%q", keyType, value)
	}

	pageBound, err := doc.Bound(0)
	if err != nil {
		t.Fatalf("bound page 0: %v", err)
	}
	if pageBound.Min.X != 0 || pageBound.Min.Y != 0 || pageBound.Dx() != 615 || pageBound.Dy() != 852 {
		t.Fatalf("unexpected page bound after xref box edits: %v", pageBound)
	}

	media, err := doc.BoundBox(0, PageBoxMedia)
	if err != nil {
		t.Fatalf("media box page 0: %v", err)
	}
	if media.Min.X != -10 || media.Min.Y != 0 || media.Dx() != 625 || media.Dy() != 862 {
		t.Fatalf("unexpected media box after xref box edits: %v", media)
	}
}

// Mirrors xref_is_* helper behavior from PyMuPDF.
func TestPyMuPDFParity_XrefPredicates(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	xrefLen, err := doc.XrefLength()
	if err != nil {
		t.Fatalf("xref length: %v", err)
	}

	findBy := func(pred func(int) (bool, error)) int {
		t.Helper()
		for i := 1; i < xrefLen; i++ {
			ok, err := pred(i)
			if err != nil {
				t.Fatalf("predicate xref %d: %v", i, err)
			}
			if ok {
				return i
			}
		}
		return 0
	}

	fontXref := findBy(doc.XrefIsFont)
	if fontXref <= 0 {
		t.Fatalf("expected at least one font xref")
	}
	imageXref := findBy(doc.XrefIsImage)
	if imageXref <= 0 {
		t.Fatalf("expected at least one image xref")
	}

	isStream, err := doc.XrefIsStream(imageXref)
	if err != nil {
		t.Fatalf("xref is stream image: %v", err)
	}
	if !isStream {
		t.Fatalf("expected image xref %d to be a stream", imageXref)
	}
	streamRaw, err := doc.XrefStreamRaw(imageXref)
	if err != nil {
		t.Fatalf("xref stream raw image: %v", err)
	}
	if len(streamRaw) == 0 {
		t.Fatalf("expected non-empty raw stream for image xref %d", imageXref)
	}
	streamDecoded, err := doc.XrefStream(imageXref)
	if err != nil {
		t.Fatalf("xref stream decoded image: %v", err)
	}
	if len(streamDecoded) == 0 {
		t.Fatalf("expected non-empty decoded stream for image xref %d", imageXref)
	}

	catalogXref, err := doc.CatalogXref()
	if err != nil {
		t.Fatalf("catalog xref: %v", err)
	}
	isStream, err = doc.XrefIsStream(catalogXref)
	if err != nil {
		t.Fatalf("xref is stream catalog: %v", err)
	}
	if isStream {
		t.Fatalf("expected catalog xref %d to not be a stream", catalogXref)
	}
	streamRaw, err = doc.XrefStreamRaw(catalogXref)
	if err != nil {
		t.Fatalf("xref stream raw catalog: %v", err)
	}
	if streamRaw != nil {
		t.Fatalf("expected nil raw stream for non-stream catalog object")
	}
	streamDecoded, err = doc.XrefStream(catalogXref)
	if err != nil {
		t.Fatalf("xref stream decoded catalog: %v", err)
	}
	if streamDecoded != nil {
		t.Fatalf("expected nil decoded stream for non-stream catalog object")
	}

	formDoc, err := Open(resourcePath(t, "1.pdf"))
	if err != nil {
		t.Fatalf("open 1.pdf: %v", err)
	}
	defer formDoc.Close()
	formXrefLen, err := formDoc.XrefLength()
	if err != nil {
		t.Fatalf("xref length 1.pdf: %v", err)
	}
	foundForm := false
	for i := 1; i < formXrefLen; i++ {
		ok, err := formDoc.XrefIsForm(i)
		if err != nil {
			t.Fatalf("xref is form xref %d: %v", i, err)
		}
		if ok {
			foundForm = true
			break
		}
	}
	if !foundForm {
		t.Fatalf("expected at least one Form xref in 1.pdf")
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()
	if _, err := epub.XrefIsFont(1); err == nil {
		t.Fatalf("expected xref is font on non-pdf to return error")
	}
	if _, err := epub.XrefIsImage(1); err == nil {
		t.Fatalf("expected xref is image on non-pdf to return error")
	}
	if _, err := epub.XrefIsForm(1); err == nil {
		t.Fatalf("expected xref is form on non-pdf to return error")
	}
	if _, err := epub.XrefIsStream(1); err == nil {
		t.Fatalf("expected xref is stream on non-pdf to return error")
	}
	if _, err := epub.XrefStream(1); err == nil {
		t.Fatalf("expected xref stream on non-pdf to return error")
	}
	if _, err := epub.XrefStreamRaw(1); err == nil {
		t.Fatalf("expected xref stream raw on non-pdf to return error")
	}
}
