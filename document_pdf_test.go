package gopherpdf

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPyMuPDFParity_IsRepairedForRecoveredPDF(t *testing.T) {
	// Use a fixture that MuPDF reports as "repaired" on open.
	doc, err := Open(resourcePath(t, "test2238.pdf"))
	if err != nil {
		t.Fatalf("open repaired: %v", err)
	}
	defer doc.Close()

	repaired, err := doc.IsRepaired()
	if err != nil {
		t.Fatalf("is repaired: %v", err)
	}
	if !repaired {
		t.Fatalf("expected fixture to report repaired=true")
	}
}

func TestPyMuPDFParity_HasLinksAndHasAnnotations(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	hasLinks, err := doc.HasLinks()
	if err != nil {
		t.Fatalf("has links: %v", err)
	}
	if !hasLinks {
		t.Fatalf("expected 2.pdf to report has_links=true")
	}

	hasAnnots, err := doc.HasAnnotations()
	if err != nil {
		t.Fatalf("has annotations: %v", err)
	}
	if hasAnnots {
		t.Fatalf("expected 2.pdf to report has_annotations=false for this fixture")
	}

	annotDoc, err := Open(resourcePath(t, "test_annot_file_info.pdf"))
	if err != nil {
		t.Fatalf("open annotated fixture: %v", err)
	}
	defer annotDoc.Close()

	widgetHasAnnots, err := annotDoc.HasAnnotations()
	if err != nil {
		t.Fatalf("annotated fixture has annotations: %v", err)
	}
	if !widgetHasAnnots {
		t.Fatalf("expected annotated fixture to have annotations")
	}
}

