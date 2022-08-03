package main

import (
	"errors"
	"os"
	"time"

	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
	"github.com/youpy/go-wav"
)

func callback(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
	samples := out.Float32()
	for i := 0; i < len(samples)/2; i++ {
		s := stringsSample()
		s = limitSample(s)

		l, r := applyFx(s, s)
		l, r = adjustVolume(l, r)
		l, r = playLoops(l, r)

		recordLoop(l, r)
		recordToFile(l, r)

		samples[i*2], samples[i*2+1] = l, r
	}
	return 0
}

func stringsSample() float32 {
	var sample float32
	for _, s := range sState.ringingStrings {
		sample += s.Sample() * 0.25
		s.Tic()
	}
	return sample
}

func limitSample(sample float32) float32 {
	if sample > 1 {
		sample = 1
	}
	if sample < -1 {
		sample = -1
	}
	return sample
}

func applyFx(l, r float32) (float32, float32) {
	for _, f := range fState.activeFx {
		l, r = f.apply(l, r)
	}
	return l, r
}

func adjustVolume(l, r float32) (float32, float32) {
	v := vState.volume
	return l * v, r * v
}

func recordLoop(l, r float32) {
	if lState.recordLoop {
		lState.loop[0].Append(l)
		lState.loop[1].Append(r)
	}
}

func playLoops(l, r float32) (float32, float32) {
	if len(lState.loops) > 0 && lState.playLoop {
		baseLoopIndex := lState.loops[0][0].Tic()
		lState.loops[0][1].Tic()
		for i := 0; i < len(lState.loops); i++ {
			loop := lState.loops[i]
			l += loop[0].Get(baseLoopIndex)
			r += loop[1].Get(baseLoopIndex)
		}
	}
	return l, r
}

func recordToFile(l, r float32) {
	if rState.record && rState.writer != nil {
		s := toWavSample(l, r)
		err := rState.writer.WriteSamples(s)
		if err != nil {
			if !errors.Is(err, os.ErrClosed) {
				panic(err)
			}
		}
	}
}

func toWavSample(l, r float32) []wav.Sample {
	return []wav.Sample{
		{Values: [2]int{int(l * floatToInt), int(r * floatToInt)}},
	}
}
