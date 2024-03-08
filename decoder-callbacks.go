package opus

// #include <opusfile.h>
import "C"

import (
	"fmt"
	"io"
	"sync"
	"unsafe"
)

var decoderId int
var decoderInstanceMap = make(map[int]*Decoder)
var decoderLock sync.Mutex

func getDecoderForId(id unsafe.Pointer) *Decoder {
	decoderLock.Lock()
	defer decoderLock.Unlock()

	return decoderInstanceMap[int(uintptr(id))]
}

func getDecoderId(decoder *Decoder) unsafe.Pointer {
	decoderLock.Lock()
	defer decoderLock.Unlock()

	decoderId++
	decoderInstanceMap[decoderId] = decoder

	return intToUnsafePointer(decoderId)
}

//export goFileRead
func goFileRead(stream unsafe.Pointer, ptr *C.uchar, nbytes C.int) C.int {
	d := getDecoderForId(stream)

	goPtr := (*byte)(ptr)
	data := unsafe.Slice(goPtr, nbytes)
	n, err := d.stream.Read(data)
	if err != nil {
		d.setErr(fmt.Errorf("failed to read from stream, %w", err))
		return -1
	}

	return C.int(n)
}

//export goFileSeek
func goFileSeek(stream unsafe.Pointer, offset C.opus_int64, whence C.int) C.int {
	d := getDecoderForId(stream)

	var goWhence int
	switch whence {
	case C.SEEK_SET:
		goWhence = io.SeekStart
	case C.SEEK_CUR:
		goWhence = io.SeekCurrent
	case C.SEEK_END:
		goWhence = io.SeekEnd
	default:
		d.setErr(fmt.Errorf("seek failed, unrecognised whence value %d", whence))
		return -1
	}
	_, err := d.stream.Seek(int64(offset), goWhence)
	if err != nil {
		d.setErr(fmt.Errorf("failed to seek stream, %w", err))
		return -1
	}

	return 0
}

//export goFileTell
func goFileTell(stream unsafe.Pointer) C.opus_int64 {
	d := getDecoderForId(stream)
	pos, err := d.stream.Seek(0, io.SeekCurrent)
	if err != nil {
		d.setErr(fmt.Errorf("failed to seek stream for tell, %w", err))
	}
	return C.opus_int64(pos)
}

//export goFileClose
func goFileClose(_ unsafe.Pointer) C.int {
	return 0
}
