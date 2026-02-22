# gopher-pdf Parity Matrix (PyMuPDF -> Go)

This document maps PyMuPDF API/test areas to the local `gopherpdf` package surface and tests.

## Open / Close / File Types

- PyMuPDF: `fitz.open()`, `Document` lifecycle, `open(stream=..., filetype=...)`
- Go API: `Open`, `OpenWithPassword`, `OpenBytes`, `OpenBytesWithFileType`, `OpenReaderWithFileType*`, `(*Document).Close`
- Tests:
  - `TestOpenBytesAndCoreAPIs` (`document_open_test.go`)
  - `TestPyMuPDFParity_Open2_MultiFormatFileOpen` (`document_open_test.go`)
  - `TestPyMuPDFParity_Open2_StreamWithFileTypeHint` (`document_open_test.go`)
  - `TestPyMuPDFParity_OpenExceptions` (`document_open_test.go`)
  - `TestCloseIsIdempotent` (`document_open_test.go`)

## Passwords / Permissions

- PyMuPDF: `needs_pass`, `authenticate`, `permissions`
- Go API: `(*Document).NeedsPassword`, `(*Document).HasPermission`, `(*Document).IsEncrypted`
- Tests:
  - `TestPyMuPDFParity_DocumentPermissionsAndPasswordState` (`document_pdf_test.go`)

## Metadata

- PyMuPDF: `doc.metadata`, `doc.set_metadata`, erase semantics
- Go API: `(*Document).Metadata`, `(*Document).SetMetadata`
- Tests:
  - `TestPyMuPDFParity_Metadata` (`document_metadata_test.go`)
  - `TestPyMuPDFParity_MetadataEraseAndSet` (`document_metadata_test.go`)

## Rendering

- PyMuPDF: `page.get_pixmap()`, `page.get_svg_image()`
- Go API: `(*Document).Image`, `(*Document).ImageDPI`, `(*Document).RenderPagePNG`, `(*Document).SVG`, `(*Document).Bound`, `(*Document).BoundBox`
- Tests:
  - `TestImageAndImageDPI` (`document_pdf_test.go`)
  - `TestPyMuPDFParity_PageBoxes` (`document_pdf_test.go`)
  - `TestPyMuPDFParity_RenderPagePNG` (`document_additional_parity_test.go`)

## Text Extraction / Search

- PyMuPDF: `page.get_text(...)`, `search_for(...)`
- Go API: `(*Document).Text*`, `(*Document).GetText`, `(*Document).Search`, `(*Document).SearchCount`, `(*Document).TextBox`, `(*Document).JSON*`
- Tests:
  - `document_text_test.go`
  - `TestPyMuPDFParity_DocumentJSONSmoke` (`document_pdf_test.go`)

## Links / Destinations / Labels

- PyMuPDF: `page.get_links()`, `resolve_link`, name tree destinations, page labels
- Go API: `(*Document).Links`, `(*Document).HasLinks`, `(*Document).ResolveLink`, `(*Document).ResolveNames`, `(*Document).PageLabel*`, `(*Document).SetPageLabels`
- Tests:
  - `TestPyMuPDFParity_PageLinksGenerator` (`document_pdf_test.go`)
  - `TestPyMuPDFParity_ResolveLink` (`document_pdf_test.go`)
  - `TestPyMuPDFParity_ResolveNames` (`document_pdf_test.go`)
  - `TestPyMuPDFParity_PageLabelsAndLookup` (`document_pdf_test.go`)
  - `TestPyMuPDFParity_PageLabelRulesAndSet` (`document_pdf_test.go`)

## Xref / PDF Introspection

- PyMuPDF: `xref_*` family
- Go API: `(*Document).XrefGetKeyTyped`, `(*Document).XrefObject`, `(*Document).CatalogXref`, `(*Document).PageXref`, `(*Document).FormFieldCount`
- Tests:
  - `document_xref_test.go`
  - `TestPyMuPDFParity_CatalogPageXrefAndFormState` (`document_pdf_test.go`)
  - `document_additional_parity_test.go`

## Non-PDF Documents (EPUB, etc.)

- PyMuPDF: open non-PDF containers, chapter/page locations, layout
- Go API: `(*Document).NumChapters`, `(*Document).ChapterPageCount`, `(*Document).LastLocation`, `(*Document).Layout`, `(*Document).LocationFromPageNumber`
- Tests:
  - `TestPyMuPDFParity_NonPDF_EpubOpens` (`document_nonpdf_test.go`)
  - `TestPyMuPDFParity_NonPDF_PageIDsAndLayout` (`document_nonpdf_test.go`)

