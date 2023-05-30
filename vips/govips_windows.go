//go:build windows

package vips

// #cgo CFLAGS: -I ./include
// #cgo CFLAGS: -I ./include/glib-2.0
// #cgo LDFLAGS: -L./libs -lvips-42 -lglib-2.0-0 -lgobject-2.0-0
import "C"
