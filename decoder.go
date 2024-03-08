package opus

// #include "decoder-stream.h"
import "C"

import (
	"io"
	"strings"
	"time"
	"unsafe"
)

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

// NewDecoder creates a new opus Decoder.
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
		err := errorFromOpusFileError(opusFileErr)
		d.setErr(err)
		return nil, err
	}

	d.opusFile = opusFile
	d.link = C.op_current_link(opusFile)
	d.channelCount = C.op_channel_count(opusFile, d.link)
	d.duration = time.Millisecond * time.Duration((float64(d.Len()) / 48_000 * 1_000))

	return d, nil
}

func (d *Decoder) Destroy() {
	C.op_free(d.opusFile)
	if d.stream != nil {
		d.stream = nil
	}
	if d.callbacks != nil {
		C.free_file_callbacks(d.callbacks)
	}
}

func (d *Decoder) ChannelCount() int {
	return int(d.channelCount)
}

func (d *Decoder) Duration() time.Duration {
	return d.duration
}

func (d *Decoder) Len() int {
	return int(C.op_pcm_total(d.opusFile, d.link))
}

func (d *Decoder) TagsVendor() string {
	tags := C.op_tags(d.opusFile, d.link)
	return C.GoString(tags.vendor)
}

type UserComment struct {
	Tag   string
	Value string
}

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

type Head struct {
	Version         int
	ChannelCount    int
	PreSkip         uint
	InputSampleRate uint32
	OutputGainDb    int
}

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

func (d *Decoder) Read(pcm []int16) (int, error) {
	samplesReadPerChannel := C.op_read(
		d.opusFile,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)),
		nil,
	)

	if samplesReadPerChannel < 0 {
		err := errorFromOpusFileError(samplesReadPerChannel)
		d.setErr(err)
		return int(samplesReadPerChannel), err
	}

	return int(samplesReadPerChannel), nil
}

func (d *Decoder) ReadFloat(pcm []float32) (int, error) {
	samplesReadPerChannel := C.op_read_float(
		d.opusFile,
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)),
		nil,
	)

	if samplesReadPerChannel < 0 {
		err := errorFromOpusFileError(samplesReadPerChannel)
		d.setErr(err)
		return int(samplesReadPerChannel), err
	}

	return int(samplesReadPerChannel), nil
}

// Err returns an error which occurred during streaming. If no error occurred, nil is
// returned.
func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) setErr(err error) {
	if d.err == nil {
		d.err = err
	}
}

func (d *Decoder) Seek(pos int64) error {
	err := C.op_pcm_seek(d.opusFile, C.ogg_int64_t(pos))
	if err < 0 {
		err := errorFromOpusFileError(err)
		d.setErr(err)
		return err
	}
	return nil
}

func (d *Decoder) Position() (int64, error) {
	pos := C.op_pcm_tell(d.opusFile)
	if pos < 0 {
		err := errorFromOpusFileError(C.int(pos))
		d.setErr(err)
		return 0, err
	}
	return int64(pos), nil
}
