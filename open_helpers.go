package gopherpdf

import (
	"errors"
	"fmt"

	fitz "bank-statments/backend/pkg/gopher-pdf/internal/fitz"
)

func openFitzMaybeAuth(source string, doc *fitz.Document, openErr error, password string) (*fitz.Document, error) {
	if openErr == nil {
		return doc, nil
	}

	if errors.Is(openErr, fitz.ErrNeedsPassword) && doc != nil {
		if password == "" {
			_ = doc.Close()
			return nil, fmt.Errorf("%s: password required", source)
		}
		if ok := doc.AuthenticatePassword(password); !ok {
			_ = doc.Close()
			return nil, fmt.Errorf("%s: invalid password", source)
		}
		return doc, nil
	}

	if doc != nil {
		_ = doc.Close()
	}
	return nil, fmt.Errorf("%s: %w", source, openErr)
}
