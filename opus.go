package opus

// #include <opusfile.h>
import "C"

import (
	_ "github.com/pekim/opus/c-sources"
)

/*
Test checks to see if some data appears to be the start of an Opus stream.

For good results, you will need at least 57 bytes (for a pure Opus-only stream).
Something like 512 bytes will give more reliable results for multiplexed streams.
This function is meant to be a quick-rejection filter.
Its purpose is not to guarantee that a stream is a valid Opus stream,
but to ensure that it looks enough like Opus that it isn't going to be recognized as some other format
(except possibly an Opus stream that is also multiplexed with other codecs, such as video).
*/
func Test(data []byte) error {
	result := C.op_test(nil, (*C.uchar)(&data[0]), C.ulong(min(512, len(data))))
	if result < 0 {
		return errorFromOpusFileError(result)
	}
	return nil
}
