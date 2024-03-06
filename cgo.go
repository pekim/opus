package opus

// #cgo CFLAGS: -I${SRCDIR}/lib/opus
// #cgo CFLAGS: -I${SRCDIR}/lib/opus/src
// #cgo CFLAGS: -I${SRCDIR}/lib/opus/include
// #cgo CFLAGS: -I${SRCDIR}/lib/opus/celt
// #cgo CFLAGS: -I${SRCDIR}/lib/opus/silk
// #cgo CFLAGS: -I${SRCDIR}/lib/opus/silk/float
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
