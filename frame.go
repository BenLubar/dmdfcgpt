package main

import (
	"sync/atomic"
	"unsafe"

	"github.com/nsf/termbox-go"
)

// Each function returns one of the following patterns:
//
// - any Frame, non-nil error: fatal error occurred, program must terminate.
// - nil Frame, nil error: the action completed but the Frame is unchanged.
// - non-nil Frame, nil error: the action completed and the Frame changed.
type Frame interface {
	Render(x0, y0, x1, y1 int) (Frame, error)
	Mouse(mx, my, w, h int) (Frame, error)
	Key(key termbox.Key, mod termbox.Modifier) (Frame, error)
	Ch(ch rune, mod termbox.Modifier) (Frame, error)
}

var currentFrame = new(Frame)

func CurrentFrame() *Frame {
	return (*Frame)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&currentFrame))))
}

func SetFrame(old, new *Frame) bool {
	if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&currentFrame)), unsafe.Pointer(old), unsafe.Pointer(new)) {
		select {
		case rerender <- struct{}{}:
		default:
		}
		return true
	}
	return false
}
