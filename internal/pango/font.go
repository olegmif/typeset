package pango

/*
#include <pango/pango.h>
*/
import "C"

type Font struct {
	ptr *C.PangoFont
}
