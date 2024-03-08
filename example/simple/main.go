package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pekim/opus"
	"github.com/pekim/opus/example"
)

func main() {
	var rsc io.ReadSeeker

	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		rsc = f
	} else if len(os.Args) == 1 {
		rsc = bytes.NewReader(example.SampleFile)
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [song.opus]\n", os.Args[0])
		os.Exit(1)
	}

	decoder, err := opus.NewDecoder(rsc)
	if err != nil {
		panic(err)
	}
	defer decoder.Destroy()

	fmt.Println("channel count :", decoder.ChannelCount())
	fmt.Println("duration (seconds) :", decoder.Duration().Seconds())
	fmt.Println("len :", decoder.Len())
	fmt.Printf("tags, vendor : %#v\n", decoder.TagsVendor())
	fmt.Println("tags, user comments :", decoder.TagsUserComments())
	fmt.Printf("head : %#v\n", decoder.Head())
	p, err := decoder.Position()
	if err != nil {
		panic(err)
	}
	fmt.Println("pos", p)

	var pcm = make([]int16, 10_000)
	samplesReadPerChannel, err := decoder.Read(pcm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d samples (per channel)\n", samplesReadPerChannel)
	p, err = decoder.Position()
	if err != nil {
		panic(err)
	}
	fmt.Println("pos", p)
	fmt.Println(pcm[:10])
	pos, _ := decoder.Position()

	samplesReadPerChannel, err = decoder.Read(pcm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d samples (per channel)\n", samplesReadPerChannel)
	p, err = decoder.Position()
	if err != nil {
		panic(err)
	}
	fmt.Println("pos", p)
	fmt.Println(pcm[:10])

	err = decoder.Seek(pos)
	if err != nil {
		panic(err)
	}

	samplesReadPerChannel, err = decoder.Read(pcm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d samples (per channel)\n", samplesReadPerChannel)
	p, err = decoder.Position()
	if err != nil {
		panic(err)
	}
	fmt.Println("pos", p)
	fmt.Println(pcm[:10])

	var pcmFloat = make([]float32, 10_000)
	samplesReadPerChannel, err = decoder.ReadFloat(pcmFloat)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d float samples (per channel)\n", samplesReadPerChannel)
	p, err = decoder.Position()
	if err != nil {
		panic(err)
	}
	fmt.Println("pos", p)
	fmt.Println(pcmFloat[:10])

	err = decoder.Err()
	if err != nil {
		panic(err)
	}
}
