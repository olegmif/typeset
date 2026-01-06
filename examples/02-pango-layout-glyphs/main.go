//go:build linux

package main

import (
	"fmt"

	"github.com/olegmif/typeset/internal/pango"
)

func main() {
	// --- 1) Font map + context (без PangoFT2)
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

	// --- 2) Layout
	layout, err := pango.NewLayout(ctx)
	if err != nil {
		panic(err)
	}
	defer layout.Close()

	layout.SetText("Pango layout only\nБез Cairo / без PangoFT2.\nТут мы получаем только глифы и их позиции.")

	// --- 3) Параметры layout
	const widthPx = 520
	layout.SetWidth(widthPx)
	layout.SetWrap(pango.WrapChar)
	layout.SetAlignment(pango.AlignLeft)

	// --- 4) Шрифт (как и раньше, через PangoFontDescription)
	desc, err := pango.NewFontDescriptionFromString("DejaVu Sans 22")
	if err != nil {
		panic(err)
	}
	defer desc.Close()

	layout.SetFontDescription(desc)

	// --- 5) Размеры layout (в пикселях)
	extents := layout.GetPixelExtents()
	fmt.Printf("Layout logical size: %dx%d\n", extents.Logical.Width, extents.Logical.Height)
	fmt.Printf("PANGO_SCALE=%d (units per px)\n\n", pango.GetPangoScale())

	// --- 6) Итерируем layout И печатаем глифы с позициями
	iter, err := layout.GetIter()
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for {
		// baseline в Pango units
		baseline := iter.GetBaseline()

		if hasNext := iter.NextRun(); !hasNext {
			break
		}

		run := iter.GetRunReadonly()

		if run != nil {

			// font := run.GetItem().GetAnalysis().GetFont()

			n := run.GlyphCount()

			penX := 0

			for i := 0; i < n; i++ {
				g := run.GlyphInfoAt(i)

				glyphID := g.GetGlyph()
				geometry := g.GetGeometry()

				x := penX + int(geometry.XOffset)
				y := baseline + int(geometry.YOffset)
				scale := pango.GetPangoScale()

				fmt.Printf("glyph=%d x=%.2f y=%.2f\n",
					glyphID,
					float64(x)/float64(scale),
					float64(y)/float64(scale),
				)

				penX += int(geometry.Width)
			}
		}
	}
}
