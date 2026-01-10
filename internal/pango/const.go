package pango

/*
#include <pango/pango.h>
*/
import "C"

func GetPangoScale() int {
	return int(C.PANGO_SCALE)
}