func TestPyMuPDFParity_ResolveNames(t *testing.T) {
	doc, err := Open(resourcePath(t, "bug1945.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	got, err := doc.ResolveNames()
	if err != nil {
		t.Fatalf("resolve names: %v", err)
	}
	if len(got) == 0 {
		t.Fatalf("expected non-empty names map")
	}
	if _, ok := got["orb-modules"]; !ok {
		t.Fatalf("expected orb-modules to exist")
	}
}

func TestImageAndImageDPI(t *testing.T) {
	pdfBytes := mustReadResource(t, "test_2548.pdf")
	doc, err := OpenBytes(pdfBytes)
	if err != nil {
		t.Fatalf("open bytes: %v", err)
	}
	defer doc.Close()

	imgDefault, err := doc.Image(0)
	if err != nil {
		t.Fatalf("image page 0: %v", err)
	}
	if imgDefault.Bounds().Dx() <= 0 || imgDefault.Bounds().Dy() <= 0 {
		t.Fatalf("invalid default image bounds: %v", imgDefault.Bounds())
	}

	img120, err := doc.ImageDPI(0, 120)
	if err != nil {
		t.Fatalf("image dpi page 0: %v", err)
	}
	if img120.Bounds().Dx() <= 0 || img120.Bounds().Dy() <= 0 {
		t.Fatalf("invalid dpi image bounds: %v", img120.Bounds())
	}
	if img120.Bounds().Dx() >= imgDefault.Bounds().Dx() {
		t.Fatalf("expected 120 DPI width (%d) to be less than default width (%d)", img120.Bounds().Dx(), imgDefault.Bounds().Dx())
	}
}

func TestPyMuPDFParity_CatalogPageXrefAndFormState(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	catalogXref, err := doc.CatalogXref()
	if err != nil {
		t.Fatalf("catalog xref: %v", err)
	}
	if catalogXref <= 0 {
		t.Fatalf("unexpected catalog xref: %d", catalogXref)
	}
	catalogObj, err := doc.XrefObject(catalogXref)
	if err != nil {
		t.Fatalf("catalog object: %v", err)
	}
	if !strings.Contains(catalogObj, "/Type/Catalog") {
		t.Fatalf("expected catalog object, got: %.120s", catalogObj)
	}

	pageXref, err := doc.PageXref(0)
	if err != nil {
		t.Fatalf("page xref: %v", err)
	}
	if pageXref <= 0 {
		t.Fatalf("unexpected page xref: %d", pageXref)
	}
	pageObj, err := doc.XrefObject(pageXref)
	if err != nil {
		t.Fatalf("page object: %v", err)
	}
	if !strings.Contains(pageObj, "/Type/Page") {
		t.Fatalf("expected page object, got: %.120s", pageObj)
	}

	formCount, err := doc.FormFieldCount()
	if err != nil {
		t.Fatalf("form field count: %v", err)
	}
	if formCount != 0 {
		t.Fatalf("expected 2.pdf to have no form fields, got %d", formCount)
	}
	isForm, err := doc.IsFormPDF()
	if err != nil {
		t.Fatalf("is form pdf: %v", err)
	}
	if isForm {
		t.Fatalf("expected 2.pdf is_form_pdf=false")
	}

	widget, err := Open(resourcePath(t, "widgettest.pdf"))
	if err != nil {
		t.Fatalf("open widgettest: %v", err)
	}
	defer widget.Close()
	formCount, err = widget.FormFieldCount()
	if err != nil {
		t.Fatalf("form field count widgettest: %v", err)
	}
	if formCount <= 0 {
		t.Fatalf("expected widgettest.pdf to have form fields, got %d", formCount)
	}
	isForm, err = widget.IsFormPDF()
	if err != nil {
		t.Fatalf("is form pdf widgettest: %v", err)
	}
	if !isForm {
		t.Fatalf("expected widgettest.pdf is_form_pdf=true")
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()
	if _, err := epub.CatalogXref(); err == nil {
		t.Fatalf("expected catalog xref on non-pdf to return error")
	}
	if _, err := epub.PageXref(0); err == nil {
		t.Fatalf("expected page xref on non-pdf to return error")
	}
	if _, err := epub.FormFieldCount(); err == nil {
		t.Fatalf("expected form field count on non-pdf to return error")
	}
	if _, err := epub.IsFormPDF(); err == nil {
		t.Fatalf("expected is_form_pdf on non-pdf to return error")
	}
}

// Mirrors PyMuPDF tests/test_page_links.py.
func TestPyMuPDFParity_PageLinksGenerator(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	links, err := doc.Links(pages - 1)
	if err != nil {
		t.Fatalf("links on last page: %v", err)
	}
	if len(links) != 7 {
		t.Fatalf("expected 7 links, got %d", len(links))
	}
}

// Mirrors permission and password state checks from tests/test_crypting.py.
func TestPyMuPDFParity_DocumentPermissionsAndPasswordState(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	needsPass, err := doc.NeedsPassword()
	if err != nil {
		t.Fatalf("needs password: %v", err)
	}
	if needsPass {
		t.Fatalf("expected unencrypted fixture to not need a password")
	}

	isEncrypted, err := doc.IsEncrypted()
	if err != nil {
		t.Fatalf("is encrypted: %v", err)
	}
	if isEncrypted {
		t.Fatalf("expected unencrypted fixture to report isEncrypted=false")
	}

	for _, perm := range []Permission{
		PermissionPrint,
		PermissionCopy,
		PermissionEdit,
		PermissionAnnotate,
		PermissionForm,
		PermissionAccessibility,
		PermissionAssemble,
		PermissionPrintHQ,
	} {
		ok, err := doc.HasPermission(perm)
		if err != nil {
			t.Fatalf("has permission %d: %v", perm, err)
		}
		if !ok {
			t.Fatalf("expected permission %d to be granted", perm)
		}
	}
}

// Mirrors page label lookup behavior in tests/test_pagelabels.py.
func TestPyMuPDFParity_PageLabelsAndLookup(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	wantFirst := []string{"1", "2", "i", "ii", "1", "2"}
	for i, want := range wantFirst {
		got, err := doc.PageLabel(i)
		if err != nil {
			t.Fatalf("page label %d: %v", i, err)
		}
		if got != want {
			t.Fatalf("page label mismatch at %d: got %q want %q", i, got, want)
		}
	}

	labels, err := doc.GetPageLabels()
	if err != nil {
		t.Fatalf("get page labels: %v", err)
	}
	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	if len(labels) != pages {
		t.Fatalf("labels length mismatch: got %d want %d", len(labels), pages)
	}
	for i, want := range wantFirst {
		if labels[i] != want {
			t.Fatalf("get page labels mismatch at %d: got %q want %q", i, labels[i], want)
		}
	}

	cases := map[string][]int{
		"1":  {0, 4},
		"2":  {1, 5},
		"i":  {2},
		"ii": {3},
		"V":  {},
	}
	for label, want := range cases {
		got, err := doc.GetPageNumbers(label)
		if err != nil {
			t.Fatalf("page numbers by label %q: %v", label, err)
		}
		if len(got) != len(want) {
			t.Fatalf("label %q length mismatch: got %v want %v", label, got, want)
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("label %q mismatch: got %v want %v", label, got, want)
			}
		}
	}
}

// Mirrors tests/test_pagelabels.py::test_setlabels and label rule roundtrip.
func TestPyMuPDFParity_PageLabelRulesAndSet(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	rules := []PageLabelRule{
		{StartPage: 0, Prefix: "A-", Style: "D", FirstPageNum: 1},
		{StartPage: 4, Prefix: "", Style: "R", FirstPageNum: 1},
	}
	if err := doc.SetPageLabels(rules); err != nil {
		t.Fatalf("set page labels: %v", err)
	}

	wantLabels := []string{"A-1", "A-2", "A-3", "A-4", "I", "II", "III", "IV", "V", "VI"}
	for i, want := range wantLabels {
		got, err := doc.PageLabel(i)
		if err != nil {
			t.Fatalf("page label %d: %v", i, err)
		}
		if got != want {
			t.Fatalf("page label mismatch at %d: got %q want %q", i, got, want)
		}
	}

	pnos, err := doc.GetPageNumbers("V")
	if err != nil {
		t.Fatalf("get page numbers V: %v", err)
	}
	if len(pnos) != 1 || pnos[0] != 8 {
		t.Fatalf("expected page number [8] for label V, got %v", pnos)
	}

	gotRules, err := doc.PageLabelRules()
	if err != nil {
		t.Fatalf("get page label rules: %v", err)
	}
	if len(gotRules) != len(rules) {
		t.Fatalf("rule count mismatch: got %d want %d", len(gotRules), len(rules))
	}
	for i, want := range rules {
		if gotRules[i] != want {
			t.Fatalf("rule mismatch at %d: got %+v want %+v", i, gotRules[i], want)
		}
	}
}

// Mirrors tests/test_pagelabels.py::test_labels_styleA.
func TestPyMuPDFParity_PageLabelRulesStyleA(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	rules := []PageLabelRule{
		{StartPage: 0, Prefix: "", Style: "a", FirstPageNum: 1},
		{StartPage: 5, Prefix: "", Style: "A", FirstPageNum: 1},
	}
	if err := doc.SetPageLabels(rules); err != nil {
		t.Fatalf("set page labels style A/a: %v", err)
	}

	wantLabels := []string{"a", "b", "c", "d", "e", "A", "B", "C", "D", "E"}
	for i, want := range wantLabels {
		got, err := doc.PageLabel(i)
		if err != nil {
			t.Fatalf("page label %d: %v", i, err)
		}
		if got != want {
			t.Fatalf("page label mismatch at %d: got %q want %q", i, got, want)
		}
	}

	gotRules, err := doc.GetPageLabelRules()
	if err != nil {
		t.Fatalf("get page label rules alias: %v", err)
	}
	if len(gotRules) != len(rules) {
		t.Fatalf("rule count mismatch: got %d want %d", len(gotRules), len(rules))
	}
	for i, want := range rules {
		if gotRules[i] != want {
			t.Fatalf("rule mismatch at %d: got %+v want %+v", i, gotRules[i], want)
		}
	}

	epub, err := Open(resourcePath(t, "Bezier.epub"))
	if err != nil {
		t.Fatalf("open epub: %v", err)
	}
	defer epub.Close()
	if _, err := epub.PageLabelRules(); err == nil {
		t.Fatalf("expected get page label rules on non-pdf to return error")
	}
	if err := epub.SetPageLabels(rules); err == nil {
		t.Fatalf("expected set page labels on non-pdf to return error")
	}
}

// Mirrors internal link destination resolution behavior.
func TestPyMuPDFParity_ResolveLink(t *testing.T) {
	doc, err := Open(resourcePath(t, "2.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	pages, err := doc.NumPages()
	if err != nil {
		t.Fatalf("num pages: %v", err)
	}
	links, err := doc.Links(pages - 1)
	if err != nil {
		t.Fatalf("links: %v", err)
	}
	if len(links) < 2 {
		t.Fatalf("expected links for resolve tests")
	}

	internal, ok, err := doc.ResolveLink(links[0].URI)
	if err != nil {
		t.Fatalf("resolve internal: %v", err)
	}
	if !ok {
		t.Fatalf("expected internal link to resolve: %q", links[0].URI)
	}
	if internal.Location != (Location{Chapter: 0, Page: 31}) {
		t.Fatalf("unexpected internal location: %+v", internal.Location)
	}
	if internal.X <= 0 || internal.Y <= 0 {
		t.Fatalf("unexpected internal coordinates: %+v", internal)
	}

	external, ok, err := doc.ResolveLink(links[1].URI)
	if err != nil {
		t.Fatalf("resolve external: %v", err)
	}
	if ok {
		t.Fatalf("expected external link to be unresolved")
	}
	if external.Location != (Location{Chapter: -1, Page: -1}) {
		t.Fatalf("unexpected unresolved location: %+v", external.Location)
	}
}

// Mirrors crop/media box distinctions from tests/test_general.py::test_2710.
func TestPyMuPDFParity_PageBoxes(t *testing.T) {
	doc, err := Open(resourcePath(t, "test_2710.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	bound, err := doc.Bound(0)
	if err != nil {
		t.Fatalf("bound: %v", err)
	}
	media, err := doc.BoundBox(0, PageBoxMedia)
	if err != nil {
		t.Fatalf("media box: %v", err)
	}
	crop, err := doc.BoundBox(0, PageBoxCrop)
	if err != nil {
		t.Fatalf("crop box: %v", err)
	}
	if bound != crop {
		t.Fatalf("expected page bound to match crop box: bound=%v crop=%v", bound, crop)
	}
	if media.Dx() <= crop.Dx() || media.Dy() <= crop.Dy() {
		t.Fatalf("expected media box larger than crop box: media=%v crop=%v", media, crop)
	}

	for _, box := range []PageBox{PageBoxBleed, PageBoxTrim, PageBoxArt} {
		got, err := doc.BoundBox(0, box)
		if err != nil {
			t.Fatalf("bound box %d: %v", box, err)
		}
		if got != crop {
			t.Fatalf("expected box %d to match crop box for fixture: got=%v crop=%v", box, got, crop)
		}
	}
}

// Mirrors PyMuPDF tests/test_general.py::test_2548 for non-empty extracted outlines in JSON.
func TestPyMuPDFParity_DocumentJSONSmoke(t *testing.T) {
	doc, err := Open(resourcePath(t, "test_2548.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	jsonText, err := doc.JSON(0)
	if err != nil {
		t.Fatalf("json page 0: %v", err)
	}
	if strings.TrimSpace(jsonText) == "" {
		t.Fatalf("expected non-empty json output")
	}
	var tmp any
	if err := json.Unmarshal([]byte(jsonText), &tmp); err != nil {
		t.Fatalf("expected valid json, got error: %v", err)
	}
}

func TestPyMuPDFParity_ResolveNamesJSONFixtureLoad(t *testing.T) {
	// Ensures our resolve names API doesn't regress in (de)serialization expectations for callers.
	doc, err := Open(resourcePath(t, "bug1945.pdf"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer doc.Close()

	names, err := doc.ResolveNames()
	if err != nil {
		t.Fatalf("resolve names: %v", err)
	}
	b, err := json.Marshal(names)
	if err != nil {
		t.Fatalf("marshal names: %v", err)
	}
	if !bytes.Contains(b, []byte("orb-modules")) {
		t.Fatalf("expected orb-modules in marshaled names")
	}
}
