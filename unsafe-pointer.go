package opus

import "C"

import (
	"unsafe"
)

// intToUnsafePointer converts an int to unsafe.Pointer.
//
// The conversion is performed  without triggering the unsafeptr vet warning ("possible misuse of unsafe.pointer").
func intToUnsafePointer(i int) unsafe.Pointer {
	notReallyAnAddress := uintptr(i)
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&notReallyAnAddress))
	return ptr
}
