package pango

/*
#include <pango/pango.h>
*/
import "C"
import "errors"

type PangoContext struct {
	ptr *C.PangoContext
}

func NewContext(fm *FontMap) (*PangoContext, error) {
	if fm == nil || fm.ptr == nil {
		return nil, errors.New("nil FontMap")
	}

	ctx := C.pango_font_map_create_context(fm.ptr)
	if ctx == nil {
		return nil, errors.New("pango_font_map_create_context returned nil")
	}

	return &PangoContext{ptr: ctx}, nil
}

func (ctx *PangoContext) Close() {
	if ctx.ptr != nil {
		C.g_object_unref(C.gpointer(ctx.ptr))
		ctx.ptr = nil
	}
}
