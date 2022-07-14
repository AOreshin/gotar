package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
)

const (
	AUDIO_BUFFER          = 16
	DEFAULT_DURATION_TICS = 600000
)

func main() {
	audio, err := rtaudio.Create(rtaudio.APIUnspecified)
	if err != nil {
		panic(err)
	}
	defer audio.Destroy()

	printAvailableApis()
	printAvailableDevices(audio)

	params := rtaudio.StreamParams{
		DeviceID:     uint(audio.DefaultOutputDevice()),
		NumChannels:  2,
		FirstChannel: 0,
	}
	options := rtaudio.StreamOptions{
		Flags: rtaudio.FlagsMinimizeLatency,
	}

	strings := []*GuitarString{}

	cb := func(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
		samples := out.Float32()
		for i := 0; i < len(samples)/2; i++ {
			samples[i*2], samples[i*2+1] = stringSamples(strings)
		}
		return 0
	}
	err = audio.Open(&params, nil, rtaudio.FormatFloat32, SAMPLING_RATE, AUDIO_BUFFER, cb, &options)
	if err != nil {
		panic(err)
	}
	defer audio.Close()
	audio.Start()
	defer audio.Stop()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	decayFactor := DECAY_FACTOR

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		frequency, ok := notes[char]
		switch key {
		case keyboard.KeyPgdn:
			decayFactor -= 0.001
		case keyboard.KeyPgup:
			decayFactor += 0.001
		case keyboard.KeyArrowUp:
			for k := range notes {
				notes[k] *= 2
			}
		case keyboard.KeyArrowDown:
			for k := range notes {
				notes[k] /= 2
			}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				strings = removeDeadStrings(strings, DEFAULT_DURATION_TICS)
				strings = append(strings, NewGuitarString(frequency, decayFactor))
			}
		}
		s := fmt.Sprintf("note %s, frequency %.3f, decay factor = %.3f tics, %d ringing strings", keysToNotes[char], frequency, decayFactor, len(strings))
		fmt.Printf("\r%s", s)
	}
}

func removeDeadStrings(strings []*GuitarString, duration int) []*GuitarString {
	ringingStrings := []*GuitarString{}
	for _, s := range strings {
		if s.Time() > duration {
			continue
		}
		ringingStrings = append(ringingStrings, s)
	}
	return ringingStrings
}

func stringSamples(strings []*GuitarString) (float32, float32) {
	var sample float32
	for _, s := range strings {
		sample += s.Sample()
		s.Tic()
	}
	if sample > 1 {
		sample = 1
	}
	if sample < -1 {
		sample = -1
	}
	return sample, sample
}

func printAvailableApis() {
	fmt.Println("RtAudio version: ", rtaudio.Version())
	for _, api := range rtaudio.CompiledAPI() {
		fmt.Println("Compiled API: ", api)
	}
}

func printAvailableDevices(audio rtaudio.RtAudio) {
	devices, err := audio.Devices()
	if err != nil {
		panic(err)
	}
	for _, d := range devices {
		fmt.Printf("Audio device: %#v\n", d)
	}
}
