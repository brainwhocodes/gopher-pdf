package gopherpdf

import (
	"path/filepath"
	"runtime"
	"testing"
)

func benchResourcePath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	backendPkg := filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
	return filepath.Join(backendPkg, "gopher-pdf", "resources", name)
}

func BenchmarkGetText(b *testing.B) {
	doc, err := Open(benchResourcePath("test_2548.pdf"))
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer doc.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := doc.GetText(0, "text", 0); err != nil {
			b.Fatalf("get text: %v", err)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	doc, err := Open(benchResourcePath("test_2548.pdf"))
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer doc.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := doc.Search(0, "the", TextFlagsSearch, 512, nil); err != nil {
			b.Fatalf("search: %v", err)
		}
	}
}

func BenchmarkRenderPagePNG(b *testing.B) {
	doc, err := Open(benchResourcePath("test_2548.pdf"))
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer doc.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := doc.RenderPagePNG(0, 72); err != nil {
			b.Fatalf("render: %v", err)
		}
	}
}

func BenchmarkXrefGetKeyTyped(b *testing.B) {
	doc, err := Open(benchResourcePath("2.pdf"))
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer doc.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := doc.XrefGetKeyTyped(-1, "Root"); err != nil {
			b.Fatalf("xref get key typed: %v", err)
		}
	}
}
