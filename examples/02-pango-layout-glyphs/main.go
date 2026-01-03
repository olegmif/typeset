//go:build linux

package main

/*
#cgo pkg-config: pango pangoft2 glib-2.0 gobject-2.0
#include <stdlib.h>
#include <pango/pango.h>
#include <pango/pangoft2.h>
#include <glib.h>

static const char* font_desc_string(PangoFont* font) {
	PangoFontDescription* d = pango_font_describe(font);
	if (!d) return NULL;
	char* s = pango_font_description_to_string(d);
	pango_font_description_free(d);
	// caller must g_free(s)
	return s;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	// --- 1) Font map + context (без PangoFT2)
	// На Linux это обычно PangoFcFontMap (fontconfig) под капотом.
	fontmap := C.pango_ft2_font_map_new()
	if fontmap == nil {
		panic("pango_ft2_font_map_new returned nil")
	}
	defer C.g_object_unref(C.gpointer(fontmap))

	ctx := C.pango_font_map_create_context(fontmap)
	if ctx == nil {
		panic("pango_font_map_create_context returned nil")
	}
	defer C.g_object_unref(C.gpointer(ctx))

	// --- 2) Layout
	layout := C.pango_layout_new(ctx)
	if layout == nil {
		panic("pango_layout_new returned nil")
	}
	defer C.g_object_unref(C.gpointer(layout))

	text := "Pango layout only\nБез Cairo / без PangoFT2.\nТут мы получаем только глифы и их позиции."
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.pango_layout_set_text(layout, ctext, -1)

	// --- 3) Параметры layout
	const widthPx = 520
	C.pango_layout_set_width(layout, C.int(widthPx*C.PANGO_SCALE))
	C.pango_layout_set_wrap(layout, C.PANGO_WRAP_WORD_CHAR)
	C.pango_layout_set_alignment(layout, C.PANGO_ALIGN_LEFT)

	// --- 4) Шрифт (как и раньше, через PangoFontDescription)
	descStr := C.CString("DejaVu Sans 22")
	defer C.free(unsafe.Pointer(descStr))
	desc := C.pango_font_description_from_string(descStr)
	if desc == nil {
		panic("pango_font_description_from_string returned nil")
	}
	defer C.pango_font_description_free(desc)
	C.pango_layout_set_font_description(layout, desc)

	// --- 5) Размеры layout (в пикселях)
	var incRect C.PangoRectangle
	var logRect C.PangoRectangle
	C.pango_layout_get_pixel_extents(layout, &incRect, &logRect)

	fmt.Printf("Layout logical size: %dx%d", int(logRect.width), int(logRect.height))
	fmt.Printf("PANGO_SCALE=%d (units per px)\n\n", int(C.PANGO_SCALE))

	// --- 6) Итерируем layout И печатаем глифы с позициями
	iter := C.pango_layout_get_iter(layout)
	if iter == nil {
		panic("pango_layout_get_iter returned nil")
	}
	defer C.pango_layout_iter_free(iter)
	for {
		// baseline в Pango units
		baseline := int(C.pango_layout_iter_get_baseline(iter))

		run := C.pango_layout_iter_get_run_readonly(iter)
		if run != nil {
			gi := (*C.PangoGlyphItem)(unsafe.Pointer(run))

			font := gi.item.analysis.font
			_ = font // позже пригодится

			gs := gi.glyphs
			n := int(gs.num_glyphs)

			penX := 0
			for i := 0; i < n; i++ {
				g := (*C.PangoGlyphInfo)(unsafe.Pointer(
					uintptr(unsafe.Pointer(gs.glyphs)) + uintptr(i)*unsafe.Sizeof(*gs.glyphs),
				))

				glyphID := uint32(g.glyph)
				x := penX + int(g.geometry.x_offset)
				y := baseline + int(g.geometry.y_offset)

				fmt.Printf("glyph=%d x=%.2f y=%.2f\n",
					glyphID,
					float64(x)/float64(C.PANGO_SCALE),
					float64(y)/float64(C.PANGO_SCALE),
				)

				penX += int(g.geometry.width)
			}
		}

		if C.pango_layout_iter_next_run(iter) == 0 {
			break
		}
	}
}
