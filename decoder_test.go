package opus

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	d, err := NewDecoder(bytes.NewReader(SampleStream))
	assert.NotNil(t, d)
	assert.Nil(t, err)
}

func TestHead(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))
	assert.Equal(t,
		Head{Version: 1, ChannelCount: 2, PreSkip: 0x138, InputSampleRate: 0xbb80, OutputGainDb: 0},
		d.Head(),
	)
}

func TestTagsVendor(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))
	assert.Equal(t, "Lavf57.83.100", d.TagsVendor())
}

func TestTagsUserComments(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))
	assert.Equal(t,
		[]UserComment{
			{Tag: "encoder", Value: "Lavc57.107.100 libopus"},
		},
		d.TagsUserComments(),
	)
}

func TestChannelCount(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))
	assert.Equal(t, 2, d.ChannelCount())
}

func TestLen(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))
	assert.Equal(t, 5860491, d.Len())
}

func TestDuration(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))
	assert.Equal(t, 122.093, d.Duration().Seconds())
}

func TestPositionSeek(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))

	// initial position
	pos, err := d.Position()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), pos)

	// seek
	err = d.Seek(100)
	assert.Nil(t, err)

	// position after seek
	pos, err = d.Position()
	assert.Nil(t, err)
	assert.Equal(t, int64(100), pos)
}

func TestRead(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))

	samples := 6
	pcm := make([]int16, samples)
	n, err := d.Read(pcm)
	assert.Nil(t, err)
	assert.Equal(t, len(pcm)/d.ChannelCount(), n)
	assert.Equal(t,
		[]int16([]int16{81, 13, 41, -42, -55, -99}),
		pcm[:n*d.ChannelCount()],
	)
}

func TestReadFloat(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))

	pcm := make([]float32, 6)
	n, err := d.ReadFloat(pcm)
	assert.Nil(t, err)
	assert.Equal(t, len(pcm)/d.ChannelCount(), n)
	assert.Equal(t,
		[]float32([]float32{0.002462285, 0.0003901136, 0.0012744348, -0.0013133116, -0.0017320435, -0.0029684622}),
		pcm[:n*d.ChannelCount()],
	)
}

func TestReadUntilEOF(t *testing.T) {
	d, _ := NewDecoder(bytes.NewReader(SampleStream))

	samplesBeforeEnd := 100
	err := d.Seek(int64(d.Len()) - int64(samplesBeforeEnd))
	assert.Nil(t, err)

	samplesRead := 0
	pcm := make([]int16, 30)
	for {
		n, err := d.Read(pcm)
		if err == io.EOF {
			break
		}
		samplesRead += n
	}
	assert.Equal(t, samplesBeforeEnd, samplesRead)
}
