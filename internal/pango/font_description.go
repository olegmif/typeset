package pango

/*
#include <stdlib.h>
#include <pango/pango.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

type FontDescription struct {
	ptr *C.PangoFontDescription
}

func NewFontDescriptionFromString(descStr string) (*FontDescription, error) {
	cDescStr := C.CString(descStr)
	defer C.free(unsafe.Pointer(cDescStr))

	ptr := C.pango_font_description_from_string(cDescStr)
	if ptr == nil {
		return nil, errors.New("pango_font_description_from_string returns nil")
	}

	return &FontDescription{ptr: ptr}, nil
}

func (fd *FontDescription) Close() {
	if fd != nil && fd.ptr != nil {
		C.pango_font_description_free(fd.ptr)
		fd.ptr = nil
	}
}
