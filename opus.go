package opus

// #include <opusfile.h>
import "C"

import (
	_ "github.com/pekim/opus/c-sources"
)

func Test(data []byte) error {
	result := C.op_test(nil, (*C.uchar)(&data[0]), C.ulong(min(512, len(data))))
	if result < 0 {
		return errorFromOpusFileError(result)
	}
	return nil
}
