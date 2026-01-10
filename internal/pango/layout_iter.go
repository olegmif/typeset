package pango

/*
#include <pango/pango.h>
*/
import "C"

import (
	"unsafe"
)

type Iter struct {
	ptr *C.PangoLayoutIter
}

type Run = GlyphItem

func (i *Iter) Close() {
	if i != nil && i.ptr != nil {
		C.pango_layout_iter_free(i.ptr)
		i.ptr = nil
	}
}

func (i *Iter) GetBaseline() int {
	if i == nil || i.ptr == nil {
		panic("nil Iter")
	}

	return int(C.pango_layout_iter_get_baseline(i.ptr))
}

func (i *Iter) GetRunReadonly() *Run {
	if i == nil || i.ptr == nil {
		panic("nil Iter")
	}

	ptr := C.pango_layout_iter_get_run_readonly(i.ptr)
	if ptr == nil {
		return nil
	}

	return &GlyphItem{ptr: (*C.PangoGlyphItem)(unsafe.Pointer(ptr))}
}

func (i *Iter) NextRun() bool {
	if i == nil || i.ptr == nil {
		panic("nil Iter")
	}

	return C.pango_layout_iter_next_run(i.ptr) != 0
}
