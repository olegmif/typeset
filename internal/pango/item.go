package pango

/*
#include <pango/pango.h>
*/
import "C"

type Item struct {
	ptr *C.PangoItem
}

type ItemData struct {
	Offset   int
	Length   int
	NumChars int
	Analysis *Analysis
}

func (i *Item) GetItemData() *ItemData {
	if i == nil || i.ptr == nil {
		panic("nil Item")
	}

	return &ItemData{
		Offset:   int(i.ptr.offset),
		Length:   int(i.ptr.length),
		NumChars: int(i.ptr.num_chars),
		Analysis: &Analysis{analysis: i.ptr.analysis},
	}
}

func (i *Item) GetAnalysis() *Analysis {
	if i == nil || i.ptr == nil {
		panic("nil Item")
	}

	return &Analysis{analysis: i.ptr.analysis}
}
