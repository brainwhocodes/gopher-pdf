# gopher-pdf

Unified Go-only PDF package (no Python runtime) that merges:
- rendering/page rasterization
- document/page/text/html/svg access
- metadata, links, and table-of-contents helpers

Primary APIs:
- `Open`, `OpenWithPassword`
- `OpenBytes`, `OpenBytesWithPassword`, `OpenBytesWithFileType`, `OpenBytesWithFileTypeAndPassword`
- `OpenReader`, `OpenReaderWithPassword`, `OpenReaderWithFileType`, `OpenReaderWithFileTypeAndPassword`
- `(*Document).NumPages`, `Bound`, `Text`, `TextWithFlags`, `GetText`
- `(*Document).GetPageText`, `SearchPageFor`, `GetPagePixmap`
- `(*Document).HTML`, `XHTML`, `XML`, `JSON`, `SearchCount`, `SVG`, `Links`, `Metadata`, `ToC`
- `(*Document).NumChapters`, `ChapterPageCount`, `LastLocation`, `NextLocation`, `PrevLocation`
- `(*Document).LocationFromPageNumber`, `PageNumberFromLocation`, `MakeBookmark`, `LookupBookmark`, `Layout`
- `(*Document).Image`, `ImageDPI`
- `(*Document).RenderPagePNG`, `RenderPage`
- `NewMuPDFRenderer`, `(*MuPDFRenderer).CountPages`, `Render`

Backend dependency:
- in-repo MuPDF cgo binding at `backend/pkg/gopher-pdf/internal/fitz`
- `CGO_ENABLED=1` required for runtime document operations (`nocgo` builds compile with explicit runtime errors)

Test fixtures:
- `backend/pkg/gopher-pdf/resources`
