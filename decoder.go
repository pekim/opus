package opus

// #include <opusfile.h>
import "C"

import (
	"strings"
	"time"
	"unsafe"
)

type Decoder struct {
	file         *C.OggOpusFile
	data         []byte
	link         C.int
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
		link:         link,
		channelCount: channelCount,
		duration:     duration,
	}

	return d, nil
}

func (d *Decoder) Close() {
	C.op_free(d.file)
	d.data = nil
}

func (d *Decoder) ChannelCount() int {
	return int(d.channelCount)
}

func (d *Decoder) Duration() time.Duration {
	return d.duration
}

func (d *Decoder) TagsVendor() string {
	tags := C.op_tags(d.file, d.link)
	return C.GoString(tags.vendor)
}

type UserComment struct {
	Tag   string
	Value string
}

func (d *Decoder) TagsUserComments() []UserComment {
	tags := C.op_tags(d.file, d.link)
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

func (d *Decoder) Read(pcm []int16) (int, error) {
	samplesReadPerChannel := C.op_read(
		d.file,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)),
		nil,
	)

	if samplesReadPerChannel < 0 {
		return int(samplesReadPerChannel), errorFromOpusFileError(samplesReadPerChannel)
	}

	return int(samplesReadPerChannel), nil
}
