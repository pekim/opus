package opus

// #include <opusfile.h>
import "C"

import (
	"fmt"
	"os"

	_ "github.com/pekim/opus/c-sources"
)

func Temp() {
	// fmt.Println(C.GoString(C.opus_get_version_string()))

	// panic(errorFromOpusFileError(-133))

	data, err := os.ReadFile("/home/mike/Music/The Chicks - Travelin' Soldier (Official Video) [AbfgxznPmZM].opus")
	if err != nil {
		panic(err)
	}

	err = Test(data)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoder.channelCount, decoder.duration.Seconds())
	fmt.Println(decoder.TagsVendor())
	fmt.Println(decoder.TagsUserComments())

	var opusFileErr C.int
	oggOpusFile := C.op_open_memory((*C.uchar)(&data[0]), C.size_t(len(data)), &opusFileErr)
	if opusFileErr < 0 {
		panic(errorFromOpusFileError(opusFileErr))
	}

	li := C.op_current_link(oggOpusFile)
	channel_count := C.op_channel_count(oggOpusFile, li)
	pcmTotal := C.op_pcm_total(oggOpusFile, li)
	fmt.Println(pcmTotal, pcmTotal/48_000, pcmTotal/48_000/60, pcmTotal/48_000%60)
	// bitrate := C.op_bitrate(oggOpusFile, li)
	// fmt.Println("bitrate", bitrate)

	fmt.Println()

	var pcm = make([]int16, 20_000)
	samplesReadPerChannel := C.op_read(oggOpusFile,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)), // / channel_count,
		&li,
	)
	if samplesReadPerChannel < 0 {
		panic(errorFromOpusFileError(samplesReadPerChannel))
	}

	fmt.Println(li, samplesReadPerChannel, samplesReadPerChannel*channel_count)
	fmt.Println(pcm[:100])
	i := samplesReadPerChannel * channel_count
	fmt.Println(pcm[i-20 : i+20])

	decoder.Close()
}

func Test(data []byte) error {
	result := C.op_test(nil, (*C.uchar)(&data[0]), C.ulong(min(512, len(data))))
	if result < 0 {
		return errorFromOpusFileError(result)
	}
	return nil
}
