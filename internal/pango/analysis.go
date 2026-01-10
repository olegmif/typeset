package pango

/*
#include <pango/pango.h>
*/
import "C"

import (
	"unsafe"
)

type Analysis struct {
	analysis C.PangoAnalysis
}

type Language struct {
	ptr *C.PangoLanguage
}

func (a *Analysis) GetFont() *Font {
	if a == nil {
		panic("nil Analysis")
	}

	return &Font{ptr: a.analysis.font}
}

func (a *Analysis) GetLevel() uint8 {
	if a == nil {
		panic("nil Analysis")
	}

	return uint8(a.analysis.level)
}

func (a *Analysis) GetGravity() uint8 {
	if a == nil {
		panic("nil Analysis")
	}

	return uint8(a.analysis.gravity)
}

func (a *Analysis) GetFlags() uint8 {
	if a == nil {
		panic("nil Analysis")
	}

	return uint8(a.analysis.flags)
}

func (a *Analysis) GetScript() uint8 {
	if a == nil {
		panic("nil Analysis")
	}

	return uint8(a.analysis.script)
}

func (a *Analysis) GetLanguage() *Language {
	if a == nil {
		panic("nil Analysis")
	}

	return &Language{ptr: a.analysis.language}
}

func (a *Analysis) GetExtraAttrs() unsafe.Pointer {
	if a == nil {
		panic("nil Analysis")
	}

	return unsafe.Pointer(a.analysis.extra_attrs)
}
