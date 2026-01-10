package pango

/*
#include <pango/pango.h>
*/
import "C"

type GlyphInfo struct {
	ptr *C.PangoGlyphInfo
}

type Glyph = uint32

type GlyphGeometry struct {
	Width   int
	XOffset int
	YOffset int
}

// GetGlyph - В pango Glyph - это идентификатор глифа
func (gi *GlyphInfo) GetGlyph() Glyph {
	if gi == nil || gi.ptr == nil {
		panic("nil GlyphInfo")
	}

	return uint32(gi.ptr.glyph)
}

func (gi *GlyphInfo) GetGeometry() GlyphGeometry {
	if gi == nil || gi.ptr == nil {
		panic("nil GlyphInfo")
	}

	return GlyphGeometry{
		Width:   int(gi.ptr.geometry.width),
		XOffset: int(gi.ptr.geometry.x_offset),
		YOffset: int(gi.ptr.geometry.y_offset),
	}
}
