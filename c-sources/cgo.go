package opus

// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opus
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opus/celt
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opus/include
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opus/src
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opus/silk
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opus/silk/float
//
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opusfile/src
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/opusfile/include
//
// #cgo CFLAGS: -I${SRCDIR}/../xiph.org/ogg/src
//
// #cgo LDFLAGS: -lm
import "C"
