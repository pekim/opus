package main

import (
	"fmt"
	"os"

	"github.com/pekim/opus"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expected 1 argument, path to an opus file")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = opus.Test(data)
	if err != nil {
		panic(err)
	}

	decoder, err := opus.NewDecoder(data)
	if err != nil {
		panic(err)
	}

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

	decoder.Close()

	err = decoder.Err()
	if err != nil {
		panic(err)
	}
}
