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

type PangoLayout struct {
	ptr *C.PangoLayout
}

type WrapMode int

const (
	WrapWord     WrapMode = WrapMode(C.PANGO_WRAP_WORD)
	WrapChar     WrapMode = WrapMode(C.PANGO_WRAP_CHAR)
	WrapWordChar WrapMode = WrapMode(C.PANGO_WRAP_WORD_CHAR)
	WrapNone     WrapMode = WrapMode(C.PANGO_WRAP_NONE)
)

type Alignment int

const (
	AlignLeft   Alignment = Alignment(C.PANGO_ALIGN_LEFT)
	AlignCenter Alignment = Alignment(C.PANGO_ALIGN_CENTER)
	AlignRight  Alignment = Alignment(C.PANGO_ALIGN_RIGHT)
)

type Extents struct {
	Ink     Rectangle
	Logical Rectangle
}

func NewLayout(ctx *PangoContext) (*PangoLayout, error) {
	if ctx == nil || ctx.ptr == nil {
		return nil, errors.New("nil PangoContext")
	}

	layout := C.pango_layout_new(ctx.ptr)
	if layout == nil {
		return nil, errors.New("pango_layout_new returned nil")
	}

	return &PangoLayout{ptr: layout}, nil
}

func (l *PangoLayout) Close() {
	if l != nil && l.ptr != nil {
		C.g_object_unref(C.gpointer(l.ptr))
		l.ptr = nil
	}
}

func (l *PangoLayout) SetText(text string) error {
	if l == nil || l.ptr == nil {
		return errors.New("nil PangoLayout")
	}

	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	C.pango_layout_set_text(l.ptr, ctext, -1)
	return nil
}

func (l *PangoLayout) SetWrap(mode WrapMode) {
	if l != nil && l.ptr != nil {
		C.pango_layout_set_wrap(l.ptr, C.PangoWrapMode(mode))
	}
}

func (l *PangoLayout) SetAlignment(alignment Alignment) {
	if l != nil && l.ptr != nil {
		C.pango_layout_set_alignment(l.ptr, C.PangoAlignment(alignment))
	}
}

func (l *PangoLayout) SetJustify(justify bool) {
	if l != nil && l.ptr != nil {

		var v C.gboolean
		if justify {
			v = C.gboolean(1)
		} else {
			v = C.gboolean(0)
		}

		C.pango_layout_set_justify(l.ptr, v)
	}
}

func (l *PangoLayout) SetWidth(width int) {
	if l != nil && l.ptr != nil {
		C.pango_layout_set_width(l.ptr, C.int(width*C.PANGO_SCALE))
	}
}

func (l *PangoLayout) SetFontDescription(desc *FontDescription) error {
	if l == nil || l.ptr == nil {
		return errors.New("nil Layout")
	}

	if desc == nil || desc.ptr == nil {
		return errors.New("nil FontDescription")
	}

	C.pango_layout_set_font_description(l.ptr, desc.ptr)
	return nil
}

func (l *PangoLayout) GetPixelExtents() Extents {
	if l == nil || l.ptr == nil {
		panic("nil Layout")
	}

	var irC, lrC C.PangoRectangle

	C.pango_layout_get_pixel_extents(l.ptr, &irC, &lrC)

	return Extents{
		Ink:     rectFromC(&irC),
		Logical: rectFromC(&lrC),
	}
}

func (l *PangoLayout) GetIter() (*Iter, error) {
	if l == nil || l.ptr == nil {
		return nil, errors.New("nil Layout")
	}

	ptr := C.pango_layout_get_iter(l.ptr)
	if ptr == nil {
		return nil, errors.New("pango_layout_get_inter returns nil")
	}

	return &Iter{ptr: ptr}, nil
}
