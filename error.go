package opus

// #include <opusfile.h>
import "C"

import (
	"fmt"
	"math"
)

// OpusFileError represents an error from the opusfile C library.
type OpusFileError struct {
	code int
	text string
}

func (err *OpusFileError) Error() string {
	return fmt.Sprintf("opusfile error %d, '%s'", err.code, err.text)
}

// Code returns the numeric error code reported by the opusfile C library.
func (err *OpusFileError) Code() int {
	return err.code
}

// Text provides text that describes the nature of the error.
func (err *OpusFileError) Text() string {
	return err.text
}

// Errors from the opusfile C library.
var (
	// A request did not succeed.
	OP_FALSE = OpusFileError{
		code: C.OP_FALSE,
		text: "A request did not succeed.",
	}

	// Currently not used externally.
	OP_EOF = OpusFileError{
		code: C.OP_EOF,
		text: "Currently not used externally.",
	}

	// There was a hole in the page sequence numbers (e.g., a page was corrupt or missing).
	OP_HOLE = OpusFileError{
		code: C.OP_HOLE,
		text: "There was a hole in the page sequence numbers (e.g., a page was corrupt or missing).",
	}

	// An underlying read, seek, or tell operation failed when it should have succeeded.
	OP_EREAD = OpusFileError{
		code: C.OP_EREAD,
		text: "An underlying read, seek, or tell operation failed when it should have succeeded.",
	}

	// A NULL pointer was passed where one was unexpected, or an internal memory allocation failed, or an internal library error was encountered.
	OP_EFAULT = OpusFileError{
		code: C.OP_EFAULT,
		text: "A NULL pointer was passed where one was unexpected, or an internal memory allocation failed, or an internal library error was encountered.",
	}

	// The stream used a feature that is not implemented, such as an unsupported channel family.
	OP_EIMPL = OpusFileError{
		code: C.OP_EIMPL,
		text: "The stream used a feature that is not implemented, such as an unsupported channel family.",
	}

	// One or more parameters to a function were invalid.
	OP_EINVAL = OpusFileError{
		code: C.OP_EINVAL,
		text: "One or more parameters to a function were invalid.",
	}

	// A purported Ogg Opus stream did not begin with an Ogg page, a purported header packet did not start with one of the required strings, "OpusHead" or "OpusTags", or a link in a chained file was encountered that did not contain any logical Opus streams.
	OP_ENOTFORMAT = OpusFileError{
		code: C.OP_ENOTFORMAT,
		text: "A purported Ogg Opus stream did not begin with an Ogg page, a purported header packet did not start with one of the required strings, \"OpusHead\" or \"OpusTags\", or a link in a chained file was encountered that did not contain any logical Opus streams.",
	}

	// A required header packet was not properly formatted, contained illegal values, or was missing altogether.
	OP_EBADHEADER = OpusFileError{
		code: C.OP_EBADHEADER,
		text: "A required header packet was not properly formatted, contained illegal values, or was missing altogether.",
	}

	// The ID header contained an unrecognized version number.
	OP_EVERSION = OpusFileError{
		code: C.OP_EVERSION,
		text: "The ID header contained an unrecognized version number.",
	}

	// Currently not used at all.
	OP_ENOTAUDIO = OpusFileError{
		code: C.OP_ENOTAUDIO,
		text: "Currently not used at all.",
	}

	// An audio packet failed to decode properly. More...
	OP_EBADPACKET = OpusFileError{
		code: C.OP_EBADPACKET,
		text: "An audio packet failed to decode properly. More...",
	}

	// We failed to find data we had seen before, or the bitstream structure was sufficiently malformed that seeking to the target destination was impossible.
	OP_EBADLINK = OpusFileError{
		code: C.OP_EBADLINK,
		text: "We failed to find data we had seen before, or the bitstream structure was sufficiently malformed that seeking to the target destination was impossible.",
	}

	// An operation that requires seeking was requested on an unseekable stream.
	OP_ENOSEEK = OpusFileError{
		code: C.OP_ENOSEEK,
		text: "An operation that requires seeking was requested on an unseekable stream.",
	}

	// The first or last granule position of a link failed basic validity checks.
	OP_EBADTIMESTAMP = OpusFileError{
		code: C.OP_EBADTIMESTAMP,
		text: "The first or last granule position of a link failed basic validity checks.",
	}
)

var allOpusFileErrors = []OpusFileError{
	OP_FALSE,
	OP_EOF,
	OP_HOLE,
	OP_EREAD,
	OP_EFAULT,
	OP_EIMPL,
	OP_EINVAL,
	OP_ENOTFORMAT,
	OP_EBADHEADER,
	OP_EVERSION,
	OP_ENOTAUDIO,
	OP_EBADPACKET,
	OP_EBADLINK,
	OP_ENOSEEK,
	OP_EBADTIMESTAMP,
}

func errorFromOpusFileError(code C.int) error {
	for _, err := range allOpusFileErrors {
		if err.code == int(code) {
			return &err
		}
	}

	return &OpusFileError{
		code: math.MinInt,
		text: "unknown opus file error",
	}
}
