package opus

// #include <opusfile.h>
import "C"
import (
	"time"
)

type Decoder struct {
	file         *C.OggOpusFile
	data         []byte
	channelCount C.int
	duration     time.Duration
}

func NewDecoder(data []byte) (*Decoder, error) {
	var opusFileErr C.int
	file := C.op_open_memory((*C.uchar)(&data[0]), C.size_t(len(data)), &opusFileErr)
	if opusFileErr < 0 {
		return nil, errorFromOpusFileError(opusFileErr)
	}

	link := C.op_current_link(file)
	channelCount := C.op_channel_count(file, link)
	pcmTotal := C.op_pcm_total(file, link)
	duration := time.Millisecond * time.Duration((float64(pcmTotal) / 48_000 * 1_000))

	d := &Decoder{
		file:         file,
		data:         data,
		channelCount: channelCount,
		duration:     duration,
	}

	return d, nil
}

func (d *Decoder) Duration() time.Duration {
	return d.duration
}
