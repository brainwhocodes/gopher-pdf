//go:build cgo && !nocgo

package fitz

/*
#include <stdlib.h>
*/
import "C"

import "unsafe"

func withCString(s string, fn func(*C.char)) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	fn(cs)
}

