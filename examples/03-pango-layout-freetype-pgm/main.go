//go:build linux

package main

/*
#cgo pkg-config: freetype2
#include <stdlib.h>
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/olegmif/typeset/internal/pango"
)

func findDejaVuSansTTF() (string, error) {
	candidates := []string{
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
	}

	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}

	return "", errors.New("cannod find DejaVuSans.ttf")
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func blendBack(dst *uint8, alpha uint8) {
	if alpha == 0 {
		return
	}
	d := int(*dst)
	a := int(alpha)
	*dst = uint8((d * (255 - a)) / 255)
}

func render(g pango.Glyph, x int, y int, face C.FT_Face, imgW int, imgH int, buffer []uint8) error {
	if e := C.FT_Load_Glyph(face, C.FT_UInt(g), C.FT_LOAD_DEFAULT); e != 0 {
		return errors.New("FT_Load_Glyph returns error")
	}

	if e := C.FT_Render_Glyph(face.glyph, C.FT_RENDER_MODE_NORMAL); e != 0 {
		return errors.New("FT_Render_Glyph returns error")
	}

	slot := face.glyph
	bm := slot.bitmap

	bw := int(bm.width)
	bh := int(bm.rows)
	pitch := int(bm.pitch)

	left := int(slot.bitmap_left)
	top := int(slot.bitmap_top)

	dstX0 := x + left
	dstY0 := y - top

	if bw > 0 && bh > 0 && bm.buffer != nil {
		buf := unsafe.Slice((*uint8)(unsafe.Pointer(bm.buffer)), bh*pitch)

		for row := 0; row < bh; row++ {
			dstY := dstY0 + row
			if dstY < 0 || dstY >= imgH {
				continue
			}

			srcRow := buf[row*pitch : row*pitch+bw]

			for col := 0; col < bw; col++ {
				dstX := dstX0 + col
				if dstX < 0 || dstX >= imgW {
					continue
				}

				alpha := srcRow[col]
				if alpha == 0 {
					continue
				}

				idx := dstY*imgW + dstX
				blendBack(&buffer[idx], alpha)
			}
		}
	}

	return nil
}

func doWithGlyph(g pango.GlyphInfo, baseline int, penX int, face C.FT_Face, imgW int, imgH int, buffer []uint8) int {
	glyph := g.GetGlyph()
	fmt.Printf("GlyphID: %d, ", glyph)

	scale := pango.GetPangoScale()
	geom := g.GetGeometry()
	nextOffset := penX + geom.Width

	x := penX + int(geom.XOffset)
	y := baseline + int(geom.YOffset)

	fmt.Printf("(x, y): (%d, %d)\n", x/scale, y/scale)

	if err := render(glyph, x/scale, y/scale, face, imgW, imgH, buffer); err != nil {
		return nextOffset
	}

	return nextOffset
}

func doWithRun(r *pango.Run, i *pango.Iter, face C.FT_Face, imgW int, imgH int, buffer []uint8) {
	gc := r.GlyphCount()
	b := i.GetBaseline()
	bPx := b / pango.GetPangoScale()
	fmt.Printf("run glyphs: %d, baseline: %d\n", gc, bPx)
	offset := 0
	for i := 0; i < gc; i++ {
		gi := r.GlyphInfoAt(i)
		offset = doWithGlyph(gi, b, offset, face, imgW, imgH, buffer)
	}
}

func main() {
	outPath := "out.pgm"

	//
	// 1) Init Pango layout
	//

	fontmap, err := pango.NewFontMap()
	if err != nil {
		panic(err)
	}
	defer fontmap.Close()

	ctx, err := pango.NewContext(fontmap)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	layout, err := pango.NewLayout(ctx)
	if err != nil {
		panic(err)
	}
	defer layout.Close()

	layout.SetText("Один Два Три Четыре Пять Шесть Семь Восемь Девять Десять\nОдиннадцать Двенадцать")

	const widthPx = 520
	layout.SetWidth(widthPx)
	layout.SetWrap(pango.WrapWordChar)
	layout.SetAlignment(pango.AlignLeft)

	desc, err := pango.NewFontDescriptionFromString("DejaVu Sans 22")
	if err != nil {
		panic(err)
	}
	defer desc.Close()

	layout.SetFontDescription(desc)

	ext := layout.GetPixelExtents()

	imgW := ext.Logical.Width
	imgH := ext.Logical.Height

	pix := make([]uint8, imgW*imgH)
	for i := range pix {
		pix[i] = 255
	}

	//
	// 2) Init FreeType
	//

	fontPath, err := findDejaVuSansTTF()
	if err != nil {
		panic(err)
	}
	fmt.Println("Using font:", fontPath)

	var ftLib C.FT_Library
	if e := C.FT_Init_FreeType(&ftLib); e != 0 {
		panic(fmt.Errorf("FT_Init_FreeType failed: %d", int(e)))
	}
	defer C.FT_Done_FreeType(ftLib)

	cFontPath := C.CString(fontPath)
	defer C.free(unsafe.Pointer(cFontPath))

	var face C.FT_Face
	if e := C.FT_New_Face(ftLib, cFontPath, 0, &face); e != 0 {
		panic(fmt.Errorf("FT_New_Face failed: %d", int(e)))
	}
	defer C.FT_Done_Face(face)

	// 22px по высоте
	if e := C.FT_Set_Pixel_Sizes(face, 0, 22); e != 0 {
		panic(fmt.Errorf("FT_Set_Pixel_Sizes failed: %d", int(e)))
	}

	//
	// 3) Интеграция Pango runs -> glyphs -> FT raster -> blit
	//

	iter, err := layout.GetIter()
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for {
		run := iter.GetRunReadonly()

		if run != nil {
			doWithRun(run, iter, face, imgW, imgH, pix)
		}

		hr := iter.NextRun()
		if !hr {
			break
		}
	}

	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	header := fmt.Sprintf("P5\n%d %d\n255\n", imgW, imgH)
	if _, err := f.WriteString(header); err != nil {
		panic(err)
	}
	if _, err := f.Write(pix); err != nil {
		panic(err)
	}

	fmt.Println("Wrote:", outPath)
	_ = clamp // (на случай, если понадобится расширить)
}
