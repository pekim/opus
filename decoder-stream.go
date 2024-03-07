package opus

// #include "decoder-stream.h"
import "C"

import (
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

func NewStreamDecoder(stream io.ReadSeekCloser) (*Decoder, error) {
	d := &Decoder{stream: stream}

	callbacks := C.create_file_callbacks()
	var opusFileErr C.int
	opusFile := C.op_open_callbacks(
		getDecoderId(d),
		callbacks,
		nil, 0,
		&opusFileErr,
	)
	if opusFileErr < 0 {
		return nil, errorFromOpusFileError(opusFileErr)
	}

	d.init(opusFile)
	d.callbacks = callbacks

	return d, nil
}

//export goFileRead
func goFileRead(stream unsafe.Pointer, ptr *C.uchar, nbytes C.int) C.int {
	d := getDecoderForId(stream)

	goPtr := (*byte)(ptr)
	data := unsafe.Slice(goPtr, nbytes)
	n, err := d.stream.Read(data)
	if err != nil {
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
		return -1
	}
	_, err := d.stream.Seek(int64(offset), goWhence)
	if err != nil {
		return -1
	}

	return 0
}

//export goFileTell
func goFileTell(stream unsafe.Pointer) C.opus_int64 {
	d := getDecoderForId(stream)
	pos, _ := d.stream.Seek(0, io.SeekCurrent)
	return C.opus_int64(pos)
}

//export goFileClose
func goFileClose(stream unsafe.Pointer) C.int {
	d := getDecoderForId(stream)
	err := d.stream.Close()
	if err != nil {
		return -1
	}
	return 0
}
