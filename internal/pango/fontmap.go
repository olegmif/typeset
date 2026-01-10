package pango

/*
#include <pango/pangoft2.h>
*/
import "C"

import "errors"

type FontMap struct {
	ptr *C.PangoFontMap
}

// На Linux это обычно PangoFcFontMap (fontconfig) под капотом.
/*
	fontmap := C.pango_ft2_font_map_new()
	if fontmap == nil {
		panic("pango_ft2_font_map_new returned nil")
	}
	defer C.g_object_unref(C.gpointer(fontmap))
*/

func NewFontMap() (*FontMap, error) {
	fm := C.pango_ft2_font_map_new()
	if fm == nil {
		return nil, errors.New("pango_ft2_font_map_new returned nil")
	}

	return &FontMap{
		ptr: (*C.PangoFontMap)(fm),
	}, nil
}

func (fm *FontMap) Close() {
	if fm.ptr != nil {
		C.g_object_unref(C.gpointer(fm.ptr))
		fm.ptr = nil
	}
}
