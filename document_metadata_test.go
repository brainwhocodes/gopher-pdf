package gopherpdf

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Mirrors PyMuPDF tests/test_metadata.py against fixture 001003ED.pdf.
func TestPyMuPDFParity_Metadata(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	got, err := doc.Metadata()
	if err != nil {
		t.Fatalf("metadata: %v", err)
	}

	var expected map[string]any
	if err := json.Unmarshal(mustReadResource(t, "metadata.txt"), &expected); err != nil {
		t.Fatalf("unmarshal expected metadata: %v", err)
	}

	for _, key := range []string{"format", "title", "author", "subject", "keywords", "creator", "producer", "creationDate", "modDate", "trapped"} {
		want, ok := expected[key]
		if !ok {
			t.Fatalf("expected metadata missing key %q", key)
		}
		if got[key] != fmt.Sprint(want) {
			t.Fatalf("metadata mismatch for %q: got %q want %q", key, got[key], fmt.Sprint(want))
		}
	}
	// PyMuPDF metadata reports encryption as null when not encrypted.
	if got["encryption"] != "" && got["encryption"] != "None" {
		t.Fatalf("unexpected encryption metadata: %q", got["encryption"])
	}
}

// Mirrors tests/test_metadata.py metadata erase and update expectations.
func TestPyMuPDFParity_MetadataEraseAndSet(t *testing.T) {
	doc, err := Open(resourcePath(t, "001003ED.pdf"))
	if err != nil {
		t.Fatalf("open path: %v", err)
	}
	defer doc.Close()

	before, err := doc.Metadata()
	if err != nil {
		t.Fatalf("metadata before: %v", err)
	}
	if strings.TrimSpace(before["title"]) == "" {
		t.Fatalf("expected fixture title to be non-empty before erase")
	}

	if err := doc.SetMetadata(map[string]string{}); err != nil {
		t.Fatalf("set empty metadata: %v", err)
	}

	afterErase, err := doc.Metadata()
	if err != nil {
		t.Fatalf("metadata after erase: %v", err)
	}
	for _, key := range []string{"title", "author", "subject", "keywords", "creator", "producer", "creationDate", "modDate", "trapped"} {
		if afterErase[key] != "" {
			t.Fatalf("expected metadata key %q to be erased, got %q", key, afterErase[key])
		}
	}
	if strings.TrimSpace(afterErase["format"]) == "" {
		t.Fatalf("expected format metadata to remain populated")
	}
	infoValue, err := doc.XrefGetKey(-1, "Info")
	if err != nil {
		t.Fatalf("xref get key trailer Info: %v", err)
	}
	trailerKeys, err := doc.XrefGetKeys(-1)
	if err != nil {
		t.Fatalf("xref get keys trailer: %v", err)
	}
	hasInfo := false
	for _, k := range trailerKeys {
		if k == "Info" {
			hasInfo = true
			break
		}
	}
	// PyMuPDF accepts either /Info=null or complete absence of /Info.
	statement1 := infoValue == "null"
	statement2 := !hasInfo
	if !(statement1 || statement2) {
		t.Fatalf("expected trailer /Info to be null or absent after erase; info=%q hasInfo=%v", infoValue, hasInfo)
	}

	update := map[string]string{
		"title":        "Parity Title",
		"author":       "Parity Author",
		"subject":      "Parity Subject",
		"keywords":     "k1,k2",
		"creator":      "Parity Creator",
		"producer":     "Parity Producer",
		"creationDate": "D:20260101000000+00'00'",
		"modDate":      "D:20260102000000+00'00'",
		"trapped":      "False",
	}
	if err := doc.SetMetadata(update); err != nil {
		t.Fatalf("set updated metadata: %v", err)
	}
	afterSet, err := doc.Metadata()
	if err != nil {
		t.Fatalf("metadata after set: %v", err)
	}
	for k, want := range update {
		if got := afterSet[k]; got != want {
			t.Fatalf("metadata mismatch for %q: got %q want %q", k, got, want)
		}
	}
}

