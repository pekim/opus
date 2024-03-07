package opus

// #include <opusfile.h>
import "C"

import (
	"fmt"
	"os"

	_ "github.com/pekim/opus/c-sources"
)

func Temp() {
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

	var pcm = make([]int16, 10_000)
	samplesReadPerChannel, err := decoder.Read(pcm)
	if err != nil {
		panic(err)
	}
	fmt.Println(samplesReadPerChannel)

	// var i int
	// for {
	// 	var pcm = make([]int16, 10_000)
	// 	samplesReadPerChannel, err := decoder.Read(pcm)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	i++
	// 	fmt.Println(i, samplesReadPerChannel)

	// 	if samplesReadPerChannel == 0 {
	// 		break
	// 	}
	// }

	decoder.Close()
}

func Test(data []byte) error {
	result := C.op_test(nil, (*C.uchar)(&data[0]), C.ulong(min(512, len(data))))
	if result < 0 {
		return errorFromOpusFileError(result)
	}
	return nil
}
