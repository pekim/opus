package opus

// #include <opus.h>
// #include <opusfile.h>
import "C"

import (
	"fmt"
	"os"
)

func Test() {
	fmt.Println(C.GoString(C.opus_get_version_string()))

	data, err := os.ReadFile("/home/mike/Music/The Chicks - Travelin' Soldier (Official Video) [AbfgxznPmZM].opus")
	if err != nil {
		panic(err)
	}

	var opusFileErr C.int
	oggOpusFile := C.op_open_memory((*C.uchar)(&data[0]), C.size_t(len(data)), &opusFileErr)
	fmt.Println(opusFileErr, oggOpusFile)
	li := C.op_current_link(oggOpusFile)
	fmt.Println(li)
	channel_count := C.op_channel_count(oggOpusFile, li)
	fmt.Println(channel_count)
	pcmTotal := C.op_pcm_total(oggOpusFile, li)
	fmt.Println(pcmTotal, pcmTotal/48_000, pcmTotal/48_000/60, pcmTotal/48_000%60)
	bitrate := C.op_bitrate(oggOpusFile, li)
	fmt.Println(bitrate)
	fmt.Println()

	var pcm = make([]int16, 20_000)
	samplesReadPerChannel := C.op_read(oggOpusFile,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)), // / channel_count,
		&li,
	)
	fmt.Println(li, samplesReadPerChannel, samplesReadPerChannel*channel_count)
	fmt.Println(pcm[:100])
	i := samplesReadPerChannel * channel_count
	fmt.Println(pcm[i-20 : i+20])

	// numChannels := C.opus_packet_get_nb_channels((*C.uchar)(&data[0]))
	// fmt.Println("numChannels", numChannels)
	// numFrames := C.opus_packet_get_nb_frames((*C.uchar)(&data[0]), C.opus_int32(len(data)))
	// fmt.Println("numFrames", numFrames)

	// var error C.int
	// decoder := C.opus_decoder_create(48000, numChannels, &error)
	// fmt.Println(error == C.OPUS_OK, decoder)

	// var pcm = make([]int16, 200)

	// numDecodedSamples := C.opus_decode(decoder,
	// 	(*C.uchar)(&data[0]),
	// 	C.opus_int32(len(data)),
	// 	(*C.opus_int16)(&pcm[0]),
	// 	C.int(cap(pcm))/numChannels,
	// 	0,
	// )

	// if numDecodedSamples > 0 {
	// 	fmt.Println(numDecodedSamples)
	// } else {
	// 	fmt.Println(errorToString(numDecodedSamples))
	// }
}

func errorToString(err C.int) string {
	return C.GoString(C.opus_strerror(err))
}
