// Package pango ... package comment here
package pango

/*
#include <pango/pango.h>
*/
import "C"

import "unsafe"

type GlyphItem struct {
	ptr *C.PangoGlyphItem
}

func (gi *GlyphItem) GetItem() *Item {
	if gi == nil || gi.ptr == nil {
		panic("nil GlyphItem")
	}

	return &Item{ptr: gi.ptr.item}
}

// GetGlyphs - то аллоцирующий метод, лучше не использовать часто.
// Если нужно обойти все GlyphInfo, воспользуйтесь GlyphCount и GlyphInfoAt
func (gi *GlyphItem) GetGlyphs() []GlyphInfo {
	if gi == nil || gi.ptr == nil {
		panic("nil GlyphItem")
	}

	gs := gi.ptr.glyphs
	if gs == nil {
		return nil
	}

	n := int(gs.num_glyphs)
	if n == 0 || gs.glyphs == nil {
		return nil
	}

	infos := unsafe.Slice((*C.PangoGlyphInfo)(unsafe.Pointer(gs.glyphs)), n)

	out := make([]GlyphInfo, n)
	for i := range infos {
		out[i] = GlyphInfo{ptr: &infos[i]}
	}

	return out
}

func (gi *GlyphItem) GlyphCount() int {
	if gi == nil || gi.ptr == nil {
		panic("nil GlyphItem")
	}
	gs := gi.ptr.glyphs
	if gs == nil {
		return 0
	}
	return int(gs.num_glyphs)
}

func (gi *GlyphItem) GlyphInfoAt(i int) GlyphInfo {
	if gi == nil || gi.ptr == nil {
		panic("nil GlyphItem")
	}

	gs := gi.ptr.glyphs
	if gs == nil || gs.glyphs == nil {
		panic("no glyphs")
	}

	n := int(gs.num_glyphs)
	if i < 0 || i >= n {
		panic("glyph index out of range")
	}

	base := (*C.PangoGlyphInfo)(unsafe.Pointer(gs.glyphs))
	return GlyphInfo{ptr: (*C.PangoGlyphInfo)(unsafe.Add(unsafe.Pointer(base), uintptr(i)*unsafe.Sizeof(*base)))}
}
