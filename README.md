# gopher-pdf

Go package for PDF rendering and document extraction via MuPDF (no Python runtime).

## Install

```bash
go get github.com/brainwhocodes/gopher-pdf@latest
```

## Requirements

- Go 1.24+
- CGO enabled (`CGO_ENABLED=1`) with a working C toolchain (clang/gcc)
  - macOS: `xcode-select --install`
  - Debian/Ubuntu: `sudo apt-get install build-essential`
  - Alpine: `apk add build-base`
  - Windows: a `gcc` toolchain in `PATH` (MinGW-w64)

This module vendors MuPDF headers and prebuilt libraries under:
- `internal/fitz/include`
- `internal/fitz/libs`

Prebuilt libraries are included for:
- `darwin/amd64`, `darwin/arm64`
- `linux/amd64` (glibc + musl), `linux/arm64` (glibc + musl)
- `windows/amd64`, `windows/arm64`
- `android/arm64`

If CGO is disabled (or you build with `-tags nocgo`), the package still compiles, but runtime operations return `fitz.ErrCGODisabled`.

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/brainwhocodes/gopher-pdf"
)

func main() {
	pdfBytes, err := os.ReadFile("statement.pdf")
	if err != nil {
		panic(err)
	}

	doc, err := gopherpdf.OpenBytes(pdfBytes)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	pages, _ := doc.NumPages()
	text, _ := doc.Text(0)
	png, _ := doc.RenderPagePNG(0, 200)

	fmt.Println("pages:", pages, "text_len:", len(text), "png_bytes:", len(png))
}
```

## API surface

Primary entry points:
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

## Development

```bash
go test ./...
go test -bench . ./...
```

Test fixtures live in `resources/`.

## Local development (replace)

If you want to consume a local checkout from another module, add a `replace` to your app's `go.mod`:

```go
replace github.com/brainwhocodes/gopher-pdf => /absolute/path/to/gopher-pdf
```

## Licensing

MuPDF is licensed under AGPL-3.0 (see `internal/fitz/COPYING`).
