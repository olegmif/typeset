package pango

/*
#include <pango/pango.h>
*/
import "C"

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r *Rectangle) setFromC(cr *C.PangoRectangle) {
	if cr == nil {
		*r = Rectangle{}
		return
	}
	r.X = int(cr.x)
	r.Y = int(cr.y)
	r.Width = int(cr.width)
	r.Height = int(cr.height)
}

func rectFromC(cr *C.PangoRectangle) Rectangle {
	var r Rectangle
	r.setFromC(cr)
	return r
}
