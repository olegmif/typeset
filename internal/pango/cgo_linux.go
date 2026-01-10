//go:build linux

package pango

/*
#cgo pkg-config: pango
#cgo pkg-config: pangoft2
#include <stdlib.h>
#include <pango/pango.h>
#include <pango/pangoft2.h>
*/
import "C"
