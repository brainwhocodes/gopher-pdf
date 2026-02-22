# gopher-pdf Refactor Items

## P0 - Structure and Maintainability

- [ ] Split `internal/fitz/fitz_cgo.go` into domain files
  - target: `fitz_cgo_open.go`, `fitz_cgo_nav.go`, `fitz_cgo_labels.go`, `fitz_cgo_text.go`, `fitz_cgo_render.go`, `fitz_cgo_toc.go`, `fitz_cgo_metadata.go`
  - keep existing `fitz_cgo_xref.go` as-is
  - acceptance: no behavior change, `CGO_ENABLED=1 go test ./pkg/gopher-pdf -count=1` passes

- [ ] Split `document.go` into domain files
  - target: `document_open.go`, `document_nav.go`, `document_pdf.go`, `document_text.go`, `document_render.go`, `document_meta.go`, `document_aliases.go`, `document_guards.go`
  - acceptance: public API signatures unchanged, package tests green

- [ ] Split large test file `document_test.go`
  - target: `document_metadata_test.go`, `document_toc_test.go`, `document_nonpdf_test.go`, `document_text_test.go`, `document_open_test.go`
  - acceptance: no duplicate tests, same or better coverage, tests green

## P1 - API/Behavior Consistency

- [ ] Unify password/open error behavior between `document.go` and `renderer.go`
  - current duplication: `wrapOpenResult(...)` vs `openMemoryWithPassword(...)`
  - acceptance: one shared open/auth path and consistent error wording

- [ ] Promote common runtime errors to typed/sentinel errors
  - candidates: document closed, page out of range, chapter out of range
  - acceptance: callers can check via `errors.Is(...)`, existing tests updated accordingly

- [ ] Centralize cgo string and buffer conversion helpers
  - reduce repeated `C.CString` / `C.free` / `C.GoString` patterns in `internal/fitz/fitz_cgo*.go`
  - acceptance: helper wrappers used in all touched call sites, no leaks/regressions

## P2 - Parity Tracking and Safety Nets

- [ ] Add parity matrix doc mapping PyMuPDF test cases to Go tests
  - target: `backend/pkg/gopher-pdf/PARITY_MATRIX.md`
  - acceptance: each implemented API area has upstream test references and local test names

- [ ] Add benchmark suite for hot paths
  - target: `document_benchmark_test.go`
  - focus: `GetText`, `Search`, `RenderPagePNG`, `XrefGetKeyTyped`
  - acceptance: benchmarks run with `go test -bench . ./pkg/gopher-pdf`

- [ ] Add CI matrix jobs for cgo/non-cgo package validation
  - `CGO_ENABLED=1 go test ./pkg/gopher-pdf -count=1`
  - `CGO_ENABLED=0 go test ./pkg/gopher-pdf -run '^$' -count=1`
  - acceptance: both lanes required for merge

## Already Completed

- [x] Removed go-fitz dependency path and consolidated on local MuPDF cgo binding
- [x] Added PDF sentinel error (`ErrNoPDF`) and normalized PDF-only guards
- [x] Split xref internals into `internal/fitz/fitz_cgo_xref.go`
- [x] Added dedicated xref parity tests in `document_xref_test.go`
- [x] Added additional parity tests for reflow/chapter/page-label/render/wrapper APIs
