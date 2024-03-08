package main

import (
	"fmt"
	"io"

	"github.com/gopxl/beep"
	"github.com/pekim/opus"
	"github.com/pkg/errors"
)

// Decode takes a ReadSeekCloser containing audio data in opus format and returns a StreamSeekCloser,
// which streams that audio. The Seek method will panic if rc is not io.Seeker.
//
// Do not close the supplied ReadSeekCloser, instead, use the Close method of the returned
// StreamSeekCloser when you want to release the resources.
func opusDecode(rsc io.ReadSeekCloser) (s beep.StreamSeekCloser, format beep.Format, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "opus")
		}
	}()

	d, err := opus.NewDecoder(rsc)
	if err != nil {
		return nil, beep.Format{}, err
	}
	if d.ChannelCount() > 2 {
		return nil, beep.Format{}, fmt.Errorf("opus: unsupported number of channels, %d", d.ChannelCount())
	}

	return &decoder{
			rsc:     rsc,
			decoder: d,
		},
		beep.Format{
			SampleRate:  beep.SampleRate(opus.SampleRate),
			NumChannels: d.ChannelCount(),
			Precision:   2,
		},
		nil
}

type decoder struct {
	rsc     io.ReadSeekCloser
	decoder *opus.Decoder
	err     error
}

func (d *decoder) Stream(samples [][2]float64) (n int, ok bool) {
	if d.err != nil {
		return 0, false
	}

	tmp := make([]float32, len(samples)*d.decoder.ChannelCount())
	n, err := d.decoder.ReadFloat(tmp)
	if err != nil {
		if err == io.EOF {
			return n, n > 0
		}
		d.err = errors.Wrap(err, "opus")
		return n, false
	}

	if d.decoder.ChannelCount() == 1 {
		for i := 0; i < len(tmp); i++ {
			samples[i][0] = float64(tmp[i])
			samples[i][1] = float64(tmp[i])
		}
	} else if d.decoder.ChannelCount() == 2 {
		for i := 0; i < len(tmp); i += 2 {
			samples[i/2][0] = float64(tmp[i])
			samples[i/2][1] = float64(tmp[i+1])
		}
	} else {
		d.err = fmt.Errorf("opus: unsupported number of channels, %d", d.decoder.ChannelCount())
		return 0, false
	}

	return n, true
}

func (d *decoder) Err() error {
	return d.err
}

func (d *decoder) Len() int {
	return int(d.decoder.Len())
}

func (d *decoder) Position() int {
	pos, err := d.decoder.Position()
	if err != nil {
		d.err = errors.Wrap(err, "opus")
	}
	return int(pos)
}

func (d *decoder) Seek(p int) error {
	if p < 0 || d.Len() < p {
		return fmt.Errorf("opus: seek position %v out of range [%v, %v]", p, 0, d.Len())
	}
	err := d.decoder.Seek(int64(p))
	if err != nil {
		return errors.Wrap(err, "opus")
	}
	return nil
}

func (d *decoder) Close() error {
	d.decoder.Destroy()
	return d.rsc.Close()
}
