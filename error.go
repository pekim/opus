package opus

// #include <opusfile.h>
import "C"

import (
	"fmt"
)

type OpusError struct {
	code   int
	reason string
}

func (err *OpusError) Error() string {
	return fmt.Sprintf("opus error %d, '%s'", err.code, err.reason)
}

func (err *OpusError) Code() int {
	return err.code
}

func (err *OpusError) Reason() string {
	return err.reason
}

func errorFromOpusError(err C.int) error {
	return &OpusError{
		code:   int(err),
		reason: C.GoString(C.opus_strerror(err)),
	}
}
