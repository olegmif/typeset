//go:build linux

package main

/*
#cgo pkg-config: pangoft2 pango glib-2.0 gobject-2.0
#include <stdlib.h>
#include <string.h>

#include <pango/pango.h>
#include <pango/pangoft2.h>
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func writePGM(path string, w, h int, buf []byte) error {
	// P5 = binary PGM
	header := fmt.Sprintf("P5\n%d %d\n255\n", w, h)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(header); err != nil {
		return err
	}

	_, err = f.Write(buf)
	return err
}

func main() {
	// --- 1) Font map + context (PangoFT2)
	fontmap := C.pango_ft2_font_map_new()
	if fontmap == nil {
		panic("pango_ft2_font_map_new returned nil")
	}
	defer C.g_object_unref(C.gpointer(fontmap))

	ctx := C.pango_font_map_create_context((*C.PangoFontMap)(fontmap))
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

	text := "Pango + FreeType (FT2)\nБез Cairo.\nВыравнивание по ширине/переносы - на стороне Pango."
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	C.pango_layout_set_text(layout, ctext, -1)

	// --- 3) Настройки layout (ширина, перенос, выравнивание)
	// ширина в Pango units (1px = PANGO_SCALE)
	const widthPx = 520
	C.pango_layout_set_width(layout, C.int(widthPx*C.PANGO_SCALE))
	C.pango_layout_set_wrap(layout, C.PANGO_WRAP_WORD_CHAR)
	C.pango_layout_set_alignment(layout, C.PANGO_ALIGN_LEFT)

	// --- 4) Шрифт (через PangoFontDescription)
	// Можно заменить на другой шрифт
	descStr := C.CString("DejaVu Sans 22")
	defer C.free(unsafe.Pointer(descStr))
	desc := C.pango_font_description_from_string(descStr)
	if desc == nil {
		panic("pango_font_description_from_string returned nil")
	}
	defer C.pango_font_description_free(desc)
	C.pango_layout_set_font_description(layout, desc)

	// --- 5) Узнаем пиксельный размер layout
	var inkRect C.PangoRectangle
	var logRect C.PangoRectangle
	C.pango_layout_get_pixel_extents(layout, &inkRect, &logRect)

	// Берем логический прямоугольник - это "сколько места занял layout"
	w := int(logRect.width)
	h := int(logRect.height)
	if w <= 0 || h <= 0 {
		panic(fmt.Sprintf("bad layout size: %dx%d", w, h))
	}

	// --- 6) Готовим FT_Bitmap (grayscale)
	// Важно: pitch может быть >= width (выравнивание), поэтому считаем pitch = width.
	var bmp C.FT_Bitmap
	C.memset(unsafe.Pointer(&bmp), 0, C.size_t(unsafe.Sizeof(bmp)))
	bmp.width = C.uint(w)
	bmp.rows = C.uint(h)
	bmp.pitch = C.int(w)
	bmp.pixel_mode = C.FT_PIXEL_MODE_GRAY
	bmp.num_grays = 256

	bufSize := w * h
	buf := C.malloc(C.size_t(bufSize))
	if buf == nil {
		panic("malloc failed")
	}
	defer C.free(buf)
	C.memset(buf, 0, C.size_t(bufSize))
	bmp.buffer = (*C.uchar)(buf)

	// --- 7) Рендерим layout в bitmap
	// Сдвиг (x,y) - в пикселях относительно буфера. Рендерим с (0,0).
	// PangoFT2 рисует альфу в grayscale-буфер.
	C.pango_ft2_render_layout(&bmp, layout, 0, 0)

	// --- 8) Копируем буфер в Go и пишем PGM
	goBuf := unsafe.Slice((*byte)(buf), bufSize)
	out := "out.pgm"
	check(writePGM(out, w, h, goBuf))

	fmt.Printf("Wrote %s (%dx%d)\n", out, w, h)
	fmt.Println("Preview (пример):")
	fmt.Println("   feh out.pgm")
	fmt.Println("Convert to PNG:")
	fmt.Println("   magick out.pgm out.png")
}
