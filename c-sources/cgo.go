package opus

// #cgo CFLAGS: -I${SRCDIR}/../lib/opus
// #cgo CFLAGS: -I${SRCDIR}/../lib/opus/celt
// #cgo CFLAGS: -I${SRCDIR}/../lib/opus/include
// #cgo CFLAGS: -I${SRCDIR}/../lib/opus/src
// #cgo CFLAGS: -I${SRCDIR}/../lib/opus/silk
// #cgo CFLAGS: -I${SRCDIR}/../lib/opus/silk/float
//
// #cgo CFLAGS: -I${SRCDIR}/../lib/opusfile/src
// #cgo CFLAGS: -I${SRCDIR}/../lib/opusfile/include
//
// #cgo CFLAGS: -I${SRCDIR}/../lib/ogg/src
//
// #cgo LDFLAGS: -lm
import "C"
