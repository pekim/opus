package opus

// #cgo CFLAGS: -I${SRCDIR}/lib
// #cgo CFLAGS: -I${SRCDIR}/lib/src
// #cgo CFLAGS: -I${SRCDIR}/lib/include
// #cgo CFLAGS: -I${SRCDIR}/lib/celt
// #cgo CFLAGS: -I${SRCDIR}/lib/silk
// #cgo CFLAGS: -I${SRCDIR}/lib/silk/float
//
// #cgo CFLAGS: -I${SRCDIR}/lib/opusfile
// #cgo CFLAGS: -I${SRCDIR}/lib/opusfile/src
// #cgo CFLAGS: -I${SRCDIR}/lib/opusfile/include
//
// #cgo CFLAGS: -I${SRCDIR}/lib/ogg
// #cgo CFLAGS: -I${SRCDIR}/lib/ogg/src
// #cgo CFLAGS: -I${SRCDIR}/lib/ogg/include
//
// #cgo LDFLAGS: -lm
import "C"
