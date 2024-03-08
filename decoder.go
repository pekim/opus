package opus

// #include "decoder-callbacks.h"
// #include <opusfile.h>
import "C"

import (
	"io"
	"strings"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

// SampleRate is number of samples per second in streams.
// It is always 48,000 for opus streams.
const SampleRate = 48_000

// Decoder decodes an opus bitstream into PCM.
type Decoder struct {
	opusFile     *C.OggOpusFile
	callbacks    *C.OpusFileCallbacks
	stream       io.ReadSeeker
	link         C.int
	channelCount C.int
	duration     time.Duration
	err          error
}

// NewDecoder creates a new opus Decoder
// for a stream that implements io.Reader and io.Seeker.
func NewDecoder(stream io.ReadSeeker) (*Decoder, error) {
	d := &Decoder{
		stream:    stream,
		callbacks: C.create_file_callbacks(),
	}

	var opusFileErr C.int
	opusFile := C.op_open_callbacks(
		getDecoderId(d),
		d.callbacks,
		nil, 0,
		&opusFileErr,
	)
	if opusFileErr < 0 {
		return nil, errorFromOpusFileError(opusFileErr)
	}

	d.opusFile = opusFile
	d.link = C.op_current_link(opusFile)
	d.channelCount = C.op_channel_count(opusFile, d.link)
	d.duration = time.Millisecond * time.Duration((float64(d.Len()) / SampleRate * 1_000))

	return d, nil
}

// Destroy releases C heap memory used by the decoder.
// After Destroy has been called no methods on the decoder should be called.
func (d *Decoder) Destroy() {
	if d.opusFile != nil {
		C.op_free(d.opusFile)
		d.opusFile = nil
	}
	if d.callbacks != nil {
		C.free_file_callbacks(d.callbacks)
		d.callbacks = nil
	}
}

// ChannelCount returns the number of channels in the stream.
// It is typically 2, representing a stereo stream.
func (d *Decoder) ChannelCount() int {
	return int(d.channelCount)
}

// Duration returns the total duration of the stream.
// It is independent of the current position in the stream.
func (d *Decoder) Duration() time.Duration {
	return d.duration
}

// Len returns the number of samples in the stream.
//
// There will be SampleRate (48,000) samples per second, per channel.
func (d *Decoder) Len() int {
	return int(C.op_pcm_total(d.opusFile, d.link))
}

// TagsVendor returns a string identifying the software used to encode the stream.
func (d *Decoder) TagsVendor() string {
	tags := C.op_tags(d.opusFile, d.link)
	return C.GoString(tags.vendor)
}

// UserComment represents a comment from the stream's meta data.
// A comment comprises a tag name and a value.
type UserComment struct {
	// Tag is the comment's tag name.
	Tag string
	// Value is the comment value.
	Value string
}

// TagsUserComments returns comments from the stream's metadata.
//
// The comments are ordered.
// Comments may have duplicate tag names.
func (d *Decoder) TagsUserComments() []UserComment {
	tags := C.op_tags(d.opusFile, d.link)
	cUserComments := (**C.char)(tags.user_comments)
	userComments := unsafe.Slice(cUserComments, tags.comments)

	splitUserComments := make([]UserComment, tags.comments)
	for i, cComment := range userComments {
		comment := C.GoString(cComment)
		parts := strings.SplitN(comment, "=", 2)
		tag := parts[0]
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}
		splitUserComments[i] = UserComment{
			Tag:   tag,
			Value: value,
		}
	}

	return splitUserComments
}

// Head contains playback parameters for a stream.
type Head struct {
	// Version holds the Ogg Opus format version, in the range 0...255.
	//
	// The top 4 bits represent a "major" version, and the bottom four bits represent backwards-compatible "minor" revisions. The current specification describes version 1. This library will recognize versions up through 15 as backwards compatible with the current specification. An earlier draft of the specification described a version 0, but the only difference between version 1 and version 0 is that version 0 did not specify the semantics for handling the version field.
	Version int
	// ChannelCount is the number of channels, in the range 1...255.
	ChannelCount int
	// PreSkip is the number of samples that should be discarded from the beginning of the stream.
	PreSkip uint
	// InputSampleRate is sampling rate of the original input.
	//
	// All Opus audio is coded at 48 kHz, and should also be decoded at 48 kHz for playback (unless the target hardware does not support this sampling rate). However, this field may be used to resample the audio back to the original sampling rate, for example, when saving the output to a file.
	InputSampleRate uint32
	// OutputGainDb is the gain to apply to the decoded output, in dB, as a Q8 value in the range -32768...32767.
	//
	// The decoder will automatically scale the output by pow(10,output_gain/(20.0*256)).
	OutputGainDb int
}

// Head gets the stream's header information.
func (d *Decoder) Head() Head {
	head := C.op_head(d.opusFile, d.link)
	return Head{
		Version:         int(head.version),
		ChannelCount:    int(head.channel_count),
		PreSkip:         uint(head.pre_skip),
		InputSampleRate: uint32(head.input_sample_rate),
		OutputGainDb:    int(head.output_gain),
	}
}

// Read will read up to len(pcm) samples in to the pcm argument.
//
// The number of samples actually read, per channel, is returned.
// When there are no more samples to read, 0 and an io.EOF error will be returned.
func (d *Decoder) Read(pcm []int16) (int, error) {
	samplesReadPerChannel := C.op_read(
		d.opusFile,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)),
		nil,
	)

	if samplesReadPerChannel < 0 {
		return int(samplesReadPerChannel), d.errorFromOpusFileError(samplesReadPerChannel)
	}
	if samplesReadPerChannel == 0 {
		return 0, io.EOF
	}

	return int(samplesReadPerChannel), nil
}

// Read will read up to len(pcm) samples in to the pcm argument.
//
// The number of samples actually read, per channel, is returned.
// When there are no more samples to read, 0 and an io.EOF error will be returned.
func (d *Decoder) ReadFloat(pcm []float32) (int, error) {
	samplesReadPerChannel := C.op_read_float(
		d.opusFile,
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)),
		nil,
	)

	if samplesReadPerChannel < 0 {
		return int(samplesReadPerChannel), d.errorFromOpusFileError(samplesReadPerChannel)
	}
	if samplesReadPerChannel == 0 {
		return 0, io.EOF
	}

	return int(samplesReadPerChannel), nil
}

// Err returns an error which may have occurred during streaming.
// If no error has ever occurred, nil is returned.
func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) setErr(err error) {
	if d.err == nil {
		d.err = err
	}
}

// Seek will seek to the specified PCM offset,
// such that decoding (using the Read or ReadFloat method)
// will begin at exactly the requested position.
func (d *Decoder) Seek(pos int64) error {
	err := C.op_pcm_seek(d.opusFile, C.ogg_int64_t(pos))
	if err < 0 {
		return d.errorFromOpusFileError(err)
	}
	return nil
}

// Position obtains the PCM offset of the next sample to be read.
//
// If the stream is not properly timestamped,
// this might not increment by the proper amount between reads,
// or even return monotonically increasing values.
func (d *Decoder) Position() (int64, error) {
	pos := C.op_pcm_tell(d.opusFile)
	if pos < 0 {
		return 0, d.errorFromOpusFileError(C.int(pos))
	}
	return int64(pos), nil
}

func (d *Decoder) errorFromOpusFileError(code C.int) error {
	if d.err != nil {
		return errors.Wrap(d.err, errorFromOpusFileError(code).Error())
	} else {
		return errorFromOpusFileError(code)
	}
}
